package awslambdago

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

func makeSureSQSTransportHasCorrectSignature() {
	h := NewSQSTransport(
		true,
		func(events.SQSMessage) (string, error) { return "", nil },
		func(context.Context, string) (string, error) { return "", nil },
	)
	lambda.StartHandlerFunc(h)
}

func Test_SQS_Handle(t *testing.T) {
	testCases := []struct {
		desc       string
		fifo       bool
		decodeFunc func(events.SQSMessage) (string, error)
		ep         func(context.Context, string) (string, error)
		input      *events.SQSEvent
		want       *events.SQSEventResponse
	}{
		{
			desc: "test successfully processing a single event",
			fifo: false,
			decodeFunc: func(events.SQSMessage) (string, error) {
				return "", nil
			},
			ep: func(context.Context, string) (string, error) {
				return "", nil
			},
			input: &events.SQSEvent{
				Records: []events.SQSMessage{{}},
			},
			want: &events.SQSEventResponse{},
		},
		{
			desc: "test failing to decode doesn't call endpoint",
			fifo: false,
			decodeFunc: func(ev events.SQSMessage) (string, error) {
				return "", errors.New("boom")
			},
			ep: func(context.Context, string) (string, error) {
				panic("don't call this")
			},
			input: &events.SQSEvent{
				Records: []events.SQSMessage{
					{MessageId: "first", Body: "first"},
				},
			},
			want: &events.SQSEventResponse{BatchItemFailures: []events.SQSBatchItemFailure{
				{ItemIdentifier: "first"},
			}},
		},
		{
			desc: "test failing a single event",
			fifo: false,
			decodeFunc: func(ev events.SQSMessage) (string, error) {
				return ev.Body, nil
			},
			ep: func(_ context.Context, input string) (string, error) {
				if input == "first" {
					return "", errors.New("boom")
				}
				return "", nil
			},
			input: &events.SQSEvent{
				Records: []events.SQSMessage{
					{MessageId: "first", Body: "first"},
					{MessageId: "second", Body: "second"},
				},
			},
			want: &events.SQSEventResponse{BatchItemFailures: []events.SQSBatchItemFailure{
				{ItemIdentifier: "first"},
			}},
		},
		{
			desc: "test failing a all events in fifo queue",
			fifo: true,
			decodeFunc: func(ev events.SQSMessage) (string, error) {
				return ev.Body, nil
			},
			ep: func(_ context.Context, input string) (string, error) {
				if input == "first" {
					return "", errors.New("boom")
				}
				return "", nil
			},
			input: &events.SQSEvent{
				Records: []events.SQSMessage{
					{MessageId: "first", Body: "first"},
					{MessageId: "second", Body: "second"},
				},
			},
			want: &events.SQSEventResponse{BatchItemFailures: []events.SQSBatchItemFailure{
				{ItemIdentifier: "first"},
				{ItemIdentifier: "second"},
			}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			h := NewSQSTransport(
				tC.fifo,
				tC.decodeFunc,
				tC.ep,
			)
			got, err := h(context.Background(), tC.input)

			assert.Equal(t, tC.want, got)
			assert.NoError(t, err)
		})
	}
}

func Test_SQS_TracingMiddleware(t *testing.T) {
	// given
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	h := NewSQSTransport(
		true,
		func(events.SQSMessage) (string, error) { return "", nil },
		func(context.Context, string) (string, error) { return "", nil },
	)
	mw := NewSQSTracingMiddleware(h, tp)

	// when
	got, err := mw(context.Background(), &events.SQSEvent{})

	// then
	assert.NoError(t, err)
	want := &events.SQSEventResponse{}
	assert.Equal(t, want, got)

	spans := sr.Ended()
	wantSpanAttributes := []attribute.KeyValue{
		semconv.FaaSTriggerPubsub,
		semconv.MessagingOperationProcess,
		semconv.MessagingSystem("AmazonSQS"),
		semconv.MessagingSourceKindQueue,
	}
	assert.Len(t, spans, 1)
	assert.ElementsMatch(t, wantSpanAttributes, spans[0].Attributes())
	assert.Equal(t, "multiple_sources process", spans[0].Name())
}

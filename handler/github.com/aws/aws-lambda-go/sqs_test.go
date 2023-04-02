package awslambdago

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/stretchr/testify/assert"
)

func makeSureSQSHandlerHasCorrectSignature() {
	h := NewSQSHandler(
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
			h := NewSQSHandler(
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

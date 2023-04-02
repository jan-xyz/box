package awslambdago

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/stretchr/testify/assert"
)

func makeSureKinesisHandlerHasCorrectSignature() {
	h := NewKinesisHandler(
		func(events.KinesisEventRecord) (string, error) { return "", nil },
		func(context.Context, string) (string, error) { return "", nil },
	)
	lambda.StartHandlerFunc(h)
}

func Test_Kinesis_Handle(t *testing.T) {
	testCases := []struct {
		desc       string
		decodeFunc func(events.KinesisEventRecord) (string, error)
		ep         func(context.Context, string) (string, error)
		input      *events.KinesisEvent
		want       *events.KinesisEventResponse
	}{
		{
			desc: "test successfully processing a single event",
			decodeFunc: func(events.KinesisEventRecord) (string, error) {
				return "", nil
			},
			ep: func(context.Context, string) (string, error) {
				return "", nil
			},
			input: &events.KinesisEvent{
				Records: []events.KinesisEventRecord{{}},
			},
			want: &events.KinesisEventResponse{},
		},
		{
			desc: "test failing to decode doesn't call endpoint",
			decodeFunc: func(ev events.KinesisEventRecord) (string, error) {
				return "", errors.New("boom")
			},
			ep: func(context.Context, string) (string, error) {
				panic("don't call this")
			},
			input: &events.KinesisEvent{
				Records: []events.KinesisEventRecord{
					{EventID: "first", Kinesis: events.KinesisRecord{Data: []byte("first")}},
				},
			},
			want: &events.KinesisEventResponse{BatchItemFailures: []events.KinesisBatchItemFailure{
				{ItemIdentifier: "first"},
			}},
		},
		{
			desc: "test failing a single event",
			decodeFunc: func(ev events.KinesisEventRecord) (string, error) {
				return string(ev.Kinesis.Data), nil
			},
			ep: func(_ context.Context, input string) (string, error) {
				if input == "first" {
					return "", errors.New("boom")
				}
				return "", nil
			},
			input: &events.KinesisEvent{
				Records: []events.KinesisEventRecord{
					{EventID: "first", Kinesis: events.KinesisRecord{Data: []byte("first")}},
				},
			},
			want: &events.KinesisEventResponse{BatchItemFailures: []events.KinesisBatchItemFailure{
				{ItemIdentifier: "first"},
			}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			h := NewKinesisHandler(
				tC.decodeFunc,
				tC.ep,
			)
			got, err := h(context.Background(), tC.input)

			assert.Equal(t, tC.want, got)
			assert.NoError(t, err)
		})
	}
}

package awslambdago

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/stretchr/testify/assert"
)

func makeSureCloudwatchTransportHasCorrectSignature() {
	h := NewClouadWatchEventTransport(
		func(*events.CloudWatchEvent) (string, error) { return "", nil },
		func(context.Context, string) (string, error) { return "", nil },
	)
	lambda.StartHandlerFunc(h)
}

func Test_CloudWatchEvent_Handle(t *testing.T) {
	testCases := []struct {
		desc       string
		decodeFunc func(*events.CloudWatchEvent) (string, error)
		ep         func(context.Context, string) (string, error)
		input      *events.CloudWatchEvent
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			desc: "test successfully processing a single event",
			decodeFunc: func(*events.CloudWatchEvent) (string, error) {
				return "", nil
			},
			ep: func(context.Context, string) (string, error) {
				return "", nil
			},
			input:   &events.CloudWatchEvent{},
			wantErr: assert.NoError,
		},
		{
			desc: "test failing to decode doesn't call endpoint",
			decodeFunc: func(ev *events.CloudWatchEvent) (string, error) {
				return "", errors.New("boom")
			},
			ep: func(context.Context, string) (string, error) {
				panic("don't call this")
			},
			input: &events.CloudWatchEvent{
				ID:     "first",
				Detail: []byte("first"),
			},
			wantErr: assert.Error,
		},
		{
			desc: "test failing when the endpoint fails",
			decodeFunc: func(ev *events.CloudWatchEvent) (string, error) {
				return string(ev.Detail), nil
			},
			ep: func(_ context.Context, input string) (string, error) {
				return "", errors.New("boom")
			},
			input: &events.CloudWatchEvent{
				ID:     "first",
				Detail: []byte("first"),
			},
			wantErr: assert.Error,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			h := NewClouadWatchEventTransport(
				tC.decodeFunc,
				tC.ep,
			)
			_, err := h(context.Background(), tC.input)

			tC.wantErr(t, err)
		})
	}
}

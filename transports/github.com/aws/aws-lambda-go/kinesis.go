package awslambdago

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jan-xyz/box"
)

type KinesisTransport = func(ctx context.Context, e *events.KinesisEvent) (*events.KinesisEventResponse, error)

func NewKinesisTransport[TIn, TOut any](
	decode func(events.KinesisEventRecord) (TIn, error),
	endpoint box.Endpoint[TIn, TOut],
) KinesisTransport {
	return func(ctx context.Context, e *events.KinesisEvent) (*events.KinesisEventResponse, error) {
		resp := &events.KinesisEventResponse{}
		for _, r := range e.Records {
			in, err := decode(r)
			if err != nil {
				resp.BatchItemFailures = append(resp.BatchItemFailures, events.KinesisBatchItemFailure{ItemIdentifier: r.EventID})
				continue
			}
			_, err = endpoint(ctx, in)
			if err != nil {
				resp.BatchItemFailures = append(resp.BatchItemFailures, events.KinesisBatchItemFailure{ItemIdentifier: r.EventID})
				continue
			}
		}
		return resp, nil
	}
}

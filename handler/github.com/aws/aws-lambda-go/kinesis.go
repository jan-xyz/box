package awslambdago

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jan-xyz/box"
)

func NewKinesisHandler[TIn, TOut any](
	decode func(events.KinesisEventRecord) (TIn, error),
	endpoint box.Endpoint[TIn, TOut],
) KinesisHandler {
	return kinesisHandler[TIn, TOut]{
		decode:   decode,
		endpoint: endpoint,
	}
}

type KinesisHandler interface {
	Handle(ctx context.Context, e *events.KinesisEvent) *events.KinesisEventResponse
}

type kinesisHandler[TIn, TOut any] struct {
	decode   func(events.KinesisEventRecord) (TIn, error)
	endpoint box.Endpoint[TIn, TOut]
}

func (s kinesisHandler[TIn, TOut]) Handle(ctx context.Context, e *events.KinesisEvent) *events.KinesisEventResponse {
	resp := &events.KinesisEventResponse{}
	for _, r := range e.Records {
		in, err := s.decode(r)
		if err != nil {
			resp.BatchItemFailures = append(resp.BatchItemFailures, events.KinesisBatchItemFailure{ItemIdentifier: r.EventID})
			continue
		}
		_, err = s.endpoint(ctx, in)
		if err != nil {
			resp.BatchItemFailures = append(resp.BatchItemFailures, events.KinesisBatchItemFailure{ItemIdentifier: r.EventID})
			continue
		}
	}
	return resp
}

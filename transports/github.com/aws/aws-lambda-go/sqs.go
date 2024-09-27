package awslambdago

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jan-xyz/box"
)

type SQSTransport = func(ctx context.Context, e *events.SQSEvent) (*events.SQSEventResponse, error)

func NewSQSTransport[TIn any](
	fifo bool,
	decode func(events.SQSMessage) (TIn, error),
	endpoint box.Endpoint[TIn, any],
) SQSTransport {
	return func(ctx context.Context, e *events.SQSEvent) (*events.SQSEventResponse, error) {
		resp := &events.SQSEventResponse{}
		for _, r := range e.Records {
			// in a FIFO queue all other items in the batch after a failure should also be failed
			// to preserve ordering.
			if fifo && len(resp.BatchItemFailures) > 0 {
				resp.BatchItemFailures = append(resp.BatchItemFailures, events.SQSBatchItemFailure{ItemIdentifier: r.MessageId})
				continue
			}
			in, err := decode(r)
			if err != nil {
				resp.BatchItemFailures = append(resp.BatchItemFailures, events.SQSBatchItemFailure{ItemIdentifier: r.MessageId})
				continue
			}
			_, err = endpoint(ctx, in)
			if err != nil {
				resp.BatchItemFailures = append(resp.BatchItemFailures, events.SQSBatchItemFailure{ItemIdentifier: r.MessageId})
				continue
			}
		}
		return resp, nil
	}
}

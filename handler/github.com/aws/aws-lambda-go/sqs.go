package awslambdago

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)
func NewSQSHandler[TIn, TOut any](
	fifo bool,
	decode func(events.SQSMessage) (TIn, error),
	encode func(TOut) error,
	endpoint func(context.Context, TIn) (TOut, error),
) SQSHandler {
	return sqsHandler[TIn, TOut]{
		fifo:     fifo,
		decode:   decode,
		encode:   encode,
		endpoint: endpoint,
	}
}

type SQSHandler interface {
	Handle(ctx context.Context, e *events.SQSEvent) *events.SQSEventResponse
}

type sqsHandler[TIn, TOut any] struct {
	// specify if the events are coming from a fifo queue
	// and the order should be preserved on record failures
	fifo     bool
	decode   func(events.SQSMessage) (TIn, error)
	encode   func(TOut) error
	endpoint func(context.Context, TIn) (TOut, error)
}

func (s sqsHandler[TIn, TOut]) Handle(ctx context.Context, e *events.SQSEvent) *events.SQSEventResponse {
	resp := &events.SQSEventResponse{}
	for _, r := range e.Records {
		// in a FIFO queue all other items in the batch after a failure should also be failed
		// to preserve ordering.
		if s.fifo && len(resp.BatchItemFailures) > 0 {
			resp.BatchItemFailures = append(resp.BatchItemFailures, events.SQSBatchItemFailure{ItemIdentifier: r.MessageId})
			continue
		}
		in, err := s.decode(r)
		if err != nil {
			resp.BatchItemFailures = append(resp.BatchItemFailures, events.SQSBatchItemFailure{ItemIdentifier: r.MessageId})
		}
		out, err := s.endpoint(ctx, in)
		if err != nil {
			resp.BatchItemFailures = append(resp.BatchItemFailures, events.SQSBatchItemFailure{ItemIdentifier: r.MessageId})
		}
		err = s.encode(out)
		if err != nil {
			resp.BatchItemFailures = append(resp.BatchItemFailures, events.SQSBatchItemFailure{ItemIdentifier: r.MessageId})
		}
	}
	return resp
}
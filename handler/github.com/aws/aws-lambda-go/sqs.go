package awslambdago

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jan-xyz/box"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

func NewSQSHandler[TIn, TOut any](
	fifo bool,
	decode func(events.SQSMessage) (TIn, error),
	endpoint box.Endpoint[TIn, TOut],
) SQSHandler {
	return sqsHandler[TIn, TOut]{
		fifo:     fifo,
		decode:   decode,
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
	endpoint box.Endpoint[TIn, TOut]
}

func (s sqsHandler[Tin, TOut]) Handle(ctx context.Context, e *events.SQSEvent) *events.SQSEventResponse {
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
			continue
		}
		_, err = s.endpoint(ctx, in)
		if err != nil {
			resp.BatchItemFailures = append(resp.BatchItemFailures, events.SQSBatchItemFailure{ItemIdentifier: r.MessageId})
			continue
		}
	}
	return resp
}

func NewSQSTracingMiddleware(handler SQSHandler, tracer trace.Tracer) SQSHandler {
	return sqsTracingMiddleware{handler: handler, tracer: tracer}
}

type sqsTracingMiddleware struct {
	handler SQSHandler
	tracer  trace.Tracer
}

// implementation of https://opentelemetry.io/docs/reference/specification/trace/semantic_conventions/instrumentation/aws-lambda/#sqs
func (s sqsTracingMiddleware) Handle(ctx context.Context, e *events.SQSEvent) *events.SQSEventResponse {
	ctx, span := s.tracer.Start(ctx, "multiple_sources process", trace.WithAttributes(
		semconv.FaaSTriggerPubsub,
		semconv.MessagingOperationProcess,
		semconv.MessagingSystem("AmazonSQS"),
		semconv.MessagingSourceKindQueue,
	))
	defer span.End()
	resp := s.handler.Handle(ctx, e)
	return resp
}

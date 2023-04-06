package awslambdago

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jan-xyz/box"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"
)

type SQSHandler = func(ctx context.Context, e *events.SQSEvent) (*events.SQSEventResponse, error)

func NewSQSHandler[TIn, TOut any](
	fifo bool,
	decode func(events.SQSMessage) (TIn, error),
	endpoint box.Endpoint[TIn, TOut],
) SQSHandler {
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

// implementation of https://opentelemetry.io/docs/reference/specification/trace/semantic_conventions/instrumentation/aws-lambda/#sqs
func NewSQSTracingMiddleware(handler SQSHandler, tp trace.TracerProvider) SQSHandler {
	tracer := tp.Tracer(box.TracerName)
	return func(ctx context.Context, e *events.SQSEvent) (*events.SQSEventResponse, error) {
		ctx, span := tracer.Start(ctx, "multiple_sources process", trace.WithAttributes(
			semconv.FaaSTriggerPubsub,
			semconv.MessagingOperationProcess,
			semconv.MessagingSystem("AmazonSQS"),
			semconv.MessagingSourceKindQueue,
		))
		defer span.End()
		return handler(ctx, e)
	}
}

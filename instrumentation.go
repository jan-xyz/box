package box

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func EndpointTracing[TIn, TOut any](tp trace.TracerProvider) Middleware[TIn, TOut] {
	tracer := tp.Tracer("box")
	return func(next Endpoint[TIn, TOut]) Endpoint[TIn, TOut] {
		return func(ctx context.Context, req TIn) (TOut, error) {
			ctx, span := tracer.Start(ctx, "Endpoint")
			defer span.End()

			// first need to convert to `any` to use the type assertion
			switch val := any(req).(type) {
			case fmt.Stringer, string:
				span.SetAttributes(attribute.String("req", fmt.Sprintf("%s", val)))
			default:
				span.SetAttributes(attribute.String("req", fmt.Sprintf("%v", val)))
			}

			resp, err := next(ctx, req)
			span.RecordError(err)
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
			}
			return resp, err
		}
	}
}

func EndpointLogging[TIn, TOut any]() Middleware[TIn, TOut] {
	return func(next Endpoint[TIn, TOut]) Endpoint[TIn, TOut] {
		return func(ctx context.Context, req TIn) (TOut, error) {
			// first need to convert to `any` to use the type assertion
			switch val := any(req).(type) {
			case fmt.Stringer, string:
				log.Printf("incoming: %s", val)
			default:
				log.Printf("incoming: %v", val)
			}

			resp, err := next(ctx, req)

			// first need to convert to `any` to use the type assertion
			switch val := any(resp).(type) {
			case fmt.Stringer, string:
				log.Printf("incoming: %s", val)
			default:
				log.Printf("outgoing: %v", val)
			}
			return resp, err
		}
	}
}

package box

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var TracerName = "github.com/jan-xyz/box"

// EndpointTracing adds tracing for the endpoint. The Span name defaults to "Endpoint".
func EndpointTracing[TIn, TOut any](tp trace.TracerProvider) Middleware[TIn, TOut] {
	tracer := tp.Tracer(TracerName)
	spanName := "Endpoint"
	return func(next Endpoint[TIn, TOut]) Endpoint[TIn, TOut] {
		return func(ctx context.Context, req TIn) (TOut, error) {
			ctx, span := tracer.Start(ctx, spanName)
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

// EndpointLogging logs the requests and response and errors.
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

// EndpointMetrics records RED metrics for your Endpoint. It adds an Endpoint dimension to it.
// the default dimension value is "Endpoint".
func EndpointMetrics[TIn, TOut any](mp metric.MeterProvider) Middleware[TIn, TOut] {
	meter := mp.Meter(TracerName)
	endpointName := "Endpoint"
	attrs := attribute.String("endpoint", endpointName)
	return func(next Endpoint[TIn, TOut]) Endpoint[TIn, TOut] {
		return func(ctx context.Context, req TIn) (TOut, error) {
			start := time.Now()
			if counter, Merr := meter.Int64Counter("requests"); Merr == nil {
				counter.Add(ctx, 1, attrs)
			}
			if hist, Merr := meter.Float64Histogram("latency"); Merr == nil {
				defer hist.Record(ctx, time.Since(start).Seconds(), attrs)
			}
			resp, err := next(ctx, req)
			if err != nil {
				if counter, Merr := meter.Int64Counter("errors"); Merr == nil {
					counter.Add(ctx, 1, attrs)
				}
			}
			return resp, err
		}
	}
}

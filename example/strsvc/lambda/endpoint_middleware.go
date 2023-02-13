package lambda

import (
	"context"
	"log"

	"github.com/jan-xyz/box"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var EndpointLogging = func() box.Middleware[*request, *response] {
	return func(next box.Endpoint[*request, *response]) box.Endpoint[*request, *response] {
		return box.EndpointFunc[*request, *response](func(ctx context.Context, req *request) (*response, error) {
			log.Printf("incoming: %s", req)

			resp, err := next.EP(ctx, req)
			log.Printf("outgoing: %s", resp)
			return resp, err
		})
	}
}

var EndpointTracing = func(tracer trace.Tracer) box.Middleware[*request, *response] {
	return func(next box.Endpoint[*request, *response]) box.Endpoint[*request, *response] {
		return box.EndpointFunc[*request, *response](func(ctx context.Context, req *request) (*response, error) {
			ctx, span := tracer.Start(ctx, "Handle", trace.WithAttributes(attribute.String("req", req.String())))
			defer span.End()

			return next.EP(ctx, req)
		})
	}
}

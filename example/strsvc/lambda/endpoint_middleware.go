package lambda

import (
	"context"
	"log"

	"github.com/jan-xyz/box"
	"github.com/jan-xyz/box/example/strsvc/lambda/proto/strsvcv1"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var EndpointLogging = func() box.Middleware[*strsvcv1.Request, *strsvcv1.Response] {
	return func(next box.Endpoint[*strsvcv1.Request, *strsvcv1.Response]) box.Endpoint[*strsvcv1.Request, *strsvcv1.Response] {
		return func(ctx context.Context, req *strsvcv1.Request) (*strsvcv1.Response, error) {
			log.Printf("incoming: %s", req)

			resp, err := next(ctx, req)
			log.Printf("outgoing: %s", resp)
			return resp, err
		}
	}
}

var EndpointTracing = func(tracer trace.Tracer) box.Middleware[*strsvcv1.Request, *strsvcv1.Response] {
	return func(next box.Endpoint[*strsvcv1.Request, *strsvcv1.Response]) box.Endpoint[*strsvcv1.Request, *strsvcv1.Response] {
		return func(ctx context.Context, req *strsvcv1.Request) (*strsvcv1.Response, error) {
			ctx, span := tracer.Start(ctx, "Handle", trace.WithAttributes(attribute.String("req", req.String())))
			defer span.End()

			return next(ctx, req)
		}
	}
}

package strsvc

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tracingMiddleware struct {
svc Service
tracer trace.Tracer
}

func (t tracingMiddleware) UpperCase(ctx context.Context, s string) (string, error){
	ctx, span := t.tracer.Start(ctx, "UpperCase")
	defer span.End()
	span.SetAttributes(attribute.String("input", s))

	return t.svc.UpperCase(ctx, s)

}

func (t tracingMiddleware) LowerCase(ctx context.Context, s string) (string, error){
	ctx, span := t.tracer.Start(ctx, "LowerCase")
	defer span.End()
	span.SetAttributes(attribute.String("input", s))

	return t.svc.LowerCase(ctx, s)

}

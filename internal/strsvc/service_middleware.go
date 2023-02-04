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
	t.tracer.Start(ctx, "UpperCase", attribute.String("input", s))

	return t.svc.UpperCase(ctx, s)

}

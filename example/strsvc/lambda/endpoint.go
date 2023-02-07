package lambda

import (
	"context"
	"errors"
	"log"

	"github.com/jan-xyz/box/example/strsvc/lambda/proto/strsvcv1"
	"go.opentelemetry.io/otel"
)

type StringEndpoint interface {
	EP(ctx context.Context, req *strsvcv1.Request) (*strsvcv1.Response, error)
}

func NewEndpoint() StringEndpoint {
	var svc Service = &service{}
	svc = tracingMiddleware{svc: svc, tracer: otel.Tracer("strsvc")}
	svc = validationMiddleware{svc: svc}
	return endpoint{svc: svc}
}

type endpoint struct {
	svc Service
}

var errUknown = errors.New("unknown message")

func (e endpoint) EP(ctx context.Context, req *strsvcv1.Request) (*strsvcv1.Response, error) {
	switch m := req.GetMessage().(type) {
	case *strsvcv1.Request_LowerCase:
		upper, err := e.svc.LowerCase(ctx, m.LowerCase.GetInput())
		if err != nil {
			return &strsvcv1.Response{}, err
		}
		return &strsvcv1.Response{Result: upper}, nil
	case *strsvcv1.Request_UpperCase:
		upper, err := e.svc.UpperCase(ctx, m.UpperCase.GetInput())
		if err != nil {
			return &strsvcv1.Response{}, err
		}
		return &strsvcv1.Response{Result: upper}, nil
	default:
		log.Printf("unhandled message: %T", m)
	}
	return nil, errUknown
}

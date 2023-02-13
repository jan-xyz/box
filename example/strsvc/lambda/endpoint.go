package lambda

import (
	"context"
	"errors"
	"log"

	"github.com/jan-xyz/box/example/strsvc/lambda/proto/strsvcv1"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("strsvc")

func NewEndpoint() *Endpoint {
	var svc Service = &service{}
	svc = tracingMiddleware{svc: svc, tracer: tracer}
	svc = validationMiddleware{svc: svc}
	return &Endpoint{svc: svc}
}

type Endpoint struct {
	svc Service
	strsvcv1.StringServiceServer
}

var errUknown = errors.New("unknown message")

func (e Endpoint) EP(ctx context.Context, req *request) (*response, error) {
	switch m := req.GetMessage().(type) {
	case *strsvcv1.CasingRequest_LowerCase:
		upper, err := e.svc.LowerCase(ctx, m.LowerCase.GetInput())
		if err != nil {
			return &response{}, err
		}
		return &response{Result: upper}, nil
	case *strsvcv1.CasingRequest_UpperCase:
		upper, err := e.svc.UpperCase(ctx, m.UpperCase.GetInput())
		if err != nil {
			return &response{}, err
		}
		return &response{Result: upper}, nil
	default:
		log.Printf("unhandled request: %T", m)
	}
	return nil, errUknown
}

// implementing the GRPC method
func (e Endpoint) Casing(ctx context.Context, req *request) (*response, error) {
	return e.EP(ctx, req)
}

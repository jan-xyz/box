package lambda

import (
	"context"
	"errors"
	"log"

	"github.com/jan-xyz/box"
	"github.com/jan-xyz/box/example/strsvc/lambda/proto/strsvcv1"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("strsvc")

var errUknown = errors.New("unknown message")

func NewEndpoint() box.Endpoint[*strsvcv1.Request, *strsvcv1.Response] {
	var svc Service = &service{}
	svc = tracingMiddleware{svc: svc, tracer: tracer}
	svc = validationMiddleware{svc: svc}

	return func(ctx context.Context, req *strsvcv1.Request) (*strsvcv1.Response, error) {
		switch m := req.GetMessage().(type) {
		case *strsvcv1.Request_LowerCase:
			upper, err := svc.LowerCase(ctx, m.LowerCase.GetInput())
			if err != nil {
				return &strsvcv1.Response{}, err
			}
			return &strsvcv1.Response{Result: upper}, nil
		case *strsvcv1.Request_UpperCase:
			upper, err := svc.UpperCase(ctx, m.UpperCase.GetInput())
			if err != nil {
				return &strsvcv1.Response{}, err
			}
			return &strsvcv1.Response{Result: upper}, nil
		default:
			log.Printf("unhandled request: %T", m)
		}
		return nil, errUknown
	}
}

package strsvc

import (
	"context"
	"log"

	"github.com/jan-xyz/box"
)

type endpointWrapper func(ctx context.Context, req StringRequest) (StringResponse, error)

func (ep endpointWrapper) EP(ctx context.Context, req StringRequest) (StringResponse, error) {
	return ep(ctx, req)
}

type middlewareWrapper func(next StringEndpoint) StringEndpoint

func (mw middlewareWrapper) MW(next box.Endpoint[StringRequest, StringResponse]) box.Endpoint[StringRequest, StringResponse] {
	return mw(next)
}

var LoggingMiddleware = middlewareWrapper(func(next StringEndpoint) StringEndpoint {
	return endpointWrapper(func(ctx context.Context, req StringRequest) (StringResponse, error) {
		log.Printf("incoming: %s", req.Name)

		resp, err := next.EP(ctx, req)
		log.Printf("outgoing: %s", resp.UpperCaseName)
		return resp, err
	})
})


package strsvc

import (
	"context"
	"log"

	"github.com/jan-xyz/box"
	"github.com/jan-xyz/box/internal/strsvc/proto/strsvcv1"
)

type endpointWrapper func(ctx context.Context, req *strsvcv1.Request) (*strsvcv1.Response, error)

func (ep endpointWrapper) EP(ctx context.Context, req *strsvcv1.Request) (*strsvcv1.Response, error) {
	return ep(ctx, req)
}

type middlewareWrapper func(next StringEndpoint) StringEndpoint

func (mw middlewareWrapper) MW(next box.Endpoint[*strsvcv1.Request, *strsvcv1.Response]) box.Endpoint[*strsvcv1.Request, *strsvcv1.Response] {
	return mw(next)
}

var LoggingMiddleware = middlewareWrapper(func(next StringEndpoint) StringEndpoint {
	return endpointWrapper(func(ctx context.Context, req *strsvcv1.Request) (*strsvcv1.Response, error) {
		log.Printf("incoming: %s", req)

		resp, err := next.EP(ctx, req)
		log.Printf("outgoing: %s", resp)
		return resp, err
	})
})

package lambda

import (
	"context"
	"log"

	"github.com/jan-xyz/box"
	"github.com/jan-xyz/box/example/strsvc/lambda/proto/strsvcv1"
)

var LoggingMiddleware = box.Middleware[*strsvcv1.Request, *strsvcv1.Response](func(next box.Endpoint[*strsvcv1.Request, *strsvcv1.Response]) box.Endpoint[*strsvcv1.Request, *strsvcv1.Response] {
	return box.Endpoint[*strsvcv1.Request, *strsvcv1.Response](func(ctx context.Context, req *strsvcv1.Request) (*strsvcv1.Response, error) {
		log.Printf("incoming: %s", req)

		resp, err := next(ctx, req)
		log.Printf("outgoing: %s", resp)
		return resp, err
	})
})

package box

import (
	"context"
)

// Endpoint is the generic representation of any endpoint
// which takes a request and returns a reponse.
type Endpoint[TIn, TOut any] interface {
	EP(ctx context.Context, req TIn) (TOut, error)
}

// Helper function to create an endpoint from a function
type EndpointFunc[TIn, TOut any] func(ctx context.Context, req TIn) (TOut, error)

func (ep EndpointFunc[TIn, TOut]) EP(ctx context.Context, req TIn) (TOut, error) {
	return ep(ctx, req)
}

// Middleware is an [Endpoint] middleware which can be used to wrap
// around endpoints and decorate them with auxillary functionality, like
// request logging, instrumentation, context enrichment etc.
type Middleware[TIn, TOut any] func(next Endpoint[TIn, TOut]) Endpoint[TIn, TOut]

func Chain[TIn, TOut any](outer Middleware[TIn, TOut], others ...Middleware[TIn, TOut]) Middleware[TIn, TOut] {
	return func(next Endpoint[TIn, TOut]) Endpoint[TIn, TOut] {
		for i := len(others) - 1; i >= 0; i-- { // reverse
			next = others[i](next)
		}
		return outer(next)
	}
}

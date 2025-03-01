package box

import (
	"context"
)

// Endpoint is the generic representation of any endpoint
// which takes a request and returns a reponse.
type Endpoint[TIn, TOut any] = func(ctx context.Context, req TIn) (TOut, error)

// Middleware is an [Endpoint] middleware which can be used to wrap
// around endpoints and decorate them with auxillary functionality, like
// request logging, instrumentation, context enrichment etc.
type Middleware[TIn, TOut any] = func(next Endpoint[TIn, TOut]) Endpoint[TIn, TOut]

// Chain is a convenience function to chain multiple [Middleware]s together
// and finally wraps an [Endpoint].
func Chain[TIn, TOut any](outer Middleware[TIn, TOut], others ...Middleware[TIn, TOut]) Middleware[TIn, TOut] {
	return func(next Endpoint[TIn, TOut]) Endpoint[TIn, TOut] {
		for i := len(others) - 1; i >= 0; i-- { // reverse
			next = others[i](next)
		}
		return outer(next)
	}
}

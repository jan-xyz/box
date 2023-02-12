package box

import (
	"context"
)

type Endpoint[TIn, TOut any] interface {
	EP(ctx context.Context, req TIn) (TOut, error)
}

type Middleware[TIn, TOut any] interface {
	MW(next Endpoint[TIn, TOut]) Endpoint[TIn, TOut]
}

func NewChainBuilder[TIn, TOut any](outer Middleware[TIn, TOut], others ...Middleware[TIn, TOut]) Chain[TIn, TOut] {
	return Chain[TIn, TOut]{
		outer:  outer,
		others: others,
	}
}

// EndpointFunc is a helper function to create an [Endpoint] from just a function.
type EndpointFunc[TIn, TOut any] func(ctx context.Context, req TIn) (TOut, error)

func (ep EndpointFunc[TIn, TOut]) EP(ctx context.Context, req TIn) (TOut, error) {
	return ep(ctx, req)
}

// MiddlewareFunc is a helper function to create a Middelware from just a function
// when combined with an [EndpointFunc]
type MiddlewareFunc[TIn, TOut any] func(next Endpoint[TIn, TOut]) Endpoint[TIn, TOut]

func (mw MiddlewareFunc[TIn, TOut]) MW(next Endpoint[TIn, TOut]) Endpoint[TIn, TOut] {
	return mw(next)
}

type Chain[TIn, TOut any] struct {
	outer  Middleware[TIn, TOut]
	others []Middleware[TIn, TOut]
}

func (c Chain[TIn, TOut]) Build(next Endpoint[TIn, TOut]) Endpoint[TIn, TOut] {
	for i := len(c.others) - 1; i >= 0; i-- { // reverse
		next = c.others[i].MW(next)
	}
	return c.outer.MW(next)
}

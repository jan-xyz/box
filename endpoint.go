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

type EndpointWrapper[TIn, TOut any] func(ctx context.Context, req TIn) (TOut, error)

func (ep EndpointWrapper[TIn, TOut]) EP(ctx context.Context, req TIn) (TOut, error) {
	return ep(ctx, req)
}

type MiddlewareWrapper[TIn, TOut any] func(next Endpoint[TIn, TOut]) Endpoint[TIn, TOut]

func (mw MiddlewareWrapper[TIn, TOut]) MW(next Endpoint[TIn, TOut]) Endpoint[TIn, TOut] {
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

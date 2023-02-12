package awslambdago

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jan-xyz/box"
)

func NewClouadWatchEventHandler[TIn, TOut any](
	decode func(*events.CloudWatchEvent) (TIn, error),
	encode func(TOut) error,
	encodeError func(error) error,
	endpoint box.Endpoint[TIn, TOut],
) CloudWatchEventHandler {
	return cloudWatchEventHandler[TIn, TOut]{
		decode:      decode,
		encode:      encode,
		encodeError: encodeError,
		endpoint:    endpoint,
	}
}

type CloudWatchEventHandler interface {
	Handle(ctx context.Context, e *events.CloudWatchEvent) error
}

type cloudWatchEventHandler[TIn, TOut any] struct {
	decode      func(*events.CloudWatchEvent) (TIn, error)
	encode      func(TOut) error
	encodeError func(error) error
	endpoint    box.Endpoint[TIn, TOut]
}

func (s cloudWatchEventHandler[TIn, TOut]) Handle(ctx context.Context, req *events.CloudWatchEvent) error {
	in, err := s.decode(req)
	if err != nil {
		return s.encodeError(err)
	}
	out, err := s.endpoint.EP(ctx, in)
	if err != nil {
		return s.encodeError(err)
	}
	err = s.encode(out)
	if err != nil {
		return s.encodeError(err)
	}
	return err
}

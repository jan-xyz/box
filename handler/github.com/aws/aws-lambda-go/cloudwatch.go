package awslambdago

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jan-xyz/box"
)

type CloudWatchEventHandler = func(ctx context.Context, e *events.CloudWatchEvent) (any, error)

func NewClouadWatchEventHandler[TIn, TOut any](
	decode func(*events.CloudWatchEvent) (TIn, error),
	endpoint box.Endpoint[TIn, TOut],
) CloudWatchEventHandler {
	return func(ctx context.Context, req *events.CloudWatchEvent) (any, error) {
		in, err := decode(req)
		if err != nil {
			return nil, err
		}
		_, err = endpoint(ctx, in)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

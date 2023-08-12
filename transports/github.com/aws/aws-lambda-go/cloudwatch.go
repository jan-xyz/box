package awslambdago

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jan-xyz/box"
)

type CloudWatchEventTransport = func(ctx context.Context, e *events.CloudWatchEvent) (any, error)

func NewClouadWatchEventTransport[TIn, TOut any](
	decode func(*events.CloudWatchEvent) (TIn, error),
	endpoint box.Endpoint[TIn, TOut],
) CloudWatchEventTransport {
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

package awslambdago

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jan-xyz/box"
)

type APIGatewayTransport = func(ctx context.Context, req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)

func NewAPIGatewayTransport[TIn, TOut any](
	decode func(*events.APIGatewayProxyRequest) (TIn, error),
	encode func(TOut) (*events.APIGatewayProxyResponse, error),
	encodeError func(error) *events.APIGatewayProxyResponse,
	endpoint box.Endpoint[TIn, TOut],
) APIGatewayTransport {
	return func(ctx context.Context, req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		in, err := decode(req)
		if err != nil {
			return encodeError(err), nil
		}
		out, err := endpoint(ctx, in)
		if err != nil {
			return encodeError(err), nil
		}
		resp, err := encode(out)
		if err != nil {
			return encodeError(err), nil
		}
		return resp, err
	}
}

// adds a default HSTS header to the APIGateway response. If the maxAge is 0, it will default to `63072000` seconds
func NewAPIGatewayHSTSMiddleware(transport APIGatewayTransport, maxAge time.Duration) APIGatewayTransport {
	if maxAge == 0 {
		maxAge = 63072000 * time.Second
	}
	return func(ctx context.Context, req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		resp, err := transport(ctx, req)
		if resp.Headers == nil {
			resp.Headers = map[string]string{}
		}
		resp.Headers["strict-transport-security"] = fmt.Sprintf("max-age=%d", int(maxAge.Seconds()))
		return resp, err
	}
}

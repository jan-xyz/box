package awslambdago

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jan-xyz/box"
)

func NewAPIGatewayHandler[TIn, TOut any](
	decode func(*events.APIGatewayProxyRequest) (TIn, error),
	encode func(TOut) (*events.APIGatewayProxyResponse, error),
	encodeError func(error) (*events.APIGatewayProxyResponse, error),
	endpoint box.Endpoint[TIn, TOut],
) APIGatewayHandler {
	return apiGatewayHandler[TIn, TOut]{
		decode:      decode,
		encode:      encode,
		encodeError: encodeError,
		endpoint:    endpoint,
	}
}

type APIGatewayHandler interface {
	Handle(ctx context.Context, e *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)
}

type apiGatewayHandler[TIn, TOut any] struct {
	decode      func(*events.APIGatewayProxyRequest) (TIn, error)
	encode      func(TOut) (*events.APIGatewayProxyResponse, error)
	encodeError func(error) (*events.APIGatewayProxyResponse, error)
	endpoint    box.Endpoint[TIn, TOut]
}

func (s apiGatewayHandler[TIn, TOut]) Handle(ctx context.Context, req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	in, err := s.decode(req)
	if err != nil {
		return s.encodeError(err)
	}
	out, err := s.endpoint.EP(ctx, in)
	if err != nil {
		return s.encodeError(err)
	}
	resp, err := s.encode(out)
	if err != nil {
		return s.encodeError(err)
	}
	return resp, err
}

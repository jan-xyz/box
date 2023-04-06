package awslambdago

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jan-xyz/box"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"
)

type APIGatewayHandler = func(ctx context.Context, req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)

func NewAPIGatewayHandler[TIn, TOut any](
	decode func(*events.APIGatewayProxyRequest) (TIn, error),
	encode func(TOut) (*events.APIGatewayProxyResponse, error),
	encodeError func(error) (*events.APIGatewayProxyResponse, error),
	endpoint box.Endpoint[TIn, TOut],
) APIGatewayHandler {
	return func(ctx context.Context, req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		in, err := decode(req)
		if err != nil {
			return encodeError(err)
		}
		out, err := endpoint(ctx, in)
		if err != nil {
			return encodeError(err)
		}
		resp, err := encode(out)
		if err != nil {
			return encodeError(err)
		}
		return resp, err
	}
}

// implements https://opentelemetry.io/docs/reference/specification/trace/semantic_conventions/instrumentation/aws-lambda/#api-gateway
func NewAPIGatewayTracingMiddleware(handler APIGatewayHandler, tp trace.TracerProvider) APIGatewayHandler {
	tracer := tp.Tracer(box.TracerName)

	return func(ctx context.Context, req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		ctx, span := tracer.Start(ctx, req.Resource, trace.WithAttributes(
			semconv.HTTPRoute(req.Resource),
			semconv.FaaSTriggerHTTP,
			semconv.HTTPScheme(req.Headers["x-forwarded-proto"]),
			semconv.HTTPMethod(req.HTTPMethod),
		))
		defer span.End()

		resp, err := handler(ctx, req)
		span.SetAttributes(semconv.HTTPStatusCode(resp.StatusCode))
		if resp.StatusCode >= 500 {
			span.SetStatus(codes.Error, "")
		}
		return resp, err
	}
}

package awslambdago

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jan-xyz/box"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
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
	Handle(ctx context.Context, req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)
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
	out, err := s.endpoint(ctx, in)
	if err != nil {
		return s.encodeError(err)
	}
	resp, err := s.encode(out)
	if err != nil {
		return s.encodeError(err)
	}
	return resp, err
}

func NewAPIGatewayTracingMiddleware(handler APIGatewayHandler, tp trace.TracerProvider) APIGatewayHandler {
	tracer := tp.Tracer(box.TracerName)
	return apiGatewayTracingMiddleware{handler: handler, tracer: tracer}
}

type apiGatewayTracingMiddleware struct {
	handler APIGatewayHandler
	tracer  trace.Tracer
}

// implements https://opentelemetry.io/docs/reference/specification/trace/semantic_conventions/instrumentation/aws-lambda/#api-gateway
func (s apiGatewayTracingMiddleware) Handle(ctx context.Context, req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	ctx, span := s.tracer.Start(ctx, req.Resource, trace.WithAttributes(
		semconv.HTTPRoute(req.Resource),
		semconv.FaaSTriggerHTTP,
		semconv.HTTPScheme(req.Headers["x-forwarded-proto"]),
		semconv.HTTPMethod(req.HTTPMethod),
	))
	defer span.End()

	resp, err := s.handler.Handle(ctx, req)
	span.SetAttributes(semconv.HTTPStatusCode(resp.StatusCode))
	if resp.StatusCode >= 500 {
		span.SetStatus(codes.Error, "")
	}
	return resp, err
}

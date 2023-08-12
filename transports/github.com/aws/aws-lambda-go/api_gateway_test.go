package awslambdago

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

func makeSureAPIGatewayTransportHasCorrectSignature() {
	h := NewAPIGatewayTransport(
		func(*events.APIGatewayProxyRequest) (string, error) { return "", nil },
		func(string) (*events.APIGatewayProxyResponse, error) { return nil, nil },
		func(error) (*events.APIGatewayProxyResponse, error) { return nil, nil },
		func(context.Context, string) (string, error) { return "", nil },
	)
	lambda.StartHandlerFunc(h)
}

func Test_APIGateway_Handle(t *testing.T) {
	testCases := []struct {
		desc            string
		decodeFunc      func(*events.APIGatewayProxyRequest) (string, error)
		encodeFunc      func(string) (*events.APIGatewayProxyResponse, error)
		encodeErrorFunc func(error) (*events.APIGatewayProxyResponse, error)
		ep              func(context.Context, string) (string, error)
		input           *events.APIGatewayProxyRequest
		want            *events.APIGatewayProxyResponse
		wantErr         bool
	}{
		{
			desc: "test successfully processing request",
			decodeFunc: func(r *events.APIGatewayProxyRequest) (string, error) {
				return r.Body, nil
			},
			encodeFunc: func(s string) (*events.APIGatewayProxyResponse, error) {
				return &events.APIGatewayProxyResponse{Body: s}, nil
			},
			ep: func(context.Context, string) (string, error) {
				return "foo", nil
			},
			input:   &events.APIGatewayProxyRequest{},
			want:    &events.APIGatewayProxyResponse{Body: "foo"},
			wantErr: false,
		},
		{
			desc: "test failing to decode doesn't call endpoint",
			decodeFunc: func(ev *events.APIGatewayProxyRequest) (string, error) {
				return "", errors.New("boom")
			},
			encodeErrorFunc: func(err error) (*events.APIGatewayProxyResponse, error) {
				return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
			},
			input:   &events.APIGatewayProxyRequest{Body: "first"},
			want:    &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError},
			wantErr: true,
		},
		{
			desc: "test failing when the endpoint fails",
			decodeFunc: func(ev *events.APIGatewayProxyRequest) (string, error) {
				return ev.Body, nil
			},
			ep: func(_ context.Context, input string) (string, error) {
				return "", errors.New("boom")
			},
			encodeErrorFunc: func(err error) (*events.APIGatewayProxyResponse, error) {
				return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
			},
			input: &events.APIGatewayProxyRequest{
				Body: "first",
			},
			want:    &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError},
			wantErr: true,
		},
		{
			desc: "test failing when encoding fails",
			decodeFunc: func(r *events.APIGatewayProxyRequest) (string, error) {
				return r.Body, nil
			},
			encodeFunc: func(s string) (*events.APIGatewayProxyResponse, error) {
				return &events.APIGatewayProxyResponse{}, errors.New("boom")
			},
			encodeErrorFunc: func(err error) (*events.APIGatewayProxyResponse, error) {
				return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
			},
			ep: func(context.Context, string) (string, error) {
				return "foo", nil
			},
			input:   &events.APIGatewayProxyRequest{},
			want:    &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError},
			wantErr: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			h := NewAPIGatewayTransport(
				tC.decodeFunc,
				tC.encodeFunc,
				tC.encodeErrorFunc,
				tC.ep,
			)
			resp, err := h(context.Background(), tC.input)

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tC.want, resp)
			}
		})
	}
}

func Test_APIGateway_TracingMiddleware(t *testing.T) {
	// given
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	h := NewAPIGatewayTransport(
		func(*events.APIGatewayProxyRequest) (string, error) { return "", nil },
		func(string) (*events.APIGatewayProxyResponse, error) { return &events.APIGatewayProxyResponse{}, nil },
		func(error) (*events.APIGatewayProxyResponse, error) { return &events.APIGatewayProxyResponse{}, nil },
		func(context.Context, string) (string, error) { return "", nil },
	)
	mw := NewAPIGatewayTracingMiddleware(h, tp)

	// when
	input := &events.APIGatewayProxyRequest{
		Resource:   "some-resource",
		HTTPMethod: "POST",
		Headers:    map[string]string{"x-forwarded-proto": "https"},
	}
	got, err := mw(context.Background(), input)

	// then
	assert.NoError(t, err)
	want := &events.APIGatewayProxyResponse{}
	assert.Equal(t, want, got)

	spans := sr.Ended()
	wantSpanAttributes := []attribute.KeyValue{
		semconv.HTTPRoute("some-resource"),
		semconv.FaaSTriggerHTTP,
		semconv.HTTPScheme("https"),
		semconv.HTTPMethod("POST"),
		semconv.HTTPStatusCode(0),
	}
	assert.Len(t, spans, 1)
	assert.ElementsMatch(t, wantSpanAttributes, spans[0].Attributes())
	assert.Equal(t, "some-resource", spans[0].Name())
}

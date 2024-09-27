package awslambdago

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/stretchr/testify/assert"
)

func makeSureAPIGatewayTransportHasCorrectSignature() {
	h := NewAPIGatewayTransport(
		func(*events.APIGatewayProxyRequest) (string, error) { return "", nil },
		func(string) (*events.APIGatewayProxyResponse, error) { return nil, nil },
		func(error) *events.APIGatewayProxyResponse { return nil },
		func(context.Context, string) (string, error) { return "", nil },
	)
	lambda.StartHandlerFunc(h)
}

func Test_APIGateway_Handle(t *testing.T) {
	testCases := []struct {
		desc            string
		decodeFunc      func(*events.APIGatewayProxyRequest) (string, error)
		encodeFunc      func(string) (*events.APIGatewayProxyResponse, error)
		encodeErrorFunc func(error) *events.APIGatewayProxyResponse
		ep              func(context.Context, string) (string, error)
		input           *events.APIGatewayProxyRequest
		want            *events.APIGatewayProxyResponse
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
			input: &events.APIGatewayProxyRequest{},
			want:  &events.APIGatewayProxyResponse{Body: "foo"},
		},
		{
			desc: "test failing to decode doesn't call endpoint",
			decodeFunc: func(ev *events.APIGatewayProxyRequest) (string, error) {
				return "", errors.New("boom")
			},
			encodeErrorFunc: func(err error) *events.APIGatewayProxyResponse {
				return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
			},
			input: &events.APIGatewayProxyRequest{Body: "first"},
			want:  &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError},
		},
		{
			desc: "test failing when the endpoint fails",
			decodeFunc: func(ev *events.APIGatewayProxyRequest) (string, error) {
				return ev.Body, nil
			},
			ep: func(_ context.Context, input string) (string, error) {
				return "", errors.New("boom")
			},
			encodeErrorFunc: func(err error) *events.APIGatewayProxyResponse {
				return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
			},
			input: &events.APIGatewayProxyRequest{
				Body: "first",
			},
			want: &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError},
		},
		{
			desc: "test failing when encoding fails",
			decodeFunc: func(r *events.APIGatewayProxyRequest) (string, error) {
				return r.Body, nil
			},
			encodeFunc: func(s string) (*events.APIGatewayProxyResponse, error) {
				return &events.APIGatewayProxyResponse{}, errors.New("boom")
			},
			encodeErrorFunc: func(err error) *events.APIGatewayProxyResponse {
				return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
			},
			ep: func(context.Context, string) (string, error) {
				return "foo", nil
			},
			input: &events.APIGatewayProxyRequest{},
			want:  &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError},
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

			assert.NoError(t, err)
			assert.Equal(t, tC.want, resp)
		})
	}
}

func Test_APIGateway_HSTSMiddleware(t *testing.T) {
	testCases := []struct {
		desc   string
		maxAge time.Duration
		want   *events.APIGatewayProxyResponse
	}{
		{
			desc:   "enrich with default value",
			maxAge: time.Duration(0),
			want: &events.APIGatewayProxyResponse{
				StatusCode:      0,
				Headers:         map[string]string{"strict-transport-security": "max-age=63072000"},
				Body:            "",
				IsBase64Encoded: false,
			},
		},
		{
			desc:   "enrich with value",
			maxAge: time.Minute,
			want: &events.APIGatewayProxyResponse{
				StatusCode:      0,
				Headers:         map[string]string{"strict-transport-security": "max-age=60"},
				Body:            "",
				IsBase64Encoded: false,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			h := NewAPIGatewayTransport(
				func(*events.APIGatewayProxyRequest) (string, error) { return "", nil },
				func(string) (*events.APIGatewayProxyResponse, error) { return &events.APIGatewayProxyResponse{}, nil },
				func(error) *events.APIGatewayProxyResponse { return &events.APIGatewayProxyResponse{} },
				func(context.Context, string) (string, error) { return "", nil },
			)
			mw := NewAPIGatewayHSTSMiddleware(h, tC.maxAge)

			got, err := mw(context.Background(), nil)

			assert.NoError(t, err)
			assert.Equal(t, tC.want, got)
		})
	}
}

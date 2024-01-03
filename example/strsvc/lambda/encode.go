package lambda

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jan-xyz/box/example/strsvc/lambda/proto/strsvcv1"
	"google.golang.org/protobuf/proto"
)

var errNoUpper = errors.New("unable to uppercase name")

func EncodeAPIGateway(m *strsvcv1.Response) (*events.APIGatewayProxyResponse, error) {
	resp, err := proto.Marshal(m)
	if err != nil {
		return nil, err
	}
	body := base64.StdEncoding.EncodeToString(resp)
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       body,
	}, nil
}

func EncodeErrorAPIGateway(err error) *events.APIGatewayProxyResponse {
	if errors.Is(err, errNoUpper) {
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: "cannot uppercase string"}
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       "oops, an error happened",
	}
}

func EncodeHTTP(m *strsvcv1.Response, w http.ResponseWriter) {
	resp, err := proto.Marshal(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

func EncodeErrorHTTP(err error, w http.ResponseWriter) {
	if errors.Is(err, errNoUpper) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

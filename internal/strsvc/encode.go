package strsvc

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jan-xyz/box/internal/strsvc/proto/strsvcv1"
	"google.golang.org/protobuf/proto"
)

var errNoUpper = errors.New("unable to uppercase name")

func EncodeSQS(resp *strsvcv1.Response) error {
	if resp.GetResult() == "" {
		return errNoUpper
	}
	return nil
}

func EncodeAPIGateway(m *strsvcv1.Response) (*events.APIGatewayProxyResponse, error) {
	resp, err := proto.Marshal(m)
	if err != nil {
		return nil, err
	}
	body := base64.StdEncoding.EncodeToString(resp)
	return &events.APIGatewayProxyResponse{
		StatusCode:        http.StatusOK,
		Body:              body,
	}, nil
}

func EncodeErrorAPIGateway(err error) (*events.APIGatewayProxyResponse, error) {
	log.Printf("failed: %s", err)
	if errors.Is(err, errNoUpper) {
return &events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: "cannot uppercase string"}, nil
	}
	return &events.APIGatewayProxyResponse{
		StatusCode:        http.StatusInternalServerError,
		Body:              "oops, an error happened",
	}, nil
}

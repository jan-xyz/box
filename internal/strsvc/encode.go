package strsvc

import (
	"errors"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

var errNoUpper = errors.New("unable to uppercase name")

func EncodeSQS(resp StringResponse) error {
	if resp.UpperCaseName != "" {
		return nil
	}
	return errNoUpper
}

func EncodeAPIGateway(resp StringResponse) (*events.APIGatewayProxyResponse, error) {
	if resp.UpperCaseName == "" {
		return nil, errNoUpper
	}
	return &events.APIGatewayProxyResponse{
		StatusCode:        http.StatusOK,
		Body:              resp.UpperCaseName,
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

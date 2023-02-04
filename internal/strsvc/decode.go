package strsvc

import "github.com/aws/aws-lambda-go/events"

func DecodeSQS(m events.SQSMessage) (StringRequest, error) {
	return StringRequest{
		Name: m.Body,
	}, nil
}

func DecodeAPIGateway(r *events.APIGatewayProxyRequest) (StringRequest, error) {
	return StringRequest{
		Name: r.Body,
	}, nil
}

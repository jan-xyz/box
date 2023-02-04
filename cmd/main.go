package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jan-xyz/box"
	"github.com/jan-xyz/box/internal/strsvc"
)

func main() {
	// setup handler
	c := box.NewChainBuilder[strsvc.StringRequest, strsvc.StringResponse](
		strsvc.LoggingMiddleware,
	)
	ep := c.Build(strsvc.NewEndpoint()).EP
	sqsHandler := box.NewSQSHandler(
		false,
		strsvc.DecodeSQS,
		strsvc.EncodeSQS,
		ep,
	)

	// Setup Lambda
	sqsResp := sqsHandler.Handle(
		context.Background(),
		&events.SQSEvent{Records: []events.SQSMessage{
			{Body: "foo"},
		}},
	)
	log.Printf("sqs: %#v", sqsResp)

	apiGWHandler := box.NewAPIGatewayHandler(
		strsvc.DecodeAPIGateway,
		strsvc.EncodeAPIGateway,
		strsvc.EncodeErrorAPIGateway,
		ep,
	)

	apiGWResp, err := apiGWHandler.Handle(
		context.Background(),
		&events.APIGatewayProxyRequest{Body: "bar"},
	)
	if err != nil {
		panic(err)
	}
	log.Printf("sqs: %#v", apiGWResp)
}

package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jan-xyz/box"
	"github.com/jan-xyz/box/internal/strsvc"
)

func main() {
	// setup endpoint
	c := box.NewChainBuilder[strsvc.StringRequest, strsvc.StringResponse](
		strsvc.LoggingMiddleware,
	)
	ep := c.Build(strsvc.NewEndpoint()).EP

	// connect endpoint to SQS
	sqsHandler := box.NewSQSHandler(
		false,
		strsvc.DecodeSQS,
		strsvc.EncodeSQS,
		ep,
	)

	// simulate SQS invocation
	sqsResp := sqsHandler.Handle(
		context.Background(),
		&events.SQSEvent{Records: []events.SQSMessage{
			{Body: "foo"},
		}},
	)
	log.Printf("sqs: %#v", sqsResp)

	// connect endpoint to APIGateway
	apiGWHandler := box.NewAPIGatewayHandler(
		strsvc.DecodeAPIGateway,
		strsvc.EncodeAPIGateway,
		strsvc.EncodeErrorAPIGateway,
		ep,
	)

	// simulate APIGateway invocation
	apiGWResp, err := apiGWHandler.Handle(
		context.Background(),
		&events.APIGatewayProxyRequest{Body: "bar"},
	)
	if err != nil {
		panic(err)
	}
	log.Printf("sqs: %#v", apiGWResp)
}

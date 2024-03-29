package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"log"
	"net/http/httptest"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jan-xyz/box"
	strsvc "github.com/jan-xyz/box/example/strsvc/lambda"
	"github.com/jan-xyz/box/example/strsvc/lambda/proto/strsvcv1"
	awslambdago "github.com/jan-xyz/box/transports/github.com/aws/aws-lambda-go"
	boxhttp "github.com/jan-xyz/box/transports/net/http"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/proto"
)

func main() {
	// setup endpoint with it's middlewares
	mw := box.Chain(
		box.EndpointLogging[*strsvcv1.Request, *strsvcv1.Response](),
		box.EndpointTracing[*strsvcv1.Request, *strsvcv1.Response](otel.GetTracerProvider()),
		box.EndpointMetrics[*strsvcv1.Request, *strsvcv1.Response](otel.GetMeterProvider()),
	)
	ep := mw(strsvc.NewEndpoint())

	// connect endpoint to SQS
	sqsTransport := awslambdago.NewSQSTransport(
		false,
		strsvc.DecodeSQS,
		ep,
	)

	// connect endpoint to APIGateway
	apiGWTransport := awslambdago.NewAPIGatewayTransport(
		strsvc.DecodeAPIGateway,
		strsvc.EncodeAPIGateway,
		strsvc.EncodeErrorAPIGateway,
		ep,
	)
	apiGWTransport = awslambdago.NewAPIGatewayTracingMiddleware(apiGWTransport, otel.GetTracerProvider())

	// connect to HTTP
	httpServer := boxhttp.NewHTTPServer(
		strsvc.DecodeHTTP,
		strsvc.EncodeHTTP,
		strsvc.EncodeErrorHTTP,
		ep,
	)

	// test input
	requests := []*strsvcv1.Request{
		{
			Message: &strsvcv1.Request_UpperCase{
				UpperCase: &strsvcv1.UpperCase{
					Input: "Foo",
				},
			},
		},
		{
			Message: &strsvcv1.Request_LowerCase{
				LowerCase: &strsvcv1.LowerCase{
					Input: "Bar",
				},
			},
		},
		{
			Message: &strsvcv1.Request_LowerCase{
				LowerCase: &strsvcv1.LowerCase{
					Input: "",
				},
			},
		},
	}

	// simular incoming events via lambda.Start
	for _, m := range requests {
		marshalledM, err := proto.Marshal(m)
		if err != nil {
			panic(err)
		}
		body := base64.StdEncoding.EncodeToString(marshalledM)

		// simulate SQS invocation
		sqsResp, _ := sqsTransport(
			context.Background(),
			&events.SQSEvent{Records: []events.SQSMessage{
				{Body: body, MessageId: "the message"},
			}},
		)
		log.Printf("sqs: %#v", sqsResp)

		// simulate APIGateway invocation
		apiGWResp, err := apiGWTransport(
			context.Background(),
			&events.APIGatewayProxyRequest{Body: body},
		)
		if err != nil {
			panic(err)
		}
		log.Printf("api gw: %#v", apiGWResp)

		// simulate HTTP invocation
		srv := httptest.NewServer(httpServer)
		httpResp, err := srv.Client().Post(srv.URL, "application/octetstream", bytes.NewReader(marshalledM))
		if err != nil {
			panic(err)
		}
		log.Printf("http: %#v", httpResp)
	}
}

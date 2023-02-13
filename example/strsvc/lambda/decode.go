package lambda

import (
	"encoding/base64"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jan-xyz/box/example/strsvc/lambda/proto/strsvcv1"
	"google.golang.org/protobuf/proto"
)

type request = strsvcv1.CasingRequest

func DecodeSQS(m events.SQSMessage) (*request, error) {
	body, err := base64.StdEncoding.DecodeString(m.Body)
	if err != nil {
		return nil, err
	}

	msg := &request{}
	err = proto.Unmarshal(body, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func DecodeAPIGateway(r *events.APIGatewayProxyRequest) (*request, error) {
	body, err := base64.StdEncoding.DecodeString(r.Body)
	if err != nil {
		return nil, err
	}

	msg := &request{}
	err = proto.Unmarshal(body, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func DecodeHTTP(r *http.Request) (*request, error) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			panic(err)
		}
	}()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	msg := &request{}
	err = proto.Unmarshal(body, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

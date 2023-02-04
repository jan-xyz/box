package strsvc

import (
	"encoding/base64"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jan-xyz/box/internal/strsvc/proto/strsvcv1"
	"google.golang.org/protobuf/proto"
)

func DecodeSQS(m events.SQSMessage) (*strsvcv1.Request, error) {
	body, err := base64.StdEncoding.DecodeString(m.Body)
	if err != nil {
		return nil, err
	}

	msg := &strsvcv1.Request{}
	err = proto.Unmarshal(body, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func DecodeAPIGateway(r *events.APIGatewayProxyRequest) (*strsvcv1.Request, error) {
	body, err := base64.StdEncoding.DecodeString(r.Body)
	if err != nil {
		return nil, err
	}

	msg := &strsvcv1.Request{}
	err = proto.Unmarshal(body, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

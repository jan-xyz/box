package strsvc

import (
	"context"
)

type StringRequest struct {
	Name string
}
type StringResponse struct {
	UpperCaseName string
}

func NewEndpoint() StringEndpoint {
	svc := &service{}
	return endpointer{svc: svc}
}

type endpointer struct {
	svc Service
}

func (e endpointer) EP(ctx context.Context, req StringRequest) (StringResponse, error) {
	upper, err := e.svc.UpperCase(ctx, req.Name)
	if err != nil {
		return StringResponse{}, err
	}
	return StringResponse{UpperCaseName: upper}, nil
}

type StringEndpoint interface {
	EP(ctx context.Context, req StringRequest) (StringResponse, error)
}

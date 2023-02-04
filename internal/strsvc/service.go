package strsvc

import (
	"context"
	"strings"
)

type Service interface {
	UpperCase(ctx context.Context, s string) (string, error)
}

type service struct{}

func (svc *service) UpperCase(ctx context.Context, s string) (string, error) {
	return strings.ToUpper(s), nil
}

package lambda

import (
	"context"
	"strings"
)

type Service interface {
	UpperCase(ctx context.Context, s string) (string, error)
	LowerCase(ctx context.Context, s string) (string, error)
}

type service struct{}

func (svc *service) UpperCase(ctx context.Context, s string) (string, error) {
	return strings.ToUpper(s), nil
}

func (svc *service) LowerCase(ctx context.Context, s string) (string, error) {
	return strings.ToLower(s), nil
}

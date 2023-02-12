package box

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	records := []string{}

	recordingMiddleware := func(name string) Middleware[string, string] {
		return MiddlewareFunc[string, string](func(next Endpoint[string, string]) Endpoint[string, string] {
			return EndpointFunc[string, string](func(ctx context.Context, req string) (string, error) {
				records = append(records, fmt.Sprintf("inc-%s", name))
				resp, err := next.EP(ctx, req)
				records = append(records, fmt.Sprintf("out-%s", name))
				return resp, err
			})
		})
	}

	// setup endpoint with it's middlewares
	c := NewChainBuilder(
		recordingMiddleware("first"),
		recordingMiddleware("second"),
		recordingMiddleware("third"),
	)

	var ep Endpoint[string, string] = EndpointFunc[string, string](func(_ context.Context, req string) (string, error) {
		records = append(records, req)
		return "response", nil
	})

	ep = c.Build(ep)

	resp, err := ep.EP(context.Background(), "request")

	assert.NoError(t, err)
	assert.Equal(t, "response", resp)

	expectedRecords := []string{
		"inc-first",
		"inc-second",
		"inc-third",
		"request",
		"out-third",
		"out-second",
		"out-first",
	}
	assert.Equal(t, expectedRecords, records)
}

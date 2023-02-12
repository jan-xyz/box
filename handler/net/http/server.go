package http

import (
	"net/http"

	"github.com/jan-xyz/box"
)

func NewHTTPServer[TIn, TOut any](
	decode func(*http.Request) (TIn, error),
	encode func(TOut, http.ResponseWriter),
	encodeError func(error, http.ResponseWriter),
	endpoint box.Endpoint[TIn, TOut],
) httpServer[TIn, TOut] {
	return httpServer[TIn, TOut]{
		decode:      decode,
		encode:      encode,
		encodeError: encodeError,
		endpoint:    endpoint,
	}
}

type httpServer[TIn, TOut any] struct {
	decode      func(*http.Request) (TIn, error)
	encode      func(TOut, http.ResponseWriter)
	encodeError func(error, http.ResponseWriter)
	endpoint    box.Endpoint[TIn, TOut]
}

// ServeHTTP implements the http.Handler
func (s httpServer[TIn, TOut]) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	in, err := s.decode(req)
	if err != nil {
		s.encodeError(err, w)
		return
	}
	out, err := s.endpoint(ctx, in)
	if err != nil {
		s.encodeError(err, w)
		return
	}
	s.encode(out, w)
}

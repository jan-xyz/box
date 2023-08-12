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
) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		in, err := decode(req)
		if err != nil {
			encodeError(err, w)
			return
		}
		out, err := endpoint(ctx, in)
		if err != nil {
			encodeError(err, w)
			return
		}
		encode(out, w)
	}
}

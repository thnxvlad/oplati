package hmiddlewares

import "net/http"

func UseMiddlewares(handler http.Handler, mws []func(next http.Handler) http.Handler) http.Handler {
	for _, v := range mws {
		handler = v(handler)
	}

	return handler
}

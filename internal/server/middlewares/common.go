package middlewares

import "net/http"

type HandlerFunc func(http.ResponseWriter, *http.Request) error

type Middleware func(HandlerFunc) HandlerFunc

func ChainMiddleware(h HandlerFunc, middlewares ...Middleware) HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

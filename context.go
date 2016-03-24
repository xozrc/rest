package rest

import (
	"net/http"
)

import (
	"fmt"
	"golang.org/x/net/context"
)

var _ = fmt.Print

//context handler
type ContextHandler interface {
	ServeHTTPContext(ctx context.Context, rw http.ResponseWriter, req *http.Request)
}

//context handler func
type ContextHandlerFunc func(ctx context.Context, rw http.ResponseWriter, req *http.Request)

//implement context handler interface
func (chf ContextHandlerFunc) ServeHTTPContext(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	chf(ctx, rw, req)
}

//middleware context handler
type MiddlewareContextHandler interface {
	ServeHTTPContext(ctx context.Context, rw http.ResponseWriter, req *http.Request, next ContextHandler)
}

type MiddlewareContextHandlerFunc func(ctx context.Context, rw http.ResponseWriter, req *http.Request, next ContextHandler)

func (chf MiddlewareContextHandlerFunc) ServeHTTPContext(ctx context.Context, rw http.ResponseWriter, req *http.Request, next ContextHandler) {
	chf(ctx, rw, req, next)
}

func chainHandlers(ch MiddlewareContextHandler, nh ContextHandler) ContextHandler {
	return ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		ch.ServeHTTPContext(ctx, rw, req, nh)
	})
}

//http handler
func HTTPHandlers(ctx context.Context, tailHandler ContextHandler, handlers ...MiddlewareContextHandler) http.Handler {
	lenOfHandlers := len(handlers)

	f := func(rw http.ResponseWriter, req *http.Request) {

		nh := tailHandler
		for i := lenOfHandlers - 1; i >= 0; i-- {
			ch := handlers[i]
			nh = chainHandlers(ch, nh)
		}
		nh.ServeHTTPContext(ctx, rw, req)
	}
	return http.HandlerFunc(f)
}

func HTTPHandler(ctx context.Context, h ContextHandler) http.Handler {
	f := func(rw http.ResponseWriter, req *http.Request) {
		h.ServeHTTPContext(ctx, rw, req)
	}
	return http.HandlerFunc(f)
}

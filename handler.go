package rest

import (
	"encoding/json"
	"golang.org/x/net/context"
	"net/http"
)

type Handler interface {
	Handle(ctx context.Context, rw http.ResponseWriter, req *http.Request) (code int, result interface{}, err error)
}

type HandleFunc func(ctx context.Context, rw http.ResponseWriter, req *http.Request) (code int, result interface{}, err error)

func (h HandleFunc) Handle(ctx context.Context, rw http.ResponseWriter, req *http.Request) (code int, result interface{}, err error) {
	return h(ctx, rw, req)
}

func RestHandler(ctx context.Context, h Handler) http.Handler {
	f := func(rw http.ResponseWriter, req *http.Request) {
		code, result, err := h.Handle(ctx, rw, req)
		restReturn(code, result, err, rw, req)
	}
	return http.HandlerFunc(f)
}

func restReturn(code int, result interface{}, err error, rw http.ResponseWriter, req *http.Request) {
	if result == nil && err == nil {
		panic("rest return error")
	}

	rw.WriteHeader(code)
	rw.Header().Set("Content-Type", "application/json")
	if code == http.StatusOK && result != nil {
		rb, err := json.Marshal(result)
		if err != nil {
			panic(err.Error())
		}
		_, err = rw.Write(rb)
		if err != nil {
			panic(err.Error())
		}
		return
	}

	if err != nil {
		rb, err := json.Marshal(err)
		if err != nil {
			panic(err.Error())
		}
		_, err = rw.Write(rb)
		if err != nil {
			panic(err.Error())
		}
		return
	}
	panic("rest return error")
}

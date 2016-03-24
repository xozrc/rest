package rest

import (
	"net/http"
)

import (
	"fmt"
	"github.com/xozrc/pkg/httputils"
	"golang.org/x/net/context"
)

var _ = fmt.Print

const (
	formKey = "xozrc.rest.form"
)

//form context
func WithForm(parent context.Context, form interface{}) context.Context {
	return context.WithValue(parent, formKey, form)
}

func FormInContext(ctx context.Context) interface{} {
	return ctx.Value(formKey)
}

//bind form middleware handler
func BindFormHandler(form interface{}) MiddlewareContextHandler {
	f := func(ctx context.Context, rw http.ResponseWriter, req *http.Request, next ContextHandler) {
		err := httputils.Bind(req, form)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		cctx := WithForm(ctx, form)
		next.ServeHTTPContext(cctx, rw, req)
	}
	return MiddlewareContextHandlerFunc(f)
}

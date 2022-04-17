package chai

import (
	"net/http"

	"github.com/zhamlin/chi-openapi/pkg/openapi/operations"
)

type ResHandlerFunc[Res any, Err ErrType] func(http.ResponseWriter, *http.Request) (Res, int, Err)

type resHandler[Res any, Err ErrType] struct {
	method    string
	pattern   string
	fn        ResHandlerFunc[Res, Err]
	res       *Res
	err       *Err
	op        *operations.Operation
	respondFn ResponderFunc[Res]
	errorFn   ErrorResponderFunc
}

func NewResHandler[Res any, Err ErrType](method string, pattern string, fn ResHandlerFunc[Res, Err]) ResHandler[Res, Err] {
	return &resHandler[Res, Err]{
		method:    method,
		pattern:   pattern,
		fn:        fn,
		respondFn: defaultResponder[Res],
		errorFn:   DefaultErrorResponder,
		op:        &operations.Operation{},
	}
}

func (h *resHandler[Res, Err]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, code, err := h.fn(w, r)
	if handleErr(w, r, err, code, h.errorFn) {
		return
	}
	h.respondFn(w, r, code, res)
}

type ResHandler[Res any, Err ErrType] interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	ID(id string) ResHandler[Res, Err]
	Tags(tags ...string) ResHandler[Res, Err]
	Summary(summary string) ResHandler[Res, Err]
	Deprecated() ResHandler[Res, Err]
	Security(name string, scopes ...string) ResHandler[Res, Err]
	NoSecurity() ResHandler[Res, Err]
	Extensions(data operations.ExtensionData) ResHandler[Res, Err]
	WithResponder(respondFn ResponderFunc[Res]) ResHandler[Res, Err]
	WithErrorResponder(errorFn ErrorResponderFunc) ResHandler[Res, Err]
}

func (h *resHandler[Res, Err]) WithResponder(respondFn ResponderFunc[Res]) ResHandler[Res, Err] {
	h.respondFn = respondFn
	return h
}

func (h *resHandler[Res, Err]) WithErrorResponder(errorFn ErrorResponderFunc) ResHandler[Res, Err] {
	h.errorFn = errorFn
	return h
}

func (h *resHandler[Res, Err]) Extensions(data operations.ExtensionData) ResHandler[Res, Err] {
	var err error
	*h.op, err = operations.Extensions(data)(nil, *h.op)
	requireValidSpec(err, h.method, h.pattern)
	return h
}

func (h *resHandler[Res, Err]) NoSecurity() ResHandler[Res, Err] {
	var err error
	*h.op, err = operations.NoSecurity()(nil, *h.op)
	requireValidSpec(err, h.method, h.pattern)
	return h
}

func (h *resHandler[Res, Err]) Security(name string, scopes ...string) ResHandler[Res, Err] {
	var err error
	*h.op, err = operations.Security(name, scopes...)(nil, *h.op)
	requireValidSpec(err, h.method, h.pattern)
	return h
}

func (h *resHandler[Res, Err]) ID(id string) ResHandler[Res, Err] {
	var err error
	*h.op, err = operations.ID(id)(nil, *h.op)
	requireValidSpec(err, h.method, h.pattern)
	return h
}

func (h *resHandler[Res, Err]) Tags(tags ...string) ResHandler[Res, Err] {
	var err error
	*h.op, err = operations.Tags(tags...)(nil, *h.op)
	requireValidSpec(err, h.method, h.pattern)
	return h
}

func (h *resHandler[Res, Err]) Deprecated() ResHandler[Res, Err] {
	var err error
	*h.op, err = operations.Deprecated()(nil, *h.op)
	requireValidSpec(err, h.method, h.pattern)
	return h
}

func (h *resHandler[Res, Err]) Summary(summary string) ResHandler[Res, Err] {
	var err error
	*h.op, err = operations.Summary(summary)(nil, *h.op)
	requireValidSpec(err, h.method, h.pattern)
	return h
}
func (h *resHandler[Res, Err]) WithSpec(op *operations.Operation) ResHandler[Res, Err] {
	h.op = op
	return h
}

func (h *resHandler[Res, Err]) Res() any {
	return h.res
}

func (h *resHandler[Res, Err]) Err() any {
	return h.err
}

func (h *resHandler[Res, Err]) Handler() any {
	return h.fn
}

func (h *resHandler[Res, Err]) Op() *operations.Operation {
	return h.op
}

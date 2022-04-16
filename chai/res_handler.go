package chai

import (
	"net/http"

	"github.com/zhamlin/chi-openapi/pkg/openapi/operations"
)

type ResHandlerFunc[Res any, Err ErrType] func(http.ResponseWriter, *http.Request) (Res, int, Err)

type ResHandler[Res any, Err ErrType] struct {
	method    string
	pattern   string
	fn        ResHandlerFunc[Res, Err]
	res       *Res
	err       *Err
	op        operations.Operation
	respondFn ResponderFunc[Res]
	errorFn   ErrorResponderFunc
}

func NewResHandler[Res any, Err ErrType](method string, pattern string, fn ResHandlerFunc[Res, Err]) *ResHandler[Res, Err] {
	return &ResHandler[Res, Err]{
		method:    method,
		pattern:   pattern,
		fn:        fn,
		respondFn: defaultResponder[Res],
		errorFn:   DefaultErrorResponder,
		op:        operations.Operation{},
	}
}

func (h *ResHandler[Res, Err]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, code, err := h.fn(w, r)
	if handleErr(w, r, err, code, h.errorFn) {
		return
	}
	h.respondFn(w, r, code, res)
}

func (h *ResHandler[Res, Err]) WithSpec(op operations.Operation) *ResHandler[Res, Err] {
	h.op = op
	return h
}

func (h *ResHandler[Res, Err]) WithResponder(respondFn ResponderFunc[Res]) *ResHandler[Res, Err] {
	h.respondFn = respondFn
	return h
}

func (h *ResHandler[Res, Err]) WithErrorResponder(errorFn ErrorResponderFunc) *ResHandler[Res, Err] {
	h.errorFn = errorFn
	return h
}

func (h *ResHandler[Res, Err]) Res() any {
	return h.res
}

func (h *ResHandler[Res, Err]) Err() any {
	return h.err
}

func (h *ResHandler[Res, Err]) Handler() any {
	return h.fn
}

func (h *ResHandler[Res, Err]) Op() operations.Operation {
	return h.op
}

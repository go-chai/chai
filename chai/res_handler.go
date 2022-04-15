package chai

import (
	"net/http"

	"github.com/go-openapi/spec"
)

type ResHandlerFunc[Res any, Err ErrType] func(http.ResponseWriter, *http.Request) (Res, int, Err)

type ResHandler[Res any, Err ErrType] struct {
	f               ResHandlerFunc[Res, Err]
	res             *Res
	err             *Err
	swagAnnotations string
	spec            *spec.Operation
	respondFn       ResponderFunc[Res]
	errorFn         ErrorResponderFunc
}

func NewResHandler[Res any, Err ErrType](h ResHandlerFunc[Res, Err]) *ResHandler[Res, Err] {
	return &ResHandler[Res, Err]{
		f: h,
	}
}

func (h *ResHandler[Res, Err]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, code, err := h.f(w, r)
	if handleErr(w, r, err, code, h.errorFn) {
		return
	}
	h.respondFn(w, r, code, res)
}

func (h *ResHandler[Res, Err]) WithSwagAnnotations(swagAnnotations string) *ResHandler[Res, Err] {
	h.swagAnnotations = swagAnnotations
	return h
}

func (h *ResHandler[Res, Err]) WithSpec(spec *spec.Operation) *ResHandler[Res, Err] {
	h.spec = spec
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
	return h.f
}

func (h *ResHandler[Res, Err]) Docs() string {
	return h.swagAnnotations
}

package chai

import (
	"net/http"

	"github.com/go-openapi/spec"
)

type ReqResHandlerFunc[Req any, Res any, Err ErrType] func(Req, http.ResponseWriter, *http.Request) (Res, int, Err)

type ReqResHandler[Req any, Res any, Err ErrType] struct {
	f               ReqResHandlerFunc[Req, Res, Err]
	req             *Req
	res             *Res
	err             *Err
	swagAnnotations string
	spec            *spec.Operation
	decodeFn        DecoderFunc[Req]
	validateFn      ValidatorFunc[Req]
	respondFn       ResponderFunc[Res]
	errorFn         ErrorResponderFunc
}

func NewReqResHandler[Req any, Res any, Err ErrType](h ReqResHandlerFunc[Req, Res, Err]) *ReqResHandler[Req, Res, Err] {
	return &ReqResHandler[Req, Res, Err]{
		f:          h,
		decodeFn:   defaultDecoder[Req],
		respondFn:  defaultResponder[Res],
		validateFn: defaultValidator[Req],
		errorFn:    DefaultErrorResponder,
	}
}

func (h *ReqResHandler[Req, Res, Err]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := h.decodeFn(r)
	if handleErr(w, r, err, http.StatusBadRequest, h.errorFn) {
		return
	}
	err = h.validateFn(req)
	if handleErr(w, r, err, http.StatusBadRequest, h.errorFn) {
		return
	}
	res, code, err := h.f(req, w, r)
	if handleErr(w, r, err, code, h.errorFn) {
		return
	}
	h.respondFn(w, r, code, res)
}

func (h *ReqResHandler[Req, Res, Err]) WithSwagAnnotations(swagAnnotations string) *ReqResHandler[Req, Res, Err] {
	h.swagAnnotations = swagAnnotations
	return h
}

func (h *ReqResHandler[Req, Res, Err]) WithDecoder(decodeFn DecoderFunc[Req]) *ReqResHandler[Req, Res, Err] {
	h.decodeFn = decodeFn
	return h
}

func (h *ReqResHandler[Req, Res, Err]) WithValidator(validateFn ValidatorFunc[Req]) *ReqResHandler[Req, Res, Err] {
	h.validateFn = validateFn
	return h
}

func (h *ReqResHandler[Req, Res, Err]) WithResponder(respondFn ResponderFunc[Res]) *ReqResHandler[Req, Res, Err] {
	h.respondFn = respondFn
	return h
}

func (h *ReqResHandler[Req, Res, Err]) WithErrorResponder(errorFn ErrorResponderFunc) *ReqResHandler[Req, Res, Err] {
	h.errorFn = errorFn
	return h
}

func (h *ReqResHandler[Req, Res, Err]) WithSpec(spec *spec.Operation) *ReqResHandler[Req, Res, Err] {
	h.spec = spec
	return h
}

func (h *ReqResHandler[Req, Res, Err]) Req() any {
	return h.req
}

func (h *ReqResHandler[Req, Res, Err]) Res() any {
	return h.res
}

func (h *ReqResHandler[Req, Res, Err]) Err() any {
	return h.err
}

func (h *ReqResHandler[Req, Res, Err]) Handler() any {
	return h.f
}

func (h *ReqResHandler[Req, Res, Err]) Docs() string {
	return h.swagAnnotations
}

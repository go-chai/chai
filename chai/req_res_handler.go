package chai

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zhamlin/chi-openapi/pkg/openapi/operations"
)

type ReqResHandlerFunc[Req any, Res any, Err ErrType] func(Req, http.ResponseWriter, *http.Request) (Res, int, Err)

type ReqResHandler[Req any, Res any, Err ErrType] struct {
	method     string
	pattern    string
	fn         ReqResHandlerFunc[Req, Res, Err]
	req        *Req
	res        *Res
	err        *Err
	op         operations.Operation
	decodeFn   DecoderFunc[Req]
	validateFn ValidatorFunc[Req]
	respondFn  ResponderFunc[Res]
	errorFn    ErrorResponderFunc
}

func NewReqResHandler[Req any, Res any, Err ErrType](method string, pattern string, fn ReqResHandlerFunc[Req, Res, Err]) *ReqResHandler[Req, Res, Err] {
	return &ReqResHandler[Req, Res, Err]{
		method:     method,
		pattern:    pattern,
		fn:         fn,
		decodeFn:   defaultDecoder[Req],
		validateFn: defaultValidator[Req],
		respondFn:  defaultResponder[Res],
		errorFn:    DefaultErrorResponder,
		op:         operations.Operation{},
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
	res, code, err := h.fn(req, w, r)
	if handleErr(w, r, err, code, h.errorFn) {
		return
	}
	h.respondFn(w, r, code, res)
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

func (h *ReqResHandler[Req, Res, Err]) Operation(op operations.Operation) *ReqResHandler[Req, Res, Err] {
	h.op = op
	return h
}

func (h *ReqResHandler[Req, Res, Err]) Extensions(data operations.ExtensionData) *ReqResHandler[Req, Res, Err] {
	var err error
	h.op, err = operations.Extensions(data)(nil, h.op)
	requireValidSpec(err, h.method, h.pattern)
	return h
}

func (h *ReqResHandler[Req, Res, Err]) NoSecurity() *ReqResHandler[Req, Res, Err] {
	var err error
	h.op, err = operations.NoSecurity()(nil, h.op)
	requireValidSpec(err, h.method, h.pattern)
	return h
}

func (h *ReqResHandler[Req, Res, Err]) Security(name string, scopes ...string) *ReqResHandler[Req, Res, Err] {
	var err error
	h.op, err = operations.Security(name, scopes...)(nil, h.op)
	requireValidSpec(err, h.method, h.pattern)
	return h
}

func (h *ReqResHandler[Req, Res, Err]) ID(id string) *ReqResHandler[Req, Res, Err] {
	var err error
	h.op, err = operations.ID(id)(nil, h.op)
	requireValidSpec(err, h.method, h.pattern)
	return h
}

func (h *ReqResHandler[Req, Res, Err]) Tags(tags ...string) *ReqResHandler[Req, Res, Err] {
	var err error
	h.op, err = operations.Tags(tags...)(nil, h.op)
	requireValidSpec(err, h.method, h.pattern)
	return h
}

func (h *ReqResHandler[Req, Res, Err]) Deprecated() *ReqResHandler[Req, Res, Err] {
	var err error
	h.op, err = operations.Deprecated()(nil, h.op)
	requireValidSpec(err, h.method, h.pattern)
	return h
}

func (h *ReqResHandler[Req, Res, Err]) Summary(summary string) *ReqResHandler[Req, Res, Err] {
	var err error
	h.op, err = operations.Summary(summary)(nil, h.op)
	requireValidSpec(err, h.method, h.pattern)
	return h
}

func requireValidSpec(err error, method string, pattern string) {
	if err != nil {
		log.Fatal(fmt.Sprintf("invalid spec for handler %s %s: %v", method, pattern, err))
	}
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
	return h.fn
}

func (h *ReqResHandler[Req, Res, Err]) Op() operations.Operation {
	return h.op
}

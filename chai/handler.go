package chai

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
)

type HandlerFunc[Req any, Res any, Err ErrType] func(Req, http.ResponseWriter, *http.Request) (Res, int, Err)

type Handler[Req any, Res any, Err ErrType] struct {
	method     string
	pattern    string
	fn         HandlerFunc[Req, Res, Err]
	req        *Req
	res        *Res
	err        *Err
	op         *openapi3.Operation
	decodeFn   DecoderFunc[Req]
	validateFn ValidatorFunc[Req]
	respondFn  ResponderFunc[Res]
	errorFn    ErrorResponderFunc
}

func NewHandler[Req any, Res any, Err ErrType](method string, pattern string, fn HandlerFunc[Req, Res, Err]) *Handler[Req, Res, Err] {
	// TODO? panic if the Req type has a path param that is not specified in the pattern

	return &Handler[Req, Res, Err]{
		method:     method,
		pattern:    pattern,
		fn:         fn,
		decodeFn:   defaultDecoder[Req],
		validateFn: defaultValidator[Req],
		respondFn:  defaultResponder[Res],
		errorFn:    DefaultErrorResponder,
		op:         &openapi3.Operation{},
	}
}

func (h *Handler[Req, Res, Err]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := h.decodeFn(r)
	if handleErr(w, r, err, http.StatusBadRequest, h.errorFn) {
		return
	}
	err = h.validateFn(req)
	if handleErr(w, r, err, http.StatusBadRequest, h.errorFn) {
		return
	}
	// Note: err is of type error, while err2 is of type Err
	// the error check inside handleErr would be incorrect if we pass Err wrapped in error
	// due to how Go handles comparing nil values of different types
	res, code, err2 := h.fn(req, w, r)
	if handleErr(w, r, err2, code, h.errorFn) {
		return
	}
	h.respondFn(w, r, code, res)
}

func (h *Handler[Req, Res, Err]) WithDecoder(decodeFn DecoderFunc[Req]) *Handler[Req, Res, Err] {
	h.decodeFn = decodeFn
	return h
}

func (h *Handler[Req, Res, Err]) WithValidator(validateFn ValidatorFunc[Req]) *Handler[Req, Res, Err] {
	h.validateFn = validateFn
	return h
}

func (h *Handler[Req, Res, Err]) WithResponder(respondFn ResponderFunc[Res]) *Handler[Req, Res, Err] {
	h.respondFn = respondFn
	return h
}

func (h *Handler[Req, Res, Err]) WithErrorResponder(errorFn ErrorResponderFunc) *Handler[Req, Res, Err] {
	h.errorFn = errorFn
	return h
}

func (h *Handler[Req, Res, Err]) Extensions(data map[string]interface{}) *Handler[Req, Res, Err] {
	h.op.Extensions = data
	return h
}

func (h *Handler[Req, Res, Err]) NoSecurity() *Handler[Req, Res, Err] {
	h.op.Security = openapi3.NewSecurityRequirements()
	return h
}

func (h *Handler[Req, Res, Err]) Security(name string, scopes ...string) *Handler[Req, Res, Err] {
	if h.op.Security == nil {
		h.op.Security = openapi3.NewSecurityRequirements()
	}

	h.op.Security = h.op.Security.With(openapi3.
		NewSecurityRequirement().
		Authenticate(name, scopes...))
	return h
}

func (h *Handler[Req, Res, Err]) ID(id string) *Handler[Req, Res, Err] {
	h.op.OperationID = id
	return h
}

func (h *Handler[Req, Res, Err]) Tags(tags ...string) *Handler[Req, Res, Err] {
	h.op.Tags = tags
	return h
}

func (h *Handler[Req, Res, Err]) Deprecated() *Handler[Req, Res, Err] {
	h.op.Deprecated = true
	return h
}

func (h *Handler[Req, Res, Err]) Summary(summary string) *Handler[Req, Res, Err] {
	h.op.Summary = summary
	return h
}

func (h *Handler[Req, Res, Err]) Description(description string) *Handler[Req, Res, Err] {
	h.op.Description = description
	return h
}

func (h *Handler[Req, Res, Err]) AddResponse(code int, description string) *Handler[Req, Res, Err] {
	addResponse(h.op, code, openapi3.NewResponse().WithDescription(description))
	return h
}

func (h *Handler[Req, Res, Err]) ResponseCodes(description string, codes ...int) *Handler[Req, Res, Err] {
	for _, code := range codes {
		addResponse(h.op, code, openapi3.NewResponse().WithDescription(description))
	}
	return h
}

func (h *Handler[Req, Res, Err]) GetMetadata() *Metadata {
	return &Metadata{
		Req:            h.req,
		Res:            h.res,
		Err:            h.err,
		Op:             h.op,
		HandlerFunc:    h.fn,
		HandlerWrapper: h,
	}
}

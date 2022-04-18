package chai

import (
	"net/http"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
)

type ResHandlerFunc[Res any, Err ErrType] func(http.ResponseWriter, *http.Request) (Res, int, Err)

type ResHandler[Res any, Err ErrType] struct {
	method    string
	pattern   string
	fn        ResHandlerFunc[Res, Err]
	res       *Res
	err       *Err
	op        *openapi3.Operation
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
		op:        &openapi3.Operation{},
	}
}

func (h *ResHandler[Res, Err]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, code, err := h.fn(w, r)
	if handleErr(w, r, err, code, h.errorFn) {
		return
	}
	h.respondFn(w, r, code, res)
}

func (h *ResHandler[Res, Err]) WithResponder(respondFn ResponderFunc[Res]) *ResHandler[Res, Err] {
	h.respondFn = respondFn
	return h
}

func (h *ResHandler[Res, Err]) WithErrorResponder(errorFn ErrorResponderFunc) *ResHandler[Res, Err] {
	h.errorFn = errorFn
	return h
}

func (h *ResHandler[Res, Err]) Extensions(data map[string]interface{}) *ResHandler[Res, Err] {
	h.op.Extensions = data
	return h
}

func (h *ResHandler[Res, Err]) NoSecurity() *ResHandler[Res, Err] {
	h.op.Security = openapi3.NewSecurityRequirements()
	return h
}

func (h *ResHandler[Res, Err]) Security(name string, scopes ...string) *ResHandler[Res, Err] {
	if h.op.Security == nil {
		h.op.Security = openapi3.NewSecurityRequirements()
	}

	h.op.Security = h.op.Security.With(openapi3.
		NewSecurityRequirement().
		Authenticate(name, scopes...))
	return h
}

func (h *ResHandler[Res, Err]) ID(id string) *ResHandler[Res, Err] {
	h.op.OperationID = id
	return h
}

func (h *ResHandler[Res, Err]) Tags(tags ...string) *ResHandler[Res, Err] {
	h.op.Tags = tags
	return h
}

func (h *ResHandler[Res, Err]) Deprecated() *ResHandler[Res, Err] {
	h.op.Deprecated = true
	return h
}

func (h *ResHandler[Res, Err]) Summary(summary string) *ResHandler[Res, Err] {
	h.op.Summary = summary
	return h
}

func (h *ResHandler[Res, Err]) Description(description string) *ResHandler[Res, Err] {
	h.op.Description = description
	return h
}

func AddResponse(operation *openapi3.Operation, status int, response *openapi3.Response) {
	responses := operation.Responses
	if responses == nil {
		responses = make(openapi3.Responses)
		operation.Responses = responses
	}
	code := "default"
	if status != 0 {
		code = strconv.FormatInt(int64(status), 10)
	}
	responses[code] = &openapi3.ResponseRef{Value: response}
}

func (h *ResHandler[Res, Err]) AddResponse(code int, description string) *ResHandler[Res, Err] {
	AddResponse(h.op, code, openapi3.NewResponse().WithDescription(description))
	return h
}

func (h *ResHandler[Res, Err]) ResponseCodes(description string, codes ...int) *ResHandler[Res, Err] {
	for _, code := range codes {
		AddResponse(h.op, code, openapi3.NewResponse().WithDescription(description))
	}
	return h
}

type Metadata struct {
	Req            any
	Res            any
	Err            any
	Op             *openapi3.Operation
	HandlerFunc    any
	HandlerWrapper http.Handler
}

func (h *ResHandler[Res, Err]) GetMetadata() *Metadata {
	return &Metadata{
		Req:            nil,
		Res:            h.res,
		Err:            h.err,
		Op:             h.op,
		HandlerFunc:    h.fn,
		HandlerWrapper: h,
	}
}

package chai

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
)

type BasicHandlerFunc func(http.ResponseWriter, *http.Request)

type BasicHandler struct {
	method  string
	pattern string
	fn      BasicHandlerFunc
	req     any
	res     any
	err     any
	op      *openapi3.Operation
}

func NewBasicHandler(method string, pattern string, fn BasicHandlerFunc) *BasicHandler {
	return &BasicHandler{
		method:  method,
		pattern: pattern,
		fn:      fn,
		op:      &openapi3.Operation{},
	}
}

func (h *BasicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.fn(w, r)
}

func (h *BasicHandler) Extensions(data map[string]interface{}) *BasicHandler {
	h.op.Extensions = data
	return h
}

func (h *BasicHandler) NoSecurity() *BasicHandler {
	h.op.Security = openapi3.NewSecurityRequirements()
	return h
}

func (h *BasicHandler) Security(name string, scopes ...string) *BasicHandler {
	if h.op.Security == nil {
		h.op.Security = openapi3.NewSecurityRequirements()
	}

	h.op.Security = h.op.Security.With(openapi3.
		NewSecurityRequirement().
		Authenticate(name, scopes...))
	return h
}

func (h *BasicHandler) ID(id string) *BasicHandler {
	h.op.OperationID = id
	return h
}

func (h *BasicHandler) Tags(tags ...string) *BasicHandler {
	h.op.Tags = tags
	return h
}

func (h *BasicHandler) Deprecated() *BasicHandler {
	h.op.Deprecated = true
	return h
}

func (h *BasicHandler) Summary(summary string) *BasicHandler {
	h.op.Summary = summary
	return h
}

func (h *BasicHandler) Description(description string) *BasicHandler {
	h.op.Description = description
	return h
}

func (h *BasicHandler) RequestBodyType(req any) *BasicHandler {
	h.req = req
	return h
}

func (h *BasicHandler) ResponseType(res any) *BasicHandler {
	h.res = res
	return h
}

func (h *BasicHandler) AddResponse(res any, code int, description string) *BasicHandler {
	AddResponse(h.op, code, openapi3.NewResponse().WithDescription(description))
	h.res = res
	return h
}

func (h *BasicHandler) Operation(op *openapi3.Operation) *BasicHandler {
	h.op = op
	return h
}

func (h *BasicHandler) GetMetadata() *Metadata {
	return &Metadata{
		Req:            h.req,
		Res:            h.res,
		Err:            h.err,
		Op:             h.op,
		HandlerFunc:    h.fn,
		HandlerWrapper: h,
	}
}

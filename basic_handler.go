package chai

import (
	"go/ast"
	"net/http"
)

func NewBasicHandler(h http.HandlerFunc) *BasicHandler {
	return &BasicHandler{
		f: h,
	}
}

type BasicHandler struct {
	f       http.HandlerFunc
	comment string
	astFile *ast.File
}

func (h *BasicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.f.ServeHTTP(w, r)
}

func (h *BasicHandler) Handler() any {
	return h.f
}

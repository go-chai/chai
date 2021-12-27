package chai

import (
	"go/ast"
	"net/http"
)

func NewBasicHandler(h http.HandlerFunc) *BasicHandler {
	return &BasicHandler{
		f:       h,
		comment: GetFuncInfo(h).Comment,
		astFile: GetFuncInfo(h).ASTFile,
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

func (h *BasicHandler) Comment() string {
	return h.comment
}

func (h *BasicHandler) ASTFile() *ast.File {
	return h.astFile
}

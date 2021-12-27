package chai

import (
	"go/ast"
	"net/http"

	"errors"
)

type ResHandlerFunc[Res any, Err ErrType] func(http.ResponseWriter, *http.Request) (Res, int, Err)

func NewResHandler[Res any, Err ErrType](h ResHandlerFunc[Res, Err]) *ResHandler[Res, Err] {
	return &ResHandler[Res, Err]{
		f:       h,
		res:     new(Res),
		err:     new(Err),
		comment: GetFuncInfo(h).Comment,
		astFile: GetFuncInfo(h).ASTFile,
	}
}

type ResHandler[Res any, Err ErrType] struct {
	f       ResHandlerFunc[Res, Err]
	res     *Res
	err     *Err
	comment string
	astFile *ast.File
}

func (h *ResHandler[Res, Err]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, code, err := h.f(w, r)
	if !errors.Is(err, nil) {
		if code == 0 {
			code = http.StatusInternalServerError
		}

		write(w, code, err)
		return
	}

	if code == 0 {
		code = http.StatusOK
	}

	write(w, code, res)
}

func (h *ResHandler[Res, Err]) Res() any {
	return h.res
}

func (h *ResHandler[Res, Err]) Err() any {
	return h.err
}

func (h *ResHandler[Res, Err]) Comment() string {
	return h.comment
}
func (h *ResHandler[Res, Err]) ASTFile() *ast.File {
	return h.astFile
}

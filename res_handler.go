package chai

import (
	"go/ast"
	"net/http"

	"errors"
)

type resHandlerFunc[Res any, Err ErrType] func(http.ResponseWriter, *http.Request) (Res, int, Err)

func newResHandlerFunc[Res any, Err ErrType](h resHandlerFunc[Res, Err]) *resHandler[Res, Err] {
	return &resHandler[Res, Err]{
		f:       h,
		res:     new(Res),
		err:     new(Err),
		comment: GetFuncInfo(h).Comment,
		astFile: GetFuncInfo(h).ASTFile,
	}
}

type resHandler[Res any, Err ErrType] struct {
	f       resHandlerFunc[Res, Err]
	res     *Res
	err     *Err
	comment string
	astFile *ast.File
}

func (h *resHandler[Res, Err]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, code, err := h.f(w, r)
	if !errors.Is(err, nil) {
		if code == 0 {
			code = http.StatusInternalServerError
		}

		write(w, code, JSONError{Message: err.Error()})
		return
	}

	if code == 0 {
		code = http.StatusOK
	}

	write(w, code, res)
}

func (h *resHandler[Res, Err]) Res() any {
	return h.res
}

func (h *resHandler[Res, Err]) Err() any {
	return h.err
}

func (h *resHandler[Res, Err]) Comment() string {
	return h.comment
}
func (h *resHandler[Res, Err]) ASTFile() *ast.File {
	return h.astFile
}

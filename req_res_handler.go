package chai

import (
	"encoding/json"
	"errors"
	"go/ast"
	"net/http"
)

type reqResHandlerFunc[Req any, Res any, Err ErrType] func(Req, http.ResponseWriter, *http.Request) (Res, int, Err)

func newReqResHandlerFunc[Req any, Res any, Err ErrType](h reqResHandlerFunc[Req, Res, Err]) *reqResHandler[Req, Res, Err] {
	return &reqResHandler[Req, Res, Err]{
		f:       h,
		req:     new(Req),
		res:     new(Res),
		comment: GetFuncInfo(h).Comment,
		astFile: GetFuncInfo(h).ASTFile,
	}
}

type reqResHandler[Req any, Res any, Err ErrType] struct {
	f       reqResHandlerFunc[Req, Res, Err]
	req     *Req
	res     *Res
	comment string
	astFile *ast.File
}

func (h *reqResHandler[Req, Res, Err]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req *Req

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		write(w, http.StatusBadRequest, DefaultFromErrorer.FromError(err))
		return
	}

	res, code, err := h.f(*req, w, r)

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

func (h *reqResHandler[Req, Res, Err]) Req() any {
	return h.req
}

func (h *reqResHandler[Req, Res, Err]) Res() any {
	return h.res
}

func (h *reqResHandler[Req, Res, Err]) Comment() string {
	return h.comment
}

func (h *reqResHandler[Req, Res, Err]) ASTFile() *ast.File {
	return h.astFile
}

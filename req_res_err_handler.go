package chai

import (
	"encoding/json"
	"go/ast"
	"net/http"
)

type reqResErrHandlerFunc[Req any, Res any, Err ErrType] func(*Request[Req]) (*Response[Res], *Error[Err])

func newReqResErrHandlerFunc[Req any, Res any, Err ErrType](h reqResErrHandlerFunc[Req, Res, Err]) *reqResErrHandler[Req, Res, Err] {
	return &reqResErrHandler[Req, Res, Err]{
		f:       h,
		req:     new(Req),
		res:     new(Res),
		comment: GetFuncInfo(h).Comment,
		astFile: GetFuncInfo(h).ASTFile,
	}
}

type reqResErrHandler[Req any, Res any, Err ErrType] struct {
	f       reqResErrHandlerFunc[Req, Res, Err]
	req     *Req
	res     *Res
	err     error
	comment string
	astFile *ast.File
}

type FromErrorer interface {
	FromError(error) any
}

type defaultFromErrorer struct{}

func (defaultFromErrorer) FromError(err error) any {
	return &APIError{Message: err.Error()}
}

var DefaultFromErrorer = &defaultFromErrorer{}

func (h *reqResErrHandler[Req, Res, Err]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req *Request[Req]

	if err := json.NewDecoder(r.Body).Decode(&req.Req); err != nil {
		write(w, http.StatusBadRequest, DefaultFromErrorer.FromError(err))
		return
	}

	res, err2 := h.f(req)
	if err2 != nil {
		if err2.Code == 0 {
			err2.Code = http.StatusInternalServerError
		}

		write(w, err2.Code, err2.Error)
		return
	}

	if res.Code == 0 {
		res.Code = http.StatusOK
	}

	write(w, res.Code, res.Res)
}

func (h *reqResErrHandler[Req, Res, Err]) Req() any {
	return h.req
}

func (h *reqResErrHandler[Req, Res, Err]) Res() any {
	return h.res
}

func (h *reqResErrHandler[Req, Res, Err]) Err() any {
	return h.err
}

func (h *reqResErrHandler[Req, Res, Err]) Comment() string {
	return h.comment
}

func (h *reqResErrHandler[Req, Res, Err]) ASTFile() *ast.File {
	return h.astFile
}

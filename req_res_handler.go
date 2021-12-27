package chai

import (
	"encoding/json"
	"errors"
	"go/ast"
	"net/http"
	"reflect"
)

type ReqResHandlerFunc[Req any, Res any, Err ErrType] func(Req, http.ResponseWriter, *http.Request) (Res, int, Err)

func NewReqResHandler[Req any, Res any, Err ErrType](h ReqResHandlerFunc[Req, Res, Err]) *ReqResHandler[Req, Res, Err] {
	return &ReqResHandler[Req, Res, Err]{
		f:       h,
		req:     new(Req),
		res:     new(Res),
		err:     new(Err),
		ht:      reflect.TypeOf(h),
		comment: GetFuncInfo(h).Comment,
		astFile: GetFuncInfo(h).ASTFile,
	}
}

type ReqResHandler[Req any, Res any, Err ErrType] struct {
	f       ReqResHandlerFunc[Req, Res, Err]
	req     *Req
	res     *Res
	err     *Err
	ht      reflect.Type
	comment string
	astFile *ast.File
}

func (h *ReqResHandler[Req, Res, Err]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (h *ReqResHandler[Req, Res, Err]) Req() any {
	return h.req
}

func (h *ReqResHandler[Req, Res, Err]) Res() any {
	return h.res
}

func (h *ReqResHandler[Req, Res, Err]) Err() any {
	return h.err
}

func (h *ReqResHandler[Req, Res, Err]) HT() reflect.Type {
	return h.ht
}

func (h *ReqResHandler[Req, Res, Err]) Comment() string {
	return h.comment
}

func (h *ReqResHandler[Req, Res, Err]) ASTFile() *ast.File {
	return h.astFile
}

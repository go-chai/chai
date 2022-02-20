package chai

import (
	"encoding/json"
	"net/http"
	"reflect"
)

type ReqResHandlerFunc[Req any, Res any, Err ErrType] func(Req, http.ResponseWriter, *http.Request) (Res, int, Err)

func NewReqResHandler[Req any, Res any, Err ErrType](h ReqResHandlerFunc[Req, Res, Err]) *ReqResHandler[Req, Res, Err] {
	return &ReqResHandler[Req, Res, Err]{
		f: h,
	}
}

type ReqResHandler[Req any, Res any, Err ErrType] struct {
	f   ReqResHandlerFunc[Req, Res, Err]
	req *Req
	res *Res
	err *Err
}

func isErr[Err ErrType](err Err) bool {
	return !reflect.ValueOf(&err).Elem().IsZero()
}

func (h *ReqResHandler[Req, Res, Err]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req *Req

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}

	res, code, err := h.f(*req, w, r)
	if isErr(err) {
		if code == 0 {
			code = http.StatusInternalServerError
		}

		writeErr(w, code, err)
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

func (h *ReqResHandler[Req, Res, Err]) Handler() any {
	return h.f
}

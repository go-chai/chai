package chai

import (
	"net/http"
)

func Get[Res any, Err ErrType](r Methoder, path string, fn ResHandlerFunc[Res, Err]) {
	r.Method(http.MethodGet, path, NewResHandler(fn))
}

func Connect[Res any, Err ErrType](r Methoder, path string, fn ResHandlerFunc[Res, Err]) {
	r.Method(http.MethodConnect, path, NewResHandler(fn))
}

func Options[Res any, Err ErrType](r Methoder, path string, fn ResHandlerFunc[Res, Err]) {
	r.Method(http.MethodOptions, path, NewResHandler(fn))
}

func Post[Req any, Res any, Err ErrType](r Methoder, path string, fn ReqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodPost, path, NewReqResHandler(fn))
}

func Put[Req any, Res any, Err ErrType](r Methoder, path string, fn ReqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodPut, path, NewReqResHandler(fn))
}

func Patch[Req any, Res any, Err ErrType](r Methoder, path string, fn ReqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodPatch, path, NewReqResHandler(fn))
}

func Delete[Req any, Res any, Err ErrType](r Methoder, path string, fn ReqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodDelete, path, NewReqResHandler(fn))
}

package chai

import (
	"net/http"
)

func Get[Res any, Err ErrType](r Methoder, path string, fn ResHandlerFunc[Res, Err]) {
	r.Method(http.MethodGet, path, NewResHandlerFunc(fn))
}

func Connect[Res any, Err ErrType](r Methoder, path string, fn ResHandlerFunc[Res, Err]) {
	r.Method(http.MethodConnect, path, NewResHandlerFunc(fn))
}

func Options[Res any, Err ErrType](r Methoder, path string, fn ResHandlerFunc[Res, Err]) {
	r.Method(http.MethodOptions, path, NewResHandlerFunc(fn))
}

func Post[Req any, Res any, Err ErrType](r Methoder, path string, fn ReqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodPost, path, NewReqResHandlerFunc(fn))
}

func Put[Req any, Res any, Err ErrType](r Methoder, path string, fn ReqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodPut, path, NewReqResHandlerFunc(fn))
}

func Patch[Req any, Res any, Err ErrType](r Methoder, path string, fn ReqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodPatch, path, NewReqResHandlerFunc(fn))
}

func Delete[Req any, Res any, Err ErrType](r Methoder, path string, fn ReqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodDelete, path, NewReqResHandlerFunc(fn))
}

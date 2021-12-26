package chai

import (
	"net/http"
)

func Get[Res any, Err ErrType](r Methoder, path string, fn resHandlerFunc[Res, Err]) {
	r.Method(http.MethodGet, path, newResHandlerFunc(fn))
}

func Connect[Res any, Err ErrType](r Methoder, path string, fn resHandlerFunc[Res, Err]) {
	r.Method(http.MethodConnect, path, newResHandlerFunc(fn))
}

func Options[Res any, Err ErrType](r Methoder, path string, fn resHandlerFunc[Res, Err]) {
	r.Method(http.MethodOptions, path, newResHandlerFunc(fn))
}

func Post[Req any, Res any, Err ErrType](r Methoder, path string, fn reqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodPost, path, newReqResHandlerFunc(fn))
}

func Put[Req any, Res any, Err ErrType](r Methoder, path string, fn reqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodPut, path, newReqResHandlerFunc(fn))
}

func Patch[Req any, Res any, Err ErrType](r Methoder, path string, fn reqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodPatch, path, newReqResHandlerFunc(fn))
}

func Delete[Req any, Res any, Err ErrType](r Methoder, path string, fn reqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodDelete, path, newReqResHandlerFunc(fn))
}

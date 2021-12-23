package chai

import (
	"net/http"
)

func Get[Res any](r Methoder, path string, fn resHandlerFunc[Res]) {
	r.Method(http.MethodGet, path, newResHandlerFunc(fn))
}

func Post[Req any, Res any, Err ErrType](r Methoder, path string, fn reqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodPost, path, newReqResHandlerFunc(fn))
}

func Post2[Req any, Res any, Err ErrType](r Methoder, path string, fn reqResErrHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodPost, path, newReqResErrHandlerFunc(fn))
}

func Put[Req any, Res any, Err ErrType](r Methoder, path string, fn reqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodPut, path, newReqResHandlerFunc(fn))
}

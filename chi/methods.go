package chai

import (
	"net/http"

	"github.com/go-chai/chai/chai"
)

func Get[Res any, Err chai.ErrType](r chai.Methoder, pattern string, fn chai.ResHandlerFunc[Res, Err]) *chai.ResHandler[Res, Err] {
	method := http.MethodGet
	h := chai.NewResHandler(method, pattern, fn)
	r.Method(method, pattern, h)
	return h
}

func Connect[Res any, Err chai.ErrType](r chai.Methoder, pattern string, fn chai.ResHandlerFunc[Res, Err]) *chai.ResHandler[Res, Err] {
	method := http.MethodConnect
	h := chai.NewResHandler(method, pattern, fn)
	r.Method(method, pattern, h)
	return h
}

func Options[Res any, Err chai.ErrType](r chai.Methoder, pattern string, fn chai.ResHandlerFunc[Res, Err]) *chai.ResHandler[Res, Err] {
	method := http.MethodOptions
	h := chai.NewResHandler(method, pattern, fn)
	r.Method(method, pattern, h)
	return h
}

func Post[Req any, Res any, Err chai.ErrType](r chai.Methoder, pattern string, fn chai.ReqResHandlerFunc[Req, Res, Err]) *chai.ReqResHandler[Req, Res, Err] {
	method := http.MethodPost
	h := chai.NewReqResHandler(method, pattern, fn)
	r.Method(method, pattern, h)
	return h
}

func Put[Req any, Res any, Err chai.ErrType](r chai.Methoder, pattern string, fn chai.ReqResHandlerFunc[Req, Res, Err]) *chai.ReqResHandler[Req, Res, Err] {
	method := http.MethodPut
	h := chai.NewReqResHandler(method, pattern, fn)
	r.Method(method, pattern, h)
	return h
}

func Patch[Req any, Res any, Err chai.ErrType](r chai.Methoder, pattern string, fn chai.ReqResHandlerFunc[Req, Res, Err]) *chai.ReqResHandler[Req, Res, Err] {
	method := http.MethodPatch
	h := chai.NewReqResHandler(method, pattern, fn)
	r.Method(method, pattern, h)
	return h
}

func Delete[Req any, Res any, Err chai.ErrType](r chai.Methoder, pattern string, fn chai.ReqResHandlerFunc[Req, Res, Err]) *chai.ReqResHandler[Req, Res, Err] {
	method := http.MethodDelete
	h := chai.NewReqResHandler(method, pattern, fn)
	r.Method(method, pattern, h)
	return h
}

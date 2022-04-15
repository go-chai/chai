package chai

import (
	"net/http"

	"github.com/go-chai/chai/chai"
)

func Get[Res any, Err chai.ErrType](r chai.Methoder, path string, fn chai.ResHandlerFunc[Res, Err]) *chai.ResHandler[Res, Err] {
	h := chai.NewResHandler(fn)
	r.Method(http.MethodGet, path, h)
	return h
}

func Connect[Res any, Err chai.ErrType](r chai.Methoder, path string, fn chai.ResHandlerFunc[Res, Err]) *chai.ResHandler[Res, Err] {
	h := chai.NewResHandler(fn)
	r.Method(http.MethodConnect, path, h)
	return h
}

func Options[Res any, Err chai.ErrType](r chai.Methoder, path string, fn chai.ResHandlerFunc[Res, Err]) *chai.ResHandler[Res, Err] {
	h := chai.NewResHandler(fn)
	r.Method(http.MethodOptions, path, h)
	return h
}

func Post[Req any, Res any, Err chai.ErrType](r chai.Methoder, path string, fn chai.ReqResHandlerFunc[Req, Res, Err]) *chai.ReqResHandler[Req, Res, Err] {
	h := chai.NewReqResHandler(fn)
	r.Method(http.MethodPost, path, h)
	return h
}

func Put[Req any, Res any, Err chai.ErrType](r chai.Methoder, path string, fn chai.ReqResHandlerFunc[Req, Res, Err]) *chai.ReqResHandler[Req, Res, Err] {
	h := chai.NewReqResHandler(fn)
	r.Method(http.MethodPut, path, h)
	return h
}

func Patch[Req any, Res any, Err chai.ErrType](r chai.Methoder, path string, fn chai.ReqResHandlerFunc[Req, Res, Err]) *chai.ReqResHandler[Req, Res, Err] {
	h := chai.NewReqResHandler(fn)
	r.Method(http.MethodPatch, path, h)
	return h
}

func Delete[Req any, Res any, Err chai.ErrType](r chai.Methoder, path string, fn chai.ReqResHandlerFunc[Req, Res, Err]) *chai.ReqResHandler[Req, Res, Err] {
	h := chai.NewReqResHandler(fn)
	r.Method(http.MethodDelete, path, h)
	return h
}

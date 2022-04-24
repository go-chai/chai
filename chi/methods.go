package chai

import (
	"net/http"

	"github.com/go-chai/chai/chai"
)

func Get[Req any, Res any, Err chai.ErrType](r chai.Methoder, pattern string, fn chai.HandlerFunc[Req, Res, Err]) *chai.Handler[Req, Res, Err] {
	method := http.MethodGet
	h := chai.NewHandler(method, pattern, fn)
	r.Method(method, pattern, h)
	return h
}

func Connect[Req any, Res any, Err chai.ErrType](r chai.Methoder, pattern string, fn chai.HandlerFunc[Req, Res, Err]) *chai.Handler[Req, Res, Err] {
	method := http.MethodConnect
	h := chai.NewHandler(method, pattern, fn)
	r.Method(method, pattern, h)
	return h
}

func Options[Req any, Res any, Err chai.ErrType](r chai.Methoder, pattern string, fn chai.HandlerFunc[Req, Res, Err]) *chai.Handler[Req, Res, Err] {
	method := http.MethodOptions
	h := chai.NewHandler(method, pattern, fn)
	r.Method(method, pattern, h)
	return h
}

func Post[Req any, Res any, Err chai.ErrType](r chai.Methoder, pattern string, fn chai.HandlerFunc[Req, Res, Err]) *chai.Handler[Req, Res, Err] {
	method := http.MethodPost
	h := chai.NewHandler(method, pattern, fn)
	r.Method(method, pattern, h)
	return h
}

func Put[Req any, Res any, Err chai.ErrType](r chai.Methoder, pattern string, fn chai.HandlerFunc[Req, Res, Err]) *chai.Handler[Req, Res, Err] {
	method := http.MethodPut
	h := chai.NewHandler(method, pattern, fn)
	r.Method(method, pattern, h)
	return h
}

func Patch[Req any, Res any, Err chai.ErrType](r chai.Methoder, pattern string, fn chai.HandlerFunc[Req, Res, Err]) *chai.Handler[Req, Res, Err] {
	method := http.MethodPatch
	h := chai.NewHandler(method, pattern, fn)
	r.Method(method, pattern, h)
	return h
}

func Delete[Req any, Res any, Err chai.ErrType](r chai.Methoder, pattern string, fn chai.HandlerFunc[Req, Res, Err]) *chai.Handler[Req, Res, Err] {
	method := http.MethodDelete
	h := chai.NewHandler(method, pattern, fn)
	r.Method(method, pattern, h)
	return h
}

func GetB(r chai.Methoder, pattern string, fn chai.BasicHandlerFunc) *chai.BasicHandler {
	method := http.MethodGet
	h := chai.NewBasicHandler(method, pattern, fn)
	r.Method(method, pattern, h)
	return h
}

func ConnectB(r chai.Methoder, pattern string, fn chai.BasicHandlerFunc) *chai.BasicHandler {
	method := http.MethodConnect
	h := chai.NewBasicHandler(method, pattern, fn)
	r.Method(method, pattern, h)
	return h
}

func OptionsB(r chai.Methoder, pattern string, fn chai.BasicHandlerFunc) *chai.BasicHandler {
	method := http.MethodOptions
	h := chai.NewBasicHandler(method, pattern, fn)
	r.Method(method, pattern, h)
	return h
}

func PostB(r chai.Methoder, pattern string, fn chai.BasicHandlerFunc) *chai.BasicHandler {
	method := http.MethodPost
	h := chai.NewBasicHandler(method, pattern, fn)
	r.Method(method, pattern, h)
	return h
}

func PutB(r chai.Methoder, pattern string, fn chai.BasicHandlerFunc) *chai.BasicHandler {
	method := http.MethodPut
	h := chai.NewBasicHandler(method, pattern, fn)
	r.Method(method, pattern, h)
	return h
}

func PatchB(r chai.Methoder, pattern string, fn chai.BasicHandlerFunc) *chai.BasicHandler {
	method := http.MethodPatch
	h := chai.NewBasicHandler(method, pattern, fn)
	r.Method(method, pattern, h)
	return h
}

func DeleteB(r chai.Methoder, pattern string, fn chai.BasicHandlerFunc) *chai.BasicHandler {
	method := http.MethodDelete
	h := chai.NewBasicHandler(method, pattern, fn)
	r.Method(method, pattern, h)
	return h
}

package chai

import (
	"net/http"

	"github.com/go-chai/chai/chai"
)

func Get[Res any, Err chai.ErrType](r chai.Methoder, path string, fn chai.ResHandlerFunc[Res, Err]) {
	r.Method(http.MethodGet, path, chai.NewResHandler(fn))
}

func Connect[Res any, Err chai.ErrType](r chai.Methoder, path string, fn chai.ResHandlerFunc[Res, Err]) {
	r.Method(http.MethodConnect, path, chai.NewResHandler(fn))
}

func Options[Res any, Err chai.ErrType](r chai.Methoder, path string, fn chai.ResHandlerFunc[Res, Err]) {
	r.Method(http.MethodOptions, path, chai.NewResHandler(fn))
}

func Post[Req any, Res any, Err chai.ErrType](r chai.Methoder, path string, fn chai.ReqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodPost, path, chai.NewReqResHandler(fn))
}

func Put[Req any, Res any, Err chai.ErrType](r chai.Methoder, path string, fn chai.ReqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodPut, path, chai.NewReqResHandler(fn))
}

func Patch[Req any, Res any, Err chai.ErrType](r chai.Methoder, path string, fn chai.ReqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodPatch, path, chai.NewReqResHandler(fn))
}

func Delete[Req any, Res any, Err chai.ErrType](r chai.Methoder, path string, fn chai.ReqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodDelete, path, chai.NewReqResHandler(fn))
}

package chai

import (
	"net/http"

	"github.com/go-chai/chai/chai"
	"github.com/gorilla/mux"
)

func Get[Res any, Err chai.ErrType](r *mux.Router, path string, fn chai.ResHandlerFunc[Res, Err]) {
	r.Methods(http.MethodGet).Path(path).Handler(chai.NewResHandler(fn))
}

func Connect[Res any, Err chai.ErrType](r *mux.Router, path string, fn chai.ResHandlerFunc[Res, Err]) {
	r.Methods(http.MethodConnect).Path(path).Handler(chai.NewResHandler(fn))
}

func Options[Res any, Err chai.ErrType](r *mux.Router, path string, fn chai.ResHandlerFunc[Res, Err]) {
	r.Methods(http.MethodOptions).Path(path).Handler(chai.NewResHandler(fn))
}

func Post[Req any, Res any, Err chai.ErrType](r *mux.Router, path string, fn chai.ReqResHandlerFunc[Req, Res, Err]) {
	r.Methods(http.MethodPost).Path(path).Handler(chai.NewReqResHandler(fn))
}
func Put[Req any, Res any, Err chai.ErrType](r *mux.Router, path string, fn chai.ReqResHandlerFunc[Req, Res, Err]) {
	r.Methods(http.MethodPut).Path(path).Handler(chai.NewReqResHandler(fn))
}

func Patch[Req any, Res any, Err chai.ErrType](r *mux.Router, path string, fn chai.ReqResHandlerFunc[Req, Res, Err]) {
	r.Methods(http.MethodPatch).Path(path).Handler(chai.NewReqResHandler(fn))
}

func Delete[Req any, Res any, Err chai.ErrType](r *mux.Router, path string, fn chai.ReqResHandlerFunc[Req, Res, Err]) {
	r.Methods(http.MethodDelete).Path(path).Handler(chai.NewReqResHandler(fn))
}

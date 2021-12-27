package chai

import (
	"net/http"
)

func GetG[Res any, Err ErrType](r Methoder, path string, fn ResHandlerFunc[Res, Err]) {
	r.Method(http.MethodGet, path, NewResHandler(fn))
}

func ConnectG[Res any, Err ErrType](r Methoder, path string, fn ResHandlerFunc[Res, Err]) {
	r.Method(http.MethodConnect, path, NewResHandler(fn))
}

func OptionsG[Res any, Err ErrType](r Methoder, path string, fn ResHandlerFunc[Res, Err]) {
	r.Method(http.MethodOptions, path, NewResHandler(fn))
}

func PostG[Req any, Res any, Err ErrType](r Methoder, path string, fn ReqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodPost, path, NewReqResHandler(fn))
}

func PutG[Req any, Res any, Err ErrType](r Methoder, path string, fn ReqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodPut, path, NewReqResHandler(fn))
}

func PatchG[Req any, Res any, Err ErrType](r Methoder, path string, fn ReqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodPatch, path, NewReqResHandler(fn))
}

func DeleteG[Req any, Res any, Err ErrType](r Methoder, path string, fn ReqResHandlerFunc[Req, Res, Err]) {
	r.Method(http.MethodDelete, path, NewReqResHandler(fn))
}

func Get(r Methoder, path string, fn http.HandlerFunc) {
	r.Method(http.MethodGet, path, NewBasicHandler(fn))
}

func Connect(r Methoder, path string, fn http.HandlerFunc) {
	r.Method(http.MethodConnect, path, NewBasicHandler(fn))
}

func Options(r Methoder, path string, fn http.HandlerFunc) {
	r.Method(http.MethodOptions, path, NewBasicHandler(fn))
}

func Post(r Methoder, path string, fn http.HandlerFunc) {
	r.Method(http.MethodPost, path, NewBasicHandler(fn))
}

func Put(r Methoder, path string, fn http.HandlerFunc) {
	r.Method(http.MethodPut, path, NewBasicHandler(fn))
}

func Patch(r Methoder, path string, fn http.HandlerFunc) {
	r.Method(http.MethodPatch, path, NewBasicHandler(fn))
}

func Delete(r Methoder, path string, fn http.HandlerFunc) {
	r.Method(http.MethodDelete, path, NewBasicHandler(fn))
}

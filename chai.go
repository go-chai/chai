package chai

import (
	"encoding/json"
	"net/http"
)

type Methoder interface {
	Method(method, pattern string, h http.Handler)
}

type Reqer interface {
	Req() any
}

type ResErrer interface {
	Res() any
	Err() any
}

type Handlerer interface {
	Handler() any
}

type Error struct {
	Message          string `json:"error"`
	ErrorDebug       string `json:"error_debug,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
	StatusCode       int    `json:"status_code,omitempty"`
}

func (e Error) Error() string {
	return e.Message
}

func write(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func writeBytes(w http.ResponseWriter, code int, bytes []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(bytes)
}

type ErrType = error

func PtrTo[T any](t T) *T {
	return &t
}

type FromErrorer interface {
	FromError(error) any
}

type defaultFromErrorer struct{}

func (defaultFromErrorer) FromError(err error) any {
	return &Error{Message: err.Error(), StatusCode: http.StatusBadRequest}
}

var DefaultFromErrorer = &defaultFromErrorer{}

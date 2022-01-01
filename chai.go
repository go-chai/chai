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

func writeErr(w http.ResponseWriter, code int, e ErrType) {
	DefaultErrorWriter.WriteError(w, code, e)
}

func writeBytes(w http.ResponseWriter, code int, bytes []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(bytes)
}

type ErrType = error

type ErrorWriter interface {
	WriteError(w http.ResponseWriter, code int, e ErrType)
}

type defaultErrorWriter struct{}

func (defaultErrorWriter) WriteError(w http.ResponseWriter, code int, e ErrType) {
	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	if string(b) == "{}" {
		b, err = json.Marshal(&Error{Message: e.Error(), StatusCode: code})
		if err != nil {
			panic(err)
		}
	}

	writeBytes(w, code, b)
}

var DefaultErrorWriter = &defaultErrorWriter{}

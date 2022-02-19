package chai

import (
	"net/http"

	"encoding/json"
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

type Err error

type ErrWrap struct {
	Err        error
	StatusCode int
	Error      string
}

// TODO figure out how to do this without multiple json.Marshal/Unmarshal calls
func (ew *ErrWrap) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"error":       ew.Error,
		"status_code": ew.StatusCode,
	}

	b, err := json.Marshal(ew.Err)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	return json.Marshal(m)
}

type ErrorWriter interface {
	WriteError(w http.ResponseWriter, code int, e ErrType)
}

type defaultErrorWriter struct{}

func (defaultErrorWriter) WriteError(w http.ResponseWriter, code int, e ErrType) {
	ew := &ErrWrap{
		Err:        e,
		StatusCode: code,
		Error:      e.Error(),
	}

	b, err := json.Marshal(ew)
	if err != nil {
		panic(err)
	}

	writeBytes(w, code, b)
}

var DefaultErrorWriter = &defaultErrorWriter{}

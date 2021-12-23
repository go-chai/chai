package httputil

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

// NewError example
func NewError(w http.ResponseWriter, code int, err error) bool {
	if err == nil {
		return false
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)

	err = enc.Encode(HTTPError{Code: code, Message: err.Error()})

	if err != nil {
		panic(err) // If this happens, it's a programmer mistake so we panic
	}

	return true
}

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

func Respond(w http.ResponseWriter, r *http.Request, v interface{}) {
	JSON(w, r, v)
}

func JSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	bytes, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if status, ok := r.Context().Value(render.StatusCtxKey).(int); ok {
		w.WriteHeader(status)
	}
	w.Write(bytes)
}

func Decode(r *http.Request, v interface{}) error {
	return DecodeJSON(r.Body, v)
}

func DecodeFormValue(r *http.Request, key string, v interface{}) error {
	err := json.Unmarshal([]byte(r.FormValue(key)), v)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

func DecodeJSON(r io.Reader, v interface{}) error {
	defer io.Copy(ioutil.Discard, r)
	return json.NewDecoder(r).Decode(v)
}

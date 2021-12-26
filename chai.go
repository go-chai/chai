package chai

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"net/http"
	"reflect"

	"github.com/ghodss/yaml"
)

type Methoder interface {
	Method(method, pattern string, h http.Handler)
}

type Reqer interface {
	Req() any
}

type Reser interface {
	Res() any
}

type Errer interface {
	Err() any
}

type Commenter interface {
	Comment() string
}

type ASTFiler interface {
	ASTFile() *ast.File
}

type HTer interface {
	HT() reflect.Type
}

type JSONError struct {
	Message          string `json:"error"`
	ErrorDebug       string `json:"error_debug,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
	StatusCode       int    `json:"status_code,omitempty"`
}

func (e JSONError) Error() string {
	return e.Message
}

func writeBytes(w http.ResponseWriter, code int, bytes []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(bytes)
}

func write(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

type ErrType = error

func PtrTo[T any](t T) *T {
	return &t
}

func LogYAML(v any) {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

	return
}

func LogJSON(v any) {
	bytes, err := json.MarshalIndent(v, "", "  ")

	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

	return
}

type FromErrorer interface {
	FromError(error) any
}

type defaultFromErrorer struct{}

func (defaultFromErrorer) FromError(err error) any {
	return &JSONError{Message: err.Error()}
}

var DefaultFromErrorer = &defaultFromErrorer{}

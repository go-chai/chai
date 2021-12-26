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

type APIError struct {
	Message          string `json:"error"`
	ErrorDebug       string `json:"error_debug"`
	ErrorDescription string `json:"error_description"`
	StatusCode       int    `json:"status_code"`
}

func (e APIError) Error() string {
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

type Request[Req any] struct {
	Req Req
}

type Result[Res any, Err any] struct {
	Code  int
	Res   Res
	Error Err
}

type Response[Res any] struct {
	Code int
	Res  Res
}

type E interface {
	error
}

// type Error struct {
// 	Code  int
// 	Error error
// }

type ErrType = error

// type ErrType = any

type Error[Err ErrType] struct {
	Code  int
	Error Err
}

func PtrTo[T any](t T) *T {
	return &t
}

func LogYAML(v any, label string) {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s:\n", label)
	fmt.Println(string(bytes))

	return
}

func LogJSON(v any, label string) {
	bytes, err := json.MarshalIndent(v, "", "  ")

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s:\n", label)
	fmt.Println(string(bytes))

	return
}

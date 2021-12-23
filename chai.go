package chai

import (
	"encoding/json"
	"go/ast"
	"net/http"
)

type Methoder interface {
	Method(method, pattern string, h http.Handler)
}

type Reqer interface {
	Req() interface{}
}

type Reser interface {
	Res() interface{}
}

type Commenter interface {
	Comment() string
}

type ASTFiler interface {
	ASTFile() *ast.File
}

type APIError struct {
	Message string `json:"msg"`
}

func writeBytes(w http.ResponseWriter, code int, bytes []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(bytes)
}

func write(w http.ResponseWriter, code int, v interface{}) {
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

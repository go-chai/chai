package chai

import (
	"fmt"
	"go/ast"
	"net/http"

	"github.com/go-chi/docgen"
	"gopkg.in/yaml.v2"
)

type resHandlerFunc[Res any] func(http.ResponseWriter, *http.Request) (*Res, int, error)

func newResHandlerFunc[Res any](h resHandlerFunc[Res]) *resHandler[Res] {
	return &resHandler[Res]{
		f:       h,
		res:     new(Res),
		comment: docgen.GetFuncInfo(h).Comment,
		astFile: GetFuncInfo(h).ASTFile,
	}
}

type resHandler[Res any] struct {
	f       resHandlerFunc[Res]
	res     *Res
	comment string
	astFile *ast.File
}

func (h *resHandler[Res]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, code, err := h.f(w, r)
	if err != nil {
		if code == 0 {
			code = http.StatusInternalServerError
		}

		write(w, code, APIError{Message: err.Error()})
		return
	}

	if code == 0 {
		code = http.StatusOK
	}

	write(w, code, res)
}

func (h *resHandler[Res]) Res() any {
	return h.res
}

func (h *resHandler[Res]) Comment() string {
	return h.comment
}
func (h *resHandler[Res]) ASTFile() *ast.File {
	return h.astFile
}

func logg2(v interface{}) {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

	return
}

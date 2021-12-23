package openapi3

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"github.com/go-chai/chai"
	"github.com/go-chi/chi/v5"
	"gopkg.in/yaml.v2"
)

func NewDocs() *openapi3.T {
	return &openapi3.T{}
}

func Docs(r chi.Router) (*openapi3.T, error) {
	t := NewDocs()

	gen := openapi3gen.NewGenerator()
	schemas := make(openapi3.Schemas)

	err := chi.Walk(r, func(method string, route string, h http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		op := openapi3.NewOperation()

		if reqer, ok := h.(chai.Reqer); ok {
			rref, err := gen.NewSchemaRefForValue(reqer.Req(), schemas)
			if err != nil {
				return err
			}
			op.RequestBody = &openapi3.RequestBodyRef{
				Value: openapi3.NewRequestBody().WithJSONSchemaRef(rref),
			}
		}

		if reser, ok := h.(chai.Reser); ok {
			rref, err := gen.NewSchemaRefForValue(reser.Res(), schemas)
			if err != nil {
				return err
			}
			op.AddResponse(http.StatusOK, openapi3.NewResponse().WithJSONSchemaRef(rref))
		}

		t.AddOperation(route, method, op)

		// t.Paths.

		return nil
	})

	return t, err
}

func log(v interface{}) {
	bytes, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

	return
}

func logg2(v interface{}) {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

	return
}

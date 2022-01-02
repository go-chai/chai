package chai

import (
	"github.com/go-chai/chai/openapi2"
	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/spec"
)

func OpenAPI2(r chi.Routes) (*spec.Swagger, error) {
	return openapi2.Docs(func(parseOperationFn openapi2.OperationParserFunc) error {
		return chi.Walk(r, chi.WalkFunc(parseOperationFn))
	})
}

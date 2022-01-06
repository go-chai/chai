package chai

import (
	kinopenapi3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chai/chai/openapi3"
	"github.com/go-chi/chi/v5"
)

func OpenAPI3(r chi.Routes) (*kinopenapi3.T, error) {
	routes, err := getChiRoutes(r)
	if err != nil {
		return nil, err
	}

	return openapi3.Docs(routes)
}

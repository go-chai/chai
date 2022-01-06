package chai

import (
	kinopenapi3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chai/chai/openapi3"
	"github.com/gorilla/mux"
)

func OpenAPI3(r *mux.Router) (*kinopenapi3.T, error) {
	routes, err := getGorillaRoutes(r)
	if err != nil {
		return nil, err
	}
	return openapi3.Docs(routes)
}

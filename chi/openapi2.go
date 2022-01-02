package chai

import (
	"net/http"

	"github.com/go-chai/chai/openapi2"
	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/spec"
)

func OpenAPI2(r chi.Routes) (*spec.Swagger, error) {
	routes := make([]*openapi2.Route, 0)

	err := chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		routes = append(routes, &openapi2.Route{
			Method:  method,
			Path:    route,
			Handler: handler,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return openapi2.Docs(routes)
}

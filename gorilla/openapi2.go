package chai

import (
	"strings"

	"github.com/go-chai/chai/openapi2"
	"github.com/go-openapi/spec"
	"github.com/gorilla/mux"
)

func OpenAPI2(r *mux.Router) (*spec.Swagger, error) {
	routes := make([]*openapi2.Route, 0)

	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		methods, err := route.GetMethods()
		if err != nil && !strings.Contains(err.Error(), "route doesn't have methods") {
			return err
		}

		for _, method := range methods {
			routes = append(routes, &openapi2.Route{
				Method:  method,
				Path:    path,
				Handler: route.GetHandler(),
			})
		}

		return nil
	})

	return openapi2.Docs(routes)
}

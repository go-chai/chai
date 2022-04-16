package chai

import (
	"strings"

	chai "github.com/go-chai/chai/chi"
	"github.com/go-chai/chai/openapi2"
	"github.com/gorilla/mux"
	"github.com/zhamlin/chi-openapi/pkg/openapi/operations"
)

func OpenAPI2(r *mux.Router) (operations.OpenAPI, error) {
	routes, err := getGorillaRoutes(r)

	if err != nil {
		return nil, err
	}

	return openapi2.Docs(routes)
}

func getGorillaRoutes(r *mux.Router) ([]*openapi2.Route, error) {
	routes := make([]*openapi2.Route, 0)

	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		methods, err := route.GetMethods()
		if err != nil && !strings.Contains(err.Error(), "route doesn't have methods") {
			return err
		}

		for _, method := range methods {
			params, regexlessPath := chai.ParsePathParams(path)

			routes = append(routes, &openapi2.Route{
				Method:  method,
				Path:    regexlessPath,
				Params:  params,
				Handler: route.GetHandler(),
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return routes, nil
}

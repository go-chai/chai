package chai

import (
	"strings"

	"github.com/go-chai/chai/openapi2"
	"github.com/go-openapi/spec"
	"github.com/gorilla/mux"
)

func OpenAPI2(r *mux.Router) (*spec.Swagger, error) {
	return openapi2.Docs(func(parseOperationFn openapi2.OperationParserFunc) error {
		return r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			path, err := route.GetPathTemplate()
			if err != nil {
				return err
			}
			methods, err := route.GetMethods()
			if err != nil && !strings.Contains(err.Error(), "route doesn't have methods") {
				return err
			}

			for _, method := range methods {
				err := parseOperationFn(method, path, route.GetHandler())
				if err != nil {
					return err
				}
			}

			return nil
		})
	})
}

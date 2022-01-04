package chai

import (
	"net/http"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-chai/chai/openapi2"
	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/spec"
)

func OpenAPI2(r chi.Routes) (*spec.Swagger, error) {
	routes := make([]*openapi2.Route, 0)

	err := chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		params, newPath := getParams(route)
		routes = append(routes, &openapi2.Route{
			Method:  method,
			Path:    newPath,
			Params:  params,
			Handler: handler,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return openapi2.Docs(routes)
}

func getParams(path string) ([]spec.Parameter, string) {
	res := make([]spec.Parameter, 0)
	newPath := ""

	for {
		param, before, rest := nextParam(path)
		newPath += before

		if param == nil {
			break
		}

		newPath += "{" + param.Name + "}"

		res = append(res, *param)
		path = rest
	}

	if len(res) != 0 {
		// openapi2.LogYAML(res)
	}

	return res, newPath
}

func nextParam(pattern string) (param *spec.Parameter, path string, rest string) {
	path, rest, found := strings.Cut(pattern, "{")
	if !found {
		return nil, path, rest
	}

	// Read to closing } taking into account opens and closes in curl count (cc)
	cc := 1
	pe := 0

	for i, c := range rest {
		if c == '{' {
			cc++
		} else if c == '}' {
			cc--

			if cc == 0 {
				pe = i
				break
			}
		}
	}
	spew.Fdump(os.Stderr, rest[:pe])

	key := rest[:pe]
	pe++
	rest = rest[pe:]

	key, rexpat, _ := strings.Cut(key, ":")

	if len(rexpat) > 0 {
		if rexpat[0] != '^' {
			rexpat = "^" + rexpat
		}
		if rexpat[len(rexpat)-1] != '$' {
			rexpat += "$"
		}
	}

	return &spec.Parameter{
		CommonValidations: spec.CommonValidations{
			Pattern: rexpat,
		},
		ParamProps: spec.ParamProps{
			Name:     key,
			In:       "path",
			Required: true,
		},
		SimpleSchema: spec.SimpleSchema{
			Type: "string",
		},
	}, path, rest
}

func getParams2(route string) []spec.Parameter {
	res := make([]spec.Parameter, 0)

	for _, sect := range strings.Split(route, "/") {
		if strings.Contains(sect, "{") {
			param := spec.Parameter{}
			param.Name = strings.Trim(sect, "{}")
			param.In = "path"
			param.Required = true
			param.Type = "string"
			param.Pattern = "^[a-zA-Z0-9]+$"

			res = append(res, param)
		}
	}

	return res
}

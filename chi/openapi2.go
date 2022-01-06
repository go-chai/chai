package chai

import (
	"net/http"
	"strings"

	"github.com/go-chai/chai/openapi2"
	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/spec"
)

var integerSchema = spec.SimpleSchema{Type: "integer"}
var numberSchema = spec.SimpleSchema{Type: "number"}

var RegexPatternSchemas = map[string]spec.SimpleSchema{
	"/^(0|-*[1-9]+[0-9]*)$/":  integerSchema,
	"^[0-9]+$":                integerSchema,
	"[+-]?([0-9]*[.])?[0-9]+": numberSchema,
}

func OpenAPI2(r chi.Routes) (*spec.Swagger, error) {
	routes, err := getChiRoutes(r)

	if err != nil {
		return nil, err
	}

	return openapi2.Docs(routes)
}

func getChiRoutes(r chi.Routes) ([]*openapi2.Route, error) {
	routes := make([]*openapi2.Route, 0)

	err := chi.Walk(r, func(method, path string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		params, regexlessPath := ParsePathParams(path)
		routes = append(routes, &openapi2.Route{
			Method:  method,
			Path:    regexlessPath,
			Params:  params,
			Handler: handler,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return routes, nil
}

func ParsePathParams(path string) ([]spec.Parameter, string) {
	res := make([]spec.Parameter, 0)
	regexlessPath := ""

	for {
		param, before, after := nextParam(path)
		regexlessPath += before

		if param == nil {
			break
		}

		regexlessPath += "{" + param.Name + "}"

		res = append(res, *param)
		path = after
	}

	return res, regexlessPath
}

func nextParam(pattern string) (param *spec.Parameter, before string, after string) {
	before, after, found := strings.Cut(pattern, "{")
	if !found {
		return nil, before, after
	}

	// Read to closing } taking into account opens and closes in curl count (cc)
	cc := 1
	pe := 0

	for i, c := range after {
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

	key := after[:pe]
	after = after[pe+1:]

	key, rexpat, _ := strings.Cut(key, ":")

	if len(rexpat) > 0 {
		if rexpat[0] != '^' {
			rexpat = "^" + rexpat
		}
		if rexpat[len(rexpat)-1] != '$' {
			rexpat += "$"
		}
	}

	schema, ok := RegexPatternSchemas[rexpat]
	if !ok {
		schema = spec.SimpleSchema{
			Type: "string",
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
		SimpleSchema: schema,
	}, before, after
}

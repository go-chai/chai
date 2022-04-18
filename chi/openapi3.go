package chai

import (
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chai/chai/chai"
	chaiopenapi "github.com/go-chai/chai/openapi3"
	"github.com/go-chi/chi/v5"
)

func OpenAPI3(r chi.Routes) (*openapi3.T, error) {
	routes, err := getChiRoutes(r)

	if err != nil {
		return nil, err
	}

	return chaiopenapi.Docs(routes)
}

type Metadater interface {
	GetMetadata() *chai.Metadata
}

func GetHandlerInfo(fn http.Handler) *chai.Metadata {
	if fn, ok := fn.(Metadater); ok {
		return fn.GetMetadata()
	}

	return nil
}

func getChiRoutes(r chi.Routes) ([]*chaiopenapi.Route, error) {
	routes := make([]*chaiopenapi.Route, 0)

	err := chi.Walk(r, func(method, path string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		params, regexlessPath := ParsePathParams(path)

		routes = append(routes, &chaiopenapi.Route{
			Method:   method,
			Path:     regexlessPath,
			Params:   params,
			Handler:  handler,
			Metadata: GetHandlerInfo(handler),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return routes, nil
}

func ParsePathParams(path string) (openapi3.Parameters, string) {
	res := make(openapi3.Parameters, 0)
	regexlessPath := ""

	for {
		param, before, after := nextParam(path)
		regexlessPath += before

		if param == nil {
			break
		}

		regexlessPath += "{" + param.Value.Name + "}"

		res = append(res, param)
		path = after
	}

	return res, regexlessPath
}

func nextParam(pattern string) (param *openapi3.ParameterRef, before string, after string) {
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

	return &openapi3.ParameterRef{
		Ref: "",
		Value: &openapi3.Parameter{
			Name:     key,
			In:       "path",
			Required: true,
			Schema: &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:    "string",
					Pattern: rexpat,
				},
			}},
	}, before, after
}

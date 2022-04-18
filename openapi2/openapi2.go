package openapi2

import (
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chai/chai/chai"
	"github.com/go-chai/swag/gen"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
	"github.com/zhamlin/chi-openapi/pkg/openapi"
	"github.com/zhamlin/chi-openapi/pkg/openapi/operations"
)

type GenConfig = gen.GenConfig

func WriteDocs(docs *spec.Swagger, cfg *GenConfig) error {
	return gen.New().Generate(docs, cfg)
}

type Route struct {
	Method      string
	Path        string
	Params      openapi3.Parameters
	Handler     http.Handler
	Middlewares []func(http.Handler) http.Handler
}

func Docs(routes []*Route) (*openapi3.T, error) {
	var err error

	var spec = &openapi.OpenAPI{
		RegisteredTypes: openapi.RegisteredTypes{},
		T: &openapi3.T{
			Info: &openapi3.Info{
				Version: "0.0.1",
				Title:   "Title",
			},
			Servers: openapi3.Servers{},
			OpenAPI: "3.0.0",
			Paths:   openapi3.Paths{},
			Components: openapi3.Components{
				Schemas:         openapi3.Schemas{},
				Parameters:      openapi3.ParametersMap{},
				Responses:       map[string]*openapi3.ResponseRef{},
				SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{},
			},
		},
	}

	for _, route := range routes {
		err = RegisterRoute(spec, route)
		if err != nil {
			return nil, err
		}
	}

	return spec.T, nil
}

type HandlerInfo struct {
	IsReqer        bool
	Req            any
	IsReser        bool
	Res            any
	IsErrer        bool
	Err            any
	IsOper         bool
	Op             *operations.Operation
	IsChaiHandler  bool
	HandlerFunc    any
	HandlerWrapper http.Handler
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

func RegisterRoute(spec operations.OpenAPI, route *Route) error {
	var err error
	hi := GetHandlerInfo(route.Handler)

	op := hi.Op

	err = updateRequests(spec, op, hi, route.Params)
	if err != nil {
		return err
	}

	err = updateResponses(spec, op, hi)
	if err != nil {
		return err
	}

	spec.AddOperation(route.Path, route.Method, &op.Operation)

	return nil
}

func updateRequests(spec operations.OpenAPI, op *operations.Operation, hi *chai.Metadata, pathParams openapi3.Parameters) error {
	var err error

	if reflect.TypeOf(hi.Req) == nil {
		// log.Dump(op)
		op.Parameters = mergeSlices(makeKey, cmpKeys, mergeParamsFn, pathParams, op.Parameters)

		return nil
	}

	// inferredParams2 := openapi.SchemaFromObj(hi.Req, openapi.Schemas(spec.Components.Schemas), spec.RegisteredTypes)

	inferredParams, err := openapi.ParamsFromType(reflect.TypeOf(hi.Req).Elem().Elem(), openapi.Schemas(spec.Components.Schemas), spec.RegisteredTypes)
	if err != nil {
		return errors.Wrap(err, "failed to parse schema")
	}
	// log.Dump(inferredParams2)
	// log.Dump(inferredParams)
	// log.Dump(reflect.TypeOf(hi.Req).Elem().Elem())
	// log.Dump(reflect.TypeOf(hi.Req).Elem())
	// log.Dump(reflect.TypeOf(hi.Req))

	op.Parameters = mergeSlices(makeKey, cmpKeys, mergeParamsFn, pathParams, inferredParams, op.Parameters)

	return nil
}

type key struct {
	In   string
	Name string
}

func makeKey(p *openapi3.ParameterRef) key {
	return key{p.Value.In, p.Value.Name}
}
func cmpKeys(a, b key) bool {
	if a.In == b.In {
		return a.Name < b.Name
	}

	return a.In < b.In
}

func mergeParamsFn(a, b *openapi3.ParameterRef) *openapi3.ParameterRef {
	return b
}

func mergeSlices[K comparable, V any](keyFn func(V) K, cmp func(K, K) bool, mergeFn func(V, V) V, slices ...[]V) []V {
	m := make(map[K]V)

	for _, slice := range slices {
		m = mergeMaps(mergeFn, m, associateBy(slice, keyFn))
	}

	return sortedValues(m, cmp)
}

func mergeMaps[K comparable, V any](mergeFn func(V, V) V, maps ...map[K]V) map[K]V {
	res := make(map[K]V)

	for _, m := range maps {
		for k, v := range m {
			res[k] = mergeFn(res[k], v)
		}
	}

	return res
}

func associateBy[K comparable, V any](slice []V, keyFn func(V) K) map[K]V {
	m := make(map[K]V)

	for _, t := range slice {
		m[keyFn(t)] = t
	}

	return m
}

func sortedValues[K comparable, V any](m map[K]V, cmp func(K, K) bool) []V {
	res := make([]V, len(m))

	for i, k := range sortedKeys(m, cmp) {
		res[i] = m[k]
	}

	return res
}

func sortedKeys[K comparable, V any](m map[K]V, less func(K, K) bool) []K {
	keys := make([]K, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	if less == nil {
		return keys
	}

	sort.Slice(keys, func(i, j int) bool { return less(keys[i], keys[j]) })

	return keys
}

func updateResponses(spec operations.OpenAPI, op *operations.Operation, hi *chai.Metadata) error {
	if reflect.TypeOf(hi.Res) == nil {
		return nil
	}

	resSchema := openapi.SchemaFromObj(hi.Res, openapi.Schemas(spec.Components.Schemas), spec.RegisteredTypes)

	errSchema := openapi.SchemaFromObj(hi.Err, openapi.Schemas(spec.Components.Schemas), spec.RegisteredTypes)

	responses := op.Responses
	if responses == nil {
		responses = openapi3.Responses{}
		op.Responses = responses
	}
	noErrors := true
	noResponses := true
	for code := range op.Responses {
		code, err := strconv.Atoi(code)
		if err != nil {
			return err
		}

		if code < http.StatusBadRequest {
			noResponses = false
			updateResponseSchema(&op.Operation, responses, code, resSchema)
		}

		if code >= http.StatusBadRequest {
			noErrors = false
			updateResponseSchema(&op.Operation, responses, code, errSchema)
		}
	}
	if noResponses {
		op.AddResponse(http.StatusOK, openapi3.NewResponse().WithJSONSchemaRef(resSchema))
	}
	if noErrors {
		op.AddResponse(0, openapi3.NewResponse().WithJSONSchemaRef(errSchema))
	}

	return nil
}

func typeName(i any) string {
	t := reflect.TypeOf(i)

	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	s := strings.ReplaceAll(t.String(), " ", "")
	s = strings.ReplaceAll(s, "*", "")

	if s == "error" {
		return "string"
	}

	return s
}

func updateResponseSchema(op *openapi3.Operation, responses openapi3.Responses, code int, schema *openapi3.SchemaRef) {
	s := op.Responses.Get(code)
	if s.Value != nil {
		return
	}
	s.Value.WithJSONSchemaRef(schema)
}

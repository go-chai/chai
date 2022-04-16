package openapi2

import (
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
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

func Docs(routes []*Route) (operations.OpenAPI, error) {
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

	return spec, nil
}

type HandlerInfo struct {
	IsReqer        bool
	Req            any
	IsReser        bool
	Res            any
	IsErrer        bool
	Err            any
	IsOper         bool
	Op             operations.Operation
	IsChaiHandler  bool
	HandlerFunc    any
	HandlerWrapper http.Handler
}

func GetHandlerInfo(fn http.Handler) *HandlerInfo {
	hi := new(HandlerInfo)
	hi.HandlerWrapper = fn

	ch, ok := fn.(chai.Handlerer)
	if ok {
		hi.IsChaiHandler = true
		hi.HandlerFunc = ch.Handler()
	}

	reqer, ok := fn.(chai.Reqer)
	if ok {
		hi.IsReqer = true
		hi.Req = reqer.Req()
	}

	resErrer, ok := fn.(chai.ResErrer)
	if ok {
		hi.IsReser = true
		hi.Res = resErrer.Res()
		hi.IsErrer = true
		hi.Err = resErrer.Err()
	}

	oper, ok := fn.(chai.Oper)
	if ok {
		hi.IsOper = true
		hi.Op = oper.Op()
	}

	return hi
}

func SpecGen(value any) (*spec.Schema, error) {
	schemas := openapi3.Schemas{}

	g := openapi3gen.NewGenerator()

	// schemaRef, err := openapi3gen.NewSchemaRefForValue(value, schemas)
	ref, err := g.NewSchemaRefForValue(value, schemas)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create schema ref")
	}

	LogYAML(ref)
	LogYAML(schemas)

	return nil, nil
}

func SpecGen2(value any) (*spec.Schema, error) {
	schemas := openapi.Schemas{}
	schemaRef := openapi.SchemaFromObj(value, schemas, nil)

	LogYAML(schemaRef)
	LogYAML(schemas)

	return nil, nil
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

func updateRequests(spec operations.OpenAPI, op operations.Operation, hi *HandlerInfo, params openapi3.Parameters) error {
	var err error

	if !hi.IsReqer {
		op.Parameters = mergeParameters(params, op.Parameters)

		return nil
	}

	reqParams, err := openapi.ParamsFromType(reflect.TypeOf(hi.Req).Elem().Elem(), openapi.Schemas(spec.Components.Schemas), spec.RegisteredTypes)
	if err != nil {
		return errors.Wrap(err, "failed to parse schema")
	}

	op.Parameters = mergeParameters(params, reqParams, op.Parameters)

	return nil
}

type pk struct {
	In   string
	Name string
}

func less(pk, pk2 pk) bool {
	if pk.In == pk2.In {
		return pk.Name < pk2.Name
	}

	return pk.In < pk2.In
}

func mergeParameters(paramsList ...openapi3.Parameters) openapi3.Parameters {
	m := make(map[pk]*openapi3.ParameterRef)

	for _, params := range paramsList {
		m = mergeMaps(m, associateBy(params, func(p *openapi3.ParameterRef) pk {
			return pk{p.Value.In, p.Value.Name}
		}))
	}

	return sortedValues(m, less)
}

func mergeMaps[K comparable, V any](maps ...map[K]V) map[K]V {
	res := make(map[K]V)

	for _, m := range maps {
		for k, v := range m {
			res[k] = v
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

func sortedValues[K comparable, V any](m map[K]V, less func(K, K) bool) []V {
	res := make([]V, len(m))

	for i, k := range sortedKeys(m, less) {
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

func updateResponses(spec operations.OpenAPI, op operations.Operation, hi *HandlerInfo) error {
	if !hi.IsReser {
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

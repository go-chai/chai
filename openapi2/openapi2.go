package openapi2

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/go-chai/chai/chai"
	"github.com/go-chai/swag"
	"github.com/go-chai/swag/gen"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
)

type GenConfig = gen.GenConfig

func WriteDocs(docs *spec.Swagger, cfg *GenConfig) error {
	return gen.New().Generate(docs, cfg)
}

type Route struct {
	Method      string
	Path        string
	Params      []spec.Parameter
	Handler     http.Handler
	Middlewares []func(http.Handler) http.Handler
}

func Docs(routes []*Route) (*spec.Swagger, error) {
	var err error

	parser := swag.New(swag.SetDebugger(log.Default()), func(p *swag.Parser) {
		p.ParseDependency = true
	})

	for _, route := range routes {
		err = RegisterRoute(parser, route)
		if err != nil {
			return nil, err
		}
	}

	return parser.GetSwagger(), nil
}

func RegisterRoute(parser *swag.Parser, route *Route) error {
	var h = route.Handler
	var hh any = h

	ch, ok := h.(chai.Handlerer)
	if ok {
		hh = ch.Handler()
	}

	fi := GetFuncInfo(hh)

	op, err := parseSwaggoAnnotations(fi, parser)
	if err != nil {
		return err
	}

	err = updateRequests(fi, op, h, route.Params)
	if err != nil {
		return err
	}

	err = updateResponses(fi, op, h)
	if err != nil {
		return err
	}

	addOperation(parser.GetSwagger(), route.Path, route.Method, op)

	return nil
}

func parseSwaggoAnnotations(fi FuncInfo, parser *swag.Parser) (*swag.Operation, error) {
	var err error
	ops := swag.NewOperation(parser)

	pkg, err := getPkgPath(fi.File)
	if err != nil {
		return nil, err
	}

	err = parser.GetAllGoFileInfoAndParseTypes(pkg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse docs spec")
	}

	for _, line := range strings.Split(fi.Comment, "\n") {
		err := ops.ParseComment(line, fi.ASTFile)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse comment")
		}
	}

	return ops, nil
}

func getPkgPath(file string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "failed to get working directory")
	}

	file, err = filepath.Rel(wd, file)
	if err != nil {
		return "", errors.Wrap(err, "failed to get relative path")
	}

	return filepath.Dir(file), nil
}

func updateRequests(fi FuncInfo, op *swag.Operation, h http.Handler, params []spec.Parameter) error {
	var err error

	reqer, ok := h.(chai.Reqer)
	if !ok {
		op.Parameters = mergeParameters(params, op.Parameters)

		return nil
	}

	if len(op.Consumes) == 0 {
		op.Consumes = append(op.Consumes, "application/json")
	}

	schema, err := op.ParseAPIObjectSchema("object", typeName(reqer.Req()), fi.ASTFile)
	if err != nil {
		return err
	}

	noBody := true
	for i := range op.Parameters {
		if op.Parameters[i].In == "body" {
			noBody = false
			if op.Parameters[i].Schema == nil {
				op.Parameters[i].Schema = schema
			}
		}
	}
	if noBody {
		op.AddParam(spec.BodyParam("body", schema))
	}

	op.Parameters = mergeParameters(params, op.Parameters)

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

func mergeParameters(paramsList ...[]spec.Parameter) []spec.Parameter {
	m := make(map[pk]spec.Parameter)

	for _, params := range paramsList {
		m = mergeMaps(m, assoc(params, func(p spec.Parameter) pk {
			return pk{p.In, p.Name}
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

func assoc[K comparable, V any](slice []V, keyFn func(V) K) map[K]V {
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

func updateResponses(fi FuncInfo, op *swag.Operation, h http.Handler) error {
	resErrer, ok := h.(chai.ResErrer)
	if !ok {
		return nil
	}

	if len(op.Produces) == 0 {
		op.Produces = append(op.Produces, "application/json")
	}

	resSchema, err := op.ParseAPIObjectSchema("object", typeName(resErrer.Res()), fi.ASTFile)
	if err != nil {
		return err
	}

	errSchema, err := op.ParseAPIObjectSchema("object", typeName(resErrer.Err()), fi.ASTFile)
	if err != nil {
		return err
	}

	responses := op.Responses
	if responses == nil {
		responses = &spec.Responses{}
		op.Responses = responses
	}
	noErrors := true
	noResponses := true
	for code := range op.Responses.StatusCodeResponses {
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
		op.RespondsWith(http.StatusOK, spec.NewResponse().WithSchema(resSchema))
	}
	if noErrors {
		op.RespondsWith(0, spec.NewResponse().WithSchema(errSchema))
	}

	return nil
}

func typeName(i any) string {
	t := reflect.TypeOf(i)

	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	s := strings.ReplaceAll(t.String(), " ", "")

	if s == "error" {
		return "string"
	}

	return s
}

func updateResponseSchema(op *spec.Operation, responses *spec.Responses, code int, schema *spec.Schema) {
	s := op.Responses.StatusCodeResponses[code]

	if s.Schema != nil {
		return
	}

	s.Schema = schema

	op.Responses.StatusCodeResponses[code] = s
}

package openapi2

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/go-chai/chai/chai"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
	"github.com/swaggo/swag"
	"github.com/swaggo/swag/gen"
)

type GenConfig = gen.GenConfig

func WriteDocs(docs *spec.Swagger, cfg *GenConfig) error {
	return gen.New().Generate(docs, cfg)
}

func ParseOperation(docs *spec.Swagger, parser *swag.Parser, method string, route string, h http.Handler, middlewares ...func(http.Handler) http.Handler) error {
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

	err = updateRequests(fi, op, h)
	if err != nil {
		return err
	}

	err = updateResponses(fi, op, h)
	if err != nil {
		return err
	}

	addOperation(docs, route, method, op)

	return nil
}

type OperationParserFunc func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error

func Docs(walkRoutes func(OperationParserFunc) error) (*spec.Swagger, error) {
	var err error

	docs := newSpec()

	parser := swag.New(swag.SetDebugger(log.Default()), func(p *swag.Parser) {
		p.ParseDependency = true
	})

	operationParser := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		return ParseOperation(docs, parser, method, route, handler, middlewares...)
	}

	err = walkRoutes(operationParser)

	docs.Definitions = parser.GetSwagger().Definitions

	return docs, err
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

func updateRequests(fi FuncInfo, op *swag.Operation, h http.Handler) error {
	var err error

	reqer, ok := h.(chai.Reqer)
	if !ok {
		return nil
	}

	if len(op.Consumes) == 0 {
		op.Consumes = append(op.Consumes, "application/json")
	}

	schema, err := op.ParseAPIObjectSchema("object", typeName(reqer.Req()), fi.ASTFile)
	if err != nil {
		return err
	}

	if len(op.Parameters) == 0 {
		op.AddParam(spec.BodyParam("body", schema))
		return nil
	}

	for i := range op.Parameters {
		if op.Parameters[i].In != "body" {
			continue
		}

		if op.Parameters[i].Schema == nil {
			op.Parameters[i].Schema = schema
		}
	}

	return nil
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

package openapi2

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/go-chai/chai"
	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
	"github.com/swaggo/swag"
)

func Docs(r chi.Router) (*spec.Swagger, error) {
	var err error

	t := newSpec()

	parser := swag.New(swag.SetDebugger(log.Default()), func(p *swag.Parser) {
		p.ParseDependency = true
	})

	err = chi.Walk(r, func(method string, route string, h http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		ch, ok := h.(chai.Handlerer)
		if !ok {
			// ignore non-chai handlers
			return nil
		}

		fi := GetFuncInfo(ch.Handler())

		op, err := parseSwaggoAnnotations(fi, ch, parser)
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

		addOperation(t, route, method, op)

		return nil
	})

	t.Definitions = parser.GetSwagger().Definitions

	return t, err
}

func parseSwaggoAnnotations(fi FuncInfo, ch chai.Handlerer, parser *swag.Parser) (*swag.Operation, error) {
	var err error
	ops := swag.NewOperation(parser)

	pkg, err := getPkgPath(fi.File)
	if err != nil {
		return nil, err
	}

	err = parser.GetAllGoFileInfoAndParseTypes(pkg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse swagger spec")
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

	for strings.HasPrefix(t.String(), "*") {
		t = t.Elem()
	}

	if t.String() == "error" {
		return "string"
	}

	return t.String()
}

func updateResponseSchema(op *spec.Operation, responses *spec.Responses, code int, schema *spec.Schema) {
	s := op.Responses.StatusCodeResponses[code]

	if s.Schema != nil {
		return
	}

	s.Schema = schema

	op.Responses.StatusCodeResponses[code] = s
}

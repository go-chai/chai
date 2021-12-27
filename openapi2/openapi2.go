package openapi2

import (
	"constraints"
	"go/ast"
	"log"
	"net/http"
	"strings"

	"github.com/go-chai/chai"
	"github.com/go-chai/chai/specc"
	"github.com/go-chai/chai/specgen"
	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
	"github.com/swaggo/swag"
)

func Docs(r chi.Router) (*specc.Swagger, error) {
	var err error

	t := specc.New()

	gen := specgen.NewGenerator()
	schemas := make(map[string]spec.Schema)

	parser := swag.New(swag.SetDebugger(log.Default()), func(p *swag.Parser) {
		p.ParseDependency = true
	})

	err = chi.Walk(r, func(method string, route string, h http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if _, ok := h.(chai.Commenter); !ok {
			return nil
		}

		var astFile *ast.File

		if astFiler, ok := h.(chai.ASTFiler); ok {
			astFile = astFiler.ASTFile()
		}

		op := spec.NewOperation("")

		if commenter, ok := h.(chai.Commenter); ok {
			comment := commenter.Comment()
			ops := swag.NewOperation(parser)

			for _, line := range strings.Split(comment, "\n") {
				err := ops.ParseComment(line, astFile)
				if err != nil {
					return errors.Wrap(err, "failed to parse comment")
				}
			}

			op = &ops.Operation
		}

		if reqer, ok := h.(chai.Reqer); ok {
			if len(op.Consumes) == 0 {
				op.Consumes = append(op.Consumes, "application/json")
			}

			schema, err := gen.NewSchemaRefForValue(reqer.Req(), schemas)
			if err != nil {
				return err
			}

			if len(op.Parameters) == 0 {
				op.AddParam(spec.BodyParam("body", schema))
			} else {
				for i := range op.Parameters {
					if op.Parameters[i].In != "body" {
						continue
					}

					if op.Parameters[i].Schema == nil {
						op.Parameters[i].Schema = schema
					}
				}
			}
		}

		if reser, ok := h.(chai.Reser); ok {
			if len(op.Produces) == 0 {
				op.Produces = append(op.Produces, "application/json")
			}

			schema, err := gen.NewSchemaRefForValue(reser.Res(), schemas)
			if err != nil {
				return err
			}

			responses := op.Responses
			if responses == nil {
				responses = &spec.Responses{}
				op.Responses = responses
			}
			found := false
			for code := range op.Responses.StatusCodeResponses {
				if code < http.StatusBadRequest {
					found = true
					s := op.Responses.StatusCodeResponses[code]

					if s.Schema != nil {
						continue
					}

					s.Schema = schema

					op.Responses.StatusCodeResponses[code] = s

				}
			}
			if !found {
				op.RespondsWith(http.StatusOK, spec.NewResponse().WithSchema(schema))
			}
		}

		if errer, ok := h.(chai.Errer); ok {

			responses := op.Responses
			if responses == nil {
				responses = &spec.Responses{}
				op.Responses = responses
			}
			schema, err := gen.NewSchemaRefForValue(errer.Err(), schemas)
			if err != nil {
				return err
			}
			found := false
			for code := range op.Responses.StatusCodeResponses {
				if code >= http.StatusBadRequest {
					found = true
					s := op.Responses.StatusCodeResponses[code]

					if s.Schema != nil {
						continue
					}

					s.Schema = schema

					op.Responses.StatusCodeResponses[code] = s

				}
			}
			if !found {
				op.RespondsWith(0, spec.NewResponse().WithSchema(schema))
			}
		}

		t.AddOperation(route, method, op)

		return nil
	})

	t.Definitions = parser.GetSwagger().Definitions

	return t, err
}

type number interface {
	constraints.Integer | constraints.Float
}

func ptrTo[T any](t T) *T {
	return &t
}

func convNumPtr[U number, T number](t *T) *U {
	if t == nil {
		return nil
	}

	u := U(*t)

	return &u
}

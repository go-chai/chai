package openapi2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"github.com/ghodss/yaml"
	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
	"github.com/swaggo/swag"
	"github.com/go-chai/chai"
)

func Docs(r chi.Router, t *openapi3.T) error {
	gen := openapi3gen.NewGenerator()
	schemas := make(openapi3.Schemas)

	// spew.Dump(schemas)

	err := chi.Walk(r, func(method string, route string, h http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		op := &openapi3.Operation{}
		op2 := &openapi2.Operation{}
		op3 := spec.NewOperation("")

		// var astFile *ast.File

		// if astFiler, ok := h.(chai.ASTFiler); ok {
		// 	astFile = astFiler.ASTFile()

		// 	// spew.Dump(astFile)
		// }

		if commenter, ok := h.(chai.Commenter); ok {
			comment := commenter.Comment()

			fmt.Printf("comment: %s\n", comment)

			logg2("zzzzzzzzzzzzzz")
			parser := swag.New()

			// err := parser.Packages.CollectAstFile("zz", "zz", astFile)
			// if err != nil {
			// 	return err
			// }

			// parser

			ops := swag.NewOperation(parser)

			// op.

			for _, line := range strings.Split(comment, "\n") {
				err := ops.ParseCommentChai(line, nil)
				if err != nil {
					return errors.Wrap(err, "failed to parse comment")
				}
			}

			logg2("op")
			logg2(op)
			logg2(ops.Operation)
			logg2(op2)
			logg2(op3)

		}

		if reqer, ok := h.(chai.Reqer); ok {
			spec.RefSchema("")

			// op := spec.NewOperation("")

			// op.Parameters[0].In

			rref, err := gen.NewSchemaRefForValue(reqer.Req(), schemas)
			if err != nil {
				return err
			}
			op.RequestBody = &openapi3.RequestBodyRef{
				Value: openapi3.NewRequestBody().WithJSONSchemaRef(rref),
			}
			op2.Parameters = append(op2.Parameters, &openapi2.Parameter{
				In:          "",
				Name:        "",
				Description: "",
				Type:        "",
			})
			param := spec.BodyParam("", spec.RefSchema(""))
			op3.Parameters = append(op3.Parameters, *param)

		}

		if reser, ok := h.(chai.Reser); ok {
			rref, err := gen.NewSchemaRefForValue(reser.Res(), schemas)
			if err != nil {
				return err
			}
			op.AddResponse(http.StatusOK, openapi3.NewResponse().WithJSONSchemaRef(rref))
		}

		t.AddOperation(route, method, op)

		return nil
	})

	// log(gen.SchemaRefs)
	// log(gen.Types)
	log("schemas")
	log(schemas)

	return err
}

func log(v interface{}) {
	bytes, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

	return
}

func logg2(v interface{}) {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

	return
}

package raml

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/bxcodec/faker/v3"
	"github.com/go-chai/chai"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/docgen"
	"github.com/go-chi/docgen/raml"
	"github.com/pkg/errors"
)

func Docs(r chi.Router, ramlDocs *raml.RAML) error {
	err := chi.Walk(r, func(method string, route string, h http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fi := docgen.GetFuncInfo(h)

		resource := &raml.Resource{
			DisplayName:     fi.Func,
			Description:     fi.Comment,
			Is:              []string{},
			Type:            "",
			SecuredBy:       []string{},
			UriParameters:   []string{},
			QueryParameters: []string{},
			Resources:       map[string]*raml.Resource{},
		}

		if reqer, ok := h.(chai.Reqer); ok {
			req, err := exampleJSON(reqer.Req())
			if err != nil {
				return err
			}

			resource.Body = map[string]raml.Example{
				"application/json": {
					Example: req,
				},
			}
		}

		if reser, ok := h.(chai.Reser); ok {
			res, err := exampleJSON(reser.Res())
			if err != nil {
				return err
			}

			resource.Responses = map[int]raml.Response{
				http.StatusOK: {
					Body: raml.Body{
						"application/json": {
							Example: res,
						},
					}},
			}
		}

		return ramlDocs.Add(method, route, resource)
	})
	if err != nil {
		return err
	}

	return nil
}

func exampleJSON(v interface{}) (string, error) {
	newV := reflect.New(reflect.TypeOf(v)).Interface()

	err := faker.FakeData(newV)
	if err != nil {
		return "", errors.Wrap(err, "failed to fake data")
	}

	b, err := json.MarshalIndent(newV, "", "  ")
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal request")
	}

	return string(b), nil
}

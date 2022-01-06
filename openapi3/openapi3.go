package openapi3

import (
	kinopenapi2 "github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
	"github.com/go-chai/chai/openapi2"
	"github.com/go-chai/swag/gen"
)

type GenConfig = gen.GenConfig

func WriteDocs(docs *openapi3.T, cfg *GenConfig) error {
	return NewGen().Generate(docs, cfg)
}

type Route = openapi2.Route

func Docs(routes []*Route) (*openapi3.T, error) {
	docs, err := openapi2.Docs(routes)
	if err != nil {
		return nil, err
	}
	docsJSON, err := openapi2.MarshalYAML(docs)
	if err != nil {
		return nil, err
	}
	kinOpenAPI2 := new(kinopenapi2.T)
	err = yaml.Unmarshal(docsJSON, kinOpenAPI2)
	if err != nil {
		return nil, err
	}
	return openapi2conv.ToV3(kinOpenAPI2)
}

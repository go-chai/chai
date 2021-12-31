package openapi2

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/spec"
	"github.com/swaggo/swag"
)

func newSpec() *spec.Swagger {
	return &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Info: &spec.Info{
				InfoProps: spec.InfoProps{
					Contact: &spec.ContactInfo{},
					License: nil,
				},
				VendorExtensible: spec.VendorExtensible{
					Extensions: spec.Extensions{},
				},
			},
			Paths: &spec.Paths{
				Paths: make(map[string]spec.PathItem),
			},
			Definitions:         make(map[string]spec.Schema),
			SecurityDefinitions: make(map[string]*spec.SecurityScheme),
		},
	}
}

func addOperation(doc *spec.Swagger, path string, method string, operation *swag.Operation) {
	paths := doc.Paths
	if paths == nil {
		paths = &spec.Paths{}
		doc.Paths = paths
	}

	if paths.Paths == nil {
		paths.Paths = make(map[string]spec.PathItem)
	}

	// doc.Paths.Paths

	pathItem := paths.Paths[path]
	SetOperation(&pathItem, method, &operation.Operation)
	paths.Paths[path] = pathItem
}

func SetOperation(pathItem *spec.PathItem, method string, operation *spec.Operation) {
	switch method {
	case http.MethodDelete:
		pathItem.Delete = operation
	case http.MethodGet:
		pathItem.Get = operation
	case http.MethodHead:
		pathItem.Head = operation
	case http.MethodOptions:
		pathItem.Options = operation
	case http.MethodPatch:
		pathItem.Patch = operation
	case http.MethodPost:
		pathItem.Post = operation
	case http.MethodPut:
		pathItem.Put = operation
	default:
		panic(fmt.Errorf("unsupported HTTP method %q", method))
	}
}

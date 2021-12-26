package specc

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/spec"
)

type Swagger struct {
	*spec.Swagger
}

func (doc *Swagger) AddOperation(path string, method string, operation *spec.Operation) {
	// spew.Dump("operation")
	// spew.Dump(operation)

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
	SetOperation(&pathItem, method, operation)
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

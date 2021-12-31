package openapi2

import (
	"reflect"

	"github.com/go-chai/chai/specgen"
	"github.com/go-openapi/spec"
	"github.com/swaggo/swag"
)

func schemaCustomizer(p *swag.Parser) specgen.SchemaCustomizerFn {
	return func(name string, t reflect.Type, tag reflect.StructTag, schema *spec.Schema) error {
		return newTagBaseFieldParser(name, p, tag).ComplementSchema(schema)
	}
}

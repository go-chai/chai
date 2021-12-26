// Package specgen generates OpenAPIv3 JSON schemas from Go types.
package specgen

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/jsoninfo"
	"github.com/go-openapi/spec"
)

type Schemas map[string]spec.Schema

// CycleError indicates that a type graph has one or more possible cycles.
type CycleError struct{}

func (err *CycleError) Error() string { return "detected cycle" }

// Option allows tweaking SchemaRef generation
type Option func(*generatorOpt)

// SchemaCustomizerFn is a callback function, allowing
// the OpenAPI schema definition to be updated with additional
// properties during the generation process, based on the
// name of the field, the Go type, and the struct tags.
// name will be "_root" for the top level object, and tag will be ""
type SchemaCustomizerFn func(name string, t reflect.Type, tag reflect.StructTag, schema *spec.Schema) error

type generatorOpt struct {
	useAllExportedFields bool
	throwErrorOnCycle    bool
	schemaCustomizer     SchemaCustomizerFn
}

// UseAllExportedFields changes the default behavior of only
// generating schemas for struct fields with a JSON tag.
func UseAllExportedFields() Option {
	return func(x *generatorOpt) { x.useAllExportedFields = true }
}

// ThrowErrorOnCycle changes the default behavior of creating cycle
// refs to instead error if a cycle is detected.
func ThrowErrorOnCycle() Option {
	return func(x *generatorOpt) { x.throwErrorOnCycle = true }
}

// SchemaCustomizer allows customization of the schema that is generated
// for a field, for example to support an additional tagging scheme
func SchemaCustomizer(sc SchemaCustomizerFn) Option {
	return func(x *generatorOpt) { x.schemaCustomizer = sc }
}

// NewSchemaRefForValue is a shortcut for NewGenerator(...).NewSchemaRefForValue(...)
func NewSchemaRefForValue(value interface{}, schemas map[string]spec.Schema, opts ...Option) (*spec.Schema, error) {
	g := NewGenerator(opts...)
	return g.NewSchemaRefForValue(value, schemas)
}

type Generator struct {
	opts generatorOpt

	Types map[reflect.Type]*spec.Schema

	// SchemaRefs contains all references and their counts.
	// If count is 1, it's not ne
	// An OpenAPI identifier has been assigned to each.
	SchemaRefs map[*spec.Schema]int

	// componentSchemaRefs is a set of schemas that must be defined in the components to avoid cycles
	componentSchemaRefs map[string]struct{}
}

func NewGenerator(opts ...Option) *Generator {
	gOpt := &generatorOpt{}
	for _, f := range opts {
		f(gOpt)
	}
	return &Generator{
		Types:               make(map[reflect.Type]*spec.Schema),
		SchemaRefs:          make(map[*spec.Schema]int),
		componentSchemaRefs: make(map[string]struct{}),
		opts:                *gOpt,
	}
}

func (g *Generator) GenerateSchemaRef(t reflect.Type) (*spec.Schema, error) {
	//check generatorOpt consistency here
	return g.generateSchemaRefFor(nil, t, "_root", "")
}

// NewSchemaRefForValue uses reflection on the given value to produce a SchemaRef, and updates a supplied map with any dependent component schemas if they lead to cycles
func (g *Generator) NewSchemaRefForValue(value interface{}, schemas map[string]spec.Schema) (*spec.Schema, error) {
	ref, err := g.GenerateSchemaRef(reflect.TypeOf(value))
	if err != nil {
		return nil, err
	}
	for ref := range g.SchemaRefs {
		if _, ok := g.componentSchemaRefs[ref.Ref.String()]; ok && schemas != nil {
			schemas[ref.Ref.String()] = *ref
		}
		if strings.HasPrefix(ref.Ref.String(), "#/definitions/schemas/") {
			// ref.Value = nil
			// ref.Ref = spec.MustCreateRef("bb")
		} else {
			// TODO??
			// ref.Ref = spec.Ref{}
			// ref.Ref = spec.MustCreateRef("aa")
		}
	}
	return ref, nil
}

func (g *Generator) generateSchemaRefFor(parents []*jsoninfo.TypeInfo, t reflect.Type, name string, tag reflect.StructTag) (*spec.Schema, error) {
	if ref := g.Types[t]; ref != nil && g.opts.schemaCustomizer == nil {
		g.SchemaRefs[ref]++
		return ref, nil
	}
	ref, err := g.generateWithoutSaving(parents, t, name, tag)
	if err != nil {
		return nil, err
	}
	if ref != nil {
		g.Types[t] = ref
		g.SchemaRefs[ref]++
	}
	return ref, nil
}

func getStructField(t reflect.Type, fieldInfo jsoninfo.FieldInfo) reflect.StructField {
	var ff reflect.StructField
	// fieldInfo.Index is an array of indexes starting from the root of the type
	for i := 0; i < len(fieldInfo.Index); i++ {
		ff = t.Field(fieldInfo.Index[i])
		t = ff.Type
	}
	return ff
}

func (g *Generator) generateWithoutSaving(parents []*jsoninfo.TypeInfo, t reflect.Type, name string, tag reflect.StructTag) (*spec.Schema, error) {
	typeInfo := jsoninfo.GetTypeInfo(t)
	for _, parent := range parents {
		if parent == typeInfo {
			return nil, &CycleError{}
		}
	}

	if cap(parents) == 0 {
		parents = make([]*jsoninfo.TypeInfo, 0, 4)
	}
	parents = append(parents, typeInfo)

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if strings.HasSuffix(t.Name(), "Ref") {
		_, a := t.FieldByName("Ref")
		v, b := t.FieldByName("Value")
		if a && b {
			vs, err := g.generateSchemaRefFor(parents, v.Type, name, tag)
			if err != nil {
				if _, ok := err.(*CycleError); ok && !g.opts.throwErrorOnCycle {
					g.SchemaRefs[vs]++
					return vs, nil
				}
				return nil, err
			}
			refSchemaRef := RefSchemaRef
			g.SchemaRefs[refSchemaRef]++
			ref := NewSchemaRef(t.Name(), &spec.Schema{
				SchemaProps: spec.SchemaProps{
					OneOf: []spec.Schema{
						*refSchemaRef,
						*vs,
					},
				},
			})

			g.SchemaRefs[ref]++
			return ref, nil
		}
	}

	schema := &spec.Schema{}

	switch t.Kind() {
	case reflect.Func, reflect.Chan:
		return nil, nil // ignore

	case reflect.Bool:
		schema.Type = []string{"boolean"}

	case reflect.Int:
		schema.Type = []string{"integer"}
	case reflect.Int8:
		schema.Type = []string{"integer"}
		schema.Minimum = &minInt8
		schema.Maximum = &maxInt8
	case reflect.Int16:
		schema.Type = []string{"integer"}
		schema.Minimum = &minInt16
		schema.Maximum = &maxInt16
	case reflect.Int32:
		schema.Type = []string{"integer"}
		schema.Format = "int32"
	case reflect.Int64:
		schema.Type = []string{"integer"}
		schema.Format = "int64"
	case reflect.Uint:
		schema.Type = []string{"integer"}
		schema.Minimum = &zeroInt
	case reflect.Uint8:
		schema.Type = []string{"integer"}
		schema.Minimum = &zeroInt
		schema.Maximum = &maxUint8
	case reflect.Uint16:
		schema.Type = []string{"integer"}
		schema.Minimum = &zeroInt
		schema.Maximum = &maxUint16
	case reflect.Uint32:
		schema.Type = []string{"integer"}
		schema.Minimum = &zeroInt
		schema.Maximum = &maxUint32
	case reflect.Uint64:
		schema.Type = []string{"integer"}
		schema.Minimum = &zeroInt
		schema.Maximum = &maxUint64

	case reflect.Float32:
		schema.Type = []string{"number"}
		schema.Format = "float"
	case reflect.Float64:
		schema.Type = []string{"number"}
		schema.Format = "double"

	case reflect.String:
		schema.Type = []string{"string"}

	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			if t == rawMessageType {
				return schema, nil
			}
			schema.Type = []string{"string"}
			schema.Format = "byte"
		} else {
			schema.Type = []string{"array"}
			items, err := g.generateSchemaRefFor(parents, t.Elem(), name, tag)
			if err != nil {
				if _, ok := err.(*CycleError); ok && !g.opts.throwErrorOnCycle {
					items = g.generateCycleSchemaRef(t.Elem(), schema)
				} else {
					return nil, err
				}
			}
			if items != nil {
				g.SchemaRefs[items]++
				schema.Items = &spec.SchemaOrArray{Schema: items}
			}
		}

	case reflect.Map:
		schema.Type = []string{"object"}
		additionalProperties, err := g.generateSchemaRefFor(parents, t.Elem(), name, tag)
		if err != nil {
			if _, ok := err.(*CycleError); ok && !g.opts.throwErrorOnCycle {
				additionalProperties = g.generateCycleSchemaRef(t.Elem(), schema)
			} else {
				return nil, err
			}
		}
		if additionalProperties != nil {
			g.SchemaRefs[additionalProperties]++
			schema.AdditionalProperties = &spec.SchemaOrBool{
				Schema: additionalProperties,
			}
		}

	case reflect.Struct:
		if t == timeType {
			schema.Type = []string{"string"}
			schema.Format = "date-time"
		} else {
			for _, fieldInfo := range typeInfo.Fields {
				// Only fields with JSON tag are considered (by default)
				if !fieldInfo.HasJSONTag && !g.opts.useAllExportedFields {
					continue
				}
				// If asked, try to use yaml tag
				fieldName, fType := fieldInfo.JSONName, fieldInfo.Type
				if !fieldInfo.HasJSONTag && g.opts.useAllExportedFields {
					// Handle anonymous fields/embedded structs
					if t.Field(fieldInfo.Index[0]).Anonymous {
						ref, err := g.generateSchemaRefFor(parents, fType, fieldName, tag)
						if err != nil {
							if _, ok := err.(*CycleError); ok && !g.opts.throwErrorOnCycle {
								ref = g.generateCycleSchemaRef(fType, schema)
							} else {
								return nil, err
							}
						}
						if ref != nil {
							g.SchemaRefs[ref]++
							WithPropertyRef(schema, fieldName, ref)
						}
					} else {
						ff := getStructField(t, fieldInfo)
						if tag, ok := ff.Tag.Lookup("yaml"); ok && tag != "-" {
							fieldName, fType = tag, ff.Type
						}
					}
				}

				// extract the field tag if we have a customizer
				var fieldTag reflect.StructTag
				if g.opts.schemaCustomizer != nil {
					ff := getStructField(t, fieldInfo)
					fieldTag = ff.Tag
				}

				ref, err := g.generateSchemaRefFor(parents, fType, fieldName, fieldTag)
				if err != nil {
					if _, ok := err.(*CycleError); ok && !g.opts.throwErrorOnCycle {
						ref = g.generateCycleSchemaRef(fType, schema)
					} else {
						return nil, err
					}
				}
				if ref != nil {
					g.SchemaRefs[ref]++
					WithPropertyRef(schema, fieldName, ref)
				}
			}

			// Object only if it has properties
			if schema.Properties != nil {
				schema.Type = []string{"object"}
			}
		}
	}

	if g.opts.schemaCustomizer != nil {
		if err := g.opts.schemaCustomizer(name, t, tag, schema); err != nil {
			return nil, err
		}
	}

	// return NewSchemaRef(t.Name()+"zzazz", schema), nil
	return schema, nil
}

func WithPropertyRef(schema *spec.Schema, name string, ref *spec.Schema) *spec.Schema {
	properties := schema.Properties
	if properties == nil {
		properties = make(map[string]spec.Schema)
		schema.Properties = properties
	}
	properties[name] = *ref
	return schema

}
func NewSchemaRef(name string, schema *spec.Schema) *spec.Schema {
	schema.SchemaProps.Ref = spec.MustCreateRef(name)

	return schema
}

func (g *Generator) generateCycleSchemaRef(t reflect.Type, schema *spec.Schema) *spec.Schema {
	var typeName string
	switch t.Kind() {
	case reflect.Ptr:
		return g.generateCycleSchemaRef(t.Elem(), schema)
	case reflect.Slice:
		ref := g.generateCycleSchemaRef(t.Elem(), schema)
		sliceSchema := &spec.Schema{}
		sliceSchema.Type = []string{"array"}
		sliceSchema.Items = &spec.SchemaOrArray{
			Schema: ref,
		}
		return sliceSchema
	case reflect.Map:
		ref := g.generateCycleSchemaRef(t.Elem(), schema)
		mapSchema := &spec.Schema{}
		mapSchema.Type = []string{"object"}
		mapSchema.AdditionalProperties = &spec.SchemaOrBool{
			Schema: ref,
		}
		return mapSchema
	default:
		typeName = t.Name()
	}

	g.componentSchemaRefs[typeName] = struct{}{}

	return NewSchemaRef(fmt.Sprintf("#/components/schemas/%s", typeName), schema)
	// return spec.NewSchemaRef(fmt.Sprintf("#/components/schemas/%s", typeName), schema)
}

func refschema() *spec.Schema {
	props := make(map[string]spec.Schema)

	props["$rreff"] = *spec.StringProperty().WithMinLength(1)

	return NewSchemaRef("Ref", (&spec.Schema{}).WithProperties(props))
}

var RefSchemaRef = refschema()

var (
	timeType       = reflect.TypeOf(time.Time{})
	rawMessageType = reflect.TypeOf(json.RawMessage{})

	zeroInt   = float64(0)
	maxInt8   = float64(math.MaxInt8)
	minInt8   = float64(math.MinInt8)
	maxInt16  = float64(math.MaxInt16)
	minInt16  = float64(math.MinInt16)
	maxUint8  = float64(math.MaxUint8)
	maxUint16 = float64(math.MaxUint16)
	maxUint32 = float64(math.MaxUint32)
	maxUint64 = float64(math.MaxUint64)
)

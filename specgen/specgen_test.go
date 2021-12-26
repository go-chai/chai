package specgen_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-chai/chai"
	"github.com/go-chai/chai/example/celler/model"
	"github.com/go-chai/chai/specgen"
	"github.com/go-openapi/jsonreference"
	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/require"
)

func ExampleGenerator_SchemaRefs() {
	type SomeOtherType string
	type SomeStruct struct {
		Bool    bool                      `json:"bool"`
		Int     int                       `json:"int"`
		Int64   int64                     `json:"int64"`
		Float64 float64                   `json:"float64"`
		String  string                    `json:"string"`
		Bytes   []byte                    `json:"bytes"`
		JSON    json.RawMessage           `json:"json"`
		Time    time.Time                 `json:"time"`
		Slice   []SomeOtherType           `json:"slice"`
		Map     map[string]*SomeOtherType `json:"map"`

		Struct struct {
			X string `json:"x"`
		} `json:"struct"`

		EmptyStruct struct {
			Y string
		} `json:"structWithoutFields"`

		Ptr *SomeOtherType `json:"ptr"`
	}

	g := specgen.NewGenerator()
	schemaRef, err := g.NewSchemaRefForValue(&SomeStruct{}, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("g.SchemaRefs: %d\n", len(g.SchemaRefs))
	var data []byte
	if data, err = json.MarshalIndent(&schemaRef, "", "  "); err != nil {
		panic(err)
	}
	fmt.Printf("schemaRef: %s\n", data)
	// Output:
	// g.SchemaRefs: 15
	// schemaRef: {
	//   "type": "object",
	//   "properties": {
	//     "bool": {
	//       "type": "boolean"
	//     },
	//     "bytes": {
	//       "type": "string",
	//       "format": "byte"
	//     },
	//     "float64": {
	//       "type": "number",
	//       "format": "double"
	//     },
	//     "int": {
	//       "type": "integer"
	//     },
	//     "int64": {
	//       "type": "integer",
	//       "format": "int64"
	//     },
	//     "json": {},
	//     "map": {
	//       "type": "object",
	//       "additionalProperties": {
	//         "type": "string"
	//       }
	//     },
	//     "ptr": {
	//       "type": "string"
	//     },
	//     "slice": {
	//       "type": "array",
	//       "items": {
	//         "type": "string"
	//       }
	//     },
	//     "string": {
	//       "type": "string"
	//     },
	//     "struct": {
	//       "type": "object",
	//       "properties": {
	//         "x": {
	//           "type": "string"
	//         }
	//       }
	//     },
	//     "structWithoutFields": {},
	//     "time": {
	//       "type": "string",
	//       "format": "date-time"
	//     }
	//   }
	// }
}

// TODO Fix recursive types
// func ExampleThrowErrorOnCycle() {
// 	type CyclicType0 struct {
// 		CyclicField *struct {
// 			CyclicField *CyclicType0 `json:"b"`
// 		} `json:"a"`
// 	}

// 	schemas := make(specgen.Schemas)
// 	schemaRef, err := specgen.NewSchemaRefForValue(&CyclicType0{}, schemas, specgen.ThrowErrorOnCycle())
// 	if schemaRef != nil || err == nil {
// 		panic(`With option ThrowErrorOnCycle, an error is returned when a schema reference cycle is found`)
// 	}
// 	if _, ok := err.(*specgen.CycleError); !ok {
// 		panic(`With option ThrowErrorOnCycle, an error of type CycleError is returned`)
// 	}
// 	if len(schemas) != 0 {
// 		panic(`No references should have been collected at this point`)
// 	}

// 	if schemaRef, err = specgen.NewSchemaRefForValue(&CyclicType0{}, schemas); err != nil {
// 		panic(err)
// 	}

// 	var data []byte
// 	if data, err = json.MarshalIndent(schemaRef, "", "  "); err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("schemaRef: %s\n", data)
// 	if data, err = json.MarshalIndent(schemas, "", "  "); err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("schemas: %s\n", data)
// 	// Output:
// 	// schemaRef: {
// 	//   "properties": {
// 	//     "a": {
// 	//       "properties": {
// 	//         "b": {
// 	//           "$ref": "#/components/schemas/CyclicType0"
// 	//         }
// 	//       },
// 	//       "type": "object"
// 	//     }
// 	//   },
// 	//   "type": "object"
// 	// }
// 	// schemas: {
// 	//   "CyclicType0": {
// 	//     "properties": {
// 	//       "a": {
// 	//         "properties": {
// 	//           "b": {
// 	//             "$ref": "#/components/schemas/CyclicType0"
// 	//           }
// 	//         },
// 	//         "type": "object"
// 	//       }
// 	//     },
// 	//     "type": "object"
// 	//   }
// 	// }
// }

func TestExportedNonTagged(t *testing.T) {
	type Bla struct {
		A          string
		Another    string `json:"another"`
		yetAnother string // unused because unexported
		EvenAYaml  string `yaml:"even_a_yaml"`
	}

	schemaRef, err := specgen.NewSchemaRefForValue(&Bla{}, nil, specgen.UseAllExportedFields())
	require.NoError(t, err)
	require.Equal(t, &spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type: []string{"object"},
			Properties: map[string]spec.Schema{
				"A":           {SchemaProps: spec.SchemaProps{Type: []string{"string"}}},
				"another":     {SchemaProps: spec.SchemaProps{Type: []string{"string"}}},
				"even_a_yaml": {SchemaProps: spec.SchemaProps{Type: []string{"string"}}},
			},
		}}, schemaRef)
}

func ExampleUseAllExportedFields() {
	type UnsignedIntStruct struct {
		UnsignedInt uint `json:"uint"`
	}

	schemaRef, err := specgen.NewSchemaRefForValue(&UnsignedIntStruct{}, nil, specgen.UseAllExportedFields())
	if err != nil {
		panic(err)
	}

	var data []byte
	if data, err = json.MarshalIndent(schemaRef, "", "  "); err != nil {
		panic(err)
	}
	fmt.Printf("schemaRef: %s\n", data)
	// Output:
	// schemaRef: {
	//   "type": "object",
	//   "properties": {
	//     "uint": {
	//       "type": "integer",
	//       "minimum": 0
	//     }
	//   }
	// }
}

func ExampleGenerator_GenerateSchemaRef() {
	type EmbeddedStruct struct {
		ID string
	}

	type ContainerStruct struct {
		Name string
		EmbeddedStruct
	}

	instance := &ContainerStruct{
		Name: "Container",
		EmbeddedStruct: EmbeddedStruct{
			ID: "Embedded",
		},
	}

	generator := specgen.NewGenerator(specgen.UseAllExportedFields())

	schemaRef, err := generator.GenerateSchemaRef(reflect.TypeOf(instance))
	if err != nil {
		panic(err)
	}

	var data []byte
	if data, err = json.MarshalIndent(schemaRef.Properties["Name"], "", "  "); err != nil {
		panic(err)
	}
	fmt.Printf(`schemaRef.Value.Properties["Name"].Value: %s`, data)
	fmt.Println()
	if data, err = json.MarshalIndent(schemaRef.Properties["ID"], "", "  "); err != nil {
		panic(err)
	}
	fmt.Printf(`schemaRef.Value.Properties["ID"].Value: %s`, data)
	fmt.Println()
	// Output:
	// schemaRef.Value.Properties["Name"].Value: {
	//   "type": "string"
	// }
	// schemaRef.Value.Properties["ID"].Value: {
	//   "type": "string"
	// }
}

func TestEmbeddedPointerStructs(t *testing.T) {
	type EmbeddedStruct struct {
		ID string
	}

	type ContainerStruct struct {
		Name string
		*EmbeddedStruct
	}

	instance := &ContainerStruct{
		Name: "Container",
		EmbeddedStruct: &EmbeddedStruct{
			ID: "Embedded",
		},
	}

	generator := specgen.NewGenerator(specgen.UseAllExportedFields())

	schemaRef, err := generator.GenerateSchemaRef(reflect.TypeOf(instance))
	require.NoError(t, err)

	var ok bool
	_, ok = schemaRef.Properties["Name"]
	require.Equal(t, true, ok)

	_, ok = schemaRef.Properties["ID"]
	require.Equal(t, true, ok)
}

func TestCyclicReferences(t *testing.T) {
	type ObjectDiff struct {
		FieldCycle *ObjectDiff
		SliceCycle []*ObjectDiff
		MapCycle   map[*ObjectDiff]*ObjectDiff
	}

	instance := &ObjectDiff{
		FieldCycle: nil,
		SliceCycle: nil,
		MapCycle:   nil,
	}

	generator := specgen.NewGenerator(specgen.UseAllExportedFields())

	schemaRef, err := generator.GenerateSchemaRef(reflect.TypeOf(instance))
	require.NoError(t, err)

	s := (schemaRef.Properties["FieldCycle"].Ref)

	require.NotNil(t, schemaRef.Properties["FieldCycle"])
	require.Equal(t, "#/components/schemas/ObjectDiff", s.String())

	require.NotNil(t, schemaRef.Properties["SliceCycle"])
	require.Equal(t, "array", schemaRef.Properties["SliceCycle"].Type[0])
	require.Equal(t, "#/components/schemas/ObjectDiff", schemaRef.Properties["SliceCycle"].Items.Schema.Ref.String())

	require.NotNil(t, schemaRef.Properties["MapCycle"])
	require.Equal(t, "object", schemaRef.Properties["MapCycle"].Type[0])
	require.Equal(t, "#/components/schemas/ObjectDiff", schemaRef.Properties["MapCycle"].AdditionalProperties.Schema.Ref.String())
}

func ExampleSchemaCustomizer() {
	type NestedInnerBla struct {
		Enum1Field string `json:"enum1" myenumtag:"a,b"`
	}

	type InnerBla struct {
		UntaggedStringField string
		AnonStruct          struct {
			InnerFieldWithoutTag int
			InnerFieldWithTag    int `mymintag:"-1" mymaxtag:"50"`
			NestedInnerBla
		}
		Enum2Field string `json:"enum2" myenumtag:"c,d"`
	}

	type Bla struct {
		InnerBla
		EnumField3 string `json:"enum3" myenumtag:"e,f"`
	}

	customizer := specgen.SchemaCustomizer(func(name string, ft reflect.Type, tag reflect.StructTag, schema *spec.Schema) error {
		if tag.Get("mymintag") != "" {
			minVal, err := strconv.ParseFloat(tag.Get("mymintag"), 64)
			if err != nil {
				return err
			}
			schema.Minimum = &minVal
		}
		if tag.Get("mymaxtag") != "" {
			maxVal, err := strconv.ParseFloat(tag.Get("mymaxtag"), 64)
			if err != nil {
				return err
			}
			schema.Maximum = &maxVal
		}
		if tag.Get("myenumtag") != "" {
			for _, s := range strings.Split(tag.Get("myenumtag"), ",") {
				schema.Enum = append(schema.Enum, s)
			}
		}
		return nil
	})

	schemaRef, err := specgen.NewSchemaRefForValue(&Bla{}, nil, specgen.UseAllExportedFields(), customizer)
	if err != nil {
		panic(err)
	}

	var data []byte
	if data, err = json.MarshalIndent(schemaRef, "", "  "); err != nil {
		panic(err)
	}
	fmt.Printf("schemaRef: %s\n", data)
	// Output:
	// schemaRef: {
	//   "type": "object",
	//   "properties": {
	//     "AnonStruct": {
	//       "type": "object",
	//       "properties": {
	//         "InnerFieldWithTag": {
	//           "type": "integer",
	//           "maximum": 50,
	//           "minimum": -1
	//         },
	//         "InnerFieldWithoutTag": {
	//           "type": "integer"
	//         },
	//         "enum1": {
	//           "type": "string",
	//           "enum": [
	//             "a",
	//             "b"
	//           ]
	//         }
	//       }
	//     },
	//     "UntaggedStringField": {
	//       "type": "string"
	//     },
	//     "enum2": {
	//       "type": "string",
	//       "enum": [
	//         "c",
	//         "d"
	//       ]
	//     },
	//     "enum3": {
	//       "type": "string",
	//       "enum": [
	//         "e",
	//         "f"
	//       ]
	//     }
	//   }
	// }
}

func TestSchemaCustomizerError(t *testing.T) {
	customizer := specgen.SchemaCustomizer(func(name string, ft reflect.Type, tag reflect.StructTag, schema *spec.Schema) error {
		return errors.New("test error")
	})

	type Bla struct{}
	_, err := specgen.NewSchemaRefForValue(&Bla{}, nil, specgen.UseAllExportedFields(), customizer)
	require.EqualError(t, err, "test error")
}

// TODO Fix recursive types
// func ExampleNewSchemaRefForValue_recursive() {
// 	type RecursiveType struct {
// 		Field1     string           `json:"field1"`
// 		Field2     string           `json:"field2"`
// 		Field3     string           `json:"field3"`
// 		Components []*RecursiveType `json:"children,omitempty"`
// 	}

// 	schemas := make(specgen.Schemas)
// 	schemaRef, err := specgen.NewSchemaRefForValue(&RecursiveType{}, schemas)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var data []byte
// 	if data, err = json.MarshalIndent(&schemas, "", "  "); err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("schemas: %s\n", data)
// 	if data, err = json.MarshalIndent(&schemaRef, "", "  "); err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("schemaRef: %s\n", data)
// 	// Output:
// 	// schemas: {
// 	//   "RecursiveType": {
// 	//     "properties": {
// 	//       "children": {
// 	//         "items": {
// 	//           "$ref": "#/components/schemas/RecursiveType"
// 	//         },
// 	//         "type": "array"
// 	//       },
// 	//       "field1": {
// 	//         "type": "string"
// 	//       },
// 	//       "field2": {
// 	//         "type": "string"
// 	//       },
// 	//       "field3": {
// 	//         "type": "string"
// 	//       }
// 	//     },
// 	//     "type": "object"
// 	//   }
// 	// }
// 	// schemaRef: {
// 	//   "properties": {
// 	//     "children": {
// 	//       "items": {
// 	//         "$ref": "#/components/schemas/RecursiveType"
// 	//       },
// 	//       "type": "array"
// 	//     },
// 	//     "field1": {
// 	//       "type": "string"
// 	//     },
// 	//     "field2": {
// 	//       "type": "string"
// 	//     },
// 	//     "field3": {
// 	//       "type": "string"
// 	//     }
// 	//   },
// 	//   "type": "object"
// 	// }
// }

type Reffer struct {
	Ref *spec.Ref `json:"-"`
}

func TestT(t *testing.T) {
	t.Skip()
	// ref2 := spec.MustCreateRef("zzzz")
	// ref := spec.MustCreateRef("zzzz")
	// ref := spec.RefSchema("zzzz")
	// ref := spec.Ref{Ref: jsonreference.MustCreateRef("aaa")}
	// ref := spec.SchemaProps{Ref: spec.Ref{Ref: jsonreference.MustCreateRef("aaa")}}
	ref := &spec.Schema{SchemaProps: spec.SchemaProps{Ref: spec.Ref{Ref: jsonreference.MustCreateRef("aaa")}}}

	b1, _ := json.Marshal(ref.SchemaProps)
	fmt.Println(string(b1))
	// ref2 := spec.RefSchema("zzzz")
	// ref := &spec.Ref{}

	// ref := &Reffer{
	// 	Ref: &ref2,
	// }

	chai.LogYAML(ref)
	chai.LogJSON(ref)
}

func TestT2(t *testing.T) {
	t.Skip()
	// ref := spec.MustCreateRef("zzzz")
	// ref := spec.RefSchema("zzzz")

	gen := specgen.NewGenerator()
	schemas := make(map[string]spec.Schema)

	ref, err := gen.NewSchemaRefForValue(&model.Account{}, schemas)
	require.NoError(t, err)

	gen2 := specgen.NewGenerator()
	schemas2 := make(specgen.Schemas)

	ref2, err := gen2.NewSchemaRefForValue(&model.Account{}, schemas2)
	require.NoError(t, err)

	chai.LogYAML(ref2)
	chai.LogYAML(ref)
}

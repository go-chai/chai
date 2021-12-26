package specgen_test

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-chai/chai/specgen"
)

type (
	SomeStruct struct {
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

	SomeOtherType string
)

func Example() {
	schemaRef, err := specgen.NewSchemaRefForValue(&SomeStruct{}, nil)
	if err != nil {
		panic(err)
	}

	data, err := json.MarshalIndent(schemaRef, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", data)
	// Output:
	// {
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

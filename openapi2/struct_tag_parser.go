package openapi2

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
	"github.com/swaggo/swag"
)

const (
	// CamelCase indicates using CamelCase strategy for struct field.
	CamelCase = "camelcase"

	// PascalCase indicates using PascalCase strategy for struct field.
	PascalCase = "pascalcase"

	// SnakeCase indicates using SnakeCase strategy for struct field.
	SnakeCase = "snakecase"

	acceptAttr       = "@accept"
	produceAttr      = "@produce"
	xCodeSamplesAttr = "@x-codesamples"
	scopeAttrPrefix  = "@scope."
)

type tagBaseFieldParser struct {
	p *swag.Parser
	// field reflect.StructField
	tag  reflect.StructTag
	name string
}

func newTagBaseFieldParser(name string, p *swag.Parser, tag reflect.StructTag) *tagBaseFieldParser {
	ps := &tagBaseFieldParser{
		name: name,
		p:    p,
		tag:  tag,
	}

	return ps
}

func (ps *tagBaseFieldParser) ShouldSkip() (bool, error) {
	// // Skip non-exported fields.
	// if !ast.IsExported(ps.field.Names[0].Name) {
	// 	return true, nil
	// }

	// if ps.field.IsExported() {
	// 	return true, nil
	// }

	ignoreTag := ps.tag.Get("swaggerignore")
	if strings.EqualFold(ignoreTag, "true") {
		return true, nil
	}

	// json:"tag,hoge"
	name := strings.TrimSpace(strings.Split(ps.tag.Get("json"), ",")[0])
	if name == "-" {
		return true, nil
	}

	return false, nil
}

func (ps *tagBaseFieldParser) FieldName() (string, error) {
	var name string
	// json:"tag,hoge"
	name = strings.TrimSpace(strings.Split(ps.tag.Get("json"), ",")[0])

	if name != "" {
		return name, nil
	}

	switch ps.p.PropNamingStrategy {
	case SnakeCase:
		return toSnakeCase(ps.name), nil
	case PascalCase:
		return ps.name, nil
	default:
		return toLowerCamelCase(ps.name), nil
	}
}

func toSnakeCase(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) &&
			((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}

func toLowerCamelCase(in string) string {
	runes := []rune(in)

	var out []rune
	flag := false
	for i, curr := range runes {
		if (i == 0 && unicode.IsUpper(curr)) || (flag && unicode.IsUpper(curr)) {
			out = append(out, unicode.ToLower(curr))
			flag = true
		} else {
			out = append(out, curr)
			flag = false
		}
	}

	return string(out)
}

func (ps *tagBaseFieldParser) CustomSchema() (*spec.Schema, error) {
	typeTag := ps.tag.Get("swaggertype")
	if typeTag != "" {
		return swag.BuildCustomSchema(strings.Split(typeTag, ","))
	}

	return nil, nil
}

type structField struct {
	desc         string
	schemaType   string
	arrayType    string
	formatType   string
	maximum      *float64
	minimum      *float64
	multipleOf   *float64
	maxLength    *int64
	minLength    *int64
	maxItems     *int64
	minItems     *int64
	exampleValue interface{}
	defaultValue interface{}
	extensions   map[string]interface{}
	enums        []interface{}
	readOnly     bool
	unique       bool
}

// defineType enum value define the type (object and array unsupported).
func defineType(schemaType string, value string) (v interface{}, err error) {
	schemaType = swag.TransToValidSchemeType(schemaType)
	switch schemaType {
	case swag.STRING:
		return value, nil
	case swag.NUMBER:
		v, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("enum value %s can't convert to %s err: %s", value, schemaType, err)
		}
	case swag.INTEGER:
		v, err = strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("enum value %s can't convert to %s err: %s", value, schemaType, err)
		}
	case swag.BOOLEAN:
		v, err = strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("enum value %s can't convert to %s err: %s", value, schemaType, err)
		}
	default:
		return nil, errors.Errorf("%s is unsupported type in enum value %s", schemaType, value)
	}

	return v, nil
}

// defineTypeOfExample example value define the type (object and array unsupported)
func defineTypeOfExample(schemaType, arrayType, exampleValue string) (interface{}, error) {
	switch schemaType {
	case swag.STRING, swag.ANY:
		return exampleValue, nil
	case swag.NUMBER:
		v, err := strconv.ParseFloat(exampleValue, 64)
		if err != nil {
			return nil, fmt.Errorf("example value %s can't convert to %s err: %s", exampleValue, schemaType, err)
		}

		return v, nil
	case swag.INTEGER:
		v, err := strconv.Atoi(exampleValue)
		if err != nil {
			return nil, fmt.Errorf("example value %s can't convert to %s err: %s", exampleValue, schemaType, err)
		}

		return v, nil
	case swag.BOOLEAN:
		v, err := strconv.ParseBool(exampleValue)
		if err != nil {
			return nil, fmt.Errorf("example value %s can't convert to %s err: %s", exampleValue, schemaType, err)
		}

		return v, nil
	case swag.ARRAY:
		values := strings.Split(exampleValue, ",")
		result := make([]interface{}, 0)
		for _, value := range values {
			v, err := defineTypeOfExample(arrayType, "", value)
			if err != nil {
				return nil, err
			}
			result = append(result, v)
		}

		return result, nil
	case swag.OBJECT:
		if arrayType == "" {
			return nil, errors.Errorf("%s is unsupported type in example value `%s`", schemaType, exampleValue)
		}

		values := strings.Split(exampleValue, ",")
		result := map[string]interface{}{}
		for _, value := range values {
			mapData := strings.Split(value, ":")

			if len(mapData) == 2 {
				v, err := defineTypeOfExample(arrayType, "", mapData[1])
				if err != nil {
					return nil, err
				}
				result[mapData[0]] = v
			} else {
				return nil, fmt.Errorf("example value %s should format: key:value", exampleValue)
			}
		}

		return result, nil
	}

	return nil, errors.Errorf("%s is unsupported type in example value %s", schemaType, exampleValue)
}

func (ps *tagBaseFieldParser) ComplementSchema(schema *spec.Schema) error {
	types := ps.p.GetSchemaTypePath(schema, 2)
	if len(types) == 0 {
		return fmt.Errorf("invalid type for field: %s", ps.name)
	}

	structField := &structField{
		schemaType: types[0],
		formatType: ps.tag.Get("format"),
		readOnly:   ps.tag.Get("readonly") == "true",
	}

	if len(types) > 1 && (types[0] == swag.ARRAY || types[0] == swag.OBJECT) {
		structField.arrayType = types[1]
	}

	// if ps.field.Doc != nil {
	// 	structField.desc = strings.TrimSpace(ps.field.Doc.Text())
	// }
	// if structField.desc == "" && ps.field.Comment != nil {
	// 	structField.desc = strings.TrimSpace(ps.field.Comment.Text())
	// }

	structField.desc = ps.tag.Get("description")

	jsonTag := ps.tag.Get("json")
	// json:"name,string" or json:",string"

	exampleTag := ps.tag.Get("example")
	if exampleTag != "" {
		structField.exampleValue = exampleTag
		if !strings.Contains(jsonTag, ",string") {
			example, err := defineTypeOfExample(structField.schemaType, structField.arrayType, exampleTag)
			if err != nil {
				return err
			}
			structField.exampleValue = example
		}
	}

	bindingTag := ps.tag.Get("binding")
	if bindingTag != "" {
		ps.parseValidTags(bindingTag, structField)
	}

	validateTag := ps.tag.Get("validate")
	if validateTag != "" {
		ps.parseValidTags(validateTag, structField)
	}

	extensionsTag := ps.tag.Get("extensions")
	if extensionsTag != "" {
		structField.extensions = map[string]interface{}{}
		for _, val := range strings.Split(extensionsTag, ",") {
			parts := strings.SplitN(val, "=", 2)
			if len(parts) == 2 {
				structField.extensions[parts[0]] = parts[1]
			} else {
				if len(parts[0]) > 0 && string(parts[0][0]) == "!" {
					structField.extensions[parts[0][1:]] = false
				} else {
					structField.extensions[parts[0]] = true
				}
			}
		}
	}

	enumsTag := ps.tag.Get("enums")
	if enumsTag != "" {
		enumType := structField.schemaType
		if structField.schemaType == swag.ARRAY {
			enumType = structField.arrayType
		}

		structField.enums = nil
		for _, e := range strings.Split(enumsTag, ",") {
			value, err := defineType(enumType, e)
			if err != nil {
				return err
			}
			structField.enums = append(structField.enums, value)
		}
	}

	defaultTag := ps.tag.Get("default")
	if defaultTag != "" {
		value, err := defineType(structField.schemaType, defaultTag)
		if err != nil {
			return err
		}
		structField.defaultValue = value
	}

	if swag.IsNumericType(structField.schemaType) || swag.IsNumericType(structField.arrayType) {
		maximum, err := getFloatTag(ps.tag, "maximum")
		if err != nil {
			return err
		}
		if maximum != nil {
			structField.maximum = maximum
		}

		minimum, err := getFloatTag(ps.tag, "minimum")
		if err != nil {
			return err
		}
		if minimum != nil {
			structField.minimum = minimum
		}

		multipleOf, err := getFloatTag(ps.tag, "multipleOf")
		if err != nil {
			return err
		}
		if multipleOf != nil {
			structField.multipleOf = multipleOf
		}
	}

	if structField.schemaType == swag.STRING || structField.arrayType == swag.STRING {
		maxLength, err := getIntTag(ps.tag, "maxLength")
		if err != nil {
			return err
		}
		if maxLength != nil {
			structField.maxLength = maxLength
		}

		minLength, err := getIntTag(ps.tag, "minLength")
		if err != nil {
			return err
		}
		if minLength != nil {
			structField.minLength = minLength
		}
	}

	// perform this after setting everything else (min, max, etc...)
	if strings.Contains(jsonTag, ",string") { // @encoding/json: "It applies only to fields of string, floating point, integer, or boolean types."
		defaultValues := map[string]string{
			// Zero Values as string
			swag.STRING:  "",
			swag.INTEGER: "0",
			swag.BOOLEAN: "false",
			swag.NUMBER:  "0",
		}

		defaultValue, ok := defaultValues[structField.schemaType]
		if ok {
			structField.schemaType = swag.STRING

			if structField.exampleValue == nil {
				// if exampleValue is not defined by the user,
				// we will force an example with a correct value
				// (eg: int->"0", bool:"false")
				structField.exampleValue = defaultValue
			}
		}
	}

	if structField.schemaType == swag.STRING && types[0] != swag.STRING {
		*schema = *swag.PrimitiveSchema(structField.schemaType)
	}

	schema.Description = structField.desc
	schema.ReadOnly = structField.readOnly
	if !reflect.ValueOf(schema.Ref).IsZero() && schema.ReadOnly {
		schema.AllOf = []spec.Schema{*spec.RefSchema(schema.Ref.String())}
		schema.Ref = spec.Ref{} // clear out existing ref
	}
	schema.Default = structField.defaultValue
	schema.Example = structField.exampleValue
	if structField.schemaType != swag.ARRAY {
		schema.Format = structField.formatType
	}
	schema.Extensions = structField.extensions
	eleSchema := schema
	if structField.schemaType == swag.ARRAY {
		// For Array only
		schema.MaxItems = structField.maxItems
		schema.MinItems = structField.minItems
		schema.UniqueItems = structField.unique

		eleSchema = schema.Items.Schema
		eleSchema.Format = structField.formatType
	}
	eleSchema.Maximum = structField.maximum
	eleSchema.Minimum = structField.minimum
	eleSchema.MultipleOf = structField.multipleOf
	eleSchema.MaxLength = structField.maxLength
	eleSchema.MinLength = structField.minLength
	eleSchema.Enum = structField.enums

	return nil
}

func getFloatTag(structTag reflect.StructTag, tagName string) (*float64, error) {
	strValue := structTag.Get(tagName)
	if strValue == "" {
		return nil, nil
	}

	value, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return nil, fmt.Errorf("can't parse numeric value of %q tag: %v", tagName, err)
	}

	return &value, nil
}

func getIntTag(structTag reflect.StructTag, tagName string) (*int64, error) {
	strValue := structTag.Get(tagName)
	if strValue == "" {
		return nil, nil
	}

	value, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("can't parse numeric value of %q tag: %v", tagName, err)
	}

	return &value, nil
}

func (ps *tagBaseFieldParser) IsRequired() (bool, error) {
	bindingTag := ps.tag.Get("binding")
	if bindingTag != "" {
		for _, val := range strings.Split(bindingTag, ",") {
			if val == "required" {
				return true, nil
			}
		}
	}

	validateTag := ps.tag.Get("validate")
	if validateTag != "" {
		for _, val := range strings.Split(validateTag, ",") {
			if val == "required" {
				return true, nil
			}
		}
	}

	return false, nil
}

func (ps *tagBaseFieldParser) parseValidTags(validTag string, sf *structField) {
	// `validate:"required,max=10,min=1"`
	// ps. required checked by IsRequired().
	for _, val := range strings.Split(validTag, ",") {
		var (
			valKey   string
			valValue string
		)
		kv := strings.Split(val, "=")
		switch len(kv) {
		case 1:
			valKey = kv[0]
		case 2:
			valKey = kv[0]
			valValue = kv[1]
		default:
			continue
		}
		valValue = strings.Replace(strings.Replace(valValue, utf8HexComma, ",", -1), utf8Pipe, "|", -1)

		switch valKey {
		case "max", "lte":
			sf.setMax(valValue)
		case "min", "gte":
			sf.setMin(valValue)
		case "oneof":
			sf.setOneOf(valValue)
		case "unique":
			if sf.schemaType == swag.ARRAY {
				sf.unique = true
			}
		case "dive":
			// ignore dive
			return
		default:
			continue
		}
	}
}

func (sf *structField) setOneOf(valValue string) {
	if len(sf.enums) != 0 {
		return
	}

	enumType := sf.schemaType
	if sf.schemaType == swag.ARRAY {
		enumType = sf.arrayType
	}

	valValues := parseOneOfParam2(valValue)
	for i := range valValues {
		value, err := defineType(enumType, valValues[i])
		if err != nil {
			continue
		}
		sf.enums = append(sf.enums, value)
	}
}

func (sf *structField) setMin(valValue string) {
	value, err := strconv.ParseFloat(valValue, 64)
	if err != nil {
		return
	}
	switch sf.schemaType {
	case swag.INTEGER, swag.NUMBER:
		sf.minimum = &value
	case swag.STRING:
		intValue := int64(value)
		sf.minLength = &intValue
	case swag.ARRAY:
		intValue := int64(value)
		sf.minItems = &intValue
	}
}

func (sf *structField) setMax(valValue string) {
	value, err := strconv.ParseFloat(valValue, 64)
	if err != nil {
		return
	}
	switch sf.schemaType {
	case swag.INTEGER, swag.NUMBER:
		sf.maximum = &value
	case swag.STRING:
		intValue := int64(value)
		sf.maxLength = &intValue
	case swag.ARRAY:
		intValue := int64(value)
		sf.maxItems = &intValue
	}
}

const (
	utf8HexComma = "0x2C"
	utf8Pipe     = "0x7C"
)

// These code copy from
// https://github.com/go-playground/validator/blob/d4271985b44b735c6f76abc7a06532ee997f9476/baked_in.go#L207
// ---
var oneofValsCache = map[string][]string{}
var oneofValsCacheRWLock = sync.RWMutex{}
var splitParamsRegex = regexp.MustCompile(`'[^']*'|\S+`)

func parseOneOfParam2(s string) []string {
	oneofValsCacheRWLock.RLock()
	values, ok := oneofValsCache[s]
	oneofValsCacheRWLock.RUnlock()
	if !ok {
		oneofValsCacheRWLock.Lock()
		values = splitParamsRegex.FindAllString(s, -1)
		for i := 0; i < len(values); i++ {
			values[i] = strings.Replace(values[i], "'", "", -1)
		}
		oneofValsCache[s] = values
		oneofValsCacheRWLock.Unlock()
	}
	return values
}

// ---

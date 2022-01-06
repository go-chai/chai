package openapi3

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
	"github.com/swaggo/swag"
)

var open = os.Open

// DefaultOverridesFile is the location swaggo will look for type overrides.
const DefaultOverridesFile = ".swaggo"

// Gen presents a generate tool for swag.
type Gen struct {
	jsonIndent func(data interface{}) ([]byte, error)
	jsonToYAML func(data []byte) ([]byte, error)
}

// New creates a new Gen.
func NewGen() *Gen {
	return &Gen{
		jsonIndent: func(data interface{}) ([]byte, error) {
			return json.MarshalIndent(data, "", "    ")
		},
		jsonToYAML: yaml.JSONToYAML,
	}
}

// Config presents Gen configurations.
type Config struct {
	// SearchDir the swag would be parse,comma separated if multiple
	SearchDir string

	// excludes dirs and files in SearchDir,comma separated
	Excludes string

	// OutputDir represents the output directory for all the generated files
	OutputDir string

	// MainAPIFile the Go file path in which 'swagger general API Info' is written
	MainAPIFile string

	// PropNamingStrategy represents property naming strategy like snake case,camel case,pascal case
	PropNamingStrategy string

	// MarkdownFilesDir used to find markdown files, which can be used for tag descriptions
	MarkdownFilesDir string

	// CodeExampleFilesDir used to find code example files, which can be used for x-codeSamples
	CodeExampleFilesDir string

	// InstanceName is used to get distinct names for different swagger documents in the
	// same project. The default value is "swagger".
	InstanceName string

	// ParseDepth dependency parse depth
	ParseDepth int

	// ParseVendor whether swag should be parse vendor folder
	ParseVendor bool

	// ParseDependencies whether swag should be parse outside dependency folder
	ParseDependency bool

	// ParseInternal whether swag should parse internal packages
	ParseInternal bool

	// Strict whether swag should error or warn when it detects cases which are most likely user errors
	Strict bool

	// GeneratedTime whether swag should generate the timestamp at the top of docs.go
	GeneratedTime bool

	// OverridesFile defines global type overrides.
	OverridesFile string
}

// Generate outputs a swagger spec
func (g *Gen) Generate(swagger *openapi3.T, config *GenConfig) error {
	if config.InstanceName == "" {
		config.InstanceName = swag.Name
	}

	if config.OutputDir == "" {
		config.OutputDir = "docs/"
	}

	b, err := g.jsonIndent(swagger)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(config.OutputDir, os.ModePerm); err != nil {
		return err
	}

	absOutputDir, err := filepath.Abs(config.OutputDir)
	if err != nil {
		return err
	}
	packageName := filepath.Base(absOutputDir)
	docFileName := filepath.Join(config.OutputDir, "docs.go")
	jsonFileName := filepath.Join(config.OutputDir, "swagger.json")
	yamlFileName := filepath.Join(config.OutputDir, "swagger.yaml")

	docs, err := os.Create(docFileName)
	if err != nil {
		return err
	}
	defer docs.Close()

	err = g.writeFile(b, jsonFileName)
	if err != nil {
		return err
	}

	y, err := g.jsonToYAML(b)
	if err != nil {
		return fmt.Errorf("cannot convert json to yaml error: %s", err)
	}

	err = g.writeFile(y, yamlFileName)
	if err != nil {
		return err
	}

	// Write doc
	err = g.writeGoDoc(packageName, docs, swagger, config)
	if err != nil {
		return err
	}

	log.Printf("create docs.go at %+v", docFileName)
	log.Printf("create swagger.json at %+v", jsonFileName)
	log.Printf("create swagger.yaml at %+v", yamlFileName)

	return nil
}

func (g *Gen) writeFile(b []byte, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	return err
}

func (g *Gen) formatSource(src []byte) []byte {
	code, err := format.Source(src)
	if err != nil {
		code = src // Output the unformatted code anyway
	}
	return code
}

// Read the swaggo overrides
func parseOverrides(r io.Reader) (map[string]string, error) {
	overrides := make(map[string]string)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments
		if len(line) > 1 && line[0:2] == "//" {
			continue
		}

		parts := strings.Fields(line)

		switch len(parts) {
		case 0:
			// only whitespace
			continue
		case 2:
			// either a skip or malformed
			if parts[0] != "skip" {
				return nil, fmt.Errorf("could not parse override: '%s'", line)
			}
			overrides[parts[1]] = ""
		case 3:
			// either a replace or malformed
			if parts[0] != "replace" {
				return nil, fmt.Errorf("could not parse override: '%s'", line)
			}
			overrides[parts[1]] = parts[2]
		default:
			return nil, fmt.Errorf("could not parse override: '%s'", line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading overrides file: %w", err)
	}

	return overrides, nil
}

func (g *Gen) writeGoDoc(packageName string, output io.Writer, swagger *openapi3.T, config *GenConfig) error {
	generator, err := template.New("swagger_info").Funcs(template.FuncMap{
		"printDoc": func(v string) string {
			// Add servers
			v = "{\n    \"servers\": {{ marshal .Servers }}," + v[1:]
			// Sanitize backticks
			return strings.Replace(v, "`", "`+\"`\"+`", -1)
		},
	}).Parse(packageTemplate)
	if err != nil {
		return err
	}

	swaggerSpec := &openapi3.T{
		ExtensionProps: swagger.ExtensionProps,
		OpenAPI:        swagger.OpenAPI,
		Components:     swagger.Components,
		Info: &openapi3.Info{
			ExtensionProps: swagger.Info.ExtensionProps,
			Description:    "{{escape .Description}}",
			Title:          "{{.Title}}",
			TermsOfService: swagger.Info.TermsOfService,
			Contact:        swagger.Info.Contact,
			License:        swagger.Info.License,
			Version:        "{{.Version}}",
		},
		Paths:        swagger.Paths,
		Security:     swagger.Security,
		Tags:         swagger.Tags,
		ExternalDocs: swagger.ExternalDocs,
	}

	// crafted docs.json
	buf, err := g.jsonIndent(swaggerSpec)
	if err != nil {
		return err
	}

	buffer := &bytes.Buffer{}
	err = generator.Execute(buffer, struct {
		Timestamp     time.Time
		Doc           string
		Servers       openapi3.Servers
		PackageName   string
		Title         string
		Description   string
		Version       string
		InstanceName  string
		GeneratedTime bool
	}{
		Timestamp:     time.Now(),
		GeneratedTime: config.GeneratedTime,
		Doc:           string(buf),
		Servers:       swagger.Servers,
		PackageName:   packageName,
		Title:         swagger.Info.Title,
		Description:   swagger.Info.Description,
		Version:       swagger.Info.Version,
		InstanceName:  config.InstanceName,
	})
	if err != nil {
		return err
	}

	code := g.formatSource(buffer.Bytes())

	// write
	_, err = output.Write(code)
	return err
}

var packageTemplate = `// Package {{.PackageName}} GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag{{ if .GeneratedTime }} at
// {{ .Timestamp }}{{ end }}
package {{.PackageName}}

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/swaggo/swag"
)

var doc = ` + "`{{ printDoc .Doc}}`" + `

type swaggerInfo struct {
	Version     string
	Servers		openapi3.Servers
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     {{ printf "%q" .Version}},
	Servers:     openapi3.Servers{ 
		{{ range $index, $server := .Servers}}{{if gt $index 0}},{{end}} &openapi3.Server{ 
			URL: {{ printf "%q" $server.URL }},
			Description: {{ printf "%q" $server.Description }},
		},
		{{end}} 
	},
	Title:       {{ printf "%q" .Title}},
	Description: {{ printf "%q" .Description}},
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register({{ printf "%q" .InstanceName }}, &s{})
}
`

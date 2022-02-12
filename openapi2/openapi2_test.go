package openapi2

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/go-chai/chai/chai"
	"github.com/go-chai/chai/internal/tests"
	"github.com/go-chai/swag"
	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAPIObjectSchema(t *testing.T) {
	type args struct {
		val any
	}
	type want struct {
		typeName        string
		ref             string
		definitionsJSON string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "string",
			args: args{
				val: "12345asdf",
			},
			want: want{typeName: "string"},
		},
		{
			name: "int",
			args: args{
				val: 123,
			},
			want: want{typeName: "integer"},
		},
		{
			name: "obj",
			args: args{
				val: &tests.TestStruct{},
			},
			want: want{ref: "#/definitions/tests.TestStruct", definitionsJSON: `{"tests.TestStruct": {"type": "object","properties": {"bar": {"type": "integer"},"foo": {"type": "string"}}}}`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := swag.New(swag.SetDebugger(log.Default()), func(p *swag.Parser) {
				p.ParseDependency = true
			})
			fi := getFuncInfo(RegisterRoute)
			op := swag.NewOperation(parser)

			err := parser.GetAllGoFileInfoAndParseTypes("../internal")
			require.NoError(t, err)

			schema, err := op.ParseAPIObjectSchema("object", typeName(tt.args.val), fi.ASTFile)
			require.NoError(t, err)

			LogYAML(parser.GetSwagger().Definitions)

			if tt.want.typeName != "" {
				require.Equal(t, tt.want.typeName, schema.Type[0])
			}

			if tt.want.ref != "" {
				require.Equal(t, tt.want.ref, schema.Ref.String())
			}

			if tt.want.definitionsJSON != "" {
				require.JSONEq(t, tt.want.definitionsJSON, js(parser.GetSwagger().Definitions))
			}
		})
	}
}

func TestMergeParameters(t *testing.T) {
	type args struct {
		params [][]spec.Parameter
	}
	tests := []struct {
		name string
		args args
		want []spec.Parameter
	}{
		{
			name: "test 1",
			args: args{
				params: [][]spec.Parameter{
					{
						{ParamProps: spec.ParamProps{Name: "p1", In: "path", Description: "d1", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t1", Format: "f1"}},
						{ParamProps: spec.ParamProps{Name: "p3", In: "path", Description: "d3", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t3", Format: "f3"}},
					},
					{
						{ParamProps: spec.ParamProps{Name: "p1", In: "path", Description: "d11", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t11", Format: "f11"}},
					},
				},
			},
			want: []spec.Parameter{
				{ParamProps: spec.ParamProps{Name: "p1", In: "path", Description: "d11", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t11", Format: "f11"}},
				{ParamProps: spec.ParamProps{Name: "p3", In: "path", Description: "d3", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t3", Format: "f3"}},
			},
		},
		{
			name: "t2",
			args: args{
				params: [][]spec.Parameter{
					{
						{ParamProps: spec.ParamProps{Name: "p1", In: "path", Description: "d1", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t1", Format: "f1"}},
						{ParamProps: spec.ParamProps{Name: "p3", In: "path", Description: "d3", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t3", Format: "f3"}},
					},
					{
						{ParamProps: spec.ParamProps{Name: "p1", In: "body", Description: "d11", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t11", Format: "f11"}},
					},
				},
			},
			want: []spec.Parameter{
				{ParamProps: spec.ParamProps{Name: "p1", In: "body", Description: "d11", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t11", Format: "f11"}},
				{ParamProps: spec.ParamProps{Name: "p1", In: "path", Description: "d1", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t1", Format: "f1"}},
				{ParamProps: spec.ParamProps{Name: "p3", In: "path", Description: "d3", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t3", Format: "f3"}},
			},
		},
		{
			name: "t3",
			args: args{
				params: [][]spec.Parameter{
					{
						{ParamProps: spec.ParamProps{Name: "p1", In: "path", Description: "d1", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t1", Format: "f1"}},
						{ParamProps: spec.ParamProps{Name: "p3", In: "path", Description: "d3", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t3", Format: "f3"}},
					},
					{{ParamProps: spec.ParamProps{Name: "p1", In: "body", Description: "d11", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t11", Format: "f11"}}},
				},
			},
			want: []spec.Parameter{
				{ParamProps: spec.ParamProps{Name: "p1", In: "body", Description: "d11", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t11", Format: "f11"}},
				{ParamProps: spec.ParamProps{Name: "p1", In: "path", Description: "d1", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t1", Format: "f1"}},
				{ParamProps: spec.ParamProps{Name: "p3", In: "path", Description: "d3", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t3", Format: "f3"}},
			},
		},
		{
			name: "t4",
			args: args{
				params: [][]spec.Parameter{
					{
						{ParamProps: spec.ParamProps{Name: "p1", In: "path", Description: "d1", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t1", Format: "f1"}},
						{ParamProps: spec.ParamProps{Name: "p3", In: "path", Description: "d3", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t3", Format: "f3"}},
					},
					{{ParamProps: spec.ParamProps{Name: "p1", In: "body", Description: "d11", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t11", Format: "f11"}}},
				},
			},
			want: []spec.Parameter{
				{ParamProps: spec.ParamProps{Name: "p1", In: "body", Description: "d11", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t11", Format: "f11"}},
				{ParamProps: spec.ParamProps{Name: "p1", In: "path", Description: "d1", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t1", Format: "f1"}},
				{ParamProps: spec.ParamProps{Name: "p3", In: "path", Description: "d3", Required: true}, SimpleSchema: spec.SimpleSchema{Type: "t3", Format: "f3"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantJSON := js(tt.want)
			got := mergeParameters(tt.args.params...)
			gotJSON := js(got)

			assert.JSONEq(t, string(wantJSON), string(gotJSON))
		})
	}
}

func js(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func TestAssociateBy(t *testing.T) {
	type args struct {
		ts []spec.Parameter
		fn func(spec.Parameter) pk
	}
	tests := []struct {
		name string
		args args
		want map[pk]spec.Parameter
	}{
		{
			name: "t1",
			args: args{
				ts: []spec.Parameter{
					{ParamProps: spec.ParamProps{Name: "p1", In: "path", Description: "d1", Required: true}},
					{ParamProps: spec.ParamProps{Name: "p2", In: "path", Description: "d2", Required: true}},
					{ParamProps: spec.ParamProps{Name: "p1", In: "body", Description: "d11", Required: true}},
				},
				fn: func(p spec.Parameter) pk {
					return pk{p.In, p.Name}
				},
			},
			want: map[pk]spec.Parameter{
				{In: "path", Name: "p1"}: {ParamProps: spec.ParamProps{Name: "p1", In: "path", Description: "d1", Required: true}},
				{In: "path", Name: "p2"}: {ParamProps: spec.ParamProps{Name: "p2", In: "path", Description: "d2", Required: true}},
				{In: "body", Name: "p1"}: {ParamProps: spec.ParamProps{Name: "p1", In: "body", Description: "d11", Required: true}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := associateBy(tt.args.ts, tt.args.fn)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSortedKeys(t *testing.T) {
	type args struct {
		m    map[pk]string
		less func(pk, pk) bool
	}
	tests := []struct {
		name string
		args args
		want []pk
	}{
		{
			name: "t1",
			args: args{
				m: map[pk]string{
					{"path", "p1"}: "1",
					{"path", "p2"}: "2",
					{"body", "p3"}: "3",
					{"body", "p2"}: "4",
				},
				less: less,
			},
			want: []pk{
				{"body", "p2"},
				{"body", "p3"},
				{"path", "p1"},
				{"path", "p2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortedKeys(tt.args.m, tt.args.less)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDocs(t *testing.T) {
	type args struct {
		routes []*Route
	}
	tests := []struct {
		name     string
		args     args
		filePath string
		wantErr  bool
	}{
		{
			name: "t1",
			args: args{
				routes: []*Route{
					{
						Method: "GET",
						Path:   "/test1/{p1}/{p2}",
						Params: []spec.Parameter{
							{ParamProps: spec.ParamProps{Name: "p1", In: "path", Description: "d1", Required: true}},
							{ParamProps: spec.ParamProps{Name: "p2", In: "path", Description: "d2", Required: true}},
						},

						// ShowBottle godoc
						// @Summary      Test Handler
						// @Description  get string by ID
						// @ID           get-string-by-int
						// @Tags         bottles
						// @Accept       json
						// @Produce      json
						// @Success      200
						// @Failure      400,404,500
						Handler: chai.NewReqResHandler(func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, error) {
							return nil, 0, nil
						}),
					},
				},
			},
			filePath: "testdata/t1.json",
			wantErr:  false,
		},
		{
			name: "t2",
			args: args{
				routes: []*Route{
					{
						Method: "GET",
						Path:   "/test1/{p1}/{p2}",
						Params: []spec.Parameter{
							{ParamProps: spec.ParamProps{Name: "p1", In: "path", Description: "d1", Required: true}},
							{ParamProps: spec.ParamProps{Name: "p2", In: "path", Description: "d2", Required: true}},
						},

						// ShowBottles godoc
						// @Summary      Test Handler
						// @Description  get string by ID
						// @ID           get-string-by-int
						// @Tags         bottles
						// @Success      200
						// @Failure      400,404,500
						Handler: chai.NewReqResHandler(func(req []*tests.TestRequest, w http.ResponseWriter, r *http.Request) ([]*tests.TestResponse, int, error) {
							return nil, 0, nil
						}),
					},
				},
			},
			filePath: "testdata/t2.json",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Docs(tt.args.routes)

			LogJSON(got)

			if (err != nil) != tt.wantErr {
				t.Errorf("Docs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.JSONEq(t, load(t, tt.filePath), js(got))
		})
	}
}

func load(t *testing.T, path string) string {
	b, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	return string(b)
}

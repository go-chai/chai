package openapi3

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chai/chai/chai"
	"github.com/go-chai/chai/internal/tests"
	"github.com/go-chai/chai/log"
	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMergeParameters(t *testing.T) {
	type args struct {
		params [][]*openapi3.ParameterRef
	}
	tests := []struct {
		name string
		args args
		want []*openapi3.ParameterRef
	}{
		{
			name: "test 1",
			args: args{
				params: [][]*openapi3.ParameterRef{
					{
						{Value: &openapi3.Parameter{Name: "p1", In: "path", Description: "d1", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t1", Format: "f1"}}}},
						{Value: &openapi3.Parameter{Name: "p3", In: "path", Description: "d3", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t3", Format: "f3"}}}},
					},
					{
						{Value: &openapi3.Parameter{Name: "p1", In: "path", Description: "d11", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t11", Format: "f11"}}}},
					},
				},
			},
			want: []*openapi3.ParameterRef{
				{Value: &openapi3.Parameter{Name: "p1", In: "path", Description: "d11", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t11", Format: "f11"}}}},
				{Value: &openapi3.Parameter{Name: "p3", In: "path", Description: "d3", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t3", Format: "f3"}}}},
			},
		},
		{
			name: "t2",
			args: args{
				params: [][]*openapi3.ParameterRef{
					{
						{Value: &openapi3.Parameter{Name: "p1", In: "path", Description: "d1", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t1", Format: "f1"}}}},
						{Value: &openapi3.Parameter{Name: "p3", In: "path", Description: "d3", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t3", Format: "f3"}}}},
					},
					{
						{Value: &openapi3.Parameter{Name: "p1", In: "body", Description: "d11", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t11", Format: "f11"}}}},
					},
				},
			},
			want: []*openapi3.ParameterRef{
				{Value: &openapi3.Parameter{Name: "p1", In: "body", Description: "d11", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t11", Format: "f11"}}}},
				{Value: &openapi3.Parameter{Name: "p1", In: "path", Description: "d1", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t1", Format: "f1"}}}},
				{Value: &openapi3.Parameter{Name: "p3", In: "path", Description: "d3", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t3", Format: "f3"}}}},
			},
		},
		{
			name: "t3",
			args: args{
				params: [][]*openapi3.ParameterRef{
					{
						{Value: &openapi3.Parameter{Name: "p1", In: "path", Description: "d1", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t1", Format: "f1"}}}},
						{Value: &openapi3.Parameter{Name: "p3", In: "path", Description: "d3", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t3", Format: "f3"}}}},
					},
					{
						{Value: &openapi3.Parameter{Name: "p1", In: "body", Description: "d11", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t11", Format: "f11"}}}},
					},
				},
			},
			want: []*openapi3.ParameterRef{
				{Value: &openapi3.Parameter{Name: "p1", In: "body", Description: "d11", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t11", Format: "f11"}}}},
				{Value: &openapi3.Parameter{Name: "p1", In: "path", Description: "d1", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t1", Format: "f1"}}}},
				{Value: &openapi3.Parameter{Name: "p3", In: "path", Description: "d3", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t3", Format: "f3"}}}},
			},
		},
		{
			name: "t4",
			args: args{
				params: [][]*openapi3.ParameterRef{
					{
						{Value: &openapi3.Parameter{Name: "p1", In: "path", Description: "d1", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t1", Format: "f1"}}}},
						{Value: &openapi3.Parameter{Name: "p3", In: "path", Description: "d3", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t3", Format: "f3"}}}},
					},
					{
						{Value: &openapi3.Parameter{Name: "p1", In: "body", Description: "d11", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t11", Format: "f11"}}}},
					},
				},
			},
			want: []*openapi3.ParameterRef{
				{Value: &openapi3.Parameter{Name: "p1", In: "body", Description: "d11", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t11", Format: "f11"}}}},
				{Value: &openapi3.Parameter{Name: "p1", In: "path", Description: "d1", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t1", Format: "f1"}}}},
				{Value: &openapi3.Parameter{Name: "p3", In: "path", Description: "d3", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "t3", Format: "f3"}}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantJSON := js(tt.want)
			got := mergeSlices(makeKey, cmpKeys, mergeParamsFn, tt.args.params...)
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
		fn func(spec.Parameter) key
	}
	tests := []struct {
		name string
		args args
		want map[key]spec.Parameter
	}{
		{
			name: "t1",
			args: args{
				ts: []spec.Parameter{
					{ParamProps: spec.ParamProps{Name: "p1", In: "path", Description: "d1", Required: true}},
					{ParamProps: spec.ParamProps{Name: "p2", In: "path", Description: "d2", Required: true}},
					{ParamProps: spec.ParamProps{Name: "p1", In: "body", Description: "d11", Required: true}},
				},
				fn: func(p spec.Parameter) key {
					return key{p.In, p.Name}
				},
			},
			want: map[key]spec.Parameter{
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
		m    map[key]string
		less func(key, key) bool
	}
	tests := []struct {
		name string
		args args
		want []key
	}{
		{
			name: "t1",
			args: args{
				m: map[key]string{
					{"path", "p1"}: "1",
					{"path", "p2"}: "2",
					{"body", "p3"}: "3",
					{"body", "p2"}: "4",
				},
				less: cmpKeys,
			},
			want: []key{
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
	h1 := chai.NewReqResHandler("", "", func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, error) {
		return nil, 0, nil
	}).
		Summary("Test Handler").
		Description("get string by ID").
		ID("get-string-by-int").
		Tags("bottles").
		ResponseCodes("test", 200, 400, 404, 500)

	h2 := chai.NewReqResHandler("", "", func(req *tests.TestRequests, w http.ResponseWriter, r *http.Request) ([]*tests.TestResponse, int, error) {
		return nil, 0, nil
	}).
		Summary("Test Handler").
		Description("get string by ID").
		ID("get-string-by-int").
		Tags("bottles").
		ResponseCodes("test", 200, 400, 404, 500)
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
						Method:   "POST",
						Path:     "/test1/{p1}/{p2}",
						Metadata: h1.GetMetadata(),
						Params: []*openapi3.ParameterRef{
							{Value: &openapi3.Parameter{Name: "p1", In: "path", Description: "d1", Required: true}},
							{Value: &openapi3.Parameter{Name: "p2", In: "path", Description: "d2", Required: true, Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}}}},
						},
						Handler: h1,
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
						Metadata: h2.GetMetadata(),
						Method:   "POST",
						Path:     "/test1/{p1}/{p2}",
						Params: []*openapi3.ParameterRef{
							{Value: &openapi3.Parameter{Name: "p1", In: "path", Description: "d1", Required: true}},
							{Value: &openapi3.Parameter{Name: "p2", In: "path", Description: "d2", Required: true}},
						},
						Handler: h2,
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

			log.JSON(got)

			if (err != nil) != tt.wantErr {
				t.Errorf("Docs() error = %+v, wantErr %v", err, tt.wantErr)
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

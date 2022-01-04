package openapi2

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

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
			LogYAML(tt.want)

			wantJSON := js(tt.want)

			got := mergeParameters(tt.args.params...)
			LogYAML(got)

			gotJSON := js(got)

			assert.Equal(t, string(wantJSON), string(gotJSON))

			fmt.Println("-----------------")
		})
	}
}

func js(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func TestAssoc(t *testing.T) {
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
			got := assoc(tt.args.ts, tt.args.fn)

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

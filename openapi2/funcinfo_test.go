// wrong comment 1

package openapi2

import (
	"fmt"
	"go/parser"
	"go/token"
	"net/http"
	"testing"

	// "github.com/go-chai/chai/chai"

	"github.com/go-chai/chai/openapi2/internal/tests"
	"github.com/go-chi/docgen"
	"github.com/stretchr/testify/require"
)

func TT() float32 {
	fifio := func(mm int) (string, int, string) {
		println(mm)

		if true {
			return "", http.StatusAccepted, "123"
		}
		return "45m", http.StatusForbidden, "4435"
	}

	fifio(123)

	return 1.23
}

type ZZ struct {
	A string `json:"a"`
}

func TestFixFuncLine(t *testing.T) {
	type args struct {
		filePath string
		line     int
	}
	type want struct {
		fixedLine int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "simple",
			args: args{
				filePath: "internal/tests/testfile.go",
				line:     24,
			},
			want: want{fixedLine: 12},
		},
		{
			name: "simple2",
			args: args{
				filePath: "internal/tests/testfile.go",
				line:     44,
			},
			want: want{fixedLine: 32},
		},
		{
			name: "not simple",
			args: args{
				filePath: "internal/tests/testfile.go",
				line:     52,
			},
			want: want{fixedLine: 52},
		},
		{
			name: "not simple 2",
			args: args{
				filePath: "internal/tests/testfile.go",
				line:     74,
			},
			want: want{fixedLine: 74},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			astFile, err := parser.ParseFile(fset, tt.args.filePath, nil, parser.ParseComments)
			require.NoError(t, err)

			fixedLine := fixFuncLine(tt.args.line, fset, astFile)
			require.Equal(t, tt.want.fixedLine, fixedLine)
		})
	}
}

func TestGetFuncInfo(t *testing.T) {
	type args struct {
		fn any
	}
	type want struct {
		comment string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "simple",
			args: args{
				fn: tests.Simple,
			},
			want: want{comment: "Simple correct comment\n"},
		},
		{
			name: "not simple",
			args: args{
				fn: tests.NotSimple,
			},
			want: want{comment: "NotSimple correct comment\n"},
		},
		{
			name: "simple2",
			args: args{
				fn: tests.Simple2,
			},
			want: want{comment: "Simple2 correct comment\n"},
		},
		{
			name: "not simple2",
			args: args{
				fn: tests.NotSimple2,
			},
			want: want{comment: "NotSimple2 correct comment\n"},
		},
		{
			name: "simple3",
			args: args{
				fn: tests.Simple3,
			},
			want: want{comment: "Simple3 correct comment\n"},
		},
		{
			name: "simple4 outer",
			args: args{
				fn: tests.Simple4,
			},
			want: want{comment: "Simple4 outer comment\n"},
		},
		{
			name: "simple4 inner",
			args: args{
				fn: tests.Simple4(),
			},
			want: want{comment: "Simple4 outer comment\n"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want.comment, getFuncInfo(tt.args.fn).Comment)
			fmt.Println("docgen.GetFuncInfo(tt.args.fn).Line")
			fmt.Println(docgen.GetFuncInfo(tt.args.fn).Line)
			// require.Equal(t, tt.want.comment, docgen.GetFuncInfo(tt.args.fn).Comment)
		})
	}
}

// wrong comment 2
func TestGetFuncInfoLocal(t *testing.T) {
	// t.Skip()

	// wrong comment 3
	var fn3 func() (int, int)

	//fn comment
	fn := func(a, b int) (int, int) {

		//

		return a, b
	}

	// fn2 comment
	fn2 := func() (int, int) {

		// fn3 comment
		fn3 = func() (int, int) {
			return 1, 2
		}

		return 3, 4
	}

	type args struct {
		fn any
	}
	type want struct {
		comment      string
		unresolvable bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "fn",
			args: args{
				fn: fn,
			},
			want: want{comment: "fn comment\n", unresolvable: false},
		},
		{
			name: "fn2",
			args: args{
				fn: fn2,
			},
			want: want{comment: "fn2 comment\n", unresolvable: false},
		},
		{
			name: "fn3",
			args: args{
				fn: fn3,
			},
			want: want{comment: "fn3 comment\n", unresolvable: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fi := getFuncInfo(tt.args.fn)
			require.Equal(t, tt.want.unresolvable, fi.Unresolvable)

			if !tt.want.unresolvable {
				require.Equal(t, tt.want.comment, fi.Comment)
			}
		})
	}
}

// wrong comment 2
func TestGetFuncInfo3(t *testing.T) {
	// t.Skip()
	tt(t)
}

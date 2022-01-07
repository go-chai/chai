// wrong comment 1

package openapi2

import (
	"testing"

	// "github.com/go-chai/chai/chai"

	"github.com/stretchr/testify/require"
)

//tt - wrong comment 1
func tt(t *testing.T) (int, int) {

	//

	//

	//tt - wrong comment 2
	if false {

		return 1, 2
	}
	//tt - correct comment
	fn := func() (int, int) { return 3, 4 };	fn2 := func() (int, int) { return 3, 4 };	fn3 := func() (int, int) { return 1, 2 }

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
			want: want{comment: "tt - correct comment\n", unresolvable: false},
		},
		{
			name: "fn2",
			args: args{
				fn: fn2,
			},
			want: want{comment: "tt - correct comment\n", unresolvable: false},
		},
		{
			name: "fn3",
			args: args{
				fn: fn3,
			},
			want: want{comment: "tt - correct comment\n", unresolvable: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// fi := docgen.GetFuncInfo(tt.args.fn)
			fi := getFuncInfo(tt.args.fn)
			require.Equal(t, tt.want.unresolvable, fi.Unresolvable)

			if !tt.want.unresolvable {
				require.Equal(t, tt.want.comment, fi.Comment)
			}
		})
	}

	return 3, 3
}

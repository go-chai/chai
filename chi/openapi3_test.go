package chai_test

import (
	"testing"

	chai "github.com/go-chai/chai/chi"
	"github.com/go-chai/chai/examples/shared/controller"
	"github.com/go-chai/chai/internal/tests"
	"github.com/go-chai/chai/log"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestOpenAPI3(t *testing.T) {
	type args struct {
		r chi.Routes
	}
	tcs := []struct {
		name     string
		args     args
		filePath string
		wantErr  bool
	}{
		{
			name: "celler",
			args: args{
				r: controller.NewController().ChiRoutes(),
			},
			filePath: "testdata/celler.json",
		},
	}
	for _, tt := range tcs {
		t.Run(tt.name, func(t *testing.T) {
			got, err := chai.OpenAPI3(tt.args.r)
			require.NoError(t, err)
			log.JSON(got)
			require.JSONEq(t, tests.LoadFile(t, tt.filePath), tests.JS(got))
		})
	}
}

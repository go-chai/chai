package chai_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chai/chai/chai"
	"github.com/go-chai/chai/internal/tests"
	"github.com/stretchr/testify/require"
)

func newRes() *tests.TestResponse {
	return &tests.TestResponse{
		Foo: "f",
		Bar: "b",
		TestInnerResponse: tests.TestInnerResponse{
			FooFoo: 123,
			BarBar: 12,
		},
	}
}

func newReq() io.Reader {
	buf := new(bytes.Buffer)

	json.NewEncoder(buf).Encode(&tests.TestRequest{
		Foo: "312",
		Bar: "31321",
		TestInnerResponse: tests.TestInnerResponse{
			FooFoo: 4432,
			BarBar: 321,
		},
	})

	return buf
}

func TestReqResHandler(t *testing.T) {
	tests := []struct {
		name        string
		makeHandler func(t *testing.T) http.Handler
		response    string
	}{
		{
			name: "req res handler",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewReqResHandler(func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, error) {
					return newRes(), http.StatusOK, nil
				})
			},
			response: `{"foo":"f","bar":"b","test_inner_response":{"foo_foo":123,"bar_bar":12}}`,
		},
		{
			name: "req res handler with error",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewReqResHandler(func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, error) {
					return nil, http.StatusInternalServerError, errors.New("zz")
				})
			},
			response: `{"error":"zz", "status_code":500}`,
		},
		{
			name: "res handler",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewResHandler(func(w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, error) {
					return newRes(), http.StatusOK, nil
				})
			},
			response: `{"foo":"f","bar":"b","test_inner_response":{"foo_foo":123,"bar_bar":12}}`,
		},
		{
			name: "req res handler with error",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewResHandler(func(w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, error) {
					return nil, http.StatusInternalServerError, errors.New("zz")
				})
			},
			response: `{"error":"zz", "status_code":500}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotPanics(t, func() {
				h := tt.makeHandler(t)

				w := httptest.NewRecorder()

				h.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/", newReq()))

				require.JSONEq(t, tt.response, w.Body.String())
			})
		})
	}
}

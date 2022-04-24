package chai_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chai/chai/chai"
	"github.com/go-chai/chai/internal/tests"
	"github.com/go-chai/chai/internal/tests/xrequire"
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

func TestHandlers(t *testing.T) {
	tests := []struct {
		name        string
		makeHandler func(t *testing.T) http.Handler
		response    string
	}{
		{
			name: "req res handler",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewHandler("", "", func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, error) {
					return newRes(), http.StatusOK, nil
				})
			},
			response: `{"foo":"f","bar":"b","test_inner_response":{"foo_foo":123,"bar_bar":12}}`,
		},
		{
			name: "req res handler with error",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewHandler("", "", func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, error) {
					return nil, http.StatusInternalServerError, errors.New("zz")
				})
			},
			response: `{"error":"zz", "status_code":500}`,
		},
		{
			name: "req res handler with custom struct error type with a pointer receiver with no error",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewHandler("", "", func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, *tests.TestErrorPtr) {
					return newRes(), http.StatusOK, nil
				})
			},
			response: `{"foo":"f","bar":"b","test_inner_response":{"foo_foo":123,"bar_bar":12}}`,
		},
		{
			name: "req res handler with custom struct error type with a pointer receiver with error",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewHandler("", "", func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, *tests.TestErrorPtr) {
					return nil, http.StatusInternalServerError, &tests.TestErrorPtr{Message: "zz"}
				})
			},
			response: `{"error":"zz", "message":"zz", "status_code":500}`,
		},
		{
			name: "req res handler with custom struct error type with no error",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewHandler("", "", func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, tests.TestError) {
					return newRes(), http.StatusOK, tests.TestError{}
				})
			},
			response: `{"foo":"f","bar":"b","test_inner_response":{"foo_foo":123,"bar_bar":12}}`,
		},
		{
			name: "req res handler with custom struct error type with error",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewHandler("", "", func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, tests.TestError) {
					return nil, http.StatusInternalServerError, tests.TestError{Message: "zz"}
				})
			},
			response: `{"error":"zz", "message":"zz", "status_code":500}`,
		},
		{
			name: "req res handler with custom map error type with no error",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewHandler("", "", func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, tests.TestErrorMap) {
					return newRes(), http.StatusOK, nil
				})
			},
			response: `{"foo":"f","bar":"b","test_inner_response":{"foo_foo":123,"bar_bar":12}}`,
		},
		{
			name: "req res handler with custom map error type with error",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewHandler("", "", func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, tests.TestErrorMap) {
					return nil, http.StatusInternalServerError, tests.TestErrorMap{"message": "zz"}
				})
			},
			response: `{"error":"test error map", "message":"zz", "status_code":500}`,
		},
		{
			name: "req res handler with custom map error type with a pointer receiver with no error",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewHandler("", "", func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, *tests.TestErrorMapPtr) {
					return newRes(), http.StatusOK, nil
				})
			},
			response: `{"foo":"f","bar":"b","test_inner_response":{"foo_foo":123,"bar_bar":12}}`,
		},
		{
			name: "req res handler with custom map error type with a pointer receiver with error",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewHandler("", "", func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, *tests.TestErrorMapPtr) {
					return nil, http.StatusInternalServerError, &tests.TestErrorMapPtr{"message": "zz"}
				})
			},
			response: `{"error":"test error map ptr", "message":"zz", "status_code":500}`,
		},
		// TODO add test cases for parsing the requests (body, query, path)
		// TODO add test cases for validating the requests
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotPanics(t, func() {
				h := tt.makeHandler(t)

				w := httptest.NewRecorder()

				r := httptest.NewRequest(http.MethodPost, "/", newReq())
				r.Header.Set("Content-Type", "application/json")

				h.ServeHTTP(w, r)

				fmt.Printf("%q\n", w.Body.String())

				xrequire.JSONEq(t, tt.response, w.Body.String())
			})
		})
	}
}

package chai_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chai/chai/chai"
	chaichi "github.com/go-chai/chai/chi"
	"github.com/go-chai/chai/internal/tests"
	"github.com/go-chai/chai/internal/tests/xrequire"
	"github.com/go-chai/chai/log"
	"github.com/go-chi/chi/v5"
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

	json.NewEncoder(buf).Encode(map[string]any{
		"foob": "312",
		"barb": "31321",
		"test_inner_responseb": map[string]any{
			"foo_foo": 4432,
			"bar_bar": 321,
		},
	})

	return buf
}

func TestDecode(t *testing.T) {
	s := &tests.TestRequest{
		P1:  123,
		MM:  "asd",
		GG:  "gg",
		Foo: "mm",
		Bar: "bb",
		TestInnerResponse: tests.TestInnerResponse{
			FooFoo: 13,
			BarBar: 12,
		},
	}

	// json.NewDecoder(bytes.NewReader([]byte(`{"p1":124}`))).Decode(&s)
	json.NewDecoder(bytes.NewReader([]byte(`{}`))).Decode(&s)

	log.Dump(s)
}

func TestHandlers(t *testing.T) {
	tests := []struct {
		name        string
		makeHandler func(t *testing.T) http.Handler
		response    string
	}{
		{
			name: "handler",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewHandler("", "", func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, error) {
					return newRes(), http.StatusOK, nil
				})
			},
			response: `{"foo":"f","bar":"b","test_inner_response":{"foo_foo":123,"bar_bar":12}}`,
		},
		{
			name: "handler with error",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewHandler("", "", func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, error) {
					return nil, http.StatusInternalServerError, errors.New("zz")
				})
			},
			response: `{"error":"zz", "status_code":500}`,
		},
		{
			name: "handler with custom struct error type with a pointer receiver with no error",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewHandler("", "", func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, *tests.TestErrorPtr) {
					return newRes(), http.StatusOK, nil
				})
			},
			response: `{"foo":"f","bar":"b","test_inner_response":{"foo_foo":123,"bar_bar":12}}`,
		},
		{
			name: "handler with custom struct error type with a pointer receiver with error",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewHandler("", "", func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, *tests.TestErrorPtr) {
					return nil, http.StatusInternalServerError, &tests.TestErrorPtr{Message: "zz"}
				})
			},
			response: `{"error":"zz", "message":"zz", "status_code":500}`,
		},
		{
			name: "handler with custom struct error type with no error",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewHandler("", "", func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, tests.TestError) {
					return newRes(), http.StatusOK, tests.TestError{}
				})
			},
			response: `{"foo":"f","bar":"b","test_inner_response":{"foo_foo":123,"bar_bar":12}}`,
		},
		{
			name: "handler with custom struct error type with error",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewHandler("", "", func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, tests.TestError) {
					return nil, http.StatusInternalServerError, tests.TestError{Message: "zz"}
				})
			},
			response: `{"error":"zz", "message":"zz", "status_code":500}`,
		},
		{
			name: "handler with custom map error type with no error",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewHandler("", "", func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, tests.TestErrorMap) {
					return newRes(), http.StatusOK, nil
				})
			},
			response: `{"foo":"f","bar":"b","test_inner_response":{"foo_foo":123,"bar_bar":12}}`,
		},
		{
			name: "handler with custom map error type with error",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewHandler("", "", func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, tests.TestErrorMap) {
					return nil, http.StatusInternalServerError, tests.TestErrorMap{"message": "zz"}
				})
			},
			response: `{"error":"test error map", "message":"zz", "status_code":500}`,
		},
		{
			name: "handler with custom map error type with a pointer receiver with no error",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewHandler("", "", func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, *tests.TestErrorMapPtr) {
					return newRes(), http.StatusOK, nil
				})
			},
			response: `{"foo":"f","bar":"b","test_inner_response":{"foo_foo":123,"bar_bar":12}}`,
		},
		{
			name: "handler with custom map error type with a pointer receiver with error",
			makeHandler: func(t *testing.T) http.Handler {
				r := chi.NewRouter()

				return chaichi.Get(r, "/{p1}/test", func(req *tests.TestRequest, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, *tests.TestErrorMapPtr) {
					log.Dump(req)

					require.Equal(t, "312", req.Foo)
					require.Equal(t, "query_value", req.MM)
					require.Equal(t, 14, req.P1)
					return nil, http.StatusInternalServerError, &tests.TestErrorMapPtr{"message": "zz"}
				})
			},
			response: `{"error":"test error map ptr", "message":"zz", "status_code":500}`,
		},
		{
			name: "handler with any request with custom map error type with a pointer receiver with error",
			makeHandler: func(t *testing.T) http.Handler {
				return chai.NewHandler("", "", func(req any, w http.ResponseWriter, r *http.Request) (*tests.TestResponse, int, *tests.TestErrorMapPtr) {
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

				r := httptest.NewRequest(http.MethodPost, "/?mm=query_value", newReq())
				r.Header.Set("Content-Type", "application/json")

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("p1", "14")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				log.Dump(r.URL)

				h.ServeHTTP(w, r)

				log.Dump(w.Body.String())

				xrequire.JSONEq(t, tt.response, w.Body.String())
			})
		})
	}
}

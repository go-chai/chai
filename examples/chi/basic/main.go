package main

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	chai "github.com/go-chai/chai/chi"
	_ "github.com/go-chai/chai/examples/docs/basic" // This is required to be able to serve the stored swagger spec in prod
	"github.com/go-chai/chai/examples/shared/httputil"
	"github.com/go-chai/chai/examples/shared/model"
	chaiopenapi "github.com/go-chai/chai/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/examples", func(r chi.Router) {
			chai.Post(r, "/post", PostHandler).
				ID("123123123").
				Deprecated().
				Summary(`some text used as a summary
				for this example post handler`)

			chai.Post(r, "/{pathParam}/post2", PostHandler).
				// WithValidator(func(a *model.Address) error {
				// 	err := validation.ValidateStruct(&a,
				// 		validation.Field(a.Zip, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))),
				// 	)
				// 	if err != nil {
				// 		return &httputil.Error{Message: err.Error(), StatusCode: http.StatusBadRequest}
				// 	}
				// 	return nil
				// }).
				ID("123123123")
			// WithValidator((*model.Address).ValidateStep1).
			// WithValidator((*model.Address).ValidateStep2)

			chai.Get(r, "/uuid2", UUIDHandler)
			chai.Get(r, "/calc2", CalcHandler)
			chai.Get(r, "/ping", PingHandler)
			chai.Get(r, "/groups/{group_id}/accounts/{account_id}", PathParamsHandler)
			chai.Get(r, "/header", HeaderHandler)
			chai.Get(r, "/securities", SecuritiesHandler)
			chai.Get(r, "/attribute", AttributeHandler)
		})
	})

	// This must be used only during development to generate the swagger spec
	docs, err := chai.OpenAPI3(r)
	if err != nil {
		panic(fmt.Sprintf("failed to generate the swagger spec: %+v", err))
	}

	// This should be used in prod to serve the swagger spec
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	addCustomDocs(docs)

	chaiopenapi.LogYAML(docs)

	// This must be used only during development to store the swagger spec
	// err = openapi2.WriteDocs(docs, &openapi2.GenConfig{
	// OutputDir: "examples/docs/basic",
	// })
	if err != nil {
		panic(fmt.Sprintf("failed to write the swagger spec: %+v", err))
	}

	fmt.Println("The swagger spec is available at http://localhost:8080/swagger/")

	http.ListenAndServe(":8080", r)
}

type Int struct {
	*big.Int
}

func (i *Int) Validate() error {
	return nil
}

func PostHandler(account ***model.Address, w http.ResponseWriter, r *http.Request) (*model.Address, int, *httputil.Error) {
	return **account, http.StatusOK, nil
}

func CalcHandler2(w http.ResponseWriter, r *http.Request) (string, int, error) {
	val1, err := strconv.Atoi(r.URL.Query().Get("val1"))
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	val2, err := strconv.Atoi(r.URL.Query().Get("val2"))
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	return fmt.Sprintf("%d", val1*val2), http.StatusOK, nil
}

func CalcHandler(w http.ResponseWriter, r *http.Request) (*big.Int, int, error) {
	val1, err := strconv.Atoi(r.URL.Query().Get("val1"))
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	val2, err := strconv.Atoi(r.URL.Query().Get("val2"))
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return big.NewInt(int64(val1 + val2)), http.StatusOK, nil
}

func UUIDHandler(w http.ResponseWriter, r *http.Request) (uuid.UUID, int, error) {
	return uuid.Must(uuid.NewV4()), http.StatusOK, nil
}

func PingHandler(w http.ResponseWriter, r *http.Request) (string, int, error) {
	return "pong", http.StatusOK, nil
}

func PathParamsHandler(w http.ResponseWriter, r *http.Request) (string, int, error) {
	groupID, err := strconv.Atoi(chi.URLParam(r, "group_id"))
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	accountID, err := strconv.Atoi(chi.URLParam(r, "account_id"))
	if err != nil {
		return "", http.StatusBadRequest, err
	}

	return fmt.Sprintf("group_id=%d account_id=%d", groupID, accountID), http.StatusOK, nil
}

func HeaderHandler(w http.ResponseWriter, r *http.Request) (string, int, error) {
	return r.Header.Get("Authorization"), http.StatusOK, nil
}

func SecuritiesHandler(w http.ResponseWriter, r *http.Request) (string, int, error) {
	return "ok", http.StatusOK, nil
}

func AttributeHandler(w http.ResponseWriter, r *http.Request) (string, int, error) {
	return fmt.Sprintf("enumstring=%s enumint=%s enumnumber=%s string=%s int=%s default=%s",
		r.URL.Query().Get("enumstring"),
		r.URL.Query().Get("enumint"),
		r.URL.Query().Get("enumnumber"),
		r.URL.Query().Get("string"),
		r.URL.Query().Get("int"),
		r.URL.Query().Get("default"),
	), http.StatusOK, nil
}

func addCustomDocs(docs *openapi3.T) {
	docs.OpenAPI = "3.0"
	docs.Servers = openapi3.Servers{{URL: "localhost:8080"}}

	docs.Info = &openapi3.Info{
		Description:    "This is a sample celler server.",
		Title:          "Swagger Example API",
		TermsOfService: "http://swagger.io/terms/",
		Contact: &openapi3.Contact{
			Name:  "API Support",
			URL:   "http://www.swagger.io/support",
			Email: "support@swagger.io",
		},
		License: &openapi3.License{
			Name: "Apache 2.0",
			URL:  "http://www.apache.org/licenses/LICENSE-2.0.html",
		},
		Version: "1.0",
	}
}

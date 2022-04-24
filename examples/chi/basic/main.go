package main

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	chai "github.com/go-chai/chai/chi"
	"github.com/go-chai/chai/examples/shared/httputil"
	"github.com/go-chai/chai/examples/shared/model"
	"github.com/go-chai/chai/log"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

func main() {
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/examples", func(r chi.Router) {
			chai.Post(r, "/post", PostHandler).
				ID("example").
				Tags("Examples").
				Deprecated().
				Summary(`some text used as a summary
				for this example post handler`)

			chai.Post(r, "/{pathParam}/post2", PostHandler).
				Tags("Examples").
				ID("example2")
			chai.Get(r, "/uuid2", UUIDHandler)
			chai.Get(r, "/calc2", CalcHandler)
			chai.Get(r, "/ping", PingHandler)
			chai.Get(r, "/groups/{group_id}/accounts/{account_id}", PathParamsHandler)
			chai.Get(r, "/header", HeaderHandler)
			chai.Get(r, "/securities", SecuritiesHandler)
			chai.Get(r, "/attribute", AttributeHandler)
		})
	})

	docs, err := chai.OpenAPI3(r)
	if err != nil {
		panic(fmt.Sprintf("failed to generate the swagger spec: %+v", err))
	}
	addCustomDocs(docs)
	log.YAML(docs)

	// Serve the swagger spec
	r.Get("/swagger/*", chai.SwaggerHandler(docs))

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

func CalcHandler2(req any, w http.ResponseWriter, r *http.Request) (string, int, error) {
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

func CalcHandler(req any, w http.ResponseWriter, r *http.Request) (*big.Int, int, error) {
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

func UUIDHandler(req any, w http.ResponseWriter, r *http.Request) (uuid.UUID, int, error) {
	return uuid.Must(uuid.NewV4()), http.StatusOK, nil
}

func PingHandler(req any, w http.ResponseWriter, r *http.Request) (string, int, error) {
	return "pong", http.StatusOK, nil
}

func PathParamsHandler(req any, w http.ResponseWriter, r *http.Request) (string, int, error) {
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

func HeaderHandler(req any, w http.ResponseWriter, r *http.Request) (string, int, error) {
	return r.Header.Get("Authorization"), http.StatusOK, nil
}

func SecuritiesHandler(req any, w http.ResponseWriter, r *http.Request) (string, int, error) {
	return "ok", http.StatusOK, nil
}

func AttributeHandler(req any, w http.ResponseWriter, r *http.Request) (string, int, error) {
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

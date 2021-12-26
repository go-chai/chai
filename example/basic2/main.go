package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ghodss/yaml"
	"github.com/go-chai/chai"
	"github.com/go-chai/chai/example/celler/model"
	"github.com/go-chai/chai/openapi2"
	"github.com/go-chai/chai/specc"
	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/spec"
)

func main() {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/examples", func(r chi.Router) {
			chai.Post(r, "/post", PostExample)
			chai.Get(r, "/calc", CalcExample)
			chai.Get(r, "/ping", Ping)
			chai.Get(r, "/groups/{group_id}/accounts/{account_id}", PathParamsExample)
			chai.Get(r, "/header", HeaderExample)
			chai.Get(r, "/securities", SecuritiesExample)
			chai.Get(r, "/attribute", AttributeExample)
		})
	})

	docs, err := openapi2.Docs(r)
	if err != nil {
		panic(err)
	}

	addCustomDocs(docs)

	LogYAML(docs)

	http.ListenAndServe(":8080", r)
}

// PostExample godoc
// @Summary      post request example
// @Description  post request example
// @Accept       json
// @Produce      plain
// @Param        message  body      model.Account  true  "Account Info"
// @Success      200      {string}  string         "success"
// @Param        enumnumber  query     number  false  "int enums"       Enums(1.1, 1.2, 1.3)
func PostExample(account *model.Account, w http.ResponseWriter, r *http.Request) (*model.Account, int, *chai.JSONError) {
	return account, http.StatusOK, nil
}

// CalcExampleee godoc
// @Summary      calc example
// @Description  plus
// @Tags         example
// @Accept       json
// @Produce      json
// @Param        val1  query      int     true  "used for calc"
// @Param        val2  query      int     true  "used for calc"
// @Success      200
// @Failure      400,404,500
func CalcExample(w http.ResponseWriter, r *http.Request) (string, int, error) {
	val1, err := strconv.Atoi(r.URL.Query().Get("val1"))
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	val2, err := strconv.Atoi(r.URL.Query().Get("val2"))
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	ans := val1 + val2
	return fmt.Sprintf("%d", ans), http.StatusOK, nil
}

// PingExample godoc
// @Summary      ping example
// @Description  do ping
// @Tags         example
// @Produce      json
// @Success      200
// @Failure      400,404,500
func Ping(w http.ResponseWriter, r *http.Request) (string, int, error) {
	return "pong", http.StatusOK, nil
}

// PathParamsExample godoc
// @Summary      path params example
// @Description  path params
// @Tags         example
// @Accept       json
// @Produce      json
// @Param        group_id    path      int     true  "Group ID"
// @Param        account_id  path      int     true  "Account ID"
// @Success      200
// @Failure      400,404,500
func PathParamsExample(w http.ResponseWriter, r *http.Request) (string, int, error) {
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

// HeaderExample godoc
// @Summary      custome header example
// @Description  custome header
// @Tags         example
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Authentication header"
// @Success      200
// @Failure      400,404,500
func HeaderExample(w http.ResponseWriter, r *http.Request) (string, int, error) {
	return r.Header.Get("Authorization"), http.StatusOK, nil
}

// SecuritiesExample godoc
// @Summary      custome header example
// @Description  custome header
// @Tags         example
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Authentication header"
// @Success      200
// @Failure      400,404,500
// @Security     ApiKeyAuth
func SecuritiesExample(w http.ResponseWriter, r *http.Request) (string, int, error) {
	return "ok", http.StatusOK, nil
}

// AttributeExample godoc
// @Summary      attribute example
// @Description  attribute
// @Tags         example
// @Accept       json
// @Produce      json
// @Param        enumstring  query     string  false  "string enums"    Enums(A, B, C)
// @Param        enumint     query     int     false  "int enums"       Enums(1, 2, 3)
// @Param        enumnumber  query     number  false  "int enums"       Enums(1.1, 1.2, 1.3)
// @Param        string      query     string  false  "string valid"    minlength(5)  maxlength(10)
// @Param        int         query     int     false  "int valid"       minimum(1)    maximum(10)
// @Param        default     query     string  false  "string default"  default(A)
// @Success      200 "answer"
// @Failure      400,404,500 "ok"
func AttributeExample(w http.ResponseWriter, r *http.Request) (string, int, error) {
	return fmt.Sprintf("enumstring=%s enumint=%s enumnumber=%s string=%s int=%s default=%s",
		r.URL.Query().Get("enumstring"),
		r.URL.Query().Get("enumint"),
		r.URL.Query().Get("enumnumber"),
		r.URL.Query().Get("string"),
		r.URL.Query().Get("int"),
		r.URL.Query().Get("default"),
	), http.StatusOK, nil
}

func addCustomDocs(docs *specc.Swagger) {
	docs.Swagger.Swagger = "2.0"
	docs.Host = "localhost:8080"
	docs.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Description:    "This is a sample celler server.",
			Title:          "Swagger Example API",
			TermsOfService: "http://swagger.io/terms/",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "API Support",
					URL:   "http://www.swagger.io/support",
					Email: "support@swagger.io",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "Apache 2.0",
					URL:  "http://www.apache.org/licenses/LICENSE-2.0.html",
				},
			},
			Version: "1.0",
		},
	}
	docs.SecurityDefinitions = map[string]*spec.SecurityScheme{
		"ApiKeyAuth": {
			SecuritySchemeProps: spec.SecuritySchemeProps{
				Type: "apiKey",
				In:   "header",
				Name: "Authorization",
			},
		},
	}
}

func LogYAML(v interface{}) {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

	return
}

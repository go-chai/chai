# Chai

## Description

Chai is an extension for the [chi](https://github.com/go-chi/chi) router that can automatically generate a Swagger 2.0 spec from a router. 

Go 1.18 introduces generics which are used by chai for inferring the types of requests and responses without losing type safety.

[Swaggo](https://github.com/swaggo/swag) annotations are used as a fallback for the parts of the swagger spec that cannot be inferred automatically.

## Project status
Chai is still a work in progress

## Gotchas

- Compiler optimizations

    The way the swaggo annotations are obtained from the source code currently relies on the caller frames and go's compiler optimizations sometimes change those which results in the annotations for "small" handlers to get ignored. As a workaround, the optimizations can be disabled for the binary that generates the docs by passing `-gcflags '-N'` to the compiler, e.g. `go run -gcflags='-N' ./example/basic2/main.go`

## Usage

This [swagger.yaml](https://editor.swagger.io/?url=https://raw.githubusercontent.com/go-chai/chai/main/example/basic2/swagger.yaml) is generated by the program below. Notice that the spec for the `PostExample` handler was generated without any annotations. The request/response types and the route were detected automatically from the router itself.

![image](https://user-images.githubusercontent.com/1100051/147469137-18f47ce2-d43f-4b36-94fc-8c0e0793a9e0.png)

```go
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
			chai.PostG(r, "/post", PostExample)
			chai.GetG(r, "/calc", CalcExample)
			chai.GetG(r, "/ping", Ping)
			chai.GetG(r, "/groups/{group_id}/accounts/{account_id}", PathParamsExample)
			chai.GetG(r, "/header", HeaderExample)
			chai.GetG(r, "/securities", SecuritiesExample)
			chai.GetG(r, "/attribute", AttributeExample)
		})
	})

	docs, err := openapi2.Docs(r)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	addCustomDocs(docs)

	LogYAML(docs)

	http.ListenAndServe(":8080", r)
}

func PostExample(account *model.Account, w http.ResponseWriter, r *http.Request) (*model.Account, int, *chai.JSONError) {
	return account, http.StatusOK, nil
}

// @Param        val1  query      int     true  "used for calc"
// @Param        val2  query      int     true  "used for calc"
// @Success      203
// @Failure      400,404
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
func Ping(w http.ResponseWriter, r *http.Request) (string, int, error) {
	return "pong", http.StatusOK, nil
}

// PathParamsExample godoc
// @Summary      path params example
// @Description  path params
// @Tags         example
// @Param        group_id    path      int     true  "Group ID"
// @Param        account_id  path      int     true  "Account ID"
// @Failure      400,404
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
// @Param        Authorization  header    string  true  "Authentication header"
// @Failure      400,404
func HeaderExample(w http.ResponseWriter, r *http.Request) (string, int, error) {
	return r.Header.Get("Authorization"), http.StatusOK, nil
}

// SecuritiesExample godoc
// @Summary      custome header example
// @Description  custome header
// @Tags         example
// @Param        Authorization  header    string  true  "Authentication header"
// @Failure      400,404
// @Security     ApiKeyAuth
func SecuritiesExample(w http.ResponseWriter, r *http.Request) (string, int, error) {
	return "ok", http.StatusOK, nil
}

// AttributeExample godoc
// @Summary      attribute example
// @Description  attribute
// @Tags         example
// @Param        enumstring  query     string  false  "string enums"    Enums(A, B, C)
// @Param        enumint     query     int     false  "int enums"       Enums(1, 2, 3)
// @Param        enumnumber  query     number  false  "int enums"       Enums(1.1, 1.2, 1.3)
// @Param        string      query     string  false  "string valid"    minlength(5)  maxlength(10)
// @Param        int         query     int     false  "int valid"       minimum(1)    maximum(10)
// @Param        default     query     string  false  "string default"  default(A)
// @Success      200 "answer"
// @Failure      400,404 "ok"
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
```

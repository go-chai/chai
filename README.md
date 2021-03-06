# chai

## Description

`chai` is an extension for a few popular http routers that adds support for type safe http handlers via Go 1.18's generics. This allows it to also generate a swagger spec by automatically detecting the request/response types, http methods, route paths and path params.

`chai` uses [swaggo/swag](https://github.com/swaggo/swag) annotations for the parts of the swagger spec that cannot be automatically inferred.


## Supported http routers

- [chi](https://github.com/go-chi/chi)
- [gorilla/mux](https://github.com/gorilla/mux)

## Project status
`chai` is still a work in progress

## Gotchas

- YAML marshalling

	Currently only https://github.com/ghodss/yaml is supported as a yaml marshaller for the generated swagger spec, which is also provided via `openapi2.MarshalYAML()` as an alias

## Examples

- chi - [./examples/chi](./examples/chi)
- gorilla/mux - [./examples/gorilla](./examples/gorilla)
- standalone repo - https://github.com/go-chai/examples

## Usage

This [swagger.yaml](https://editor.swagger.io/?url=https://raw.githubusercontent.com/go-chai/chai/main/examples/docs/basic/swagger.yaml) is generated by the program below. Notice that the spec for the `PostHandler` handler was generated without any annotations. The request/response types and the route were detected automatically from the router itself.

![image](https://user-images.githubusercontent.com/1100051/147469383-b257f396-be7c-45d9-bf55-a1f7454bf5bf.png)

```go
package main

import (
	"fmt"
	"net/http"
	"strconv"

	chai "github.com/go-chai/chai/chi"
	_ "github.com/go-chai/chai/examples/docs/basic" // This is required to be able to serve the stored swagger spec in prod
	"github.com/go-chai/chai/examples/shared/model"
	"github.com/go-chai/chai/openapi2"
	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/spec"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/examples", func(r chi.Router) {
			chai.Post(r, "/post", PostHandler)
			chai.Get(r, "/calc", CalcHandler)
			chai.Get(r, "/ping", PingHandler)
			chai.Get(r, "/groups/{group_id}/accounts/{account_id}", PathParamsHandler)
			chai.Get(r, "/header", HeaderHandler)
			chai.Get(r, "/securities", SecuritiesHandler)
			chai.Get(r, "/attribute", AttributeHandler)
		})
	})

	// This must be used only during development to generate the swagger spec
	docs, err := chai.OpenAPI2(r)
	if err != nil {
		panic(fmt.Sprintf("failed to generate the swagger spec: %+v", err))
	}

	// This should be used in prod to serve the swagger spec
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	addCustomDocs(docs)

	openapi2.LogYAML(docs)

	// This must be used only during development to store the swagger spec
	err = openapi2.WriteDocs(docs, &openapi2.GenConfig{
		OutputDir: "examples/docs/basic",
	})
	if err != nil {
		panic(fmt.Sprintf("failed to write the swagger spec: %+v", err))
	}

	fmt.Println("The swagger spec is available at http://localhost:8080/swagger/")

	http.ListenAndServe(":8080", r)
}

type Error struct {
	Message          string `json:"error"`
	ErrorDebug       string `json:"error_debug,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
	StatusCode       int    `json:"status_code,omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}

func PostHandler(account *model.Account, w http.ResponseWriter, r *http.Request) (*model.Account, int, *Error) {
	return account, http.StatusOK, nil
}

// @Param        val1  query      int     true  "used for calc"
// @Param        val2  query      int     true  "used for calc"
// @Success      203
// @Failure      400,404
func CalcHandler(w http.ResponseWriter, r *http.Request) (string, int, error) {
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

// PingExample godoc
// @Summary      ping example
// @Description  do ping
// @Tags         example
func PingHandler(w http.ResponseWriter, r *http.Request) (string, int, error) {
	return "pong", http.StatusOK, nil
}

// PathParamsHandler godoc
// @Summary      path params example
// @Description  path params
// @Tags         example
// @Param        group_id    path      int     true  "Group ID"
// @Param        account_id  path      int     true  "Account ID"
// @Failure      400,404
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

// HeaderHandler godoc
// @Summary      custome header example
// @Description  custome header
// @Tags         example
// @Param        Authorization  header    string  true  "Authentication header"
// @Failure      400,404
func HeaderHandler(w http.ResponseWriter, r *http.Request) (string, int, error) {
	return r.Header.Get("Authorization"), http.StatusOK, nil
}

// SecuritiesHandler godoc
// @Summary      custome header example
// @Description  custome header
// @Tags         example
// @Param        Authorization  header    string  true  "Authentication header"
// @Failure      400,404
// @Security     ApiKeyAuth
func SecuritiesHandler(w http.ResponseWriter, r *http.Request) (string, int, error) {
	return "ok", http.StatusOK, nil
}

// AttributeHandler godoc
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

func addCustomDocs(docs *spec.Swagger) {
	docs.Swagger = "2.0"
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

```

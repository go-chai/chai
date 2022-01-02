package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ghodss/yaml"
	"github.com/gorilla/mux"

	_ "github.com/go-chai/chai/examples/docs/celler" // This is required to be able to serve the stored swagger spec in prod
	"github.com/go-chai/chai/examples/shared/controller"
	"github.com/go-chai/chai/examples/shared/httputil"
	chai "github.com/go-chai/chai/gorilla"
	"github.com/go-chai/chai/openapi2"
	"github.com/go-openapi/spec"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	r := mux.NewRouter()

	c := controller.NewController()

	chai.Get(r, "/api/v1/accounts/{id}", c.ShowAccount)
	chai.Get(r, "/api/v1/accounts/", c.ListAccounts)
	chai.Post(r, "/api/v1/accounts/", c.AddAccount)
	r.HandleFunc("/api/v1/accounts/{id}", c.DeleteAccount).Methods(http.MethodDelete)
	r.HandleFunc("/api/v1/accounts/{id}", c.UpdateAccount).Methods(http.MethodPatch)
	r.HandleFunc("/api/v1/accounts/{id}/images", c.UploadAccountImage).Methods(http.MethodPost)
	chai.Get(r, "/api/v1/bottles/{id}", c.ShowBottle)
	chai.Get(r, "/api/v1/bottles/", c.ListBottles)
	chai.Get(r, "/api/v1/bottles/", c.ListBottles)

	authMux := r.Path("/").Subrouter()
	authMux.Use(auth)

	chai.Post(authMux, "/api/v1/admin/auth", c.Auth)

	chai.Get(r, "/api/v1/examples/ping", c.PingExample)
	chai.Get(r, "/api/v1/examples/calc", c.CalcExample)
	// chai.Get(r, "/api/v1/examples/group{s/{gro}up_id}/accounts/{account_id}", c.CalcExample)
	chai.Get(r, "/api/v1/examples/groups/{group_id}/accounts/{account_id}", c.PathParamsExample)
	chai.Get(r, "/api/v1/examples/header", c.HeaderExample)
	chai.Get(r, "/api/v1/examples/securities", c.SecuritiesExample)
	chai.Get(r, "/api/v1/examples/attribute", c.AttributeExample)
	chai.Post(r, "/api/v1/examples/attribute", c.PostExample)

	// This must be used only during development to generate the swagger spec
	docs, err := chai.OpenAPI2(r)
	if err != nil {
		panic(fmt.Sprintf("chai.OpenAPI2() failed: %+v", err))
	}

	// This should be used in prod to serve the swagger spec
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	addCustomDocs(docs)

	LogYAML(docs)

	// This must be used only during development to store the swagger spec
	err = openapi2.WriteDocs(docs, &openapi2.GenConfig{
		OutputDir: "examples/docs/celler",
	})
	if err != nil {
		panic(fmt.Sprintf("gen.New().Generate() failed: %+v", err))
	}

	fmt.Println("Find the swagger spec at http://localhost:8080/swagger/")

	http.ListenAndServe(":8080", r)
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
		"BasicAuth": {
			SecuritySchemeProps: spec.SecuritySchemeProps{
				Type: "basic",
			},
		},
		"ApiKeyAuth": {
			SecuritySchemeProps: spec.SecuritySchemeProps{
				Type: "apiKey",
				In:   "header",
				Name: "Authorization",
			},
		},
		"OAuth2Implicit": {
			SecuritySchemeProps: spec.SecuritySchemeProps{
				Description:      "Use with the OAuth2 Implicit Grant to retrieve a token",
				Type:             "oauth2",
				Flow:             "implicit",
				AuthorizationURL: "https://example.com/oauth/authorize",
				TokenURL:         "",
				Scopes: map[string]string{
					"admin": "Grants read and write access to administrative information",
					"write": "Grants write access",
				},
			},
		},
		"OAuth2Application": {
			SecuritySchemeProps: spec.SecuritySchemeProps{
				Description: "Use with the OAuth2 Implicit Grant to retrieve a token",
				Type:        "oauth2",
				Flow:        "application",
				TokenURL:    "https://example.com/oauth/token",
				Scopes: map[string]string{
					"admin": "Grants read and write access to administrative information",
					"write": "Grants write access",
				},
			},
		},

		"OAuth2Password": {
			SecuritySchemeProps: spec.SecuritySchemeProps{
				Type:     "oauth2",
				Flow:     "password",
				TokenURL: "https://example.com/oauth/token",
				Scopes: map[string]string{
					"admin": "Grants read and write access to administrative information",
					"write": "Grants write access",
					"read":  "Grants read access",
				},
			},
		},
		"OAuth2AccessToken": {
			SecuritySchemeProps: spec.SecuritySchemeProps{
				Type:             "oauth2",
				Flow:             "accessCode",
				AuthorizationURL: "https://example.com/oauth/authorize",
				TokenURL:         "https://example.com/oauth/token",
				Scopes: map[string]string{
					"admin": "Grants read and write access to administrative information",
				},
			},
		},
	}
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.Header.Get("Authorization")) == 0 {
			httputil.NewError(w, http.StatusUnauthorized, errors.New("Authorization is required Header"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func LogYAML(v any) {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

	return
}

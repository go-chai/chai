package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ghodss/yaml"
	"github.com/go-chai/chai"
	"github.com/go-chai/chai/examples/celler/controller"
	_ "github.com/go-chai/chai/examples/celler/docs"
	"github.com/go-chai/chai/examples/celler/httputil"
	"github.com/go-chai/chai/openapi2"
	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/spec"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	r := chi.NewRouter()

	c := controller.NewController()

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/accounts", func(r chi.Router) {
			chai.GetG(r, "/{id}", c.ShowAccount)
			chai.GetG(r, "/", c.ListAccounts)
			chai.PostG(r, "/", c.AddAccount)
			chai.Delete(r, "/{id}", c.DeleteAccount)
			chai.Patch(r, "/{id}", c.UpdateAccount)
			chai.Post(r, "/{id}/images", c.UploadAccountImage)
		})

		r.Route("/bottles", func(r chi.Router) {
			chai.GetG(r, "/{id}", c.ShowBottle)
			chai.GetG(r, "/", c.ListBottles)
		})

		r.Route("/admin", func(r chi.Router) {
			r.Use(auth)

			chai.PostG(r, "/auth", c.Auth)
		})

		r.Route("/examples", func(r chi.Router) {
			chai.GetG(r, "/ping", c.PingExample)
			chai.GetG(r, "/calc", c.CalcExample)
			// chai.GetG(r, "/group{s/{gro}up_id}/accounts/{account_id}", c.PathParamsExample)
			chai.GetG(r, "/groups/{group_id}/accounts/{account_id}", c.PathParamsExample)
			chai.GetG(r, "/header", c.HeaderExample)
			chai.GetG(r, "/securities", c.SecuritiesExample)
			chai.GetG(r, "/attribute", c.AttributeExample)
			chai.PostG(r, "/attribute", c.PostExample)
		})
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	docs, err := openapi2.Docs(r)
	if err != nil {
		panic(fmt.Sprintf("openapi2.Docs() failed: %+v", err))
	}

	addCustomDocs(docs)

	// LogJSON(docs, "swagger")
	LogYAML(docs, "")

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

func LogYAML(v interface{}, label string) {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}

	if label != "" {
		fmt.Printf("%s:\n", label)
	}
	fmt.Println(string(bytes))

	return
}

func LogJSON(v interface{}, label string) {
	bytes, err := json.MarshalIndent(v, "", "  ")

	if err != nil {
		panic(err)
	}

	if label != "" {
		fmt.Printf("%s:\n", label)
	}
	fmt.Println(string(bytes))

	return
}

package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ghodss/yaml"

	"github.com/go-chai/chai"
	"github.com/go-chai/chai/examples/celler/controller"
	_ "github.com/go-chai/chai/examples/celler/docs" // This is required to be able to serve the stored swagger spec in prod
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
			chai.Get(r, "/{id}", c.ShowAccount)
			chai.Get(r, "/", c.ListAccounts)
			chai.Post(r, "/", c.AddAccount)
			r.Delete("/{id}", c.DeleteAccount)
			r.Patch("/{id}", c.UpdateAccount)
			r.Post("/{id}/images", c.UploadAccountImage)
		})

		r.Route("/bottles", func(r chi.Router) {
			chai.Get(r, "/{id}", c.ShowBottle)
			chai.Get(r, "/", c.ListBottles)
		})

		r.Route("/admin", func(r chi.Router) {
			r.Use(auth)

			chai.Post(r, "/auth", c.Auth)
		})

		r.Route("/examples", func(r chi.Router) {
			chai.Get(r, "/ping", c.PingExample)
			chai.Get(r, "/calc", c.CalcExample)
			// chai.Get(r, "/group{s/{gro}up_id}/accounts/{account_id}", c.PathParamsExample)
			chai.Get(r, "/groups/{group_id}/accounts/{account_id}", c.PathParamsExample)
			chai.Get(r, "/header", c.HeaderExample)
			chai.Get(r, "/securities", c.SecuritiesExample)
			chai.Get(r, "/attribute", c.AttributeExample)
			chai.Post(r, "/attribute", c.PostExample)
		})
	})

	// This must be used only during development to generate the swagger spec
	docs, err := openapi2.ChiDocs(r)
	if err != nil {
		panic(fmt.Sprintf("openapi2.Docs() failed: %+v", err))
	}

	// This should be used in prod to serve the swagger spec
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	addCustomDocs(docs)

	LogYAML(docs)

	// This must be used only during development to store the swagger spec
	err = openapi2.WriteDocs(docs, &openapi2.GenConfig{
		OutputDir: "examples/celler/docs",
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
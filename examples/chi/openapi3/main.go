package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ghodss/yaml"

	chai "github.com/go-chai/chai/chi"
	"github.com/go-chai/chai/examples/shared/controller"

	kinopenapi3 "github.com/getkin/kin-openapi/openapi3"

	_ "github.com/go-chai/chai/examples/docs/openapi3" // This is required to be able to serve the stored swagger spec in prod
	"github.com/go-chai/chai/examples/shared/httputil"
	"github.com/go-chai/chai/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	c := controller.NewController()

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/accounts", func(r chi.Router) {
			chai.Get(r, "/{id}", c.ShowAccount)
			chai.Get(r, "/", c.ListAccounts)
			chai.Post(r, "/", c.AddAccount)
			r.Delete("/{id:[0-9]+}", c.DeleteAccount)
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
	docs, err := chai.OpenAPI3(r)
	if err != nil {
		panic(fmt.Sprintf("failed to generate the swagger spec: %+v", err))
	}

	// This should be used in prod to serve the swagger spec
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	addCustomDocs(docs)

	LogYAML(docs)

	// This must be used only during development to store the swagger spec
	err = openapi3.WriteDocs(docs, &openapi3.GenConfig{
		OutputDir: "examples/docs/openapi3",
	})
	if err != nil {
		panic(fmt.Sprintf("failed to write the swagger spec: %+v", err))
	}

	fmt.Println("The swagger spec is available at http://localhost:8080/swagger/")

	http.ListenAndServe(":8080", r)
}

func addCustomDocs(docs *kinopenapi3.T) {
	docs.Servers = []*kinopenapi3.Server{
		{
			URL: "localhost:8080",
		},
	}
	docs.Info = &kinopenapi3.Info{
		ExtensionProps: kinopenapi3.ExtensionProps{},
		Description:    "This is a sample celler server.",
		Title:          "Swagger Example API",
		TermsOfService: "http://swagger.io/terms/",
		Contact: &kinopenapi3.Contact{
			Name:  "API Support",
			URL:   "http://www.swagger.io/support",
			Email: "support@swagger.io",
		},
		License: &kinopenapi3.License{
			Name: "Apache 2.0",
			URL:  "http://www.apache.org/licenses/LICENSE-2.0.html",
		},
		Version: "1.0",
	}

	docs.Components.SecuritySchemes = kinopenapi3.SecuritySchemes{
		"BasicAuth": {
			Value: kinopenapi3.NewJWTSecurityScheme(),
		},
		"ApiKeyAuth": {
			Value: kinopenapi3.NewCSRFSecurityScheme(),
		},
		"OAuth2Implicit": {
			Value: kinopenapi3.NewCSRFSecurityScheme(),
		},
		"OAuth2Application": {
			Value: kinopenapi3.NewCSRFSecurityScheme(),
		},
		"OAuth2Password": {
			Value: kinopenapi3.NewCSRFSecurityScheme(),
		},
		"OAuth2AccessToken": {
			Value: kinopenapi3.NewCSRFSecurityScheme(),
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

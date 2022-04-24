package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"

	chai "github.com/go-chai/chai/chi"
	"github.com/go-chai/chai/examples/shared/controller"
	"github.com/go-chai/chai/log"

	"github.com/go-chai/chai/examples/shared/httputil"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	c := controller.NewController()

	r.Mount("/", c.ChiRoutes())

	docs, err := chai.OpenAPI3(r)
	if err != nil {
		panic(fmt.Sprintf("failed to generate the swagger spec: %+v", err))
	}
	addCustomDocs(docs)

	// Serve the swagger spec
	r.Get("/swagger/*", chai.SwaggerHandler(docs))

	log.YAML(docs)

	fmt.Println("The swagger spec is available at http://localhost:8080/swagger/")

	http.ListenAndServe(":8080", r)
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
	docs.Components.SecuritySchemes = map[string]*openapi3.SecuritySchemeRef{
		"BasicAuth": {Value: &openapi3.SecurityScheme{Type: "basic"}},
		"ApiKeyAuth": {Value: &openapi3.SecurityScheme{
			Type: "apiKey",
			In:   "header",
			Name: "Authorization",
		}},
	}
	// "OAuth2Implicit": {
	// 	SecuritySchemeProps: spec.SecuritySchemeProps{
	// 		Description:      "",
	// 		Type:             "oauth2",
	// 		Flow:             "implicit",
	// 		AuthorizationURL: "https://example.com/oauth/authorize",
	// 		TokenURL:         "",
	// 		Scopes: map[string]string{
	// 			"admin": "Grants read and write access to administrative information",
	// 			"write": "Grants write access",
	// 		},
	// 	},
	// },
	// "OAuth2Application": {
	// 	SecuritySchemeProps: spec.SecuritySchemeProps{
	// 		Description: "Use with the OAuth2 Implicit Grant to retrieve a token",
	// 		Type:        "oauth2",
	// 		Flow:        "application",
	// 		TokenURL:    "https://example.com/oauth/token",
	// 		Scopes: map[string]string{
	// 			"admin": "Grants read and write access to administrative information",
	// 			"write": "Grants write access",
	// 		},
	// 	},
	// },

	// "OAuth2Password": {
	// 	SecuritySchemeProps: spec.SecuritySchemeProps{
	// 		Type:     "oauth2",
	// 		Flow:     "password",
	// 		TokenURL: "https://example.com/oauth/token",
	// 		Scopes: map[string]string{
	// 			"admin": "Grants read and write access to administrative information",
	// 			"write": "Grants write access",
	// 			"read":  "Grants read access",
	// 		},
	// 	},
	// },
	// "OAuth2AccessToken": {
	// 	SecuritySchemeProps: spec.SecuritySchemeProps{
	// 		Type:             "oauth2",
	// 		Flow:             "accessCode",
	// 		AuthorizationURL: "https://example.com/oauth/authorize",
	// 		TokenURL:         "https://example.com/oauth/token",
	// 		Scopes: map[string]string{
	// 			"admin": "Grants read and write access to administrative information",
	// 		},
	// 	},
	// },
	// }
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

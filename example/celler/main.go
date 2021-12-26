package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ghodss/yaml"
	"github.com/go-chai/chai"
	"github.com/go-chai/chai/example/celler/controller"
	_ "github.com/go-chai/chai/example/celler/docs"
	"github.com/go-chai/chai/example/celler/httputil"
	"github.com/go-chai/chai/openapi2"
	"github.com/go-chai/chai/specc"
	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/spec"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization

// @securitydefinitions.oauth2.application  OAuth2Application
// @tokenUrl                                https://example.com/oauth/token
// @scope.write                             Grants write access
// @scope.admin                             Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit  OAuth2Implicit
// @authorizationUrl                     https://example.com/oauth/authorize
// @scope.write                          Grants write access
// @scope.admin                          Grants read and write access to administrative information

// @securitydefinitions.oauth2.password  OAuth2Password
// @tokenUrl                             https://example.com/oauth/token
// @scope.read                           Grants read access
// @scope.write                          Grants write access
// @scope.admin                          Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode  OAuth2AccessCode
// @tokenUrl                               https://example.com/oauth/token
// @authorizationUrl                       https://example.com/oauth/authorize
// @scope.admin                            Grants read and write access to administrative information
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
			chai.Get(r, "/groups/{group_id}/accounts/{account_id}", c.PathParamsExample)
			chai.Get(r, "/header", c.HeaderExample)
			chai.Get(r, "/securities", c.SecuritiesExample)
			chai.Get(r, "/attribute", c.AttributeExample)
			chai.Post(r, "/attribute", c.PostExample)
		})
	})

	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	docs, err := openapi2.Docs(r)
	if err != nil {
		panic(err)
	}

	addCustomDocs(docs)

	// LogJSON(docs, "swagger")
	LogYAML(docs, "")

	http.ListenAndServe(":8080", r)
}

func addCustomDocs(docs *specc.Swagger) {
	docs.Host = "localhost:8080"
	docs.Swagger.Swagger = "2.0"
	docs.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Description:    "This is a sample server celler server.",
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
	// spew.Dump(v)

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
	// spew.Dump(v)

	return
}

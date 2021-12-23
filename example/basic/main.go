package main

import (
	"net/http"

	"github.com/go-chai/chai"
	"github.com/go-chai/chai/example/basic/api"
	"github.com/go-chi/chi/v5"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func main() {
	r := chi.NewRouter()

	chai.Get(r, "/testapi/get-string-by-int/", api.GetStringByInt)
	chai.Get(r, "//testapi/get-struct-array-by-string/", api.GetStructArrayByString)
	r.Post("/testapi/upload", api.Upload)

	http.ListenAndServe(":8080", r)
}

package main

import (
	"net/http"

	"github.com/go-chai/chai"
	"github.com/go-chai/chai/examples/basic/api"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	chai.GetG(r, "/testapi/get-string-by-int/", api.GetStringByInt)
	chai.GetG(r, "//testapi/get-struct-array-by-string/", api.GetStructArrayByString)
	r.Post("/testapi/upload", api.Upload)

	http.ListenAndServe(":8080", r)
}

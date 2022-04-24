package chai

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag"
)

type openAPI struct {
	*openapi3.T
}

func (o *openAPI) ReadDoc() string {
	res, _ := o.T.MarshalJSON()
	return string(res)
}

func SwaggerHandler(docs *openapi3.T) http.HandlerFunc {
	swag.Register(swag.Name, &openAPI{docs})
	return httpSwagger.Handler()
}

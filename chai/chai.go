package chai

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

func init() {
	queryParamDecoder = schema.NewDecoder()
	queryParamDecoder.SetAliasTag("query")

	pathParamDecoder = schema.NewDecoder()
	pathParamDecoder.SetAliasTag("path")

	Validate = validator.New()
}

var queryParamDecoder *schema.Decoder
var pathParamDecoder *schema.Decoder
var Validate *validator.Validate

var DefaultDecoder = func(req any, r *http.Request) ErrType {
	// panic(spew.Sdump(req))
	req = Indirect(req, true)
	// panic(spew.Sdump(req))
	// panic(reflect.TypeOf(req).Elem().Kind() == reflect.Interface)
	if reflect.TypeOf(req).Elem().Kind() == reflect.Interface {
		return nil
	}

	err := decodeQueryParams(r, req)
	if err != nil {
		return err
	}
	err = decodePathParams(r, req)
	if err != nil {
		return err
	}
	err = render.Decode(r, req)
	if err != nil {
		return err
	}

	// TODO populate the header/cookie params
	return nil
}

// modified version of reflect.Indirect
func Indirect(vv any, ptr bool) any {
	v := reflect.ValueOf(vv)
	// If v is a named type and is addressable,
	// start with its address, so that if the type has pointer methods,
	// we find them.
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}
	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() {
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.IsNil() {
			if v.CanSet() {
				v.Set(reflect.New(v.Type().Elem()))
			} else {
				v = reflect.New(v.Type().Elem())
			}
		}
		v = v.Elem()
	}

	if ptr {
		return v.Addr().Interface()
	}

	return v.Interface()
}

func decodeQueryParams(r *http.Request, req any) error {
	queryParams, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return err
	}
	if err := queryParamDecoder.Decode(req, queryParams); err != nil {
		return err
	}
	return nil
}

func decodePathParams(r *http.Request, req any) error {
	routeContext := chi.RouteContext(r.Context())
	if routeContext == nil {
		return nil
	}

	pathParams := make(url.Values)
	for i, key := range routeContext.URLParams.Keys {
		if key == "*" {
			continue
		}
		pathParams[key] = append(pathParams[key], routeContext.URLParams.Values[i])
	}
	if err := pathParamDecoder.Decode(req, pathParams); err != nil {
		return err
	}
	return nil
}

type ValidationError struct {
	Message string                                 `json:"error"`
	Fields  validator.ValidationErrorsTranslations `json:"fields"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

var DefaultValidator = func(req any) ErrType {
	if !reflect.ValueOf(req).IsValid() {
		return nil
	}
	req = Indirect(req, true)
	if reflect.TypeOf(req).Elem().Kind() == reflect.Interface {
		return nil
	}

	err := Validate.Struct(req)
	if err != nil {
		err := err.(validator.ValidationErrors)
		return &ValidationError{
			Message: "validation error",
			Fields:  err.Translate(nil),
		}
	}
	return nil
}

var DefaultResponder = func(w http.ResponseWriter, r *http.Request, code int, res any) {
	if code == 0 {
		code = http.StatusOK
	}
	render.Status(r, code)
	render.Respond(w, r, res)
}

var DefaultErrorResponder = func(w http.ResponseWriter, r *http.Request, code int, e ErrType) {
	if code == 0 {
		code = http.StatusInternalServerError
	}
	ew := &ErrWrap{
		Err:        e,
		StatusCode: code,
		Message:    e.Error(),
	}
	render.Status(r, code)
	render.Respond(w, r, ew)
}

type Methoder interface {
	Method(method, pattern string, h http.Handler)
}

type ErrType = error

type Err error

func handleErr[Err ErrType](w http.ResponseWriter, r *http.Request, err Err, code int, errorFn ErrorResponderFunc) bool {
	if !isErr(err) {
		return false
	}
	errorFn(w, r, code, err)
	return true
}

func isErr[Err ErrType](err Err) bool {
	return !reflect.ValueOf(&err).Elem().IsZero()
}

type ErrWrap struct {
	Err        error
	StatusCode int
	Message    string
}

// TODO figure out how to do this without multiple json.Marshal/Unmarshal calls
func (ew *ErrWrap) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"error":       ew.Message,
		"status_code": ew.StatusCode,
	}
	b, err := json.Marshal(ew.Err)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

type DecoderFunc[Req any] func(*http.Request) (Req, ErrType)
type ResponderFunc[Res any] func(w http.ResponseWriter, r *http.Request, code int, res Res)
type ErrorResponderFunc func(w http.ResponseWriter, r *http.Request, code int, e ErrType)
type ValidatorFunc[Req any] func(req Req) ErrType

func defaultDecoder[Req any](r *http.Request) (Req, ErrType) {
	req := new(Req)
	err := DefaultDecoder(req, r)
	return *req, err
}

func defaultValidator[Req any](req Req) ErrType {
	return DefaultValidator(req)
}

func defaultResponder[Res any](w http.ResponseWriter, r *http.Request, code int, res Res) {
	DefaultResponder(w, r, code, res)
}

func defaultErrorResponder[Err ErrType](w http.ResponseWriter, r *http.Request, code int, err Err) {
	DefaultErrorResponder(w, r, code, err)
}

type Metadata struct {
	Req            any
	Res            any
	Err            any
	Op             *openapi3.Operation
	HandlerFunc    any
	HandlerWrapper http.Handler
}

func addResponse(operation *openapi3.Operation, status int, response *openapi3.Response) {
	responses := operation.Responses
	if responses == nil {
		responses = make(openapi3.Responses)
		operation.Responses = responses
	}
	code := "default"
	if status != 0 {
		code = strconv.FormatInt(int64(status), 10)
	}
	responses[code] = &openapi3.ResponseRef{Value: response}
}

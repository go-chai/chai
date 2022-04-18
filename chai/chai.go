package chai

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"

	"github.com/go-chai/chai/internal/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/zhamlin/chi-openapi/pkg/openapi/operations"
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

	// TODO populate the path/header/cookie params
	return nil
}

// modified version of reflect.indirect
func indirect(vv any) any {
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

	return v.Addr().Interface()
}

func decodeQueryParams(r *http.Request, req any) error {
	queryParams, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return err
	}
	log.Dump(req)
	req = indirect(req)
	log.Dump(req)
	log.Dump(queryParams)

	if err := queryParamDecoder.Decode(req, queryParams); err != nil {
		return err
	}
	return nil
}

func decodePathParams(r *http.Request, req any) error {
	routeContextParams := chi.RouteContext(r.Context()).URLParams
	pathParams := make(url.Values)
	for i, key := range routeContextParams.Keys {
		if key == "*" {
			continue
		}
		pathParams[key] = append(pathParams[key], routeContextParams.Values[i])
	}
	req = indirect(req)
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
	req = indirect(req)
	err := Validate.Struct(req)
	if err != nil {
		log.Dump(err)
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

type Reqer interface {
	Req() any
}

type ResErrer interface {
	Res() any
	Err() any
}

type Handlerer interface {
	Handler() any
}

type Oper interface {
	Op() *operations.Operation
}

type ErrType = error

type Err error

func handleErr[Err ErrType](w http.ResponseWriter, r *http.Request, err Err, code int, errorFn ErrorResponderFunc) bool {
	log.Dump(err)
	if !isErr(err) {
		return false
	}
	log.Dump(err)
	errorFn(w, r, code, err)
	return true
}

func isErr[Err ErrType](err Err) bool {
	log.Dump(reflect.ValueOf(&err))
	log.Dump(reflect.ValueOf(&err).Elem())
	log.Dump(reflect.ValueOf(&err).Elem().IsZero())

	if !reflect.ValueOf(&err).Elem().IsZero() {
		log.Dump(err)
		log.Dump(reflect.ValueOf(err))
		log.Dump(reflect.ValueOf(err).IsZero())
	}
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

func defaultResponder[Res any](w http.ResponseWriter, r *http.Request, code int, res Res) {
	DefaultResponder(w, r, code, res)
}

func defaultErrorResponder[Err ErrType](w http.ResponseWriter, r *http.Request, code int, err Err) {
	DefaultErrorResponder(w, r, code, err)
}

func defaultValidator[Req any](req Req) ErrType {
	return DefaultValidator(req)
}

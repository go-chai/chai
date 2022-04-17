package chai

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"

	"github.com/go-chai/chai/internal/log"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/zhamlin/chi-openapi/pkg/openapi/operations"
)

func init() {
	schemaDecoder = schema.NewDecoder()
	schemaDecoder.SetAliasTag("query")

	Validate = validator.New()
}

var schemaDecoder *schema.Decoder
var Validate *validator.Validate

var DefaultDecoder = func(req any, r *http.Request) ErrType {
	queryParams, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return err
	}
	v := reflect.New(reflect.ValueOf(req).Elem().Type().Elem()).Interface()
	log.Dump(v)
	log.Dump(reflect.ValueOf(v).Elem().Kind())
	if err := schemaDecoder.Decode(v, queryParams); err != nil {
		return err
	}
	err = render.Decode(r, req)
	if err != nil {
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

func handleErr(w http.ResponseWriter, r *http.Request, err ErrType, code int, errorFn ErrorResponderFunc) bool {
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

func defaultResponder[Res any](w http.ResponseWriter, r *http.Request, code int, res Res) {
	DefaultResponder(w, r, code, res)
}

func defaultErrorResponder[Err ErrType](w http.ResponseWriter, r *http.Request, code int, err Err) {
	DefaultErrorResponder(w, r, code, err)
}

func defaultValidator[Req any](req Req) ErrType {
	return DefaultValidator(req)
}

package model

import (
	"net/http"
	"regexp"

	"github.com/go-chai/chai/examples/shared/httputil"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	Zip        string `json:"zip"`
	QueryParam int    `query:"queryParam"`
	PathParam  string `path:"pathParam"`
}

func (a *Address) ValidateStep1() error {
	err := validation.ValidateStruct(a,
		// State cannot be empty, and must be a string consisting of five digits
		validation.Field(&a.Zip, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))),
	)
	if err != nil {
		return &httputil.Error{Message: err.Error(), StatusCode: http.StatusBadRequest}
	}
	return nil
}

func (a *Address) ValidateStep2() error {
	err := validation.ValidateStruct(a,
		// Street cannot be empty, and the length must between 5 and 50
		validation.Field(&a.Street, validation.Required, validation.Length(5, 50)),
		// City cannot be empty, and the length must between 5 and 50
		validation.Field(&a.City, validation.Required, validation.Length(5, 50)),
		// State cannot be empty, and must be a string consisting of two letters in upper case
		validation.Field(&a.State, validation.Required, validation.Match(regexp.MustCompile("^[A-Z]{2}$"))),
		// State cannot be empty, and must be a string consisting of five digits
		validation.Field(&a.Zip, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))),
	)
	if err != nil {
		return &httputil.Error{Message: err.Error(), StatusCode: http.StatusBadRequest}
	}
	return nil
}

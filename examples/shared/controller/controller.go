package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	chai "github.com/go-chai/chai/chi"
	"github.com/go-chai/chai/examples/shared/httputil"
	"github.com/go-chai/chai/examples/shared/model"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Controller example
type Controller struct {
}

// NewController example
func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) ChiRoutes() chi.Router {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/accounts", func(r chi.Router) {
			chai.Get(r, "/{id}", c.ShowAccount).
				ID("show-account").
				Tags("accounts").
				Summary("Show an account")
			chai.Get(r, "/", c.ListAccounts)
			chai.Post(r, "/", c.AddAccount).
				WithValidator(func(a *model.AddAccount) error {
					err := validation.ValidateStruct(&a,
						validation.Field(a.Name, validation.Required),
					)
					if err != nil {
						return &httputil.Error{Message: err.Error(), StatusCode: http.StatusBadRequest}
					}
					return nil
				})
			chai.Post(r, "/v2", c.AddAccount).
				WithValidator((*model.AddAccount).Validation)
			chai.DeleteB(r, "/{id:[0-9]+}", c.DeleteAccount)
			chai.PatchB(r, "/{id}", c.UpdateAccount)
			chai.PostB(r, "/{id}/images", c.UploadAccountImage).
				Operation(UploadAccountImageOperation())
		})

		r.Route("/bottles", func(r chi.Router) {
			chai.Get(r, "/{id}", func(req any, w http.ResponseWriter, r *http.Request) (*model.Bottle, int, error) {
				id := chi.URLParam(r, "id")
				bid, err := strconv.Atoi(id)
				if err != nil {
					return nil, http.StatusBadRequest, err
				}
				bottle, err := model.BottleOne(bid)
				if err != nil {
					return nil, http.StatusNotFound, err
				}
				return bottle, http.StatusOK, nil
			})
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

	return r
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

// Message example
type Message struct {
	Message string `json:"message" example:"message"`
}

type Message2 struct {
	Message string `json:"message" example:"message"`
}

func UploadAccountImageOperation() *openapi3.Operation {
	op := openapi3.NewOperation()
	op.Summary = "Upload an image"
	op.Description = "Upload file"
	op.Tags = []string{"accounts"}

	return op
}

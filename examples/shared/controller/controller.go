package controller

import (
	"errors"
	"net/http"
	"strconv"

	chai "github.com/go-chai/chai/chi"
	"github.com/go-chai/chai/examples/shared/httputil"
	"github.com/go-chai/chai/examples/shared/model"
	"github.com/go-chi/chi/v5"
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
			chai.Get(r, "/{id}", c.ShowAccount)
			chai.Get(r, "/", c.ListAccounts)
			chai.Post(r, "/", c.AddAccount)
			r.Delete("/{id:[0-9]+}", c.DeleteAccount)
			r.Patch("/{id}", c.UpdateAccount)
			r.Post("/{id}/images", c.UploadAccountImage)
		})

		r.Route("/bottles", func(r chi.Router) {
			// ShowBottle godoc
			// @Summary      Show a bottle
			// @Description  get string by ID
			// @ID           get-string-by-int
			// @Tags         bottles
			// @Accept       json
			// @Produce      json
			// @Param        id   path      int  true  "Bottle ID"
			// @Success      200  {object}  model.Bottle
			// @Failure      400  {object}  httputil.Error
			// @Failure      404  {object}  httputil.Error
			// @Failure      500  {object}  httputil.Error
			chai.Get(r, "/{id}", func(w http.ResponseWriter, r *http.Request) (*model.Bottle, int, error) {
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

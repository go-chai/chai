package controller

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/go-chai/chai/examples/shared/httputil"
	"github.com/go-chai/chai/examples/shared/model"
)

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
func (c *Controller) ShowBottle(w http.ResponseWriter, r *http.Request) (*model.Bottle, int, error) {
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
}

// ListBottles godoc
// @Summary      List bottles
// @Description  get bottles
// @Tags         bottles
// @Accept       json
// @Produce      json
// @Success      200,201,202
// @Failure      400,404,500
func (c *Controller) ListBottles(w http.ResponseWriter, r *http.Request) (*[]model.Bottle, int, *httputil.Error) {
	bottles, err := model.BottlesAll()
	if err != nil {
		return nil, http.StatusNotFound, &httputil.Error{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		}
	}
	return &bottles, http.StatusOK, nil
}

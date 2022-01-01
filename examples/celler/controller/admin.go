package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chai/chai/examples/celler/model"
)

// Auth godoc
// @Summary      Auth admin
// @Description  get admin info
// @Tags         accounts,admin
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400  {object}  httputil.HTTPError
// @Failure      401  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Security     ApiKeyAuth
// @Router       /admin/auth [post]
func (c *Controller) Auth(m any, w http.ResponseWriter, r *http.Request) (*model.Admin, int, error) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		return nil, http.StatusBadRequest, errors.New("please set Header Authorization")
	}
	if authHeader != "admin" {
		return nil, http.StatusUnauthorized, fmt.Errorf("this user isn't authorized to operation key=%s expected=admin", authHeader)
	}

	return &model.Admin{
		ID:   1,
		Name: "admin",
	}, http.StatusOK, nil
}

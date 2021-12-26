package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chai/chai"
	"github.com/go-chai/chai/example/celler/model"
	"github.com/go-chi/chi/v5"
)

// PingExample godoc
// @Summary      ping example
// @Description  do ping
// @Tags         example
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "pong"
// @Failure      400  {string}  string  "ok"
// @Failure      404  {string}  string  "ok"
// @Failure      500  {string}  string  "ok"
// @Router       /examples/ping [get]
func (c *Controller) PingExample(w http.ResponseWriter, r *http.Request) (string, int, error) {
	return "pong", http.StatusOK, nil
}

// CalcExample godoc
// @Summary      calc example
// @Description  plus
// @Tags         example
// @Accept       json
// @Produce      json
// @Param        val1  query      int     true  "used for calc"
// @Param        val2  query      int     true  "used for calc"
// @Success      200   {integer}  string  "answer"
// @Failure      400   {string}   string  "ok"
// @Failure      404   {string}   string  "ok"
// @Failure      500   {string}   string  "ok"
// @Router       /examples/calc [get]
func (c *Controller) CalcExample(w http.ResponseWriter, r *http.Request) (string, int, error) {
	val1, err := strconv.Atoi(r.URL.Query().Get("val1"))
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	val2, err := strconv.Atoi(r.URL.Query().Get("val2"))
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	ans := val1 + val2
	return fmt.Sprintf("%d", ans), http.StatusOK, nil
}

// PathParamsExample godoc
// @Summary      path params example
// @Description  path params
// @Tags         example
// @Accept       json
// @Produce      json
// @Param        group_id    path      int     true  "Group ID"
// @Param        account_id  path      int     true  "Account ID"
// @Success      200         {string}  string  "answer"
// @Failure      400         {string}  string  "ok"
// @Failure      404         {string}  string  "ok"
// @Failure      500         {string}  string  "ok"
// @Router       /examples/groups/{group_id}/accounts/{account_id} [get]
func (c *Controller) PathParamsExample(w http.ResponseWriter, r *http.Request) (string, int, error) {
	groupID, err := strconv.Atoi(chi.URLParam(r, "group_id"))
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	accountID, err := strconv.Atoi(chi.URLParam(r, "account_id"))
	if err != nil {
		return "", http.StatusBadRequest, err
	}

	return fmt.Sprintf("group_id=%d account_id=%d", groupID, accountID), http.StatusOK, nil
}

// HeaderExample godoc
// @Summary      custome header example
// @Description  custome header
// @Tags         example
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Authentication header"
// @Success      200            {string}  string  "answer"
// @Failure      400            {string}  string  "ok"
// @Failure      404            {string}  string  "ok"
// @Failure      500            {string}  string  "ok"
// @Router       /examples/header [get]
func (c *Controller) HeaderExample(w http.ResponseWriter, r *http.Request) (string, int, error) {
	return r.Header.Get("Authorization"), http.StatusOK, nil
}

// SecuritiesExample godoc
// @Summary      custome header example
// @Description  custome header
// @Tags         example
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Authentication header"
// @Success      200            {string}  string  "answer"
// @Failure      400            {string}  string  "ok"
// @Failure      404            {string}  string  "ok"
// @Failure      500            {string}  string  "ok"
// @Security     ApiKeyAuth
// @Security     OAuth2Implicit[admin, write]
// @Router       /examples/securities [get]
func (c *Controller) SecuritiesExample(w http.ResponseWriter, r *http.Request) (string, int, error) {
	return "ok", http.StatusOK, nil
}

// AttributeExample godoc
// @Summary      attribute example
// @Description  attribute
// @Tags         example
// @Accept       json
// @Produce      json
// @Param        enumstring  query     string  false  "string enums"    Enums(A, B, C)
// @Param        enumint     query     int     false  "int enums"       Enums(1, 2, 3)
// @Param        enumnumber  query     number  false  "int enums"       Enums(1.1, 1.2, 1.3)
// @Param        string      query     string  false  "string valid"    minlength(5)  maxlength(10)
// @Param        int         query     int     false  "int valid"       minimum(1)    maximum(10)
// @Param        default     query     string  false  "string default"  default(A)
// @Success      200         {string}  string  "answer"
// @Failure      400         {string}  string  "ok"
// @Failure      404         {string}  string  "ok"
// @Failure      500         {string}  string  "ok"
// @Router       /examples/attribute [get]
func (c *Controller) AttributeExample(w http.ResponseWriter, r *http.Request) (string, int, error) {
	return fmt.Sprintf("enumstring=%s enumint=%s enumnumber=%s string=%s int=%s default=%s",
		r.URL.Query().Get("enumstring"),
		r.URL.Query().Get("enumint"),
		r.URL.Query().Get("enumnumber"),
		r.URL.Query().Get("string"),
		r.URL.Query().Get("int"),
		r.URL.Query().Get("default"),
	), http.StatusOK, nil
}

// PostExample godoc
// @Summary      post request example
// @Description  post request example
// @Accept       json
// @Produce      plain
// @Param        message  body      model.Account  true  "Account Info"
// @Success      200      {string}  string         "success"
// @Router       /examples/post [post]
func (c *Controller) PostExample(account *model.Account, w http.ResponseWriter, r *http.Request) (string, int, *chai.APIError) {
	return account.Name, http.StatusOK, nil
}

package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chai/chai/example/celler/httputil"
	"github.com/go-chai/chai/example/celler/model"
	"github.com/go-chi/chi/v5"
)

// ShowAccount godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  model.Account
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /accounts/{id} [get]
func (c *Controller) ShowAccount(w http.ResponseWriter, r *http.Request) (*model.Account, int, error) {
	id := chi.URLParam(r, "id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	account, err := model.AccountOne(aid)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return &account, http.StatusOK, nil
}

// ListAccounts godoc
// @Summary      List accounts
// @Description  get accounts
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        q    query     string  false  "name search by q"  Format(email)
// @Success      201  {array}   model.Account
// @Success      202  {array}   model.Account
// @Success      203  {array}   model.Account
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /accounts [get]
func (c *Controller) ListAccounts(w http.ResponseWriter, r *http.Request) (*[]model.Account, int, error) {
	q := r.URL.Query().Get("q")
	accounts, err := model.AccountsAll(q)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return &accounts, http.StatusOK, nil
}

// AddAccount godoc
// @Summary      Add an account
// @Description  add by json account
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        account  body      model.AddAccount  true  "Add account"
// @Success      200      {object}  model.Account
// @Failure      400      {object}  httputil.HTTPError
// @Failure      404      {object}  httputil.HTTPError
// @Failure      500      {object}  httputil.HTTPError
// @Router       /accounts [post]
func (c *Controller) AddAccount(addAccount *model.AddAccount, w http.ResponseWriter, r *http.Request) (*model.Account, int, error) {
	if err := addAccount.Validation(); err != nil {
		return nil, http.StatusBadRequest, err
	}
	account := model.Account{
		Name: addAccount.Name,
	}
	lastID, err := account.Insert()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	account.ID = lastID

	return &account, http.StatusOK, nil
}

// UpdateAccount godoc
// @Summary      Update an account
// @Description  Update by json account
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id       path      int                  true  "Account ID"
// @Param        account  body      model.UpdateAccount  true  "Update account"
// @Success      200      {object}  model.Account
// @Failure      400      {object}  httputil.HTTPError
// @Failure      404      {object}  httputil.HTTPError
// @Failure      500      {object}  httputil.HTTPError
// @Router       /accounts/{id} [patch]
func (c *Controller) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	aid, err := strconv.Atoi(id)
	if httputil.NewError(w, http.StatusBadRequest, err) {
		return
	}
	var updateAccount model.UpdateAccount
	err = httputil.Decode(r, updateAccount)
	if httputil.NewError(w, http.StatusBadRequest, err) {
		return
	}
	account := model.Account{
		ID:   aid,
		Name: updateAccount.Name,
	}
	err = account.Update()
	if httputil.NewError(w, http.StatusNotFound, err) {
		return
	}
	httputil.Respond(w, r, account)
}

// DeleteAccount godoc
// @Summary      Delete an account
// @Description  Delete by account ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"  Format(int64)
// @Success      204  {object}  model.Account
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /accounts/{id} [delete]
func (c *Controller) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	aid, err := strconv.Atoi(id)
	if httputil.NewError(w, http.StatusBadRequest, err) {
		return
	}
	err = model.Delete(aid)
	if httputil.NewError(w, http.StatusNotFound, err) {
		return
	}
	httputil.Respond(w, r, struct{}{})
}

// UploadAccountImage godoc
// @Summary      Upload account image
// @Description  Upload file
// @Tags         accounts
// @Accept       multipart/form-data
// @Produce      json
// @Param        id    path      int   true  "Account ID"
// @Param        file  formData  file  true  "account image"
// @Success      200   {object}  controller.Message
// @Failure      400   {object}  httputil.HTTPError
// @Failure      404   {object}  httputil.HTTPError
// @Failure      500   {object}  httputil.HTTPError
// @Router       /accounts/{id}/images [post]
func (c *Controller) UploadAccountImage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if httputil.NewError(w, http.StatusBadRequest, err) {
		return
	}
	_, h, err := r.FormFile("file")
	if httputil.NewError(w, http.StatusBadRequest, err) {
		return
	}

	httputil.Respond(w, r, Message{Message: fmt.Sprintf("upload complete userID=%d filename=%s", id, h.Filename)})
}

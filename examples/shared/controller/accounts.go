package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chai/chai/examples/shared/httputil"
	"github.com/go-chai/chai/examples/shared/model"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

func (c *Controller) ShowAccount(req any, w http.ResponseWriter, r *http.Request) (*model.Account, int, error) {
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

func (c *Controller) ListAccounts(req any, w http.ResponseWriter, r *http.Request) (*[]model.Account, int, error) {
	q := r.URL.Query().Get("q")
	accounts, err := model.AccountsAll(q)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return &accounts, http.StatusOK, nil
}

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

func (c *Controller) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	aid, err := strconv.Atoi(id)
	if httputil.NewError(w, http.StatusBadRequest, err) {
		return
	}
	updateAccount := new(model.UpdateAccount)
	err = httputil.Decode(r, updateAccount)
	if httputil.NewError(w, http.StatusBadRequest, err) {
		return
	}
	account := model.Account{
		ID:   aid,
		Name: updateAccount.Name,
		UUID: uuid.Must(uuid.NewV4()),
	}
	err = account.Update()
	if httputil.NewError(w, http.StatusNotFound, err) {
		return
	}
	httputil.Respond(w, r, account)
}

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

// Admin example
type Admin struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"admin name"`
}

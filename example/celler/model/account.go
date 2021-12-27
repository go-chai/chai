package model

import (
	"errors"
	"fmt"

	uuid "github.com/gofrs/uuid"
)

type CustomType456 struct {
	ZZ    int    `json:"zz" example:"34"`
	ZZStr string `json:"zzStr" example:"zz string value"`
}

// Account example
type Account2 struct {
	ID             int             `json:"id" example:"1" format:"int64"`
	Name           string          `json:"name" example:"account name"`
	SomeNumber     int             `json:"some_number" example:"1234"`
	CustomVal123   string          `json:"custom_val_123" example:"custom val 123"`
	CustomType4567 []CustomType456 `json:"custom_type_4567"`
	UUID           uuid.UUID       `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
}

type Account struct {
	ID             int             `json:"id" example:"1" format:"int64"`
	Name           string          `json:"name" example:"account name"`
	SomeNumber     int             `json:"some_number" example:"1234"`
	CustomVal123   string          `json:"custom_val_123" example:"custom val 123"`
	CustomType4567 []CustomType456 `json:"custom_type_4567"`
	UUID           uuid.UUID       `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
}

//  example
var (
	ErrNameInvalid = errors.New("name is empty")
)

// AddAccount example
type AddAccount struct {
	Name string `json:"name" example:"account name"`
}

// Validation example
func (a AddAccount) Validation() error {
	switch {
	case len(a.Name) == 0:
		return ErrNameInvalid
	default:
		return nil
	}
}

// UpdateAccount example
type UpdateAccount struct {
	Name string `json:"name" example:"account name"`
}

// Validation example
func (a UpdateAccount) Validation() error {
	switch {
	case len(a.Name) == 0:
		return ErrNameInvalid
	default:
		return nil
	}
}

// AccountsAll example
func AccountsAll(q string) ([]Account, error) {
	if q == "" {
		return accounts, nil
	}
	as := []Account{}
	for k, v := range accounts {
		if q == v.Name {
			as = append(as, accounts[k])
		}
	}
	return as, nil
}

// AccountOne example
func AccountOne(id int) (Account, error) {
	for _, v := range accounts {
		if id == v.ID {
			return v, nil
		}
	}
	return Account{}, ErrNoRow
}

// Insert example
func (a Account) Insert() (int, error) {
	accountMaxID++
	a.ID = accountMaxID
	a.Name = fmt.Sprintf("account_%d", accountMaxID)
	accounts = append(accounts, a)
	return accountMaxID, nil
}

// Delete example
func Delete(id int) error {
	for k, v := range accounts {
		if id == v.ID {
			accounts = append(accounts[:k], accounts[k+1:]...)
			return nil
		}
	}
	return fmt.Errorf("account id=%d is not found", id)
}

// Update example
func (a Account) Update() error {
	for k, v := range accounts {
		if a.ID == v.ID {
			accounts[k].Name = a.Name
			return nil
		}
	}
	return fmt.Errorf("account id=%d is not found", a.ID)
}

var accountMaxID = 3
var accounts = []Account{
	{ID: 1, Name: "account_1"},
	{ID: 2, Name: "account_2"},
	{ID: 3, Name: "account_3"},
}

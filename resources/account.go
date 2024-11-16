package resources

import (
	"errors"
	"net/http"
	"transactions/models"
)

type AccountResponse struct {
	*models.Account
}

type NewAccountRequest struct {
	*models.Account
	ProtectedID string `json:"id"` // override
}

func (a *NewAccountRequest) Bind(r *http.Request) error {
	if a.Account == nil {
		return errors.New("missing required Article fields")
	}

	if a.Account.DocumentNumber == "" {
		return errors.New("missing required DocumentNumber field")
	}
	a.ProtectedID = ""

	return nil
}

func NewAccountResponse(account *models.Account) *AccountResponse {
	resp := &AccountResponse{Account: account}
	return resp
}
func (rd *AccountResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"transactions/models"

	"github.com/stretchr/testify/assert"
)

func createAccountFixture(t *testing.T, ts *httptest.Server) *models.Account {
	res, err := http.Post(
		ts.URL+"/account", "application/json",
		bytes.NewBufferString(`{"document_number":"1234567890"}`),
	)
	if err != nil {
		t.Fatal(err.Error())
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, res.StatusCode)
	}
	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	account := models.Account{}
	json.Unmarshal(resBody, &account)
	assert.NoError(t, json.Unmarshal(resBody, &account))

	return &account
}

func createTransactionFixture(t *testing.T, ts *httptest.Server, accountID int64) models.Transaction {
	res, err := http.Post(
		ts.URL+"/transaction", "application/json",
		bytes.NewBufferString(
			fmt.Sprintf(`{"account_id":%d,"operation_type_id":1,"amount":100.0}`, accountID),
		),
	)
	if err != nil {
		t.Fatal(err.Error())
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, res.StatusCode)
	}
	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	transaction := models.Transaction{}
	json.Unmarshal(resBody, &transaction)
	assert.NoError(t, json.Unmarshal(resBody, &transaction))

	return transaction
}

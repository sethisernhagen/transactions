package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"transactions/models"
	"transactions/stores"

	"github.com/stretchr/testify/assert"
)

func TestTransaction_Create(t *testing.T) {
	s := NewServer()
	_, err := stores.GetTestDB(t)
	assert.NoError(t, err)

	ts := httptest.NewServer(s.Router)
	defer ts.Close()
	acct := createAccountFixture(t, ts)

	res, err := http.Post(
		ts.URL+"/transaction", "application/json",
		bytes.NewBufferString(fmt.Sprintf(`{"account_id":%d,"operation_type": "Purchase","amount":100.0}`, acct.AccountID)),
	)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	transaction := models.Transaction{}
	assert.NoError(t, json.Unmarshal(resBody, &transaction))
	assert.IsType(t, int64(0), transaction.TransactionID)
	assert.Equal(t, acct.AccountID, transaction.AccountID)
	assert.Equal(t, models.OperationTypeID(1), transaction.OperationTypeID)
	assert.Equal(t, 100.0, transaction.Amount)
	assert.True(t, transaction.CreatedAt.After(time.Unix(0, 0)), "transaction.CreatedAt should be after Unix epoch, got %v", transaction.CreatedAt)

}

func TestTransaction_Retrieve(t *testing.T) {
	s := NewServer()
	_, err := stores.GetTestDB(t)
	assert.NoError(t, err)
	ts := httptest.NewServer(s.Router)
	defer ts.Close()
	acct := createAccountFixture(t, ts)
	trx := createTransactionFixture(t, ts, acct.AccountID)

	res, err := http.Get(ts.URL + "/transaction/" + fmt.Sprint(trx.TransactionID))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	transaction := models.Transaction{}
	assert.NoError(t, json.Unmarshal(resBody, &transaction))
	assert.Equal(t, trx, transaction)
}

func TestTransaction_Retrieve_NotFound(t *testing.T) {
	s := NewServer()
	_, err := stores.GetTestDB(t)
	assert.NoError(t, err)
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/transaction/1")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestTransaction_Retrieve_BadRequest(t *testing.T) {
	s := NewServer()
	_, err := stores.GetTestDB(t)
	assert.NoError(t, err)
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/transaction/abc")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestTransaction_Create_BadRequest(t *testing.T) {
	s := NewServer()
	_, err := stores.GetTestDB(t)
	assert.NoError(t, err)
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	res, err := http.Post(ts.URL+"/transaction", "application/json", bytes.NewBufferString(`{}`))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestTransaction_Create_BadOperationType(t *testing.T) {
	s := NewServer()
	_, err := stores.GetTestDB(t)
	assert.NoError(t, err)
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	res, err := http.Post(ts.URL+"/transaction", "application/json", bytes.NewBufferString(`{"account_id":1,"operation_type":"Unknown","amount":100.0}`))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

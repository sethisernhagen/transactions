package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
	"transactions/models"
	"transactions/stores"

	"github.com/stretchr/testify/assert"
)

func TestAccount_Create(t *testing.T) {
	s := NewServer()
	_, err := stores.GetTestDB(t)
	assert.NoError(t, err)
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	res, err := http.Post(
		ts.URL+"/account", "application/json",
		bytes.NewBufferString(`{"document_number":"1234567890"}`),
	)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	account := models.Account{}
	assert.NoError(t, json.Unmarshal(resBody, &account))
	assert.IsType(t, int64(0), account.AccountID)
	assert.Equal(t, "1234567890", account.DocumentNumber)
	assert.True(t, account.CreatedAt.After(time.Unix(0, 0)))
}

func TestAccount_Retrieve(t *testing.T) {
	s := NewServer()
	_, err := stores.GetTestDB(t)
	assert.NoError(t, err)
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	res, err := http.Post(
		ts.URL+"/account", "application/json",
		bytes.NewBufferString(`{"document_number":"1234567890"}`),
	)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	account := models.Account{}
	assert.NoError(t, json.Unmarshal(resBody, &account))

	res, err = http.Get(ts.URL + "/account/" + strconv.FormatInt(account.AccountID, 10))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	resBody, err = io.ReadAll(res.Body)
	assert.NoError(t, err)
	resultAccount := models.Account{}
	assert.NoError(t, json.Unmarshal(resBody, &resultAccount))
	assert.Equal(t, account, resultAccount)
}

package stores

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountStore_Create(t *testing.T) {
	db, err := GetTestDB(t)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	assert.NoError(t, err)
	store := NewAccountStore(db)

	createdAccount, err := store.Create(accountFixture1)
	assert.NoError(t, err)
	assert.Equal(t, accountFixture1.AccountID, createdAccount.AccountID)
}

func TestAccountStore_Retrieve(t *testing.T) {
	db, err := GetTestDB(t)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	assert.NoError(t, err)
	store := NewAccountStore(db)

	acct, err := store.Create(accountFixture1)
	assert.NoError(t, err)

	retrievedAccount, err := store.Retrieve(acct.AccountID)
	assert.NoError(t, err)
	assert.Equal(t, accountFixture1, retrievedAccount)
}

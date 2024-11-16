package server

import (
	"database/sql"
	"testing"
)

// TODO: maybe should be in a fixtures package or stores
func truncateTables(t *testing.T, db *sql.DB) {
	_, err := db.Exec("DELETE FROM transaction")
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = db.Exec("DELETE FROM account")
	if err != nil {
		t.Fatal(err.Error())
	}
}

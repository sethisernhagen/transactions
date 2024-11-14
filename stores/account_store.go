package stores

import (
	"database/sql"

	"transactions/models"

	sq "github.com/Masterminds/squirrel"
)

type AccountStore struct {
	db          *sql.DB
	insertQuery sq.InsertBuilder
	selectQuery sq.SelectBuilder
}

func NewAccountStore(db *sql.DB) *AccountStore {
	return &AccountStore{
		db: db,
		insertQuery: psql.Insert("account").
			Columns("document_number").
			Suffix("RETURNING \"id\", \"created_at\""),
		selectQuery: psql.Select("id", "document_number, created_at").
			From("account"),
	}
}

func (s *AccountStore) Create(account *models.Account) (*models.Account, error) {
	err := s.insertQuery.Values(
		account.DocumentNumber,
	).RunWith(s.db).QueryRow().Scan(&account.AccountID, &account.CreatedAt)

	return account, err
}

func (s *AccountStore) Retrieve(id int64) (*models.Account, error) {
	row := s.selectQuery.Where(sq.Eq{"id": id}).RunWith(s.db).QueryRow()
	account := &models.Account{}
	err := row.Scan(&account.AccountID, &account.DocumentNumber, &account.CreatedAt)

	return account, err
}

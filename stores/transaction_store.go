package stores

import (
	"database/sql"
	"transactions/models"

	sq "github.com/Masterminds/squirrel"
)

type TransactionStore struct {
	db          *sql.DB
	insertQuery sq.InsertBuilder
	selectQuery sq.SelectBuilder
	updateQuery sq.UpdateBuilder
}

func NewTransactionStore(db *sql.DB) *TransactionStore {
	return &TransactionStore{
		db: db,
		insertQuery: psql.Insert("transaction").
			Columns("account_id", "amount", "operation_type", "balance").
			Suffix("RETURNING \"id\", \"created_at\""),
		selectQuery: psql.Select("id", "account_id", "amount", "operation_type, created_at, balance").
			From("transaction"),
		updateQuery: psql.Update("transaction"),
	}
}

func (s *TransactionStore) Create(transaction *models.Transaction) (*models.Transaction, error) {
	transaction.Balance = transaction.Amount
	if transaction.OperationTypeID != 1 {
		// calc balance
		s.getPaymentBalance(transaction)
	}

	err := s.insertQuery.Values(
		transaction.AccountID,
		transaction.Amount,
		transaction.OperationTypeID,
		transaction.Balance,
	).RunWith(s.db).QueryRow().Scan(&transaction.TransactionID, &transaction.CreatedAt)

	return transaction, err
}

func (s *TransactionStore) Retrieve(id int64) (*models.Transaction, error) {
	row := s.selectQuery.Where(sq.Eq{"id": id}).RunWith(s.db).QueryRow()
	transaction := &models.Transaction{}
	err := row.Scan(&transaction.TransactionID, &transaction.AccountID, &transaction.Amount, &transaction.OperationTypeID, &transaction.CreatedAt, &transaction.Balance)
	return transaction, err
}

func (s *TransactionStore) Search(account_id int64) ([]models.Transaction, error) {

	q := s.selectQuery.Where(sq.Eq{"account_id": account_id}, sq.Lt{"balance": 0}).OrderBy("id")
	sql, params, err := q.ToSql()
	if err != nil {
		return []models.Transaction{}, err
	}

	rows, err := s.db.Query(sql, params...)
	result := []models.Transaction{}
	for rows.Next() {
		t := models.Transaction{}
		err := rows.Scan(&t.TransactionID, &t.AccountID, &t.Amount, &t.OperationTypeID, &t.CreatedAt, &t.Balance)
		if err != nil {
			return result, err
		}
		result = append(result, t)
	}
	return result, err
}

func (s *TransactionStore) Update(trans models.Transaction) (models.Transaction, error) {
	_, err := s.updateQuery.Set("balance", trans.Balance).Where(sq.Eq{"id": trans.TransactionID}).RunWith(s.db).Exec()

	return trans, err
}

func (s *TransactionStore) getPaymentBalance(transaction *models.Transaction) (*models.Transaction, error) {
	undischarged, err := s.Search(transaction.AccountID)
	if err != nil {
		return transaction, err
	}

	// iterate through transactions that need discharge, sub from trans
	for _, undischarged_item := range undischarged {

		var dischargeAmount float64
		if transaction.Balance+undischarged_item.Balance > 0 {
			dischargeAmount = undischarged_item.Balance * -1
		} else {
			dischargeAmount = transaction.Balance
		}
		undischarged_item.Balance = undischarged_item.Balance + dischargeAmount
		s.Update(undischarged_item)
		transaction.Balance = transaction.Balance - dischargeAmount
		if transaction.Balance == 0 {
			break
		}
	}

	return transaction, err
}

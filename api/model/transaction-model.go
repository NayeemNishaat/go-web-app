package model

import (
	"context"
	"time"
)

// Transaction is the type for transactions
type Transactions struct {
	ID                  int       `json:"id"`
	Amount              int       `json:"amount"`
	Currency            string    `json:"currency"`
	LastFour            string    `json:"last_four"`
	ExpiryMonth         int       `json:"expiry_month"`
	ExpiryYear          int       `json:"expiry_year"`
	PaymentIntent       string    `json:"payment_intent"`
	PaymentMethod       string    `json:"payment_method"`
	BankReturnCode      string    `json:"bank_return_code"`
	TransactionStatusID int       `json:"transaction_status_id"`
	CreatedAt           time.Time `json:"-"`
	UpdatedAt           time.Time `json:"-"`
}

// InsertTransaction inserts a new txn, and returns its id
func (m *SqlDB) InsertTransaction(txn Transactions) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	stmt := `
		INSERT INTO transactions
			(amount, currency, last_four, bank_return_code, expiry_month, expiry_year, payment_intent, payment_method, transaction_status_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id;`

	err := m.QueryRow(ctx, stmt, txn.Amount, txn.Currency, txn.LastFour, txn.BankReturnCode, txn.ExpiryMonth, txn.ExpiryYear, txn.PaymentIntent, txn.PaymentMethod, txn.TransactionStatusID, time.Now(), time.Now()).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

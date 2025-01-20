package data

import (
	"context"
	"database/sql"
	"time"
)

type WalletTransaction struct {
	Id               int64     `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	Title            string    `json:"title"`
	Category         string    `json:"category"`
	Account          string    `json:"account"`
	Amount           float64   `json:"amount"`
	Currency         string    `json:"currency"`
	TransactionType  string    `json:"type"`
	TransferAmount   float64   `json:"transferAmount"`
	TransferCurrency string    `json:"transferCurrency"`
	ToAccount        string    `json:"toAccount"`
	ReceiveAmount    float64   `json:"receiveAmount"`
	ReceiveCurrency  string    `json:"receiveCurrency"`
	Description      string    `json:"description"`
	DueDate          time.Time `json:"dueDate"`
	WalletId         string    `json:"walletId"`
}

type WalletTransactionModel struct {
	DB *sql.DB
}

func (m WalletTransactionModel) Insert(walletTransaction *WalletTransaction) error {
	query := `
		INSERT INTO WalletTransaction (created_at, title, category, account, amount, currency,
		                               transactionType, transferamount, transfercurrency, toaccount,
		                               receiveamount, receivecurrency, description,
		                               duedate, walletid)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id`

	args := []any{
		walletTransaction.CreatedAt,
		walletTransaction.Title,
		walletTransaction.Category,
		walletTransaction.Account,
		walletTransaction.Amount,
		walletTransaction.Currency,
		walletTransaction.TransactionType,
		walletTransaction.TransferAmount,
		walletTransaction.TransferCurrency,
		walletTransaction.ToAccount,
		walletTransaction.ReceiveAmount,
		walletTransaction.ReceiveCurrency,
		walletTransaction.Description,
		walletTransaction.DueDate,
		walletTransaction.WalletId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&walletTransaction.Id)
}

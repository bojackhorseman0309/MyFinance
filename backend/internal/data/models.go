package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	WalletTransactions WalletTransactionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		WalletTransactions: WalletTransactionModel{DB: db},
	}
}

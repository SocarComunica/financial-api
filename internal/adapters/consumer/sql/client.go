package sql

import "github.com/socarcomunica/financial-api/internal/application/transaction"

type Client interface {
	transaction.TransactionsDatabase
}

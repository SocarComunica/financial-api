package sql

import (
	"github.com/socarcomunica/financial-api/internal/application/account"
	"github.com/socarcomunica/financial-api/internal/application/transaction"
)

type Client interface {
	transaction.TransactionsDatabase
	account.AccountsDatabase
}

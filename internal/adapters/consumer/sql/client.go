package sql

import (
	"github.com/socarcomunica/financial-api/internal/application/account"
	"github.com/socarcomunica/financial-api/internal/application/transaction"
	"github.com/socarcomunica/financial-api/internal/application/user"
	"github.com/socarcomunica/financial-api/internal/domain"

	"github.com/labstack/gommon/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Client interface {
	transaction.TransactionsDatabase
	account.AccountsDatabase
	user.UsersDatabase
}

type client struct {
	DB *gorm.DB
}

func NewClient(database string) Client {
	db, err := gorm.Open(sqlite.Open(database), &gorm.Config{})
	if err != nil {
		panic("failed to connect to local database")
	}

	if err := db.AutoMigrate(&domain.Account{}, &domain.Tag{}, &domain.Transaction{}, &domain.User{}); err != nil {
		log.Error(err)
	}

	return &client{
		DB: db,
	}
}

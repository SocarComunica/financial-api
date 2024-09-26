package sql

import (
	"errors"

	"github.com/labstack/gommon/log"
	"github.com/socarcomunica/financial-api/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	AddTransactionError = "AddTransactionError Local Client: "
)

type localClient struct {
	DB *gorm.DB
}

func NewLocalClient() Client {
	db, err := gorm.Open(sqlite.Open("local.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to local database")
	}

	if err := db.AutoMigrate(&domain.Transaction{}, &domain.Tag{}); err != nil {
		log.Error(err)
	}

	return &localClient{
		DB: db,
	}
}

func (l *localClient) AddTransaction(model *domain.Transaction) (*domain.Transaction, error) {
	result := l.DB.Create(model)

	if result.Error != nil {
		return nil, errors.New(AddTransactionError + result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return nil, errors.New(AddTransactionError + "no new transactions were created")
	}

	return model, nil
}

package sql

import (
	"errors"

	"github.com/labstack/gommon/log"
	"github.com/socarcomunica/financial-api/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	AddTransactionError         = "AddTransactionError Local Client: "
	CreateAccountError          = "CreateAccountError Local Client: "
	ErrorUpdatingAccountBalance = "error updating account balance: "
)

type localClient struct {
	DB *gorm.DB
}

func NewLocalClient() Client {
	db, err := gorm.Open(sqlite.Open("local.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to local database")
	}

	if err := db.AutoMigrate(&domain.Transaction{}, &domain.Tag{}, &domain.Account{}); err != nil {
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

func (l *localClient) AddAccount(model *domain.Account) (*domain.Account, error) {
	result := l.DB.Create(model)

	if result.Error != nil {
		return nil, errors.New(CreateAccountError + result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return nil, errors.New(CreateAccountError + "no new accounts were created")
	}

	return model, nil
}

func (l *localClient) GetAccount(id uint) (*domain.Account, error) {
	var account domain.Account

	result := l.DB.Where("id = ?", id).First(&account)

	if result.Error != nil {
		return nil, result.Error
	}

	return &account, nil
}

func (l *localClient) UpdateAccountBalance(account *domain.Account) error {
	result := l.DB.Save(account)
	if result.Error != nil {
		return errors.New(ErrorUpdatingAccountBalance + result.Error.Error())
	}

	return nil
}

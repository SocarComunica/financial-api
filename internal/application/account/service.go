package account

import (
	"errors"
	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
	"github.com/socarcomunica/financial-api/internal/domain"
)

const (
	CreateAccountError = "CreateAccount Service Error: "
)

type AccountsDatabase interface {
	AddAccount(model *domain.Account) (*domain.Account, error)
	GetAccount(id uint) (*domain.Account, error)
	GetAccountsByUser(userID uint) ([]*domain.Account, error)
}

type Service struct {
	Database AccountsDatabase
}

func NewAccountsDatabase(database AccountsDatabase) *Service {
	return &Service{
		Database: database,
	}
}

func (a *Service) AddAccount(request request.CreateAccount) (*domain.Account, error) {
	accountModel := &domain.Account{
		Name:           request.Name,
		Type:           request.Type,
		Balance:        request.Balance,
		InitialBalance: request.Balance,
	}

	account, err := a.Database.AddAccount(accountModel)
	if err != nil {
		return nil, errors.New(CreateAccountError + err.Error())
	}
	return account, nil
}

func (a *Service) GetAccountsByUser(userID uint) ([]*domain.Account, error) {
	return a.Database.GetAccountsByUser(userID)
}

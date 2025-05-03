package account

import (
	"errors"

	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
	"github.com/socarcomunica/financial-api/internal/domain"
)

const (
	CreateAccountError = "CreateAccount Service Error: "
	DeleteAccountError = "DeleteAccount Service Error: "
)

type AccountsDatabase interface {
	AddAccount(model *domain.Account) (*domain.Account, error)
	GetAccount(id uint) (*domain.Account, error)
	GetAccountsByUser(userID uint) ([]*domain.Account, error)
	DeleteAccount(id uint, userID uint) error
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
		UserID:         request.UserID,
	}

	account, err := a.Database.AddAccount(accountModel)
	if err != nil {
		return nil, errors.New(CreateAccountError + err.Error())
	}
	return account, nil
}

func (a *Service) GetAccountsByUser(userID uint) ([]*domain.Account, error) {
	accounts, err := a.Database.GetAccountsByUser(userID)
	if err != nil {
		return nil, err // Propagate database errors
	}
	if accounts == nil {
		return []*domain.Account{}, nil // Return empty slice if no accounts found
	}
	return accounts, nil
}

func (a *Service) DeleteAccount(id uint, userID uint) error {
	// First verify that the account exists and belongs to the user
	account, err := a.Database.GetAccount(id)
	if err != nil {
		return errors.New(DeleteAccountError + err.Error())
	}

	if account.UserID != userID {
		return errors.New(DeleteAccountError + "account does not belong to the specified user")
	}

	return a.Database.DeleteAccount(id, userID)
}

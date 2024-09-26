package transaction

import (
	"errors"
	"github.com/labstack/gommon/log"
	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
	"github.com/socarcomunica/financial-api/internal/domain"
)

const (
	AddTransactionError = "AddTransaction Service Error: "
)

type TransactionsDatabase interface {
	AddTransaction(model *domain.Transaction) (*domain.Transaction, error)
	GetAccount(id uint) (*domain.Account, error)
	UpdateAccountBalance(account *domain.Account) error
}

type Service struct {
	Database TransactionsDatabase
}

func NewTransactionService(database TransactionsDatabase) *Service {
	return &Service{
		Database: database,
	}
}

func (t *Service) AddTransaction(request request.CreateTransaction) (*domain.Transaction, error) {

	account, err := t.Database.GetAccount(request.OriginID)
	if err != nil {
		return nil, errors.New(AddTransactionError + err.Error())
	}

	if account.Balance < request.Amount && request.Type.Name == "debit" {
		return nil, errors.New(AddTransactionError + "insufficient funds")
	}

	transactionModel := &domain.Transaction{
		Amount:      request.Amount,
		Description: request.Description,
		Tags: func() []domain.Tag {
			var tags []domain.Tag
			for _, tag := range request.Tags {
				tags = append(tags, domain.Tag{
					Name: tag.Name,
				})
			}
			return tags
		}(),
		Type:          domain.Type{Name: request.Type.Name},
		OriginID:      request.OriginID,
		Origin:        *account,
		DestinationID: request.DestinationID,
	}

	var destination *domain.Account
	if request.DestinationID != nil {
		destination, err = t.Database.GetAccount(*request.DestinationID)
		if err != nil {
			return nil, errors.New(AddTransactionError + err.Error())
		}
		transactionModel.Destination = destination
	}

	transaction, err := t.Database.AddTransaction(transactionModel)
	if err != nil {
		return nil, errors.New(AddTransactionError + err.Error())
	}

	go t.updateOriginAndDestinationBalance(request, account, destination)

	return transaction, nil
}

func (t *Service) updateOriginAndDestinationBalance(request request.CreateTransaction, account *domain.Account, destination *domain.Account) {
	if request.Type.Name == "credit" {
		account.Balance += request.Amount
	} else {
		account.Balance -= request.Amount
	}
	if err := t.Database.UpdateAccountBalance(account); err != nil {
		log.Error("error updating origin account balance: ", err)
	}

	if destination != nil && request.Type.Name == "transfer" {
		destination.Balance += request.Amount
		if err := t.Database.UpdateAccountBalance(destination); err != nil {
			log.Error("error updating destination account balance: ", err)
		}
	}
}

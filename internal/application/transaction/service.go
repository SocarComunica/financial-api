package transaction

import (
	"errors"
	"github.com/labstack/gommon/log"
	"github.com/socarcomunica/financial-api/common"
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
	GetTransactionsByAccount(accountID uint, offset int) ([]*domain.Transaction, error)
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

	if account.Balance < request.Amount && request.Type.Name == common.TransactionTypeDebit {
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
	switch request.Type.Name {
	case common.TransactionTypeCredit:
		account.Balance += request.Amount
		break
	case common.TransactionTypeDebit:
		account.Balance -= request.Amount
		break
	case common.TransactionTypeTransfer:
		account.Balance -= request.Amount
		destination.Balance += request.Amount
		if err := t.Database.UpdateAccountBalance(destination); err != nil {
			log.Error("error updating destination account balance: ", err)
		}
		break
	}
	if err := t.Database.UpdateAccountBalance(account); err != nil {
		log.Error("error updating origin account balance: ", err)
	}
}

func (t *Service) GetTransactionsByAccount(accountID uint, offset int) ([]*domain.Transaction, error) {
	return t.Database.GetTransactionsByAccount(accountID, offset)
}

package transaction

import (
	"errors"

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
	UpdateAccountBalances(origin *domain.Account, destination *domain.Account) error
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
	origin, err := t.Database.GetAccount(request.OriginID)
	if err != nil {
		return nil, errors.New(AddTransactionError + err.Error())
	}

	if origin.Balance < request.Amount && request.Type == common.TransactionTypeDebit {
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
		Type:          request.Type,
		OriginID:      request.OriginID,
		Origin:        *origin,
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

	// Update balances in a single transaction
	if err := t.updateBalances(request, origin, destination); err != nil {
		return nil, err
	}

	transaction, err := t.Database.AddTransaction(transactionModel)
	if err != nil {
		return nil, errors.New(AddTransactionError + err.Error())
	}

	return transaction, nil
}

func (t *Service) updateBalances(request request.CreateTransaction, origin *domain.Account, destination *domain.Account) error {
	switch request.Type {
	case common.TransactionTypeCredit:
		origin.Balance += request.Amount
		if destination != nil {
			destination.Balance -= request.Amount
		}
	case common.TransactionTypeDebit:
		origin.Balance -= request.Amount
		if destination != nil {
			destination.Balance += request.Amount
		}
	case common.TransactionTypeTransfer:
		if destination == nil {
			return errors.New("destination account is required for transfer transactions")
		}
		origin.Balance -= request.Amount
		destination.Balance += request.Amount
	}

	return t.Database.UpdateAccountBalances(origin, destination)
}

func (t *Service) GetTransactionsByAccount(accountID uint, offset int) ([]*domain.Transaction, error) {
	return t.Database.GetTransactionsByAccount(accountID, offset)
}

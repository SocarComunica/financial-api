package transaction

import (
	"errors"

	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
	"github.com/socarcomunica/financial-api/internal/domain"
)

const (
	AddTransactionError = "AddTransaction Service Error: "
)

type TransactionsDatabase interface {
	AddTransaction(model *domain.Transaction) (*domain.Transaction, error)
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
		DestinationID: request.DestinationID,
	}

	// TODO: Adds transaction to the database
	transaction, err := t.Database.AddTransaction(transactionModel)
	if err != nil {
		return nil, errors.New(AddTransactionError + err.Error())
	}
	return transaction, nil
}

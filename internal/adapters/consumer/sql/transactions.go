package sql

import (
	"errors"
	"github.com/socarcomunica/financial-api/internal/domain"
)

const (
	AddTransactionError         = "AddTransactionError DB Client: "
	ErrorUpdatingAccountBalance = "error updating account balance: "
)

func (c *client) AddTransaction(model *domain.Transaction) (*domain.Transaction, error) {
	result := c.DB.Create(model)

	if result.Error != nil {
		return nil, errors.New(AddTransactionError + result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return nil, errors.New(AddTransactionError + "no new transactions were created")
	}

	return model, nil
}

func (c *client) UpdateAccountBalance(account *domain.Account) error {
	result := c.DB.Save(account)
	if result.Error != nil {
		return errors.New(ErrorUpdatingAccountBalance + result.Error.Error())
	}

	return nil
}

func (c *client) GetTransactionsByAccount(accountID uint, offset int) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction

	//validate account exists
	var account domain.Account
	if result := c.DB.Where("id = ?", accountID).First(&account); result.Error != nil {
		return nil, result.Error
	}

	result := c.DB.Where("origin_id = ?", accountID).
		Preload("Origin").
		Preload("Destination").
		Preload("Tags").
		Limit(10).
		Offset(offset).
		Order("created_at desc").
		Find(&transactions)

	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}

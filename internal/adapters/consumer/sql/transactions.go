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

func (c *client) UpdateAccountBalances(origin *domain.Account, destination *domain.Account) error {
	// Start a transaction
	tx := c.DB.Begin()
	if tx.Error != nil {
		return errors.New(ErrorUpdatingAccountBalance + tx.Error.Error())
	}

	// Update origin account
	if err := tx.Save(origin).Error; err != nil {
		tx.Rollback()
		return errors.New(ErrorUpdatingAccountBalance + "origin: " + err.Error())
	}

	// Update destination account if provided
	if destination != nil {
		if err := tx.Save(destination).Error; err != nil {
			tx.Rollback()
			return errors.New(ErrorUpdatingAccountBalance + "destination: " + err.Error())
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New(ErrorUpdatingAccountBalance + "commit: " + err.Error())
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

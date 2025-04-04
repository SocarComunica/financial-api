package sql

import (
	"errors"

	"github.com/socarcomunica/financial-api/internal/domain"
)

const (
	CreateAccountError = "CreateAccountError DB Client: "
	DeleteAccountError = "DeleteAccountError DB Client: "
)

func (c *client) AddAccount(model *domain.Account) (*domain.Account, error) {
	result := c.DB.Create(model)

	if result.Error != nil {
		return nil, errors.New(CreateAccountError + result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return nil, errors.New(CreateAccountError + "no new accounts were created")
	}

	return model, nil
}

func (c *client) GetAccount(id uint) (*domain.Account, error) {
	var account domain.Account

	result := c.DB.Where("id = ?", id).First(&account)

	if result.Error != nil {
		return nil, result.Error
	}

	return &account, nil
}

func (c *client) GetAccountsByUser(userID uint) ([]*domain.Account, error) {
	var accounts []*domain.Account

	// validate user exists
	var user domain.User
	if result := c.DB.Where("id = ?", userID).First(&user); result.Error != nil {
		return nil, result.Error
	}

	result := c.DB.Where("user_id = ?", userID).
		Order("user_id asc").
		Find(&accounts)

	if result.Error != nil {
		return nil, result.Error
	}

	return accounts, nil
}

func (c *client) DeleteAccount(id uint, userID uint) error {
	// First verify that the account exists and belongs to the user
	var account domain.Account
	result := c.DB.Where("id = ? AND user_id = ?", id, userID).First(&account)
	if result.Error != nil {
		return errors.New(DeleteAccountError + result.Error.Error())
	}

	// Delete the account
	result = c.DB.Delete(&account)
	if result.Error != nil {
		return errors.New(DeleteAccountError + result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return errors.New(DeleteAccountError + "no account was deleted")
	}

	return nil
}

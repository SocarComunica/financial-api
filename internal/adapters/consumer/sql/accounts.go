package sql

import (
	"errors"
	"github.com/socarcomunica/financial-api/internal/domain"
)

const (
	CreateAccountError = "CreateAccountError DB Client: "
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

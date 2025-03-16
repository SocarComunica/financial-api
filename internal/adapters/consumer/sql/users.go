package sql

import (
	"errors"

	"github.com/socarcomunica/financial-api/internal/domain"
)

const (
	CreateUserError = "CreateUserError DB Client: "
)

func (c *client) AddUser(model *domain.User) (*domain.User, error) {
	result := c.DB.Create(model)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New(CreateUserError + "no new accounts were created")
	}

	return model, nil
}

func (c *client) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := c.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

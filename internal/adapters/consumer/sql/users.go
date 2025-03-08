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

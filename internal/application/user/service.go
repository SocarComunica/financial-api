package user

import (
	"errors"
	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
	"github.com/socarcomunica/financial-api/internal/domain"
)

const (
	AddUserError = "AddUser Service Error: "
)

type UsersDatabase interface {
	AddUser(model *domain.User) (*domain.User, error)
}

type Service struct {
	Database UsersDatabase
}

func NewUserService(database UsersDatabase) *Service {
	return &Service{
		Database: database,
	}
}

func (u *Service) AddUser(request request.CreateUser) (*domain.User, error) {
	userModel := &domain.User{
		Username: request.Username,
		Email:    request.Email,
		Accounts: []domain.Account{},
	}

	user, err := u.Database.AddUser(userModel)
	if err != nil {
		return nil, errors.New(AddUserError + err.Error())
	}

	return user, nil
}

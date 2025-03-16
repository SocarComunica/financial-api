package user

import (
	"errors"

	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
	"github.com/socarcomunica/financial-api/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

const (
	AddUserError = "AddUser Service Error: "
	LoginError   = "Login Service Error: "
)

type UsersDatabase interface {
	AddUser(model *domain.User) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
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
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New(AddUserError + "error hashing password: " + err.Error())
	}

	userModel := &domain.User{
		Username: request.Username,
		Email:    request.Email,
		Password: string(hashedPassword),
		Accounts: []domain.Account{},
	}

	user, err := u.Database.AddUser(userModel)
	if err != nil {
		return nil, errors.New(AddUserError + err.Error())
	}

	return user, nil
}

func (u *Service) Login(request request.Login) (*domain.User, error) {
	user, err := u.Database.GetUserByEmail(request.Email)
	if err != nil {
		return nil, errors.New(LoginError + err.Error())
	}

	if !u.VerifyPassword(user.Password, request.Password) {
		return nil, errors.New(LoginError + "invalid password")
	}

	return user, nil
}

// VerifyPassword compares a hashed password with a plain text password
func (u *Service) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

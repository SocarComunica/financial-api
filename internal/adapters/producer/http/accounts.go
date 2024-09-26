package http

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
	"github.com/socarcomunica/financial-api/internal/domain"
	"net/http"
)

const (
	CreateAccountError = "CreateAccountError Handler: "
)

type accountService interface {
	AddAccount(request request.CreateAccount) (*domain.Account, error)
}

type AccountsHandler struct {
	accountService accountService
}

func NewAccountsHandler(accountService accountService) *AccountsHandler {
	return &AccountsHandler{
		accountService: accountService,
	}
}

func (a *AccountsHandler) AddRoutes(router *echo.Router) {
	router.Add(echo.POST, "accounts", a.createAccount)
}

func (a *AccountsHandler) createAccount(c echo.Context) error {
	r := new(request.CreateAccount)
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(CreateAccountError+err.Error()).Error())
	}
	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(CreateAccountError+err.Error()).Error())
	}

	account, err := a.accountService.AddAccount(*r)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, &account)
}

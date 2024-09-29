package http

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
	"github.com/socarcomunica/financial-api/internal/domain"
	"net/http"
	"strconv"
)

const (
	CreateAccountError     = "CreateAccountError Handler: "
	GetAccountsByUserError = "GetAccountsByUserError Handler: "
)

type accountService interface {
	AddAccount(request request.CreateAccount) (*domain.Account, error)
	GetAccountsByUser(userID uint) ([]*domain.Account, error)
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
	router.Add(echo.GET, "accounts/:userID", a.getAccountsByUser)
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

func (a *AccountsHandler) getAccountsByUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(GetAccountsByUserError+err.Error()).Error())
	}

	accounts, err := a.accountService.GetAccountsByUser(uint(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"layer": GetAccountsByUserError,
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &accounts)
}

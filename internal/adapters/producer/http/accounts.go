package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/response"
	"github.com/socarcomunica/financial-api/internal/domain"
)

const (
	CreateAccountError     = "CreateAccountError Handler: "
	GetAccountsByUserError = "GetAccountsByUserError Handler: "
	DeleteAccountError     = "DeleteAccountError Handler: "
	UserIDHeader           = "X-User-ID"
)

type accountService interface {
	AddAccount(request request.CreateAccount) (*domain.Account, error)
	GetAccountsByUser(userID uint) ([]*domain.Account, error)
	DeleteAccount(id uint, userID uint) error
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
	router.Add(echo.DELETE, "accounts/:id", a.deleteAccount)
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

	accountResponse := response.MapAccountToResponse(account)
	return c.JSON(http.StatusCreated, accountResponse)
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

	// Initialize accountResponses as an empty slice, not nil
	accountResponses := make([]*response.Account, 0)
	for _, account := range accounts {
		accountResponses = append(accountResponses, response.MapAccountToResponse(account))
	}

	return c.JSON(http.StatusOK, accountResponses)
}

func (a *AccountsHandler) deleteAccount(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(DeleteAccountError+err.Error()).Error())
	}

	userIDStr := c.Request().Header.Get(UserIDHeader)
	if userIDStr == "" {
		return c.JSON(http.StatusBadRequest, errors.New(DeleteAccountError+"user ID header is required").Error())
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(DeleteAccountError+"invalid user ID format").Error())
	}

	err = a.accountService.DeleteAccount(uint(id), uint(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"layer": DeleteAccountError,
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
	"github.com/socarcomunica/financial-api/internal/domain"
)

const (
	CreateTransactionError = "CreateTransactionError Handler: "
)

type transactionService interface {
	AddTransaction(request request.CreateTransaction) (*domain.Transaction, error)
}

type TransactionsHandler struct {
	transactionService transactionService
}

func NewTransactionsHandler(transactionsService transactionService) *TransactionsHandler {
	return &TransactionsHandler{
		transactionService: transactionsService,
	}
}

func (t *TransactionsHandler) AddRoutes(router *echo.Router) {
	router.Add(echo.POST, "transactions", t.createTransaction)
}

func (t *TransactionsHandler) createTransaction(c echo.Context) error {
	r := new(request.CreateTransaction)
	if err := c.Bind(r); err != nil {
		return c.String(http.StatusBadRequest, errors.New(CreateTransactionError+err.Error()).Error())
	}
	if err := c.Validate(r); err != nil {
		return c.String(http.StatusBadRequest, errors.New(CreateTransactionError+err.Error()).Error())
	}

	transaction, err := t.transactionService.AddTransaction(*r)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, &transaction)
}

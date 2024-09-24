package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
	"github.com/socarcomunica/financial-api/internal/domain"
)

type transactionService interface {
	AddTransaction(request.CreateTransaction) (domain.Transaction, error)
}

type TransactionsHandler struct {
	transactionService transactionService
}

func NewTransactionsHandler() *TransactionsHandler {
	return &TransactionsHandler{}
}

func (t *TransactionsHandler) AddRoutes(router *echo.Router) {
	router.Add(echo.POST, "transactions", t.createTransaction)
}

func (t *TransactionsHandler) createTransaction(c echo.Context) error {
	request := new(request.CreateTransaction)
	if err := c.Bind(request); err != nil {
		return err
	}

	transaction, err := t.transactionService.AddTransaction(*request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, transaction)
}

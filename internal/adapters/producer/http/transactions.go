package http

import (
	"net/http"
	"strconv"

	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
	"github.com/socarcomunica/financial-api/internal/domain"

	"github.com/labstack/echo/v4"
)

const (
	CreateTransactionError = "CreateTransactionError Handler: "
)

type transactionService interface {
	AddTransaction(request request.CreateTransaction) (*domain.Transaction, error)
	GetTransactionsByAccount(accountID uint, offset int) ([]*domain.Transaction, error)
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
	router.Add(echo.GET, "transactions/:accountID", t.getTransactionsByAccount)
}

func (t *TransactionsHandler) createTransaction(c echo.Context) error {
	r := new(request.CreateTransaction)
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"layer": CreateTransactionError,
			"error": err.Error(),
		})
	}
	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"layer": CreateTransactionError,
			"error": err.Error(),
		})
	}

	transaction, err := t.transactionService.AddTransaction(*r)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"layer": CreateTransactionError,
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, &transaction)
}

func (t *TransactionsHandler) getTransactionsByAccount(c echo.Context) error {
	accountIDParam := c.Param("accountID")
	accountID, err := strconv.ParseUint(accountIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid accountID"})
	}

	offsetParam := c.QueryParam("offset")
	if offsetParam == "" {
		offsetParam = "0"
	}
	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid offset"})
	}

	transactions, err := t.transactionService.GetTransactionsByAccount(uint(accountID), offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error retrieving transactions"})
	}

	return c.JSON(http.StatusOK, transactions)
}

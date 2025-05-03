package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
	"github.com/socarcomunica/financial-api/internal/domain"
	"gorm.io/gorm"
)

const (
	CreateTransactionError = "CreateTransactionError Handler: "
)

type transactionService interface {
	AddTransaction(request request.CreateTransaction) (*domain.Transaction, error)
	GetTransactionsByAccount(accountID uint, offset int) ([]*domain.Transaction, error)
	GetLatestTransactionByUser(userID uint) (*domain.Transaction, error)
	GetTransactionsByUser(userID uint) ([]*domain.Transaction, error)
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
	router.Add(echo.GET, "transactions/latest", t.getLatestTransactionByUser)
	router.Add(echo.GET, "transactions/user/:userID", t.getAllTransactionsByUser)
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

// getLatestTransactionByUser handles requests to fetch the latest transaction for a user
func (t *TransactionsHandler) getLatestTransactionByUser(c echo.Context) error {
	userIDStr := c.QueryParam("user_id")
	if userIDStr == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "user_id query parameter is required"})
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user_id format"})
	}

	transaction, err := t.transactionService.GetLatestTransactionByUser(uint(userID))
	if err != nil {
		// Check if the error is 'record not found'
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "No transactions found for this user"})
		}
		// Handle other potential errors
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error retrieving latest transaction"})
	}

	// We might want a specific response DTO later, but for now return the domain model
	return c.JSON(http.StatusOK, transaction)
}

// getAllTransactionsByUser handles requests to fetch all transactions for a user
func (t *TransactionsHandler) getAllTransactionsByUser(c echo.Context) error {
	userIDStr := c.Param("userID") // Get userID from path parameter
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid userID format"})
	}

	transactions, err := t.transactionService.GetTransactionsByUser(uint(userID))
	if err != nil {
		// Since the service/db layer doesn't return a specific 'not found' for Find,
		// we just return a general server error here. An empty list is the expected 'not found' result.
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error retrieving transactions"})
	}

	// Service layer ensures transactions is an empty slice, not nil, if none are found.
	// The JSON marshaler will correctly output [] for an empty slice.
	return c.JSON(http.StatusOK, transactions)
}

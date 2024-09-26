package internal

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/socarcomunica/financial-api/common"
	"github.com/socarcomunica/financial-api/internal/adapters/consumer/sql"
	"github.com/socarcomunica/financial-api/internal/adapters/producer/http"
	"github.com/socarcomunica/financial-api/internal/application/transaction"
)

func Run() error {
	e := echo.New()

	common.InitConfig()
	config, err := common.GetConfig()
	if err != nil {
		return err
	}

	// add validator
	e.Validator = common.NewCustomValidator()

	// setup database
	var database sql.Client
	switch config.Environment {
	case "development":
		database = sql.NewLocalClient()

	default:
		return errors.New("undefined environment")
	}

	// add middlewares here

	// Config transactions
	transactionsService := transaction.NewTransactionService(database)
	transactionHandler := http.NewTransactionsHandler(transactionsService)

	router := e.Router()

	// add here endpoints
	handlers := []common.Handler{
		transactionHandler,
	}

	for _, handler := range handlers {
		handler.AddRoutes(router)
	}

	return e.Start(":8088")
}

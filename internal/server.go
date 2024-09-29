package internal

import (
	"github.com/socarcomunica/financial-api/common"
	"github.com/socarcomunica/financial-api/internal/adapters/consumer/sql"
	"github.com/socarcomunica/financial-api/internal/adapters/producer/http"
	"github.com/socarcomunica/financial-api/internal/application/account"
	"github.com/socarcomunica/financial-api/internal/application/transaction"
	"github.com/socarcomunica/financial-api/internal/application/user"

	"github.com/labstack/echo/v4"
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
	database := sql.NewClient(config.Database)

	// add middlewares here

	// Config users
	usersService := user.NewUserService(database)
	usersHandler := http.NewUsersHandler(usersService)

	// Config accounts
	accountsService := account.NewAccountsDatabase(database)
	accountsHandler := http.NewAccountsHandler(accountsService)

	// Config transactions
	transactionsService := transaction.NewTransactionService(database)
	transactionHandler := http.NewTransactionsHandler(transactionsService)

	router := e.Router()

	// add here endpoints
	handlers := []common.Handler{
		transactionHandler,
		accountsHandler,
		usersHandler,
	}

	for _, handler := range handlers {
		handler.AddRoutes(router)
	}

	return e.Start(":8088")
}

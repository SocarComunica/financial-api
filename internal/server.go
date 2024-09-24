package internal

import (
	"github.com/labstack/echo/v4"
	"github.com/socarcomunica/financial-api/common"
	"github.com/socarcomunica/financial-api/internal/adapters/producer/http"
)

func Run() error {
	e := echo.New()

	// add middlewares here

	router := e.Router()

	// add here endpoints
	handlers := []common.Handler{
		http.NewTransactionsHandler(),
	}

	for _, handler := range handlers {
		handler.AddRoutes(router)
	}

	return e.Start(":8088")
}

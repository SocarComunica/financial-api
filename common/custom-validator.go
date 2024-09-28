package common

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	v := validator.New()
	_ = v.RegisterValidation("positive", positive)
	_ = v.RegisterValidation("transaction_type", transactionType)
	return &CustomValidator{
		validator: v,
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func positive(fl validator.FieldLevel) bool {
	value := fl.Field().Float()
	return value >= 0
}

func transactionType(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return value == TransactionTypeDebit || value == TransactionTypeCredit || value == TransactionTypeTransfer
}

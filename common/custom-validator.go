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
	_ = v.RegisterValidation("password", password)
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

func password(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if len(value) < 8 {
		return false
	}
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, ch := range value {
		switch {
		case ch >= 'A' && ch <= 'Z':
			hasUpper = true
		case ch >= 'a' && ch <= 'z':
			hasLower = true
		case ch >= '0' && ch <= '9':
			hasNumber = true
		case (ch >= ' ' && ch <= '/') || (ch >= ':' && ch <= '@') ||
			(ch >= '[' && ch <= '`') || (ch >= '{' && ch <= '~'):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial
}

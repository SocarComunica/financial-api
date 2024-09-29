package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
	"github.com/socarcomunica/financial-api/internal/domain"
)

const (
	CreateUserError = "CreateUserError Handler: "
)

type userService interface {
	AddUser(request request.CreateUser) (*domain.User, error)
}

type UsersHandler struct {
	userService userService
}

func NewUsersHandler(userService userService) *UsersHandler {
	return &UsersHandler{
		userService: userService,
	}
}

func (u *UsersHandler) AddRoutes(router *echo.Router) {
	router.Add(echo.POST, "users", u.createUser)
}

func (u *UsersHandler) createUser(c echo.Context) error {
	r := new(request.CreateUser)
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"layer": CreateUserError,
			"error": err.Error(),
		})
	}
	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"layer": CreateUserError,
			"error": err.Error(),
		})
	}

	user, err := u.userService.AddUser(*r)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"layer": CreateUserError,
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, &user)
}

package http

import (
	"net/http"

	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/response"
	"github.com/socarcomunica/financial-api/internal/domain"

	"github.com/labstack/echo/v4"
)

const (
	CreateUserError = "CreateUserError Handler: "
	LoginError      = "LoginError Handler: "
)

type userService interface {
	AddUser(request request.CreateUser) (*domain.User, error)
	Login(request request.Login) (*domain.User, error)
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
	router.Add(echo.POST, "users/login", u.login)
}

func (u *UsersHandler) login(c echo.Context) error {
	r := new(request.Login)
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"layer": LoginError,
			"error": err.Error(),
		})
	}

	user, err := u.userService.Login(*r)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"layer": LoginError,
			"error": err.Error(),
		})
	}

	response := response.FromDomain(user)

	return c.JSON(http.StatusOK, response)
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

	response := response.FromDomain(user)

	return c.JSON(http.StatusCreated, response)
}

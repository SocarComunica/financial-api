package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
	"github.com/socarcomunica/financial-api/internal/domain"
)

const (
	CreateGoalError     = "CreateGoalError Handler: "
	GetGoalError        = "GetGoalError Handler: "
	GetGoalsError       = "GetGoalsError Handler: "
	GetLatestGoalsError = "GetLatestGoalsError Handler: "
)

type goalService interface {
	AddGoal(request request.CreateGoal) (*domain.Goal, error)
	GetGoal(id uint) (*domain.Goal, error)
	GetGoalsByUser(userID uint) ([]*domain.Goal, error)
	GetLatestThreeGoalsByUser(userID uint) ([]*domain.Goal, error)
}

type GoalsHandler struct {
	goalService goalService
}

func NewGoalsHandler(goalService goalService) *GoalsHandler {
	return &GoalsHandler{
		goalService: goalService,
	}
}

func (g *GoalsHandler) AddRoutes(router *echo.Router) {
	router.Add(echo.POST, "goals", g.createGoal)
	router.Add(echo.GET, "goals/:id", g.getGoal)
	router.Add(echo.GET, "users/:userID/goals", g.getGoalsByUser)
	router.Add(echo.GET, "users/:userID/goals/latest", g.getLatestThreeGoalsByUser)
}

func (g *GoalsHandler) createGoal(c echo.Context) error {
	r := new(request.CreateGoal)
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(CreateGoalError+err.Error()).Error())
	}
	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(CreateGoalError+err.Error()).Error())
	}

	goal, err := g.goalService.AddGoal(*r)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, goal)
}

func (g *GoalsHandler) getGoal(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(GetGoalError+err.Error()).Error())
	}

	goal, err := g.goalService.GetGoal(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"layer": GetGoalError,
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, goal)
}

func (g *GoalsHandler) getGoalsByUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(GetGoalsError+err.Error()).Error())
	}

	goals, err := g.goalService.GetGoalsByUser(uint(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"layer": GetGoalsError,
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, goals)
}

func (g *GoalsHandler) getLatestThreeGoalsByUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(GetLatestGoalsError+err.Error()).Error())
	}

	goals, err := g.goalService.GetLatestThreeGoalsByUser(uint(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"layer": GetLatestGoalsError,
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, goals)
}

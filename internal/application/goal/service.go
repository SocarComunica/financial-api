package goal

import (
	"errors"

	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
	"github.com/socarcomunica/financial-api/internal/domain"
)

const (
	AddGoalError = "AddGoal Service Error: "
)

type GoalsDatabase interface {
	AddGoal(model *domain.Goal) (*domain.Goal, error)
	GetGoal(id uint) (*domain.Goal, error)
	UpdateGoal(goal *domain.Goal) error
	GetGoalsByUser(userID uint) ([]*domain.Goal, error)
}

type Service struct {
	Database GoalsDatabase
}

func NewGoalService(database GoalsDatabase) *Service {
	return &Service{
		Database: database,
	}
}

func (g *Service) AddGoal(request request.CreateGoal) (*domain.Goal, error) {
	if request.Initial > request.Target {
		return nil, errors.New(AddGoalError + "initial amount cannot be greater than target amount")
	}

	goalModel := &domain.Goal{
		UserID:      request.UserID,
		Title:       request.Title,
		Description: request.Description,
		Target:      request.Target,
		Initial:     request.Initial,
		Current:     request.Initial,
	}

	goal, err := g.Database.AddGoal(goalModel)
	if err != nil {
		return nil, errors.New(AddGoalError + err.Error())
	}

	return goal, nil
}

func (g *Service) GetGoal(id uint) (*domain.Goal, error) {
	goal, err := g.Database.GetGoal(id)
	if err != nil {
		return nil, errors.New(AddGoalError + err.Error())
	}
	return goal, nil
}

func (g *Service) GetGoalsByUser(userID uint) ([]*domain.Goal, error) {
	return g.Database.GetGoalsByUser(userID)
}

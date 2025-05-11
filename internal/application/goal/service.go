package goal

import (
	"errors"

	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
	"github.com/socarcomunica/financial-api/internal/domain"
)

const (
	AddGoalError               = "AddGoal Service Error: "
	GetLatestGoalsServiceError = "GetLatestGoals Service Error: "
	UpdateGoalServiceError     = "UpdateGoal Service Error: "
)

type GoalsDatabase interface {
	AddGoal(model *domain.Goal) (*domain.Goal, error)
	GetGoal(id uint) (*domain.Goal, error)
	UpdateGoal(goal *domain.Goal) error
	GetGoalsByUser(userID uint) ([]*domain.Goal, error)
	GetLatestThreeGoalsByUser(userID uint) ([]*domain.Goal, error)
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

func (g *Service) GetLatestThreeGoalsByUser(userID uint) ([]*domain.Goal, error) {
	goals, err := g.Database.GetLatestThreeGoalsByUser(userID)
	if err != nil {
		return nil, errors.New(GetLatestGoalsServiceError + err.Error())
	}
	return goals, nil
}

func (g *Service) UpdateGoal(id uint, req request.UpdateGoal) (*domain.Goal, error) {
	goal, err := g.Database.GetGoal(id)
	if err != nil {
		return nil, errors.New(UpdateGoalServiceError + "goal not found: " + err.Error())
	}

	// Update fields if provided in the request
	if req.Title != nil {
		goal.Title = *req.Title
	}
	if req.Description != nil {
		goal.Description = *req.Description
	}
	if req.Target != nil {
		if *req.Target < goal.Initial {
			return nil, errors.New(UpdateGoalServiceError + "target amount cannot be less than initial amount")
		}
		goal.Target = *req.Target
	}
	if req.Current != nil {
		if *req.Current < 0 {
			return nil, errors.New(UpdateGoalServiceError + "current amount cannot be negative")
		}
		if *req.Current > goal.Target {
			// Potentially allow current to exceed target, or cap it. For now, allowing.
			// Or, return errors.New(UpdateGoalServiceError + "current amount cannot exceed target amount")
		}
		goal.Current = *req.Current
	}

	// Validate that current amount is not greater than target amount after all updates
	if goal.Current > goal.Target {
		//This validation might be too restrictive depending on product requirements
		//Consider if a goal can be "overachieved"
		//return nil, errors.New(UpdateGoalServiceError + "current amount cannot be greater than target amount after update")
	}

	if err := g.Database.UpdateGoal(goal); err != nil {
		return nil, errors.New(UpdateGoalServiceError + "failed to update goal: " + err.Error())
	}

	return goal, nil
}

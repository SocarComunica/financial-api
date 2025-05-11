package sql

import (
	"errors"

	"github.com/socarcomunica/financial-api/internal/domain"
)

const (
	CreateGoalError     = "CreateGoalError DB Client: "
	GetGoalError        = "GetGoalError DB Client: "
	UpdateGoalError     = "UpdateGoalError DB Client: "
	GetLatestGoalsError = "GetLatestGoalsError DB Client: "
)

func (c *client) AddGoal(model *domain.Goal) (*domain.Goal, error) {
	// Validate user exists
	var user domain.User
	if result := c.DB.Where("id = ?", model.UserID).First(&user); result.Error != nil {
		return nil, errors.New(CreateGoalError + "user not found")
	}

	result := c.DB.Create(model)

	if result.Error != nil {
		return nil, errors.New(CreateGoalError + result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return nil, errors.New(CreateGoalError + "no new goals were created")
	}

	return model, nil
}

func (c *client) GetGoal(id uint) (*domain.Goal, error) {
	var goal domain.Goal

	result := c.DB.Where("id = ?", id).First(&goal)

	if result.Error != nil {
		return nil, errors.New(GetGoalError + result.Error.Error())
	}

	return &goal, nil
}

func (c *client) UpdateGoal(goal *domain.Goal) error {
	result := c.DB.Save(goal)

	if result.Error != nil {
		return errors.New(UpdateGoalError + result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return errors.New(UpdateGoalError + "no goals were updated")
	}

	return nil
}

func (c *client) GetGoalsByUser(userID uint) ([]*domain.Goal, error) {
	// Validate user exists
	var user domain.User
	if result := c.DB.Where("id = ?", userID).First(&user); result.Error != nil {
		return nil, errors.New(GetGoalError + "user not found")
	}

	var goals []*domain.Goal
	result := c.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&goals)

	if result.Error != nil {
		return nil, result.Error
	}

	return goals, nil
}

func (c *client) GetLatestThreeGoalsByUser(userID uint) ([]*domain.Goal, error) {
	// Validate user exists
	var user domain.User
	if result := c.DB.Where("id = ?", userID).First(&user); result.Error != nil {
		return nil, errors.New(GetLatestGoalsError + "user not found")
	}

	var goals []*domain.Goal
	result := c.DB.Where("user_id = ?", userID).Order("updated_at desc").Limit(3).Find(&goals)

	if result.Error != nil {
		return nil, errors.New(GetLatestGoalsError + result.Error.Error())
	}

	// Return empty slice if no goals found, not an error
	if len(goals) == 0 {
		return []*domain.Goal{}, nil
	}

	return goals, nil
}

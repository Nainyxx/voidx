package userservice

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"example.com/m/models"
	"example.com/m/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserProfileByID(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user_id cannot be empty")
	}

	profile, err := s.repo.GetUserProfileByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

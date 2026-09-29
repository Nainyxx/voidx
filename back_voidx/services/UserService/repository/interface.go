package repository

import (
	"context"
	"uuid"

	"example.com/m/models"
)

type UserRepository interface {
	GetUserProfileByID(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error)
	CreateUserProfile(ctx context.Context, profile *models.UserProfile) error
	UpdateUserProfile(ctx context.Context, profile *models.UserProfile) error
	DeleteUserProfile(ctx context.Context, userID uuid.UUID) error
}

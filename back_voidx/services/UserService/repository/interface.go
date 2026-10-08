package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"example.com/m/models"
)

var (
	ErrNotFound      = errors.New("profile not found")
	ErrAlreadyExists = errors.New("profile already exists")
)

type UserRepository interface {
	GetUserProfileByID(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error)
	CreateUserProfile(ctx context.Context, profile *models.UserProfile) error
	UpdateUserProfile(ctx context.Context, profile *models.UserProfile) error
	DeleteUserProfile(ctx context.Context, userID uuid.UUID) error
}

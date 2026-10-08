package userservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"example.com/m/models"
	"example.com/m/repository"
)

var (
	ErrInvalidUserID = errors.New("user_id cannot be empty")
	ErrNilUpdate     = errors.New("update cannot be nil")
)

type UserService struct {
	repo repository.UserRepository
}

type UpdateUserProfileInput struct {
	UserID         uuid.UUID
	Username       *string
	Name           *string
	Surname        *string
	Phone          *string
	Description    *string
	AvatarImageURL *string
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func formatErr(err error) error {
	return fmt.Errorf("UserService: %w", err)
}

// GET
func (s *UserService) GetUserProfileByID(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	if userID == uuid.Nil {
		return nil, formatErr(ErrInvalidUserID)
	}

	profile, err := s.repo.GetUserProfileByID(ctx, userID)
	if err != nil {
		return nil, formatErr(err)
	}

	return profile, nil
}

// POST
func (s *UserService) CreateUserProfile(ctx context.Context, userID uuid.UUID, username, name, surname, phone string) (*models.UserProfile, error) {
	if userID == uuid.Nil {
		return nil, formatErr(ErrInvalidUserID)
	}

	userProfile, err := models.CreateUserProfile(userID, username, name, surname, phone)
	if err != nil {
		return nil, formatErr(err)
	}
	err = s.repo.CreateUserProfile(ctx, userProfile)
	if err != nil {
		return nil, formatErr(err)
	}
	return userProfile, nil
}

// PATCH
func (s *UserService) UpdateUserProfile(ctx context.Context, update *UpdateUserProfileInput) (*models.UserProfile, error) {
	if update == nil {
		return nil, formatErr(ErrNilUpdate)
	}
	if update.UserID == uuid.Nil {
		return nil, formatErr(ErrInvalidUserID)
	}

	userProfile, err := s.repo.GetUserProfileByID(ctx, update.UserID)
	if err != nil {
		return nil, formatErr(err)
	}

	changed := false
	if update.Username != nil && *update.Username != userProfile.Username {
		err = userProfile.ChangeUsername(*update.Username)
		if err != nil {
			return nil, formatErr(err)
		}
		changed = true
	}
	if update.Name != nil && *update.Name != userProfile.Name {
		err = userProfile.ChangeName(*update.Name)
		if err != nil {
			return nil, formatErr(err)
		}
		changed = true
	}
	if update.Surname != nil && *update.Surname != userProfile.Surname {
		err = userProfile.ChangeSurname(*update.Surname)
		if err != nil {
			return nil, formatErr(err)
		}
		changed = true
	}
	if update.Phone != nil && *update.Phone != userProfile.Phone {
		err = userProfile.ChangePhone(*update.Phone)
		if err != nil {
			return nil, formatErr(err)
		}
		changed = true
	}
	if update.Description != nil && *update.Description != userProfile.Description {
		err = userProfile.ChangeDescription(*update.Description)
		if err != nil {
			return nil, formatErr(err)
		}
		changed = true
	}
	if update.AvatarImageURL != nil && *update.AvatarImageURL != userProfile.AvatarImageURL {
		err = userProfile.ChangeAvatarImageURL(*update.AvatarImageURL)
		if err != nil {
			return nil, formatErr(err)
		}
		changed = true
	}
	if !changed {
		return userProfile, nil
	}

	err = s.repo.UpdateUserProfile(ctx, userProfile)
	if err != nil {
		return nil, formatErr(err)
	}
	return userProfile, nil
}

// DELETE
func (s *UserService) DeleteUserProfile(ctx context.Context, userID uuid.UUID) error {
	if userID == uuid.Nil {
		return formatErr(ErrInvalidUserID)
	}

	err := s.repo.DeleteUserProfile(ctx, userID)
	if err != nil {
		return formatErr(err)
	}
	return nil
}

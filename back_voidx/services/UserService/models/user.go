package models

import (
	"errors"
	"fmt"
	"time"
	"uuid"

	"example.com/m/utils"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
	Surname  string    `json:"surname"`

	Login        string `json:"login"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	PasswordHash string `json:"-"`

	Description    string `json:"description"`
	AvatarImageURL string `json:"avatarImageUrl"`

	IsOnline  bool      `json:"isOnline"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// RegisterUser creates new user with validation
func RegisterUser(name, username, surname, login, email, phone, password string) (*User, error) {
	// Validate all fields
	if err := utils.IsEmailValid(email); err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}
	if err := utils.IsLoginValid(login); err != nil {
		return nil, fmt.Errorf("invalid login: %w", err)
	}
	if err := utils.IsNameValid(name); err != nil {
		return nil, fmt.Errorf("invalid name: %w", err)
	}
	if err := utils.IsNameValid(username); err != nil {
		return nil, fmt.Errorf("invalid username: %w", err)
	}
	if err := utils.IsPhoneValid(phone); err != nil {
		return nil, fmt.Errorf("invalid phone: %w", err)
	}\

	// Hash password
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	return &User{
		ID:           uuid.New(),
		Name:         name,
		Username:     username,
		Surname:      surname,
		Login:        login,
		Email:        email,
		Phone:        phone,
		PasswordHash: passwordHash,
		Description:  "",
		IsOnline:     false,
		UpdatedAt:    time.Now(),
		CreatedAt:    time.Now(),
	}, nil
}

// ChangeName updates user name
func (u *User) ChangeName(newName string) error {
	if err := utils.IsNameValid(newName); err != nil {
		return fmt.Errorf("invalidd name: %w", err)
	}
	u.Name = newName
	u.UpdatedAt = time.Now()
	return nil
}

// ChangeSurname updates user surname
func (u *User) ChangeSurname(newSurname string) error {
	if err := utils.IsNameValid(newSurname); err != nil {
		return fmt.Errorf("invalid surname: %w", err)
	}
	u.Surname = newSurname
	u.UpdatedAt = time.Now()
	return nil
}

// ChangeDescription updates user description
func (u *User) ChangeDescription(newDescription string) error {
	if len(newDescription) > 500 {
		return errors.New("description too long (max 500 chars)")
	}
	u.Description = newDescription
	u.UpdatedAt = time.Now()
	return nil
}

// ChangeAvatar updates user avatar
func (u *User) ChangeAvatar(newAvatarURL string) error {
	if newAvatarURL == "" {
		return errors.New("avatar URL cannot be empty")
	}
	u.AvatarImageURL = newAvatarURL
	u.UpdatedAt = time.Now()
	return nil
}

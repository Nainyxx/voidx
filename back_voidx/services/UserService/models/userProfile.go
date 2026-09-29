package models

import (
	"errors"
	"time"
	"uuid"

	"example.com/m/utils"
)

type UserProfile struct {
	UserID   uuid.UUID
	Username string
	Name     string
	Surname  string
	Phone    string

	Description    string
	AvatarImageURL string

	CreatedAt time.Time
	UpdatedAt time.Time

	// IsOnline bool
	// LastSeen time.Time
}

func CreateUserProfile(userID uuid.UUID, username, name, surname, phone string) (*UserProfile, error) {
	var creationErr error
	creationErr = utils.IsLoginValid(username)
	if creationErr != nil {
		return nil, creationErr
	}

	creationErr = utils.IsNameValid(name)
	if creationErr != nil {
		return nil, creationErr
	}

	if len(surname) != 0 {
		creationErr = utils.IsNameValid(surname)
		if creationErr != nil {
			return nil, creationErr
		}
	}

	creationErr = utils.IsPhoneValid(phone)
	if creationErr != nil {
		return nil, creationErr
	}

	return &UserProfile{
		UserID:         userID,
		Username:       username,
		Name:           name,
		Surname:        surname,
		Phone:          phone,
		Description:    "",
		AvatarImageURL: "",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

func (p *UserProfile) ChangeUsername(newUsername string) error {
	updateErr := utils.IsLoginValid(newUsername)
	if updateErr != nil {
		return updateErr
	}
	if newUsername == p.Username {
		return errors.New("new username must be different")
	}

	p.Username = newUsername
	p.UpdatedAt = time.Now()

	return nil
}

func (p *UserProfile) ChangeName(newName string) error {
	updateErr := utils.IsNameValid(newName)
	if updateErr != nil {
		return updateErr
	}
	if newName == p.Name {
		return errors.New("new name must be different")
	}

	p.Name = newName
	p.UpdatedAt = time.Now()

	return nil
}

func (p *UserProfile) ChangeSurname(newSurname string) error {
	updateErr := utils.IsSurnameValid(newSurname)
	if updateErr != nil {
		return updateErr
	}
	if newSurname == p.Surname {
		return errors.New("new surname must be different")
	}

	p.Surname = newSurname
	p.UpdatedAt = time.Now()

	return nil
}

func (p *UserProfile) ChangeDescription(newDesc string) error {
	updateErr := utils.IsDescriptionValid(newDesc)
	if updateErr != nil {
		return updateErr
	}
	if newDesc == p.Description {
		return errors.New("new description must be different")
	}

	p.Description = newDesc
	p.UpdatedAt = time.Now()

	return nil
}

func (p *UserProfile) ChangeAvatarImageURL(imageURL string) error {
	updateErr := utils.IsImageURLValid(imageURL)
	if updateErr != nil {
		return updateErr
	}
	if imageURL == p.AvatarImageURL {
		return errors.New("new avatar must be different")
	}

	p.AvatarImageURL = imageURL
	p.UpdatedAt = time.Now()

	return nil
}

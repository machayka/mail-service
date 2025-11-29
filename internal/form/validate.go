package form

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrEmptyForm    = errors.New("empty form")
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidUUID  = errors.New("invalid uuid")
)

func ValidateFormData(d *FormData) error {
	if d.Email == "" || d.Message == "" {
		return ErrEmptyForm
	}

	return nil
}

func ValidateID(id string) error {
	_, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidUUID
	}
	return nil
}

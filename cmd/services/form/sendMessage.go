package form

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/machayka/mail-service/cmd/models"
)

func (s *FormService) SendMessage(id string, d *models.FormData) error {
	if d.Email == "" || d.Message == "" {
		return errors.New("empty form")
	}

	uuidValue, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid uuid")
	}

	form, err := s.repo.GetByID(uuidValue) // form from database
	if err != nil {
		return err
	}

	if form == nil {
		return errors.New("not found")
	}

	fmt.Println(form.Email)
	return nil
}

// Package form jest odpowiedzialny za logikę biznesową obsługi formularza
package form

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SendMessage(id string, d *FormData) (form *Form, error error) {
	if d.Email == "" || d.Message == "" {
		return nil, errors.New("empty form")
	}

	uuidValue, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid uuid")
	}

	form, err = s.repo.GetByID(uuidValue) // form from database
	if err != nil {
		return nil, err
	}

	// TODO: w prod do usunięcia linijka niżej
	fmt.Println(form.Email)
	return form, nil
}

// Package form jest odpowiedzialny za logikę biznesową obsługi formularza
package form

import "github.com/machayka/mail-service/internal/mail"

type Service struct {
	repo       *Repository
	mailSender *mail.MailService
}

func NewService(repo *Repository, mailSender *mail.MailService) *Service {
	return &Service{repo: repo, mailSender: mailSender}
}

func (s *Service) SendMessage(id string, d *FormData) (f *Form, error error) {
	if err := ValidateFormData(d); err != nil {
		return nil, err
	}
	if err := ValidateID(id); err != nil {
		return nil, err
	}

	f, err := s.repo.GetByID(id) // form from database
	if err != nil {
		return nil, err
	}

	err = s.mailSender.SendMessageFromContactForm(f.Email, d.Email, d.Message)
	if err != nil {
		return f, err
	}

	return f, nil
}

func (s *Service) RegisterNewForm(form *Form) error {
	err := s.repo.CreateNewForm(form)
	// TODO: tu robimy obsługę stripe'a?
	return err
}

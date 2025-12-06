// Package form jest odpowiedzialny za logikę biznesową obsługi formularza
package form

import (
	"database/sql"

	"github.com/machayka/mail-service/internal/mail"
	"github.com/machayka/mail-service/internal/payments"
)

type Service struct {
	repo          *Repository
	mailSender    *mail.Service
	paymentClient *payments.Payment
}

func NewService(repo *Repository, mailSender *mail.Service, paymentClient *payments.Payment) *Service {
	return &Service{repo: repo, mailSender: mailSender, paymentClient: paymentClient}
}

func (s *Service) SendMessage(id string, d *FormData) error {
	if err := ValidateFormData(d); err != nil {
		return err
	}
	if err := ValidateID(id); err != nil {
		return err
	}

	f, err := s.repo.GetByID(id) // form from database
	if err != nil {
		return err
	}

	err = s.mailSender.SendMessageFromContactForm(f.Email, d.Email, d.Message)
	return err
}

func (s *Service) CreateCheckout(f *Form) (string, error) {
	customerID, err := s.repo.GetStripeCustomerID(f.Email)
	if err == sql.ErrNoRows {
		// Użytkownik nie ma customerID w stripe -> musimy go stworzyć
		customerID, err = s.paymentClient.CreateCustomer(f.Email)
		if err != nil {
			return "", err
		}
		err = s.repo.CreateCustomer(f.Email, customerID)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	checkoutURL, err := s.paymentClient.CreatePayment(customerID, f.ID.String(), f.Email)
	if err != nil {
		return "", err
	}

	return checkoutURL, err
}

func (s *Service) HandleCheckoutCompleted(formID, email, customerID, subscriptionID string) error {
	if formID == "" || email == "" || customerID == "" {
		return ErrMissingMetadata
	}
	if subscriptionID == "" {
		return ErrMissingSubscriptionID
	}

	return s.repo.CreateNewForm(formID, email, subscriptionID)
}

func (s *Service) HandleSubscriptionDeleted(subscriptionID string) error {
	if subscriptionID == "" {
		return ErrMissingSubscriptionID
	}
	return s.repo.DeleteForm(subscriptionID)
}

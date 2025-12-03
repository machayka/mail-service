// Package mail jest odpowiedzialne za obsługę maili
package mail

import (
	"errors"
	"fmt"

	"github.com/machayka/mail-service/config"
	"github.com/wneessen/go-mail"
)

var (
	ErrFailedToSetFromAddr      = errors.New("failed to set FROM address")
	ErrFailedToSetToAddr        = errors.New("failed to set TO address")
	ErrFailedToCreateMailClient = errors.New("failed to create mail delivery client")
	ErrFailedToDeliverMail      = errors.New("failed to deliver mail")
)

type Service struct {
	smtpHost string
	smptUser string
	smtpPass string
}

func NewService(cfg *config.Config) *Service {
	return &Service{
		smtpHost: "smtp.gmail.com",
		smptUser: cfg.SMTP.User,
		smtpPass: cfg.SMTP.Pass,
	}
}

func (s *Service) SendMessageFromContactForm(to, from, msg string) error {
	message := mail.NewMsg()
	if err := message.From(s.smptUser); err != nil {
		return ErrFailedToSetFromAddr
	}
	if err := message.To(to); err != nil {
		return ErrFailedToSetToAddr
	}
	message.Subject("Nowa wiadomość z formularza kontaktowego")
	message.SetBodyString(mail.TypeTextPlain, fmt.Sprintf("%v: %v", from, msg))

	client, err := mail.NewClient(s.smtpHost,
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover), mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(s.smptUser), mail.WithPassword(s.smtpPass),
	)
	if err != nil {
		return ErrFailedToCreateMailClient
	}
	if err := client.DialAndSend(message); err != nil {
		return ErrFailedToDeliverMail
	}
	return nil
}

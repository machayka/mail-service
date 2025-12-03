// Package form jest odpowiedzialny za obsługę formularza
package form

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/machayka/mail-service/config"
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/webhook"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) FormSubmit(c *fiber.Ctx) error {
	var formData FormData
	if err := c.BodyParser(&formData); err != nil {
		return err
	}

	id := c.Params("id")
	err := h.service.SendMessage(id, &formData)
	if err != nil {
		if err == ErrNotFound {
			// Przekierowanie do formularza rejestracji nowego formularza
			return c.Redirect(fmt.Sprintf("/add/%v", id))
		}
		return err
	}

	// TODO: Tutaj zwrócimy success.html z elegancką informacją
	return c.SendString("Wiadomość przesłana")
}

func (h *Handler) NewForm(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.Render("new-form", fiber.Map{
		"ID": id,
	})
}

func (h *Handler) AddForm(c *fiber.Ctx) error {
	var form Form
	if err := c.BodyParser(&form); err != nil {
		return err
	}
	checkoutURL, err := h.service.CreateCheckout(&form)
	if err != nil {
		return err
	}
	return c.Redirect(checkoutURL)
}

func (h *Handler) PaymentSuccess(c *fiber.Ctx) error {
	// TODO: 	Tu bym chciał zwrócić customer manage portal link
	return c.Render("add-form-success", fiber.Map{"stripeUrl": "https://link.com"})
}

func (h *Handler) HandleWebhook(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload := c.Body()

		const MaxBodyBytes = int64(65536)
		if len(payload) > int(MaxBodyBytes) {
			return c.Status(fiber.StatusRequestEntityTooLarge).SendString("Too large")
		}

		signatureHeader := c.Get("Stripe-Signature")

		endpointSecret := cfg.Stripe.WebhookSecret

		event, err := webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid signature")
		}
		switch event.Type {
		case "checkout.session.completed!":
			log.Println("checkout.session.completed!")
			// TODO: Tutaj zmieniamy is_paid na true
			var subscription stripe.Subscription
			err := json.Unmarshal(event.Data.Raw, &subscription)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
				return err
			}
			log.Println(subscription.ID)
		default:
			log.Printf("Unhandled event type: %s", event.Type)
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

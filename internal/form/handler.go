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

	return c.Render("static/submitted", fiber.Map{})
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
	return c.Render("static/success", fiber.Map{})
}

func (h *Handler) CustomerPortal(redirect string) fiber.Handler {
	// inny link jest dla produkcji a inny dla testu
	return func(c *fiber.Ctx) error {
		return c.Redirect(redirect)
	}
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
		case "checkout.session.completed":
			var session stripe.CheckoutSession
			err := json.Unmarshal(event.Data.Raw, &session)
			if err != nil {
				log.Println("Error parsing checkout session:", err)
				return err
			}

			formID := session.Metadata["form_id"]
			if formID == "" {
				log.Println("Missing form_id in metadata")
				return c.Status(fiber.StatusBadRequest).SendString("Missing form_id")
			}

			subscriptionID := session.Subscription.ID
			if subscriptionID == "" {
				return c.Status(fiber.StatusBadRequest).SendString("Missing subscription ID")
			}

			err = h.service.repo.UpdateSubscriptionID(formID, subscriptionID)
			if err != nil {
				log.Println("Error updating subscription ID:", err)
				return err
			}

		case "customer.subscription.created":
			subscriptionID, err := getSubscriptionIDFromStripe(event)
			if err != nil {
				log.Println(err)
				return err
			}
			err = h.service.repo.ChangePaymentStatus(subscriptionID, true)
			if err != nil {
				log.Println(err)
				return err
			}
		case "customer.subscription.deleted":
			subscriptionID, err := getSubscriptionIDFromStripe(event)
			if err != nil {
				return err
			}
			err = h.service.repo.DeleteForm(subscriptionID)
			if err != nil {
				return err
			}
		default:
			log.Printf("Unhandled event type: %s", event.Type)
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

func getSubscriptionIDFromStripe(event stripe.Event) (string, error) {
	var subscription stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &subscription)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
		return "", err
	}
	return subscription.ID, nil
}

func (h *Handler) Regulamin(c *fiber.Ctx) error {
	return c.Render("static/regulamin", fiber.Map{})
}

func (h *Handler) PolitykaPrywatnosci(c *fiber.Ctx) error {
	return c.Render("static/polityka-prywatnosci", fiber.Map{})
}

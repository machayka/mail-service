// Package form jest odpowiedzialny za obsługę formularza
package form

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/machayka/mail-service/config"
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/webhook"
)

const (
	stripeEventCheckoutCompleted   = "checkout.session.completed"
	stripeEventSubscriptionDeleted = "customer.subscription.deleted"
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

	// Sprawdź czy formularz już istnieje
	existingForm, err := h.service.CheckIfFormExists(id)
	if err != nil {
		return err
	}

	if existingForm != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Ten formularz jest już aktywny")
	}

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
		if err == ErrFormAlreadyExists {
			return c.Status(fiber.StatusConflict).SendString("Ten formularz jest już aktywny")
		}
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
			return c.Status(fiber.StatusRequestEntityTooLarge).SendString("Payload too large")
		}

		event, err := webhook.ConstructEvent(payload, c.Get("Stripe-Signature"), cfg.Stripe.WebhookSecret)
		if err != nil {
			log.Println("Webhook signature verification failed:", err)
			return c.Status(fiber.StatusBadRequest).SendString("Invalid signature")
		}

		switch event.Type {
		case stripeEventCheckoutCompleted:
			if err := h.handleCheckoutCompleted(c, event); err != nil {
				return err
			}
		case stripeEventSubscriptionDeleted:
			if err := h.handleSubscriptionDeleted(c, event); err != nil {
				return err
			}
		default:
			log.Printf("Unhandled event type: %s", event.Type)
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

func (h *Handler) handleCheckoutCompleted(c *fiber.Ctx, event stripe.Event) error {
	var session stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		log.Println("Error parsing checkout session:", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid session data")
	}

	formID := session.Metadata["form_id"]
	email := session.Metadata["email"]
	customerID := session.Customer.ID
	subscriptionID := session.Subscription.ID

	err := h.service.HandleCheckoutCompleted(formID, email, customerID, subscriptionID)
	if err != nil {
		log.Println("Error handling checkout completed:", err)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return nil
}

func (h *Handler) handleSubscriptionDeleted(c *fiber.Ctx, event stripe.Event) error {
	var subscription stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
		log.Println("Error parsing subscription:", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid subscription data")
	}

	err := h.service.HandleSubscriptionDeleted(subscription.ID)
	if err != nil {
		log.Println("Error handling subscription deleted:", err)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return nil
}

func (h *Handler) Regulamin(c *fiber.Ctx) error {
	return c.Render("static/regulamin", fiber.Map{})
}

func (h *Handler) PolitykaPrywatnosci(c *fiber.Ctx) error {
	return c.Render("static/polityka-prywatnosci", fiber.Map{})
}

func (h *Handler) Index(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

func (h *Handler) NotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).Render("404", fiber.Map{})
}

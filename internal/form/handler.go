// Package form jest odpowiedzialny za obsługę formularza
package form

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) FormHandler(c *fiber.Ctx) error {
	var formData FormData
	if err := c.BodyParser(&formData); err != nil {
		return err
	}

	id := c.Params("id")
	err := h.service.SendMessage(id, &formData)
	if err != nil {
		if err == ErrFormNotFound {
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
	err := h.service.RegisterNewForm(&form)
	if err != nil {
		return err
	}

	// TODO: Dopiero jak wszystko pójdzie dobrze z płatnością to zwróć sukces
	return c.Render("add-form-success", fiber.Map{"stripeUrl": "https://link.com"}) // Link do zarządzanie subskrypcją
}

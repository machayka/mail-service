// Package form jest odpowiedzialny za obsługę formularza
package form

import (
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
		return fiber.ErrBadRequest
	}

	id := c.Params("id")

	form, err := h.service.SendMessage(id, &formData)
	if err != nil {
		if form == nil {
			// TODO: zrobić przekierowanie na rejestracje formularza
			return fiber.ErrLocked
		}
		return err
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Form submitted successfully",
	})
}

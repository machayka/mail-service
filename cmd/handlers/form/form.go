// Package form jest odpowiedzialny za obsługę formularza
package form

import (
	"github.com/gofiber/fiber/v2"
	"github.com/machayka/mail-service/cmd/models"
)

func (h *FormHandler) FormHandler(c *fiber.Ctx) error {
	var formData models.FormData
	if err := c.BodyParser(&formData); err != nil {
		return fiber.ErrBadRequest
	}

	id := c.Params("id")

	err := h.service.SendMessage(id, &formData)
	if err != nil {
		if err.Error() == "not found" {
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

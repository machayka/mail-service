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
		return err
	}
	id := c.Params("id")

	form, err := h.service.SendMessage(id, &formData)
	if err != nil {
		switch err {
		case ErrInvalidUUID:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   ErrInvalidUUID.Error(),
				"message": "Id must be valid UUID",
			})
		case ErrEmptyForm:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   ErrEmptyForm.Error(),
				"message": "Email and message are required",
			})
		case ErrFormNotFound:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   ErrFormNotFound.Error(),
				"message": "Form is not registered",
			})

		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Form submitted successfully",
		"data": fiber.Map{
			"formId": form.ID,
			"email":  form.Email,
			"formData": fiber.Map{
				"email":   formData.Email,
				"message": formData.Message,
			},
		},
	})
}

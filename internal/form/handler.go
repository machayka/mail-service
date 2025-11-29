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
			return c.Redirect(fmt.Sprintf("/add/%v", id))
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
	}

	// TODO: Tutaj zwrócimy success.html z elegancką informacją zamiast tego co jest teraz
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
	// TODO: to na switchach, bo może być dużo rodzajów błedów?
	if err != nil {
		return err
	}

	// TODO: Dopiero jak wsystko pójdzie dobrze z płatnością to zwróć sukces
	return c.Render("add-form-success", fiber.Map{"stripeUrl": "https://link.com"}) // Link do zarządzanie subskrypcją
}

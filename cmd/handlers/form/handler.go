package form

import FormService "github.com/machayka/mail-service/cmd/services/form"

type FormHandler struct {
	service *FormService.FormService
}

func NewFormHandler(service *FormService.FormService) *FormHandler {
	return &FormHandler{service: service}
}

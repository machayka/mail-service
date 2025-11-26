// Package form jest odpowiedzialny za logikę biznesową obsługi formularza
package form

import FormRepository "github.com/machayka/mail-service/cmd/repositories/form"

type FormService struct {
	repo *FormRepository.FormRepository
}

func NewFormService(repo *FormRepository.FormRepository) *FormService {
	return &FormService{repo: repo}
}

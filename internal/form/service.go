// Package form jest odpowiedzialny za logikę biznesową obsługi formularza
package form

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SendMessage(id string, d *FormData) (form *Form, error error) {
	if err := ValidateFormData(d); err != nil {
		return nil, err
	}
	if err := ValidateID(id); err != nil {
		return nil, err
	}

	form, err := s.repo.GetByID(id) // form from database
	if err != nil {
		return nil, err
	}

	return form, nil
}

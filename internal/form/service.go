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

	// TODO: stworzyć pakiet do wysyłania maili i wysłac rzeczywiście tego maila + dodać obsługę błędów

	return form, nil
}

func (s *Service) RegisterNewForm(form *Form) error {
	err := s.repo.CreateNewForm(form)
	// TODO: tu robimy obsługę stripe'a?
	return err
}

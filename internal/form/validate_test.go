package form

import (
	"testing"
)

func TestID(t *testing.T) {
	tests := []struct {
		id            string
		expectedError error
	}{
		{
			id:            "xyz-323",
			expectedError: ErrInvalidUUID,
		}, {
			id:            "c15563b2-3bcb-4117-9a17-0770684aac98",
			expectedError: nil,
		}, {
			id:            "",
			expectedError: ErrInvalidUUID,
		},
	}
	for _, tt := range tests {
		err := ValidateID(tt.id)
		if err != tt.expectedError {
			t.Errorf("ValidateID() error: %v, wanted error: %v", err, tt.expectedError)
		}

	}
}

func TestValidateFormData(t *testing.T) {
	tests := []struct {
		d             *FormData
		expectedError error
	}{
		{
			d: &FormData{
				Email:   "test@example.com",
				Message: "Test message",
			},
			expectedError: nil,
		}, {
			d: &FormData{
				Email:   "",
				Message: "Test",
			},
			expectedError: ErrEmptyForm,
		}, {
			d: &FormData{
				Email:   "",
				Message: "",
			},
			expectedError: ErrEmptyForm,
		},
	}

	for _, tt := range tests {
		err := ValidateFormData(tt.d)

		if err != tt.expectedError {
			t.Errorf("ValidateFormData() error: %v, wanted error: %v", err, tt.expectedError)
		}
	}
}

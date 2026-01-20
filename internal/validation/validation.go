package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	FailedField string `json:"field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
	Message     string `json:"message,omitempty"`
}

type XValidator struct {
	validator *validator.Validate
}

var Validator = &XValidator{
	validator: validator.New(),
}

func (v *XValidator) Validate(data interface{}) []ErrorResponse {
	var validationErrors []ErrorResponse

	err := v.validator.Struct(data)
	if err == nil {
		return nil
	}

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			validationErrors = append(validationErrors, ErrorResponse{
				FailedField: e.Field(),
				Tag:         e.Tag(),
				Value:       fmt.Sprintf("%v", e.Value()),
				Message:     "Validation failed for " + e.Field(),
			})
		}
	}

	return validationErrors
}

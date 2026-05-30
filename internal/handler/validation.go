package handler

import (
	"encoding/json"
	"errors"

	"github.com/go-playground/validator/v10"
)

func handleValidationError(err error) string {

	var unmarshalErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalErr) {
		return unmarshalErr.Field + " must be " + unmarshalErr.Type.String()
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			switch fieldErr.Tag() {
			case "required":
				return fieldErr.Field() + " is required"
			case "max":
				return fieldErr.Field() + " exceeds maximum length"
			case "min":
				return fieldErr.Field() + " does not meet minimum length"
			case "email":
				return fieldErr.Field() + " must be a valid email"
			case "oneof":
				return fieldErr.Field() + " must be one of " + fieldErr.Param()
			default:
				return fieldErr.Field() + " is invalid"
			}
		}
	}
	return err.Error()
}

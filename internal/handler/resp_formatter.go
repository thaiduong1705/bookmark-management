package handler

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// FieldErrorResponse is the 400 response body returned when request binding or validation fails.
type FieldErrorResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func parseBindingError(err error) FieldErrorResponse {
	var ve validator.ValidationErrors
	if ok := errors.As(err, &ve); !ok {
		return FieldErrorResponse{
			Message: "Invalid request",
		}
	}

	var errStr string
	for _, err := range ve {
		errStr = errStr + err.Field() + " is invalid (" + err.Tag() + ") "
	}

	return FieldErrorResponse{
		Message: "Invalid request",
		Detail:  strings.TrimSpace(errStr),
	}
}

package response

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

const invalidRequestMessage = "Invalid request"

// ErrorResponse is the common error body; Detail is optional.
type ErrorResponse struct {
	Message
	Detail string `json:"detail,omitempty"`
}

// InternalServerError is the body returned for unexpected server errors.
var InternalServerError = NewErrorResponse("Internal server error", "")

// NewErrorResponse builds an ErrorResponse; detail may be empty.
func NewErrorResponse(message, detail string) ErrorResponse {
	return ErrorResponse{Message: NewMessage(message), Detail: detail}
}

// ParseBindingError converts a request binding/validation error into a 400 ErrorResponse.
func ParseBindingError(err error) ErrorResponse {
	var ve validator.ValidationErrors
	if ok := errors.As(err, &ve); !ok {
		return NewErrorResponse(invalidRequestMessage, "")
	}

	var errStr string
	for _, err := range ve {
		errStr = errStr + err.Field() + " is invalid (" + err.Tag() + ") "
	}

	return NewErrorResponse(invalidRequestMessage, strings.TrimSpace(errStr))
}

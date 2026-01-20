package domain

import "errors"

// Domain errors that can be returned by core services.
var (
	// Link errors
	ErrLinkNotFound      = errors.New("link not found")
	ErrLinkExpired       = errors.New("link has expired")
	ErrLinkDeleted       = errors.New("link has been deleted")
	ErrSlugTaken         = errors.New("slug is already taken")
	ErrSlugInvalid       = errors.New("slug contains invalid characters")
	ErrSlugTooShort      = errors.New("slug is too short")
	ErrSlugTooLong       = errors.New("slug is too long")
	ErrURLInvalid        = errors.New("URL is invalid")
	ErrURLUnreachable    = errors.New("URL is unreachable")
	ErrPasswordRequired  = errors.New("password is required for this link")
	ErrPasswordIncorrect = errors.New("password is incorrect")

	// Auth errors
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInvalidAPIKey     = errors.New("invalid API key")
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token has expired")

	// Validation errors
	ErrValidation        = errors.New("validation error")
	ErrMissingField      = errors.New("required field is missing")

	// Storage errors
	ErrDatabase          = errors.New("database error")
	ErrNotImplemented    = errors.New("not implemented")
)

// ValidationError wraps a validation error with field details.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a validation error for a specific field.
func NewValidationError(field, message string) ValidationError {
	return ValidationError{Field: field, Message: message}
}

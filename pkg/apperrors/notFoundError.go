package apperrors

import "fmt"

type NotFoundError struct {
	ResourceType string
	Identifier   string
}

func NewNotFoundError(resourceType, identifier string) *NotFoundError {
	return &NotFoundError{
		ResourceType: resourceType,
		Identifier:   identifier,
	}
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with identifier '%s' not found", e.ResourceType, e.Identifier)
}

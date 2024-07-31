package apiErrors

import (
	"encoding/json"
	"fmt"
)

type ErrorCode string

type APIError struct {
	Code    string      `json:"code"`
	Status  int         `json:"-"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func New(code string, message string, status int) *APIError {
	return &APIError{
		Code:    code,
		Status:  status,
		Message: message,
	}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *APIError) JSON() ([]byte, error) {
	return json.Marshal(e)
}

// Custom error types (example)
type NotFoundError struct {
	Resource string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("resource not found: %s", e.Resource)
}

func IsNotFound(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

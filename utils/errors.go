package utils

import "fmt"

type HTTPError struct {
	Status  int
	Message string
}

func (h HTTPError) Error() string {
	return fmt.Sprintf("%d: %s", h.Status, h.Message)
}

func NewHTTPError(status int, msg string) error {
	return HTTPError{status, msg}
}

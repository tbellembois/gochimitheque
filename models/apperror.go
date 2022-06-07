package models

import (
	"fmt"
	"net/http"
)

// AppError is the error type returned by the custom handlers.
type AppError struct {
	OriginalError error
	Message       string
	Code          int
}

func (e *AppError) Error() string {
	originalError := ""

	if e.OriginalError != nil {
		originalError = e.OriginalError.Error()
	}

	return fmt.Sprintf("%s %s", originalError, e.Message)
}

// AppHandlerFunc is an HandlerFunc returning an AppError.
type AppHandlerFunc func(http.ResponseWriter, *http.Request) *AppError

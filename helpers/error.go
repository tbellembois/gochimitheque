package helpers

// AppError is the error type returned by the custom handlers
type AppError struct {
	Error   error
	Message string
	Code    int
}

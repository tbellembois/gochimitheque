package handlers

// Env is used to pass variables throughout the application.
type Env struct {
	// "http://localhost:8081" by default
	AppURL string
	// BuildID is a compile time variable
	BuildID string
}

func NewEnv() Env {
	return Env{}
}

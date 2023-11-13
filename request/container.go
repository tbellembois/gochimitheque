package request

import (
	"net/http"
)

// ChimithequeContextKey is the Go request context
// used in each request.
type ChimithequeContextKey string

// Container is a struct passed to the view.
type Container struct {
	PersonEmail    string `json:"PersonEmail"`
	PersonLanguage string `json:"PersonLanguage"`
	PersonID       int    `json:"PersonID"`
	AppURL         string `json:"AppURL"`
	AppPath        string `json:"AppPath"`
	BuildID        string `json:"BuildID"`
	DisableCache   bool   `json:"DisableCache"`
}

// ContainerFromRequestContext returns a ViewContainer from the request context
// initialized in the AuthenticateMiddleware and AuthorizeMiddleware middlewares.
func ContainerFromRequestContext(r *http.Request) Container {
	// getting the request context
	var (
		container Container
	)

	ctx := r.Context()
	ctxcontainer := ctx.Value(ChimithequeContextKey("container"))

	if ctxcontainer != nil {
		container = ctxcontainer.(Container)
	}

	return container
}

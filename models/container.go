package models

import (
	"net/http"
)

// ViewContainer is a struct passed to the view
type ViewContainer struct {
	PersonEmail    string `json:"PersonEmail"`
	PersonLanguage string `json:"PersonLanguage"`
	PersonID       int    `json:"PersonID"`
	ProxyPath      string `json:"ProxyPath"`
	BuildID        string `json:"BuildID"`
	DisableCache   bool   `json:"DisableCache"`
}

// ContainerFromRequestContext returns a ViewContainer from the request context
// initialized in the AuthenticateMiddleware and AuthorizeMiddleware middlewares
func ContainerFromRequestContext(r *http.Request) ViewContainer {
	// getting the request context
	var (
		container ViewContainer
	)
	ctx := r.Context()
	ctxcontainer := ctx.Value(ChimithequeContextKey("container"))
	if ctxcontainer != nil {
		container = ctxcontainer.(ViewContainer)
	}
	return container
}

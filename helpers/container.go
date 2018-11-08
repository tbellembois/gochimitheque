package helpers

import (
	"net/http"
	"net/url"
)

// ViewContainer is a struct passed to the view
type ViewContainer struct {
	PersonEmail string
	PersonID    int
	URLValues   url.Values
}

// helpers.ContainerFromRequestContext returns a ViewContainer from the request context
// initialized in the AuthenticateMiddleware and AuthorizeMiddleware middlewares
func ContainerFromRequestContext(r *http.Request) ViewContainer {
	// getting the request context
	var (
		container ViewContainer
	)
	ctx := r.Context()
	ctxcontainer := ctx.Value("container")
	if ctxcontainer != nil {
		container = ctxcontainer.(ViewContainer)
	}
	//container.URLValues = r.URL.Query()
	return container
}

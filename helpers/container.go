package helpers

import (
	"net/http"
	"net/url"

	"github.com/tbellembois/gochimitheque/global"
)

// ViewContainer is a struct passed to the view
type ViewContainer struct {
	PersonEmail    string
	PersonLanguage string
	PersonID       int
	URLValues      url.Values
	ProxyPath      string
	BuildID        string
}

// helpers.ContainerFromRequestContext returns a ViewContainer from the request context
// initialized in the AuthenticateMiddleware and AuthorizeMiddleware middlewares
func ContainerFromRequestContext(r *http.Request) ViewContainer {
	// getting the request context
	var (
		container ViewContainer
	)
	ctx := r.Context()
	ctxcontainer := ctx.Value(global.ChimithequeContextKey("container"))
	if ctxcontainer != nil {
		container = ctxcontainer.(ViewContainer)
	}
	return container
}

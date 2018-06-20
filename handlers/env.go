package handlers

import (
	"github.com/tbellembois/gochimitheque/models"
	"html/template"
	"net/http"
)

// Env is a structure used to pass objects throughout the application.
type Env struct {
	DB        models.Datastore              // application DB connection
	Templates map[string]*template.Template // application templates
}

// containerFromRequestContext returns a ViewContainer from the request context
// initialized in the AuthenticateMiddleware and AuthorizeMiddleware middlewares
func containerFromRequestContext(r *http.Request) models.ViewContainer {
	// getting the request context
	var (
		container models.ViewContainer
	)
	ctx := r.Context()
	ctxcontainer := ctx.Value("container")
	if ctxcontainer != nil {
		container = ctxcontainer.(models.ViewContainer)
	}
	return container
}

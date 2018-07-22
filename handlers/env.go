package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/tbellembois/gochimitheque/models"
)

// Env is a structure used to pass objects throughout the application.
type Env struct {
	DB        models.Datastore              // application DB connection
	Templates map[string]*template.Template // application templates
}

// FakeHandler returns true
func (env *Env) FakeHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("true")
	return nil
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

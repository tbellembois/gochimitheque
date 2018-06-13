package handlers

import (
	"net/http"

	"github.com/tbellembois/gochimitheque/models"
)

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

// HomeHandler serve the main page
func (env *Env) HomeHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := containerFromRequestContext(r)

	if e := env.Templates["home"].Execute(w, c); e != nil {
		return &models.AppError{
			Error:   e,
			Code:    http.StatusInternalServerError,
			Message: "error executing template home",
		}
	}

	return nil
}

package handlers

import (
	"github.com/tbellembois/gochimitheque/helpers"
	"net/http"
)

/*
	views handlers
*/

// HomeHandler serves the main page
func (env *Env) HomeHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	if e := env.Templates["home"].Execute(w, c); e != nil {
		return &helpers.AppError{
			Error:   e,
			Code:    http.StatusInternalServerError,
			Message: "error executing template home",
		}
	}

	return nil
}

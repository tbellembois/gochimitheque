package handlers

import (
	"net/http"

	"github.com/tbellembois/gochimitheque/helpers"
)

func (env *Env) VTestHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	if e := env.Templates["test"].ExecuteTemplate(w, "BASE", c); e != nil {
		return &helpers.AppError{
			Error:   e,
			Code:    http.StatusInternalServerError,
			Message: "error executing template base",
		}
	}
	return nil
}

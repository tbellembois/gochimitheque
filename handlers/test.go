package handlers

import (
	"github.com/tbellembois/gochimitheque/helpers"
	"net/http"
)

func (env *Env) VTestHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	if e := env.Templates["test"].Execute(w, c); e != nil {
		return &helpers.AppError{
			Error:   e,
			Code:    http.StatusInternalServerError,
			Message: "error executing template base",
		}
	}
	return nil
}

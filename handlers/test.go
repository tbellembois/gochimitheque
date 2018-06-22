package handlers

import (
	"github.com/tbellembois/gochimitheque/models"
	"net/http"
)

func (env *Env) VTestHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := containerFromRequestContext(r)

	if e := env.Templates["test"].Execute(w, c); e != nil {
		return &models.AppError{
			Error:   e,
			Code:    http.StatusInternalServerError,
			Message: "error executing template base",
		}
	}
	return nil
}

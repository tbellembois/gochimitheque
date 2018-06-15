package handlers

import (
	"github.com/tbellembois/gochimitheque/models"
	"net/http"
)

func (env *Env) VTestHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	if e := env.Templates["jadetest"].ExecuteTemplate(w, "menu", nil); e != nil {
		return &models.AppError{
			Error:   e,
			Code:    http.StatusInternalServerError,
			Message: "error executing template menu",
		}
	}
	return nil
}

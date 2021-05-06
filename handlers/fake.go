package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tbellembois/gochimitheque/models"
)

// FakeHandler returns true.
func (env *Env) FakeHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode("true"); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil

}

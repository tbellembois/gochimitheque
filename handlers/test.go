package handlers

import (
	"errors"
	"net/http"

	"github.com/tbellembois/gochimitheque/models"
)

// VTestHandler is a test handler of course
func (env *Env) VTestHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	e := errors.New("test error")

	return &models.AppError{
		Code:    http.StatusInternalServerError,
		Message: "error running the test",
		Error:   e,
	}

}

package handlers

import (
	"errors"
	"net/http"

	"github.com/tbellembois/gochimitheque/helpers"
)

// VTestHandler is a test handler of course
func (env *Env) VTestHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	e := errors.New("test error")

	return &helpers.AppError{
		Code:    http.StatusInternalServerError,
		Message: "error running the test",
		Error:   e,
	}

}

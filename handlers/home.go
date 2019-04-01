package handlers

import (
	"net/http"

	"github.com/tbellembois/gochimitheque/jade"

	"github.com/tbellembois/gochimitheque/helpers"
)

/*
	views handlers
*/

// HomeHandler serves the main page
func (env *Env) HomeHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	jade.Home(c, w)

	return nil
}

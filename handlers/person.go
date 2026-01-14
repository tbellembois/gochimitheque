package handlers

import (
	"net/http"

	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/request"
	"github.com/tbellembois/gochimitheque/static/jade"
)

/*
	views handlers
*/

// VCreatePersonHandler handles the person creation page.
func (env *Env) VCreatePersonHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Personcreate(c, w)

	return nil
}

// VGetPeopleHandler handles the people list page.
func (env *Env) VGetPeopleHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Personindex(c, w)

	return nil
}

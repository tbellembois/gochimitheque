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

// VGetEntitiesHandler handles the entity list page.
func (env *Env) VGetEntitiesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Entityindex(c, w)

	return nil
}

// VCreateEntityHandler handles the entity creation page.
func (env *Env) VCreateEntityHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Entitycreate(c, w)

	return nil
}

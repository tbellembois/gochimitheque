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

// VGetStoreLocationsHandler handles the store location list page.
func (env *Env) VGetStoreLocationsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Storelocationindex(c, w)

	return nil
}

// VCreateStoreLocationHandler handles the store location creation page.
func (env *Env) VCreateStoreLocationHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Storelocationcreate(c, w)

	return nil
}

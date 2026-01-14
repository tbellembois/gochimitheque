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

// VGetStoragesHandler handles the store location list page.
func (env *Env) VGetStoragesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Storageindex(c, w)

	return nil
}

// VCreateStorageHandler handles the storage creation page.
func (env *Env) VCreateStorageHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Storagecreate(c, w)

	return nil
}

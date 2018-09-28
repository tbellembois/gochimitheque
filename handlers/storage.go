package handlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/helpers"
	"github.com/tbellembois/gochimitheque/models"
	"net/http"
)

/*
	views handlers
*/

// VGetStoragesHandler handles the store location list page
func (env *Env) VGetStoragesHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	if e := env.Templates["storageindex"].Execute(w, c); e != nil {
		return &helpers.AppError{
			Error:   e,
			Code:    http.StatusInternalServerError,
			Message: "error executing template base",
		}
	}
	return nil
}

/*
	REST handlers
*/

// GetStoragesHandler returns a json list of the storages matching the search criteria
func (env *Env) GetStoragesHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetStoragesHandler")

	var (
		err  error
		aerr *helpers.AppError
		dsps helpers.DbselectparamStorage
	)

	// init db request parameters
	if dsps, aerr = helpers.NewdbselectparamStorage(r); err != nil {
		return aerr
	}

	storages, count, err := env.DB.GetStorages(dsps)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the storages",
		}
	}

	type resp struct {
		Rows  []models.Storage `json:"rows"`
		Total int              `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp{Rows: storages, Total: count})
	return nil
}

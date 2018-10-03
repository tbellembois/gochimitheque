package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/helpers"
	"github.com/tbellembois/gochimitheque/models"
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

// VCreateStorageHandler handles the storage creation page
func (env *Env) VCreateStorageHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	if e := env.Templates["storagecreate"].Execute(w, c); e != nil {
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

// GetStorageHandler returns a json of the entity with the requested id
func (env *Env) GetStorageHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	storage, err := env.DB.GetStorage(id)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the storage",
		}
	}
	log.WithFields(log.Fields{"storage": storage}).Debug("GetStorageHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(storage)
	return nil
}

// UpdateStorageHandler updates the storage from the request form
func (env *Env) UpdateStorageHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
		s   models.Storage
	)
	if err := r.ParseForm(); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form parsing error",
			Code:    http.StatusBadRequest}
	}

	// retrieving the logged user id from request context
	c := helpers.ContainerFromRequestContext(r)

	var decoder = schema.NewDecoder()
	if err := decoder.Decode(&s, r.PostForm); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}
	log.WithFields(log.Fields{"s": s}).Debug("UpdateStorageHandler")

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}
	updateds, _ := env.DB.GetStorage(id)
	updateds.StorageComment = s.StorageComment
	updateds.StoreLocation = s.StoreLocation
	updateds.PersonID = c.PersonID
	log.WithFields(log.Fields{"updateds": updateds}).Debug("UpdateStorageHandler")

	if err := env.DB.UpdateStorage(updateds); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "update storage error",
			Code:    http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updateds)
	return nil
}

// DeleteStorageHandler deletes the storage with the requested id
func (env *Env) DeleteStorageHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	env.DB.DeleteStorage(id)
	return nil
}

// CreateStorageHandler creates the storage from the request form
func (env *Env) CreateStorageHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("CreateStorageHandler")
	var (
		s models.Storage
	)
	if err := r.ParseForm(); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form parsing error",
			Code:    http.StatusBadRequest}
	}

	// retrieving the logged user id from request context
	c := helpers.ContainerFromRequestContext(r)

	var decoder = schema.NewDecoder()
	if err := decoder.Decode(&s, r.PostForm); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}
	s.StorageCreationDate = time.Now()
	s.PersonID = c.PersonID
	log.WithFields(log.Fields{"s": s}).Debug("CreateStorageHandler")

	if err, _ := env.DB.CreateStorage(s); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "create storage error",
			Code:    http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)
	return nil
}

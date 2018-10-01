package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/helpers"
	"github.com/tbellembois/gochimitheque/models"
)

/*
	views handlers
*/

// VGetStoreLocationsHandler handles the store location list page
func (env *Env) VGetStoreLocationsHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	if e := env.Templates["storelocationindex"].Execute(w, c); e != nil {
		return &helpers.AppError{
			Error:   e,
			Code:    http.StatusInternalServerError,
			Message: "error executing template base",
		}
	}
	return nil
}

// VCreateStoreLocationHandler handles the store location creation page
func (env *Env) VCreateStoreLocationHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	if e := env.Templates["storelocationcreate"].Execute(w, c); e != nil {
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

// GetStoreLocationsHandler returns a json list of the store locations matching the search criteria
func (env *Env) GetStoreLocationsHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetStoreLocationsHandler")

	var (
		err   error
		aerr  *helpers.AppError
		dspsl helpers.DbselectparamStoreLocation
	)

	// init db request parameters
	if dspsl, aerr = helpers.NewdbselectparamStoreLocation(r); err != nil {
		return aerr
	}

	storelocations, count, err := env.DB.GetStoreLocations(dspsl)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the store locations",
		}
	}

	type resp struct {
		Rows  []models.StoreLocation `json:"rows"`
		Total int                    `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp{Rows: storelocations, Total: count})
	return nil
}

// GetStoreLocationHandler returns a json of the store location with the requested id
func (env *Env) GetStoreLocationHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
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

	storelocation, err := env.DB.GetStoreLocation(id)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the store location",
		}
	}
	log.WithFields(log.Fields{"storelocation": storelocation}).Debug("GetStoreLocationHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(storelocation)
	return nil
}

// CreateStoreLocationHandler creates the store location from the request form
func (env *Env) CreateStoreLocationHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("CreateStoreLocationHandler")
	var (
		sl models.StoreLocation
	)
	if err := r.ParseForm(); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form parsing error",
			Code:    http.StatusBadRequest}
	}

	var decoder = schema.NewDecoder()
	if err := decoder.Decode(&sl, r.PostForm); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}
	log.WithFields(log.Fields{"sl": sl}).Debug("CreateStoreLocationHandler")

	if err, _ := env.DB.CreateStoreLocation(sl); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "create store location error",
			Code:    http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sl)
	return nil
}

// UpdateStoreLocationHandler updates the store location from the request form
func (env *Env) UpdateStoreLocationHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
		sl  models.StoreLocation
	)
	if err := r.ParseForm(); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form parsing error",
			Code:    http.StatusBadRequest}
	}
	var decoder = schema.NewDecoder()
	if err := decoder.Decode(&sl, r.PostForm); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}
	log.WithFields(log.Fields{"sl": sl}).Debug("UpdateStoreLocationHandler")

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	updatedsl, _ := env.DB.GetStoreLocation(id)
	updatedsl.StoreLocationName = sl.StoreLocationName
	updatedsl.Entity = sl.Entity
	log.WithFields(log.Fields{"updatedsl": updatedsl}).Debug("UpdateStoreLocationHandler")

	if err := env.DB.UpdateStoreLocation(updatedsl); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "update store location error",
			Code:    http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedsl)
	return nil
}

// DeleteStoreLocationHandler deletes the store location with the requested id
func (env *Env) DeleteStoreLocationHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
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

	env.DB.DeleteStoreLocation(id)
	return nil
}

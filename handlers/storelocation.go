package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/global"
	"github.com/tbellembois/gochimitheque/helpers"
	"github.com/tbellembois/gochimitheque/jade"
	"github.com/tbellembois/gochimitheque/models"
)

/*
	views handlers
*/

// VGetStoreLocationsHandler handles the store location list page
func (env *Env) VGetStoreLocationsHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	jade.Storelocationindex(c, w)

	return nil
}

// VCreateStoreLocationHandler handles the store location creation page
func (env *Env) VCreateStoreLocationHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	jade.Storelocationcreate(c, w)

	return nil
}

/*
	REST handlers
*/

// GetStoreLocationsHandler returns a json list of the store locations matching the search criteria
func (env *Env) GetStoreLocationsHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	global.Log.Debug("GetStoreLocationsHandler")

	var (
		err   error
		aerr  *helpers.AppError
		dspsl helpers.DbselectparamStoreLocation
	)

	// init db request parameters
	if dspsl, aerr = helpers.NewdbselectparamStoreLocation(r, nil); err != nil {
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
	global.Log.WithFields(logrus.Fields{"storelocation": storelocation}).Debug("GetStoreLocationHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(storelocation)
	return nil
}

// CreateStoreLocationHandler creates the store location from the request form
func (env *Env) CreateStoreLocationHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	global.Log.Debug("CreateStoreLocationHandler")
	var (
		sl  models.StoreLocation
		err error
		id  int
	)
	if err := r.ParseForm(); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form parsing error",
			Code:    http.StatusBadRequest}
	}

	if err := global.Decoder.Decode(&sl, r.PostForm); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}
	// processing storelocation not processed by Decode
	if r.FormValue("storelocation.storelocation.storelocation_id") != "" {
		var slid int
		slname := r.FormValue("storelocation.storelocation.storelocation_name")
		if slid, err = strconv.Atoi(r.FormValue("storelocation.storelocation.storelocation_id")); err != nil {
			return &helpers.AppError{
				Error:   err,
				Message: "slid atoi conversion",
				Code:    http.StatusInternalServerError}
		}
		sl.StoreLocation = &models.StoreLocation{
			StoreLocationID:   sql.NullInt64{Valid: true, Int64: int64(slid)},
			StoreLocationName: sql.NullString{Valid: true, String: slname},
		}
	}
	global.Log.WithFields(logrus.Fields{"sl": sl}).Debug("CreateStoreLocationHandler")

	if id, err = env.DB.CreateStoreLocation(sl); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "create store location error",
			Code:    http.StatusInternalServerError}
	}
	sl.StoreLocationID = sql.NullInt64{Valid: true, Int64: int64(id)}

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
	if err := global.Decoder.Decode(&sl, r.PostForm); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}
	// processing storelocation not processed by Decode
	if r.FormValue("storelocation.storelocation.storelocation_id") != "" {
		var slid int
		slname := r.FormValue("storelocation.storelocation.storelocation_name")
		if slid, err = strconv.Atoi(r.FormValue("storelocation.storelocation.storelocation_id")); err != nil {
			return &helpers.AppError{
				Error:   err,
				Message: "slid atoi conversion",
				Code:    http.StatusInternalServerError}
		}
		sl.StoreLocation = &models.StoreLocation{
			StoreLocationID:   sql.NullInt64{Valid: true, Int64: int64(slid)},
			StoreLocationName: sql.NullString{Valid: true, String: slname},
		}
	}
	global.Log.WithFields(logrus.Fields{"sl": sl}).Debug("UpdateStoreLocationHandler")

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	updatedsl, err := env.DB.GetStoreLocation(id)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "get store location error",
			Code:    http.StatusInternalServerError}
	}
	updatedsl.StoreLocationName = sl.StoreLocationName
	updatedsl.StoreLocationColor = sl.StoreLocationColor
	updatedsl.StoreLocationCanStore = sl.StoreLocationCanStore
	updatedsl.StoreLocation = sl.StoreLocation
	updatedsl.Entity = sl.Entity
	global.Log.WithFields(logrus.Fields{"updatedsl": updatedsl}).Debug("UpdateStoreLocationHandler")

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

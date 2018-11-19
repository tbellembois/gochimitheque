package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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

	if e := env.Templates["storageindex"].ExecuteTemplate(w, "BASE", c); e != nil {
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

	if e := env.Templates["storagecreate"].ExecuteTemplate(w, "BASE", c); e != nil {
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

// GetStoragesUnitsHandler returns a json list of the units matching the search criteria
func (env *Env) GetStoragesUnitsHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetStoragesUnitsHandler")

	var (
		err  error
		aerr *helpers.AppError
		dsp  helpers.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = helpers.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	units, count, err := env.DB.GetStoragesUnits(dsp)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the units",
		}
	}

	type resp struct {
		Rows  []models.Unit `json:"rows"`
		Total int           `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp{Rows: units, Total: count})
	return nil
}

// GetStoragesSuppliersHandler returns a json list of the suppliers matching the search criteria
func (env *Env) GetStoragesSuppliersHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetStoragesSuppliersHandler")

	var (
		err  error
		aerr *helpers.AppError
		dsp  helpers.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = helpers.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	suppliers, count, err := env.DB.GetStoragesSuppliers(dsp)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the suppliers",
		}
	}

	type resp struct {
		Rows  []models.Supplier `json:"rows"`
		Total int               `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp{Rows: suppliers, Total: count})
	return nil
}

// GetStoragesHandler returns a json list of the storages matching the search criteria
func (env *Env) GetStoragesHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetStoragesHandler")

	var (
		err  error
		aerr *helpers.AppError
		dsps helpers.DbselectparamStorage
	)

	// init db request parameters
	if dsps, aerr = helpers.NewdbselectparamStorage(r, nil); err != nil {
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

	if err := Decoder.Decode(&s, r.PostForm); err != nil {
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
	updateds.StorageModificationDate = time.Now()
	updateds.StorageBarecode = s.StorageBarecode
	updateds.StorageQuantity = s.StorageQuantity
	updateds.Supplier = s.Supplier
	updateds.Unit = s.Unit
	updateds.StorageComment = s.StorageComment
	updateds.StoreLocation = s.StoreLocation
	updateds.PersonID = c.PersonID
	updateds.StorageEntryDate = s.StorageEntryDate
	updateds.StorageExitDate = s.StorageExitDate
	updateds.StorageOpeningDate = s.StorageOpeningDate
	updateds.StorageExpirationDate = s.StorageExpirationDate
	updateds.StorageReference = s.StorageReference
	updateds.StorageBatchNumber = s.StorageBatchNumber
	updateds.StorageToDestroy = s.StorageToDestroy
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
		s   models.Storage
		err error
		id  int
	)
	if err := r.ParseForm(); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form parsing error",
			Code:    http.StatusBadRequest}
	}

	// retrieving the logged user id from request context
	c := helpers.ContainerFromRequestContext(r)

	if err := Decoder.Decode(&s, r.PostForm); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}

	// retrieving the full store location
	// we need its entity id to compute the barecode
	if s.StoreLocation, err = env.DB.GetStoreLocation(int(s.StoreLocationID.Int64)); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "error retrieving the storage store location",
			Code:    http.StatusInternalServerError}
	}

	s.StorageCreationDate = time.Now()
	s.StorageModificationDate = time.Now()
	s.PersonID = c.PersonID
	log.WithFields(log.Fields{"s": s}).Debug("CreateStorageHandler")

	for i := 1; i <= s.StorageNbItem; i++ {
		if err, id = env.DB.CreateStorage(s); err != nil {
			return &helpers.AppError{
				Error:   err,
				Message: "create storage error",
				Code:    http.StatusInternalServerError}
		}
		// TODO: move it in the DB CreateStorage method?
		env.DB.GenerateAndUpdateStorageBarecode(&s)
	}
	s.StorageID = sql.NullInt64{Valid: true, Int64: int64(id)}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)
	return nil
}

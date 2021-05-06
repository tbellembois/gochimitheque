package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/static/jade"
)

/*
	views handlers
*/

// VGetStoragesHandler handles the store location list page
func (env *Env) VGetStoragesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := models.ContainerFromRequestContext(r)

	jade.Storageindex(c, w)

	return nil
}

// VCreateStorageHandler handles the storage creation page
func (env *Env) VCreateStorageHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := models.ContainerFromRequestContext(r)

	jade.Storagecreate(c, w)

	return nil
}

/*
	REST handlers
*/

// ToogleStorageBorrowingHandler (un)borrow the storage with id passed in the request vars
// for the logged user.
func (env *Env) ToogleStorageBorrowingHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	var (
		err error
		s   models.Storage
	)

	if err = json.NewDecoder(r.Body).Decode(&s); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "JSON decoding error",
			Code:    http.StatusInternalServerError}
	}

	// retrieving the logged user id from request context
	c := models.ContainerFromRequestContext(r)
	s.Borrowing.Person = &models.Person{}
	s.Borrowing.Person.PersonID = c.PersonID

	// toggling the borrowing
	err = env.DB.ToogleStorageBorrowing(s)

	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error creating the borrowing",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(s); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

// GetStoragesUnitsHandler returns a json list of the units matching the search criteria
func (env *Env) GetStoragesUnitsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetStoragesUnitsHandler")

	var (
		err  error
		aerr *models.AppError
		dsp  models.DbselectparamUnit
	)

	// init db request parameters
	if dsp, aerr = models.NewdbselectparamUnit(r, nil); err != nil {
		return aerr
	}

	units, count, err := env.DB.GetStoragesUnits(dsp)
	if err != nil {
		return &models.AppError{
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
	if err = json.NewEncoder(w).Encode(resp{Rows: units, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetStoragesSuppliersHandler returns a json list of the suppliers matching the search criteria
func (env *Env) GetStoragesSuppliersHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetStoragesSuppliersHandler")

	var (
		err  error
		aerr *models.AppError
		dsp  models.Dbselectparam
	)

	// init db request parameters
	if dsp, aerr = models.Newdbselectparam(r, nil); err != nil {
		return aerr
	}

	suppliers, count, err := env.DB.GetStoragesSuppliers(dsp)
	if err != nil {
		return &models.AppError{
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
	if err = json.NewEncoder(w).Encode(resp{Rows: suppliers, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetOtherStoragesHandler returns a json list of the storages matching the search criteria
// in other entities with no storage details
func (env *Env) GetOtherStoragesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetOtherStoragesHandler")

	var (
		err      error
		aerr     *models.AppError
		dsps     models.DbselectparamStorage
		exportfn string
	)

	// init db request parameters
	if dsps, aerr = models.NewdbselectparamStorage(r, nil); err != nil {
		return aerr
	}

	entities, count, err := env.DB.GetOtherStorages(dsps)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the storages",
		}
	}

	type resp struct {
		Rows     []models.Entity `json:"rows"`
		Total    int             `json:"total"`
		ExportFN string          `json:"exportfn"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp{Rows: entities, Total: count, ExportFN: exportfn}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

// GetStoragesHandler returns a json list of the storages matching the search criteria
func (env *Env) GetStoragesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetStoragesHandler")

	var (
		err      error
		aerr     *models.AppError
		dsps     models.DbselectparamStorage
		exportfn string
	)

	// init db request parameters
	if dsps, aerr = models.NewdbselectparamStorage(r, nil); err != nil {
		return aerr
	}

	storages, count, err := env.DB.GetStorages(dsps)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the storages",
		}
	}

	// export?
	if _, export := r.URL.Query()["export"]; export {
		if exportfn, err = models.StoragesToCSV(storages); err != nil {
			return &models.AppError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
		// emptying results on exports
		storages = []models.Storage{}
		count = 0
	}

	type resp struct {
		Rows     []models.Storage `json:"rows"`
		Total    int              `json:"total"`
		ExportFN string           `json:"exportfn"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp{Rows: storages, Total: count, ExportFN: exportfn}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

// GetStorageHandler returns a json of the entity with the requested id
func (env *Env) GetStorageHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	storage, err := env.DB.GetStorage(id)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the storage",
		}
	}
	logger.Log.WithFields(logrus.Fields{"storage": storage}).Debug("GetStorageHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(storage); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// UpdateStorageHandler updates the storage from the request form
func (env *Env) UpdateStorageHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
		s   models.Storage
	)
	if err = json.NewDecoder(r.Body).Decode(&s); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "JSON decoding error",
			Code:    http.StatusInternalServerError}
	}

	// retrieving the logged user id from request context
	c := models.ContainerFromRequestContext(r)

	logger.Log.WithFields(logrus.Fields{"s": s}).Debug("UpdateStorageHandler")

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}
	updateds, _ := env.DB.GetStorage(id)
	updateds.StorageModificationDate = time.Now()
	updateds.StorageBarecode = s.StorageBarecode
	updateds.StorageQuantity = s.StorageQuantity
	updateds.Supplier = s.Supplier
	updateds.UnitQuantity = s.UnitQuantity
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
	updateds.StorageNumberOfBag = s.StorageNumberOfBag
	updateds.StorageNumberOfCarton = s.StorageNumberOfCarton
	updateds.StorageNumberOfUnit = s.StorageNumberOfUnit
	logger.Log.WithFields(logrus.Fields{"updateds": updateds}).Debug("UpdateStorageHandler")

	if err := env.DB.UpdateStorage(updateds); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "update storage error",
			Code:    http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode([]models.Storage{updateds}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// DeleteStorageHandler deletes the storage with the requested id
func (env *Env) DeleteStorageHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	if err = env.DB.DeleteStorage(id); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// ArchiveStorageHandler archives the storage with the requested id
func (env *Env) ArchiveStorageHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	if err = env.DB.ArchiveStorage(id); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// RestoreStorageHandler restores the storage with the requested id
func (env *Env) RestoreStorageHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	if err = env.DB.RestoreStorage(id); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// CreateStorageHandler creates the storage from the request form
func (env *Env) CreateStorageHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("CreateStorageHandler")
	var (
		s   models.Storage
		err error
		id  int
	)
	if err = json.NewDecoder(r.Body).Decode(&s); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "JSON decoding error",
			Code:    http.StatusInternalServerError}
	}

	// retrieving the logged user id from request context
	c := models.ContainerFromRequestContext(r)

	// retrieving the full store location
	// we need its entity id to compute the barecode
	if s.StoreLocation, err = env.DB.GetStoreLocation(int(s.StoreLocationID.Int64)); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "error retrieving the storage store location",
			Code:    http.StatusInternalServerError}
	}

	s.StorageCreationDate = time.Now()
	s.StorageModificationDate = time.Now()
	s.PersonID = c.PersonID
	logger.Log.WithFields(logrus.Fields{"s": fmt.Sprintf("%+v", s)}).Debug("CreateStorageHandler")
	logger.Log.WithFields(logrus.Fields{"s.StorageNbItem": s.StorageNbItem}).Debug("CreateStorageHandler")

	if s.StorageNbItem == 0 {
		s.StorageNbItem = 1
	}

	var result []models.Storage
	for i := 1; i <= s.StorageNbItem; i++ {
		if id, err = env.DB.CreateStorage(s, i); err != nil {
			return &models.AppError{
				Error:   err,
				Message: "create storage error",
				Code:    http.StatusInternalServerError}
		}
		result = append(result, models.Storage{
			StorageID: sql.NullInt64{Valid: true, Int64: int64(id)},
		})
	}
	s.StorageID = sql.NullInt64{Valid: true, Int64: int64(id)}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(result); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

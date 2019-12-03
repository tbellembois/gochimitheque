package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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

// VGetStoragesHandler handles the store location list page
func (env *Env) VGetStoragesHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	jade.Storageindex(c, w)

	return nil
}

// VCreateStorageHandler handles the storage creation page
func (env *Env) VCreateStorageHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	jade.Storagecreate(c, w)

	return nil
}

/*
	REST handlers
*/

// ToogleStorageBorrowingHandler (un)borrow the storage with id passed in the request vars
// for the logged user.
func (env *Env) ToogleStorageBorrowingHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	var (
		err         error
		isborrowing bool
		id          int
		b           models.Borrowing
	)
	vars := mux.Vars(r)

	// parsing request form
	if err = r.ParseForm(); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form parsing error",
			Code:    http.StatusBadRequest}
	}
	// decoding request form
	if err = global.Decoder.Decode(&b, r.PostForm); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}
	// getting the storage id
	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}
	b.Storage.StorageID.Int64 = int64(id)

	// retrieving the logged user id from request context
	c := helpers.ContainerFromRequestContext(r)
	b.Person.PersonID = c.PersonID

	if b.Borrower != nil {
		if isborrowing, err = env.DB.IsStorageBorrowing(b); err != nil {
			return &helpers.AppError{
				Error:   err,
				Code:    http.StatusInternalServerError,
				Message: "error getting borrowing status",
			}
		}
	} else {
		isborrowing = true
	}

	// toggling the borrowing
	if isborrowing {
		err = env.DB.DeleteStorageBorrowing(b)
	} else {
		err = env.DB.CreateStorageBorrowing(b)
	}
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error creating the borrowing",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(b.Storage)
	w.WriteHeader(http.StatusOK)
	return nil
}

// GetStoragesUnitsHandler returns a json list of the units matching the search criteria
func (env *Env) GetStoragesUnitsHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	global.Log.Debug("GetStoragesUnitsHandler")

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
	global.Log.Debug("GetStoragesSuppliersHandler")

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

// GetOtherStoragesHandler returns a json list of the storages matching the search criteria
// in other entities with no storage details
func (env *Env) GetOtherStoragesHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	global.Log.Debug("GetOtherStoragesHandler")

	var (
		err      error
		aerr     *helpers.AppError
		dsps     helpers.DbselectparamStorage
		exportfn string
	)

	// init db request parameters
	if dsps, aerr = helpers.NewdbselectparamStorage(r, nil); err != nil {
		return aerr
	}

	entities, count, err := env.DB.GetOtherStorages(dsps)
	if err != nil {
		return &helpers.AppError{
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
	json.NewEncoder(w).Encode(resp{Rows: entities, Total: count, ExportFN: exportfn})

	return nil
}

// GetStoragesHandler returns a json list of the storages matching the search criteria
func (env *Env) GetStoragesHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	global.Log.Debug("GetStoragesHandler")

	var (
		err      error
		aerr     *helpers.AppError
		dsps     helpers.DbselectparamStorage
		exportfn string
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

	// export?
	if _, export := r.URL.Query()["export"]; export {
		exportfn = models.StoragesToCSV(storages)
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
	json.NewEncoder(w).Encode(resp{Rows: storages, Total: count, ExportFN: exportfn})

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
	global.Log.WithFields(logrus.Fields{"storage": storage}).Debug("GetStorageHandler")

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

	if err := global.Decoder.Decode(&s, r.PostForm); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}
	global.Log.WithFields(logrus.Fields{"s": s}).Debug("UpdateStorageHandler")

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
	global.Log.WithFields(logrus.Fields{"updateds": updateds}).Debug("UpdateStorageHandler")

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

// ArchiveStorageHandler archives the storage with the requested id
func (env *Env) ArchiveStorageHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
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

	env.DB.ArchiveStorage(id)
	return nil
}

// RestoreStorageHandler restores the storage with the requested id
func (env *Env) RestoreStorageHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
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

	env.DB.RestoreStorage(id)
	return nil
}

// CreateStorageHandler creates the storage from the request form
func (env *Env) CreateStorageHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	global.Log.Debug("CreateStorageHandler")
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

	if err := global.Decoder.Decode(&s, r.PostForm); err != nil {
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
	global.Log.WithFields(logrus.Fields{"s": s}).Debug("CreateStorageHandler")

	for i := 1; i <= s.StorageNbItem; i++ {
		if id, err = env.DB.CreateStorage(s); err != nil {
			return &helpers.AppError{
				Error:   err,
				Message: "create storage error",
				Code:    http.StatusInternalServerError}
		}
		// generating barecode if not specified
		if s.StorageBarecode.String == "" {
			s.StorageID = sql.NullInt64{Valid: true, Int64: int64(id)}
			if err = env.DB.GenerateAndUpdateStorageBarecode(&s); err != nil {
				return &helpers.AppError{
					Error:   err,
					Message: "generate storage barecode error",
					Code:    http.StatusInternalServerError}
			}
		}
	}
	s.StorageID = sql.NullInt64{Valid: true, Int64: int64(id)}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)
	return nil
}

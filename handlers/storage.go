package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/barweiss/go-tuple"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/request"
	"github.com/tbellembois/gochimitheque/static/jade"
	"github.com/tbellembois/gochimitheque/zmqclient"
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
			OriginalError: err,
			Message:       "JSON decoding error",
			Code:          http.StatusInternalServerError,
		}
	}

	// retrieving the logged user id from request context
	c := request.ContainerFromRequestContext(r)
	s.Borrowing.Person = &models.Person{}
	s.Borrowing.Person.PersonID = c.PersonID

	// toggling the borrowing
	err = env.DB.ToogleStorageBorrowing(s)

	if err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error creating the borrowing",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(s); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

// GetStoragesUnitsHandler godoc
// @Summary Get units.
// @Description `unit_type` can be `temperature`, `concentration` or `quantity`.
// @tags unit
// @Produce json
// @Success 200 {object} models.UnitsResp
// @Failure 500
// @Failure 403
// @Router /units/ [get].
func (env *Env) GetStoragesUnitsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetStoragesUnitsHandler")

	var (
		jsonRawMessage json.RawMessage
		err            error
	)

	if jsonRawMessage, err = zmqclient.DBGetUnits("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetUnits",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)
	if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
		return appErr
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil

}

// GetOtherStoragesHandler returns a json list of the storages matching the search criteria
// in other entities with no storage details.
func (env *Env) GetOtherStoragesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetOtherStoragesHandler")

	var (
		err      error
		filter   zmqclient.RequestFilter
		exportfn string
	)

	c := request.ContainerFromRequestContext(r)

	// init db request parameters
	if filter, err = zmqclient.RequestFilterFromRawString("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.Request_filter",
		}
	}

	entities, count, err := env.DB.GetOtherStorages(filter, c.PersonID)
	if err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error getting the storages",
		}
	}

	type resp struct {
		Rows     []models.Entity `json:"rows"`
		Total    int             `json:"total"`
		ExportFN string          `json:"exportfn"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(resp{Rows: entities, Total: count, ExportFN: exportfn}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

// GetStoragesHandler returns a json list of the storages matching the search criteria.
func (env *Env) GetStoragesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetStoragesHandler")

	var (
		err      error
		filter   zmqclient.RequestFilter
		exportfn string
	)

	c := request.ContainerFromRequestContext(r)

	// init db request parameters
	if filter, err = zmqclient.RequestFilterFromRawString("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.Request_filter",
		}
	}

	storages, count, err := env.DB.GetStorages(filter, c.PersonID)
	if err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error getting the storages",
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

	if err = json.NewEncoder(w).Encode(resp{Rows: storages, Total: count, ExportFN: exportfn}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

// GetStorageHandler godoc
// @Summary Get a storage.
// @tags storage
// @Accept plain
// @Produce json
// @Param id path int true "Storage id."
// @Success 200 {object} models.Storage
// @Failure 500
// @Failure 403
// @Router /storage/{id} [get].
func (env *Env) GetStorageHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)

	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	storage, err := env.DB.GetStorage(id)
	if err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error getting the storage",
		}
	}

	logger.Log.WithFields(logrus.Fields{"storage": storage}).Debug("GetStorageHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(storage); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// UpdateStorageHandler updates the storage from the request form.
func (env *Env) UpdateStorageHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)

	var (
		id  int
		err error
		s   models.Storage
	)

	if err = json.NewDecoder(r.Body).Decode(&s); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "JSON decoding error",
			Code:          http.StatusInternalServerError,
		}
	}

	// retrieving the logged user id from request context
	c := request.ContainerFromRequestContext(r)

	logger.Log.WithFields(logrus.Fields{"s": s}).Debug("UpdateStorageHandler")

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	var updateds models.Storage

	if updateds, err = env.DB.GetStorage(id); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "get storage",
			Code:          http.StatusInternalServerError,
		}
	}

	s.StorageModificationDate = models.MyTime{time.Now()}
	s.StorageID = updateds.StorageID
	s.PersonID = c.PersonID

	logger.Log.WithFields(logrus.Fields{"updateds": updateds}).Debug("UpdateStorageHandler")

	if _, err := env.DB.CreateUpdateStorage(s, 0, true); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "update storage error",
			Code:          http.StatusInternalServerError,
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode([]models.Storage{s}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// DeleteStorageHandler deletes the storage with the requested id.
func (env *Env) DeleteStorageHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)

	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	if err = env.DB.DeleteStorage(id); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       err.Error(),
		}
	}
	return nil
}

// ArchiveStorageHandler archives the storage with the requested id.
func (env *Env) ArchiveStorageHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)

	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	if err = env.DB.ArchiveStorage(id); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// RestoreStorageHandler restores the storage with the requested id.
func (env *Env) RestoreStorageHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)

	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	if err = env.DB.RestoreStorage(id); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// CreateStorageHandler creates the storage from the request form.
func (env *Env) CreateStorageHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("CreateStorageHandler")

	var (
		s              models.Storage
		err            error
		id             int64
		jsonRawMessage json.RawMessage
	)

	if err = json.NewDecoder(r.Body).Decode(&s); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "JSON decoding error",
			Code:          http.StatusInternalServerError,
		}
	}

	// retrieving the logged user id from request context
	c := request.ContainerFromRequestContext(r)

	// getting the store location matching the id
	if jsonRawMessage, err = zmqclient.DBGetStorelocations("http://localhost/?store_location="+strconv.Itoa(int(s.StoreLocationID.Int64)), c.PersonID); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetStorelocations",
		}
	}

	// unmarshalling response
	var tuple tuple.T2[[]models.StoreLocation, int]
	if err = json.Unmarshal(jsonRawMessage, &tuple); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error unmarshalling jsonRawMessage",
		}
	}

	s.StoreLocation = tuple.V1[0]

	s.StorageCreationDate = models.MyTime{time.Now()}
	s.StorageModificationDate = models.MyTime{time.Now()}
	s.PersonID = c.PersonID

	logger.Log.WithFields(logrus.Fields{"s": fmt.Sprintf("%+v", s)}).Debug("CreateStorageHandler")
	logger.Log.WithFields(logrus.Fields{"s.StorageNbItem": s.StorageNbItem}).Debug("CreateStorageHandler")

	if s.StorageNbItem == 0 {
		s.StorageNbItem = 1
	}

	var result []models.Storage

	for i := 1; i <= s.StorageNbItem; i++ {
		if id, err = env.DB.CreateUpdateStorage(s, i, false); err != nil {
			return &models.AppError{
				OriginalError: err,
				Message:       "create storage error",
				Code:          http.StatusInternalServerError,
			}
		}

		result = append(result, models.Storage{
			StorageID: sql.NullInt64{Valid: true, Int64: int64(id)},
		})
	}
	s.StorageID = sql.NullInt64{Valid: true, Int64: id}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(result); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

package handlers

import (
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
		err               error
		storage_id        int
		borrower_id       int
		borrowing_comment *string
	)

	vars := mux.Vars(r)

	logger.Log.WithFields(logrus.Fields{"vars": vars}).Debug("ToogleStorageBorrowingHandler")

	if storage_id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	if var_borrower_id := r.URL.Query().Get("borrower_id"); len(var_borrower_id) > 0 {
		if borrower_id, err = strconv.Atoi(var_borrower_id); err != nil {
			return &models.AppError{
				OriginalError: err,
				Message:       "id atoi conversion",
				Code:          http.StatusInternalServerError,
			}
		}
	}

	if var_borrowing_comment := r.URL.Query().Get("borrowing_comment"); len(var_borrowing_comment) > 0 {
		borrowing_comment = &var_borrowing_comment
	}

	// retrieving the logged user id from request context
	c := request.ContainerFromRequestContext(r)

	logger.Log.WithFields(logrus.Fields{"storage_id": storage_id, "borrower_id": borrower_id}).Debug("ToogleStorageBorrowingHandler")

	if _, err = zmqclient.DBToggleStorageBorrowing(c.PersonID, storage_id, borrower_id, borrowing_comment); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBToggleStorageBorrowing",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode("ok"); err != nil {
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

// GetStoragesHandler returns a json list of the storages matching the search criteria.
func (env *Env) GetStoragesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetStoragesHandler")

	var (
		err            error
		jsonRawMessage json.RawMessage
	)

	c := request.ContainerFromRequestContext(r)

	if jsonRawMessage, err = zmqclient.DBGetStorages("http://localhost"+r.RequestURI, c.PersonID); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetStorages",
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

	s.PersonID = &c.PersonID

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	var updateds *models.Storage

	// getting the storage
	var (
		jsonRawMessage json.RawMessage
	)
	if jsonRawMessage, err = zmqclient.DBGetStorages("http://localhost/"+strconv.Itoa(id), *s.PersonID); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "zmqclient.DBGetStorages",
			Code:          http.StatusInternalServerError,
		}
	}

	if updateds, err = zmqclient.ConvertDBJSONToStorage(jsonRawMessage); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "ConvertDBJSONToStorage",
			Code:          http.StatusInternalServerError,
		}
	}

	s.StorageModificationDate = time.Now()
	s.StorageID = updateds.StorageID
	s.PersonID = &c.PersonID

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
	if jsonRawMessage, err = zmqclient.DBGetStorelocations("http://localhost/store_locations/"+strconv.Itoa(int(*s.StoreLocationID)), c.PersonID); err != nil {
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

	s.StorageCreationDate = time.Now()
	s.StorageModificationDate = time.Now()
	s.PersonID = &c.PersonID

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

		var storage_id int64 = int64(id)
		result = append(result, models.Storage{
			StorageID: &storage_id,
		})
	}
	var storage_id int64 = int64(id)
	s.StorageID = &storage_id

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(result); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/request"
	"github.com/tbellembois/gochimitheque/static/jade"
)

/*
	views handlers
*/

// VGetStoreLocationsHandler handles the store location list page.
func (env *Env) VGetStoreLocationsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Storelocationindex(c, w)

	return nil
}

// VCreateStoreLocationHandler handles the store location creation page.
func (env *Env) VCreateStoreLocationHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Storelocationcreate(c, w)

	return nil
}

/*
	REST handlers
*/

// GetStoreLocationsHandler godoc
// @Summary Get store locations. Only store locations visible by the authenticated user are returned.
// @tags storelocation
// @Produce json
// @Success 200 {object} models.StoreLocationsResp
// @Failure 500
// @Failure 403
// @Router /storelocations/ [get].
func (env *Env) GetStoreLocationsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetStoreLocationsHandler")

	var (
		err    error
		aerr   *models.AppError
		filter *request.Filter
	)

	// init db request parameters
	if filter, aerr = request.NewFilter(r, nil); err != nil {
		return aerr
	}

	storelocations, count, err := env.DB.GetStoreLocations(*filter)
	if err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error getting the store locations",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(models.StoreLocationsResp{Rows: storelocations, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetStoreLocationHandler returns a json of the store location with the requested id.
func (env *Env) GetStoreLocationHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
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

	storelocation, err := env.DB.GetStoreLocation(id)
	if err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error getting the store location",
		}
	}

	logger.Log.WithFields(logrus.Fields{"storelocation": storelocation}).Debug("GetStoreLocationHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(storelocation); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// CreateStoreLocationHandler creates the store location from the request form.
func (env *Env) CreateStoreLocationHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("CreateStoreLocationHandler")

	var (
		sl  models.StoreLocation
		err error
		id  int64
	)

	if err = json.NewDecoder(r.Body).Decode(&sl); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "JSON decoding error",
			Code:          http.StatusInternalServerError,
		}
	}

	logger.Log.WithFields(logrus.Fields{"sl": sl}).Debug("CreateStoreLocationHandler")

	if id, err = env.DB.CreateStoreLocation(sl); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "create store location error",
			Code:          http.StatusInternalServerError,
		}
	}
	sl.StoreLocationID = sql.NullInt64{Valid: true, Int64: id}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(sl); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// UpdateStoreLocationHandler updates the store location from the request form.
func (env *Env) UpdateStoreLocationHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)

	var (
		id  int
		err error
		sl  models.StoreLocation
	)

	if err = json.NewDecoder(r.Body).Decode(&sl); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "JSON decoding error",
			Code:          http.StatusInternalServerError,
		}
	}

	logger.Log.WithFields(logrus.Fields{"sl": sl}).Debug("UpdateStoreLocationHandler")

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	updatedsl, err := env.DB.GetStoreLocation(id)
	if err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "get store location error",
			Code:          http.StatusInternalServerError,
		}
	}
	updatedsl.StoreLocationName = sl.StoreLocationName
	updatedsl.StoreLocationColor = sl.StoreLocationColor
	updatedsl.StoreLocationCanStore = sl.StoreLocationCanStore
	updatedsl.StoreLocation = sl.StoreLocation
	updatedsl.Entity = sl.Entity

	logger.Log.WithFields(logrus.Fields{"updatedsl": updatedsl}).Debug("UpdateStoreLocationHandler")

	if err := env.DB.UpdateStoreLocation(updatedsl); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "update store location error",
			Code:          http.StatusInternalServerError,
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(updatedsl); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// DeleteStoreLocationHandler deletes the store location with the requested id.
func (env *Env) DeleteStoreLocationHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
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

	if err = env.DB.DeleteStoreLocation(id); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

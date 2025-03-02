package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

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

func (env *Env) GetStoreLocationsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetStoreLocationsHandler")

	var (
		err            error
		jsonRawMessage json.RawMessage
	)

	c := request.ContainerFromRequestContext(r)

	if jsonRawMessage, err = zmqclient.DBGetStorelocations("http://localhost/?"+r.URL.RawQuery, c.PersonID); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetStorelocations",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)

	if r.URL.Query().Get("store_location") != "" {
		if jsonresp, appErr = zmqclient.ConvertDBJSONToStorelocationJSON(jsonRawMessage); appErr != nil {
			logger.Log.WithFields(logrus.Fields{"ConvertDBJSONToStorelocationJSON appErr": fmt.Sprintf("%+v", appErr)}).Debug("GetStoreLocationsHandler")

			return appErr
		}
	} else {
		if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
			logger.Log.WithFields(logrus.Fields{"ConvertDBJSONToBSTableJSON appErr": fmt.Sprintf("%+v", appErr)}).Debug("GetStoreLocationsHandler")

			return appErr
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil
}

// CreateStoreLocationHandler creates the store location from the request form.
func (env *Env) CreateStoreLocationHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("CreateStoreLocationHandler")

	var (
		jsonRawMessage json.RawMessage
		body           []byte
		err            error
	)

	if body, err = io.ReadAll(r.Body); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error reading request body",
		}
	}
	logger.Log.Debug("body " + string(body))

	if jsonRawMessage, err = zmqclient.DBCreateUpdateStorelocation(body); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBCreateUpdateStorelocation",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonRawMessage)

	return nil
}

// UpdateStoreLocationHandler updates the store location from the request form.
func (env *Env) UpdateStoreLocationHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("UpdateStoreLocationHandler")

	var (
		jsonRawMessage json.RawMessage
		body           []byte
		err            error
	)

	if body, err = io.ReadAll(r.Body); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error reading request body",
		}
	}
	logger.Log.Debug("body " + string(body))

	if jsonRawMessage, err = zmqclient.DBCreateUpdateStorelocation(body); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBCreateUpdateStorelocation",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonRawMessage)

	return nil
}

// DeleteStoreLocationHandler deletes the store location with the requested id.
func (env *Env) DeleteStoreLocationHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)

	var (
		id             int
		err            error
		jsonRawMessage json.RawMessage
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	if jsonRawMessage, err = zmqclient.DBDeleteStorelocation(id); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBDeleteStorelocation",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonRawMessage)

	return nil
}

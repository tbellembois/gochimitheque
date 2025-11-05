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

// VGetEntitiesHandler handles the entity list page.
func (env *Env) VGetEntitiesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Entityindex(c, w)

	return nil
}

// VCreateEntityHandler handles the entity creation page.
func (env *Env) VCreateEntityHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Entitycreate(c, w)

	return nil
}

/*
	REST handlers
*/

// GetEntitiesHandler returns a json list of the entities matching the search criteria.
func (env *Env) GetEntitiesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetEntitiesHandler")

	var (
		err            error
		jsonRawMessage json.RawMessage
	)

	c := request.ContainerFromRequestContext(r)

	if jsonRawMessage, err = zmqclient.DBGetEntities("http://localhost"+r.RequestURI, c.PersonID); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetEntities",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)

	if request.EndsPathWithDigits(r.RequestURI) || request.HasIDParam(r) {
		if jsonresp, appErr = zmqclient.ConvertDBJSONToEntityJSON(jsonRawMessage); appErr != nil {
			logger.Log.WithFields(logrus.Fields{"ConvertDBJSONToEntityJSON appErr": fmt.Sprintf("%+v", appErr)}).Debug("GetEntitiesHandler")

			return appErr
		}
	} else {
		if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
			logger.Log.WithFields(logrus.Fields{"ConvertDBJSONToBSTableJSON appErr": fmt.Sprintf("%+v", appErr)}).Debug("GetEntitiesHandler")

			return appErr
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil
}

func (env *Env) GetEntityStockHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetEntityStockHandler")

	vars := mux.Vars(r)

	var (
		err            error
		jsonRawMessage json.RawMessage
		product_id     int
	)

	c := request.ContainerFromRequestContext(r)

	if product_id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusBadRequest,
		}
	}

	if jsonRawMessage, err = zmqclient.DBComputeStock(product_id, c.PersonID); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBComputeStock",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonRawMessage)

	return nil
}

// CreateEntityHandler creates the entity from the request form.
func (env *Env) CreateEntityHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
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

	if jsonRawMessage, err = zmqclient.DBCreateUpdateEntity(body); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBCreateUpdateEntity",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonRawMessage)

	return nil
}

// UpdateEntityHandler updates the entity from the request form.
func (env *Env) UpdateEntityHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
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

	if jsonRawMessage, err = zmqclient.DBCreateUpdateEntity(body); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBCreateUpdateEntity",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonRawMessage)

	return nil
}

// DeleteEntityHandler deletes the entity with the requested id.
func (env *Env) DeleteEntityHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)

	var (
		jsonRawMessage json.RawMessage
		id64           int64
		id             int
		err            error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	id64 = int64(id)

	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("DeleteEntityHandler")

	if jsonRawMessage, err = zmqclient.DBDeleteEntity(id64); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBDeleteEntity",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonRawMessage)

	return nil
}

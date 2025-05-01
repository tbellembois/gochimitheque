package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/casbin"
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
		e   models.Entity
		err error
	)

	if err = json.NewDecoder(r.Body).Decode(&e); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "JSON decoding error",
			Code:          http.StatusInternalServerError,
		}
	}

	logger.Log.WithFields(logrus.Fields{"e": e}).Debug("CreateEntityHandler")

	if _, err = env.DB.CreateEntity(e); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "create entity error",
			Code:          http.StatusInternalServerError,
		}
	}

	env.Enforcer = casbin.InitCasbinPolicy(env.DB)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(e); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

// UpdateEntityHandler updates the entity from the request form.
func (env *Env) UpdateEntityHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)

	var (
		id       int
		err      error
		e        models.Entity
		updatede *models.Entity
	)

	c := request.ContainerFromRequestContext(r)

	if err = json.NewDecoder(r.Body).Decode(&e); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "JSON decoding error",
			Code:          http.StatusInternalServerError,
		}
	}

	logger.Log.WithFields(logrus.Fields{"e": e}).Debug("UpdateEntityHandler")

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	// getting the entity
	var (
		jsonRawMessage json.RawMessage
	)
	if jsonRawMessage, err = zmqclient.DBGetEntities("http://localhost/entities/"+strconv.Itoa(id), c.PersonID); err != nil {
		logger.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
		return &models.AppError{
			OriginalError: err,
			Message:       "get entity error",
			Code:          http.StatusInternalServerError,
		}
	}

	if updatede, err = zmqclient.ConvertDBJSONToEntity(jsonRawMessage); err != nil {
		logger.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
		return &models.AppError{
			OriginalError: err,
			Message:       "get entity error",
			Code:          http.StatusInternalServerError,
		}
	}

	updatede.EntityName = e.EntityName
	updatede.EntityDescription = e.EntityDescription
	updatede.Managers = e.Managers

	logger.Log.WithFields(logrus.Fields{"updatede": updatede}).Debug("UpdateEntityHandler")

	if err = env.DB.UpdateEntity(*updatede); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "update entity error",
			Code:          http.StatusInternalServerError,
		}
	}

	env.Enforcer = casbin.InitCasbinPolicy(env.DB)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(updatede); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

// DeleteEntityHandler deletes the entity with the requested id.
func (env *Env) DeleteEntityHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
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

	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("DeleteEntityHandler")

	if err := env.DB.DeleteEntity(id); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "delete entity error",
			Code:          http.StatusInternalServerError,
		}
	}

	return nil
}

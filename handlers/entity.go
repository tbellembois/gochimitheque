package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/casbin"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/request"
	"github.com/tbellembois/gochimitheque/static/jade"
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
		err      error
		aerr     *models.AppError
		entities []models.Entity
		count    int
		filter   *request.Filter
	)

	// retrieving the logged user id from request context
	c := request.ContainerFromRequestContext(r)

	// init db request parameters
	if filter, aerr = request.NewFilter(r, nil); err != nil {
		return aerr
	}

	filter.LoggedPersonID = c.PersonID

	if entities, count, err = env.DB.GetEntities(*filter); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error getting the entities",
		}
	}

	type resp struct {
		Rows  []models.Entity `json:"rows"`
		Total int             `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp{Rows: entities, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

// GetEntityStockHandler returns a json of the stock of the entity with the requested id.
func (env *Env) GetEntityStockHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)

	var (
		pid int
		p   models.Product
		err error
	)

	if pid, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusBadRequest,
		}
	}

	if p, err = env.DB.GetProduct(pid); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error getting the product",
		}
	}

	m := env.DB.ComputeStockEntity(p, r)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(m); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

// GetEntityHandler returns a json of the entity with the requested id.
func (env *Env) GetEntityHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)

	var (
		id     int
		err    error
		entity models.Entity
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	if entity, err = env.DB.GetEntity(id); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error getting the entity",
		}
	}

	logger.Log.WithFields(logrus.Fields{"entity": entity}).Debug("GetEntityHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(entity); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

// GetEntityPeopleHandler return the entity managers.
func (env *Env) GetEntityPeopleHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)

	var (
		id     int
		err    error
		people []models.Person
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	if people, err = env.DB.GetEntityManager(id); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error getting the entity people",
		}
	}

	logger.Log.WithFields(logrus.Fields{"people": people}).Debug("GetEntityPeopleHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(people); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

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
	w.WriteHeader(http.StatusOK)
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
		id          int
		err         error
		e, updatede models.Entity
	)

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

	if updatede, err = env.DB.GetEntity(id); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "get entity error",
			Code:          http.StatusInternalServerError,
		}
	}
	updatede.EntityName = e.EntityName
	updatede.EntityDescription = e.EntityDescription
	updatede.Managers = e.Managers
	updatede.LDAPGroups = e.LDAPGroups

	logger.Log.WithFields(logrus.Fields{"updatede": updatede}).Debug("UpdateEntityHandler")

	if err = env.DB.UpdateEntity(updatede); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "update entity error",
			Code:          http.StatusInternalServerError,
		}
	}

	env.Enforcer = casbin.InitCasbinPolicy(env.DB)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
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

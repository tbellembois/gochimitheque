package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/constants"
	"github.com/tbellembois/gochimitheque/models"
)

/*
	views handlers
*/

// VGetEntitiesHandler handles the entity list page
func (env *Env) VGetEntitiesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := containerFromRequestContext(r)

	if e := env.Templates["entityindex"].Execute(w, c); e != nil {
		return &models.AppError{
			Error:   e,
			Code:    http.StatusInternalServerError,
			Message: "error executing template base",
		}
	}
	return nil
}

// VCreateEntityHandler handles the entity creation page
func (env *Env) VCreateEntityHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := containerFromRequestContext(r)

	if e := env.Templates["entitycreate"].Execute(w, c); e != nil {
		return &models.AppError{
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

// GetEntitiesHandler returns a json list of the entities matching the search criteria
func (env *Env) GetEntitiesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	log.Debug("GetEntitiesHandler")

	var (
		search string
		order  string
		offset uint64
		limit  uint64
		err    error
	)

	if s, ok := r.URL.Query()["search"]; !ok {
		search = ""
	} else {
		search = s[0]
	}
	if o, ok := r.URL.Query()["order"]; !ok {
		order = "asc"
	} else {
		order = o[0]
	}
	if o, ok := r.URL.Query()["offset"]; !ok {
		offset = 0
	} else {
		var of int
		if of, err = strconv.Atoi(o[0]); err != nil {
			return &models.AppError{
				Error:   err,
				Code:    http.StatusInternalServerError,
				Message: "offset atoi conversion",
			}
		}
		offset = uint64(of)
	}
	if l, ok := r.URL.Query()["limit"]; !ok {
		limit = constants.MaxUint64
	} else {
		var lm int
		if lm, err = strconv.Atoi(l[0]); err != nil {
			return &models.AppError{
				Error:   err,
				Code:    http.StatusInternalServerError,
				Message: "limit atoi conversion",
			}
		}
		limit = uint64(lm)
	}

	// retrieving the logged user id from request context
	c := containerFromRequestContext(r)
	entities, err := env.DB.GetEntities(c.PersonID, search, order, offset, limit)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the entities",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entities)
	return nil
}

// GetEntityHandler returns a json of the entity with the requested id
func (env *Env) GetEntityHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
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

	entity, err := env.DB.GetEntity(id)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the entity",
		}
	}
	log.WithFields(log.Fields{"entity": entity}).Debug("GetEntityHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entity)
	return nil
}

// GetEntityPeopleHandler return the entity managers
func (env *Env) GetEntityPeopleHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	log.Debug("GetEntityPeopleHandler")
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

	people, err := env.DB.GetEntityPeople(id)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the entity people",
		}
	}
	log.WithFields(log.Fields{"people": people}).Debug("GetEntityPeopleHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(people)
	return nil
}

// CreateEntityHandler creates the entity from the request form
func (env *Env) CreateEntityHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	log.Debug("CreateEntityHandler")
	var (
		e models.Entity
	)
	if err := r.ParseForm(); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "form parsing error",
			Code:    http.StatusBadRequest}
	}

	var decoder = schema.NewDecoder()
	if err := decoder.Decode(&e, r.PostForm); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}
	log.WithFields(log.Fields{"e": e}).Debug("CreateEntityHandler")

	if err := env.DB.CreateEntity(e); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "create entity error",
			Code:    http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(e)
	return nil
}

// UpdateEntityHandler updates the entity from the request form
func (env *Env) UpdateEntityHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
		e   models.Entity
	)
	if err := r.ParseForm(); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "form parsing error",
			Code:    http.StatusBadRequest}
	}
	var decoder = schema.NewDecoder()
	if err := decoder.Decode(&e, r.PostForm); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}
	log.WithFields(log.Fields{"e": e}).Debug("UpdateEntityHandler")

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	updatede, _ := env.DB.GetEntity(id)
	updatede.EntityName = e.EntityName
	updatede.EntityDescription = e.EntityDescription
	log.WithFields(log.Fields{"updatede": updatede}).Debug("UpdateEntityHandler")

	if err := env.DB.UpdateEntity(updatede); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "update entity error",
			Code:    http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatede)
	return nil
}

// DeleteEntityHandler deletes the entity with the requested id
func (env *Env) DeleteEntityHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
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

	env.DB.DeleteEntity(id)
	return nil
}

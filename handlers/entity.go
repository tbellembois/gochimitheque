package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/models"
)

// VGetEntitiesHandler
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

// GetEntitiesHandler
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
		limit = 0
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

	entities, err := env.DB.GetEntities(search, order, offset, limit)
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

// GetEntityHandler
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

	entity, _ := env.DB.GetEntity(id)
	log.WithFields(log.Fields{"entity": entity}).Debug("GetEntityHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entity)
	return nil
}

// CreateEntityHandler
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

// UpdateEntityHandler
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
	updatede.PersonID = e.PersonID
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

// DeleteEntityHandler
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

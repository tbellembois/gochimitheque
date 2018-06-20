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

// VGetPeopleHandler handles the people list page
func (env *Env) VGetPeopleHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := containerFromRequestContext(r)

	if e := env.Templates["personindex"].Execute(w, c); e != nil {
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

// GetPeopleHandler returns a json list of the people matching the search criteria
func (env *Env) GetPeopleHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	log.Debug("GetPeopleHandler")

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
	people, err := env.DB.GetPeople(c.PersonID, search, order, offset, limit)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the people",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(people)
	return nil
}

// GetPersonHandler returns a json of the person with the requested id
func (env *Env) GetPersonHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
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

	person, _ := env.DB.GetPerson(id)
	log.WithFields(log.Fields{"person": person}).Debug("GetPersonHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(person)
	return nil
}

// GetPersonEntitiesHandler returns a json of the entities of the person with the requested id
func (env *Env) GetPersonEntitiesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
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

	entities, err := env.DB.GetPersonEntities(id)
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

// GetPersonPermissionsHandler returns a json of the permissions of the person with the requested id
func (env *Env) GetPersonPermissionsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
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

	permissions, err := env.DB.GetPersonPermissions(id)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the entities",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(permissions)
	return nil
}

// UpdatePersonHandler updates the person from the request form
func (env *Env) UpdatePersonHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
		p   models.Person
	)
	if err := r.ParseForm(); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "form parsing error",
			Code:    http.StatusBadRequest}
	}
	var decoder = schema.NewDecoder()
	if err := decoder.Decode(&p, r.PostForm); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}
	log.WithFields(log.Fields{"p": p}).Debug("UpdatePersonHandler")

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	updatedp, _ := env.DB.GetPerson(id)
	updatedp.PersonEmail = p.PersonEmail
	log.WithFields(log.Fields{"updatedp": updatedp}).Debug("UpdatePersonHandler")

	if err := env.DB.UpdatePerson(updatedp); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "update person error",
			Code:    http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedp)
	return nil
}

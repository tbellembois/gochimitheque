package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/utils"
)

// ValidatePersonEmailHandler checks that the person email does not already exist
// if an id is given is the request the validator ignore the email of the person with this id
func (env *Env) ValidatePersonEmailHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)
	var (
		err       error
		res       bool
		resp      string
		person    models.Person
		person_id int
	)

	// retrieving the logged user id from request context
	c := containerFromRequestContext(r)

	// init db request parameters
	// FIXME: handle errors
	cp, _ := models.NewSelectParametersFromRequest(r)
	cp.LoggedPersonID = c.PersonID

	// converting the id
	if person_id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	// getting the email
	if err = r.ParseForm(); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "form parsing",
			Code:    http.StatusInternalServerError}
	}
	cp.Search = r.Form.Get("person_email")

	// getting the people matching the email
	people, count, err := env.DB.GetPeople(models.GetPeopleParameters{CP: cp, EntityID: -1})
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the people",
		}
	}

	if count == 0 {
		res = false
	} else if person_id == -1 {
		res = (count == 1)
	} else {
		// getting the person
		if person, err = env.DB.GetPerson(person_id); err != nil {
			return &models.AppError{
				Error:   err,
				Code:    http.StatusBadRequest,
				Message: "error looking for person by email",
			}
		}
		res = (person.PersonID != people[0].PersonID)
	}

	log.WithFields(log.Fields{"vars": vars, "res": res}).Debug("ValidatePersonEmailHandler")
	if res {
		resp = "person with this email already present"
	} else {
		resp = "true"
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	return nil
}

// ValidateEntityNameHandler checks that the entity name does not already exist
// if an id != -1 is given is the request the validator ignore the name of the entity with this id
func (env *Env) ValidateEntityNameHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)
	var (
		err       error
		res       bool
		resp      string
		entity    models.Entity
		entity_id int
	)

	// retrieving the logged user id from request context
	c := containerFromRequestContext(r)

	// init db request parameters
	// FIXME: handle errors
	cp, _ := models.NewSelectParametersFromRequest(r)
	cp.LoggedPersonID = c.PersonID

	// converting the id
	if entity_id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	// getting the name
	if err = r.ParseForm(); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "form parsing",
			Code:    http.StatusInternalServerError}
	}
	cp.Search = r.Form.Get("entity_name")

	// getting the entities matching the name
	entities, count, err := env.DB.GetEntities(models.GetEntitiesParameters{CP: cp})
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the entities",
		}
	}

	if count == 0 {
		res = false
	} else if entity_id == -1 {
		res = (count == 1)
	} else {
		// getting the entity
		if entity, err = env.DB.GetEntity(entity_id); err != nil {
			return &models.AppError{
				Error:   err,
				Code:    http.StatusBadRequest,
				Message: "error looking for entity by id",
			}
		}
		res = (entity.EntityID != entities[0].EntityID)
	}

	log.WithFields(log.Fields{"vars": vars, "res": res}).Debug("ValidateEntityNameHandler")
	if res {
		resp = "entity with this name already present"
	} else {
		resp = "true"
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	return nil
}

// ValidateProductNameHandler checks that a product with the name does not already exist
func (env *Env) ValidateProductNameHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	resp := "bad name"
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	return nil
}

// ValidateProductCasNumberHandler checks that a product with the cas number does not already exist
func (env *Env) ValidateProductCasNumberHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	var (
		err  error
		resp string
	)

	// getting the cas number
	if err = r.ParseForm(); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "form parsing",
			Code:    http.StatusInternalServerError}
	}
	// validating it
	v := utils.IsCasNumber(r.Form.Get("casnumber"))

	if v {
		resp = "true"
	} else {
		resp = "invalid cas number"
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	return nil
}

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/models"
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

	// converting the id
	if person_id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	if person_id == -1 {
		// querying the database
		if res, err = env.DB.IsPersonWithEmail(vars["email"]); err != nil {
			return &models.AppError{
				Error:   err,
				Code:    http.StatusBadRequest,
				Message: "error looking for person by email",
			}
		}
	} else {
		// getting the person
		if person, err = env.DB.GetPerson(person_id); err != nil {
			return &models.AppError{
				Error:   err,
				Code:    http.StatusBadRequest,
				Message: "error looking for person by email",
			}
		}
		// querying the database
		if res, err = env.DB.IsPersonWithEmailExcept(vars["email"], person.PersonEmail); err != nil {
			return &models.AppError{
				Error:   err,
				Code:    http.StatusBadRequest,
				Message: "error looking for person by email",
			}
		}
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
// if an id is given is the request the validator ignore the name of the entity with this id
func (env *Env) ValidateEntityNameHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)
	var (
		err       error
		res       bool
		resp      string
		entity    models.Entity
		entity_id int
	)

	// converting the id
	if entity_id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	if entity_id == -1 {
		// querying the database
		if res, err = env.DB.IsEntityWithName(vars["name"]); err != nil {
			return &models.AppError{
				Error:   err,
				Code:    http.StatusBadRequest,
				Message: "error looking for entity by name",
			}
		}
	} else {
		// getting the entity
		if entity, err = env.DB.GetEntity(entity_id); err != nil {
			return &models.AppError{
				Error:   err,
				Code:    http.StatusBadRequest,
				Message: "error looking for entity by name",
			}
		}
		// querying the database
		if res, err = env.DB.IsEntityWithNameExcept(vars["name"], entity.EntityName); err != nil {
			return &models.AppError{
				Error:   err,
				Code:    http.StatusBadRequest,
				Message: "error looking if entity name exist",
			}
		}
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
	vars := mux.Vars(r)
	var (
		err  error
		res  bool
		resp string
	)

	// querying the database
	if res, err = env.DB.IsProductWithName(vars["name"]); err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusBadRequest,
			Message: "error looking if product name exist",
		}
	}

	log.WithFields(log.Fields{"res": res}).Debug("ValidateProductNameHandler")
	if res {
		resp = "product with this name already present"
	} else {
		resp = "true"
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	return nil
}

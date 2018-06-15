package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/models"
	"net/http"
	"strconv"
)

type LoginNameResp struct {
	Name string `json:"name"`
}

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
				Message: "error looking for entity by name",
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

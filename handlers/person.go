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

// VCreatePersonHandler handles the person creation page.
func (env *Env) VCreatePersonHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Personcreate(c, w)

	return nil
}

// VGetPeopleHandler handles the people list page.
func (env *Env) VGetPeopleHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Personindex(c, w)

	return nil
}

/*
	REST handlers
*/

// GetPeopleHandler returns a json list of the people matching the search criteria.
func (env *Env) GetPeopleHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetPeopleHandler")

	var (
		err            error
		jsonRawMessage json.RawMessage
	)

	c := request.ContainerFromRequestContext(r)

	if jsonRawMessage, err = zmqclient.DBGetPeople("http://localhost"+r.RequestURI, c.PersonID); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBGetPeople",
		}
	}

	var (
		jsonresp []byte
		appErr   *models.AppError
	)
	if request.EndsPathWithDigits(r.RequestURI) || request.HasIDParam(r) {
		if jsonresp, appErr = zmqclient.ConvertDBJSONToPersonJSON(jsonRawMessage); appErr != nil {
			logger.Log.WithFields(logrus.Fields{"ConvertDBJSONToPersonJSON appErr": fmt.Sprintf("%+v", appErr)}).Debug("GetPeopleHandler")

			return appErr
		}
	} else {
		if jsonresp, appErr = zmqclient.ConvertDBJSONToBSTableJSON(jsonRawMessage); appErr != nil {
			logger.Log.WithFields(logrus.Fields{"ConvertDBJSONToBSTableJSON appErr": fmt.Sprintf("%+v", appErr)}).Debug("GetPeopleHandler")

			return appErr
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonresp)

	return nil
}

// UpdatePersonHandler updates the person from the request form.
func (env *Env) UpdatePersonHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)

	var (
		id  int
		err error
		p   models.Person
		es  []*models.Entity
	)

	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "JSON decoding error",
			Code:          http.StatusInternalServerError,
		}
	}

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("UpdatePersonHandler")

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	// TODO: remove 1 by connected user id.
	var (
		jsonRawMessage json.RawMessage
		updatedp       *models.Person
	)

	if jsonRawMessage, err = zmqclient.DBGetPeople("http://localhost/"+strconv.Itoa(id), 1); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "zmqclient.DBGetPeople",
			Code:          http.StatusInternalServerError,
		}
	}

	if updatedp, err = zmqclient.ConvertDBJSONToPerson(jsonRawMessage); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "ConvertDBJSONToPerson",
			Code:          http.StatusInternalServerError,
		}
	}

	updatedp.PersonEmail = p.PersonEmail
	updatedp.Entities = p.Entities
	updatedp.Permissions = p.Permissions

	es = updatedp.ManagedEntities

	logger.Log.WithFields(logrus.Fields{"es": es}).Debug("UpdatePersonHandler")

	// for the managed entities setting up the permissions
	if es != nil && len(es) != 0 {
		for _, e := range es {
			updatedp.Permissions = append(updatedp.Permissions, &models.Permission{
				PermissionName:   "all",
				PermissionItem:   "all",
				PermissionEntity: e.EntityID,
				Person:           *updatedp,
			})
		}
	}

	// product permissions are not for a given entity
	for i, p := range updatedp.Permissions {
		if p.PermissionName == "products" {
			updatedp.Permissions[i].PermissionEntity = -1
		}
	}

	logger.Log.WithFields(logrus.Fields{"updatedp": updatedp}).Debug("UpdatePersonHandler")
	logger.Log.WithFields(logrus.Fields{"updatedp.Permissions": updatedp.Permissions}).Debug("UpdatePersonHandler")

	if err = env.DB.UpdatePerson(*updatedp); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "update person error",
			Code:          http.StatusInternalServerError,
		}
	}

	env.Enforcer = casbin.InitCasbinPolicy(env.DB)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(updatedp); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// DeletePersonHandler deletes the person with the requested id.
func (env *Env) DeletePersonHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
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

	if err := env.DB.DeletePerson(id); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "delete person error",
			Code:          http.StatusInternalServerError,
		}
	}

	return nil
}

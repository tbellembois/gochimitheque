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
		err    error
		filter zmqclient.Filter
	)

	c := request.ContainerFromRequestContext(r)

	// init db request parameters
	// if filter, aerr = request.NewFilter(r); aerr != nil {
	// 	return aerr
	// }
	if filter, err = zmqclient.Request_filter("http://localhost/?" + r.URL.RawQuery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.Request_filter",
		}
	}

	people, count, err := env.DB.GetPeople(filter, c.PersonID)
	if err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error getting the people",
		}
	}

	type resp struct {
		Rows  []models.Person `json:"rows"`
		Total int             `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(resp{Rows: people, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetPersonHandler returns a json of the person with the requested id.
func (env *Env) GetPersonHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)

	var (
		id     int
		person models.Person
		err    error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	if person, err = env.DB.GetPerson(id); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "get person error",
			Code:          http.StatusInternalServerError,
		}
	}

	logger.Log.WithFields(logrus.Fields{
		"person": person,
	}).Debug("GetPersonHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(person); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetPersonManageEntitiesHandler returns a json of the entities the person with the requested id is manager of.
func (env *Env) GetPersonManageEntitiesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
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

	entities, err := env.DB.GetPersonManageEntities(id)
	if err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error getting the entities",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(entities); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetPersonEntitiesHandler returns a json of the entities of the person with the requested id.
func (env *Env) GetPersonEntitiesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
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

	// retrieving the logged user id from request context
	c := request.ContainerFromRequestContext(r)
	entities, err := env.DB.GetPersonEntities(c.PersonID, id)
	if err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error getting the entities",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(entities); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetPersonPermissionsHandler returns a json of the permissions of the person with the requested id.
func (env *Env) GetPersonPermissionsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
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

	permissions, err := env.DB.GetPersonPermissions(id)
	if err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error getting the entities",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(permissions); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// CreatePersonHandler creates the person from the request form.
func (env *Env) CreatePersonHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	var (
		p   models.Person
		err error
	)

	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "JSON decoding error",
			Code:          http.StatusInternalServerError,
		}
	}

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("CreatePersonHandler")

	if _, err := env.DB.CreatePerson(p); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "create person error",
			Code:          http.StatusInternalServerError,
		}
	}

	env.Enforcer = casbin.InitCasbinPolicy(env.DB)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(p); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// UpdatePersonHandler updates the person from the request form.
func (env *Env) UpdatePersonHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)

	var (
		id          int
		err         error
		p, updatedp models.Person
		es          []models.Entity
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

	if updatedp, err = env.DB.GetPerson(id); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "error getting the person",
			Code:          http.StatusInternalServerError,
		}
	}
	updatedp.PersonEmail = p.PersonEmail
	updatedp.Entities = p.Entities
	updatedp.Permissions = p.Permissions

	// checking if the person is a manager
	if es, err = env.DB.GetPersonManageEntities(id); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "error getting entities managers",
			Code:          http.StatusInternalServerError,
		}
	}

	logger.Log.WithFields(logrus.Fields{"es": es}).Debug("UpdatePersonHandler")

	// for the managed entities setting up the permissions
	if len(es) != 0 {
		for _, e := range es {
			updatedp.Permissions = append(updatedp.Permissions, &models.Permission{
				PermissionPermName: "all",
				PermissionItemName: "all",
				PermissionEntityID: e.EntityID,
				Person:             updatedp,
			})
		}
	}

	// product permissions are not for a given entity
	for i, p := range updatedp.Permissions {
		if p.PermissionItemName == "products" {
			updatedp.Permissions[i].PermissionEntityID = -1
		}
	}

	logger.Log.WithFields(logrus.Fields{"updatedp": updatedp}).Debug("UpdatePersonHandler")
	logger.Log.WithFields(logrus.Fields{"updatedp.Permissions": updatedp.Permissions}).Debug("UpdatePersonHandler")

	if err = env.DB.UpdatePerson(updatedp); err != nil {
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

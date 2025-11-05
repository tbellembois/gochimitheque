package handlers

import (
	"encoding/json"
	"fmt"
	"io"
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

	logger.Log.Debug("UpdatePersonHandler")

	var (
		jsonRawMessage json.RawMessage
		body           []byte
		err            error
	)

	if body, err = io.ReadAll(r.Body); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error reading request body",
		}
	}
	logger.Log.Debug("body " + string(body))

	if jsonRawMessage, err = zmqclient.DBCreateUpdatePerson(body); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBCreateUpdatePerson",
		}
	}

	env.Enforcer = casbin.InitCasbinPolicy(env.DB)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonRawMessage)

	return nil

}

// DeletePersonHandler deletes the person with the requested id.
func (env *Env) DeletePersonHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	vars := mux.Vars(r)

	var (
		jsonRawMessage json.RawMessage
		id64           int64
		id             int
		err            error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "id atoi conversion",
			Code:          http.StatusInternalServerError,
		}
	}

	id64 = int64(id)

	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("DeletePersonHandler")

	if jsonRawMessage, err = zmqclient.DBDeletePerson(id64); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error calling zmqclient.DBDeletePerson",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonRawMessage)

	return nil
}

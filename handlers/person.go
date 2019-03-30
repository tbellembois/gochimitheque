package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/tbellembois/gochimitheque/utils"

	"../jade"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/global"
	"github.com/tbellembois/gochimitheque/helpers"
	"github.com/tbellembois/gochimitheque/models"
)

/*
	views handlers
*/

// VUpdatePersonPasswordHandler handles the person password update page
func (env *Env) VUpdatePersonPasswordHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	jade.Personpupdate(c, w)

	return nil
}

// VCreatePersonHandler handles the person creation page
func (env *Env) VCreatePersonHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	jade.Personcreate(c, w)

	return nil
}

// VGetPeopleHandler handles the people list page
func (env *Env) VGetPeopleHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	jade.Personindex(c, w)

	return nil
}

/*
	REST handlers
*/

// GetPeopleHandler returns a json list of the people matching the search criteria
func (env *Env) GetPeopleHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	log.Debug("GetPeopleHandler")

	var (
		err  error
		aerr *helpers.AppError
		dspp helpers.DbselectparamPerson
	)

	// init db request parameters
	if dspp, aerr = helpers.NewdbselectparamPerson(r, nil); err != nil {
		return aerr
	}

	people, count, err := env.DB.GetPeople(dspp)
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the people",
		}
	}

	type resp struct {
		Rows  []models.Person `json:"rows"`
		Total int             `json:"total"`
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp{Rows: people, Total: count})
	return nil
}

// GetPersonHandler returns a json of the person with the requested id
func (env *Env) GetPersonHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &helpers.AppError{
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

// GetPersonManageEntitiesHandler returns a json of the entities the person with the requested id is manager of
func (env *Env) GetPersonManageEntitiesHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	entities, err := env.DB.GetPersonManageEntities(id)
	if err != nil {
		return &helpers.AppError{
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

// GetPersonEntitiesHandler returns a json of the entities of the person with the requested id
func (env *Env) GetPersonEntitiesHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	// retrieving the logged user id from request context
	c := helpers.ContainerFromRequestContext(r)
	entities, err := env.DB.GetPersonEntities(c.PersonID, id)
	if err != nil {
		return &helpers.AppError{
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
func (env *Env) GetPersonPermissionsHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	permissions, err := env.DB.GetPersonPermissions(id)
	if err != nil {
		return &helpers.AppError{
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

// CreatePersonHandler creates the person from the request form
func (env *Env) CreatePersonHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	var (
		p   models.Person
		err error
	)
	if err := r.ParseForm(); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form parsing error",
			Code:    http.StatusBadRequest}
	}
	if err := global.Decoder.Decode(&p, r.PostForm); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}
	log.WithFields(log.Fields{"p": p}).Debug("CreatePersonHandler")

	// generating a random password
	// the user will have to get a new password
	// from the login page
	p.PersonPassword = utils.RandStringBytes(64)

	if err, _ := env.DB.CreatePerson(p); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "create person error",
			Code:    http.StatusInternalServerError}
	}

	// sending the new mail
	msgbody := fmt.Sprintf(env.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "createperson_mailbody", PluralCount: 1}), global.ProxyURL+global.ProxyPath+"login", p.PersonEmail)
	msgsubject := env.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "createperson_mailsubject", PluralCount: 1})
	if err = sendMail(p.PersonEmail, msgsubject, msgbody); err != nil {
		return &helpers.AppError{
			Code:    http.StatusInternalServerError,
			Error:   err,
			Message: "error sending the new person mail",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
	return nil
}

// UpdatePersonpHandler updates the person password from the request form
func (env *Env) UpdatePersonpHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	var (
		err error
		p   models.Person
	)
	if err := r.ParseForm(); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form parsing error",
			Code:    http.StatusBadRequest}
	}
	if err := global.Decoder.Decode(&p, r.PostForm); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}
	log.WithFields(log.Fields{"p": p}).Debug("UpdatePersonpHandler")

	// retrieving the logged user id from request context
	c := helpers.ContainerFromRequestContext(r)

	updatedp, _ := env.DB.GetPerson(c.PersonID)
	updatedp.PersonPassword = p.PersonPassword
	log.WithFields(log.Fields{"updatedp": updatedp}).Debug("UpdatePersonpHandler")

	if err = env.DB.UpdatePersonPassword(updatedp); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "update person password error",
			Code:    http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedp)
	return nil
}

// UpdatePersonHandler updates the person from the request form
func (env *Env) UpdatePersonHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
		p   models.Person
	)
	if err := r.ParseForm(); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form parsing error",
			Code:    http.StatusBadRequest}
	}
	if err := global.Decoder.Decode(&p, r.PostForm); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}
	log.WithFields(log.Fields{"p": p}).Debug("UpdatePersonHandler")

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	updatedp, _ := env.DB.GetPerson(id)
	updatedp.PersonEmail = p.PersonEmail
	updatedp.Entities = p.Entities
	updatedp.Permissions = p.Permissions

	// product permissions are not for a given entity
	for i, p := range updatedp.Permissions {
		if p.PermissionItemName == "products" {
			updatedp.Permissions[i].PermissionEntityID = -1
		}
	}
	log.WithFields(log.Fields{"updatedp": updatedp}).Debug("UpdatePersonHandler")

	if err := env.DB.UpdatePerson(updatedp); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "update person error",
			Code:    http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedp)
	return nil
}

// DeletePersonHandler deletes the person with the requested id
func (env *Env) DeletePersonHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	vars := mux.Vars(r)
	var (
		id  int
		err error
	)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	if err := env.DB.DeletePerson(id); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "delete person error",
			Code:    http.StatusInternalServerError}
	}

	return nil
}

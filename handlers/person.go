package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/locales"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/mailer"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/static/jade"
)

/*
	views handlers
*/

// VUpdatePersonPasswordHandler handles the person password update page
func (env *Env) VUpdatePersonPasswordHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := models.ContainerFromRequestContext(r)

	jade.Personpupdate(c, w)

	return nil
}

// VCreatePersonHandler handles the person creation page
func (env *Env) VCreatePersonHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := models.ContainerFromRequestContext(r)

	jade.Personcreate(c, w)

	return nil
}

// VGetPeopleHandler handles the people list page
func (env *Env) VGetPeopleHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := models.ContainerFromRequestContext(r)

	jade.Personindex(c, w)

	return nil
}

/*
	REST handlers
*/

// GetPeopleHandler returns a json list of the people matching the search criteria
func (env *Env) GetPeopleHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetPeopleHandler")

	var (
		err  error
		aerr *models.AppError
		dspp models.DbselectparamPerson
	)

	// init db request parameters
	if dspp, aerr = models.NewdbselectparamPerson(r, nil); aerr != nil {
		return aerr
	}

	people, count, err := env.DB.GetPeople(dspp)
	if err != nil {
		return &models.AppError{
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
	if err = json.NewEncoder(w).Encode(resp{Rows: people, Total: count}); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
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
	logger.Log.WithFields(logrus.Fields{"person": person}).Debug("GetPersonHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(person); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetPersonManageEntitiesHandler returns a json of the entities the person with the requested id is manager of
func (env *Env) GetPersonManageEntitiesHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
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

	entities, err := env.DB.GetPersonManageEntities(id)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the entities",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(entities); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
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

	// retrieving the logged user id from request context
	c := models.ContainerFromRequestContext(r)
	entities, err := env.DB.GetPersonEntities(c.PersonID, id)
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the entities",
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(entities); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
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
	if err = json.NewEncoder(w).Encode(permissions); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// CreatePersonHandler creates the person from the request form
func (env *Env) CreatePersonHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	var (
		p   models.Person
		err error
	)

	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "JSON decoding error",
			Code:    http.StatusInternalServerError}
	}

	// if err := r.ParseForm(); err != nil {
	// 	return &models.AppError{
	// 		Error:   err,
	// 		Message: "form parsing error",
	// 		Code:    http.StatusBadRequest}
	// }
	// if err := globals.Decoder.Decode(&p, r.PostForm); err != nil {
	// 	return &models.AppError{
	// 		Error:   err,
	// 		Message: "form decoding error",
	// 		Code:    http.StatusBadRequest}
	// }
	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("CreatePersonHandler")

	// generating a random password
	// the user will have to get a new password
	// from the login page
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, 64)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	p.PersonPassword = string(b)

	if _, err := env.DB.CreatePerson(p); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "create person error",
			Code:    http.StatusInternalServerError}
	}

	// sending the new mail
	msgbody := fmt.Sprintf(locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "createperson_mailbody", PluralCount: 1}), env.ApplicationFullURL, p.PersonEmail)
	msgsubject := locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "createperson_mailsubject", PluralCount: 1})
	if err = mailer.SendMail(p.PersonEmail, msgsubject, msgbody); err != nil {
		logger.Log.Errorf("error sending email %s", err.Error())
		// return &models.AppError{
		// 	Code:    http.StatusInternalServerError,
		// 	Error:   err,
		// 	Message: "error sending the new person mail",
		// }
	}

	env.InitCasbinPolicy()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(p); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// UpdatePersonpHandler updates the person password from the request form
func (env *Env) UpdatePersonpHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	var (
		err error
		p   models.Person
	)

	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "JSON decoding error",
			Code:    http.StatusInternalServerError}
	}
	// if err = r.ParseForm(); err != nil {
	// 	return &models.AppError{
	// 		Error:   err,
	// 		Message: "form parsing error",
	// 		Code:    http.StatusBadRequest}
	// }
	// if err = globals.Decoder.Decode(&p, r.PostForm); err != nil {
	// 	return &models.AppError{
	// 		Error:   err,
	// 		Message: "form decoding error",
	// 		Code:    http.StatusBadRequest}
	// }
	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("UpdatePersonpHandler")

	// retrieving the logged user id from request context
	c := models.ContainerFromRequestContext(r)

	updatedp, _ := env.DB.GetPerson(c.PersonID)
	updatedp.PersonPassword = p.PersonPassword
	logger.Log.WithFields(logrus.Fields{"updatedp": updatedp}).Debug("UpdatePersonpHandler")

	if err = env.DB.UpdatePersonPassword(updatedp); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "update person password error",
			Code:    http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(updatedp); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// UpdatePersonHandler updates the person from the request form
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
			Error:   err,
			Message: "JSON decoding error",
			Code:    http.StatusInternalServerError}
	}
	// if err = r.ParseForm(); err != nil {
	// 	return &models.AppError{
	// 		Error:   err,
	// 		Message: "form parsing error",
	// 		Code:    http.StatusBadRequest}
	// }
	// if err = globals.Decoder.Decode(&p, r.PostForm); err != nil {
	// 	return &models.AppError{
	// 		Error:   err,
	// 		Message: "form decoding error",
	// 		Code:    http.StatusBadRequest}
	// }
	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("UpdatePersonHandler")

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "id atoi conversion",
			Code:    http.StatusInternalServerError}
	}

	if updatedp, err = env.DB.GetPerson(id); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "error getting the person",
			Code:    http.StatusInternalServerError}
	}
	updatedp.PersonEmail = p.PersonEmail
	updatedp.Entities = p.Entities
	updatedp.Permissions = p.Permissions

	// checking if the person is a manager
	if es, err = env.DB.GetPersonManageEntities(id); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "error getting entities managers",
			Code:    http.StatusInternalServerError}
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
			Error:   err,
			Message: "update person error",
			Code:    http.StatusInternalServerError}
	}

	// hidden feature
	if p.PersonPassword != "" {
		logger.Log.Debug("hidden feature person password set")
		if err = env.DB.UpdatePersonPassword(p); err != nil {
			return &models.AppError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
	}

	env.InitCasbinPolicy()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(updatedp); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// DeletePersonHandler deletes the person with the requested id
func (env *Env) DeletePersonHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
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

	if err := env.DB.DeletePerson(id); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "delete person error",
			Code:    http.StatusInternalServerError}
	}

	return nil
}

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/tbellembois/gochimitheque/aes"
	"github.com/tbellembois/gochimitheque/casbin"
	"github.com/tbellembois/gochimitheque/ldap"
	"github.com/tbellembois/gochimitheque/locales"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/mailer"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/request"
	"github.com/tbellembois/gochimitheque/static/jade"
)

/*
	views handlers
*/

// VUpdatePersonPasswordHandler handles the person qrcode update page.
func (env *Env) VUpdatePersonQRCodeHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Personqrcode(c, w)

	return nil
}

// VUpdatePersonPasswordHandler handles the person password update page.
func (env *Env) VUpdatePersonPasswordHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Personpupdate(c, w)

	return nil
}

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

func (env *Env) GetLDAPGroupsHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetLDAPGroupsHandler")

	var (
		err    error
		aerr   *models.AppError
		filter *request.Filter
		result *ldap.LDAPSearchResult
	)

	// init db request parameters
	if filter, aerr = request.NewFilter(r, nil); aerr != nil {
		return aerr
	}

	if env.LDAPConnection, err = ldap.Connect(); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "LDAP connection",
			Code:          http.StatusInternalServerError,
		}
	}

	result, err = env.LDAPConnection.SearchGroup(strings.ReplaceAll(filter.Search, "%", "*"))

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(result); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// GetPeopleHandler returns a json list of the people matching the search criteria.
func (env *Env) GetPeopleHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("GetPeopleHandler")

	var (
		err    error
		aerr   *models.AppError
		filter *request.Filter
	)

	// init db request parameters
	if filter, aerr = request.NewFilter(r, nil); aerr != nil {
		return aerr
	}

	people, count, err := env.DB.GetPeople(*filter)
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

func (env *Env) GenerateQRCodeHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
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

	if person.PersonAESKey, err = aes.GenerateAESKey(); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "generate aes key error",
			Code:          http.StatusInternalServerError,
		}
	}

	if err = env.DB.UpdatePersonAESKey(person); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "update person aes key error",
			Code:          http.StatusInternalServerError,
		}
	}

	// Encoding the password.
	// We need to keep the email unencrypted the retrieve the personnal AES key of the person.
	var (
		encryptedPassword string
	)
	if encryptedPassword, err = aes.Encrypt(person.PersonPassword, person.PersonAESKey); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "encrypt person credentials error",
			Code:          http.StatusInternalServerError,
		}
	}

	logger.Log.WithFields(logrus.Fields{
		"person":               person,
		"encryptedCredentials": encryptedPassword,
	}).Debug("GetPersonHandler")

	if person.QRCode, err = qrcode.Encode(fmt.Sprintf("%s:%s", person.PersonEmail, encryptedPassword), qrcode.Medium, 512); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(person); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

func (env *Env) IsPersonLDAPHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	var (
		result bool
		err    error
	)

	vars := mux.Vars(r)

	if env.LDAPConnection, err = ldap.Connect(); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "LDAP connection",
			Code:          http.StatusInternalServerError,
		}
	}

	if env.LDAPConnection.IsEnabled {
		var sr *ldap.LDAPSearchResult

		if env.LDAPConnection, err = ldap.Connect(); err != nil {
			return &models.AppError{
				OriginalError: err,
				Message:       "LDAP connection",
				Code:          http.StatusInternalServerError,
			}
		}

		if sr, err = env.LDAPConnection.SearchUser(vars["email"]); err != nil {
			return &models.AppError{
				OriginalError: err,
				Message:       "LDAP user bind error",
				Code:          http.StatusInternalServerError,
			}
		}

		if sr.NbResults > 0 {
			result = true
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(result); err != nil {
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

	// Encoding the password.
	// We need to keep the email unencrypted the retrieve the personnal AES key of the person.
	var (
		encryptedPassword string
	)
	if encryptedPassword, err = aes.Encrypt(person.PersonPassword, person.PersonAESKey); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "encrypt person credentials error",
			Code:          http.StatusInternalServerError,
		}
	}

	logger.Log.WithFields(logrus.Fields{
		"person":               person,
		"encryptedCredentials": encryptedPassword,
	}).Debug("GetPersonHandler")

	if person.QRCode, err = qrcode.Encode(fmt.Sprintf("%s:%s", person.PersonEmail, encryptedPassword), qrcode.Medium, 512); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

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

	if err = p.GeneratePassword(); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "password generation error",
			Code:          http.StatusInternalServerError,
		}
	}

	if p.PersonAESKey, err = aes.GenerateAESKey(); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "aeskey generation error",
			Code:          http.StatusInternalServerError,
		}
	}

	if _, err := env.DB.CreatePerson(p); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "create person error",
			Code:          http.StatusInternalServerError,
		}
	}

	// sending the new mail
	msgbody := fmt.Sprintf(locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "createperson_mailbody", PluralCount: 1}), env.AppFullURL, p.PersonEmail)
	msgsubject := locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "createperson_mailsubject", PluralCount: 1})
	if err = mailer.SendMail(p.PersonEmail, msgsubject, msgbody); err != nil {
		logger.Log.Errorf("error sending email %s", err.Error())
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

// UpdatePersonpHandler updates the person password from the request form.
func (env *Env) UpdatePersonpHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	var (
		err error
		p   models.Person
	)

	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "JSON decoding error",
			Code:          http.StatusInternalServerError,
		}
	}

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("UpdatePersonpHandler")

	// retrieving the logged user id from request context
	c := request.ContainerFromRequestContext(r)

	var (
		updatedp models.Person
	)

	if updatedp, err = env.DB.GetPerson(c.PersonID); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "get person error",
			Code:          http.StatusInternalServerError,
		}
	}
	updatedp.PersonPassword = p.PersonPassword

	logger.Log.WithFields(logrus.Fields{"updatedp": updatedp}).Debug("UpdatePersonpHandler")

	if err = env.DB.UpdatePersonPassword(updatedp); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "update person password error",
			Code:          http.StatusInternalServerError,
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err = json.NewEncoder(w).Encode(updatedp); err != nil {
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

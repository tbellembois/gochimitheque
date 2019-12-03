package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/global"
	"github.com/tbellembois/gochimitheque/helpers"
	"github.com/tbellembois/gochimitheque/jade"
	"github.com/tbellembois/gochimitheque/models"
)

/*
	views handlers
*/

// HomeHandler serves the main page
func (env *Env) HomeHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	jade.Home(c, w)

	return nil
}

// WelcomeAnnounceHandler serves the welcome announce edition page
func (env *Env) VWelcomeAnnounceHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	jade.Welcomeannounceindex(c, w)

	return nil
}

// GetWelcomeAnnounceHandler returns a json of the welcome announce
func (env *Env) GetWelcomeAnnounceHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	var (
		err error
	)

	wa, err := env.DB.GetWelcomeAnnounce()
	if err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the welcome announce",
		}
	}
	global.Log.WithFields(logrus.Fields{"wa": wa}).Debug("GetWelcomeAnnounceHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wa)
	return nil
}

// UpdateWelcomeAnnounceHandler updates the entity from the request form
func (env *Env) UpdateWelcomeAnnounceHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	var (
		err error
		wa  models.WelcomeAnnounce
	)
	if err := r.ParseForm(); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form parsing error",
			Code:    http.StatusBadRequest}
	}
	if err := global.Decoder.Decode(&wa, r.PostForm); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "form decoding error",
			Code:    http.StatusBadRequest}
	}
	global.Log.WithFields(logrus.Fields{"wa": wa}).Debug("UpdateWelcomeAnnounceHandler")

	if err = env.DB.UpdateWelcomeAnnounce(wa); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "update welcomeannounce error",
			Code:    http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wa)
	return nil
}

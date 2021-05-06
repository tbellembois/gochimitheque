package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/static/jade"
	"gopkg.in/russross/blackfriday.v2"
)

/*
	views handlers
*/

// HomeHandler serves the main page
func (env *Env) HomeHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := models.ContainerFromRequestContext(r)

	jade.Home(c, w)

	return nil
}

// AboutHandler serves the about page
func (env *Env) AboutHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := models.ContainerFromRequestContext(r)

	jade.About(c, w)

	return nil
}

// VWelcomeAnnounceHandler serves the welcome announce edition page
func (env *Env) VWelcomeAnnounceHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := models.ContainerFromRequestContext(r)

	jade.Welcomeannounceindex(c, w)

	return nil
}

// GetWelcomeAnnounceHandler returns a json of the welcome announce
func (env *Env) GetWelcomeAnnounceHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	var (
		err error
	)

	wa, err := env.DB.GetWelcomeAnnounce()
	wa.WelcomeAnnounceHTML = string(blackfriday.Run([]byte(wa.WelcomeAnnounceText)))
	if err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error getting the welcome announce",
		}
	}
	logger.Log.WithFields(logrus.Fields{"wa": wa}).Debug("GetWelcomeAnnounceHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(wa); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// UpdateWelcomeAnnounceHandler updates the entity from the request form
func (env *Env) UpdateWelcomeAnnounceHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	var (
		err error
		wa  models.WelcomeAnnounce
	)

	if err = json.NewDecoder(r.Body).Decode(&wa); err != nil {
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
	// if err = globals.Decoder.Decode(&wa, r.PostForm); err != nil {
	// 	return &models.AppError{
	// 		Error:   err,
	// 		Message: "form decoding error",
	// 		Code:    http.StatusBadRequest}
	// }
	logger.Log.WithFields(logrus.Fields{"wa": wa}).Debug("UpdateWelcomeAnnounceHandler")

	if err = env.DB.UpdateWelcomeAnnounce(wa); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "update welcomeannounce error",
			Code:    http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(wa); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dchest/passwordreset"
	"github.com/dgrijalva/jwt-go"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"github.com/steambap/captcha"
	"github.com/tbellembois/gochimitheque/locales"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/mailer"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/static/jade"

	"golang.org/x/crypto/bcrypt"
)

/*
	Views handler.
*/

// GetPasswordHash return password hash for the login.
func GetPasswordHash(login string) ([]byte, error) {

	h := hmac.New(sha256.New, []byte("secret"))
	if _, err := h.Write([]byte(login)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil

}

// VSearchHandler return the search div.
func (env *Env) VSearchHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := models.ContainerFromRequestContext(r)

	jade.Search(c, w)

	return nil
}

// VMenuHandler return the menu div.
func (env *Env) VMenuHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := models.ContainerFromRequestContext(r)

	jade.Menu(c, w)

	return nil
}

// VLoginHandler return the login page.
func (env *Env) VLoginHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := models.ContainerFromRequestContext(r)

	jade.Login(c, w)

	return nil
}

/*
	REST handler
*/

// CaptchaHandler return a captcha image with an uuid
func (env *Env) CaptchaHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	var (
		e    error
		data *captcha.Data
		b    bytes.Buffer
	)

	type captchaData struct {
		Image string `json:"image"`
		UID   string `json:"uid"`
	}
	cd := captchaData{}

	// Create a captcha.
	if data, e = captcha.New(350, 150, func(options *captcha.Options) {
		options.CharPreset = "abcdefghijklmnopqrstuvwxyz0123456789"
		options.TextLength = 4
	}); e != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: e.Error(),
		}
	}

	// Create a token.
	var uuid []byte
	if uuid, e = GetPasswordHash(time.Now().Format("20060102150405")); e != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: e.Error(),
		}
	}
	cd.UID = hex.EncodeToString(uuid)

	// Save the token.
	if e = env.DB.InsertCaptcha(cd.UID, data); e != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: e.Error(),
		}
	}

	// Create response.
	if e = data.WriteImage(&b); e != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: e.Error(),
		}
	}
	cd.Image = base64.StdEncoding.EncodeToString(b.Bytes())

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if e = json.NewEncoder(w).Encode(cd); e != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: e.Error(),
		}
	}

	return nil

}

// RequestResetPasswordHandler reset the user password from the token.
func (env *Env) RequestResetPasswordHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	var (
		tokenquery []string
		token      string
		ok         bool
		err        error
	)

	if tokenquery, ok = r.URL.Query()["token"]; !ok {
		return &models.AppError{
			Code:    http.StatusBadRequest,
			Message: "token not found in request",
		}
	}
	token = tokenquery[0]

	var login string
	if login, err = passwordreset.VerifyToken(token, GetPasswordHash, []byte("secret")); err != nil {
		return &models.AppError{
			Code:    http.StatusForbidden,
			Error:   err,
			Message: "password reset token not verified",
		}
	}

	var person models.Person
	if person, err = env.DB.GetPersonByEmail(login); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Error:   err,
			Message: "error getting user",
		}
	}

	// Generating a random password with the login.
	var ph []byte
	if ph, err = GetPasswordHash(login); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Error:   err,
			Message: "error generating password hash",
		}
	}
	person.PersonPassword = hex.EncodeToString(ph)

	// Updating the person password.
	if err = env.DB.UpdatePersonPassword(person); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Error:   err,
			Message: "error updating the user password",
		}
	}

	// Send reset password mail.
	msgbody := fmt.Sprintf(locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "resetpassword_mailbody1", PluralCount: 1}), person.PersonPassword)
	msgsubject := locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "resetpassword_mailsubject1", PluralCount: 1})
	if err = mailer.SendMail(login, msgsubject, msgbody); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Error:   err,
			Message: "error sending the new password mail",
		}
	}

	// Redirect home.
	msgdone := fmt.Sprintf(locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "resetpassword_done", PluralCount: 1}), person.PersonEmail)
	http.Redirect(w, r, env.ApplicationFullURL+"?message="+msgdone, http.StatusSeeOther)

	return nil

}

// ResetPasswordHandler send a password reinitialisation link by mail
func (env *Env) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	var (
		person models.Person
		err    error
	)

	if err = json.NewDecoder(r.Body).Decode(&person); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "JSON decoding error",
			Code:    http.StatusInternalServerError}
	}

	logger.Log.WithFields(logrus.Fields{
		"person.PersonEmail": person.PersonEmail,
		"person.CaptchaUID":  person.CaptchaUID,
		"person.CaptchaText": person.CaptchaText}).Debug("ResetPasswordHandler")

	// Validate captcha.
	var v bool
	if v, err = env.DB.ValidateCaptcha(person.CaptchaUID, person.CaptchaText); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Error:   err,
			Message: "error validating captcha",
		}
	}
	logger.Log.WithFields(logrus.Fields{"v": v}).Debug("ResetPasswordHandler")

	if !v {
		return &models.AppError{
			Code:    http.StatusBadRequest,
			Error:   nil,
			Message: "captcha not verified",
		}
	}

	// Get the person from db.
	if _, err = env.DB.GetPersonByEmail(person.PersonEmail); err != nil {
		if err == sql.ErrNoRows {
			return &models.AppError{
				Code:    http.StatusUnauthorized,
				Error:   err,
				Message: "user not found in database",
			}
		}
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Error:   err,
			Message: "error getting user",
		}
	}

	// Generate a password hash.
	var hash []byte
	if hash, err = GetPasswordHash(person.PersonEmail); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Error:   err,
			Message: "error generating the password hash",
		}
	}

	// Generate a reinitialization token.
	token := passwordreset.NewToken(person.PersonEmail, 12*time.Hour, hash, []byte("secret"))

	// Send reset password mail.
	msgbody := fmt.Sprintf(locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "resetpassword_mailbody2", PluralCount: 1}), env.ApplicationFullURL, token)
	msgsubject := locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "resetpassword_mailsubject2", PluralCount: 1})
	if err = mailer.SendMail(person.PersonEmail, msgsubject, msgbody); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Error:   err,
			Message: "error sending the reinitialization mail",
		}
	}

	w.WriteHeader(http.StatusOK)

	return nil

}

// DeleteTokenHandler reset the token cookie.
func (env *Env) DeleteTokenHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	logger.Log.Debug("DeleteTokenHandler")
	ctoken := http.Cookie{
		Name:  "token",
		Value: "",
	}
	cemail := http.Cookie{
		Name:  "email",
		Value: "",
	}
	http.SetCookie(w, &ctoken)
	http.SetCookie(w, &cemail)
	http.Redirect(w, r, env.ApplicationFullURL, 307)

	return nil

}

// GetTokenHandler authenticate the user and return a JWT token on success.
func (env *Env) GetTokenHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	var (
		person      models.Person
		personquery *models.Person
		err         error
	)

	if err = json.NewDecoder(r.Body).Decode(&personquery); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "JSON decoding error",
			Code:    http.StatusInternalServerError}
	}

	logger.Log.WithFields(logrus.Fields{"person": personquery}).Debug("GetTokenHandler")

	// Get the person from db.
	if person, err = env.DB.GetPersonByEmail(personquery.PersonEmail); err != nil {
		if err == sql.ErrNoRows {
			return &models.AppError{
				Code:    http.StatusUnauthorized,
				Error:   err,
				Message: "user not found in database",
			}
		}
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Error:   err,
			Message: "error getting user",
		}
	}
	logger.Log.WithFields(logrus.Fields{"db p": person}).Debug("GetTokenHandler")

	// Check password.
	if err = bcrypt.CompareHashAndPassword([]byte(person.PersonPassword), []byte(personquery.PersonPassword)); err != nil {
		return &models.AppError{
			Code:    http.StatusUnauthorized,
			Error:   err,
			Message: locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_password", PluralCount: 1}),
		}
	}

	// Create the JWT token.
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = personquery.PersonEmail
	claims["id"] = personquery.PersonID
	claims["exp"] = time.Now().Add(time.Hour * 8).Unix()

	// Sign the token.
	var tokenString string
	if tokenString, err = token.SignedString(env.TokenSignKey); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Error:   err,
			Message: "error signing token",
		}
	}

	// Write the token to the browser window.
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte(tokenString))
	// Write the token in a cookie.
	// further readings: https://www.calhoun.io/securing-cookies-in-go/
	ctoken := http.Cookie{
		Name:  "token",
		Value: tokenString,
	}
	cemail := http.Cookie{
		Name:  "email",
		Value: person.PersonEmail,
	}
	cid := http.Cookie{
		Name:  "id",
		Value: strconv.Itoa(person.PersonID),
	}
	http.SetCookie(w, &ctoken)
	http.SetCookie(w, &cemail)
	http.SetCookie(w, &cid)

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte(tokenString)); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil

}

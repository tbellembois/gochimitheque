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
	views handlers
*/

// GetPasswordHash return password hash for the login,
func GetPasswordHash(login string) ([]byte, error) {

	h := hmac.New(sha256.New, []byte("secret"))
	if _, err := h.Write([]byte(login)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil

}

// VSearchHandler returns the search page
func (env *Env) VSearchHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := models.ContainerFromRequestContext(r)

	jade.Search(c, w)

	return nil
}

// VMenuHandler returns the menu page
func (env *Env) VMenuHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := models.ContainerFromRequestContext(r)

	jade.Menu(c, w)

	return nil
}

// VLoginHandler returns the login page
func (env *Env) VLoginHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	c := models.ContainerFromRequestContext(r)

	jade.Login(c, w)

	return nil
}

/*
	REST handlers
*/

// CaptchaHandler returns a captcha image with an uuid
func (env *Env) CaptchaHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	var (
		e    error
		data *captcha.Data
		b    bytes.Buffer
	)

	type resp struct {
		Image string `json:"image"`
		UID   string `json:"uid"`
	}
	re := resp{}

	// create a captcha
	if data, e = captcha.New(350, 150, func(options *captcha.Options) {
		options.CharPreset = "abcdefghijklmnopqrstuvwxyz0123456789"
		options.TextLength = 4
	}); e != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: e.Error(),
		}
	}

	// create a token
	var uuid []byte
	if uuid, e = GetPasswordHash(time.Now().Format("20060102150405")); e != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: e.Error(),
		}
	}
	re.UID = hex.EncodeToString(uuid)

	// saving it, retrieving its uuid
	if e = env.DB.InsertCaptcha(re.UID, data); e != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: e.Error(),
		}
	}

	// writing response
	if e = data.WriteImage(&b); e != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: e.Error(),
		}
	}
	re.Image = base64.StdEncoding.EncodeToString(b.Bytes())

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if e = json.NewEncoder(w).Encode(re); e != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: e.Error(),
		}
	}

	return nil
}

// ResetHandler reset the user password from the token
func (env *Env) ResetHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	var (
		t     []string
		token string
		login string
		ok    bool
		err   error
		p     models.Person
	)

	if t, ok = r.URL.Query()["token"]; !ok {
		return &models.AppError{
			Code:    http.StatusBadRequest,
			Message: "token not found in request",
		}
	}
	token = t[0]

	if login, err = passwordreset.VerifyToken(token, GetPasswordHash, []byte("secret")); err != nil {
		return &models.AppError{
			Code:    http.StatusForbidden,
			Error:   err,
			Message: "password reset token not verified",
		}
	}

	// getting the person in the db
	if p, err = env.DB.GetPersonByEmail(login); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Error:   err,
			Message: "error getting user",
		}
	}

	// generating a random password using the login
	brp, _ := GetPasswordHash(login)
	p.PersonPassword = hex.EncodeToString(brp)

	// updating the person password
	if err = env.DB.UpdatePersonPassword(p); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Error:   err,
			Message: "error updating the user password",
		}
	}

	// sending the new mail
	msgbody := fmt.Sprintf(locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "resetpassword_mailbody1", PluralCount: 1}), p.PersonPassword)
	msgsubject := locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "resetpassword_mailsubject1", PluralCount: 1})
	if err = mailer.SendMail(login, msgsubject, msgbody); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Error:   err,
			Message: "error sending the new password mail",
		}
	}

	//w.WriteHeader(http.StatusOK)
	// redirecting to home page
	msgdone := fmt.Sprintf(locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "resetpassword_done", PluralCount: 1}), p.PersonEmail)
	http.Redirect(w, r, env.ApplicationFullURL+"?message="+msgdone, http.StatusSeeOther)

	return nil
}

// ResetPasswordHandler send a password reinitialisation link by mail
func (env *Env) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	var (
		e      error
		hash   []byte
		v      bool
		person models.Person
	)

	if e = json.NewDecoder(r.Body).Decode(&person); e != nil {
		return &models.AppError{
			Error:   e,
			Message: "JSON decoding error",
			Code:    http.StatusInternalServerError}
	}

	// // parsing the form
	// if e = r.ParseForm(); e != nil {
	// 	return &models.AppError{
	// 		Code:    http.StatusBadRequest,
	// 		Error:   e,
	// 		Message: "error parsing form",
	// 	}
	// }

	// // decoding the form
	// person := new(models.Person)
	// if e = globals.Decoder.Decode(person, r.PostForm); e != nil {
	// 	return &models.AppError{
	// 		Code:    http.StatusInternalServerError,
	// 		Error:   e,
	// 		Message: "error decoding form",
	// 	}
	// }
	logger.Log.WithFields(logrus.Fields{
		"person.PersonEmail": person.PersonEmail,
		"person.CaptchaUID":  person.CaptchaUID,
		"person.CaptchaText": person.CaptchaText}).Debug("ResetPasswordHandler")

	// validating captcha
	if v, e = env.DB.ValidateCaptcha(person.CaptchaUID, person.CaptchaText); e != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Error:   e,
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

	// getting the person in the db
	if _, e = env.DB.GetPersonByEmail(person.PersonEmail); e != nil {
		if e == sql.ErrNoRows {
			return &models.AppError{
				Code:    http.StatusUnauthorized,
				Error:   e,
				Message: "user not found in database",
			}
		}
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Error:   e,
			Message: "error getting user",
		}
	}

	// generating a password hash
	if hash, e = GetPasswordHash(person.PersonEmail); e != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Error:   e,
			Message: "error generating the password hash",
		}
	}

	// generating the reinitialization token
	token := passwordreset.NewToken(person.PersonEmail, 12*time.Hour, hash, []byte("secret"))
	// and the mail body
	msgbody := fmt.Sprintf(locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "resetpassword_mailbody2", PluralCount: 1}), env.ApplicationFullURL, token)
	msgsubject := locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "resetpassword_mailsubject2", PluralCount: 1})

	// sending the reinitialisation email
	if e = mailer.SendMail(person.PersonEmail, msgsubject, msgbody); e != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Error:   e,
			Message: "error sending the reinitialization mail",
		}
	}

	w.WriteHeader(http.StatusOK)

	return nil
}

// DeleteTokenHandler actually reset the token cookie
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

	//w.WriteHeader(http.StatusOK)
	http.Redirect(w, r, env.ApplicationFullURL, 307)
	return nil
}

// GetTokenHandler authenticate the user and return a JWT token on success
func (env *Env) GetTokenHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	var (
		e      error
		p      models.Person  // db person
		person *models.Person // form person
	)

	if e = json.NewDecoder(r.Body).Decode(&person); e != nil {
		return &models.AppError{
			Error:   e,
			Message: "JSON decoding error",
			Code:    http.StatusInternalServerError}
	}

	logger.Log.WithFields(logrus.Fields{"person": person}).Debug("GetTokenHandler")

	// authenticating the person
	if p, e = env.DB.GetPersonByEmail(person.PersonEmail); e != nil {
		if e == sql.ErrNoRows {
			return &models.AppError{
				Code:    http.StatusUnauthorized,
				Error:   e,
				Message: "user not found in database",
			}
		}
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Error:   e,
			Message: "error getting user",
		}
	}
	logger.Log.WithFields(logrus.Fields{"db p": p}).Debug("GetTokenHandler")

	if e = bcrypt.CompareHashAndPassword([]byte(p.PersonPassword), []byte(person.PersonPassword)); e != nil {
		return &models.AppError{
			Code:    http.StatusUnauthorized,
			Error:   e,
			Message: locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_password", PluralCount: 1}),
		}
	}

	// create the token
	token := jwt.New(jwt.SigningMethodHS256)

	// create a map to store our claims
	claims := token.Claims.(jwt.MapClaims)

	// set token claims
	claims["email"] = person.PersonEmail
	claims["id"] = person.PersonID
	claims["exp"] = time.Now().Add(time.Hour * 8).Unix()

	// sign the token with our secret
	tokenString, _ := token.SignedString(env.TokenSignKey)

	// finally, write the token to the browser window
	//w.WriteHeader(http.StatusOK)
	//w.Write([]byte(tokenString))
	// finally set the token in a cookie
	// further readings: https://www.calhoun.io/securing-cookies-in-go/
	ctoken := http.Cookie{
		Name:  "token",
		Value: tokenString,
	}
	cemail := http.Cookie{
		Name:  "email",
		Value: p.PersonEmail,
	}
	cid := http.Cookie{
		Name:  "id",
		Value: strconv.Itoa(p.PersonID),
	}
	http.SetCookie(w, &ctoken)
	http.SetCookie(w, &cemail)
	http.SetCookie(w, &cid)

	w.WriteHeader(http.StatusOK)
	if _, e = w.Write([]byte(tokenString)); e != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: e.Error(),
		}
	}

	return nil
}

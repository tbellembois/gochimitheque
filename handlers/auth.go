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
	"strings"
	"time"

	"github.com/dchest/passwordreset"
	"github.com/golang-jwt/jwt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"github.com/steambap/captcha"
	"github.com/tbellembois/gochimitheque/aes"
	"github.com/tbellembois/gochimitheque/casbin"
	"github.com/tbellembois/gochimitheque/ldap"
	"github.com/tbellembois/gochimitheque/locales"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/mailer"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/request"
	"github.com/tbellembois/gochimitheque/static/jade"
	"golang.org/x/crypto/bcrypt"
)

/*
	Views handler.
*/

// GetPasswordHash return password hash for the login.
func (env *Env) GetPasswordHash(login string) ([]byte, error) {
	h := hmac.New(sha256.New, env.TokenSignKey)
	if _, err := h.Write([]byte(login)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// VSearchHandler return the search div.
func (env *Env) VSearchHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Search(c, w)

	return nil
}

// VMenuHandler return the menu div.
func (env *Env) VMenuHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Menu(c, w)

	return nil
}

// VLoginHandler return the login page.
func (env *Env) VLoginHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	c := request.ContainerFromRequestContext(r)

	jade.Login(c, w)

	return nil
}

/*
	REST handler
*/

// CaptchaHandler return a captcha image with an uuid.
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

	if uuid, e = env.GetPasswordHash(time.Now().Format("20060102150405")); e != nil {
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

	if login, err = passwordreset.VerifyToken(token, env.GetPasswordHash, env.TokenSignKey); err != nil {
		return &models.AppError{
			Code:          http.StatusForbidden,
			OriginalError: err,
			Message:       "password reset token not verified",
		}
	}

	var person models.Person

	if person, err = env.DB.GetPersonByEmail(login); err != nil {
		return &models.AppError{
			Code:          http.StatusInternalServerError,
			OriginalError: err,
			Message:       "error getting user",
		}
	}

	// Generating a random password with the login.
	var ph []byte

	if ph, err = env.GetPasswordHash(login); err != nil {
		return &models.AppError{
			Code:          http.StatusInternalServerError,
			OriginalError: err,
			Message:       "error generating password hash",
		}
	}
	person.PersonPassword = hex.EncodeToString(ph)

	// Updating the person password.
	if err = env.DB.UpdatePersonPassword(person); err != nil {
		return &models.AppError{
			Code:          http.StatusInternalServerError,
			OriginalError: err,
			Message:       "error updating the user password",
		}
	}

	// Send reset password mail.
	msgbody := fmt.Sprintf(locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "resetpassword_mailbody1", PluralCount: 1}), person.PersonPassword)
	msgsubject := locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "resetpassword_mailsubject1", PluralCount: 1})

	if err = mailer.SendMail(login, msgsubject, msgbody); err != nil {
		return &models.AppError{
			Code:          http.StatusInternalServerError,
			OriginalError: err,
			Message:       "error sending the new password mail",
		}
	}

	// Redirect home.
	msgdone := fmt.Sprintf(locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "resetpassword_done", PluralCount: 1}), person.PersonEmail)
	http.Redirect(w, r, env.AppFullURL+"?message="+msgdone, http.StatusSeeOther)

	return nil
}

// ResetPasswordHandler send a password reinitialisation link by mail.
func (env *Env) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	var (
		person models.Person
		err    error
	)

	if err = json.NewDecoder(r.Body).Decode(&person); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "JSON decoding error",
			Code:          http.StatusInternalServerError,
		}
	}

	logger.Log.WithFields(logrus.Fields{
		"person.PersonEmail": person.PersonEmail,
		"person.CaptchaUID":  person.CaptchaUID,
		"person.CaptchaText": person.CaptchaText,
	}).Debug("ResetPasswordHandler")

	// Validate captcha.
	var v bool

	if v, err = env.DB.ValidateCaptcha(person.CaptchaUID, person.CaptchaText); err != nil {
		return &models.AppError{
			Code:          http.StatusInternalServerError,
			OriginalError: err,
			Message:       "error validating captcha",
		}
	}

	logger.Log.WithFields(logrus.Fields{"v": v}).Debug("ResetPasswordHandler")

	if !v {
		return &models.AppError{
			Code:          http.StatusBadRequest,
			OriginalError: nil,
			Message:       "captcha not verified",
		}
	}

	// Get the person from db.
	if _, err = env.DB.GetPersonByEmail(person.PersonEmail); err != nil {
		if err == sql.ErrNoRows {
			return &models.AppError{
				Code:          http.StatusUnauthorized,
				OriginalError: err,
				Message:       "user not found in database",
			}
		}

		return &models.AppError{
			Code:          http.StatusInternalServerError,
			OriginalError: err,
			Message:       "error getting user",
		}
	}

	// Generate a password hash.
	var hash []byte

	if hash, err = env.GetPasswordHash(person.PersonEmail); err != nil {
		return &models.AppError{
			Code:          http.StatusInternalServerError,
			OriginalError: err,
			Message:       "error generating the password hash",
		}
	}

	// Generate a reinitialization token.
	token := passwordreset.NewToken(person.PersonEmail, 12*time.Hour, hash, env.TokenSignKey)

	// Send reset password mail.
	msgbody := fmt.Sprintf(locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "resetpassword_mailbody2", PluralCount: 1}), env.AppFullURL, token)
	msgsubject := locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "resetpassword_mailsubject2", PluralCount: 1})
	if err = mailer.SendMail(person.PersonEmail, msgsubject, msgbody); err != nil {
		return &models.AppError{
			Code:          http.StatusInternalServerError,
			OriginalError: err,
			Message:       "error sending the reinitialization mail",
		}
	}

	return nil
}

// DeleteTokenHandler reset the token cookie.
func (env *Env) DeleteTokenHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("DeleteTokenHandler")
	http.Redirect(w, r, env.AppFullURL, http.StatusTemporaryRedirect)

	return nil
}

// GetTokenHandler godoc
// @Summary Authenticate a user.
// @Description Authenticate a user and return a JWT token on success. The token must be passed to each following request.
// @tags authentication
// @Accept json
// @Produce plain
// @Param Person body models.Person true "Person with only `person_email` and `person_password` or `qrcode` (full plain qrcode string) fields."
// @Success 200 {string} Token "qwerty"
// @Failure 500
// @Failure 403
// @Router /get-token [get].
func (env *Env) GetTokenHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	var (
		qrcodeAuthenticated             bool
		person                          models.Person
		personquery                     *models.Person
		userFoundInLDAP, userBindInLDAP bool
		err                             error
	)

	if err = json.NewDecoder(r.Body).Decode(&personquery); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "JSON decoding error",
			Code:          http.StatusInternalServerError,
		}
	}

	logger.Log.WithFields(logrus.Fields{
		"personquery.PersonEmail": personquery.PersonEmail,
		"len(personquery.QRCode)": len(personquery.QRCode),
	}).Debug("GetTokenHandler")

	// If a qrcode is present, decoding it
	if len(personquery.QRCode) > 0 {
		credentials := strings.Split(string(personquery.QRCode), ":")
		personquery.PersonEmail = credentials[0]
		encryptedPassword := credentials[1]

		var dbPerson models.Person

		if dbPerson, err = env.DB.GetPersonByEmail(personquery.PersonEmail); err != nil {
			return &models.AppError{
				Code:          http.StatusInternalServerError,
				OriginalError: err,
				Message:       "error getting user",
			}
		}
		person.PersonEmail = dbPerson.PersonEmail
		person.PersonID = dbPerson.PersonID

		decryptedPassword := aes.Decrypt(encryptedPassword, dbPerson.PersonAESKey)

		logger.Log.WithFields(logrus.Fields{
			"decryptedPassword":       decryptedPassword,
			"dbPerson.PersonPassword": dbPerson.PersonPassword,
		}).Debug("GetTokenHandler")

		if decryptedPassword == dbPerson.PersonPassword {
			qrcodeAuthenticated = true
		}
	}

	if !qrcodeAuthenticated {

		if env.LDAPConnection.IsEnabled {
			var sr *ldap.LDAPSearchResult

			if env.LDAPConnection, err = ldap.Connect(); err != nil {
				return &models.AppError{
					OriginalError: err,
					Message:       "LDAP connection",
					Code:          http.StatusInternalServerError,
				}
			}

			if sr, err = env.LDAPConnection.SearchUser(personquery.PersonEmail); err != nil {
				return &models.AppError{
					OriginalError: err,
					Message:       "LDAP user bind error",
					Code:          http.StatusInternalServerError,
				}
			}

			if sr.NbResults > 0 {
				userFoundInLDAP = true

				userdn := sr.R.Entries[0].DN

				logger.Log.WithFields(logrus.Fields{
					"userdn": userdn,
				}).Debug("GetTokenHandler")

				// Bind as the user to verify their password
				if err = env.LDAPConnection.Bind(userdn, personquery.PersonPassword); err != nil {
					logger.Log.Debug(err)
				} else {
					userBindInLDAP = true
				}
			}
		}

		// Get the person from db.
		if person, err = env.DB.GetPersonByEmail(personquery.PersonEmail); err != nil {
			if err == sql.ErrNoRows {
				if !userFoundInLDAP {
					return &models.AppError{
						Code:          http.StatusUnauthorized,
						OriginalError: err,
						Message:       "user not found in LDAP or database",
					}
				}

				if userFoundInLDAP && !userBindInLDAP {
					return &models.AppError{
						Code:          http.StatusUnauthorized,
						OriginalError: err,
						Message:       "user found in LDAP but no bind",
					}
				}

				if !env.AutoCreateUser {
					return &models.AppError{
						Code:          http.StatusUnauthorized,
						OriginalError: err,
						Message:       "user found and bind in LDAP but not present in DB and auto create user disabled",
					}
				}

				// Auto create user.
				if err = personquery.GeneratePassword(); err != nil {
					return &models.AppError{
						OriginalError: err,
						Message:       "password generation error",
						Code:          http.StatusInternalServerError,
					}
				}
				if personquery.PersonAESKey, err = aes.GenerateAESKey(); err != nil {
					return &models.AppError{
						OriginalError: err,
						Message:       "generate aes key error",
						Code:          http.StatusInternalServerError,
					}
				}
				if _, err := env.DB.CreatePerson(*personquery); err != nil {
					return &models.AppError{
						OriginalError: err,
						Message:       "auto create person error",
						Code:          http.StatusInternalServerError,
					}
				}
				if person, err = env.DB.GetPersonByEmail(personquery.PersonEmail); err != nil {
					return &models.AppError{
						OriginalError: err,
						Message:       "get person error",
						Code:          http.StatusInternalServerError,
					}
				}

				env.Enforcer = casbin.InitCasbinPolicy(env.DB)
			} else {
				return &models.AppError{
					Code:          http.StatusInternalServerError,
					OriginalError: err,
					Message:       "error getting user",
				}
			}
		}

		logger.Log.WithFields(logrus.Fields{"db p": person}).Debug("GetTokenHandler")

		// Check password in local DB.
		if !userBindInLDAP {
			if err = bcrypt.CompareHashAndPassword([]byte(person.PersonPassword), []byte(personquery.PersonPassword)); err != nil {
				return &models.AppError{
					Code:          http.StatusUnauthorized,
					OriginalError: err,
					Message:       locales.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "invalid_password", PluralCount: 1}),
				}
			}
		}

		// Updating the person password in the DB.
		// Needed in the case of LDAP authentication for QRCode generation.
		if err := env.DB.UpdatePersonPassword(*personquery); err != nil {
			return &models.AppError{
				OriginalError: err,
				Message:       "update person password error",
				Code:          http.StatusInternalServerError,
			}
		}
	}

	// Create the JWT token.
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = person.PersonEmail
	// claims["exp"] = time.Now().Add(time.Hour * 8).Unix()

	// Sign the token.
	var tokenString string

	if tokenString, err = token.SignedString(env.TokenSignKey); err != nil {
		return &models.AppError{
			Code:          http.StatusInternalServerError,
			OriginalError: err,
			Message:       "error signing token",
		}
	}

	// Write the token to the browser window.
	//
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

	if _, err = w.Write([]byte(tokenString)); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

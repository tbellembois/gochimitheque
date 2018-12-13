package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/smtp"
	"strconv"
	"time"

	"database/sql"
	"github.com/dchest/passwordreset"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/global"
	"github.com/tbellembois/gochimitheque/helpers"
	"github.com/tbellembois/gochimitheque/models"
	"golang.org/x/crypto/bcrypt"
)

/*
	views handlers
*/

// VLoginHandler returns the login page
func (env *Env) VLoginHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	c := helpers.ContainerFromRequestContext(r)

	if e := env.Templates["login"].ExecuteTemplate(w, "BASE", c); e != nil {
		return &helpers.AppError{
			Error:   e,
			Code:    http.StatusInternalServerError,
			Message: "error executing template base",
		}
	}

	return nil
}

/*
	REST handlers
*/

// getPasswordHash return password hash for the login,
func getPasswordHash(login string) ([]byte, error) {

	h := hmac.New(sha256.New, []byte("secret"))
	h.Write([]byte(login))

	return h.Sum(nil), nil
}

// sendMail send a mail
func sendMail(to string, subject string, body string) error {

	var (
		e    error
		auth smtp.Auth
	)

	if global.MailServerUser != "" {
		// authenticated smtp
		auth = smtp.PlainAuth("", global.MailServerUser, global.MailServerPassword, global.MailServerAddress)
	}
	if e = smtp.SendMail(
		global.MailServerAddress+":"+global.MailServerPort,
		auth,
		global.MailServerSender,
		[]string{to},
		[]byte(subject+body),
	); e != nil {
		return e
	}
	return nil
}

// ResetHandler reset the user password from the token
func (env *Env) ResetHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	var (
		t     []string
		token string
		login string
		ok    bool
		err   error
		p     models.Person
	)

	if t, ok = r.URL.Query()["token"]; !ok {
		return &helpers.AppError{
			Code:    http.StatusBadRequest,
			Message: "token not found in request",
		}
	}
	token = t[0]

	if login, err = passwordreset.VerifyToken(token, getPasswordHash, []byte("secret")); err != nil {
		return &helpers.AppError{
			Code:    http.StatusForbidden,
			Error:   err,
			Message: "password reset token not verified",
		}
	}

	// getting the person in the db
	if p, err = env.DB.GetPersonByEmail(login); err != nil {
		return &helpers.AppError{
			Code:    http.StatusInternalServerError,
			Error:   err,
			Message: "error getting user",
		}
	}

	// generating a random password using the login
	brp, _ := getPasswordHash(login)
	p.PersonPassword = hex.EncodeToString(brp)

	// updating the person password
	if err = env.DB.UpdatePersonPassword(p); err != nil {
		return &helpers.AppError{
			Code:    http.StatusInternalServerError,
			Error:   err,
			Message: "error updating the user password",
		}
	}

	// sending the new mail
	msgbody := "This is your new password: " + p.PersonPassword
	msgsubject := "Subject: Chimithèque new password\r\n"
	if err = sendMail(login, msgsubject, msgbody); err != nil {
		return &helpers.AppError{
			Code:    http.StatusInternalServerError,
			Error:   err,
			Message: "error sending the new password mail",
		}
	}

	//w.WriteHeader(http.StatusOK)
	// redirecting to login page
	http.Redirect(w, r, "/login?message=password%20reinitialized", http.StatusSeeOther)

	return nil
}

// ResetPasswordHandler send a password reinitialisation link by mail
func (env *Env) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	var (
		e    error
		hash []byte
	)

	// parsing the form
	if e = r.ParseForm(); e != nil {
		return &helpers.AppError{
			Code:    http.StatusBadRequest,
			Error:   e,
			Message: "error parsing form",
		}
	}

	// decoding the form
	person := new(models.Person)
	if e = global.Decoder.Decode(person, r.PostForm); e != nil {
		return &helpers.AppError{
			Code:    http.StatusInternalServerError,
			Error:   e,
			Message: "error decoding form",
		}
	}
	log.WithFields(log.Fields{"person.PersonEmail": person.PersonEmail}).Debug("ResetPasswordHandler")

	// getting the person in the db
	if _, e = env.DB.GetPersonByEmail(person.PersonEmail); e != nil {
		if e == sql.ErrNoRows {
			return &helpers.AppError{
				Code:    http.StatusUnauthorized,
				Error:   e,
				Message: "user not found in database",
			}
		}
		return &helpers.AppError{
			Code:    http.StatusInternalServerError,
			Error:   e,
			Message: "error getting user",
		}
	}

	// generating a password hash
	if hash, e = getPasswordHash(person.PersonEmail); e != nil {
		return &helpers.AppError{
			Code:    http.StatusInternalServerError,
			Error:   e,
			Message: "error generating the password hash",
		}
	}

	// generating the reinitialization token
	token := passwordreset.NewToken(person.PersonEmail, 12*time.Hour, hash, []byte("secret"))
	// and the mail body
	msgbody := "Click on this link to reinitialize your password: " + global.ProxyURL + global.ProxyPath + "reset?token=" + token
	msgsubject := "Subject: Chimithèque password reinitialization\r\n"

	// sending the reinitialisation email
	if e = sendMail(person.PersonEmail, msgsubject, msgbody); e != nil {
		return &helpers.AppError{
			Code:    http.StatusInternalServerError,
			Error:   e,
			Message: "error sending the reinitialization mail",
		}
	}

	w.WriteHeader(http.StatusOK)

	return nil
}

// GetTokenHandler authenticate the user and return a JWT token on success
func (env *Env) GetTokenHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	var (
		e      error
		p      models.Person  // db person
		person *models.Person // form person
	)

	// parsing the form
	if e = r.ParseForm(); e != nil {
		return &helpers.AppError{
			Code:    http.StatusBadRequest,
			Error:   e,
			Message: "error parsing form",
		}
	}

	// decoding the form
	person = new(models.Person)
	if e = global.Decoder.Decode(person, r.PostForm); e != nil {
		return &helpers.AppError{
			Code:    http.StatusInternalServerError,
			Error:   e,
			Message: "error decoding form",
		}
	}

	// authenticating the person
	if p, e = env.DB.GetPersonByEmail(person.PersonEmail); e != nil {
		if e == sql.ErrNoRows {
			return &helpers.AppError{
				Code:    http.StatusUnauthorized,
				Error:   e,
				Message: "user not found in database",
			}
		}
		return &helpers.AppError{
			Code:    http.StatusInternalServerError,
			Error:   e,
			Message: "error getting user",
		}
	}
	log.WithFields(log.Fields{"form person": person}).Debug("GetTokenHandler")
	log.WithFields(log.Fields{"db p": p}).Debug("GetTokenHandler")

	if e = bcrypt.CompareHashAndPassword([]byte(p.PersonPassword), []byte(person.PersonPassword)); e != nil {
		return &helpers.AppError{
			Code:    http.StatusUnauthorized,
			Error:   e,
			Message: "invalid password",
		}
	}

	// create the token
	token := jwt.New(jwt.SigningMethodHS256)

	// create a map to store our claims
	claims := token.Claims.(jwt.MapClaims)

	// set token claims
	claims["email"] = person.PersonEmail
	claims["exp"] = time.Now().Add(time.Hour * 8).Unix()

	// sign the token with our secret
	tokenString, _ := token.SignedString(global.TokenSignKey)

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
		Value: person.PersonEmail,
	}
	http.SetCookie(w, &ctoken)
	http.SetCookie(w, &cemail)

	w.WriteHeader(http.StatusOK)

	return nil
}

// HasPermissionHandler returns true if the person with id "personid" has the permission "perm" on item "item" with itemid "itemid"
func (env *Env) HasPermissionHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {
	vars := mux.Vars(r)
	var (
		personid int
		itemid   int
		perm     string
		item     string
		err      error
		p        bool
	)

	if personid, err = strconv.Atoi(vars["personid"]); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "personid atoi conversion",
			Code:    http.StatusInternalServerError}
	}
	if itemid, err = strconv.Atoi(vars["itemid"]); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "itemid atoi conversion",
			Code:    http.StatusInternalServerError}
	}
	perm = vars["perm"]
	item = vars["item"]

	if p, err = env.DB.HasPersonPermission(personid, perm, item, itemid); err != nil {
		return &helpers.AppError{
			Error:   err,
			Message: "getting permissions error",
			Code:    http.StatusInternalServerError}
	}
	log.WithFields(log.Fields{"personid": personid, "perm": perm, "item": item, "itemid": itemid}).Debug("HasPermissionHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
	return nil
}

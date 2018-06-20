package handlers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/models"
	"net/http"
	"strconv"
	"time"
)

// TokenSignKey is the JWT token signing key
var TokenSignKey = []byte("secret")

/*
	views handlers
*/

// VLoginHandler returns the login page
func (env *Env) VLoginHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	if e := env.Templates["login"].Execute(w, nil); e != nil {
		return &models.AppError{
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

// GetTokenHandler authenticate the user and return a JWT token on success
func (env *Env) GetTokenHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	var (
		e error
	)

	// parsing the form
	if e = r.ParseForm(); e != nil {
		return &models.AppError{
			Code:    http.StatusBadRequest,
			Error:   e,
			Message: "error parsing form",
		}
	}

	// decoding the form
	decoder := schema.NewDecoder()
	person := new(models.Person)
	if e = decoder.Decode(person, r.PostForm); e != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Error:   e,
			Message: "error decoding form",
		}
	}
	log.WithFields(log.Fields{"person.PersonEmail": person.PersonEmail}).Debug("GetTokenHandler")

	// authenticating the person
	// TODO: true auth
	if _, e = env.DB.GetPersonByEmail(person.PersonEmail); e != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Error:   e,
			Message: "error getting user",
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
	tokenString, _ := token.SignedString(TokenSignKey)

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

	return nil
}

// HasPermissionHandler returns true if the person with id "personid" has the permission "perm" on item "item" with itemid "itemid"
func (env *Env) HasPermissionHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
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
		return &models.AppError{
			Error:   err,
			Message: "personid atoi conversion",
			Code:    http.StatusInternalServerError}
	}
	if itemid, err = strconv.Atoi(vars["itemid"]); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "itemid atoi conversion",
			Code:    http.StatusInternalServerError}
	}
	perm = vars["perm"]
	item = vars["item"]

	if p, err = env.DB.HasPersonPermission(personid, perm, item, itemid); err != nil {
		return &models.AppError{
			Error:   err,
			Message: "getting permissions error",
			Code:    http.StatusInternalServerError}
	}
	log.WithFields(log.Fields{"personid": personid, "perm": perm, "item": item, "itemid": itemid}).Debug("GetEntityHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
	return nil
}

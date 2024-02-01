package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/net/context"

	"github.com/coreos/go-oidc"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/request"
	"github.com/tbellembois/gochimitheque/static/jade"
)

/*
	Views handler.
*/

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

// DeleteTokenHandler delete all the cookies and redirect to the OIDC logout endpoint.
func (env *Env) DeleteTokenHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("DeleteTokenHandler")

	for _, cookie := range r.Cookies() {
		logger.Log.WithFields(logrus.Fields{
			"cookie.Domain": cookie.Domain,
			"cookie.Name":   cookie.Name,
			"cookie.Path":   cookie.Path,
		}).Debug("DeleteTokenHandler")

		http.SetCookie(w, &http.Cookie{Name: cookie.Name, Value: ""})
	}

	http.Redirect(w, r, env.OIDCEndSessionEndpoint, http.StatusTemporaryRedirect)

	return nil
}

func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func setCallbackCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}

func (env *Env) UserInfoHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	// getting the request context
	ctx := r.Context()
	ctxcontainer := ctx.Value(request.ChimithequeContextKey("container"))
	container := ctxcontainer.(request.Container)

	// getting auth person informations
	userinfo := models.Person{
		PersonID:    container.PersonID,
		PersonEmail: container.PersonEmail,
	}

	logger.Log.WithFields(logrus.Fields{"userinfo": userinfo}).Debug("UserInfoHandler")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(userinfo); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

func (env *Env) CallbackHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	logger.Log.Debug("CallbackHandler")

	state, err := r.Cookie("state")
	if err != nil {
		http.Error(w, "state not found", http.StatusBadRequest)
		return &models.AppError{
			Code:          http.StatusBadRequest,
			OriginalError: err,
			Message:       "state not found",
		}
	}

	if r.URL.Query().Get("state") != state.Value {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return &models.AppError{
			Code:          http.StatusBadRequest,
			OriginalError: err,
			Message:       "state did not match",
		}
	}

	oauth2Token, err := env.OAuth2Config.Exchange(context.Background(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return &models.AppError{
			Code:          http.StatusInternalServerError,
			OriginalError: err,
			Message:       "failed to exchange token",
		}
	}

	logger.Log.WithFields(logrus.Fields{
		"oauth2Token.Expiry": oauth2Token.Expiry,
	}).Debug("CallbackHandler")

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "no id_token field in oauth2 token", http.StatusInternalServerError)
		return &models.AppError{
			Code:          http.StatusInternalServerError,
			OriginalError: err,
			Message:       "no id_token field in oauth2 token",
		}
	}

	idToken, err := env.OIDCVerifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		http.Error(w, "failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return &models.AppError{
			Code:          http.StatusInternalServerError,
			OriginalError: err,
			Message:       "failed to verify ID Token",
		}
	}

	logger.Log.WithFields(logrus.Fields{
		"idToken.Expiry": idToken.Expiry,
	}).Debug("CallbackHandler")

	nonce, err := r.Cookie("nonce")
	if err != nil {
		http.Error(w, "nonce not found", http.StatusBadRequest)
		return &models.AppError{
			Code:          http.StatusBadRequest,
			OriginalError: err,
			Message:       "nonce not found",
		}
	}
	if idToken.Nonce != nonce.Value {
		http.Error(w, "nonce did not match", http.StatusBadRequest)
		return &models.AppError{
			Code:          http.StatusBadRequest,
			OriginalError: err,
			Message:       "nonce did not match",
		}
	}

	// Check the claims.
	var claims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}

	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "error getting claims from token"+err.Error(), http.StatusBadRequest)
		return &models.AppError{
			Code:          http.StatusBadRequest,
			OriginalError: err,
			Message:       "error getting claims from token",
		}
	}

	// Check the email.
	reEmail := regexp.MustCompile(`^[0-9A-Za-z_\.]+@[0-9A-Za-z_\.]+\.[0-9A-Za-z]{2,4}$`)
	if !reEmail.MatchString(claims.Email) {
		http.Error(w, "no email in claims", http.StatusBadRequest)
		return &models.AppError{
			Code:          http.StatusBadRequest,
			OriginalError: nil,
			Message:       "no email in claims",
		}
	}

	// Insert user if DB if needed.
	if _, err = env.DB.GetPersonByEmail(claims.Email); err != nil {
		if err == sql.ErrNoRows {
			if _, err = env.DB.CreatePerson(models.Person{PersonEmail: claims.Email}); err != nil {
				http.Error(w, "error creating user"+err.Error(), http.StatusInternalServerError)
				return &models.AppError{
					Code:          http.StatusInternalServerError,
					OriginalError: nil,
					Message:       "error creating user",
				}
			}
		} else {
			http.Error(w, "error getting user"+err.Error(), http.StatusInternalServerError)
			return &models.AppError{
				Code:          http.StatusInternalServerError,
				OriginalError: nil,
				Message:       "error getting user",
			}
		}
	}

	access_token := http.Cookie{
		Name:     "access_token",
		Value:    oauth2Token.AccessToken,
		Path:     "/",
		HttpOnly: true,
	}
	refresh_token := http.Cookie{
		Name:     "refresh_token",
		Value:    oauth2Token.RefreshToken,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, &access_token)
	http.SetCookie(w, &refresh_token)

	// TEST
	// resp := struct {
	// 	OAuth2Token   *oauth2.Token
	// 	IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
	// }{oauth2Token, new(json.RawMessage)}

	// if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return nil
	// }
	// data, err := json.MarshalIndent(resp, "", "    ")
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return nil
	// }
	// w.Write(data)

	http.Redirect(w, r, env.AppFullURL, http.StatusTemporaryRedirect)

	return nil
}

func (env *Env) GetTokenHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	state, err := randString(16)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return &models.AppError{
			Code:          http.StatusInternalServerError,
			OriginalError: err,
			Message:       "error calling randString",
		}
	}

	nonce, err := randString(16)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return &models.AppError{
			Code:          http.StatusInternalServerError,
			OriginalError: err,
			Message:       "error calling randString",
		}
	}
	setCallbackCookie(w, r, "state", state)
	setCallbackCookie(w, r, "nonce", nonce)

	http.Redirect(w, r, env.OAuth2Config.AuthCodeURL(state, oidc.Nonce(nonce)), http.StatusFound)

	return nil
}

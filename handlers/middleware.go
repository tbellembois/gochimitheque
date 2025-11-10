package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/locales"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/request"
	"github.com/tbellembois/gochimitheque/zmqclient"
	"golang.org/x/oauth2"
)

// AppMiddleware is the application handlers wrapper handling the "func() *models.AppError" functions.
func (env *Env) AppMiddleware(h models.AppHandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if e := h(w, r); e != nil {

			if e.OriginalError != nil {

				logger.Log.Error(e.Message + "-" + e.OriginalError.Error())

				if e.Code == http.StatusInternalServerError {
					logger.Log.Error(e.Message + "-" + e.OriginalError.Error())
				}

				http.Error(w, e.Message+"-"+e.OriginalError.Error(), e.Code)

			}

			logger.Log.Error(e.Error())
			http.Error(w, e.Error(), e.Code)

		}
	})
}

// LogingMiddleware is the application handlers wrapper logging the requests.
func (env *Env) LogingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		h.ServeHTTP(w, req)
	})
}

// HeadersMiddleware is the application handlers wrapper managing the common response HTTP headers.
func (env *Env) HeadersMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		h.ServeHTTP(w, req)
	})
}

// ContextMiddleware initialize the request context and setup initial variables.
func (env *Env) ContextMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// localization setup
		accept := r.Header.Get("Accept-Language")
		locales.Localizer = i18n.NewLocalizer(locales.Bundle, accept)

		ctx := context.WithValue(
			r.Context(),
			request.ChimithequeContextKey("container"),
			request.Container{
				AppURL:         env.AppURL,
				AppPath:        env.AppPath,
				PersonLanguage: accept,
				BuildID:        env.BuildID,
				DisableCache:   env.DisableCache,
			},
		)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FakeMiddleware fake the authentication.
func (env *Env) FakeMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// getting the request context
		ctx := r.Context()

		// getting the request container
		ctxcontainer := ctx.Value(request.ChimithequeContextKey("container"))
		container := ctxcontainer.(request.Container)

		// setting up auth person informations
		container.PersonEmail = "admin@chimitheque.fr"
		container.PersonID = 1

		ctx = context.WithValue(
			r.Context(),
			request.ChimithequeContextKey("container"),
			container,
		)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthenticateMiddleware check that a valid token is in the request.
func (env *Env) AuthenticateMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			access_token, refresh_token *http.Cookie
			err                         error
		)

		// getting the request context
		ctx := r.Context()

		if access_token, err = r.Cookie("access_token"); err != nil || access_token == nil {
			logger.Log.Debug("access token not found in cookies")
			http.Error(w, "access token not found in cookies, please log in", http.StatusUnauthorized)
			return
		}

		if refresh_token, err = r.Cookie("refresh_token"); err != nil {
			logger.Log.Debug("refresh token not found in cookies")
			http.Error(w, "refresh token not found in cookies, please log in", http.StatusUnauthorized)
			return
		}

		oauth2Token := oauth2.Token{
			AccessToken:  access_token.Value,
			RefreshToken: refresh_token.Value,
		}

		var userInfo *oidc.UserInfo
		if userInfo, err = env.OIDCProvider.UserInfo(ctx, oauth2.StaticTokenSource(&oauth2Token)); err != nil {
			// We assume that the access token not valid anymore.
			ts := env.OAuth2Config.TokenSource(ctx, &oauth2.Token{RefreshToken: refresh_token.Value})

			// Renew token.
			var newToken *oauth2.Token
			if newToken, err = ts.Token(); err != nil {
				http.Error(w, "failed to refresh token: "+err.Error(), http.StatusInternalServerError)
			} // this actually goes and renews the tokens

			if newToken != nil {
				logger.Log.WithFields(logrus.Fields{
					"newToken": newToken,
				}).Debug("AuthenticateMiddleware")

				// Save new token.
				access_token := http.Cookie{
					Name:     "access_token",
					Value:    newToken.AccessToken,
					Path:     "/",
					HttpOnly: true,
				}
				refresh_token := http.Cookie{
					Name:     "refresh_token",
					Value:    newToken.RefreshToken,
					Path:     "/",
					HttpOnly: true,
				}

				http.SetCookie(w, &access_token)
				http.SetCookie(w, &refresh_token)

				oauth2Token = *newToken
			}

			// Trying to get the user informations with the new token.
			if userInfo, err = env.OIDCProvider.UserInfo(ctx, oauth2.StaticTokenSource(&oauth2Token)); err != nil {
				http.Error(w, "failed to get userinfo: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// logger.Log.WithFields(logrus.Fields{
		// 	"userInfo": userInfo,
		// }).Debug("AuthenticateMiddleware")

		// getting the connected user
		var (
			jsonRawMessage json.RawMessage
			person         *models.Person
		)
		if jsonRawMessage, err = zmqclient.DBGetPeople("http://localhost/?search="+userInfo.Email, 1); err != nil {
			logger.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if person, err = zmqclient.ConvertDBJSONToPerson(jsonRawMessage); err != nil {
			logger.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// getting the request container
		ctxcontainer := ctx.Value(request.ChimithequeContextKey("container"))
		container := ctxcontainer.(request.Container)

		// setting up auth person informations
		container.PersonEmail = userInfo.Email
		container.PersonID = *person.PersonID

		ctx = context.WithValue(
			r.Context(),
			request.ChimithequeContextKey("container"),
			container,
		)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// func (env *Env) AuthorizeMiddleware(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// defer utils.TimeTrack(time.Now(), "AuthorizeMiddleware")
// 		h.ServeHTTP(w, r)
// 	})
// }

// AuthorizeMiddleware check that the user extracted from the JWT token by the AuthenticateMiddleware has the permissions to access the requested resource.
func (env *Env) AuthorizeMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// defer utils.TimeTrack(time.Now(), "AuthorizeMiddleware")

		var (
			personid        int64 // logged person id
			personid_string string
			item            string // storages, products...
			itemid          string // the item id to be accessed: an int, -1, -2 or ""
			action          string
			permok          bool
			err             error
		)

		logger.Log.WithFields(logrus.Fields{"r.RequestURI": fmt.Sprintf("%+v", r.RequestURI)}).Debug("AuthorizeMiddleware")

		// extracting the logged user email from context
		ctx := r.Context()
		ctxcontainer := ctx.Value(request.ChimithequeContextKey("container"))
		container := ctxcontainer.(request.Container)
		personid = container.PersonID
		personid_string = strconv.Itoa(int(personid))

		//
		// extracting request variables
		//
		vars := mux.Vars(r)

		// item = products, storages...
		item = vars["item"]
		// id = an int or ""
		itemid = vars["id"]

		logger.Log.WithFields(logrus.Fields{"vars": fmt.Sprintf("%+v", vars), "r.Method": r.Method}).Debug("AuthorizeMiddleware")

		// action = r or w
		if r.Method == "GET" {
			action = "r"
		} else if r.Method == "DELETE" {
			action = "d"
		} else {
			action = "w"
		}

		// allow/deny access
		logger.Log.WithFields(logrus.Fields{
			"itemid":          itemid,
			"item":            item,
			"personid_string": personid_string,
			"action":          action,
		}).Debug("AuthorizeMiddleware")

		if permok, err = env.Enforcer.Enforce(personid_string, action, item, itemid); err != nil {
			http.Error(w, "enforcer error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Log.WithFields(logrus.Fields{
			"permok": permok,
		}).Debug("AuthorizeMiddleware")

		if !permok {
			logger.Log.WithFields(logrus.Fields{"unauthorized": "!permok"}).Debug("AuthorizeMiddleware")
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		h.ServeHTTP(w, r)
	})
}

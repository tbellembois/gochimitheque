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
		container.PersonID = person.PersonID

		ctx = context.WithValue(
			r.Context(),
			request.ChimithequeContextKey("container"),
			container,
		)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthorizeMiddleware check that the user extracted from the JWT token by the AuthenticateMiddleware has the permissions to access the requested resource.
func (env *Env) AuthorizeMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// defer utils.TimeTrack(time.Now(), "AuthorizeMiddleware")

		var (
			personid    int    // logged person id
			item        string // storages, products...
			itemid      string // the item id to be accessed: an int, -1, -2 or ""
			itemidInt   int
			action      string
			permok      bool
			personemail string
			err         error
		)

		logger.Log.WithFields(logrus.Fields{"r.RequestURI": fmt.Sprintf("%+v", r.RequestURI)}).Debug("AuthorizeMiddleware")

		// extracting the logged user email from context
		ctx := r.Context()
		ctxcontainer := ctx.Value(request.ChimithequeContextKey("container"))
		container := ctxcontainer.(request.Container)
		personid = container.PersonID
		personemail = container.PersonEmail
		// should not be necessary
		// AuthenticateMiddleware performs a check
		if personemail == "" {
			http.Error(w, "request container personemail empty", http.StatusUnauthorized)
			return
		}

		//
		// extracting request variables
		//
		vars := mux.Vars(r)

		// item = products, storages...
		item = vars["item"]
		// id = an int or ""
		itemid = vars["id"]

		logger.Log.WithFields(logrus.Fields{"vars": fmt.Sprintf("%+v", vars)}).Debug("AuthorizeMiddleware")

		// getting the connected user
		var (
			jsonRawMessage json.RawMessage
			person         *models.Person
		)
		if jsonRawMessage, err = zmqclient.DBGetPeople("http://localhost/"+strconv.Itoa(personid), 1); err != nil {
			logger.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if person, err = zmqclient.ConvertDBJSONToPerson(jsonRawMessage); err != nil {
			logger.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// action = r or w
		if r.Method == "GET" {
			action = "r"
			itemid = "-2"
		} else {
			action = "w"
		}

		logger.Log.WithFields(logrus.Fields{"r.Method": fmt.Sprintf("%+v", r.Method)}).Debug("AuthorizeMiddleware")

		//
		// pre checks
		//
		switch r.Method {
		case "PUT":
			// REST update,create methods
			switch item {
			case "people":
				// itemid is an int
				if itemidInt, err = strconv.Atoi(itemid); err != nil {
					logger.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				// a user can not edit himself
				if itemidInt == personid {
					http.Error(w, "can not edit/delete yourself", http.StatusBadRequest)
					return
				}
			}
		case "DELETE":
			// REST delete method
			switch item {
			case "people":
				// itemid is an int
				if itemidInt, err = strconv.Atoi(itemid); err != nil {
					logger.Log.WithFields(logrus.Fields{"err": err.Error()}).Debug("AuthorizeMiddleware")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				// a user can not delete himself
				if itemidInt == personid {
					http.Error(w, "can not edit/delete yourself", http.StatusBadRequest)
					return
				}
				// we can not delete a manager
				if person.ManagedEntities == nil || len(person.ManagedEntities) != 0 {
					http.Error(w, "can not delete a manager", http.StatusBadRequest)
					return
				}
				// we can not delete an admin
				a, e := env.DB.IsPersonAdmin(itemidInt)
				if e != nil {
					logger.Log.WithFields(logrus.Fields{"e": e.Error()}).Error("AuthorizeMiddleware")
					http.Error(w, e.Error(), http.StatusInternalServerError)
					return
				}
				if a {
					http.Error(w, "can not delete an admin", http.StatusBadRequest)
					return
				}
			case "storelocations":
				// itemid is an int
				if itemidInt, err = strconv.Atoi(itemid); err != nil {
					logger.Log.WithFields(logrus.Fields{"err": err.Error()}).Debug("AuthorizeMiddleware")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				// getting the store location
				var (
					jsonRawMessage json.RawMessage
					storeLocation  *models.StoreLocation
				)
				if jsonRawMessage, err = zmqclient.DBGetStorelocations("http://localhost/store_locations/"+strconv.Itoa(int(itemidInt)), personid); err != nil {
					logger.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if storeLocation, err = zmqclient.ConvertDBJSONToStorelocation(jsonRawMessage); err != nil {
					logger.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				// can not delete store location with children
				if storeLocation.StoreLocationNbChildren != nil && *storeLocation.StoreLocationNbChildren > 0 {
					http.Error(w, "can not delete store location with children", http.StatusBadRequest)
					return
				}

				// can not delete a non empty store location
				if storeLocation.StoreLocationNbStorage != nil && *storeLocation.StoreLocationNbStorage > 0 {
					http.Error(w, "can not delete a non empty store location", http.StatusBadRequest)
					return
				}

			case "products":
				logger.Log.Debug("DELETE -> products")

				// itemid is an int
				if itemidInt, err = strconv.Atoi(itemid); err != nil {
					logger.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				logger.Log.WithFields(logrus.Fields{"itemidInt": itemidInt}).Debug("AuthorizeMiddleware")

				// getting the product
				var (
					jsonRawMessage json.RawMessage
					product        *models.Product
					count          int
				)
				if jsonRawMessage, err = zmqclient.DBGetProducts("http://localhost/"+itemid, personid); err != nil {
					logger.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return

				}

				if product, err = zmqclient.ConvertDBJSONToProduct(jsonRawMessage); err != nil {
					logger.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				logger.Log.WithFields(logrus.Fields{"product": fmt.Sprintf("%#v", product)}).Debug("AuthorizeMiddleware")

				// we can not delete a product with storages
				count = product.ProductSC
				if count != 0 {
					http.Error(w, "can not delete a product with storages", http.StatusBadRequest)
					return
				}
				logger.Log.WithFields(logrus.Fields{"count": count}).Debug("AuthorizeMiddleware")

			case "entities":
				if r.Method == "DELETE" {
					// itemid is an int
					if itemidInt, err = strconv.Atoi(itemid); err != nil {
						logger.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					// getting the entity
					var (
						jsonRawMessage json.RawMessage
						entity         *models.Entity
					)
					if jsonRawMessage, err = zmqclient.DBGetEntities("http://localhost/entities/"+strconv.Itoa(itemidInt), personid); err != nil {
						logger.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					if entity, err = zmqclient.ConvertDBJSONToEntity(jsonRawMessage); err != nil {
						logger.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					m := (entity.EntityNbPeople != nil && *entity.EntityNbPeople == 0)
					n := (entity.EntityNbStoreLocations != nil && *entity.EntityNbStoreLocations == 0)
					if m {
						http.Error(w, "can not delete an entity with members", http.StatusBadRequest)
						return
					}
					if n {
						http.Error(w, "can not delete an entity with store locations", http.StatusBadRequest)
						return
					}
				}
			}
		}

		// allow/deny access
		logger.Log.WithFields(logrus.Fields{
			"itemid":   itemid,
			"item":     item,
			"personid": strconv.Itoa(personid),
			"action":   action,
		}).Debug("AuthorizeMiddleware")

		if permok, err = env.Enforcer.Enforce(strconv.Itoa(personid), action, item, itemid); err != nil {
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

package handlers

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/globals"
	"github.com/tbellembois/gochimitheque/models"
)

// AppMiddleware is the application handlers wrapper handling the "func() *models.AppError" functions
func (env *Env) AppMiddleware(h models.AppHandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if e := h(w, r); e != nil {
			if e.Error != nil {
				globals.Log.Error(e.Message + "-" + e.Error.Error())
				if e.Code == http.StatusInternalServerError {
					globals.LogInternal.Error(e.Message + "-" + e.Error.Error())
				}
			}
			http.Error(w, e.Message, e.Code)
		}
	})
}

// LogingMiddleware is the application handlers wrapper logging the requests
func (env *Env) LogingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		h.ServeHTTP(w, req)
	})
}

// HeadersMiddleware is the application handlers wrapper managing the common response HTTP headers
func (env *Env) HeadersMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, req)
	})
}

// ContextMiddleware initialize the request context and setup initial variables
func (env *Env) ContextMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// localization setup
		accept := r.Header.Get("Accept-Language")
		globals.Localizer = i18n.NewLocalizer(globals.Bundle, accept)

		ctx := context.WithValue(
			r.Context(),
			models.ChimithequeContextKey("container"),
			models.ViewContainer{
				ProxyPath:      globals.ProxyPath,
				PersonLanguage: accept,
				BuildID:        globals.BuildID,
				DisableCache:   globals.DisableCache,
			},
		)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthenticateMiddleware check that a valid JWT token is in the request, extract and store user informations in the Go http context
func (env *Env) AuthenticateMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			email       string
			cemail      interface{}
			claims      jwt.MapClaims
			ok          bool
			person      models.Person
			err         error
			reqToken    *http.Cookie
			reqTokenStr string
			token       *jwt.Token
		)

		// token regex header version
		//tre := regexp.MustCompile("Bearer [[:alnum:]]\\.[[:alnum:]]\\.[[:alnum:]]")
		// token regex cookie version
		//tre := regexp.MustCompile("token=[[:alnum:]]\\.[[:alnum:]]\\.[[:alnum:]]")
		tre := regexp.MustCompile("token=.+")

		// extracting the token string from Authorization header
		//reqToken := r.Header.Get("Authorization")
		// extracting the token string from cookie
		if reqToken, err = r.Cookie("token"); err != nil {
			globals.Log.Debug("token not found in cookies")
			//http.Error(w, "token not found in cookies, please log in", http.StatusUnauthorized)
			http.Redirect(w, r, globals.ApplicationFullURL+"login", 307)
			return
		}
		if !tre.MatchString(reqToken.String()) {
			globals.Log.Debug("token has an invalid format")
			http.Error(w, "token has an invalid format", http.StatusUnauthorized)
			return
		}
		// header version
		//splitToken := strings.Split(reqToken, "Bearer ")
		// cookie version
		splitToken := strings.Split(reqToken.String(), "token=")
		reqTokenStr = splitToken[1]
		token, err = jwt.Parse(reqTokenStr, func(token *jwt.Token) (interface{}, error) {
			return globals.TokenSignKey, nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// getting the claims
		if claims, ok = token.Claims.(jwt.MapClaims); ok && token.Valid {
			// then the email claim
			if cemail, ok = claims["email"]; !ok {
				globals.Log.Debug("email not found in claims")
				http.Error(w, "email not found in claims", http.StatusBadRequest)
				return
			}
			email = cemail.(string)

		} else {
			globals.Log.Debug("can not extract claims")
			http.Error(w, "can not extract claims", http.StatusBadRequest)
			return
		}

		// getting the logged user
		if person, err = env.DB.GetPersonByEmail(email); err != nil {
			http.Error(w, "can not get logged user: "+err.Error(), http.StatusBadRequest)
		}

		// getting the request context
		ctx := r.Context()
		ctxcontainer := ctx.Value(models.ChimithequeContextKey("container"))
		container := ctxcontainer.(models.ViewContainer)
		// setting up auth person informations
		container.PersonEmail = person.PersonEmail
		container.PersonID = person.PersonID
		ctx = context.WithValue(
			r.Context(),
			models.ChimithequeContextKey("container"),
			container,
		)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (env *Env) getItemEntities(personid int, item, itemid string) ([]models.Entity, error) {
	var (
		e         models.Entity
		es        []models.Entity
		err       error
		itemidInt int
	)

	if itemidInt, err = strconv.Atoi(itemid); err != nil {
		return nil, err
	}

	es = make([]models.Entity, 0)

	switch item {
	case "storages":
		if e, err = env.DB.GetStorageEntity(itemidInt); err != nil {
			return nil, err
		}
		es = append(es, e)
	case "storelocations":
		if e, err = env.DB.GetStoreLocationEntity(itemidInt); err != nil {
			return nil, err
		}
		es = append(es, e)
	case "people":
		if es, err = env.DB.GetPersonEntities(personid, itemidInt); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unexpecter permission item")
	}

	return es, nil
}

// AuthorizeMiddleware check that the user extracted from the JWT token by the AuthenticateMiddleware has the permissions to access the requested resource
func (env *Env) AuthorizeMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//defer utils.TimeTrack(time.Now(), "AuthorizeMiddleware")

		var (
			personid    int    // logged person id
			item        string // storages, products...
			itemid      string // the item id to be accessed: an int, -1, -2 or ""
			itemidInt   int
			view        string
			action      string
			permok      bool
			personemail string
			err         error
		)

		// extracting the logged user email from context
		ctx := r.Context()
		ctxcontainer := ctx.Value(models.ChimithequeContextKey("container"))
		container := ctxcontainer.(models.ViewContainer)
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
		// view = v or vc or ""
		view = vars["view"]
		// item = products, storages...
		item = vars["item"]
		// id = an int or ""
		itemid = vars["id"]

		// action = r or w
		if view == "v" {
			action = "r"
		} else if view == "vc" {
			action = "w"
		} else {
			if r.Method == "GET" {
				action = "r"
			} else {
				action = "w"
			}
		}
		globals.Log.WithFields(logrus.Fields{
			"itemid":      itemid,
			"item":        item,
			"view":        view,
			"personemail": personemail,
			"personid":    personid,
			"action":      action}).Debug("AuthorizeMiddleware")

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
					globals.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				// a user can not edit himself
				if itemidInt == personid {
					http.Error(w, "can not edit/delete yourself", http.StatusUnauthorized)
					return
				}
				// we can not edit an admin
				a, e := env.DB.IsPersonAdmin(itemidInt)
				if e != nil {
					globals.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
					http.Error(w, e.Error(), http.StatusInternalServerError)
					return
				}
				if a {
					http.Error(w, "can not delete an admin", http.StatusUnauthorized)
					return
				}
			}
		case "DELETE":
			// REST delete method
			switch item {
			case "people":
				// itemid is an int
				if itemidInt, err = strconv.Atoi(itemid); err != nil {
					globals.Log.WithFields(logrus.Fields{"err": err.Error()}).Debug("AuthorizeMiddleware")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				// a user can not delete himself
				if itemidInt == personid {
					http.Error(w, "can not edit/delete yourself", http.StatusUnauthorized)
					return
				}
				// we can not delete a manager
				m, e := env.DB.IsPersonManager(itemidInt)
				if e != nil {
					globals.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
					http.Error(w, e.Error(), http.StatusInternalServerError)
					return
				}
				if m {
					http.Error(w, "can not delete a manager", http.StatusUnauthorized)
					return
				}
				// we can not delete an admin
				a, e := env.DB.IsPersonAdmin(itemidInt)
				if e != nil {
					globals.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
					http.Error(w, e.Error(), http.StatusInternalServerError)
					return
				}
				if a {
					http.Error(w, "can not delete an admin", http.StatusUnauthorized)
					return
				}
			case "storelocations":
				// itemid is an int
				if itemidInt, err = strconv.Atoi(itemid); err != nil {
					globals.Log.WithFields(logrus.Fields{"err": err.Error()}).Debug("AuthorizeMiddleware")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				// we can not delete a non empty store location
				m, e := env.DB.IsStoreLocationEmpty(itemidInt)
				if e != nil {
					globals.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
					http.Error(w, e.Error(), http.StatusInternalServerError)
					return
				}
				if !m {
					http.Error(w, "can not delete a non empty store location", http.StatusUnauthorized)
					return
				}
			case "products":
				// itemid is an int
				if itemidInt, err = strconv.Atoi(itemid); err != nil {
					globals.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				// we can not delete a product with storages
				c, e := env.DB.CountProductStorages(itemidInt)
				if e != nil {
					globals.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
					http.Error(w, e.Error(), http.StatusInternalServerError)
					return
				}
				if c != 0 {
					http.Error(w, "can not delete a product with storages", http.StatusUnauthorized)
					return
				}
			case "entities":
				if r.Method == "DELETE" {
					// itemid is an int
					if itemidInt, err = strconv.Atoi(itemid); err != nil {
						globals.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					m, e1 := env.DB.IsEntityEmpty(itemidInt)
					n, e2 := env.DB.HasEntityNoStorelocation(itemidInt)
					if e1 != nil {
						globals.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
						http.Error(w, e1.Error(), http.StatusUnauthorized)
						return
					}
					if e2 != nil {
						globals.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("AuthorizeMiddleware")
						http.Error(w, e2.Error(), http.StatusUnauthorized)
						return
					}
					if !m {
						http.Error(w, "can not delete a non empty entity", http.StatusUnauthorized)
						return
					}
					if !n {
						http.Error(w, "can not delete an entity with store locations", http.StatusUnauthorized)
						return
					}
				}
			}
		}

		// allow/deny access
		globals.Log.WithFields(logrus.Fields{
			"itemid":   itemid,
			"item":     item,
			"personid": strconv.Itoa(personid),
			"action":   action}).Debug("AuthorizeMiddleware")

		if permok, err = globals.Enforcer.Enforce(strconv.Itoa(personid), action, item, itemid); err != nil {
			http.Error(w, "enforcer error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if !permok {
			globals.Log.WithFields(logrus.Fields{"unauthorized": "!permok"}).Debug("AuthorizeMiddleware")
			if r.RequestURI == globals.ProxyPath || r.RequestURI == "" {
				// redirect on login page for the root of the application
				http.Redirect(w, r, globals.ApplicationFullURL+"login", 307)
			} else {
				// else sending a 403
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}

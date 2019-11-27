package handlers

import (
	"context"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/global"
	"github.com/tbellembois/gochimitheque/helpers"
	"github.com/tbellembois/gochimitheque/models"
)

// AppMiddleware is the application handlers wrapper handling the "func() *models.AppError" functions
func (env *Env) AppMiddleware(h models.AppHandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if e := h(w, r); e != nil {
			if e.Error != nil {
				log.Error(e.Message + "-" + e.Error.Error())
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
		global.Localizer = i18n.NewLocalizer(global.Bundle, accept)

		ctx := context.WithValue(
			r.Context(),
			global.ChimithequeContextKey("container"),
			helpers.ViewContainer{
				ProxyPath:      global.ProxyPath,
				PersonLanguage: accept,
				BuildID:        global.BuildID,
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
			log.Debug("token not found in cookies")
			//http.Error(w, "token not found in cookies, please log in", http.StatusUnauthorized)
			http.Redirect(w, r, global.ApplicationFullURL+"login", 307)
			return
		}
		if !tre.MatchString(reqToken.String()) {
			log.Debug("token has an invalid format")
			http.Error(w, "token has an invalid format", http.StatusUnauthorized)
			return
		}
		// header version
		//splitToken := strings.Split(reqToken, "Bearer ")
		// cookie version
		splitToken := strings.Split(reqToken.String(), "token=")
		reqTokenStr = splitToken[1]
		token, err = jwt.Parse(reqTokenStr, func(token *jwt.Token) (interface{}, error) {
			return global.TokenSignKey, nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// getting the claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// then the email claim
			if cemail, ok := claims["email"]; !ok {
				log.Debug("email not found in claims")
				http.Error(w, "email not found in claims", http.StatusBadRequest)
				return
			} else {
				email = cemail.(string)
			}
		} else {
			log.Debug("can not extract claims")
			http.Error(w, "can not extract claims", http.StatusBadRequest)
			return
		}

		// getting the logged user
		if person, err = env.DB.GetPersonByEmail(email); err != nil {
			http.Error(w, "can not get logged user: "+err.Error(), http.StatusBadRequest)
		}

		// getting the request context
		ctx := r.Context()
		ctxcontainer := ctx.Value(global.ChimithequeContextKey("container"))
		container := ctxcontainer.(helpers.ViewContainer)
		// setting up auth person informations
		container.PersonEmail = person.PersonEmail
		container.PersonID = person.PersonID
		ctx = context.WithValue(
			r.Context(),
			global.ChimithequeContextKey("container"),
			container,
		)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthorizeMiddleware check that the user extracted from the JWT token by the AuthenticateMiddleware has the permissions to access the requested resource
func (env *Env) AuthorizeMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			personid    int
			personemail string
			itemid      int
			perm        string
			permok      bool
			err         error
		)

		// extracting the logged user email from context
		ctx := r.Context()
		ctxcontainer := ctx.Value(global.ChimithequeContextKey("container"))
		container := ctxcontainer.(helpers.ViewContainer)
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
		view := vars["view"]
		item := vars["item"]
		id := vars["id"]
		log.WithFields(log.Fields{
			"id":          id,
			"item":        item,
			"view":        view,
			"personemail": personemail,
			"personid":    personid,
			"r.Method":    r.Method}).Debug("AuthorizeMiddleware")

		//
		// id and item translations and setup for the HasPersonPermission method, and some bypasses
		//
		switch item {
		case "peoplepass", "peoplep", "bookmarks", "delete-token", "borrowings", "download", "validate", "format":
			// everybody can change his password
			// everybody can bookmark a product
			// everybody can borrow a storage
			// everybody can logout
			// everybody can download an export
			// everybody can validate an item
			// everybody can format an item
			h.ServeHTTP(w, r)
			return
		case "welcomeannounce":
			// welcome announcements are editable by admins only
			item = "entities"
			id = "-1"
		case "stocks":
			// to access stocks, one needs permission on at least one storage
			item = "storages"
			id = "-2"
		case "storages":
			if id != "-1" && id != "-2" && id != "" {
				// storages access are per entity
				// extracting entity id from storage
				if itemid, err = strconv.Atoi(id); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				s, e := env.DB.GetStorage(itemid)
				if e != nil {
					http.Error(w, e.Error(), http.StatusInternalServerError)
					return
				}
				id = strconv.Itoa(s.EntityID)
			}
		}

		if id == "" {
			itemid = -2
		} else {
			if itemid, err = strconv.Atoi(id); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		//
		// perm variable setup for the HasPersonPermission method
		//
		switch r.Method {
		case "GET":
			switch view {
			// "v": view methods
			// "": REST get method
			case "v", "":
				perm = "r"
			// views update, create
			case "vu", "vc":
				perm = "w"
			}
		case "PUT":
			// REST update,create methods
			switch item {
			case "people":
				// a user can not edit himself
				if itemid == personid {
					http.Error(w, "can not edit/delete yourself", http.StatusForbidden)
					return
				}
			}

			perm = "w"
		case "POST":
			switch item {
			case "storages":
				if id != "-1" && id != "-2" && id != "" {
					// checking that the connected person
					// can "w" "storages" in the entity of the storage's
					// store location
					var (
						s models.Storage
						e error
					)
					if e = r.ParseForm(); err != nil {
						http.Error(w, e.Error(), http.StatusInternalServerError)
					}
					if e = global.Decoder.Decode(&s, r.PostForm); err != nil {
						http.Error(w, e.Error(), http.StatusInternalServerError)
					}
					if s.StoreLocation, e = env.DB.GetStoreLocation(int(s.StoreLocationID.Int64)); err != nil {
						http.Error(w, e.Error(), http.StatusInternalServerError)
					}
					itemid = s.StoreLocation.EntityID
				}
			}
			perm = "w"
		case "DELETE":
			// REST delete method
			switch item {
			case "people":
				// a user can not delete himself
				if itemid == personid {
					http.Error(w, "can not edit/delete yourself", http.StatusForbidden)
					return
				}
				// we can not delete a manager
				m, e := env.DB.IsPersonManager(itemid)
				if e != nil {
					http.Error(w, e.Error(), http.StatusInternalServerError)
					return
				}
				if m {
					http.Error(w, "can not delete a manager", http.StatusForbidden)
					return
				}
			case "storelocations":
				// we can not delete a non empty store location
				m, e := env.DB.IsStoreLocationEmpty(itemid)
				if e != nil {
					http.Error(w, e.Error(), http.StatusInternalServerError)
					return
				}
				if !m {
					http.Error(w, "can not delete a non empty store location", http.StatusBadRequest)
					return
				}
			case "products":
				c, e := env.DB.CountProductStorages(itemid)
				if e != nil {
					http.Error(w, e.Error(), http.StatusInternalServerError)
					return
				}
				if c != 0 {
					http.Error(w, "can not delete a product with storages", http.StatusBadRequest)
					return
				}
			case "entities":
				if r.Method == "DELETE" {
					m, e1 := env.DB.IsEntityEmpty(itemid)
					n, e2 := env.DB.HasEntityNoStorelocation(itemid)
					if e1 != nil {
						http.Error(w, e1.Error(), http.StatusInternalServerError)
						return
					}
					if e2 != nil {
						http.Error(w, e2.Error(), http.StatusInternalServerError)
						return
					}
					if !m {
						http.Error(w, "can not delete a non empty entity", http.StatusBadRequest)
						return
					}
					if !n {
						http.Error(w, "can not delete an entity with store locations", http.StatusBadRequest)
						return
					}
				}
			}

			perm = "w"
		default:
			log.Debug("unsupported http verb")
			http.Error(w, "unsupported http verb", http.StatusBadRequest)
			return
		}

		// allow/deny access
		if permok, err = env.DB.HasPersonPermission(personid, perm, item, itemid); err != nil {
			log.WithFields(log.Fields{"unauthorized": err.Error()}).Debug("AuthorizeMiddleware")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !permok {
			log.WithFields(log.Fields{"unauthorized": "!permok"}).Debug("AuthorizeMiddleware")
			if r.RequestURI == global.ProxyPath || r.RequestURI == "" {
				// redirect on login page for the root of the application
				http.Redirect(w, r, global.ApplicationFullURL+"login", 307)
			} else {
				// else sending a 403
				http.Error(w, "forbidden", http.StatusForbidden)
			}
			return
		}

		h.ServeHTTP(w, r)
	})
}

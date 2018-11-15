package handlers

import (
	"context"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/helpers"
	"github.com/tbellembois/gochimitheque/models"
)

// AppMiddleware is the application handlers wrapper handling the "func() *models.AppError" functions
func (env *Env) AppMiddleware(h models.AppHandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if e := h(w, r); e != nil {
			log.Error(e.Message + "-" + e.Error.Error())
			http.Error(w, e.Message, e.Code)
		}
	})
}

// LogingMiddleware is the application handlers wrapper logging the requests
func (env *Env) LogingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Debug(req.RequestURI)
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

// AuthenticateMiddleware check that a valid JWT token is in the request, extract and store user informations in the Go http context
func (env *Env) AuthenticateMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			email  string
			person models.Person
			//permissions []models.Permission
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
			http.Error(w, "token not found in cookies, please log in", http.StatusUnauthorized)
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
			return TokenSignKey, nil
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
			http.Error(w, "can not get logged user", http.StatusBadRequest)
		}

		// getting the logged user permissions
		// if permissions, err = env.DB.GetPersonPermissions(person.PersonID); err != nil {
		// 	http.Error(w, "can not get logged user permissions", http.StatusBadRequest)
		// }

		ctx := context.WithValue(r.Context(), "container", helpers.ViewContainer{PersonEmail: person.PersonEmail, PersonID: person.PersonID, ProxyPath: env.ProxyPath})

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
		ctxcontainer := ctx.Value("container")
		container := ctxcontainer.(helpers.ViewContainer)
		personid = container.PersonID
		personemail = container.PersonEmail
		// should not be necessary
		// AuthenticateMiddleware performs a check
		if personemail == "" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// extraction request variables
		vars := mux.Vars(r)
		item := vars["item"]
		view := vars["view"]
		id := vars["id"]
		log.WithFields(log.Fields{"id": id, "item": item, "view": view, "personemail": personemail}).Debug("AuthorizeMiddleware")

		if id == "" {
			itemid = -2
		} else {
			if itemid, err = strconv.Atoi(id); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// depending on the request
		// preparing the HasPersonPermission parameters
		// to allow/deny access
		switch r.Method {
		case "GET":
			switch view {
			// view (list) or
			// REST view methods
			case "v", "":
				perm = "r"
			// views update, create
			case "vu", "vc":
				perm = "w"
			}
		case "POST", "PUT", "DELETE":
			//REST update, delete, create methods
			// TODO: we need to perform more permission check here
			// to ensure that values in the request body
			// are allowed for the auth user
			// ex: can the auth user set foo@bar.com in entity 3?
			switch item {
			case "people":
				// a user can not edit/delete himself
				if itemid == personid {
					http.Error(w, "can not edit/delete yourself", http.StatusForbidden)
				}
				// we can not delete a manager
				if r.Method == "DELETE" {
					m, e := env.DB.IsPersonManager(itemid)
					if e != nil {
						http.Error(w, e.Error(), http.StatusInternalServerError)
					}
					if m {
						http.Error(w, "can not delete a manager", http.StatusForbidden)
					}
				}
			case "storelocations":
				// TODO: we can not delete a non empty store location
			case "products":
				// TODO: we can not delete a stored product
			case "storages":

			case "entities":
				// we can not delete a non empty entity
				if r.Method == "DELETE" {
					m, e := env.DB.IsEntityEmpty(itemid)
					if e != nil {
						http.Error(w, e.Error(), http.StatusInternalServerError)
					}
					if !m {
						http.Error(w, "can not delete a non empty entity", http.StatusForbidden)
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
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		h.ServeHTTP(w, r)
	})
}

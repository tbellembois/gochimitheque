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
		//defer helpers.TimeTrack(time.Now(), "AuthorizeMiddleware")

		var (
			// HasPersonPermission parameters
			personid int    // logged person id
			perm     string // r, w or all
			item     string // storages, products...
			eids     []int  // entity ids

			itemid      string // the item id to be accessed: an int, -1, -2 or ""
			itemidInt   int
			view        string
			permok      bool
			personemail string
			err         error
			es          []models.Entity

			permvalue helpers.PermValue
			ok        bool
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
		// view = v or vc or ""
		view = vars["view"]
		// item = products, storages...
		item = vars["item"]
		// id = an int or ""
		itemid = vars["id"]
		log.WithFields(log.Fields{
			"itemid":      itemid,
			"item":        item,
			"view":        view,
			"personemail": personemail,
			"personid":    personid,
			"r.Method":    r.Method}).Debug("AuthorizeMiddleware")

		// full access items
		if item == "peoplepass" || item == "peoplep" || item == "bookmarks" || item == "delete-token" || item == "borrowings" || item == "download" || item == "validate" || item == "format" {
			h.ServeHTTP(w, r)
			return
		}

		// building the PermKey
		permkey := helpers.PermKey{View: view, Item: item, Verb: r.Method}
		log.WithFields(log.Fields{"permkey": permkey}).Debug("AuthorizeMiddleware")
		switch itemid {
		case "-1":
			permkey.Id = "-1"
		case "-2":
			permkey.Id = "-2"
		case "":
			permkey.Id = ""
		default:
			permkey.Id = "id"
		}

		// getting the permission definition in the PermMatrix
		if permvalue, ok = helpers.PermMatrix[permkey]; !ok {
			log.Error("key not found in PermMatrix")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.WithFields(log.Fields{"permvalue": permvalue}).Debug("AuthorizeMiddleware")

		// table translation
		if permvalue.Item != "" {
			item = permvalue.Item
		}
		if permvalue.Id != "" {
			itemid = permvalue.Id
		}

		// building the HasPersonPermission method parameters
		switch permvalue.Type {
		case "r":
			perm = "r"
			// itemid is -1 -2 or and int
			if itemidInt, err = strconv.Atoi(itemid); err != nil {
				log.WithFields(log.Fields{"err": err.Error()}).Debug("AuthorizeMiddleware")
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			eids = []int{itemidInt}
		case "w":
			perm = "w"
			// itemid is -1 -2 or and int
			if itemidInt, err = strconv.Atoi(itemid); err != nil {
				log.WithFields(log.Fields{"err": err.Error()}).Debug("AuthorizeMiddleware")
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			eids = []int{itemidInt}
		case "rall":
			perm = "r"
			eids = []int{-1}
		case "wall":
			perm = "w"
			eids = []int{-1}
		case "rany":
			perm = "r"
			eids = []int{-2}
		case "wany":
			perm = "w"
			eids = []int{-2}
		case "rent":
			perm = "r"
			// itemid is an int
			// item is storages, storelocations or people (after table translation)
			if es, err = env.getItemEntities(personid, item, itemid); err != nil {
				log.WithFields(log.Fields{"err": err.Error()}).Debug("AuthorizeMiddleware")
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			for _, e := range es {
				eids = append(eids, e.EntityID)
			}
		case "went":
			perm = "w"
			// itemid is an int
			// item is storages, storelocations or people (after table translation)
			if es, err = env.getItemEntities(personid, item, itemid); err != nil {
				log.WithFields(log.Fields{"err": err.Error()}).Debug("AuthorizeMiddleware")
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			for _, e := range es {
				eids = append(eids, e.EntityID)
			}
		default:
			log.WithFields(log.Fields{"err": err.Error()}).Debug("AuthorizeMiddleware")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// developper check
		if len(eids) == 0 {
			log.Error("eids empty")
			http.Error(w, "eids empty", http.StatusInternalServerError)
			return
		}

		// allow/deny access
		if permok, err = env.DB.HasPersonPermission(personid, perm, item, eids); err != nil {
			log.WithFields(log.Fields{"err": err.Error()}).Debug("AuthorizeMiddleware")
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

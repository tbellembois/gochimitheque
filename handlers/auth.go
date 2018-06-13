package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/models"
)

var TokenSignKey = []byte("secret")

// GetTokenHandler authenticate the user and return a token on success
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

func (env *Env) AppMiddleware(h models.AppHandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if e := h(w, r); e != nil {
			log.Error(e.Message + "-" + e.Error.Error())
			http.Error(w, e.Message, e.Code)
		}
	})
}

func (env *Env) LogingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Debug(req.RequestURI)
		h.ServeHTTP(w, req)
	})
}

func (env *Env) HeadersMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, req)
	})
}

func (env *Env) AuthenticateMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			email       string
			person      models.Person
			permissions []models.Permission
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
		if permissions, err = env.DB.GetPersonPermissions(person.PersonID); err != nil {
			http.Error(w, "can not get logged user permissions", http.StatusBadRequest)
		}

		ctx := context.WithValue(r.Context(), "container", models.ViewContainer{PersonEmail: person.PersonEmail, PersonID: person.PersonID, Permissions: permissions})

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

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
		container := ctxcontainer.(models.ViewContainer)
		personid = container.PersonID
		personemail = container.PersonEmail
		// should not be necessary
		// AuthenticateMiddleware performs a check
		if personemail == "" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)
		item := vars["item"]
		view := vars["view"]
		id := vars["id"]
		log.WithFields(log.Fields{"id": id, "item": item, "view": view, "personemail": personemail}).Debug("AuthorizeMiddleware")

		if id == "" {
			itemid = -1
		} else {
			if itemid, err = strconv.Atoi(id); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		switch r.Method {
		case "GET":
			switch view {
			// view (list)
			case "v", "":
				perm = "r"
			// update, create
			case "vu", "vc":
				perm = "w"
			}
		case "POST", "PUT", "DELETE":
			perm = "w"
		default:
			log.Debug("unsupported http verb")
			http.Error(w, "unsupported http verb", http.StatusBadRequest)
			return
		}

		if permok, err = env.DB.HasPermission(personid, perm, item, itemid); err != nil {
			log.Debug("unauthorized:" + err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !permok {
			log.Debug("unauthorized:" + perm + ":" + item + ":" + id)
			http.Error(w, perm+":"+item+":"+id, http.StatusForbidden)
			return
		}

		h.ServeHTTP(w, r)
	})
}

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

	if p, err = env.DB.HasPermission(personid, perm, item, itemid); err != nil {
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

//haspermission/{personid}/{perm}/{item}/{itemid}

package main

import (
	"flag"
	"html/template"
	"net/http"
	"os"

	"github.com/GeertJohan/go.rice"
	"github.com/Joker/jade"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/handlers"
	"github.com/tbellembois/gochimitheque/models"
)

func main() {

	var (
		err       error
		logf      *os.File
		dbname    = "./storage.db"
		datastore *models.SQLiteDataStore
	)

	// getting the program parameters
	listenPort := flag.String("port", "8081", "the port to listen")
	logfile := flag.String("logfile", "", "log to the given file")
	debug := flag.Bool("debug", false, "debug (verbose log), default is error")
	flag.Parse()

	// logging to file if logfile parameter specified
	if *logfile != "" {
		if logf, err = os.OpenFile(*logfile, os.O_WRONLY|os.O_CREATE, 0755); err != nil {
			log.Panic(err)
		} else {
			log.SetOutput(logf)
		}
	}

	// setting the log level
	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
	}

	// database initialization
	if datastore, err = models.NewDBstore(dbname); err != nil {
		log.Panic(err)
	}
	if err = datastore.CreateDatabase(); err != nil {
		log.Panic(err)
	}

	// environment creation
	env := handlers.Env{
		DB:        datastore,
		Templates: make(map[string]*template.Template),
	}

	// HasPermission used by template rendering to show/hide html elements
	funcMap := template.FuncMap{
		"HasPermission": func(id int, perm string, item string, itemid int) bool {
			p, _ := env.DB.HasPersonPermission(id, perm, item, itemid)
			return p
		},
	}

	// template compilation
	// home
	hometmpl, e := jade.ParseFile("static/templates/home/index.jade")
	if e != nil {
		log.Fatal("hometmpl jade:" + e.Error())
	}
	env.Templates["home"], err = template.New("home").Funcs(funcMap).Parse(hometmpl)
	if err != nil {
		log.Fatal("hometmpl parse:" + e.Error())
	}
	// login
	logintmpl, e := jade.ParseFile("static/templates/login/index.jade")
	if e != nil {
		log.Fatal("logintmpl jade:" + e.Error())
	}
	env.Templates["login"], err = template.New("login").Funcs(funcMap).Parse(logintmpl)
	if err != nil {
		log.Fatal("logintmpl parse:" + e.Error())
	}
	// entity
	entityindextmpl, e := jade.ParseFile("static/templates/entity/index.jade")
	if e != nil {
		log.Fatal("entityindextmpl jade:" + e.Error())
	}
	env.Templates["entityindex"], err = template.New("entityindex").Funcs(funcMap).Parse(entityindextmpl)
	if err != nil {
		log.Fatal("entityindextmpl parse:" + e.Error())
	}
	entitycreatetmpl, e := jade.ParseFile("static/templates/entity/create.jade")
	if e != nil {
		log.Fatal("entitycreatetmpl jade:" + e.Error())
	}
	env.Templates["entitycreate"], err = template.New("entitycreate").Funcs(funcMap).Parse(entitycreatetmpl)
	if err != nil {
		log.Fatal("entitycreatetmpl parse:" + e.Error())
	}
	// person
	personindextmpl, e := jade.ParseFile("static/templates/person/index.jade")
	if e != nil {
		log.Fatal("personindextmpl jade:" + e.Error())
	}
	env.Templates["personindex"], err = template.New("personindex").Funcs(funcMap).Parse(personindextmpl)
	if err != nil {
		log.Fatal("personindextmpl parse:" + e.Error())
	}

	// router definition
	r := mux.NewRouter()
	commonChain := alice.New(env.HeadersMiddleware, env.LogingMiddleware)
	securechain := alice.New(env.HeadersMiddleware, env.LogingMiddleware, env.AuthenticateMiddleware, env.AuthorizeMiddleware)
	// login
	r.Handle("/login", commonChain.Then(env.AppMiddleware(env.VLoginHandler))).Methods("GET")
	r.Handle("/get-token", commonChain.Then(env.AppMiddleware(env.GetTokenHandler))).Methods("POST")
	// developper tests
	r.Handle("/v/test", commonChain.Then(env.AppMiddleware(env.VTestHandler))).Methods("GET")
	// home page
	r.Handle("/", securechain.Then(env.AppMiddleware(env.HomeHandler))).Methods("GET")
	// entities
	r.Handle("/{view:v}/{item:entities}", securechain.Then(env.AppMiddleware(env.VGetEntitiesHandler))).Methods("GET")
	r.Handle("/{view:vc}/{item:entity}", securechain.Then(env.AppMiddleware(env.VCreateEntityHandler))).Methods("GET")
	r.Handle("/{item:entities}", securechain.Then(env.AppMiddleware(env.GetEntitiesHandler))).Methods("GET")
	r.Handle("/{item:entity}/{id}", securechain.Then(env.AppMiddleware(env.GetEntityHandler))).Methods("GET")
	r.Handle("/{item:entity}", securechain.Then(env.AppMiddleware(env.CreateEntityHandler))).Methods("POST")
	r.Handle("/{item:entity}/{id}", securechain.Then(env.AppMiddleware(env.UpdateEntityHandler))).Methods("PUT")
	r.Handle("/{item:entity}/{id}", securechain.Then(env.AppMiddleware(env.DeleteEntityHandler))).Methods("DELETE")
	// people
	r.Handle("/{view:v}/{item:people}", securechain.Then(env.AppMiddleware(env.VGetPeopleHandler))).Methods("GET")
	r.Handle("/{item:people}", securechain.Then(env.AppMiddleware(env.GetPeopleHandler))).Methods("GET")
	r.Handle("/{item:person}/{id}/entities", securechain.Then(env.AppMiddleware(env.GetPersonEntitiesHandler))).Methods("GET")
	// entity name validation
	r.Handle("/validate/entity/{id}/name/{name}", commonChain.Then(env.AppMiddleware(env.ValidateEntityNameHandler))).Methods("GET")
	// permissions checker
	r.Handle("/haspermission/{personid}/{perm}/{item}/{itemid}", commonChain.Then(env.AppMiddleware(env.HasPermissionHandler))).Methods("GET")

	// rice boxes
	cssBox := rice.MustFindBox("static/css")
	cssFileServer := http.StripPrefix("/css/", http.FileServer(cssBox.HTTPBox()))
	http.Handle("/css/", cssFileServer)

	jsBox := rice.MustFindBox("static/js")
	jsFileServer := http.StripPrefix("/js/", http.FileServer(jsBox.HTTPBox()))
	http.Handle("/js/", jsFileServer)

	imgBox := rice.MustFindBox("static/img")
	imgFileServer := http.StripPrefix("/img/", http.FileServer(imgBox.HTTPBox()))
	http.Handle("/img/", imgFileServer)

	http.Handle("/", r)

	if err = http.ListenAndServe(":"+*listenPort, nil); err != nil {
		panic(err)
	}
}

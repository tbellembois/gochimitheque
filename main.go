package main

//go:generate go run gogenerate/localejs.go
//go:generate gopherjs build gopherjs/gjs-common.go -o static/js/chim/gjs-common.js -m
//go:generate rice embed-go
//go:generate jade -writer -basedir static/templates -d ./jade welcomeannounce/index.jade home/index.jade login/index.jade entity/index.jade entity/create.jade product/index.jade product/create.jade storage/index.jade storage/create.jade storelocation/index.jade storelocation/create.jade person/index.jade person/create.jade person/pupdate.jade test.jade

// build with
//go build -ldflags "-X main.BuildID=$(git tag | head -1)" -o gochimitheque

import (
	"database/sql"
	"flag"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/global"
	"github.com/tbellembois/gochimitheque/handlers"
	"github.com/tbellembois/gochimitheque/models"
)

var (
	BuildID string // BuildID is compile-time variable
)

func main() {

	var (
		err       error
		logf      *os.File
		dbname    = "./storage.db"
		datastore models.Datastore
	)

	// getting the program parameters
	listenport := flag.String("port", "8081", "the port to listen")
	proxyurl := flag.String("proxyurl", "http://localhost:"+*listenport, "the application url (without the path) if behind a proxy, with NO trailing /")
	proxypath := flag.String("proxypath", "/", "the application path if behind a proxy, with the heading and trailing /")
	mailServerAddress := flag.String("mailserveraddress", "", "the mail server address")
	mailServerPort := flag.String("mailserverport", "", "the mail server address")
	mailServerSender := flag.String("mailserversender", "", "the mail server sender")
	mailServerUseTLS := flag.Bool("mailserverusetls", false, "use TLS? (optional)")
	mailServerTLSSkipVerify := flag.Bool("mailservertlsskipverify", false, "skip TLS verification? (optional)")
	admins := flag.String("admins", "", "the additional admins (comma separated email adresses)")
	logfile := flag.String("logfile", "", "log to the given file")
	debug := flag.Bool("debug", false, "debug (verbose log), default is error")
	importfrom := flag.String("importfrom", "", "full path of the directory containing the CSV to import")
	flag.Parse()

	// setting the log level
	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// logging to file if logfile parameter specified
	if *logfile != "" {
		if logf, err = os.OpenFile(*logfile, os.O_WRONLY|os.O_CREATE, 0755); err != nil {
			log.Fatal(err)
		} else {
			log.SetOutput(logf)
		}
	}

	// global variables init
	global.BuildID = BuildID
	global.ProxyPath = *proxypath
	global.ProxyURL = *proxyurl
	global.MailServerAddress = *mailServerAddress
	global.MailServerSender = *mailServerSender
	global.MailServerPort = *mailServerPort
	global.MailServerUseTLS = *mailServerUseTLS
	global.MailServerTLSSkipVerify = *mailServerTLSSkipVerify
	log.Info("- application version: " + global.BuildID)
	log.Info("- application endpoint: " + global.ProxyURL + global.ProxyPath)

	// database initialization
	log.Info("- opening database connection to " + dbname)
	if datastore, err = models.NewSQLiteDBstore(dbname); err != nil {
		log.Fatal(err)
	}
	log.Info("- creating database if needed")
	if err = datastore.CreateDatabase(); err != nil {
		log.Fatal(err)
	}
	if *importfrom != "" {
		log.Info("- import from csv into database")
		err := datastore.Import(*importfrom)
		if err != nil {
			log.Error("an error occured: " + err.Error())
		}
		os.Exit(0)
	}

	// adding additional admins
	var (
		p             models.Person
		formerAdmins  []models.Person
		currentAdmins []string
		isStillAdmin  bool
	)
	if *admins != "" {
		currentAdmins = strings.Split(*admins, ",")
	}
	if formerAdmins, err = datastore.GetAdmins(); err != nil {
		log.Fatal(err)
	}
	// cleaning former admins
	for _, fa := range formerAdmins {
		isStillAdmin = false
		log.Info("former admin: " + fa.PersonEmail)
		for _, ca := range currentAdmins {
			if ca == fa.PersonEmail {
				isStillAdmin = true
			}
		}
		if !isStillAdmin {
			log.Info(fa.PersonEmail + " is not an admin anymore, removing permissions")
			if err = datastore.UnsetPersonAdmin(fa.PersonID); err != nil {
				log.Fatal(err)
			}
		}
	}
	// setting up new ones
	if len(currentAdmins) > 0 {
		for _, ca := range currentAdmins {
			log.Info("additional admin: " + ca)
			if p, err = datastore.GetPersonByEmail(ca); err != nil {
				if err == sql.ErrNoRows {
					log.Fatal("user " + ca + " not found in database")
				} else {
					log.Fatal(err)
				}
			}

			datastore.SetPersonAdmin(p.PersonID)
		}
	}

	// environment creation
	env := handlers.Env{
		DB: datastore,
	}

	// router definition
	r := mux.NewRouter()

	// add the pprof routes
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	r.Handle("/debug/pprof/block", pprof.Handler("block"))
	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))

	commonChain := alice.New(env.ContextMiddleware, env.HeadersMiddleware, env.LogingMiddleware)
	securechain := alice.New(env.ContextMiddleware, env.HeadersMiddleware, env.LogingMiddleware, env.AuthenticateMiddleware, env.AuthorizeMiddleware)

	// login
	r.Handle("/login", commonChain.Then(env.AppMiddleware(env.VLoginHandler))).Methods("GET")
	r.Handle("/get-token", commonChain.Then(env.AppMiddleware(env.GetTokenHandler))).Methods("POST")
	r.Handle("/{item:delete-token}", securechain.Then(env.AppMiddleware(env.DeleteTokenHandler))).Methods("GET")
	r.Handle("/reset-password", commonChain.Then(env.AppMiddleware(env.ResetPasswordHandler))).Methods("POST")
	r.Handle("/reset", commonChain.Then(env.AppMiddleware(env.ResetHandler))).Methods("GET")
	r.Handle("/captcha", commonChain.Then(env.AppMiddleware(env.CaptchaHandler))).Methods("GET")

	// developper tests
	r.Handle("/v/test", securechain.Then(env.AppMiddleware(env.VTestHandler))).Methods("GET")
	// home page
	r.Handle("/", securechain.Then(env.AppMiddleware(env.HomeHandler))).Methods("GET")
	// welcome announce
	r.Handle("/{view:v}/{item:welcomeannounce}", securechain.Then(env.AppMiddleware(env.VWelcomeAnnounceHandler))).Methods("GET")
	r.Handle("/{item:welcomeannounce}", securechain.Then(env.AppMiddleware(env.UpdateWelcomeAnnounceHandler))).Methods("PUT")
	r.Handle("/{item:welcomeannounce}", commonChain.Then(env.AppMiddleware(env.GetWelcomeAnnounceHandler))).Methods("GET")
	// entities
	r.Handle("/{view:v}/{item:entities}", securechain.Then(env.AppMiddleware(env.VGetEntitiesHandler))).Methods("GET")
	r.Handle("/{view:vc}/{item:entities}", securechain.Then(env.AppMiddleware(env.VCreateEntityHandler))).Methods("GET")
	r.Handle("/{item:entities}", securechain.Then(env.AppMiddleware(env.GetEntitiesHandler))).Methods("GET")
	r.Handle("/{item:entities}/{id}", securechain.Then(env.AppMiddleware(env.GetEntityHandler))).Methods("GET")
	r.Handle("/{item:entities}/{id}/people", securechain.Then(env.AppMiddleware(env.GetEntityPeopleHandler))).Methods("GET")
	r.Handle("/{item:entities}", securechain.Then(env.AppMiddleware(env.CreateEntityHandler))).Methods("POST")
	r.Handle("/{item:entities}/{id}", securechain.Then(env.AppMiddleware(env.UpdateEntityHandler))).Methods("PUT")
	r.Handle("/{item:entities}/{id}", securechain.Then(env.AppMiddleware(env.DeleteEntityHandler))).Methods("DELETE")
	r.Handle("/{item:stocks}/{id}", securechain.Then(env.AppMiddleware(env.GetEntityStockHandler))).Methods("GET")

	r.Handle("/f/{view:v}/{item:entities}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{view:vc}/{item:entities}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{item:entities}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{item:entities}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{item:entities}/{id}/people", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{item:entities}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("POST")
	r.Handle("/f/{item:entities}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("PUT")
	r.Handle("/f/{item:entities}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("DELETE")
	r.Handle("/f/{item:stocks}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	// people
	r.Handle("/{view:v}/{item:people}", securechain.Then(env.AppMiddleware(env.VGetPeopleHandler))).Methods("GET")
	r.Handle("/{view:vc}/{item:people}", securechain.Then(env.AppMiddleware(env.VCreatePersonHandler))).Methods("GET")
	r.Handle("/{view:vu}/{item:peoplepass}", securechain.Then(env.AppMiddleware(env.VUpdatePersonPasswordHandler))).Methods("GET")
	r.Handle("/{item:people}", securechain.Then(env.AppMiddleware(env.GetPeopleHandler))).Methods("GET")
	r.Handle("/{item:people}/{id}", securechain.Then(env.AppMiddleware(env.GetPersonHandler))).Methods("GET")
	r.Handle("/{item:people}/{id}/entities", securechain.Then(env.AppMiddleware(env.GetPersonEntitiesHandler))).Methods("GET")
	r.Handle("/{item:people}/{id}/manageentities", securechain.Then(env.AppMiddleware(env.GetPersonManageEntitiesHandler))).Methods("GET")
	r.Handle("/{item:people}/{id}/permissions", securechain.Then(env.AppMiddleware(env.GetPersonPermissionsHandler))).Methods("GET")
	r.Handle("/{item:people}/{id}", securechain.Then(env.AppMiddleware(env.UpdatePersonHandler))).Methods("PUT")
	r.Handle("/{item:people}", securechain.Then(env.AppMiddleware(env.CreatePersonHandler))).Methods("POST")
	r.Handle("/{item:people}/{id}", securechain.Then(env.AppMiddleware(env.DeletePersonHandler))).Methods("DELETE")
	r.Handle("/{item:peoplep}", securechain.Then(env.AppMiddleware(env.UpdatePersonpHandler))).Methods("POST")

	r.Handle("/f/{view:v}/{item:people}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{view:vc}/{item:people}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{item:people}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{item:people}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{item:people}/{id}/entities", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{item:people}/{id}/manageentities", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{item:people}/{id}/permissions", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{item:people}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("PUT")
	r.Handle("/f/{item:people}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("POST")
	r.Handle("/f/{item:people}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("DELETE")
	// store locations
	r.Handle("/{view:v}/{item:storelocations}", securechain.Then(env.AppMiddleware(env.VGetStoreLocationsHandler))).Methods("GET")
	r.Handle("/{view:vc}/{item:storelocations}", securechain.Then(env.AppMiddleware(env.VCreateStoreLocationHandler))).Methods("GET")
	r.Handle("/{item:storelocations}", securechain.Then(env.AppMiddleware(env.GetStoreLocationsHandler))).Methods("GET")
	r.Handle("/{item:storelocations}/{id}", securechain.Then(env.AppMiddleware(env.GetStoreLocationHandler))).Methods("GET")
	r.Handle("/{item:storelocations}/{id}", securechain.Then(env.AppMiddleware(env.UpdateStoreLocationHandler))).Methods("PUT")
	r.Handle("/{item:storelocations}", securechain.Then(env.AppMiddleware(env.CreateStoreLocationHandler))).Methods("POST")
	r.Handle("/{item:storelocations}/{id}", securechain.Then(env.AppMiddleware(env.DeleteStoreLocationHandler))).Methods("DELETE")

	r.Handle("/f/{view:v}/{item:storelocations}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{view:vc}/{item:storelocations}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{item:storelocations}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{item:storelocations}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{item:storelocations}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("PUT")
	r.Handle("/f/{item:storelocations}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("POST")
	r.Handle("/f/{item:storelocations}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("DELETE")
	// products
	r.Handle("/{item:products}/l2eformula/{f}", securechain.Then(env.AppMiddleware(env.ConvertProductEmpiricalToLinearFormulaHandler))).Methods("GET")
	r.Handle("/{view:v}/{item:products}", securechain.Then(env.AppMiddleware(env.VGetProductsHandler))).Methods("GET")
	r.Handle("/{view:vc}/{item:products}", securechain.Then(env.AppMiddleware(env.VCreateProductHandler))).Methods("GET")
	r.Handle("/{item:products}", securechain.Then(env.AppMiddleware(env.GetProductsHandler))).Methods("GET")
	r.Handle("/{item:products}/{id}", securechain.Then(env.AppMiddleware(env.GetProductHandler))).Methods("GET")
	r.Handle("/{item:products}/{id}", securechain.Then(env.AppMiddleware(env.UpdateProductHandler))).Methods("PUT")
	r.Handle("/{item:products}", securechain.Then(env.AppMiddleware(env.CreateProductHandler))).Methods("POST")
	r.Handle("/{item:products}/{id}", securechain.Then(env.AppMiddleware(env.DeleteProductHandler))).Methods("DELETE")
	r.Handle("/{item:bookmarks}/{id}", securechain.Then(env.AppMiddleware(env.ToogleProductBookmarkHandler))).Methods("PUT")
	r.Handle("/{item:products}/magic", securechain.Then(env.AppMiddleware(env.MagicHandler))).Methods("POST")

	r.Handle("/{item:products}/casnumbers/", securechain.Then(env.AppMiddleware(env.GetProductsCasNumbersHandler))).Methods("GET")
	r.Handle("/{item:products}/casnumbers/{id}", securechain.Then(env.AppMiddleware(env.GetProductsCasNumberHandler))).Methods("GET")

	r.Handle("/{item:products}/cenumbers/", securechain.Then(env.AppMiddleware(env.GetProductsCeNumbersHandler))).Methods("GET")

	r.Handle("/{item:products}/names/", securechain.Then(env.AppMiddleware(env.GetProductsNamesHandler))).Methods("GET")
	r.Handle("/{item:products}/names/{id}", securechain.Then(env.AppMiddleware(env.GetProductsNameHandler))).Methods("GET")

	r.Handle("/{item:products}/linearformulas/", securechain.Then(env.AppMiddleware(env.GetProductsLinearFormulasHandler))).Methods("GET")

	r.Handle("/{item:products}/empiricalformulas/", securechain.Then(env.AppMiddleware(env.GetProductsEmpiricalFormulasHandler))).Methods("GET")
	r.Handle("/{item:products}/empiricalformulas/{id}", securechain.Then(env.AppMiddleware(env.GetProductsEmpiricalFormulaHandler))).Methods("GET")

	r.Handle("/{item:products}/physicalstates/", securechain.Then(env.AppMiddleware(env.GetProductsPhysicalStatesHandler))).Methods("GET")

	r.Handle("/{item:products}/signalwords/", securechain.Then(env.AppMiddleware(env.GetProductsSignalWordsHandler))).Methods("GET")
	r.Handle("/{item:products}/signalwords/{id}", securechain.Then(env.AppMiddleware(env.GetProductsSignalWordHandler))).Methods("GET")

	r.Handle("/{item:products}/synonyms/", securechain.Then(env.AppMiddleware(env.GetProductsSynonymsHandler))).Methods("GET")

	r.Handle("/{item:products}/symbols/", securechain.Then(env.AppMiddleware(env.GetProductsSymbolsHandler))).Methods("GET")
	r.Handle("/{item:products}/symbols/{id}", securechain.Then(env.AppMiddleware(env.GetProductsSymbolHandler))).Methods("GET")

	r.Handle("/{item:products}/classofcompounds/", securechain.Then(env.AppMiddleware(env.GetProductsClassOfCompoundsHandler))).Methods("GET")

	r.Handle("/{item:products}/hazardstatements/", securechain.Then(env.AppMiddleware(env.GetProductsHazardStatementsHandler))).Methods("GET")
	r.Handle("/{item:products}/hazardstatements/{id}", securechain.Then(env.AppMiddleware(env.GetProductsHazardStatementHandler))).Methods("GET")

	r.Handle("/{item:products}/precautionarystatements/", securechain.Then(env.AppMiddleware(env.GetProductsPrecautionaryStatementsHandler))).Methods("GET")
	r.Handle("/{item:products}/precautionarystatements/{id}", securechain.Then(env.AppMiddleware(env.GetProductsPrecautionaryStatementHandler))).Methods("GET")

	r.Handle("/f/{view:v}/{item:products}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{view:vc}/{item:products}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{item:products}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{item:products}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{item:products}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("PUT")
	r.Handle("/f/{item:products}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("POST")
	r.Handle("/f/{item:products}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("DELETE")
	// storages
	r.Handle("/{view:v}/{item:storages}", securechain.Then(env.AppMiddleware(env.VGetStoragesHandler))).Methods("GET")
	r.Handle("/{view:vc}/{item:storages}", securechain.Then(env.AppMiddleware(env.VCreateStorageHandler))).Methods("GET")
	r.Handle("/{item:storages}", securechain.Then(env.AppMiddleware(env.GetStoragesHandler))).Methods("GET")
	r.Handle("/{item:storages}/others", securechain.Then(env.AppMiddleware(env.GetOtherStoragesHandler))).Methods("GET")
	r.Handle("/{item:storages}/suppliers", securechain.Then(env.AppMiddleware(env.GetStoragesSuppliersHandler))).Methods("GET")
	r.Handle("/{item:storages}/units", securechain.Then(env.AppMiddleware(env.GetStoragesUnitsHandler))).Methods("GET")
	r.Handle("/{item:storages}/{id}", securechain.Then(env.AppMiddleware(env.GetStorageHandler))).Methods("GET")
	r.Handle("/{item:storages}/{id}", securechain.Then(env.AppMiddleware(env.UpdateStorageHandler))).Methods("PUT")
	r.Handle("/{item:storages}", securechain.Then(env.AppMiddleware(env.CreateStorageHandler))).Methods("POST")
	r.Handle("/{item:storages}/{id}", securechain.Then(env.AppMiddleware(env.DeleteStorageHandler))).Methods("DELETE")
	r.Handle("/{item:storages}/{id}/a", securechain.Then(env.AppMiddleware(env.ArchiveStorageHandler))).Methods("DELETE")
	r.Handle("/{item:storages}/{id}/r", securechain.Then(env.AppMiddleware(env.RestoreStorageHandler))).Methods("PUT")
	r.Handle("/{item:borrowings}/{id}", securechain.Then(env.AppMiddleware(env.ToogleStorageBorrowingHandler))).Methods("PUT")

	r.Handle("/f/{item:storages}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{item:storages}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("PUT")
	r.Handle("/f/{item:storages}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("POST")
	r.Handle("/f/{item:storages}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("DELETE")

	// validators
	r.Handle("/validate/entity/{id}/name/", securechain.Then(env.AppMiddleware(env.ValidateEntityNameHandler))).Methods("POST")
	r.Handle("/validate/person/{id}/email/", securechain.Then(env.AppMiddleware(env.ValidatePersonEmailHandler))).Methods("POST")
	r.Handle("/validate/product/{id}/casnumber/", securechain.Then(env.AppMiddleware(env.ValidateProductCasNumberHandler))).Methods("POST")
	r.Handle("/validate/product/{id}/cenumber/", securechain.Then(env.AppMiddleware(env.ValidateProductCeNumberHandler))).Methods("POST")
	r.Handle("/validate/product/{id}/name/", securechain.Then(env.AppMiddleware(env.ValidateProductNameHandler))).Methods("POST")
	r.Handle("/validate/product/{id}/empiricalformula/", securechain.Then(env.AppMiddleware(env.ValidateProductEmpiricalFormulaHandler))).Methods("POST")

	// formatters
	r.Handle("/format/product/{id}/empiricalformula/", securechain.Then(env.AppMiddleware(env.FormatProductEmpiricalFormulaHandler))).Methods("POST")

	// export download
	r.Handle("/{item:download}/{id}", securechain.Then(env.AppMiddleware(env.DownloadExportHandler))).Methods("GET")

	// rice boxes
	webfontsBox := rice.MustFindBox("static/webfonts")
	webfontsFileServer := http.StripPrefix("/webfonts/", http.FileServer(webfontsBox.HTTPBox()))
	http.Handle("/webfonts/", webfontsFileServer)

	fontsBox := rice.MustFindBox("static/fonts")
	fontsFileServer := http.StripPrefix("/fonts/", http.FileServer(fontsBox.HTTPBox()))
	http.Handle("/fonts/", fontsFileServer)

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

	log.Info("- application running")
	if err = http.ListenAndServe(":"+*listenport, nil); err != nil {
		panic(err)
	}
}

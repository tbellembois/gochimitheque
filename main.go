// +build go1.14,linux

//go:generate jade -writer -basedir static/templates -d ./static/jade welcomeannounce/index.jade home/index.jade login/index.jade about/index.jade entity/index.jade entity/create.jade product/index.jade product/create.jade storage/index.jade storage/create.jade storelocation/index.jade storelocation/create.jade person/index.jade person/create.jade person/pupdate.jade test.jade
//go:generate go run main.go -genlocalejs
//go:generate gopherjs build ./static/gjs/gjs-common.go -o ./static/js/chim/gjs-common.js -m
//go:generate rice embed-go
package main

// build with
//go build -trimpath -ldflags "-X main.BuildID=$(git tag | head -1)" -o gochimitheque

import (
	"database/sql"
	"flag"
	"fmt"

	"net/http"
	"net/http/pprof"
	"os"
	"path"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	jsonadapter "github.com/casbin/json-adapter/v2"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/datastores"
	"github.com/tbellembois/gochimitheque/globals"
	"github.com/tbellembois/gochimitheque/handlers"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/static/localejs"
	"github.com/tbellembois/gochimitheque/utils"
)

var (
	// BuildID is compile-time variable
	BuildID string
	// starting flags
	listenport, dbpath, proxyurl, proxypath, mailServerAddress, mailServerPort, mailServerSender, admins, logfile, importv1from, importfrom, mailTest    *string
	mailServerUseTLS, mailServerTLSSkipVerify, enablePublicProductsEndpoint, resetAdminPassword, updateQRCode, debug, version, genLocaleJS, disableCache *bool
)

func init() {
	// configuration parameters
	listenport = flag.String("listenport", "8081", "the port to listen")
	proxyurl = flag.String("proxyurl", "", "the application url (without the path) if behind a proxy, with NO trailing /")
	proxypath = flag.String("proxypath", "/", "the application path if behind a proxy, with the trailing /")
	dbpath = flag.String("dbpath", "./", "the application sqlite directory path")
	mailServerAddress = flag.String("mailserveraddress", "localhost", "the mail server address")
	mailServerPort = flag.String("mailserverport", "25", "the mail server port")
	mailServerSender = flag.String("mailserversender", "", "the mail server sender")
	mailServerUseTLS = flag.Bool("mailserverusetls", false, "use TLS? (optional)")
	mailServerTLSSkipVerify = flag.Bool("mailservertlsskipverify", false, "skip TLS verification? (optional)")
	enablePublicProductsEndpoint = flag.Bool("enablepublicproductsendpoint", false, "enable public products endpoint (optional)")
	admins = flag.String("admins", "", "the additional admins (comma separated email adresses) (optional) ")
	logfile = flag.String("logfile", "", "log to the given file (optional)")
	debug = flag.Bool("debug", false, "debug (verbose log), default is error")
	disableCache = flag.Bool("disablecache", false, "disable the cache (development only)")

	// one shot commands
	resetAdminPassword = flag.Bool("resetadminpassword", false, "reset the admin password to `chimitheque`")
	updateQRCode = flag.Bool("updateqrcode", false, "regenerate storages QR codes")
	version = flag.Bool("version", false, "display application version")
	mailTest = flag.String("mailtest", "", "send a test mail")
	importv1from = flag.String("importv1from", "", "full path of the directory containing the Chimithèque v1 CSV to import")
	importfrom = flag.String("importfrom", "", "base URL of the external Chimithèque instance (running with -enablepublicproductsendpoint) to import products from")
	genLocaleJS = flag.Bool("genlocalejs", false, "generate JS locales (developper target)")
	flag.Parse()
}

func main() {

	var (
		err         error
		logf, logif *os.File
		dbname      string
		datastore   datastores.Datastore
	)

	// default BuildID
	if BuildID == "" {
		BuildID = "devel"
	}

	// displaying version
	if *version {
		fmt.Println(BuildID)
		os.Exit(0)
	}

	// generate JS locales
	if *genLocaleJS {
		localejs.GenerateLocalJS()
		os.Exit(0)
	}

	// building db path
	dbname = path.Join(*dbpath, "storage.db")

	// setting up logger
	globals.Log = logrus.New()
	globals.LogInternal = logrus.New()

	// setting the log level
	if *debug {
		globals.Log.SetLevel(logrus.DebugLevel)
	} else {
		globals.Log.SetLevel(logrus.InfoLevel)
	}

	// logging to file if logfile parameter specified
	if *logfile != "" {
		if logf, err = os.OpenFile(*logfile, os.O_WRONLY|os.O_CREATE, 0755); err != nil {
			globals.Log.Fatal(err)
		} else {
			globals.Log.SetOutput(logf)
		}
	}
	defer logf.Close()

	// internal server error log file
	if logif, err = os.OpenFile("errors.log", os.O_WRONLY|os.O_CREATE, 0755); err != nil {
		globals.Log.Fatal(err)
	} else {
		globals.LogInternal.SetOutput(logif)
		globals.LogInternal.SetReportCaller(true)
	}
	defer logif.Close()

	// global variables init
	globals.BuildID = BuildID
	globals.ProxyPath = *proxypath
	globals.ProxyURL = *proxyurl
	if *proxyurl != "" {
		globals.ApplicationFullURL = *proxyurl + *proxypath
	} else {
		globals.ApplicationFullURL = "http://localhost:" + *listenport + "/"
	}
	globals.MailServerAddress = *mailServerAddress
	globals.MailServerSender = *mailServerSender
	globals.MailServerPort = *mailServerPort
	globals.MailServerUseTLS = *mailServerUseTLS
	globals.MailServerTLSSkipVerify = *mailServerTLSSkipVerify
	globals.Log.Info("- application version: " + globals.BuildID)
	globals.Log.Info("- application endpoint: " + globals.ApplicationFullURL)
	globals.Log.Info("- application db: " + dbname)
	globals.Log.Debug("- globals.MailServerAddress: " + globals.MailServerAddress)
	globals.Log.Debug("- globals.MailServerPort: " + globals.MailServerPort)
	globals.Log.Debug("- globals.MailServerSender: " + globals.MailServerSender)

	// database initialization
	globals.Log.Info("- opening database connection to " + dbname)
	if datastore, err = datastores.NewSQLiteDBstore(dbname); err != nil {
		globals.Log.Fatal(err)
	}
	globals.Log.Info("- creating database if needed")
	if err = datastore.CreateDatabase(); err != nil {
		globals.Log.Fatal(err)
	}
	globals.Log.Info("- running maintenance job")
	datastore.Maintenance()

	if *importv1from != "" {
		globals.Log.Info("- import from Chimithèque v1 csv into database")
		err := datastore.ImportV1(*importv1from)
		if err != nil {
			globals.Log.Error("an error occured: " + err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
	if *importfrom != "" {
		globals.Log.Info("- import from URL into database")
		err := datastore.Import(*importfrom)
		if err != nil {
			globals.Log.Error("an error occured: " + err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
	if *resetAdminPassword {
		globals.Log.Info("- reseting admin password to `chimitheque`")
		a, err := datastore.GetPersonByEmail("admin@chimitheque.fr")
		if err != nil {
			globals.Log.Error("an error occured: " + err.Error())
			os.Exit(1)
		}
		a.PersonPassword = "chimitheque"
		err = datastore.UpdatePersonPassword(a)
		if err != nil {
			globals.Log.Error("an error occured: " + err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
	if *updateQRCode {
		globals.Log.Info("- updating storages QR codes")
		err := datastore.UpdateAllQRCodes()
		if err != nil {
			globals.Log.Error("an error occured: " + err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
	if *mailTest != "" {
		globals.Log.Info("- sending a mail to " + *mailTest)
		err := utils.TestMail(*mailTest)
		if err != nil {
			globals.Log.Error("an error occured: " + err.Error())
			os.Exit(1)
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
		globals.Log.Fatal(err)
	}
	// cleaning former admins
	for _, fa := range formerAdmins {
		isStillAdmin = false
		globals.Log.Info("former admin: " + fa.PersonEmail)
		for _, ca := range currentAdmins {
			if ca == fa.PersonEmail {
				isStillAdmin = true
			}
		}
		if !isStillAdmin {
			globals.Log.Info(fa.PersonEmail + " is not an admin anymore, removing permissions")
			if err = datastore.UnsetPersonAdmin(fa.PersonID); err != nil {
				globals.Log.Fatal(err)
			}
		}
	}
	// setting up new ones
	if len(currentAdmins) > 0 {
		for _, ca := range currentAdmins {
			globals.Log.Info("additional admin: " + ca)
			if p, err = datastore.GetPersonByEmail(ca); err != nil {
				if err == sql.ErrNoRows {
					globals.Log.Fatal("user " + ca + " not found in database")
				} else {
					globals.Log.Fatal(err)
				}
			}

			if err = datastore.SetPersonAdmin(p.PersonID); err != nil {
				globals.Log.Fatal(err)
			}
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
	r.Handle("/reset-password", commonChain.Then(env.AppMiddleware(env.ResetPasswordHandler))).Methods("POST")
	r.Handle("/reset", commonChain.Then(env.AppMiddleware(env.ResetHandler))).Methods("GET")
	r.Handle("/captcha", commonChain.Then(env.AppMiddleware(env.CaptchaHandler))).Methods("GET")
	r.Handle("/{item:delete-token}", securechain.Then(env.AppMiddleware(env.DeleteTokenHandler))).Methods("GET")

	// about
	r.Handle("/about", commonChain.Then(env.AppMiddleware(env.AboutHandler))).Methods("GET")

	// products public
	if *enablePublicProductsEndpoint {
		r.Handle("/e/{item:products}", commonChain.Then(env.AppMiddleware(env.GetExposedProductsHandler))).Methods("GET")
	}

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
	r.Handle("/{item:entities}/stocks/{id}", securechain.Then(env.AppMiddleware(env.GetEntityStockHandler))).Methods("GET")

	r.Handle("/f/{view:v}/{item:entities}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{view:vc}/{item:entities}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{item:entities}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{item:entities}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{item:entities}/{id}/people", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
	r.Handle("/f/{item:entities}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("POST")
	r.Handle("/f/{item:entities}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("PUT")
	r.Handle("/f/{item:entities}/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("DELETE")
	r.Handle("/f/{item:entities}/stocks/{id}", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("GET")
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

	r.Handle("/{item:products}/producerrefs/", securechain.Then(env.AppMiddleware(env.GetProductsProducerRefsHandler))).Methods("GET")
	r.Handle("/{item:products}/producers/", securechain.Then(env.AppMiddleware(env.GetProductsProducersHandler))).Methods("GET")
	r.Handle("/{item:products}/supplierrefs/", securechain.Then(env.AppMiddleware(env.GetProductsSupplierRefsHandler))).Methods("GET")
	r.Handle("/{item:products}/suppliers/", securechain.Then(env.AppMiddleware(env.GetProductsSuppliersHandler))).Methods("GET")
	r.Handle("/{item:products}/categories/", securechain.Then(env.AppMiddleware(env.GetProductsCategoriesHandler))).Methods("GET")
	r.Handle("/{item:products}/tags/", securechain.Then(env.AppMiddleware(env.GetProductsTagsHandler))).Methods("GET")

	r.Handle("/{item:products}/producers", securechain.Then(env.AppMiddleware(env.CreateProducerHandler))).Methods("POST")
	r.Handle("/{item:products}/suppliers", securechain.Then(env.AppMiddleware(env.CreateSupplierHandler))).Methods("POST")

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
	r.Handle("/{item:validate}/entity/{id}/name/", securechain.Then(env.AppMiddleware(env.ValidateEntityNameHandler))).Methods("POST")
	r.Handle("/{item:validate}/person/{id}/email/", securechain.Then(env.AppMiddleware(env.ValidatePersonEmailHandler))).Methods("POST")
	r.Handle("/{item:validate}/product/{id}/casnumber/", securechain.Then(env.AppMiddleware(env.ValidateProductCasNumberHandler))).Methods("POST")
	r.Handle("/{item:validate}/product/{id}/cenumber/", securechain.Then(env.AppMiddleware(env.ValidateProductCeNumberHandler))).Methods("POST")
	r.Handle("/{item:validate}/product/{id}/name/", securechain.Then(env.AppMiddleware(env.ValidateProductNameHandler))).Methods("POST")
	r.Handle("/{item:validate}/product/{id}/empiricalformula/", securechain.Then(env.AppMiddleware(env.ValidateProductEmpiricalFormulaHandler))).Methods("POST")

	// formatters
	r.Handle("/{item:format}/product/{id}/empiricalformula/", securechain.Then(env.AppMiddleware(env.FormatProductEmpiricalFormulaHandler))).Methods("POST")

	// export download
	r.Handle("/{item:download}/{id}", securechain.Then(env.AppMiddleware(env.DownloadExportHandler))).Methods("GET")

	// rice boxes
	modelBox := rice.MustFindBox("models")
	modelf, e := modelBox.Open("model.conf")
	if e != nil {
		globals.Log.Error("model.conf load from box error: " + e.Error())
		os.Exit(1)
	}
	models, e := modelf.Stat()
	if e != nil {
		globals.Log.Error("model.conf stat error: " + e.Error())
		os.Exit(1)
	}

	modelb := make([]byte, models.Size()-1)
	_, e = modelf.Read(modelb)
	if e != nil {
		globals.Log.Error("model.conf load error: " + e.Error())
		os.Exit(1)
	}

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

	// setting up enforcer policy
	if globals.JSONAdapterData, err = datastore.ToCasbinJSONAdapter(); err != nil {
		globals.Log.Error("error getting json adapter data: " + err.Error())
		os.Exit(1)
	}

	m, e := model.NewModelFromString(string(modelb))
	if e != nil {
		globals.Log.Error("model creation error: " + e.Error())
		os.Exit(1)
	}

	a := jsonadapter.NewAdapter(&globals.JSONAdapterData)
	if globals.Enforcer, err = casbin.NewEnforcer(m, a); err != nil {
		globals.Log.Error("enforcer creation error: " + err.Error())
		os.Exit(1)
	}
	globals.Enforcer.AddFunction("matchStorage", env.MatchStorageFunc)
	globals.Enforcer.AddFunction("matchStorelocation", env.MatchStorelocationFunc)
	globals.Enforcer.AddFunction("matchPeople", env.MatchPeopleFunc)
	globals.Enforcer.AddFunction("matchEntity", env.MatchEntityFunc)

	// enforccer policy load
	if err = globals.Enforcer.LoadPolicy(); err != nil {
		globals.Log.Error("enforcer policy load error: " + err.Error())
		os.Exit(1)
	}

	globals.Log.Info("- application running")
	if err = http.ListenAndServe(":"+*listenport, nil); err != nil {
		panic("error running the server")
	}
}

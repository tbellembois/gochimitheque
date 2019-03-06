package main

//go:generate go run gogenerate/localejs.go
//go:generate gopherjs build gopherjs/gjs-common.go -o static/js/gjs-common.js
//go:generate rice embed-go

import (
	"flag"
	"html/template"
	"net/http"
	"os"

	rice "github.com/GeertJohan/go.rice"
	"github.com/Joker/jade"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/global"
	"github.com/tbellembois/gochimitheque/handlers"
	"github.com/tbellembois/gochimitheque/models"
	//"github.com/pkg/profile"
)

func main() {

	// CPU profiling by default
	//defer profile.Start().Stop()

	var (
		err       error
		logf      *os.File
		dbname    = "/mnt/ramdisk/storage.db"
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
			log.Panic(err)
		} else {
			log.SetOutput(logf)
		}
	}

	// global variables init
	global.ProxyPath = *proxypath
	global.ProxyURL = *proxyurl
	global.MailServerAddress = *mailServerAddress
	global.MailServerSender = *mailServerSender
	global.MailServerPort = *mailServerPort
	global.MailServerUseTLS = *mailServerUseTLS
	global.MailServerTLSSkipVerify = *mailServerTLSSkipVerify
	log.Info("- application endpoint: " + global.ProxyURL + global.ProxyPath)

	// database initialization
	log.Info("- opening database connection")
	if datastore, err = models.NewSQLiteDBstore(dbname); err != nil {
		log.Panic(err)
	}
	log.Info("- creating database if needed")
	if err = datastore.CreateDatabase(); err != nil {
		log.Panic(err)
	}
	if *importfrom == "" {
		log.Info("- inserting database sample values if needed")
		if err = datastore.InsertSamples(); err != nil {
			log.Panic(err)
		}
	} else {
		log.Info("- import from csv into database")
		err := datastore.Import(*importfrom)
		if err != nil {
			log.Error("an error occured: " + err.Error())
		}
		os.Exit(0)
	}

	// environment creation
	env := handlers.Env{
		DB: datastore,
		//Localizer: localizer,
		Templates: make(map[string]*template.Template),
	}

	// HasPermission used by template rendering to show/hide html elements
	funcMap := template.FuncMap{
		"HasPermission": func(id int, perm string, item string, itemid int) bool {
			p, e := env.DB.HasPersonPermission(id, perm, item, itemid)
			if e != nil {
				log.Error(e.Error())
			}
			return p
		},
		// i18n
		"T": func(messageID string, pluralCount int) string {
			return env.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: messageID, PluralCount: pluralCount})
		},
	}

	// template compilation
	b := rice.MustFindBox("static/templates")
	basejades := []string{
		"base.jade",
		"localejs.jade",
		"mixins.jade",
		"head.jade",
		"header.jade",
		"footer.jade",
		"foot.jade"}
	basenomenujades := []string{
		"basenomenu.jade",
		"localejs.jade",
		"mixins.jade",
		"head.jade",
		"header.jade",
		"footer.jade",
		"foot.jade"}
	basejadess := []byte{}
	basenomenujadess := []byte{}
	for _, s := range basejades {
		basejadess = append(basejadess, b.MustBytes(s)...)
	}
	for _, s := range basenomenujades {
		basenomenujadess = append(basenomenujadess, b.MustBytes(s)...)
	}

	// test
	testtmpl, e := jade.Parse("test", append(basenomenujadess, b.MustBytes("test.jade")...))
	if e != nil {
		log.Fatal("testtmpl jade:" + e.Error())
	}
	env.Templates["test"], err = template.New("test").Funcs(funcMap).Parse(testtmpl)
	if err != nil {
		log.Fatal("testtmpl parse:" + e.Error())
	}
	// home
	hometmpl, e := jade.Parse("home_index", append(basejadess, b.MustBytes("home/index.jade")...))
	if e != nil {
		log.Fatal("hometmpl jade:" + e.Error())
	}
	env.Templates["home"], err = template.New("home").Funcs(funcMap).Parse(hometmpl)
	if err != nil {
		log.Fatal("hometmpl parse:" + e.Error())
	}
	// login
	logintmpl, e := jade.Parse("login_index", append(basenomenujadess, b.MustBytes("login/index.jade")...))
	if e != nil {
		log.Fatal("logintmpl jade:" + e.Error())
	}
	env.Templates["login"], err = template.New("login").Funcs(funcMap).Parse(logintmpl)
	if err != nil {
		log.Fatal("logintmpl parse:" + e.Error())
	}
	// entity
	entityindextmpl, e := jade.Parse("entity_index", append(append(basejadess, b.MustString("entity/commonjs.jade")...), b.MustString("entity/index.jade")...))
	if e != nil {
		log.Fatal("entityindextmpl jade:" + e.Error())
	}
	env.Templates["entityindex"], err = template.New("entityindex").Funcs(funcMap).Parse(entityindextmpl)
	if err != nil {
		log.Fatal("entityindextmpl parse:" + e.Error())
	}
	entitycreatetmpl, e := jade.Parse("entity_create", append(append(basejadess, b.MustString("entity/commonjs.jade")...), b.MustString("entity/create.jade")...))
	if e != nil {
		log.Fatal("entitycreatetmpl jade:" + e.Error())
	}
	env.Templates["entitycreate"], err = template.New("entitycreate").Funcs(funcMap).Parse(entitycreatetmpl)
	if err != nil {
		log.Fatal("entitycreatetmpl parse:" + e.Error())
	}
	// store location
	storelocationindextmpl, e := jade.Parse("storelocation_index", append(append(basejadess, b.MustString("storelocation/commonjs.jade")...), b.MustString("storelocation/index.jade")...))
	if e != nil {
		log.Fatal("storelocationindextmpl jade:" + e.Error())
	}
	env.Templates["storelocationindex"], err = template.New("storelocationindex").Funcs(funcMap).Parse(storelocationindextmpl)
	if err != nil {
		log.Fatal("storelocationtmpl parse:" + e.Error())
	}
	storelocationcreatetmpl, e := jade.Parse("storelocation_create", append(append(basejadess, b.MustString("storelocation/commonjs.jade")...), b.MustString("storelocation/create.jade")...))
	if e != nil {
		log.Fatal("storelocationcreatetmpl jade:" + e.Error())
	}
	env.Templates["storelocationcreate"], err = template.New("storelocationcreate").Funcs(funcMap).Parse(storelocationcreatetmpl)
	if err != nil {
		log.Fatal("storelocationcreatetmpl parse:" + e.Error())
	}
	// person
	personindextmpl, e := jade.Parse("person_index", append(append(basejadess, b.MustString("person/commonjs.jade")...), b.MustString("person/index.jade")...))
	if e != nil {
		log.Fatal("personindextmpl jade:" + e.Error())
	}
	env.Templates["personindex"], err = template.New("personindex").Funcs(funcMap).Parse(personindextmpl)
	if err != nil {
		log.Fatal("personindextmpl parse:" + e.Error())
	}
	personcreatetmpl, e := jade.Parse("person_create", append(append(basejadess, b.MustString("person/commonjs.jade")...), b.MustString("person/create.jade")...))
	if e != nil {
		log.Fatal("personcreatetmpl jade:" + e.Error())
	}
	env.Templates["personcreate"], err = template.New("personcreate").Funcs(funcMap).Parse(personcreatetmpl)
	if err != nil {
		log.Fatal("personcreatetmpl parse:" + e.Error())
	}
	personpupdatetmpl, e := jade.Parse("personp_update", append(append(basejadess, b.MustString("person/commonjs.jade")...), b.MustString("person/pupdate.jade")...))
	if e != nil {
		log.Fatal("personpupdatetmpl jade:" + e.Error())
	}
	env.Templates["personpupdate"], err = template.New("personpupdate").Funcs(funcMap).Parse(personpupdatetmpl)
	if err != nil {
		log.Fatal("personpupdatetmpl parse:" + e.Error())
	}
	// product
	productindextmpl, e := jade.Parse("product_index", append(append(basejadess, b.MustString("product/commonjs.jade")...), b.MustString("product/index.jade")...))
	if e != nil {
		log.Fatal("productindextmpl jade:" + e.Error())
	}
	env.Templates["productindex"], err = template.New("productindex").Funcs(funcMap).Parse(productindextmpl)
	if err != nil {
		log.Fatal("productindextmpl parse:" + e.Error())
	}
	productcreatetmpl, e := jade.Parse("product_create", append(append(basejadess, b.MustString("product/commonjs.jade")...), b.MustString("product/create.jade")...))
	if err != nil {
		log.Fatal("productcreatetmpl jade:" + e.Error())
	}
	env.Templates["productcreate"], err = template.New("productcreate").Funcs(funcMap).Parse(productcreatetmpl)
	if err != nil {
		log.Fatal("productcreatetmpl parse:" + e.Error())
	}
	// storage
	storageindextmpl, e := jade.Parse("storage_index", append(append(basejadess, b.MustString("storage/commonjs.jade")...), b.MustString("storage/index.jade")...))
	if e != nil {
		log.Fatal("storageindextmpl jade:" + e.Error())
	}
	env.Templates["storageindex"], err = template.New("storageindex").Funcs(funcMap).Parse(storageindextmpl)
	if err != nil {
		log.Fatal("storageindextmpl parse:" + e.Error())
	}
	storagecreatetmpl, e := jade.Parse("storage_create", append(append(basejadess, b.MustString("storage/commonjs.jade")...), b.MustString("storage/create.jade")...))
	if e != nil {
		log.Fatal("storagecreatetmpl jade:" + e.Error())
	}
	env.Templates["storagecreate"], err = template.New("storagecreate").Funcs(funcMap).Parse(storagecreatetmpl)
	if err != nil {
		log.Fatal("storagecreatetmpl parse:" + e.Error())
	}

	// router definition
	r := mux.NewRouter()
	commonChain := alice.New(env.ContextMiddleware, env.HeadersMiddleware, env.LogingMiddleware)
	securechain := alice.New(env.ContextMiddleware, env.HeadersMiddleware, env.LogingMiddleware, env.AuthenticateMiddleware, env.AuthorizeMiddleware)

	// login
	r.Handle("/login", commonChain.Then(env.AppMiddleware(env.VLoginHandler))).Methods("GET")
	r.Handle("/get-token", commonChain.Then(env.AppMiddleware(env.GetTokenHandler))).Methods("POST")
	r.Handle("/delete-token", securechain.Then(env.AppMiddleware(env.DeleteTokenHandler))).Methods("GET")
	r.Handle("/reset-password", commonChain.Then(env.AppMiddleware(env.ResetPasswordHandler))).Methods("POST")
	r.Handle("/reset", commonChain.Then(env.AppMiddleware(env.ResetHandler))).Methods("GET")
	r.Handle("/captcha", commonChain.Then(env.AppMiddleware(env.CaptchaHandler))).Methods("GET")

	// developper tests
	r.Handle("/v/test", securechain.Then(env.AppMiddleware(env.VTestHandler))).Methods("GET")
	// home page
	r.Handle("/", securechain.Then(env.AppMiddleware(env.HomeHandler))).Methods("GET")
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
	r.Handle("/f/{item:storages}/", securechain.Then(env.AppMiddleware(env.FakeHandler))).Methods("POST")
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
	r.Handle("/download/{id}", securechain.Then(env.AppMiddleware(env.DownloadExportHandler))).Methods("GET")

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

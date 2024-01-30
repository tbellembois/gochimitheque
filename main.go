//go:build go1.21 && linux && amd64

//go:generate jade -writer -basedir static/templates -d ./static/jade welcomeannounce/index.jade home/index.jade login/index.jade about/index.jade entity/index.jade entity/create.jade product/index.jade product/create.jade product/pubchem.jade storage/index.jade storage/create.jade storelocation/index.jade storelocation/create.jade person/index.jade person/create.jade person/password.jade person/qrcode.jade search.jade menu.jade
//go:generate go run . -genlocalejs
package main

// compile with:
// BuildID="v2.1.0" && go build -ldflags "-X main.BuildID=$BuildID".
import (
	"database/sql"
	"embed"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	zmq "github.com/pebbe/zmq4"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/casbin"
	"github.com/tbellembois/gochimitheque/datastores"
	"github.com/tbellembois/gochimitheque/handlers"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/static/localejs"
	"github.com/tbellembois/gochimitheque/zmqclient"
)

var (
	env handlers.Env

	// Starting parameters and commands.
	paramDBPath,
	paramAdminList,
	paramAutoImportURL *string
	paramPublicProductsEndpoint,
	commandUpdateQRCode,
	paramDebug,
	commandVersion,
	commandGenLocaleJS *bool
	BuildID string

	//go:embed wasm/*
	embedWasmBox embed.FS
	//go:embed static/*
	embedStaticBox embed.FS
)

func init() {
	env = handlers.NewEnv()

	// Configuration parameters.
	flagDBPath := flag.String("dbpath", "./", "the application sqlite directory path")
	flagAppURL := flag.String("appurl", "http://localhost:8081", "the application url (without the path), with NO trailing /")
	flagAppPath := flag.String("apppath", "/", "the application path with the trailing /")
	flagDockerPort := flag.Int("dockerport", 0, "application listen port while running in docker")

	flagOIDCISSUER := flag.String("oidcissuer", "http://localhost:7001", "the OIDC issuer URL")
	flagOIDCClientID := flag.String("oidcclientid", "chimitheque", "the OIDC client ID")
	flagOIDCClientSecret := flag.String("oidcclientsecret", "chimitheque", "the OIDC client secret")
	flagOIDCTokenEndpoint := flag.String("oidctokenendpoint", "http://localhost:7001/api/login/oauth/access_token", "the OIDC token endpoint")
	flagOIDCAuthEndpoint := flag.String("oidcauthendpoint", "http://localhost:7001/login/oauth/authorize", "the OIDC authorization endpoint")
	flagOIDCDeviceEndpoint := flag.String("oidcdeviceendpoint", "http://localhost:7001", "the OIDC device endpoint")

	flagPublicProductsEndpoint := flag.Bool("enablepublicproductsendpoint", false, "enable public products endpoint (optional)")
	flagAdminList := flag.String("admins", "", "the additional admins (comma separated email adresses) (optional) ")
	flagAutoImportURL := flag.String("autoimporturl", "", "the URL of the chimitheque instance to import initial products (optional) ")
	flagDebug := flag.Bool("debug", false, "debug (verbose log), default is error")

	// One shot commands.
	flagUpdateQRCode := flag.Bool("updateqrcode", false, "regenerate storages QR codes")
	flagVersion := flag.Bool("version", false, "display application version")
	flagGenLocaleJS := flag.Bool("genlocalejs", false, "generate JS locales (developper target)")

	flag.Parse()

	env.AppURL = *flagAppURL
	env.AppPath = *flagAppPath
	env.DockerPort = *flagDockerPort
	env.OIDCIssuer = *flagOIDCISSUER
	env.OIDCClientID = *flagOIDCClientID
	env.OIDCClientSecret = *flagOIDCClientSecret
	env.OIDCTokenEndpoint = *flagOIDCTokenEndpoint
	env.OIDCAuthEndpoint = *flagOIDCAuthEndpoint
	env.OIDCDeviceEndpoint = *flagOIDCDeviceEndpoint

	paramDBPath = flagDBPath
	paramPublicProductsEndpoint = flagPublicProductsEndpoint
	paramAdminList = flagAdminList
	paramAutoImportURL = flagAutoImportURL
	paramDebug = flagDebug

	commandUpdateQRCode = flagUpdateQRCode
	commandVersion = flagVersion
	commandGenLocaleJS = flagGenLocaleJS

	env.AppFullURL = env.AppURL + env.AppPath
	env.BuildID = BuildID
}

func initLogger() {
	if *paramDebug {
		logger.Log.SetLevel(logrus.DebugLevel)
	} else {
		logger.Log.SetLevel(logrus.InfoLevel)
	}
}

func initOIDC() {

	var err error = errors.New("fake")
	for err != nil {
		env.OIDCProvider, err = oidc.NewProvider(context.Background(), env.OIDCIssuer)
		logger.Log.Info("- sleeping 2 seconds waiting for OIDC issuer " + env.OIDCIssuer)
		time.Sleep(2 * time.Second)
	}

	env.OIDCConfig = &oidc.Config{
		ClientID: env.OIDCClientID,
	}
	env.OIDCVerifier = env.OIDCProvider.Verifier(env.OIDCConfig)
	env.OAuth2Config = oauth2.Config{
		ClientID:     env.OIDCClientID,
		ClientSecret: env.OIDCClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL:      env.OIDCTokenEndpoint,
			DeviceAuthURL: env.OIDCDeviceEndpoint,
			AuthURL:       env.OIDCAuthEndpoint,
		},
		RedirectURL: env.AppURL + env.AppPath + "callback",
		Scopes:      []string{oidc.ScopeOpenID, "profile", "email"},
	}

}

func initDB() {
	var (
		err       error
		datastore datastores.Datastore
	)

	dbname := path.Join(*paramDBPath, "storage.db")
	logger.Log.Info("- opening database connection to " + dbname)
	if datastore, err = datastores.NewSQLiteDBstore(dbname); err != nil {
		logger.Log.Fatal(err)
	}

	logger.Log.Info("- creating database if needed")
	if err = datastore.CreateDatabase(); err != nil {
		logger.Log.Fatal(err)
	}

	logger.Log.Info("- running maintenance job")
	datastore.Maintenance()

	env.DB = datastore

	var productCount int

	if productCount, err = env.DB.CountProducts(); err != nil {
		logger.Log.Fatal(err)
	}

	if productCount == 0 && paramAutoImportURL != nil {
		logger.Log.Info("- importing initial product list from " + *paramAutoImportURL)
		if err = env.DB.Import(*paramAutoImportURL); err != nil {
			logger.Log.Fatal(err)
		}
	}

}

func initAdmins() {
	var (
		err           error
		p             models.Person
		formerAdmins  []models.Person
		currentAdmins []string
		isStillAdmin  bool
	)

	if *paramAdminList != "" {
		currentAdmins = strings.Split(*paramAdminList, ",")
	}

	if formerAdmins, err = env.DB.GetAdmins(); err != nil {
		logger.Log.Fatal(err)
	}

	// Cleaning former admins.
	for _, fa := range formerAdmins {
		isStillAdmin = false

		logger.Log.Info("former admin: " + fa.PersonEmail)

		for _, ca := range currentAdmins {
			if ca == fa.PersonEmail {
				isStillAdmin = true
			}
		}
		if !isStillAdmin {
			logger.Log.Info(fa.PersonEmail + " is not an admin anymore, removing permissions")
			if err = env.DB.UnsetPersonAdmin(fa.PersonID); err != nil {
				logger.Log.Fatal(err)
			}
		}
	}
	// Setting up new ones.
	if len(currentAdmins) > 0 {
		for _, ca := range currentAdmins {
			logger.Log.Info("additional admin: " + ca)
			if p, err = env.DB.GetPersonByEmail(ca); err != nil {
				if err == sql.ErrNoRows {
					logger.Log.Info("user " + ca + " not found in database, creating it")
					if _, err = env.DB.CreatePerson(models.Person{PersonEmail: ca}); err != nil {
						logger.Log.Fatal(err)
					}
				} else {
					logger.Log.Fatal(err)
				}
			}

			if err = env.DB.SetPersonAdmin(p.PersonID); err != nil {
				logger.Log.Fatal(err)
			}
		}
	}
}

func initStaticResources(router *mux.Router) {
	http.Handle("/wasm/", http.FileServer(http.FS(embedWasmBox)))
	http.Handle("/static/", http.FileServer(http.FS(embedStaticBox)))
	http.Handle("/", router)
}

func main() {
	var err error

	// Basic commands.
	if *commandVersion {
		fmt.Println(env.BuildID)
		os.Exit(0)
	}

	if *commandGenLocaleJS {
		localejs.GenerateLocalJS()
		os.Exit(0)
	}

	initLogger()

	if zmqclient.Zctx, err = zmq.NewContext(); err != nil {
		logger.Log.Fatal(err)
	}

	logger.Log.WithFields(logrus.Fields{
		"commandUpdateQRCode": commandUpdateQRCode,
		"commandVersion":      commandVersion,
		"commandGenLocaleJS":  commandGenLocaleJS,
	}).Debug("main")

	initDB()

	initOIDC()

	logger.Log.Debugf("- env: %+v", env)
	logger.Log.Info("- application version: " + env.BuildID)
	logger.Log.Info("- application endpoint: " + env.AppFullURL)

	// Advanced commands.
	if *commandUpdateQRCode {
		logger.Log.Info("- updating storages QR codes")
		err := env.DB.UpdateAllQRCodes()
		if err != nil {
			logger.Log.Error("an error occurred: " + err.Error())
			os.Exit(1)
		}

		os.Exit(0)
	}

	initAdmins()

	router := buildEndpoints(env.AppFullURL)

	initStaticResources(router)

	env.Enforcer = casbin.InitCasbinPolicy(env.DB)

	var listenAddr string
	if env.DockerPort != 0 {
		listenAddr = fmt.Sprintf(":%d", env.DockerPort)
	} else {
		listenAddr = strings.Split(env.AppURL, "//")[1]
	}

	logger.Log.Infof("- application listening on %s", listenAddr)
	if err = http.ListenAndServe(listenAddr, nil); err != nil {
		panic("error running the server:" + err.Error())
	}
}

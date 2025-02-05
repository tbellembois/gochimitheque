//go:build go1.21 && linux && amd64

//go:generate jade -writer -basedir static/templates -d ./static/jade welcomeannounce/index.jade home/index.jade login/index.jade about/index.jade entity/index.jade entity/create.jade product/index.jade product/create.jade product/pubchem.jade storage/index.jade storage/create.jade storelocation/index.jade storelocation/create.jade person/index.jade person/create.jade person/password.jade person/qrcode.jade search.jade menu.jade
//go:generate go run . -genlocalejs
package main

// compile with:
// BuildID="v2.1.0" && go build -ldflags "-X main.BuildID=$BuildID".
import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	zmq "github.com/pebbe/zmq4"
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
	commandUpdateQRCode,
	paramDebug,
	paramFakeAuth,
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

	// keycloak
	flagOIDCDiscoverURL := flag.String("oidcdiscoverurl", "http://localhost:8080/keycloak/realms/chimitheque/.well-known/openid-configuration", "the OIDC server discover URL")
	flagOIDCClientID := flag.String("oidcclientid", "chimitheque", "the OIDC client ID")
	flagOIDCClientSecret := flag.String("oidcclientsecret", "mysupersecret", "the OIDC client secret")

	flagAdminList := flag.String("admins", "", "the additional admins (comma separated email adresses) (optional) ")
	flagDebug := flag.Bool("debug", false, "debug (verbose log), default is error")
	flagFakeAuth := flag.Bool("fakeauth", false, "fake authentication (use in devel only), default is false")

	// One shot commands.
	flagUpdateQRCode := flag.Bool("updateqrcode", false, "regenerate storages QR codes")
	flagVersion := flag.Bool("version", false, "display application version")
	flagGenLocaleJS := flag.Bool("genlocalejs", false, "generate JS locales (developper target)")

	flag.Parse()

	env.AppURL = *flagAppURL
	env.AppPath = *flagAppPath
	env.DockerPort = *flagDockerPort
	env.OIDCDiscoverURL = *flagOIDCDiscoverURL
	env.OIDCClientID = *flagOIDCClientID
	env.OIDCClientSecret = *flagOIDCClientSecret

	paramDBPath = flagDBPath
	paramAdminList = flagAdminList
	paramDebug = flagDebug
	paramFakeAuth = flagFakeAuth

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

type OIDCDiscover struct {
	Issuer                      string `json:"issuer"`
	AuthorizationEndpoint       string `json:"authorization_endpoint"`
	DeviceAuthorizationEndpoint string `json:"device_authorization_endpoint"`
	TokenEndpoint               string `json:"token_endpoint"`
	EndSessionEndpoint          string `json:"end_session_endpoint"`
}

func initOIDC() {

	var err error

	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	var (
		req *http.Request
		res *http.Response
	)

	logger.Log.Info("- fetching OICD discover: " + env.OIDCDiscoverURL)

	if req, err = http.NewRequest(http.MethodGet, env.OIDCDiscoverURL, nil); err != nil {
		log.Fatal(err)
	}
	if res, err = httpClient.Do(req); err != nil {
		log.Fatal(err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	// Decoding.
	var body []byte
	if body, err = io.ReadAll(res.Body); err != nil {
		log.Fatal(err)
	}

	logger.Log.Info("- OICD body: " + string(body))

	oidcDiscover := OIDCDiscover{}
	if err = json.Unmarshal(body, &oidcDiscover); err != nil {
		log.Fatal(err)
	}

	logger.Log.Info("- OICD token endpoint: " + oidcDiscover.TokenEndpoint)

	// Creating new OIDC provider.
	if env.OIDCProvider, err = oidc.NewProvider(context.Background(), oidcDiscover.Issuer); err != nil {
		panic(err)
	}

	env.OIDCEndSessionEndpoint = oidcDiscover.EndSessionEndpoint
	env.OIDCConfig = &oidc.Config{
		ClientID: env.OIDCClientID,
	}
	env.OIDCVerifier = env.OIDCProvider.Verifier(env.OIDCConfig)
	env.OAuth2Config = oauth2.Config{
		ClientID:     env.OIDCClientID,
		ClientSecret: env.OIDCClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL:      oidcDiscover.TokenEndpoint,
			DeviceAuthURL: oidcDiscover.DeviceAuthorizationEndpoint,
			AuthURL:       oidcDiscover.AuthorizationEndpoint,
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

	logger.Log.Info("- updating GHS statements")
	zmqclient.DBUpdateGHSStatements()

	env.DB = datastore

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
	// http.Handle("/wasm/", alice.New(env.HeadersMiddleware).Then(http.FileServer(http.FS(embedWasmBox))))
	// http.Handle("/static/", alice.New(env.HeadersMiddleware).Then(http.FileServer(http.FS(embedStaticBox))))
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

	router := buildEndpoints(*paramFakeAuth)

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

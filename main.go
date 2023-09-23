//go:build go1.21 && linux && amd64

//go:generate jade -writer -basedir static/templates -d ./static/jade welcomeannounce/index.jade home/index.jade login/index.jade about/index.jade entity/index.jade entity/create.jade product/index.jade product/create.jade storage/index.jade storage/create.jade storelocation/index.jade storelocation/create.jade person/index.jade person/create.jade person/password.jade person/qrcode.jade search.jade menu.jade
//go:generate go run . -genlocalejs
package main

// compile with:
// BuildID="v2.0.8" && go build -ldflags "-X main.BuildID=$BuildID".
import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/casbin"
	"github.com/tbellembois/gochimitheque/datastores"
	"github.com/tbellembois/gochimitheque/handlers"
	"github.com/tbellembois/gochimitheque/ldap"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/mailer"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/static/localejs"
)

var (
	env handlers.Env

	// Starting parameters and commands.
	paramDBPath,
	paramAdminList,
	paramLogFile,
	commandImportFrom,
	commandMailTest,
	commandLDAPSearchUserTest,
	commandLDAPSearchGroupTest *string
	paramPublicProductsEndpoint,
	commandResetAdminPassword,
	commandUpdateQRCode,
	paramDebug,
	commandVersion,
	commandGenLocaleJS,
	paramDisableCache *bool
	BuildID string

	//go:embed wasm/*
	embedWasmBox embed.FS
	//go:embed static/*
	embedStaticBox embed.FS
)

// TimeTrack displays the run time of the function "name"
// from the start time "start"
// use: defer utils.TimeTrack(time.Now(), "GetProducts")
// at the beginning of the function to track
// func TimeTrack(start time.Time, name string) {
// 	elapsed := time.Since(start)
// 	logger.Log.Debug(fmt.Sprintf("%s took %s", name, elapsed))
// }

func init() {
	env = handlers.NewEnv()

	// Configuration parameters.
	flagDBPath := flag.String("dbpath", "./", "the application sqlite directory path")
	flagAppURL := flag.String("appurl", "http://localhost:8081", "the application url (without the path), with NO trailing /")
	flagAppPath := flag.String("apppath", "/", "the application path with the trailing /")
	flagDockerPort := flag.Int("dockerport", 0, "application listen port while running in docker")

	flagMailServerAddress := flag.String("mailserveraddress", "localhost", "the mail server address")
	flagMailServerPort := flag.String("mailserverport", "25", "the SMTP server port")
	flagMailServerSender := flag.String("mailserversender", "", "the SMTP server sender")
	flagMailServerUseTLS := flag.Bool("mailserverusetls", false, "use SMTP TLS? (optional)")
	flagMailServerTLSSkipVerify := flag.Bool("mailservertlsskipverify", false, "skip SMTP TLS verification? (optional)")
	flagPublicProductsEndpoint := flag.Bool("enablepublicproductsendpoint", false, "enable public products endpoint (optional)")

	flagLDAPServerURL := flag.String("ldapserverurl", "", "the LDAP server address - ex: ldaps://192.168.1.50:636/ou=users,dc=foo,dc=local")
	flagLDAPServerUsername := flag.String("ldapserverusername", "", "the LDAP server username - ex: CN=adminro,OU=FOO,OU=local,OU=users,DC=foo,DC=local")
	flagLDAPServerPassword := flag.String("ldapserverpassword", "", "the LDAP server password")
	flagLDAPGroupSearchBaseDN := flag.String("ldapgroupsearchbasedn", "", "the LDAP group search base DN - ex: OU=groups,DC=foo,DC=local")
	flagLDAPGroupSearchFilter := flag.String("ldapgroupsearchfilter", "", "the LDAP group search filter - ex: (cn=%s)")
	flagLDAPUserSearchBaseDN := flag.String("ldapusersearchbasedn", "", "the LDAP user search base DN - ex: OU=users,DC=foo,DC=local")
	flagLDAPUserSearchFilter := flag.String("ldapusersearchfilter", "", "the LDAP user search filter - ex: (&(mail=%s)(objectclass=user))")
	flagAutoCreateUser := flag.Bool("autocreateuser", false, "auto create user if proxy authentication is used")

	flagAdminList := flag.String("admins", "", "the additional admins (comma separated email adresses) (optional) ")
	flagLogFile := flag.String("logfile", "", "log to the given file (optional)")
	flagDebug := flag.Bool("debug", false, "debug (verbose log), default is error")
	flagDisableCache := flag.Bool("disablecache", false, "disable the cache (development only)")

	// One shot commands.
	flagResetAdminPassword := flag.Bool("resetadminpassword", false, "reset the admin password to `chimitheque`")
	flagUpdateQRCode := flag.Bool("updateqrcode", false, "regenerate storages QR codes")
	flagVersion := flag.Bool("version", false, "display application version")
	flagImportFrom := flag.String("importfrom", "", "base URL of the external Chimithèque instance (running with -enablepublicproductsendpoint) to import products from")
	flagGenLocaleJS := flag.Bool("genlocalejs", false, "generate JS locales (developper target)")

	flagMailTest := flag.String("mailtest", "", "send a test mail")
	flagLDAPSearchUserTest := flag.String("ldapsearchusertest", "", "test an LDAP user search")
	flagLDAPSearchGroupTest := flag.String("ldapsearchgrouptest", "", "test an LDAP group search")

	flag.Parse()

	env.AppURL = *flagAppURL
	env.AppPath = *flagAppPath
	env.DockerPort = *flagDockerPort
	env.AutoCreateUser = *flagAutoCreateUser
	ldap.LDAPServerURL = *flagLDAPServerURL
	ldap.LDAPServerUsername = *flagLDAPServerUsername
	ldap.LDAPServerPassword = *flagLDAPServerPassword
	ldap.LDAPGroupSearchBaseDN = *flagLDAPGroupSearchBaseDN
	ldap.LDAPGroupSearchFilter = *flagLDAPGroupSearchFilter
	ldap.LDAPUserSearchBaseDN = *flagLDAPUserSearchBaseDN
	ldap.LDAPUserSearchFilter = *flagLDAPUserSearchFilter
	paramDBPath = flagDBPath
	mailer.MailServerAddress = *flagMailServerAddress
	mailer.MailServerPort = *flagMailServerPort
	mailer.MailServerSender = *flagMailServerSender
	mailer.MailServerUseTLS = *flagMailServerUseTLS
	mailer.MailServerTLSSkipVerify = *flagMailServerTLSSkipVerify
	paramPublicProductsEndpoint = flagPublicProductsEndpoint
	paramAdminList = flagAdminList
	paramLogFile = flagLogFile
	paramDebug = flagDebug
	paramDisableCache = flagDisableCache

	commandResetAdminPassword = flagResetAdminPassword
	commandUpdateQRCode = flagUpdateQRCode
	commandVersion = flagVersion
	commandImportFrom = flagImportFrom
	commandGenLocaleJS = flagGenLocaleJS

	commandMailTest = flagMailTest
	commandLDAPSearchUserTest = flagLDAPSearchUserTest
	commandLDAPSearchGroupTest = flagLDAPSearchGroupTest

	env.AppFullURL = env.AppURL + env.AppPath
	env.BuildID = BuildID
}

func initLogger() {
	var err error

	if *paramDebug {
		logger.Log.SetLevel(logrus.DebugLevel)
	} else {
		logger.Log.SetLevel(logrus.InfoLevel)
	}

	if *paramLogFile != "" {
		var commandLineLogFile *os.File

		if commandLineLogFile, err = os.OpenFile(*paramLogFile, os.O_WRONLY|os.O_CREATE, 0o755); err != nil {
			logger.Log.Fatal(err)
		} else {
			logger.Log.SetOutput(commandLineLogFile)
		}
	}

	var internalServerErrorLogFile *os.File

	if internalServerErrorLogFile, err = os.OpenFile(path.Join(*paramDBPath, "server_errors.log"), os.O_WRONLY|os.O_CREATE, 0o755); err != nil {
		logger.Log.Fatal(err)
	} else {
		logger.LogInternal.SetOutput(internalServerErrorLogFile)
		logger.LogInternal.SetReportCaller(true)
	}
}

func initLDAP() {

	var err error

	if env.LDAPConnection, err = ldap.Connect(); err != nil {
		logger.Log.Fatal(err)
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
					logger.Log.Fatal("user " + ca + " not found in database")
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

// @title Chimithèque API
// @version 2.0
// @description Chemical product management application.
// @contact.name Thomas Bellembois
// @contact.url https://github.com/tbellembois
// @contact.email thomas.bellembois@gmail.com
// @license.name GNU General Public License v3.0
// @license.url https://www.gnu.org/licenses/gpl-3.0.html
// @host localhost:8081
// @BasePath /.
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

	logger.Log.WithFields(logrus.Fields{
		"commandResetAdminPassword": commandResetAdminPassword,
		"commandUpdateQRCode":       commandUpdateQRCode,
		"commandVersion":            commandVersion,
		"commandMailTest":           commandMailTest,
		"commandImportFrom":         commandImportFrom,
		"commandGenLocaleJS":        commandGenLocaleJS,
	}).Debug("main")

	logger.Log.Debugf("- env: %+v", env)
	logger.Log.Info("- application version: " + env.BuildID)
	logger.Log.Info("- application endpoint: " + env.AppFullURL)

	initDB()

	initLDAP()

	// Advanced commands.
	if *commandImportFrom != "" {
		logger.Log.Info("- import from URL into database")
		err := env.DB.Import(*commandImportFrom)
		if err != nil {
			logger.Log.Error("an error occurred: " + err.Error())
			os.Exit(1)
		}

		os.Exit(0)
	}

	if *commandResetAdminPassword {
		logger.Log.Info("- reseting admin password to `chimitheque`")
		a, err := env.DB.GetPersonByEmail("admin@chimitheque.fr")
		if err != nil {
			logger.Log.Error("an error occurred: " + err.Error())
			os.Exit(1)
		}
		a.PersonPassword = "chimitheque"
		err = env.DB.UpdatePersonPassword(a)
		if err != nil {
			logger.Log.Error("an error occurred: " + err.Error())
			os.Exit(1)
		}

		os.Exit(0)
	}

	if *commandUpdateQRCode {
		logger.Log.Info("- updating storages QR codes")
		err := env.DB.UpdateAllQRCodes()
		if err != nil {
			logger.Log.Error("an error occurred: " + err.Error())
			os.Exit(1)
		}

		os.Exit(0)
	}

	if *commandMailTest != "" {
		logger.Log.Info("- sending a mail to " + *commandMailTest)
		err := mailer.TestMail(*commandMailTest)
		if err != nil {
			logger.Log.Error("an error occurred: " + err.Error())
			os.Exit(1)
		}

		os.Exit(0)
	}

	if *commandLDAPSearchUserTest != "" {
		logger.Log.Info("- searching user in LDAP: " + *commandLDAPSearchUserTest)
		result, err := ldap.TestSearchUser(*commandLDAPSearchUserTest)
		if err != nil {
			logger.Log.Error("an error occurred: " + err.Error())
			os.Exit(1)
		}

		if result.NbResults > 0 {
			logger.Log.Info(fmt.Sprintf("%+v", result.R.Entries[0]))
			for _, a := range result.R.Entries[0].Attributes {
				logger.Log.Info(fmt.Sprintf("- %s: %s", a.Name, a.Values))
			}
		}

		os.Exit(0)
	}

	if *commandLDAPSearchGroupTest != "" {
		logger.Log.Info("- searching group in LDAP: " + *commandLDAPSearchGroupTest)
		result, err := ldap.TestSearchGroup(*commandLDAPSearchGroupTest)
		if err != nil {
			logger.Log.Error("an error occurred: " + err.Error())
			os.Exit(1)
		}

		if result.NbResults > 0 {
			for _, a := range result.R.Entries {
				logger.Log.Info(fmt.Sprintf("- %v", a))
			}
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

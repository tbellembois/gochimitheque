//go:build go1.24 && linux && amd64

//go:generate jade -writer -basedir static/templates -d ./static/jade welcomeannounce/index.jade home/index.jade login/index.jade about/index.jade entity/index.jade entity/create.jade product/index.jade product/create.jade product/pubchem.jade storage/index.jade storage/create.jade storelocation/index.jade storelocation/create.jade person/index.jade person/create.jade person/password.jade person/qrcode.jade search.jade menu.jade
//go:generate go run . -genlocalejs
package main

// compile with:
// BuildID="v2.1.0" && go build -ldflags "-X main.BuildID=$BuildID".
import (
	"embed"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/handlers"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/static/localejs"
)

var (
	env handlers.Env

	// Starting parameters and commands.
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
	flagAppURL := flag.String("appurl", "https://192.168.1.56:8443", "the application url (without the path), with NO trailing /")
	flagDebug := flag.Bool("debug", false, "debug (verbose log), default is error")

	// One shot commands.
	flagVersion := flag.Bool("version", false, "display application version")
	flagGenLocaleJS := flag.Bool("genlocalejs", false, "generate JS locales (developper target)")

	flag.Parse()

	env.AppURL = *flagAppURL
	env.BuildID = BuildID

	paramDebug = flagDebug

	commandVersion = flagVersion
	commandGenLocaleJS = flagGenLocaleJS

}

func initLogger() {
	if *paramDebug {
		logger.Log.SetLevel(logrus.DebugLevel)
	} else {
		logger.Log.SetLevel(logrus.InfoLevel)
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

	logger.Log.WithFields(logrus.Fields{
		"commandVersion":     commandVersion,
		"commandGenLocaleJS": commandGenLocaleJS,
	}).Debug("main")

	router := buildEndpoints()

	initStaticResources(router)

	logger.Log.Infof("- env: %+v", env)
	logger.Log.Info("- application version: " + env.BuildID)

	if err = http.ListenAndServe(":8081", nil); err != nil {
		panic("error running the server:" + err.Error())
	}
}

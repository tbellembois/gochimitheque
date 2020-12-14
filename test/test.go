package test

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/casbin/casbin/v2"
	jsonadapter "github.com/casbin/json-adapter/v2"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	. "github.com/tbellembois/gochimitheque/datastores"
	"github.com/tbellembois/gochimitheque/globals"
	"github.com/tbellembois/gochimitheque/handlers"
	. "github.com/tbellembois/gochimitheque/models"

	"github.com/justinas/alice"
)

var (
	TDatastore                                   Datastore    // test datastore
	TEnv                                         handlers.Env // test environment
	Admin                                        Person       // application default admin
	Mickey1, Pluto1, Mickey2, Pluto2, Man1, Man2 Person       // entity members and managers
	E1, E2                                       Entity
	SL1a, SL2a, SL1b, SL2b                       StoreLocation
	Pa, Pb                                       Product
	S1a, S1b, S2a, S2b                           Storage

	err       error
	dbnameOri = "./storage_test.db.ori"
	dbname    = "./storage_test.db"
	from, to  *os.File
	wd        string
)

func TestOut() {
	if err = os.Remove(dbname); err != nil {
		panic(err)
	}
	if err = os.Remove(dbname + "-shm"); err != nil {
		panic(err)
	}
	if err = os.Remove(dbname + "-wal"); err != nil {
		panic(err)
	}
}

func TestInit() {
	globals.Log = logrus.New()
	globals.Log.SetLevel(logrus.DebugLevel)

	// getting original test db full path
	if wd, err = os.Getwd(); err != nil {
		panic(err)
	}
	dbfullpathOri := path.Join(wd, "..", "test", dbnameOri)

	// database original file copy
	if from, err = os.Open(dbfullpathOri); err != nil {
		panic(err)
	}
	defer from.Close()

	if to, err = os.OpenFile(dbname, os.O_RDWR|os.O_CREATE, 0666); err != nil {
		panic(err)
	}
	defer to.Close()

	if _, err = io.Copy(to, from); err != nil {
		panic(err)
	}

	// databse creation
	if TDatastore, err = NewSQLiteDBstore(dbname); err != nil {
		panic(err)
	}

	// gathering test objects
	log.Println("gathering test objects")
	if Admin, err = TDatastore.GetPersonByEmail("admin@chimitheque.fr"); err != nil {
		panic(err)
	}
	Admin.PersonPassword = "chimitheque"
	if Man1, err = TDatastore.GetPersonByEmail("manager1@test.com"); err != nil {
		panic(err)
	}
	Man1.PersonPassword = "manager1"
	if Man2, err = TDatastore.GetPersonByEmail("manager2@test.com"); err != nil {
		panic(err)
	}
	Man2.PersonPassword = "manager2"
	if Mickey1, err = TDatastore.GetPersonByEmail("person1a@test.com"); err != nil {
		panic(err)
	}
	Mickey1.PersonPassword = "person1a"
	if Pluto1, err = TDatastore.GetPersonByEmail("person1b@test.com"); err != nil {
		panic(err)
	}
	Pluto1.PersonPassword = "person1b"
	if Mickey2, err = TDatastore.GetPersonByEmail("person2a@test.com"); err != nil {
		panic(err)
	}
	Mickey2.PersonPassword = "person2a"
	if Pluto2, err = TDatastore.GetPersonByEmail("person2b@test.com"); err != nil {
		panic(err)
	}
	Pluto2.PersonPassword = "person2b"

	if E1, err = TDatastore.GetEntity(2); err != nil {
		panic(err)
	}
	if E2, err = TDatastore.GetEntity(3); err != nil {
		panic(err)
	}

	if S1a, err = TDatastore.GetStorage(1); err != nil {
		panic(err)
	}
	if S1b, err = TDatastore.GetStorage(2); err != nil {
		panic(err)
	}
	if S2a, err = TDatastore.GetStorage(3); err != nil {
		panic(err)
	}
	if S2b, err = TDatastore.GetStorage(4); err != nil {
		panic(err)
	}

	if SL1a, err = TDatastore.GetStoreLocation(1); err != nil {
		panic(err)
	}
	if SL1b, err = TDatastore.GetStoreLocation(2); err != nil {
		panic(err)
	}
	if SL2a, err = TDatastore.GetStoreLocation(3); err != nil {
		panic(err)
	}
	if SL2b, err = TDatastore.GetStoreLocation(4); err != nil {
		panic(err)
	}

	// TEnvironment creation
	TEnv = handlers.Env{
		DB: TDatastore,
	}

	// casbin enforcer creation
	if globals.JSONAdapterData, err = TDatastore.ToCasbinJSONAdapter(); err != nil {
		globals.Log.Error("error getting json adapter data: " + err.Error())
		os.Exit(1)
	}
	a := jsonadapter.NewAdapter(&globals.JSONAdapterData)
	if globals.Enforcer, err = casbin.NewEnforcer("../model.conf", a); err != nil {
		globals.Log.Error("enforcer creation error: " + err.Error())
		os.Exit(1)
	}
	globals.Enforcer.AddFunction("matchStorage", TEnv.MatchStorageFunc)
	globals.Enforcer.AddFunction("matchStorelocation", TEnv.MatchStorelocationFunc)
	globals.Enforcer.AddFunction("matchPeople", TEnv.MatchPeopleFunc)
	globals.Enforcer.AddFunction("matchEntity", TEnv.MatchEntityFunc)

	// enforccer policy load
	if err = globals.Enforcer.LoadPolicy(); err != nil {
		globals.Log.Error("enforcer policy load error: " + err.Error())
		os.Exit(1)
	}

}

func Authenticate(pemail, ppassword string) (string, error) {

	// middleware chain
	commonChain := alice.New(TEnv.ContextMiddleware, TEnv.HeadersMiddleware, TEnv.LogingMiddleware)

	// requests definition
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	r.Handle("/get-token", commonChain.Then(TEnv.AppMiddleware(TEnv.GetTokenHandler))).Methods("POST")

	// performing the request
	req, err := http.NewRequest("POST", "/get-token", strings.NewReader(url.Values{"person_email": {pemail}, "person_password": {ppassword}}.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(rr, req)

	// returning the token
	return rr.Body.String(), nil
}

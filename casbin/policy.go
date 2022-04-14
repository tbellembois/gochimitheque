package casbin

import (
	_ "embed"
	"os"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	jsonadapter "github.com/casbin/json-adapter/v2"
	"github.com/tbellembois/gochimitheque/datastores"
	"github.com/tbellembois/gochimitheque/logger"
)

var (
	//go:embed policy.conf
	embedModel string
)

func InitCasbinPolicy(datastore datastores.Datastore) (enforcer *casbin.Enforcer) {
	var (
		jsonAdapterData []byte
		err             error
	)

	if jsonAdapterData, err = datastore.ToCasbinJSONAdapter(); err != nil {
		logger.Log.Error("error getting json adapter data: " + err.Error())
		os.Exit(1)
	}

	var (
		m model.Model
	)

	if m, err = model.NewModelFromString(embedModel); err != nil {
		logger.Log.Error("model creation error: " + err.Error())
		os.Exit(1)
	}

	a := jsonadapter.NewAdapter(&jsonAdapterData)
	if enforcer, err = casbin.NewEnforcer(m, a); err != nil {
		logger.Log.Error("enforcer creation error: " + err.Error())
		os.Exit(1)
	}

	enforcer.AddFunction("matchStorage", MatchStorageFuncWrapper(datastore))
	enforcer.AddFunction("matchStorelocation", MatchStorelocationFuncWrapper(datastore))
	enforcer.AddFunction("matchPeople", MatchPeopleFuncWrapper(datastore))
	enforcer.AddFunction("matchEntity", MatchEntityFuncWrapper(datastore))

	if err = enforcer.LoadPolicy(); err != nil {
		logger.Log.Error("enforcer policy load error: " + err.Error())
		os.Exit(1)
	}

	return
}

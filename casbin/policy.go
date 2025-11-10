package casbin

import (
	_ "embed"

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
		logger.Log.Fatal("error getting json adapter data: " + err.Error())
	}

	// logger.Log.WithFields(logrus.Fields{
	// 	"jsonAdapterData": string(jsonAdapterData),
	// }).Debug("InitCasbinPolicy")

	var (
		m model.Model
	)

	if m, err = model.NewModelFromString(embedModel); err != nil {
		logger.Log.Fatal("model creation error: " + err.Error())
	}

	a := jsonadapter.NewAdapter(&jsonAdapterData)
	if enforcer, err = casbin.NewEnforcer(m, a); err != nil {
		logger.Log.Fatal("enforcer creation error: " + err.Error())
	}

	enforcer.AddFunction("matchProductHasStorages", MatchProductHasStoragesFuncWrapper())
	enforcer.AddFunction("matchPersonIsInPersonEntity", MatchPersonIsInEntityFuncWrapper())
	enforcer.AddFunction("matchPersonIsInStorageEntity", MatchPersonIsInStorageEntityFuncWrapper())
	enforcer.AddFunction("matchPersonIsInStoreLocationEntity", MatchPersonIsInStoreLocationEntityFuncWrapper())
	enforcer.AddFunction("matchEntityHasMembers", MatchEntityHasMembersFuncWrapper())
	enforcer.AddFunction("matchEntityHasStoreLocations", MatchEntityHasStoreLocationsFuncWrapper())
	enforcer.AddFunction("matchPersonIsAdmin", MatchPersonIsAdminFuncWrapper())
	enforcer.AddFunction("matchPersonIsInEntity", MatchPersonIsInEntityFuncWrapper())
	enforcer.AddFunction("matchProductHasStorages", MatchProductHasStoragesFuncWrapper())
	enforcer.AddFunction("matchStoreLocationHasChildren", MatchStoreLocationHasChildrenFuncWrapper())
	enforcer.AddFunction("matchStoreLocationHasStorages", MatchStoreLocationHasStoragesFuncWrapper())
	enforcer.AddFunction("matchPersonIsManager", MatchPersonIsManagerFuncWrapper())
	enforcer.AddFunction("matchPersonIsInStorageEntity", MatchPersonIsInStorageEntityFuncWrapper())

	if err = enforcer.LoadPolicy(); err != nil {
		logger.Log.Fatal("enforcer policy load error: " + err.Error())
	}

	return
}

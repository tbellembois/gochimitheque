package casbin

import (
	"encoding/json"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/zmqclient"
)

func matchStoreLocationIsInEntity(storeLocationId int64, entityId int64) bool {

	var (
		err            error
		jsonRawMessage json.RawMessage
		result         bool
	)

	logger.Log.WithFields(logrus.Fields{"storageId": storeLocationId, "entityId": entityId}).Debug("matchStoreLocationIsInEntity")

	if jsonRawMessage, err = zmqclient.CasbinMatchStoreLocationIsInEntity(storeLocationId, entityId); err != nil {
		logger.Log.Error("CasbinMatchStoreLocationIsInEntity: " + err.Error())
		return false
	}

	if result, err = zmqclient.ConvertDBJSONToBool(jsonRawMessage); err != nil {
		logger.Log.Error("CasbinMatchStoreLocationIsInEntity: " + err.Error())
		return false
	}

	logger.Log.WithFields(logrus.Fields{"result": result}).Debug("matchStoreLocationIsInEntity")

	return result
}

func MatchStoreLocationIsInEntityFuncWrapper() func(args ...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		storeLocationID := args[0].(string)
		itemID := args[1].(string)

		var (
			storeLocationID_int64 int64
			itemID_int64          int64
			err                   error
		)

		if storeLocationID_int64, err = strconv.ParseInt(storeLocationID, 10, 64); err != nil {
			return false, err
		}
		if itemID_int64, err = strconv.ParseInt(itemID, 10, 64); err != nil {
			return false, err
		}

		return (bool)(matchStoreLocationIsInEntity(storeLocationID_int64, itemID_int64)), nil
	}
}
func matchStorageIsInEntity(storageId int64, entityId int64) bool {

	var (
		err            error
		jsonRawMessage json.RawMessage
		result         bool
	)

	logger.Log.WithFields(logrus.Fields{"storageId": storageId, "entityId": entityId}).Debug("matchStorageIsInEntity")

	if jsonRawMessage, err = zmqclient.CasbinMatchStorageIsInEntity(storageId, entityId); err != nil {
		logger.Log.Error("CasbinMatchStorageIsInEntity: " + err.Error())
		return false
	}

	if result, err = zmqclient.ConvertDBJSONToBool(jsonRawMessage); err != nil {
		logger.Log.Error("CasbinMatchStorageIsInEntity: " + err.Error())
		return false
	}

	logger.Log.WithFields(logrus.Fields{"result": result}).Debug("matchStorageIsInEntity")

	return result
}

func MatchStorageisInEntityFuncWrapper() func(args ...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		storageID := args[0].(string)
		itemID := args[1].(string)

		var (
			storageID_int64 int64
			itemID_int64    int64
			err             error
		)

		if storageID_int64, err = strconv.ParseInt(storageID, 10, 64); err != nil {
			return false, err
		}
		if itemID_int64, err = strconv.ParseInt(itemID, 10, 64); err != nil {
			return false, err
		}

		return (bool)(matchStorageIsInEntity(storageID_int64, itemID_int64)), nil
	}
}

func matchPersonIsInPersonEntity(personId int64, otherPersonId int64) bool {

	var (
		err            error
		jsonRawMessage json.RawMessage
		result         bool
	)

	logger.Log.WithFields(logrus.Fields{"personId": personId, "otherPersonId": otherPersonId}).Debug("matchPersonIsInPersonEntity")

	if jsonRawMessage, err = zmqclient.CasbinMatchPersonIsInPersonEntity(personId, otherPersonId); err != nil {
		logger.Log.Error("CasbinMatchPersonIsInPersonEntity: " + err.Error())
		return false
	}

	if result, err = zmqclient.ConvertDBJSONToBool(jsonRawMessage); err != nil {
		logger.Log.Error("CasbinMatchPersonIsInPersonEntity: " + err.Error())
		return false
	}

	logger.Log.WithFields(logrus.Fields{"result": result}).Debug("matchProductHasStorages")

	return result
}

func MatchPersonIsInPersonEntityFuncWrapper() func(args ...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		personID := args[0].(string)
		itemID := args[1].(string)

		var (
			personID_int64 int64
			itemID_int64   int64
			err            error
		)

		if personID_int64, err = strconv.ParseInt(personID, 10, 64); err != nil {
			return false, err
		}
		if itemID_int64, err = strconv.ParseInt(itemID, 10, 64); err != nil {
			return false, err
		}

		return (bool)(matchPersonIsInPersonEntity(personID_int64, itemID_int64)), nil
	}
}

func matchPersonIsInStorageEntity(personId int64, storageId int64) bool {
	var (
		err            error
		jsonRawMessage json.RawMessage
		result         bool
	)

	logger.Log.WithFields(logrus.Fields{"personId": personId, "storageId": storageId}).Debug("matchPersonIsInStorageEntity")

	if jsonRawMessage, err = zmqclient.CasbinMatchPersonIsInStorageEntity(personId, storageId); err != nil {
		logger.Log.Error("CasbinMatchPersonIsInStorageEntity: " + err.Error())
		return false
	}

	if result, err = zmqclient.ConvertDBJSONToBool(jsonRawMessage); err != nil {
		logger.Log.Error("CasbinMatchPersonIsInStorageEntity: " + err.Error())
		return false
	}

	logger.Log.WithFields(logrus.Fields{"result": result}).Debug("matchProductHasStorages")

	return result
}

func MatchPersonIsInStorageEntityFuncWrapper() func(args ...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		personID := args[0].(string)
		itemID := args[1].(string)

		var (
			personID_int64 int64
			itemID_int64   int64
			err            error
		)

		if personID_int64, err = strconv.ParseInt(personID, 10, 64); err != nil {
			return false, err
		}
		if itemID_int64, err = strconv.ParseInt(itemID, 10, 64); err != nil {
			return false, err
		}

		return (bool)(matchPersonIsInStorageEntity(personID_int64, itemID_int64)), nil
	}
}

func matchPersonIsInStoreLocationEntity(personId int64, storeLocationId int64) bool {
	var (
		err            error
		jsonRawMessage json.RawMessage
		result         bool
	)

	logger.Log.WithFields(logrus.Fields{"personId": personId, "storeLocationId": storeLocationId}).Debug("matchPersonIsInStoreLocationEntity")

	if jsonRawMessage, err = zmqclient.CasbinMatchPersonIsInStoreLocationEntity(personId, storeLocationId); err != nil {
		logger.Log.Error("CasbinMatchPersonIsInStoreLocationEntity: " + err.Error())
		return false
	}

	if result, err = zmqclient.ConvertDBJSONToBool(jsonRawMessage); err != nil {
		logger.Log.Error("CasbinMatchPersonIsInStoreLocationEntity: " + err.Error())
		return false
	}

	logger.Log.WithFields(logrus.Fields{"result": result}).Debug("matchProductHasStorages")

	return result
}

func MatchPersonIsInStoreLocationEntityFuncWrapper() func(args ...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		personID := args[0].(string)
		itemID := args[1].(string)

		var (
			personID_int64 int64
			itemID_int64   int64
			err            error
		)

		if personID_int64, err = strconv.ParseInt(personID, 10, 64); err != nil {
			return false, err
		}
		if itemID_int64, err = strconv.ParseInt(itemID, 10, 64); err != nil {
			return false, err
		}

		return (bool)(matchPersonIsInStoreLocationEntity(personID_int64, itemID_int64)), nil
	}
}

func matchEntityHasMembers(entityId int64) bool {
	var (
		err            error
		jsonRawMessage json.RawMessage
		result         bool
	)

	logger.Log.WithFields(logrus.Fields{"entityId": entityId}).Debug("matchEntityHasMembers")

	if jsonRawMessage, err = zmqclient.CasbinMatchEntityHasMembers(entityId); err != nil {
		logger.Log.Error("CasbinMatchEntityHasMembers: " + err.Error())
		return false
	}

	if result, err = zmqclient.ConvertDBJSONToBool(jsonRawMessage); err != nil {
		logger.Log.Error("CasbinMatchEntityHasMembers: " + err.Error())
		return false
	}

	logger.Log.WithFields(logrus.Fields{"result": result}).Debug("matchProductHasStorages")

	return result
}

func MatchEntityHasMembersFuncWrapper() func(args ...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		entityId := args[0].(string)

		var (
			entityId_int64 int64
			err            error
		)

		if entityId_int64, err = strconv.ParseInt(entityId, 10, 64); err != nil {
			return false, err
		}

		return (bool)(matchEntityHasMembers(entityId_int64)), nil
	}
}

func matchEntityHasStoreLocations(entityId int64) bool {
	var (
		err            error
		jsonRawMessage json.RawMessage
		result         bool
	)

	logger.Log.WithFields(logrus.Fields{"entityId": entityId}).Debug("matchEntityHasStoreLocations")

	if jsonRawMessage, err = zmqclient.CasbinMatchEntityHasStoreLocations(entityId); err != nil {
		logger.Log.Error("CasbinMatchEntityHasStoreLocations: " + err.Error())
		return false
	}

	if result, err = zmqclient.ConvertDBJSONToBool(jsonRawMessage); err != nil {
		logger.Log.Error("CasbinMatchEntityHasStoreLocations: " + err.Error())
		return false
	}

	logger.Log.WithFields(logrus.Fields{"result": result}).Debug("matchProductHasStorages")

	return result
}

func MatchEntityHasStoreLocationsFuncWrapper() func(args ...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		entityId := args[0].(string)

		var (
			entityId_int64 int64
			err            error
		)

		if entityId_int64, err = strconv.ParseInt(entityId, 10, 64); err != nil {
			return false, err
		}

		return (bool)(matchEntityHasStoreLocations(entityId_int64)), nil
	}
}

func matchPersonIsAdmin(personId int64) bool {
	var (
		err            error
		jsonRawMessage json.RawMessage
		result         bool
	)

	logger.Log.WithFields(logrus.Fields{"personId": personId}).Debug("matchPersonIsAdmin")

	if jsonRawMessage, err = zmqclient.CasbinMatchPersonIsAdmin(personId); err != nil {
		logger.Log.Error("CasbinMatchPersonIsAdmin: " + err.Error())
		return false
	}

	if result, err = zmqclient.ConvertDBJSONToBool(jsonRawMessage); err != nil {
		logger.Log.Error("CasbinMatchPersonIsAdmin: " + err.Error())
		return false
	}

	logger.Log.WithFields(logrus.Fields{"result": result}).Debug("matchProductHasStorages")

	return result
}

func MatchPersonIsAdminFuncWrapper() func(args ...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		personId := args[0].(string)

		var (
			personId_int64 int64
			err            error
		)

		if personId_int64, err = strconv.ParseInt(personId, 10, 64); err != nil {
			return false, err
		}

		return (bool)(matchPersonIsAdmin(personId_int64)), nil
	}
}

func matchPersonIsInEntity(personId int64, entityId int64) bool {
	var (
		err            error
		jsonRawMessage json.RawMessage
		result         bool
	)

	logger.Log.WithFields(logrus.Fields{"personId": personId, "entityId": entityId}).Debug("matchPersonIsInEntity")

	if jsonRawMessage, err = zmqclient.CasbinMatchPersonIsInEntity(personId, entityId); err != nil {
		logger.Log.Error("CasbinMatchPersonIsInEntity: " + err.Error())
		return false
	}

	if result, err = zmqclient.ConvertDBJSONToBool(jsonRawMessage); err != nil {
		logger.Log.Error("CasbinMatchPersonIsInEntity: " + err.Error())
		return false
	}

	logger.Log.WithFields(logrus.Fields{"result": result}).Debug("matchProductHasStorages")

	return result
}

func MatchPersonIsInEntityFuncWrapper() func(args ...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		personID := args[0].(string)
		entityID := args[1].(string)

		var (
			personID_int64 int64
			entityID_int64 int64
			err            error
		)

		if personID_int64, err = strconv.ParseInt(personID, 10, 64); err != nil {
			return false, err
		}
		if entityID_int64, err = strconv.ParseInt(entityID, 10, 64); err != nil {
			return false, err
		}

		return (bool)(matchPersonIsInEntity(personID_int64, entityID_int64)), nil
	}
}

func matchProductHasStorages(productId int64) bool {

	var (
		err            error
		jsonRawMessage json.RawMessage
		result         bool
	)

	logger.Log.WithFields(logrus.Fields{"productId": productId}).Debug("matchProductHasStorages")

	if jsonRawMessage, err = zmqclient.CasbinMatchProductHasStorages(productId); err != nil {
		logger.Log.Error("CasbinMatchProductHasStorages: " + err.Error())
		return false
	}

	if result, err = zmqclient.ConvertDBJSONToBool(jsonRawMessage); err != nil {
		logger.Log.Error("CasbinMatchProductHasStorages: " + err.Error())
		return false
	}

	logger.Log.WithFields(logrus.Fields{"result": result}).Debug("matchProductHasStorages")

	return result
}

func MatchProductHasStoragesFuncWrapper() func(args ...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		productId := args[0].(string)

		var (
			productId_int64 int64
			err             error
		)

		if productId_int64, err = strconv.ParseInt(productId, 10, 64); err != nil {
			return false, err
		}

		return (bool)(matchProductHasStorages(productId_int64)), nil
	}
}

func matchStoreLocationHasChildren(storeLocationId int64) bool {
	var (
		err            error
		jsonRawMessage json.RawMessage
		result         bool
	)

	logger.Log.WithFields(logrus.Fields{"storeLocationId": storeLocationId}).Debug("matchStoreLocationHasChildren")

	if jsonRawMessage, err = zmqclient.CasbinMatchStoreLocationHasChildren(storeLocationId); err != nil {
		logger.Log.Error("CasbinMatchStoreLocationHasChildren: " + err.Error())
		return false
	}

	if result, err = zmqclient.ConvertDBJSONToBool(jsonRawMessage); err != nil {
		logger.Log.Error("CasbinMatchStoreLocationHasChildren: " + err.Error())
		return false
	}

	logger.Log.WithFields(logrus.Fields{"result": result}).Debug("matchProductHasStorages")

	return result
}

func MatchStoreLocationHasChildrenFuncWrapper() func(args ...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		storeLocationId := args[0].(string)

		var (
			storeLocationId_int64 int64
			err                   error
		)

		if storeLocationId_int64, err = strconv.ParseInt(storeLocationId, 10, 64); err != nil {
			return false, err
		}

		return (bool)(matchStoreLocationHasChildren(storeLocationId_int64)), nil
	}
}

func matchStoreLocationHasStorages(storeLocationId int64) bool {
	var (
		err            error
		jsonRawMessage json.RawMessage
		result         bool
	)

	logger.Log.WithFields(logrus.Fields{"storeLocationId": storeLocationId}).Debug("matchStoreLocationHasStorages")

	if jsonRawMessage, err = zmqclient.CasbinMatchStoreLocationHasStorages(storeLocationId); err != nil {
		logger.Log.Error("CasbinMatchStoreLocationHasStorages: " + err.Error())
		return false
	}

	if result, err = zmqclient.ConvertDBJSONToBool(jsonRawMessage); err != nil {
		logger.Log.Error("CasbinMatchStoreLocationHasStorages: " + err.Error())
		return false
	}

	logger.Log.WithFields(logrus.Fields{"result": result}).Debug("matchProductHasStorages")

	return result
}

func MatchStoreLocationHasStoragesFuncWrapper() func(args ...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		storeLocationId := args[0].(string)

		var (
			storeLocationId_int64 int64
			err                   error
		)

		if storeLocationId_int64, err = strconv.ParseInt(storeLocationId, 10, 64); err != nil {
			return false, err
		}

		return (bool)(matchStoreLocationHasStorages(storeLocationId_int64)), nil
	}
}

func matchPersonIsManager(personId int64) bool {
	var (
		err            error
		jsonRawMessage json.RawMessage
		result         bool
	)

	logger.Log.WithFields(logrus.Fields{"personId": personId}).Debug("matchPersonIsManager")

	if jsonRawMessage, err = zmqclient.CasbinMatchPersonIsManager(personId); err != nil {
		logger.Log.Error("CasbinMatchPersonIsManager: " + err.Error())
		return false
	}

	if result, err = zmqclient.ConvertDBJSONToBool(jsonRawMessage); err != nil {
		logger.Log.Error("CasbinMatchPersonIsManager: " + err.Error())
		return false
	}

	logger.Log.WithFields(logrus.Fields{"result": result}).Debug("matchProductHasStorages")

	return result
}

func MatchPersonIsManagerFuncWrapper() func(args ...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		personId := args[0].(string)

		var (
			personId_int64 int64
			err            error
		)

		if personId_int64, err = strconv.ParseInt(personId, 10, 64); err != nil {
			return false, err
		}

		return (bool)(matchPersonIsManager(personId_int64)), nil
	}
}

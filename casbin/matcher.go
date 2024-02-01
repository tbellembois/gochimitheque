package casbin

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/datastores"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
)

func matchPeople(datastore datastores.Datastore, personID string, itemID string, entityID string) bool {
	var (
		orphan   bool
		pid, iid int
		err      error
		ent      []models.Entity
	)

	logger.Log.WithFields(logrus.Fields{"personId": personID, "itemId": itemID, "entityId": entityID}).Debug("matchPeople")

	if orphan, err = datastore.IsOrphanPerson(iid); err != nil {
		logger.Log.Error("matchPeople: " + err.Error())
		return false
	}
	if orphan {
		return true
	}

	if pid, err = strconv.Atoi(personID); err != nil {
		logger.Log.Error("matchPeople: " + err.Error())
		return false
	}

	if iid, err = strconv.Atoi(itemID); err != nil {
		logger.Log.Error("matchPeople: " + err.Error())
		return false
	}

	if ent, err = datastore.GetPersonEntities(pid, iid); err != nil {
		logger.Log.Error("matchPeople: " + err.Error())
		return false
	}

	found := false

	for _, e := range ent {
		if strconv.Itoa(e.EntityID) == entityID {
			found = true
			continue
		}
	}

	return found
}

func MatchPeopleFuncWrapper(datastore datastores.Datastore) func(args ...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		personID := args[0].(string)
		itemID := args[1].(string)
		entityID := args[2].(string)

		return (bool)(matchPeople(datastore, personID, itemID, entityID)), nil
	}
}

func matchStorelocation(datastore datastores.Datastore, personID string, itemID string, entityID string) bool {
	var (
		pid, iid      int
		err           error
		m             bool
		storelocation models.StoreLocation
	)

	logger.Log.WithFields(logrus.Fields{"personId": personID, "itemId": itemID, "entityId": entityID}).Debug("matchStorelocation")

	if pid, err = strconv.Atoi(personID); err != nil {
		logger.Log.Error("matchStorelocation: " + err.Error())
		return false
	}

	if iid, err = strconv.Atoi(itemID); err != nil {
		logger.Log.Error("matchStorelocation: " + err.Error())
		return false
	}

	if storelocation, err = datastore.GetStoreLocation(iid); err != nil && err != sql.ErrNoRows {
		logger.Log.Error("matchStorelocation: " + err.Error())
		return false
	}

	if err == sql.ErrNoRows {
		return false
	}

	if strconv.Itoa(storelocation.EntityID) != entityID {
		return false
	}

	if m, err = datastore.DoesPersonBelongsTo(pid, []models.Entity{storelocation.Entity}); err != nil {
		logger.Log.Error("matchStorelocation: " + err.Error())
		return false
	}

	logger.Log.WithFields(logrus.Fields{"m": m}).Debug("matchStorelocation")

	return m
}

func MatchStorelocationFuncWrapper(datastore datastores.Datastore) func(args ...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		personID := args[0].(string)
		itemID := args[1].(string)
		entityID := args[2].(string)

		return (bool)(matchStorelocation(datastore, personID, itemID, entityID)), nil
	}
}

func matchStorage(datastore datastores.Datastore, personID string, itemID string, entityID string) bool {
	var (
		pid, iid int
		err      error
		m        bool
		ent      models.Entity
	)

	if pid, err = strconv.Atoi(personID); err != nil {
		logger.Log.Error("matchStorage - pid: " + err.Error())
		return false
	}

	if iid, err = strconv.Atoi(itemID); err != nil {
		logger.Log.Error("matchStorage - iid: " + err.Error())
		return false
	}

	if ent, err = datastore.GetStorageEntity(iid); err != nil {
		logger.Log.Error(fmt.Sprintf("matchStorage: %v %s", ent, err.Error()))
		return false
	}

	if strconv.Itoa(ent.EntityID) != entityID {
		return false
	}

	if m, err = datastore.DoesPersonBelongsTo(pid, []models.Entity{ent}); err != nil {
		logger.Log.Error(fmt.Sprintf("matchStorage: %v %s", ent, err.Error()))
		return false
	}

	logger.Log.WithFields(logrus.Fields{"m": m}).Debug("matchStorage")

	return m
}

func MatchStorageFuncWrapper(datastore datastores.Datastore) func(args ...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		personID := args[0].(string)
		itemID := args[1].(string)
		entityID := args[2].(string)

		return (bool)(matchStorage(datastore, personID, itemID, entityID)), nil
	}
}

func matchEntity(datastore datastores.Datastore, personID string, entityID string) bool {
	var (
		pid, eid int
		err      error
		m        bool
	)

	if pid, err = strconv.Atoi(personID); err != nil {
		logger.Log.Error("matchEntity: " + err.Error())
		return false
	}

	if eid, err = strconv.Atoi(entityID); err != nil {
		logger.Log.Error("matchEntity: " + err.Error())
		return false
	}

	if m, err = datastore.DoesPersonBelongsTo(pid, []models.Entity{{EntityID: eid}}); err != nil {
		logger.Log.Error("matchEntity: " + err.Error())
		return false
	}

	logger.Log.WithFields(logrus.Fields{"personId": personID, "entityId": entityID, "m": m}).Debug("matchEntity")

	return m
}

func MatchEntityFuncWrapper(datastore datastores.Datastore) func(args ...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		personID := args[0].(string)
		entityID := args[1].(string)

		return (bool)(matchEntity(datastore, personID, entityID)), nil
	}
}

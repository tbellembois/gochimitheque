package casbin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/barweiss/go-tuple"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/datastores"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/zmqclient"
)

func matchPeople(datastore datastores.Datastore, personID string, itemID string, entityID string) bool {
	var (
		orphan bool
		err    error
		ent    []*models.Entity
	)

	logger.Log.WithFields(logrus.Fields{"personId": personID, "itemId": itemID, "entityId": entityID}).Debug("matchPeople")

	// TODO: remove 1 by connected user id.
	var (
		jsonRawMessage json.RawMessage
		person         *models.Person
	)

	if jsonRawMessage, err = zmqclient.DBGetPeople("http://localhost/"+personID, 1); err != nil {
		logger.Log.Error("zmqclient.DBGetPeople: " + err.Error())
		return false
	}

	if person, err = zmqclient.ConvertDBJSONToPerson(jsonRawMessage); err != nil {
		logger.Log.Error("ConvertDBJSONToPerson: " + err.Error())
		return false
	}

	orphan = len(person.Entities) == 0

	if orphan {
		return true
	}

	ent = person.Entities

	found := false

	for _, e := range ent {
		if strconv.Itoa(int(*e.EntityID)) == entityID {
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

// Return true if the person with personID and store location with itemID are both in entity with entityID.
func matchStorelocation(datastore datastores.Datastore, personID string, itemID string, entityID string) bool {
	var (
		pid, iid       int
		err            error
		m              bool
		store_location models.StoreLocation
		jsonRawMessage json.RawMessage
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

	// getting the store location matching the id
	if jsonRawMessage, err = zmqclient.DBGetStorelocations("http://localhost/store_locations/"+strconv.Itoa(iid), int64(pid)); err != nil {
		logger.Log.Error("matchStorelocation - error calling zmqclient.DBGetStorelocations: " + err.Error())
		return false
	}

	// unmarshalling response
	var tuple tuple.T2[[]models.StoreLocation, int]
	if err = json.Unmarshal(jsonRawMessage, &tuple); err != nil {
		logger.Log.Error("matchStorelocation - error calling zmqclient.DBGetStorelocations: " + err.Error())
		return false
	}

	store_location = tuple.V1[0]

	if err == sql.ErrNoRows {
		return false
	}

	if strconv.Itoa(int(*store_location.EntityID)) != entityID {
		return false
	}

	var (
		person *models.Person
	)

	if jsonRawMessage, err = zmqclient.DBGetPeople("http://localhost/"+strconv.Itoa(pid), int64(pid)); err != nil {
		logger.Log.Error("zmqclient.DBGetPeople: " + err.Error())
		return false

	}

	if person, err = zmqclient.ConvertDBJSONToPerson(jsonRawMessage); err != nil {
		logger.Log.Error("zmqclient.ConvertDBJSONToPerson: " + err.Error())
		return false
	}

	m = false
	for _, entity := range person.Entities {
		if entity.EntityID == store_location.Entity.EntityID {
			m = true
			break
		}
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

// Return true if the person with personID and storage with itemID are both in entity with entityID.
func matchStorage(datastore datastores.Datastore, personID string, itemID string, entityID string) bool {
	var (
		pid, iid, person_id int
		err                 error
		m                   bool
	)

	if pid, err = strconv.Atoi(personID); err != nil {
		logger.Log.Error("matchStorage - pid: " + err.Error())
		return false
	}

	if iid, err = strconv.Atoi(itemID); err != nil {
		logger.Log.Error("matchStorage - iid: " + err.Error())
		return false
	}

	if person_id, err = strconv.Atoi(personID); err != nil {
		logger.Log.Error("matchStorage - person_id: " + err.Error())
		return false
	}

	var storage *models.Storage

	// getting the storage
	var (
		jsonRawMessage json.RawMessage
	)
	if jsonRawMessage, err = zmqclient.DBGetStorages("http://localhost/"+strconv.Itoa(iid), int64(person_id)); err != nil {
		logger.Log.Error(fmt.Sprintf("zmqclient.DBGetStorages: %s", err.Error()))
	}

	if storage, err = zmqclient.ConvertDBJSONToStorage(jsonRawMessage); err != nil {
		logger.Log.Error(fmt.Sprintf("ConvertDBJSONToStorage: %s", err.Error()))
	}

	// Getting the person.
	var (
		person *models.Person
	)

	if jsonRawMessage, err = zmqclient.DBGetPeople("http://localhost/"+strconv.Itoa(pid), int64(pid)); err != nil {
		logger.Log.Error("zmqclient.DBGetPeople: " + err.Error())
		return false

	}

	if person, err = zmqclient.ConvertDBJSONToPerson(jsonRawMessage); err != nil {
		logger.Log.Error("zmqclient.ConvertDBJSONToPerson: " + err.Error())
		return false
	}

	logger.Log.WithFields(logrus.Fields{"matchStorage - person": person}).Debug("matchStorage")

	m = false
	for _, entity := range person.Entities {
		if entity.EntityID == storage.EntityID {
			m = true
			break
		}
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

	// Getting the person.
	var (
		person         *models.Person
		jsonRawMessage json.RawMessage
	)

	if jsonRawMessage, err = zmqclient.DBGetPeople("http://localhost/"+strconv.Itoa(pid), int64(pid)); err != nil {
		logger.Log.Error("zmqclient.DBGetPeople: " + err.Error())
		return false

	}

	if person, err = zmqclient.ConvertDBJSONToPerson(jsonRawMessage); err != nil {
		logger.Log.Error("zmqclient.ConvertDBJSONToPerson: " + err.Error())
		return false
	}

	// Finding is the person belongs to the entity.
	m = false
	for _, entity := range person.Entities {
		if *entity.EntityID == int64(eid) {
			m = true
			break
		}
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

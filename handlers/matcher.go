package handlers

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
)

func (env *Env) matchPeople(personId string, itemId string, entityId string) bool {
	var (
		pid, iid int
		err      error
		ent      []models.Entity
	)

	if pid, err = strconv.Atoi(personId); err != nil {
		logger.Log.Error("matchPeople: " + err.Error())
		return false
	}
	if iid, err = strconv.Atoi(itemId); err != nil {
		logger.Log.Error("matchPeople: " + err.Error())
		return false
	}

	if ent, err = env.DB.GetPersonEntities(pid, iid); err != nil {
		logger.Log.Error("matchPeople: " + err.Error())
		return false
	}
	found := false
	for _, e := range ent {
		if strconv.Itoa(e.EntityID) == entityId {
			found = true
			continue
		}
	}
	return found
}

func (env *Env) MatchPeopleFunc(args ...interface{}) (interface{}, error) {
	personId := args[0].(string)
	itemId := args[1].(string)
	entityId := args[2].(string)

	return (bool)(env.matchPeople(personId, itemId, entityId)), nil
}

func (env *Env) matchStorelocation(personId string, itemId string, entityId string) bool {
	var (
		pid, iid      int
		err           error
		m             bool
		storelocation models.StoreLocation
	)
	logger.Log.WithFields(logrus.Fields{"personId": personId, "itemId": itemId, "entityId": entityId}).Debug("matchStorelocation")

	if pid, err = strconv.Atoi(personId); err != nil {
		logger.Log.Error("matchStorelocation: " + err.Error())
		return false
	}
	if iid, err = strconv.Atoi(itemId); err != nil {
		logger.Log.Error("matchStorelocation: " + err.Error())
		return false
	}
	if storelocation, err = env.DB.GetStoreLocation(iid); err != nil && err != sql.ErrNoRows {
		logger.Log.Error("matchStorelocation: " + err.Error())
		return false
	}
	if err == sql.ErrNoRows {
		return false
	}
	if strconv.Itoa(storelocation.EntityID) != entityId {
		return false
	}
	if m, err = env.DB.DoesPersonBelongsTo(pid, []models.Entity{storelocation.Entity}); err != nil {
		logger.Log.Error("matchStorelocation: " + err.Error())
		return false
	}
	logger.Log.WithFields(logrus.Fields{"m": m}).Debug("matchStorelocation")

	return m
}

func (env *Env) MatchStorelocationFunc(args ...interface{}) (interface{}, error) {
	personId := args[0].(string)
	itemId := args[1].(string)
	entityId := args[2].(string)

	return (bool)(env.matchStorelocation(personId, itemId, entityId)), nil
}

func (env *Env) matchStorage(personId string, itemId string, entityId string) bool {
	var (
		pid, iid int
		err      error
		m        bool
		ent      models.Entity
	)

	if pid, err = strconv.Atoi(personId); err != nil {
		logger.Log.Error("matchStorage: " + err.Error())
		return false
	}
	if iid, err = strconv.Atoi(itemId); err != nil {
		logger.Log.Error("matchStorage: " + err.Error())
		return false
	}

	if ent, err = env.DB.GetStorageEntity(iid); err != nil {
		logger.Log.Error(fmt.Sprintf("matchStorage: %v %s", ent, err.Error()))
		return false
	}
	if strconv.Itoa(ent.EntityID) != entityId {
		return false
	}
	if m, err = env.DB.DoesPersonBelongsTo(pid, []models.Entity{ent}); err != nil {
		logger.Log.Error(fmt.Sprintf("matchStorage: %v %s", ent, err.Error()))
		return false
	}
	logger.Log.WithFields(logrus.Fields{"m": m}).Debug("matchStorage")

	return m
}

func (env *Env) MatchStorageFunc(args ...interface{}) (interface{}, error) {
	personId := args[0].(string)
	itemId := args[1].(string)
	entityId := args[2].(string)

	return (bool)(env.matchStorage(personId, itemId, entityId)), nil
}

func (env *Env) matchEntity(personId string, entityId string) bool {
	var (
		pid, eid int
		err      error
		m        bool
	)

	if pid, err = strconv.Atoi(personId); err != nil {
		logger.Log.Error("matchEntity: " + err.Error())
		return false
	}
	if eid, err = strconv.Atoi(entityId); err != nil {
		logger.Log.Error("matchEntity: " + err.Error())
		return false
	}
	if m, err = env.DB.DoesPersonBelongsTo(pid, []models.Entity{{EntityID: eid}}); err != nil {
		logger.Log.Error("matchEntity: " + err.Error())
		return false
	}
	logger.Log.WithFields(logrus.Fields{"personId": personId, "entityId": entityId, "m": m}).Debug("matchEntity")

	return m
}

func (env *Env) MatchEntityFunc(args ...interface{}) (interface{}, error) {
	personId := args[0].(string)
	entityId := args[1].(string)

	return (bool)(env.matchEntity(personId, entityId)), nil
}

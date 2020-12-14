package handlers

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/globals"
	"github.com/tbellembois/gochimitheque/models"
)

func (env *Env) ReloadPolicy() {
	var err error
	if globals.JSONAdapterData, err = env.DB.ToCasbinJSONAdapter(); err != nil {
		globals.Log.Error("error getting json adapter data: " + err.Error())
		os.Exit(1)
	}
	if err = globals.Enforcer.LoadPolicy(); err != nil {
		globals.Log.Error("enforcer policy load error: " + err.Error())
		os.Exit(1)
	}
}

func (env *Env) matchPeople(personId string, itemId string, entityId string) bool {
	var (
		pid, iid int
		err      error
		ent      []models.Entity
	)

	if pid, err = strconv.Atoi(personId); err != nil {
		globals.Log.Fatal("matchPeople: " + err.Error())
		return false
	}
	if iid, err = strconv.Atoi(itemId); err != nil {
		globals.Log.Fatal("matchPeople: " + err.Error())
		return false
	}

	if ent, err = env.DB.GetPersonEntities(pid, iid); err != nil {
		globals.Log.Fatal("matchPeople: " + err.Error())
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
		pid, iid int
		err      error
		m        bool
		ent      models.Entity
	)
	globals.Log.WithFields(logrus.Fields{"personId": personId, "itemId": itemId, "entityId": entityId}).Debug("matchStorelocation")

	if pid, err = strconv.Atoi(personId); err != nil {
		globals.Log.Fatal("matchStorelocation: " + err.Error())
		return false
	}
	if iid, err = strconv.Atoi(itemId); err != nil {
		globals.Log.Fatal("matchStorelocation: " + err.Error())
		return false
	}
	if ent, err = env.DB.GetStoreLocationEntity(iid); err != nil {
		globals.Log.Fatal("matchStorelocation: " + err.Error())
		return false
	}
	if strconv.Itoa(ent.EntityID) != entityId {
		return false
	}
	if m, err = env.DB.DoesPersonBelongsTo(pid, []models.Entity{ent}); err != nil {
		globals.Log.Fatal("matchStorelocation: " + err.Error())
		return false
	}
	globals.Log.WithFields(logrus.Fields{"m": m}).Debug("matchStorelocation")

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
		globals.Log.Fatal("matchStorage: " + err.Error())
		return false
	}
	if iid, err = strconv.Atoi(itemId); err != nil {
		globals.Log.Fatal("matchStorage: " + err.Error())
		return false
	}

	if ent, err = env.DB.GetStorageEntity(iid); err != nil {
		globals.Log.Fatal(fmt.Sprintf("matchStorage: %v %s", ent, err.Error()))
		return false
	}
	if strconv.Itoa(ent.EntityID) != entityId {
		return false
	}
	if m, err = env.DB.DoesPersonBelongsTo(pid, []models.Entity{ent}); err != nil {
		globals.Log.Fatal(fmt.Sprintf("matchStorage: %v %s", ent, err.Error()))
		return false
	}
	globals.Log.WithFields(logrus.Fields{"m": m}).Debug("matchStorage")

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
		globals.Log.Fatal("matchEntity: " + err.Error())
		return false
	}
	if eid, err = strconv.Atoi(entityId); err != nil {
		globals.Log.Fatal("matchEntity: " + err.Error())
		return false
	}
	if m, err = env.DB.DoesPersonBelongsTo(pid, []models.Entity{{EntityID: eid}}); err != nil {
		globals.Log.Fatal("matchEntity: " + err.Error())
		return false
	}
	globals.Log.WithFields(logrus.Fields{"personId": personId, "entityId": entityId, "m": m}).Debug("matchEntity")

	return m
}

func (env *Env) MatchEntityFunc(args ...interface{}) (interface{}, error) {
	personId := args[0].(string)
	entityId := args[1].(string)

	return (bool)(env.matchEntity(personId, entityId)), nil
}

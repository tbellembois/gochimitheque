package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/constants"
)

// GetPeople returns the people matching the search criteria
// order, offset and limit are passed to the sql request
func (db *SQLiteDataStore) GetPeople(personID int, search string, order string, offset uint64, limit uint64) ([]Person, error) {
	var (
		people []Person
		sqlr   string
		sqla   []interface{}
	)
	log.WithFields(log.Fields{"search": search, "order": order, "offset": offset, "limit": limit}).Debug("GetEntities")

	sbuilder := sq.Select(`p.person_id, 
		p.person_email`).
		From("person AS p").
		Where("p.person_email LIKE ?", fmt.Sprint("%", search, "%")).
		Join(buildPermissionFilter("person", "p", "person_id", "r"), fmt.Sprint(personID)).
		GroupBy("p.person_id").
		OrderBy(fmt.Sprintf("person_email %s", order))
	if limit != constants.MaxUint64 {
		sbuilder = sbuilder.Offset(offset).Limit(limit)
	}
	sqlr, sqla, db.err = sbuilder.ToSql()
	log.Debug(sqlr)
	if db.err != nil {
		return nil, db.err
	}

	if db.err = db.Select(&people, sqlr, sqla...); db.err != nil {
		return nil, db.err
	}
	return people, nil
}

// GetPerson returns the person with id "id"
func (db *SQLiteDataStore) GetPerson(id int) (Person, error) {
	var (
		person Person
		sqlr   string
	)

	sqlr = "SELECT person_id, person_email FROM person WHERE person_id = ?"
	if db.err = db.Get(&person, sqlr, id); db.err != nil {
		return Person{}, db.err
	}
	return person, nil
}

// GetPersonByEmail returns the person with email "email"
func (db *SQLiteDataStore) GetPersonByEmail(email string) (Person, error) {
	var (
		person Person
		sqlr   string
	)

	sqlr = "SELECT person_id, person_email FROM person WHERE person_email = ?"
	if db.err = db.Get(&person, sqlr, email); db.err != nil {
		return Person{}, db.err
	}
	return person, nil
}

// GetPersonPermissions returns the person (with id "id") permissions
func (db *SQLiteDataStore) GetPersonPermissions(id int) ([]Permission, error) {
	var (
		ps   []Permission
		sqlr string
	)

	sqlr = `SELECT permission_id, permission_perm_name, permission_item_name, permission_entityid 
	FROM permission
	WHERE permission_person_id = ?`
	if db.err = db.Select(&ps, sqlr, id); db.err != nil {
		return nil, db.err
	}
	log.WithFields(log.Fields{"personID": id, "ps": ps}).Debug("GetPersonPermissions")
	return ps, nil
}

// GetPersonEntities returns the person (with id "id") entities
func (db *SQLiteDataStore) GetPersonEntities(id int) ([]Entity, error) {
	var (
		es   []Entity
		sqlr string
	)

	sqlr = `SELECT entity_id, entity_name, entity_description 
	FROM entity
	INNER JOIN personentities ON personentities.personentities_entity_id = entity.entity_id
	WHERE personentities.personentities_person_id = ?`
	if db.err = db.Select(&es, sqlr, id); db.err != nil {
		return nil, db.err
	}
	log.WithFields(log.Fields{"personID": id, "es": es}).Debug("GetPersonEntities")
	return es, nil
}

// DoesPersonBelongsTo returns true if the person (with id "id") belongs to the entities
func (db *SQLiteDataStore) DoesPersonBelongsTo(id int, entities []Entity) (bool, error) {
	var (
		sqlr  string
		count int
	)

	// extracting entities ids
	var eids []int
	for _, i := range entities {
		eids = append(eids, i.EntityID)
	}

	sqlr = `SELECT count(*) 
	FROM personentities
	WHERE personentities_person_id = ? 
	AND personentities_entity_id IN (?)`
	if db.err = db.Get(&count, sqlr, id, eids); db.err != nil {
		return false, db.err
	}
	log.WithFields(log.Fields{"personID": id, "count": count}).Debug("DoesPersonBelongsTo")
	return count > 0, nil
}

// HasPersonPermission returns true if the person with id "id" has the permission "perm" on the item "item" with id "itemid"
func (db *SQLiteDataStore) HasPersonPermission(id int, perm string, item string, itemid int) (bool, error) {
	// itemid == -1 means all items
	// itemid == -2 means any items (-2 is not a database permission_entityid possible value)
	var (
		res     bool
		count   int
		sqlr    string
		sqlargs []interface{}
		err     error
		eids    []int
	)

	log.WithFields(log.Fields{
		"id":     id,
		"perm":   perm,
		"item":   item,
		"itemid": itemid}).Debug("HasPersonPermission")

	//
	// first: retrieving the entities of the item to be accessed
	//
	switch item {
	case "person":
		// retrieving the requeted person entities
		var rpe []Entity
		if rpe, err = db.GetPersonEntities(itemid); err != nil {
			return false, err
		}
		// and their ids
		for _, i := range rpe {
			eids = append(eids, i.EntityID)
		}
	case "entity":
		eids = append(eids, itemid)
	}

	//
	// second: has the logged user "perm" on the "item" of the entities in "eids"
	//
	if itemid == -2 {
		// possible matchs:
		// permission_perm_name | permission_item_name
		// all | all
		// all | ?
		// ?   | all  => no sense (look at explanation in the else section)
		// ?   | ?
		sqlr = `SELECT count(*) FROM permission WHERE 
		permission_person_id = ? AND permission_perm_name = "all" AND permission_item_name = "all"  OR 
		permission_person_id = ? AND permission_perm_name = "all" AND permission_item_name = ? OR 
		permission_person_id = ? AND permission_perm_name = ? AND permission_item_name = "all"  OR
		permission_person_id = ? AND permission_perm_name = ? AND permission_item_name = ?`
		if db.err = db.Get(&count, sqlr, id, id, item, id, perm, id, perm, item); db.err != nil {
			switch {
			case db.err == sql.ErrNoRows:
				return false, nil
			default:
				return false, db.err
			}
		}
	} else {
		// possible matchs:
		// permission_perm_name | permission_item_name | permission_entityid
		// all | ?   | -1 (ex: all permissions on all entities)
		// all | ?   | ?  (ex: all permissions on entity 3)
		// ?   | all | -1 => no sense (ex: r permission on entities, store_locations...) we will put the permissions for each item
		// ?   | all | ?  => no sense (ex: r permission on entities, store_locations... with id = 3)
		// all | all | -1 => means super admin
		// all | all | ?  => no sense (ex: all permission on entities, store_locations... with id = 3)
		// ?   | ?   | -1 => (ex: r permission on all entities)
		// ?   | ?   | ?  => (ex: r permission on entity 3)
		if sqlr, sqlargs, db.err = sqlx.In(`SELECT count(*) FROM permission WHERE 
		permission_person_id = ? AND permission_item_name = "all" AND permission_perm_name = "all" OR 
		permission_person_id = ? AND permission_item_name = "all" AND permission_perm_name = ? AND permission_entityid = -1 OR
		permission_person_id = ? AND permission_item_name = ? AND permission_perm_name = "all" AND permission_entityid IN (?) OR
		permission_person_id = ? AND permission_item_name = ? AND permission_perm_name = "all" AND permission_entityid = -1 OR 
		permission_person_id = ? AND permission_item_name = ? AND permission_perm_name = ? AND permission_entityid = -1 OR
		permission_person_id = ? AND permission_item_name = ? AND permission_perm_name = ? AND permission_entityid IN (?)
		`, id, id, perm, id, item, eids, id, item, id, item, perm, id, item, perm, eids); db.err != nil {
			return false, db.err
		}

		if db.err = db.Get(&count, sqlr, sqlargs...); db.err != nil {
			switch {
			case db.err == sql.ErrNoRows:
				return false, nil
			default:
				return false, db.err
			}
		}
	}

	log.WithFields(log.Fields{"count": count}).Debug("HasPersonPermission")
	if count == 0 {
		res = false
	} else {
		res = true
	}
	return res, nil
}

// CreatePerson creates the given person
func (db *SQLiteDataStore) CreatePerson(p Person) error {
	var (
		sqlr   string
		res    sql.Result
		lastid int64
	)
	// inserting person
	sqlr = `INSERT INTO person(person_email, person_password) VALUES (?, ?)`
	if res, db.err = db.Exec(sqlr, p.PersonEmail, p.PersonPassword); db.err != nil {
		return db.err
	}

	// getting the last inserted id
	if lastid, db.err = res.LastInsertId(); db.err != nil {
		return db.err
	}
	p.PersonID = int(lastid)

	// inserting permissions
	for _, per := range p.Permissions {
		sqlr = `INSERT INTO permission(permission_person_id, permission_perm_name, permission_item_name, permission_entityid) VALUES (?, ?, ?, ?)`
		if _, db.err = db.Exec(sqlr, p.PersonID, per.PermissionPermName, per.PermissionItemName, -1); db.err != nil {
			return db.err
		}
	}
	return nil
}

// UpdatePerson updates the given person
func (db *SQLiteDataStore) UpdatePerson(p Person) error {
	var (
		sqlr string
	)
	sqlr = `UPDATE person SET person_email = ?
	WHERE person_id = ?`
	if _, db.err = db.Exec(sqlr, p.PersonEmail, p.PersonID); db.err != nil {
		return db.err
	}
	return nil
}

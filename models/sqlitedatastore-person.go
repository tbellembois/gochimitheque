package models

import (
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/constants"
	"github.com/tbellembois/gochimitheque/helpers"
)

// GetPeople returns the people matching the search criteria
// order, offset and limit are passed to the sql request
func (db *SQLiteDataStore) GetPeople(p helpers.DbselectparamPerson) ([]Person, int, error) {
	var (
		people                             []Person
		isadmin                            bool
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)
	log.WithFields(log.Fields{"p": p}).Debug("GetPeople")

	// is the logged user an admin?
	if isadmin, err = db.IsPersonAdmin(p.GetLoggedPersonID()); err != nil {
		return nil, 0, err
	}

	// returning all people for admins
	// we need to handle admins
	// to see people with no entities
	precreq.WriteString("SELECT count(DISTINCT p.person_id)")
	presreq.WriteString("SELECT p.person_id, p.person_email")
	comreq.WriteString(" FROM person AS p, entity AS e")
	comreq.WriteString(" JOIN personentities ON personentities.personentities_person_id = p.person_id")
	if p.GetEntity() != -1 {
		comreq.WriteString(" JOIN entity ON personentities.personentities_entity_id = :entity")
	} else {
		comreq.WriteString(" JOIN entity ON personentities.personentities_entity_id = e.entity_id")
	}
	if !isadmin {
		comreq.WriteString(` JOIN permission AS perm ON
		(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
		(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
		(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
		(perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
		(perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
		(perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
		(perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "r" and perm.permission_entity_id = e.entity_id)
		`)
	}
	comreq.WriteString(" WHERE p.person_email LIKE :search")
	postsreq.WriteString(" GROUP BY p.person_id")
	postsreq.WriteString(" ORDER BY " + p.GetOrderBy() + " " + p.GetOrder())

	// limit
	if p.GetLimit() != constants.MaxUint64 {
		postsreq.WriteString(" LIMIT :limit OFFSET :offset")
	}
	log.Debug(presreq.String() + comreq.String() + postsreq.String())
	// building count and select statements
	if cnstmt, err = db.PrepareNamed(precreq.String() + comreq.String()); err != nil {
		return nil, 0, err
	}
	if snstmt, err = db.PrepareNamed(presreq.String() + comreq.String() + postsreq.String()); err != nil {
		return nil, 0, err
	}

	// building argument map
	m := map[string]interface{}{
		"entity":   p.GetEntity(),
		"search":   p.GetSearch(),
		"personid": p.GetLoggedPersonID(),
		"order":    p.GetOrder(),
		"limit":    p.GetLimit(),
		"offset":   p.GetOffset(),
	}

	// select
	if err = snstmt.Select(&people, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	log.WithFields(log.Fields{"people": people, "count": count}).Debug("GetPeople")
	return people, count, nil
}

// GetPerson returns the person with id "id"
func (db *SQLiteDataStore) GetPerson(id int) (Person, error) {
	var (
		person Person
		sqlr   string
		err    error
	)

	sqlr = "SELECT person_id, person_email FROM person WHERE person_id = ?"
	if err = db.Get(&person, sqlr, id); err != nil {
		return Person{}, err
	}
	return person, nil
}

// GetPersonByEmail returns the person with email "email"
func (db *SQLiteDataStore) GetPersonByEmail(email string) (Person, error) {
	var (
		person Person
		sqlr   string
		err    error
	)

	sqlr = "SELECT person_id, person_email FROM person WHERE person_email = ?"
	if err = db.Get(&person, sqlr, email); err != nil {
		return Person{}, err
	}
	return person, nil
}

// GetPersonPermissions returns the person (with id "id") permissions
func (db *SQLiteDataStore) GetPersonPermissions(id int) ([]Permission, error) {
	var (
		ps   []Permission
		sqlr string
		err  error
	)

	sqlr = `SELECT permission_id, permission_perm_name, permission_item_name, permission_entity_id 
	FROM permission
	WHERE person = ?`
	if err = db.Select(&ps, sqlr, id); err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{"personID": id, "ps": ps}).Debug("GetPersonPermissions")
	return ps, nil
}

// GetPersonManageEntities returns the entities the person (with id "id") if manager of
func (db *SQLiteDataStore) GetPersonManageEntities(id int) ([]Entity, error) {
	var (
		es   []Entity
		sqlr string
		err  error
	)

	sqlr = `SELECT entity_id, entity_name, entity_description 
	FROM entity
	LEFT JOIN entitypeople ON entitypeople.entitypeople_entity_id = entity.entity_id
	WHERE entitypeople.entitypeople_person_id = ?`
	if err = db.Select(&es, sqlr, id); err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{"personID": id, "es": es}).Debug("GetPersonManageEntities")
	return es, nil
}

// GetPersonEntities returns the person (with id "id") entities
func (db *SQLiteDataStore) GetPersonEntities(LoggedPersonID int, id int) ([]Entity, error) {
	var (
		entities []Entity
		isadmin  bool
		sqlr     strings.Builder
		sstmt    *sqlx.NamedStmt
		err      error
	)

	// is the logged user an admin?
	if isadmin, err = db.IsPersonAdmin(LoggedPersonID); err != nil {
		return nil, err
	}

	sqlr.WriteString("SELECT e.entity_id, e.entity_name, e.entity_description")
	sqlr.WriteString(" FROM entity AS e, person AS p, personentities as pe")
	if !isadmin {
		sqlr.WriteString(` JOIN permission AS perm ON
		(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
		(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
		(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
		(perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
		(perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
		(perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
		(perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "r" and perm.permission_entity_id = e.entity_id)
		`)
	}
	sqlr.WriteString(" WHERE pe.personentities_person_id = :personid AND e.entity_id == pe.personentities_entity_id")
	sqlr.WriteString(" GROUP BY e.entity_id")
	sqlr.WriteString(" ORDER BY e.entity_name ASC")

	// building select statement
	if sstmt, err = db.PrepareNamed(sqlr.String()); err != nil {
		return nil, err
	}

	// building argument map
	m := map[string]interface{}{
		"personid": id}

	if err = sstmt.Select(&entities, m); err != nil {
		return nil, err
	}
	return entities, nil
}

// DoesPersonBelongsTo returns true if the person (with id "id") belongs to the entities
func (db *SQLiteDataStore) DoesPersonBelongsTo(id int, entities []Entity) (bool, error) {
	var (
		sqlr  string
		count int
		err   error
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
	if err = db.Get(&count, sqlr, id, eids); err != nil {
		return false, err
	}
	log.WithFields(log.Fields{"personID": id, "count": count}).Debug("DoesPersonBelongsTo")
	return count > 0, nil
}

// HasPersonPermission returns true if the person with id "id" has the permission "perm" on the item "item" id "itemid"
func (db *SQLiteDataStore) HasPersonPermission(id int, perm string, item string, itemid int) (bool, error) {
	// itemid == -1 means all itemid
	// itemid == -2 means any itemid
	var (
		res     bool
		count   int
		sqlr    string
		sqlargs []interface{}
		err     error
		eids    []int
		isadmin bool
	)

	log.WithFields(log.Fields{
		"id":     id,
		"perm":   perm,
		"item":   item,
		"itemid": itemid}).Debug("HasPersonPermission")

	// is the user an admin?
	if isadmin, err = db.IsPersonAdmin(id); err != nil {
		return false, err
	}
	// if yes return true
	if isadmin {
		return true, nil
	}

	//
	// first: retrieving the entities of the item to be accessed
	//
	if itemid != -1 && itemid != -2 {
		switch item {
		case "people":
			// retrieving the requested person entities
			var rpe []Entity
			if rpe, err = db.GetPersonEntities(id, itemid); err != nil {
				return false, err
			}
			// and their ids
			for _, i := range rpe {
				eids = append(eids, i.EntityID)
			}
		case "entities":
			eids = append(eids, itemid)
		case "storages":
			eids = append(eids, itemid)
		case "products":
			eids = append(eids, itemid)
		case "storelocations":
			// retrieving the requested store location entity
			var rpe Entity
			if rpe, err = db.GetStoreLocationEntity(itemid); err != nil {
				return false, err
			}
			eids = append(eids, rpe.EntityID)
		}
		log.WithFields(log.Fields{"eids": eids}).Debug("HasPersonPermission")
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
		(person = ? AND permission_perm_name = "all" AND permission_item_name = "all")  OR
		(person = ? AND permission_perm_name = "all" AND permission_item_name = ?) OR
		(person = ? AND permission_perm_name = ? AND permission_item_name = "all")  OR
		(person = ? AND permission_perm_name = ? AND permission_item_name = ?)`
		if err = db.Get(&count, sqlr, id, id, item, id, perm, id, perm, item); err != nil {
			switch {
			case err == sql.ErrNoRows:
				return false, nil
			default:
				return false, err
			}
		}
	} else if itemid == -1 {
		// possible matchs:
		// permission_perm_name | permission_item_name | permission_entity_id
		// all | ?   | -1 (ex: all permissions on all entities)
		// ?   | all | -1 => no sense (ex: r permission on entities, store_locations...) we will put the permissions for each item
		// all | all | -1 => means super admin
		// ?   | ?   | -1 => (ex: r permission on all entities)
		if sqlr, sqlargs, err = sqlx.In(`SELECT count(*) FROM permission WHERE 
		person = ? AND permission_item_name = "all" AND permission_perm_name = "all" OR 
		person = ? AND permission_item_name = "all" AND permission_perm_name = ? AND permission_entity_id = -1 OR
		person = ? AND permission_item_name = ? AND permission_perm_name = "all" AND permission_entity_id = -1 OR 
		person = ? AND permission_item_name = ? AND permission_perm_name = ? AND permission_entity_id = -1
		`, id, id, perm, id, item, id, item, perm); err != nil {
			return false, err
		}

		if err = db.Get(&count, sqlr, sqlargs...); err != nil {
			switch {
			case err == sql.ErrNoRows:
				return false, nil
			default:
				return false, err
			}
		}
	} else {
		// possible matchs:
		// permission_perm_name | permission_item_name | permission_entity_id
		// all | ?   | -1 (ex: all permissions on all entities)
		// all | ?   | ?  (ex: all permissions on entity 3)
		// ?   | all | -1 => no sense (ex: r permission on entities, store_locations...) we will put the permissions for each item
		// ?   | all | ?  => no sense (ex: r permission on entities, store_locations... with id = 3)
		// all | all | -1 => means super admin
		// all | all | ?  => no sense (ex: all permission on entities, store_locations... with id = 3)
		// ?   | ?   | -1 => (ex: r permission on all entities)
		// ?   | ?   | ?  => (ex: r permission on entity 3)
		if sqlr, sqlargs, err = sqlx.In(`SELECT count(*) FROM permission WHERE 
		person = ? AND permission_item_name = "all" AND permission_perm_name = "all" OR 
		person = ? AND permission_item_name = "all" AND permission_perm_name = ? AND permission_entity_id = -1 OR
		person = ? AND permission_item_name = "all" AND permission_perm_name = ? AND permission_entity_id IN (?) OR
		person = ? AND permission_item_name = ? AND permission_perm_name = "all" AND permission_entity_id IN (?) OR
		person = ? AND permission_item_name = ? AND permission_perm_name = "all" AND permission_entity_id = -1 OR 
		person = ? AND permission_item_name = ? AND permission_perm_name = ? AND permission_entity_id = -1 OR
		person = ? AND permission_item_name = ? AND permission_perm_name = ? AND permission_entity_id IN (?)
		`, id, id, perm, id, perm, eids, id, item, eids, id, item, id, item, perm, id, item, perm, eids); err != nil {
			return false, err
		}

		if err = db.Get(&count, sqlr, sqlargs...); err != nil {
			switch {
			case err == sql.ErrNoRows:
				return false, nil
			default:
				return false, err
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

// DeletePerson deletes the person with id "id"
func (db *SQLiteDataStore) DeletePerson(id int) error {
	var (
		sqlr string
		err  error
	)
	sqlr = `DELETE FROM personentities 
	WHERE personentities_person_id = ?`
	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	sqlr = `DELETE FROM entitypeople 
	WHERE entitypeople_person_id = ?`
	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	sqlr = `DELETE FROM permission 
	WHERE person = ?`
	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	sqlr = `DELETE FROM person 
	WHERE person_id = ?`
	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}
	return nil
}

// CreatePerson creates the given person
func (db *SQLiteDataStore) CreatePerson(p Person) (error, int) {
	var (
		sqlr   string
		res    sql.Result
		lastid int64
		err    error
	)

	// inserting person
	// FIXME: use a transaction here
	sqlr = `INSERT INTO person(person_email, person_password) VALUES (?, ?)`
	if res, err = db.Exec(sqlr, p.PersonEmail, p.PersonPassword); err != nil {
		return err, 0
	}

	// getting the last inserted id
	if lastid, err = res.LastInsertId(); err != nil {
		return err, 0
	}
	p.PersonID = int(lastid)

	// inserting permissions
	for _, per := range p.Permissions {
		sqlr = `INSERT INTO permission(person, permission_perm_name, permission_item_name, permission_entity_id) VALUES (?, ?, ?, ?)`
		if _, err = db.Exec(sqlr, p.PersonID, per.PermissionPermName, per.PermissionItemName, per.PermissionEntityID); err != nil {
			return err, 0
		}
		// adding r permission for w permissions
		if per.PermissionPermName == "w" {
			if _, err = db.Exec(sqlr, p.PersonID, "r", per.PermissionItemName, per.PermissionEntityID); err != nil {
				return err, 0
			}
		}
	}

	// inserting entities
	for _, e := range p.Entities {
		sqlr = `INSERT INTO personentities(personentities_person_id, personentities_entity_id) 
			VALUES (?, ?)`
		if _, err = db.Exec(sqlr, p.PersonID, e.EntityID); err != nil {
			return err, 0
		}
	}
	return nil, p.PersonID
}

// UpdatePerson updates the given person
func (db *SQLiteDataStore) UpdatePerson(p Person) error {
	var (
		sqlr string
		err  error
	)
	// updating person
	// FIXME: use a transaction here
	sqlr = `UPDATE person SET person_email = ?
	WHERE person_id = ?`
	if _, err = db.Exec(sqlr, p.PersonEmail, p.PersonID); err != nil {
		return err
	}

	// lazily deleting former entities
	sqlr = `DELETE FROM personentities 
	WHERE personentities_person_id = ?`
	if _, err = db.Exec(sqlr, p.PersonID); err != nil {
		return err
	}

	// updating person entities
	for _, e := range p.Entities {
		sqlr = `INSERT INTO personentities(personentities_person_id, personentities_entity_id) 
		VALUES (?, ?)`
		if _, err = db.Exec(sqlr, p.PersonID, e.EntityID); err != nil {
			return err
		}
	}

	// lazily deleting former permissions
	sqlr = `DELETE FROM permission 
		WHERE person = ?`
	if _, err = db.Exec(sqlr, p.PersonID); err != nil {
		return err
	}

	// updating person permissions
	for _, perm := range p.Permissions {
		sqlr = `INSERT INTO permission(person, permission_perm_name, permission_item_name, permission_entity_id) 
		VALUES (?, ?, ?, ?)`
		if _, err = db.Exec(sqlr, p.PersonID, perm.PermissionPermName, perm.PermissionItemName, perm.PermissionEntityID); err != nil {
			return err
		}
		// adding r permission for w permissions
		if perm.PermissionPermName == "w" {
			if _, err = db.Exec(sqlr, p.PersonID, "r", perm.PermissionItemName, perm.PermissionEntityID); err != nil {
				return err
			}
		}
	}

	return nil
}

// IsPersonAdmin returns true is the person with id "id" is an admin
func (db *SQLiteDataStore) IsPersonAdmin(id int) (bool, error) {
	var (
		res   bool
		count int
		sqlr  string
		err   error
	)
	sqlr = `SELECT count(*) from permission WHERE 
	permission.person = ? AND
	permission.permission_perm_name = "all" AND
	permission.permission_item_name = "all" AND
	permission_entity_id = -1`
	if err = db.Get(&count, sqlr, id); err != nil {
		return false, err
	}
	log.WithFields(log.Fields{"id": id, "count": count}).Debug("IsPersonAdmin")
	if count == 0 {
		res = false
	} else {
		res = true
	}
	return res, nil
}

// IsPersonWithEmail returns true is the person with id "id" is a manager
func (db *SQLiteDataStore) IsPersonManager(id int) (bool, error) {
	var (
		res   bool
		count int
		sqlr  string
		err   error
	)
	sqlr = "SELECT count(*) from entitypeople WHERE entitypeople.entitypeople_person_id = ?"
	if err = db.Get(&count, sqlr, id); err != nil {
		return false, err
	}
	log.WithFields(log.Fields{"id": id, "count": count}).Debug("IsPersonManager")
	if count == 0 {
		res = false
	} else {
		res = true
	}
	return res, nil
}

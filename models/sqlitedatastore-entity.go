package models

import (
	"fmt"

	"database/sql"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/constants"
)

// GetEntities returns the entities matching the search criteria
// order, offset and limit are passed to the sql request
func (db *SQLiteDataStore) GetEntities(personID int, search string, order string, offset uint64, limit uint64) ([]Entity, int, error) {
	var (
		entities []Entity
		count    int
		sqlr     string
		sqla     []interface{}
	)
	log.WithFields(log.Fields{"search": search, "order": order, "offset": offset, "limit": limit}).Debug("GetEntities")

	// count query
	cbuilder := sq.Select(`count(DISTINCT e.entity_id)`).
		From("entity AS e, person as p").
		Where(`e.entity_name LIKE ?`, fmt.Sprint("%", search, "%")).
		// join to filter entities personID can access to
		Join(`permission AS perm on
			(perm.permission_person_id = ? and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
			(perm.permission_person_id = ? and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
			(perm.permission_person_id = ? and perm.permission_item_name = "all" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
			(perm.permission_person_id = ? and perm.permission_item_name = "entities" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
			(perm.permission_person_id = ? and perm.permission_item_name = "entities" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
			(perm.permission_person_id = ? and perm.permission_item_name = "entities" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
			(perm.permission_person_id = ? and perm.permission_item_name = "entities" and perm.permission_perm_name = "r" and perm.permission_entity_id = e.entity_id)
			`, personID, personID, personID, personID, personID, personID, personID)
	// select query
	sbuilder := sq.Select(`e.entity_id, 
		e.entity_id,
		e.entity_name, 
		e.entity_description`).
		From("entity AS e, person as p").
		Where(`e.entity_name LIKE ?`, fmt.Sprint("%", search, "%")).
		// join to filter entities personID can access to
		Join(`permission AS perm on
			(perm.permission_person_id = ? and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
			(perm.permission_person_id = ? and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
			(perm.permission_person_id = ? and perm.permission_item_name = "all" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
			(perm.permission_person_id = ? and perm.permission_item_name = "entities" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
			(perm.permission_person_id = ? and perm.permission_item_name = "entities" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
			(perm.permission_person_id = ? and perm.permission_item_name = "entities" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
			(perm.permission_person_id = ? and perm.permission_item_name = "entities" and perm.permission_perm_name = "r" and perm.permission_entity_id = e.entity_id)
			`, personID, personID, personID, personID, personID, personID, personID).
		// join to get managers
		// Join(`entitypeople ON entitypeople.entitypeople_entity_id = e.entity_id`).
		// Join(`person ON entitypeople.entitypeople_person_id = p.person_id`).
		GroupBy("e.entity_id").
		OrderBy(fmt.Sprintf("entity_name %s", order))
	if limit != constants.MaxUint64 {
		sbuilder = sbuilder.Offset(offset).Limit(limit)
	}
	// select
	sqlr, sqla, db.err = sbuilder.ToSql()
	if db.err != nil {
		return nil, 0, db.err
	}
	if db.err = db.Select(&entities, sqlr, sqla...); db.err != nil {
		return nil, 0, db.err
	}
	// count
	sqlr, sqla, db.err = cbuilder.ToSql()
	if db.err != nil {
		return nil, 0, db.err
	}
	if db.err = db.Get(&count, sqlr, sqla...); db.err != nil {
		return nil, 0, db.err
	}
	return entities, count, nil
}

// GetEntity returns the entity with id "id"
func (db *SQLiteDataStore) GetEntity(id int) (Entity, error) {
	var (
		entity Entity
		sqlr   string
	)

	sqlr = `SELECT e.entity_id, e.entity_name, e.entity_description
	FROM entity AS e
	WHERE e.entity_id = ?`
	if db.err = db.Get(&entity, sqlr, id); db.err != nil {
		return Entity{}, db.err
	}
	log.WithFields(log.Fields{"ID": id, "entity": entity}).Debug("GetEntity")
	return entity, nil
}

// GetEntityPeople returns the entity (with id "id") managers
func (db *SQLiteDataStore) GetEntityPeople(id int) ([]Person, error) {
	var (
		people []Person
		sqlr   string
	)

	sqlr = `SELECT p.person_id, p.person_email
	FROM person AS p, entitypeople
	WHERE entitypeople.entitypeople_person_id == p.person_id AND entitypeople.entitypeople_entity_id = ?`
	if db.err = db.Select(&people, sqlr, id); db.err != nil {
		return []Person{}, db.err
	}
	log.WithFields(log.Fields{"ID": id, "people": people}).Debug("GetEntityPeople")
	return people, nil
}

// DeleteEntity deletes the entity with id "id"
func (db *SQLiteDataStore) DeleteEntity(id int) error {
	var (
		sqlr string
	)
	sqlr = `DELETE FROM entity 
	WHERE entity_id = ?`
	if _, db.err = db.Exec(sqlr, id); db.err != nil {
		return db.err
	}
	return nil
}

// CreateEntity creates the given entity
func (db *SQLiteDataStore) CreateEntity(e Entity) (error, int) {
	var (
		sqlr   string
		res    sql.Result
		lastid int64
	)
	sqlr = `INSERT INTO entity(entity_name, entity_description) VALUES (?, ?)`
	if res, db.err = db.Exec(sqlr, e.EntityName, e.EntityDescription); db.err != nil {
		return db.err, 0
	}

	// getting the last inserted id
	if lastid, db.err = res.LastInsertId(); db.err != nil {
		return db.err, 0
	}
	e.EntityID = int(lastid)

	// adding the new managers
	for _, m := range e.Managers {
		sqlr = `insert into entitypeople (entitypeople_entity_id, entitypeople_person_id) values (?, ?)`
		if _, db.err = db.Exec(sqlr, e.EntityID, m.PersonID); db.err != nil {
			return db.err, 0
		}
	}

	return nil, e.EntityID
}

// UpdateEntity updates the given entity
func (db *SQLiteDataStore) UpdateEntity(e Entity) error {
	var (
		sqlr string
		sqla []interface{}
	)
	log.WithFields(log.Fields{"e": e}).Debug("UpdateEntity")

	// updating the entity
	sqlr = `UPDATE entity SET entity_name = ?, entity_description = ?
	WHERE entity_id = ?`
	if _, db.err = db.Exec(sqlr, e.EntityName, e.EntityDescription, e.EntityID); db.err != nil {
		return db.err
	}

	// removing former managers
	notin := sq.Or{}
	// ex: AND (entitypeople_person_id <> ? OR entitypeople_person_id <> ?)
	for _, m := range e.Managers {
		notin = append(notin, sq.NotEq{"entitypeople_person_id": m.PersonID})
	}
	// ex: DELETE FROM entitypeople WHERE (entitypeople_entity_id = ? AND (entitypeople_person_id <> ? OR entitypeople_person_id <> ?)
	sbuilder := sq.Delete(`entitypeople`).Where(
		sq.And{
			sq.Eq{`entitypeople_entity_id`: e.EntityID},
			notin})
	sqlr, sqla, db.err = sbuilder.ToSql()
	if db.err != nil {
		return db.err
	}
	if _, db.err = db.Exec(sqlr, sqla...); db.err != nil {
		return db.err
	}

	// adding the new ones
	for _, m := range e.Managers {
		// adding the manager
		sqlr = `INSERT OR IGNORE INTO entitypeople (entitypeople_entity_id, entitypeople_person_id) VALUES (?, ?)`
		if _, db.err = db.Exec(sqlr, e.EntityID, m.PersonID); db.err != nil {
			return db.err
		}

		for _, man := range e.Managers {
			// setting the manager in the entity
			sqlr = `INSERT OR IGNORE INTO personentities(personentities_person_id, personentities_entity_id) 
			VALUES (?, ?)`
			if _, db.err = db.Exec(sqlr, man.PersonID, e.EntityID); db.err != nil {
				return db.err
			}

			// setting the manager permissions in the entity
			// 1. lazily deleting former permissions
			sqlr = `DELETE FROM permission 
			WHERE permission_person_id = ?`
			if _, db.err = db.Exec(sqlr, man.PersonID); db.err != nil {
				return db.err
			}
			// 2. inserting manager permissions
			sqlr = `INSERT INTO permission(permission_person_id, permission_perm_name, permission_item_name, permission_entity_id) 
			VALUES (?, ?, ?, ?)`
			if _, db.err = db.Exec(sqlr, man.PersonID, "all", "all", e.EntityID); db.err != nil {
				return db.err
			}

		}
	}

	return nil
}

// IsEntityWithName returns true is the entity "name" exists
func (db *SQLiteDataStore) IsEntityWithName(name string) (bool, error) {
	var (
		res   bool
		count int
		sqlr  string
	)

	sqlr = "SELECT count(*) from entity WHERE entity.entity_name = ?"
	if db.err = db.Get(&count, sqlr, name); db.err != nil {
		return false, db.err
	}
	log.WithFields(log.Fields{"name": name, "count": count}).Debug("HasEntityWithName")
	if count == 0 {
		res = false
	} else {
		res = true
	}
	return res, nil
}

// IsEntityWithNameExcept returns true is the entity "name" exists ignoring the "except" names
func (db *SQLiteDataStore) IsEntityWithNameExcept(name string, except ...string) (bool, error) {
	var (
		res   bool
		count int
		sqlr  sq.SelectBuilder
		w     sq.And
	)

	w = append(w, sq.Eq{"entity.entity_name": name})
	for _, e := range except {
		w = append(w, sq.NotEq{"entity.entity_name": e})
	}

	sqlr = sq.Select("count(*)").From("entity").Where(w)
	sql, args, _ := sqlr.ToSql()
	log.WithFields(log.Fields{"sql": sql, "args": args}).Debug("HasEntityWithNameExcept")

	if db.err = db.Get(&count, sql, args...); db.err != nil {
		return false, db.err
	}
	log.WithFields(log.Fields{"name": name, "count": count}).Debug("HasEntityWithNameExcept")
	if count == 0 {
		res = false
	} else {
		res = true
	}
	return res, nil
}

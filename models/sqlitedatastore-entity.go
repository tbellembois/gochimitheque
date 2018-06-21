package models

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/constants"
)

// GetEntities returns the entities matching the search criteria
// order, offset and limit are passed to the sql request
func (db *SQLiteDataStore) GetEntities(personID int, search string, order string, offset uint64, limit uint64) ([]Entity, error) {
	var (
		entities []Entity
		sqlr     string
		sqla     []interface{}
	)
	log.WithFields(log.Fields{"search": search, "order": order, "offset": offset, "limit": limit}).Debug("GetEntities")

	sbuilder := sq.Select(`e.entity_id, 
		e.entity_id,
		e.entity_name, 
		e.entity_description, 
		p.person_id, 
		p.person_email, 
		p.person_password`).
		From("entity AS e, person AS p").
		Where("e.entity_person_id = p.person_id AND e.entity_name LIKE ?", fmt.Sprint("%", search, "%")).
		Join(buildJoinFilterForItem("entity", "e", "entity_id", "r"), fmt.Sprint(personID)).
		GroupBy("e.entity_id").
		OrderBy(fmt.Sprintf("entity_name %s", order))
	if limit != constants.MaxUint64 {
		sbuilder = sbuilder.Offset(offset).Limit(limit)
	}
	sqlr, sqla, db.err = sbuilder.ToSql()

	if db.err != nil {
		return nil, db.err
	}

	if db.err = db.Select(&entities, sqlr, sqla...); db.err != nil {
		return nil, db.err
	}
	return entities, nil
}

// GetEntity returns the entity with id "id"
func (db *SQLiteDataStore) GetEntity(id int) (Entity, error) {
	var (
		entity Entity
		sqlr   string
	)

	sqlr = "SELECT e.entity_id, e.entity_name, e.entity_description, p.person_id, p.person_email, p.person_password FROM entity AS e, person AS p WHERE e.entity_person_id = p.person_id AND e.entity_id = ?"
	if db.err = db.Get(&entity, sqlr, id); db.err != nil {
		return Entity{}, db.err
	}
	log.WithFields(log.Fields{"ID": id, "entity": entity}).Debug("GetEntity")
	return entity, nil
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
func (db *SQLiteDataStore) CreateEntity(e Entity) error {
	var (
		sqlr string
	)
	sqlr = `INSERT INTO entity(entity_name, entity_description, entity_person_id) VALUES (?, ?, ?)`
	if _, db.err = db.Exec(sqlr, e.EntityName, e.EntityDescription, e.Person.PersonID); db.err != nil {
		return db.err
	}
	return nil
}

// UpdateEntity updates the given entity
func (db *SQLiteDataStore) UpdateEntity(e Entity) error {
	var (
		sqlr string
	)
	sqlr = `UPDATE entity SET entity_name = ?, entity_description = ?, entity_person_id = ?
	WHERE entity_id = ?`
	if _, db.err = db.Exec(sqlr, e.EntityName, e.EntityDescription, e.PersonID, e.EntityID); db.err != nil {
		return db.err
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

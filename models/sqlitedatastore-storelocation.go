package models

import (
	"fmt"

	"database/sql"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/constants"
)

// GetStoreLocations returns the store locations matching the search criteria
// order, offset and limit are passed to the sql request
func (db *SQLiteDataStore) GetStoreLocations(loggedpersonID int, search string, order string, offset uint64, limit uint64) ([]StoreLocation, int, error) {
	var (
		storelocations []StoreLocation
		count          int
		sqlr           string
		sqla           []interface{}
	)
	log.WithFields(log.Fields{"search": search, "order": order, "offset": offset, "limit": limit}).Debug("GetStoreLocations")

	// count query
	cbuilder := sq.Select(`count(DISTINCT s.storelocation_id)`).
		From("storelocation AS s").
		Where(`s.storelocation_name LIKE ?`, fmt.Sprint("%", search, "%")).
		// join to filter store locations personID can access to
		Join(`permission AS perm on
			(perm.permission_person_id = ? and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = s.storelocation_entity_id) OR
			(perm.permission_person_id = ? and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
			(perm.permission_person_id = ? and perm.permission_item_name = "all" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
			(perm.permission_person_id = ? and perm.permission_item_name = "entities" and perm.permission_perm_name = "all" and perm.permission_entity_id = s.storelocation_entity_id) OR
			(perm.permission_person_id = ? and perm.permission_item_name = "entities" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
			(perm.permission_person_id = ? and perm.permission_item_name = "entities" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
			(perm.permission_person_id = ? and perm.permission_item_name = "entities" and perm.permission_perm_name = "r" and perm.permission_entity_id = s.storelocation_entity_id)
			`, loggedpersonID, loggedpersonID, loggedpersonID, loggedpersonID, loggedpersonID, loggedpersonID, loggedpersonID)
	// select query
	sbuilder := sq.Select(`s.storelocation_id, 
		s.storelocation_name, 
		entity.entity_id AS "storelocation_entity_id.entity_id",
		entity.entity_name AS "storelocation_entity_id.entity_name"`).
		From("storelocation AS s").
		Join("entity ON s.storelocation_entity_id = entity.entity_id").
		// join to filter entities personID can access to
		Join(`permission AS perm on
		(perm.permission_person_id = ? and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = s.storelocation_entity_id) OR
		(perm.permission_person_id = ? and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
		(perm.permission_person_id = ? and perm.permission_item_name = "all" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
		(perm.permission_person_id = ? and perm.permission_item_name = "entities" and perm.permission_perm_name = "all" and perm.permission_entity_id = s.storelocation_entity_id) OR
		(perm.permission_person_id = ? and perm.permission_item_name = "entities" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
		(perm.permission_person_id = ? and perm.permission_item_name = "entities" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
		(perm.permission_person_id = ? and perm.permission_item_name = "entities" and perm.permission_perm_name = "r" and perm.permission_entity_id = s.storelocation_entity_id)
		`, loggedpersonID, loggedpersonID, loggedpersonID, loggedpersonID, loggedpersonID, loggedpersonID, loggedpersonID).
		Where(`s.storelocation_name LIKE ?`, fmt.Sprint("%", search, "%")).
		GroupBy("s.storelocation_id").
		OrderBy(fmt.Sprintf("storelocation_name %s", order))
	if limit != constants.MaxUint64 {
		sbuilder = sbuilder.Offset(offset).Limit(limit)
	}
	// select
	sqlr, sqla, db.err = sbuilder.ToSql()
	if db.err != nil {
		return nil, 0, db.err
	}
	if db.err = db.Select(&storelocations, sqlr, sqla...); db.err != nil {
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
	return storelocations, count, nil
}

// GetStoreLocation returns the store location with id "id"
func (db *SQLiteDataStore) GetStoreLocation(id int) (StoreLocation, error) {
	var (
		storelocation StoreLocation
		sqlr          string
	)

	sqlr = `SELECT s.storelocation_id, 
	s.storelocation_name, 
	entity.entity_id AS "storelocation_entity_id.entity_id",
	entity.entity_name AS "storelocation_entity_id.entity_name"
	FROM storelocation AS s
	JOIN entity ON s.storelocation_entity_id = entity.entity_id
	WHERE s.storelocation_id = ?`
	if db.err = db.Get(&storelocation, sqlr, id); db.err != nil {
		return StoreLocation{}, db.err
	}
	log.WithFields(log.Fields{"ID": id, "storelocation": storelocation}).Debug("GetStoreLocation")
	return storelocation, nil
}

// DeleteStoreLocation deletes the store location with id "id"
func (db *SQLiteDataStore) DeleteStoreLocation(id int) error {
	var (
		sqlr string
	)
	sqlr = `DELETE FROM storelocation 
	WHERE storelocation_id = ?`
	if _, db.err = db.Exec(sqlr, id); db.err != nil {
		return db.err
	}
	return nil
}

// CreateStoreLocation creates the given store location
func (db *SQLiteDataStore) CreateStoreLocation(s StoreLocation) (error, int) {
	var (
		sqlr   string
		res    sql.Result
		lastid int64
	)
	sqlr = `INSERT INTO storelocation(storelocation_name, storelocation_entity_id) VALUES (?, ?)`
	if res, db.err = db.Exec(sqlr, s.StoreLocationName, s.Entity.EntityID); db.err != nil {
		return db.err, 0
	}

	// getting the last inserted id
	if lastid, db.err = res.LastInsertId(); db.err != nil {
		return db.err, 0
	}
	s.StoreLocationID = int(lastid)

	return nil, s.StoreLocationID
}

// UpdateStoreLocation updates the given store location
func (db *SQLiteDataStore) UpdateStoreLocation(s StoreLocation) error {
	var (
		sqlr string
	)
	log.WithFields(log.Fields{"s": s}).Debug("UpdateStoreLocation")

	// updating the store location
	sqlr = `UPDATE storelocation SET storelocation_name = ?, storelocation_entity_id = ?
	WHERE storelocation_id = ?`
	if _, db.err = db.Exec(sqlr, s.StoreLocationName, s.Entity.EntityID, s.StoreLocationID); db.err != nil {
		return db.err
	}

	return nil
}

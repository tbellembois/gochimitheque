package models

import (
	"fmt"
	"strings"

	"database/sql"
	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/constants"
)

// GetStoreLocations returns the store locations matching the search criteria
// order, offset and limit are passed to the sql request
func (db *SQLiteDataStore) GetStoreLocations(p GetStoreLocationsParameters) ([]StoreLocation, int, error) {
	var (
		storelocations                     []StoreLocation
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
	)
	log.WithFields(log.Fields{"search": p.Search, "order": p.Order, "offset": p.Offset, "limit": p.Limit}).Debug("GetStoreLocations")

	precreq.WriteString(" SELECT count(DISTINCT s.storelocation_id)")
	presreq.WriteString(` SELECT s.storelocation_id, s.storelocation_name, 
	entity.entity_id AS "entity.entity_id", 
	entity.entity_name AS "entity.entity_name"`)
	comreq.WriteString(" FROM storelocation AS s, entity as e")

	if p.EntityID != -1 {
		comreq.WriteString(" JOIN entity ON s.entity = :entityid")
	} else {
		comreq.WriteString(" JOIN entity ON s.entity = entity.entity_id")
	}
	comreq.WriteString(` JOIN permission AS perm ON
	(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
	(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
	(perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "r" and perm.permission_entity_id = e.entity_id)
	`)
	comreq.WriteString(" WHERE s.storelocation_name LIKE :search")
	postsreq.WriteString(" GROUP BY s.storelocation_id")
	postsreq.WriteString(" ORDER BY s.storelocation_name " + p.Order)

	// limit
	if p.Limit != constants.MaxUint64 {
		postsreq.WriteString(" LIMIT :limit OFFSET :offset")
	}

	// building count and select statements
	if cnstmt, db.err = db.PrepareNamed(precreq.String() + comreq.String()); db.err != nil {
		return nil, 0, db.err
	}
	if snstmt, db.err = db.PrepareNamed(presreq.String() + comreq.String() + postsreq.String()); db.err != nil {
		return nil, 0, db.err
	}

	// building argument map
	m := map[string]interface{}{
		"search":   fmt.Sprint("%", p.Search, "%"),
		"personid": p.LoggedPersonID,
		"entityid": p.EntityID,
		"order":    p.Order,
		"limit":    p.Limit,
		"offset":   p.Offset}

	// select
	if db.err = snstmt.Select(&storelocations, m); db.err != nil {
		return nil, 0, db.err
	}
	// count
	if db.err = cnstmt.Get(&count, m); db.err != nil {
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
	entity.entity_id AS "entity.entity_id",
	entity.entity_name AS "entity.entity_name"
	FROM storelocation AS s
	JOIN entity ON s.entity = entity.entity_id
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
	sqlr = `INSERT INTO storelocation(storelocation_name, entity) VALUES (?, ?)`
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
	sqlr = `UPDATE storelocation SET storelocation_name = ?, entity = ?
	WHERE storelocation_id = ?`
	if _, db.err = db.Exec(sqlr, s.StoreLocationName, s.Entity.EntityID, s.StoreLocationID); db.err != nil {
		return db.err
	}

	return nil
}

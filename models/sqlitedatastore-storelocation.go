package models

import (
	"strings"

	"database/sql"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/constants"
	"github.com/tbellembois/gochimitheque/helpers"
)

// GetStoreLocations returns the store locations matching the search criteria
// order, offset and limit are passed to the sql request
func (db *SQLiteDataStore) GetStoreLocations(p helpers.DbselectparamStoreLocation) ([]StoreLocation, int, error) {
	var (
		storelocations                     []StoreLocation
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
	)
	log.WithFields(log.Fields{"p": p}).Debug("GetStoreLocations")

	precreq.WriteString(" SELECT count(DISTINCT s.storelocation_id)")
	presreq.WriteString(` SELECT s.storelocation_id, s.storelocation_name, 
	entity.entity_id AS "entity.entity_id", 
	entity.entity_name AS "entity.entity_name"`)
	comreq.WriteString(" FROM storelocation AS s")
	comreq.WriteString(" JOIN entity ON s.entity = entity.entity_id")
	comreq.WriteString(` JOIN permission AS perm ON
	(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = entity.entity_id) OR
	(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "all" and perm.permission_entity_id = entity.entity_id) OR
	(perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "r" and perm.permission_entity_id = entity.entity_id)
	`)
	comreq.WriteString(" WHERE s.storelocation_name LIKE :search")
	if p.GetEntity() != -1 {
		comreq.WriteString(" AND s.entity = :entity")
	}
	postsreq.WriteString(" GROUP BY s.storelocation_id")
	postsreq.WriteString(" ORDER BY " + p.GetOrderBy() + " " + p.GetOrder())

	// limit
	if p.GetLimit() != constants.MaxUint64 {
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
		"search":   p.GetSearch(),
		"personid": p.GetLoggedPersonID(),
		"order":    p.GetOrder(),
		"limit":    p.GetLimit(),
		"offset":   p.GetOffset(),
		"entity":   p.GetEntity(),
	}

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

// GetStoreLocationEntity returns the entity of the store location with id "id"
func (db *SQLiteDataStore) GetStoreLocationEntity(id int) (Entity, error) {
	var (
		entity Entity
		sqlr   string
	)

	sqlr = `SELECT 
	entity.entity_id AS "entity_id",
	entity.entity_name AS "entity_name"
	FROM storelocation AS s
	JOIN entity ON s.entity = entity.entity_id
	WHERE s.storelocation_id = ?`
	if db.err = db.Get(&entity, sqlr, id); db.err != nil {
		return Entity{}, db.err
	}
	log.WithFields(log.Fields{"ID": id, "entity": entity}).Debug("GetStoreLocationEntity")
	return entity, nil
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

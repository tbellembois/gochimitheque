package models

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"reflect"
	"strconv"
	"strings"

	sq "github.com/Masterminds/squirrel"
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
		err                                error
	)
	log.WithFields(log.Fields{"p": p}).Debug("GetStoreLocations")

	precreq.WriteString(" SELECT count(DISTINCT s.storelocation_id)")
	presreq.WriteString(` SELECT s.storelocation_id AS "storelocation_id", s.storelocation_name AS "storelocation_name", s.storelocation_canstore, s.storelocation_color, 
	storelocation.storelocation_id AS "storelocation.storelocation_id",
	storelocation.storelocation_name AS "storelocation.storelocation_name",
	entity.entity_id AS "entity.entity_id", 
	entity.entity_name AS "entity.entity_name"`)
	comreq.WriteString(" FROM storelocation AS s")
	comreq.WriteString(" JOIN entity ON s.entity = entity.entity_id")
	comreq.WriteString(" LEFT JOIN storelocation on s.storelocation = storelocation.storelocation_id")
	comreq.WriteString(` JOIN permission AS perm ON
	(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = entity.entity_id) OR
	(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "storelocations" and perm.permission_perm_name = "all" and perm.permission_entity_id = entity.entity_id) OR
	(perm.person = :personid and perm.permission_item_name = "storelocations" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "storelocations" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "storelocations" and perm.permission_perm_name = "r" and perm.permission_entity_id = entity.entity_id)
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
	if cnstmt, err = db.PrepareNamed(precreq.String() + comreq.String()); err != nil {
		return nil, 0, err
	}
	if snstmt, err = db.PrepareNamed(presreq.String() + comreq.String() + postsreq.String()); err != nil {
		return nil, 0, err
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
	if err = snstmt.Select(&storelocations, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}
	return storelocations, count, nil
}

// GetStoreLocation returns the store location with id "id"
func (db *SQLiteDataStore) GetStoreLocation(id int) (StoreLocation, error) {
	var (
		storelocation StoreLocation
		sqlr          string
		err           error
	)

	sqlr = `SELECT s.storelocation_id, s.storelocation_name, s.storelocation_canstore, s.storelocation_color,
	storelocation.storelocation_id AS "storelocation.storelocation_id",
	storelocation.storelocation_name AS "storelocation.storelocation_name",
	entity.entity_id AS "entity.entity_id",
	entity.entity_name AS "entity.entity_name"
	FROM storelocation AS s
	JOIN entity ON s.entity = entity.entity_id
	JOIN storelocation on s.storelocation = storelocation.storelocation_id
	WHERE s.storelocation_id = ?`
	if err = db.Get(&storelocation, sqlr, id); err != nil {
		return StoreLocation{}, err
	}
	log.WithFields(log.Fields{"ID": id, "storelocation": storelocation}).Debug("GetStoreLocation")
	return storelocation, nil
}

// GetStoreLocationEntity returns the entity of the store location with id "id"
func (db *SQLiteDataStore) GetStoreLocationEntity(id int) (Entity, error) {
	var (
		entity Entity
		sqlr   string
		err    error
	)

	sqlr = `SELECT 
	entity.entity_id AS "entity_id",
	entity.entity_name AS "entity_name"
	FROM storelocation AS s
	JOIN entity ON s.entity = entity.entity_id
	WHERE s.storelocation_id = ?`
	if err = db.Get(&entity, sqlr, id); err != nil {
		return Entity{}, err
	}
	log.WithFields(log.Fields{"ID": id, "entity": entity}).Debug("GetStoreLocationEntity")
	return entity, nil
}

// DeleteStoreLocation deletes the store location with id "id"
func (db *SQLiteDataStore) DeleteStoreLocation(id int) error {
	var (
		sqlr string
		err  error
	)
	sqlr = `DELETE FROM storelocation 
	WHERE storelocation_id = ?`
	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}
	return nil
}

// CreateStoreLocation creates the given store location
func (db *SQLiteDataStore) CreateStoreLocation(s StoreLocation) (error, int) {
	var (
		sqlr     string
		res      sql.Result
		lastid   int64
		err      error
		sqla     []interface{}
		tx       *sql.Tx
		ibuilder sq.InsertBuilder
	)

	m := make(map[string]interface{})
	if s.StoreLocationCanStore.Valid {
		m["storelocation_canstore"] = s.StoreLocationCanStore.Bool
	}
	if s.StoreLocationColor.Valid {
		m["storelocation_color"] = s.StoreLocationColor.String
	}
	m["storelocation_name"] = s.StoreLocationName
	m["entity"] = s.EntityID

	// building column names/values
	col := make([]string, 0, len(m))
	val := make([]interface{}, 0, len(m))
	for k, v := range m {
		col = append(col, k)
		rt := reflect.TypeOf(v)
		rv := reflect.ValueOf(v)
		switch rt.Kind() {
		case reflect.Int:
			val = append(val, strconv.Itoa(int(rv.Int())))
		case reflect.String:
			val = append(val, rv.String())
		case reflect.Bool:
			val = append(val, rv.Bool())
		default:
			panic("unknown type:" + rt.String())
		}
	}

	// beginning transaction
	if tx, err = db.Begin(); err != nil {
		return err, 0
	}

	ibuilder = sq.Insert("storelocation").Columns(col...).Values(val...)
	if sqlr, sqla, err = ibuilder.ToSql(); err != nil {
		tx.Rollback()
		return err, 0
	}

	if res, err = tx.Exec(sqlr, sqla...); err != nil {
		tx.Rollback()
		return err, 0
	}

	// committing changes
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err, 0
	}

	// getting the last inserted id
	if lastid, err = res.LastInsertId(); err != nil {
		return err, 0
	}

	return nil, int(lastid)
}

// UpdateStoreLocation updates the given store location
func (db *SQLiteDataStore) UpdateStoreLocation(s StoreLocation) error {
	var (
		sqlr     string
		sqla     []interface{}
		tx       *sql.Tx
		err      error
		ubuilder sq.UpdateBuilder
	)
	log.WithFields(log.Fields{"s": s}).Debug("UpdateStoreLocation")

	m := make(map[string]interface{})
	if s.StoreLocationCanStore.Valid {
		m["storelocation_canstore"] = s.StoreLocationCanStore.Bool
	}
	if s.StoreLocationColor.Valid {
		m["storelocation_color"] = s.StoreLocationColor.String
	}
	m["storelocation_name"] = s.StoreLocationName
	m["storelocation"] = s.StoreLocation.StoreLocationID
	m["entity"] = s.EntityID

	// building column names/values
	col := make([]string, 0, len(m))
	val := make([]interface{}, 0, len(m))
	for k, v := range m {
		col = append(col, k)
		rt := reflect.TypeOf(v)
		rv := reflect.ValueOf(v)
		switch rt.Kind() {
		case reflect.Int:
			val = append(val, strconv.Itoa(int(rv.Int())))
		case reflect.String:
			val = append(val, rv.String())
		case reflect.Bool:
			val = append(val, rv.Bool())
		default:
			panic("unknown type:" + rt.String())
		}
	}

	// beginning transaction
	if tx, err = db.Begin(); err != nil {
		return err
	}

	ubuilder = sq.Update("storelocation").
		SetMap(m).
		Where(sq.Eq{"storelocation_id": s.StoreLocationID})
	if sqlr, sqla, err = ubuilder.ToSql(); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec(sqlr, sqla...); err != nil {
		tx.Rollback()
		return err
	}

	// committing changes
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

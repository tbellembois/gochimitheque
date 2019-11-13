package models

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/constants"
	"github.com/tbellembois/gochimitheque/helpers"
)

// buildFullPath builds the store location full path
// the caller is responsible of opening and commiting the tx transaction
func (db *SQLiteDataStore) buildFullPath(s StoreLocation, tx *sqlx.Tx) string {
	// parent
	var (
		pp  StoreLocation
		err error
	)

	log.WithFields(log.Fields{"s": s}).Debug("buildFullPath")

	// getting the parent
	if s.StoreLocation != nil && s.StoreLocation.StoreLocationID.Valid {
		// retrieving the parent from db
		sqlr := `SELECT s.storelocation_id, s.storelocation_name,
		storelocation.storelocation_id AS "storelocation.storelocation_id",
		storelocation.storelocation_name AS "storelocation.storelocation_name" 
		FROM storelocation AS s
		LEFT JOIN storelocation on s.storelocation = storelocation.storelocation_id
		WHERE s.storelocation_id = ?`
		r := tx.QueryRowx(sqlr, s.StoreLocation.StoreLocationID.Int64)
		if err = r.StructScan(&pp); err != nil {
			log.Error(err)
			return ""
		}

		// prepending the path with the parent name
		return db.buildFullPath(pp, tx) + "/" + s.StoreLocationName.String
	}

	return s.StoreLocationName.String
}

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
	presreq.WriteString(` SELECT s.storelocation_id AS "storelocation_id", 
	s.storelocation_name AS "storelocation_name", 
	s.storelocation_canstore, 
	s.storelocation_color, 
	s.storelocation_fullpath AS "storelocation_fullpath",
	storelocation.storelocation_id AS "storelocation.storelocation_id",
	storelocation.storelocation_name AS "storelocation.storelocation_name",
	entity.entity_id AS "entity.entity_id", 
	entity.entity_name AS "entity.entity_name"`)
	comreq.WriteString(" FROM storelocation AS s")
	comreq.WriteString(" JOIN entity ON s.entity = entity.entity_id")
	comreq.WriteString(" LEFT JOIN storelocation on s.storelocation = storelocation.storelocation_id")

	// filter by permissions
	// comreq.WriteString(` JOIN permission AS perm ON
	// (perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = entity.entity_id) OR
	// (perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
	// (perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
	// (perm.person = :personid and perm.permission_item_name = "storages" and perm.permission_perm_name = "all" and perm.permission_entity_id = entity.entity_id) OR
	// (perm.person = :personid and perm.permission_item_name = "storages" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
	// (perm.person = :personid and perm.permission_item_name = "storages" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
	// (perm.person = :personid and perm.permission_item_name = "storages" and perm.permission_perm_name = "r" and perm.permission_entity_id = entity.entity_id)
	// `)
	comreq.WriteString(` JOIN permission AS perm ON
	perm.person = :personid and (perm.permission_item_name in ("all", "storages")) and (perm.permission_perm_name in ("all", :permission)) and (perm.permission_entity_id in (-1, entity.entity_id))
	`)
	comreq.WriteString(" WHERE s.storelocation_name LIKE :search")
	if p.GetEntity() != -1 {
		comreq.WriteString(" AND s.entity = :entity")
	}
	if p.GetStoreLocationCanStore() {
		comreq.WriteString(" AND s.storelocation_canstore = :storelocation_canstore")
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
		"search":                 p.GetSearch(),
		"storelocation_canstore": p.GetStoreLocationCanStore(),
		"personid":               p.GetLoggedPersonID(),
		"order":                  p.GetOrder(),
		"limit":                  p.GetLimit(),
		"offset":                 p.GetOffset(),
		"entity":                 p.GetEntity(),
		"permission":             p.GetPermission(),
	}
	//log.Debug(presreq.String() + comreq.String() + postsreq.String())
	//log.Debug(m)

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
	log.WithFields(log.Fields{"id": id}).Debug("GetStoreLocation")

	sqlr = `SELECT s.storelocation_id, s.storelocation_name, s.storelocation_canstore, s.storelocation_color, s.storelocation_fullpath,
	storelocation.storelocation_id AS "storelocation.storelocation_id",
	storelocation.storelocation_name AS "storelocation.storelocation_name",
	entity.entity_id AS "entity.entity_id",
	entity.entity_name AS "entity.entity_name"
	FROM storelocation AS s
	JOIN entity ON s.entity = entity.entity_id
	LEFT JOIN storelocation on s.storelocation = storelocation.storelocation_id
	WHERE s.storelocation_id = ?`
	if err = db.Get(&storelocation, sqlr, id); err != nil {
		return StoreLocation{}, err
	}

	log.WithFields(log.Fields{"ID": id, "storelocation": storelocation}).Debug("GetStoreLocation")
	return storelocation, nil
}

// GetStoreLocationChildren returns the children of the store location with id "id"
func (db *SQLiteDataStore) GetStoreLocationChildren(id int) ([]StoreLocation, error) {
	var (
		storelocations []StoreLocation
		sqlr           string
		err            error
	)

	sqlr = `SELECT s.storelocation_id, s.storelocation_name, s.storelocation_canstore, s.storelocation_color, s.storelocation_fullpath,
	storelocation.storelocation_id AS "storelocation.storelocation_id",
	storelocation.storelocation_name AS "storelocation.storelocation_name",
	entity.entity_id AS "entity.entity_id",
	entity.entity_name AS "entity.entity_name"
	FROM storelocation AS s
	JOIN entity ON s.entity = entity.entity_id
	LEFT JOIN storelocation on s.storelocation = storelocation.storelocation_id
	WHERE s.storelocation = ?`
	if err = db.Select(&storelocations, sqlr, id); err != nil {
		return []StoreLocation{}, err
	}

	log.WithFields(log.Fields{"id": id, "storelocations": storelocations}).Debug("GetStoreLocationChildren")
	return storelocations, nil
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
func (db *SQLiteDataStore) CreateStoreLocation(s StoreLocation) (int, error) {
	var (
		sqlr     string
		res      sql.Result
		lastid   int64
		err      error
		sqla     []interface{}
		tx       *sqlx.Tx
		ibuilder sq.InsertBuilder
	)

	// beginning transaction
	if tx, err = db.Beginx(); err != nil {
		return 0, nil
	}

	// building full path
	s.StoreLocationFullPath = db.buildFullPath(s, tx)

	m := make(map[string]interface{})
	if s.StoreLocationCanStore.Valid {
		m["storelocation_canstore"] = s.StoreLocationCanStore.Bool
	}
	if s.StoreLocationColor.Valid {
		m["storelocation_color"] = s.StoreLocationColor.String
	}
	m["storelocation_name"] = s.StoreLocationName.String
	if s.StoreLocation != nil {
		m["storelocation"] = s.StoreLocation.StoreLocationID.Int64
	}
	m["entity"] = s.EntityID
	m["storelocation_fullpath"] = s.StoreLocationFullPath

	// building column names/values
	col := make([]string, 0, len(m))
	val := make([]interface{}, 0, len(m))
	// for k, v := range m {
	// 	col = append(col, k)
	// 	rt := reflect.TypeOf(v)
	// 	rv := reflect.ValueOf(v)
	// 	switch rt.Kind() {
	// 	case reflect.Int, reflect.Int64:
	// 		val = append(val, strconv.Itoa(int(rv.Int())))
	// 	case reflect.String:
	// 		val = append(val, rv.String())
	// 	case reflect.Bool:
	// 		val = append(val, rv.Bool())
	// 	default:
	// 		panic("unknown type:" + rt.String())
	// 	}
	// }
	for k, v := range m {
		col = append(col, k)

		switch t := v.(type) {
		case int:
			val = append(val, v.(int))
		case string:
			val = append(val, v.(string))
		case bool:
			val = append(val, v.(bool))
		default:
			panic(fmt.Sprintf("unknown type: %T", t))
		}
	}

	ibuilder = sq.Insert("storelocation").Columns(col...).Values(val...)
	if sqlr, sqla, err = ibuilder.ToSql(); err != nil {
		tx.Rollback()
		return 0, nil
	}

	if res, err = tx.Exec(sqlr, sqla...); err != nil {
		tx.Rollback()
		return 0, nil
	}

	// committing changes
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return 0, nil
	}

	// getting the last inserted id
	if lastid, err = res.LastInsertId(); err != nil {
		return 0, nil
	}

	return int(lastid), nil
}

// UpdateStoreLocation updates the given store location
func (db *SQLiteDataStore) UpdateStoreLocation(s StoreLocation) error {
	var (
		sqlr     string
		sqla     []interface{}
		tx       *sqlx.Tx
		err      error
		ubuilder sq.UpdateBuilder
	)

	// beginning new transaction
	if tx, err = db.Beginx(); err != nil {
		return err
	}

	// building full path
	s.StoreLocationFullPath = db.buildFullPath(s, tx)

	m := make(map[string]interface{})
	if s.StoreLocationCanStore.Valid {
		m["storelocation_canstore"] = s.StoreLocationCanStore.Bool
	}
	if s.StoreLocationColor.Valid {
		m["storelocation_color"] = s.StoreLocationColor.String
	}
	m["storelocation_name"] = s.StoreLocationName.String
	if s.StoreLocation != nil {
		m["storelocation"] = s.StoreLocation.StoreLocationID.Int64
	}
	m["entity"] = s.EntityID
	m["storelocation_fullpath"] = s.StoreLocationFullPath

	// building column names/values
	col := make([]string, 0, len(m))
	val := make([]interface{}, 0, len(m))
	for k, v := range m {
		col = append(col, k)
		rt := reflect.TypeOf(v)
		rv := reflect.ValueOf(v)
		switch rt.Kind() {
		case reflect.Int, reflect.Int64:
			val = append(val, strconv.Itoa(int(rv.Int())))
		case reflect.String:
			val = append(val, rv.String())
		case reflect.Bool:
			val = append(val, rv.Bool())
		default:
			panic("unknown type:" + rt.String())
		}
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

// IsStoreLocationEmpty returns true is the store location is empty
func (db *SQLiteDataStore) IsStoreLocationEmpty(id int) (bool, error) {
	var (
		res   bool
		count int
		sqlr  string
		err   error
	)

	sqlr = "SELECT count(*) from storages WHERE  storelocation = ?"
	if err = db.Get(&count, sqlr, id); err != nil {
		return false, err
	}
	log.WithFields(log.Fields{"id": id, "count": count}).Debug("IsStoreLocationEmpty")
	if count == 0 {
		res = true
	} else {
		res = false
	}
	return res, nil
}

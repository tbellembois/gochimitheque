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

func (db *SQLiteDataStore) GetStorages(p helpers.DbselectparamStorage) ([]Storage, int, error) {
	var (
		storages                           []Storage
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
	)
	log.WithFields(log.Fields{"p": p}).Debug("GetStorages")

	// pre request: select or count
	precreq.WriteString(" SELECT count(DISTINCT storage.storage_id)")
	presreq.WriteString(` SELECT storage.storage_id,
		storage.storage_creationdate AS "storage_creationdate",
		storage.storage_comment,
		person.person_email AS "person.person_email", 
		product.product_id AS "product.product_id",
		name.name_label AS "product.name.name_label",	 
		storelocation.storelocation_name AS "storelocation.storelocation_name"
		`)

	// common parts
	comreq.WriteString(" FROM storage")
	// get product
	comreq.WriteString(" JOIN product ON storage.product = product.product_id")
	// get names
	comreq.WriteString(" JOIN name ON product.name = name.name_id")
	// get person
	comreq.WriteString(" JOIN person ON storage.person = person.person_id")
	// get store location
	comreq.WriteString(" JOIN storelocation ON storage.storelocation = storelocation.storelocation_id")
	// filter by permissions
	comreq.WriteString(` JOIN permission AS perm, entity as e ON
		(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
		(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
		(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
		(perm.person = :personid and perm.permission_item_name = "storages" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
		(perm.person = :personid and perm.permission_item_name = "storages" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
		(perm.person = :personid and perm.permission_item_name = "storages" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
		(perm.person = :personid and perm.permission_item_name = "storages" and perm.permission_perm_name = "r" and perm.permission_entity_id = e.entity_id)
		`)
	comreq.WriteString(" WHERE storage.storage_id LIKE :search")
	if p.GetProduct() != -1 {
		comreq.WriteString(" AND product.product_id = :product")
	}

	// post select request
	postsreq.WriteString(" GROUP BY storage.storage_id")
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
		"product":  p.GetProduct(),
	}

	// select
	if db.err = snstmt.Select(&storages, m); db.err != nil {
		return nil, 0, db.err
	}
	// count
	if db.err = cnstmt.Get(&count, m); db.err != nil {
		return nil, 0, db.err
	}

	return storages, count, nil
}

// GetStorage returns the storage with id "id"
func (db *SQLiteDataStore) GetStorage(id int) (Storage, error) {
	var (
		storage Storage
		sqlr    string
	)
	log.WithFields(log.Fields{"id": id}).Debug("GetStorage")

	sqlr = `SELECT storage.storage_id, storage.storage_creationdate, storage.storage_comment,
	person.person_email AS "person.person_email",
	name.name_label AS "product.name.name_label",
	casnumber.casnumber_label AS "product.casnumber.casnumber_label",
	storelocation.storelocation_id AS "storelocation.storelocation_id",
	storelocation.storelocation_name AS "storelocation.storelocation_name"
	FROM storage
	JOIN storelocation ON storage.storelocation = storelocation.storelocation_id
	JOIN person ON storage.person = person.person_id
	JOIN product ON storage.product = product.product_id
	JOIN casnumber ON product.casnumber = casnumber.casnumber_id
	JOIN name ON product.name = name.name_id
	WHERE storage.storage_id = ?`
	if db.err = db.Get(&storage, sqlr, id); db.err != nil {
		return Storage{}, db.err
	}
	log.WithFields(log.Fields{"ID": id, "storage": storage}).Debug("GetStorage")
	return storage, nil
}

func (db *SQLiteDataStore) DeleteStorage(id int) error {
	var (
		sqlr string
	)
	sqlr = `DELETE FROM storage 
	WHERE storage_id = ?`
	if _, db.err = db.Exec(sqlr, id); db.err != nil {
		return db.err
	}
	return nil
}

func (db *SQLiteDataStore) CreateStorage(s Storage) (error, int) {

	var (
		sqlr   string
		res    sql.Result
		lastid int64
	)
	// FIXME: use a transaction here
	sqlr = `INSERT INTO storage(storage_creationdate, storage_comment, person, product, storelocation) VALUES (?, ?, ?, ?, ?)`
	if res, db.err = db.Exec(sqlr, s.StorageCreationDate, s.StorageComment, s.PersonID, s.ProductID, s.StoreLocationID); db.err != nil {
		return db.err, 0
	}

	// getting the last inserted id
	if lastid, db.err = res.LastInsertId(); db.err != nil {
		return db.err, 0
	}

	return nil, int(lastid)
}
func (db *SQLiteDataStore) UpdateStorage(s Storage) error {

	var (
		sqlr string
	)

	// updating the storage - product not supposed to be changed
	sqlr = `UPDATE storage SET storage_comment = ?, person = ?, storelocation = ?
	WHERE storage_id = ?`
	if _, db.err = db.Exec(sqlr, s.StorageComment, s.PersonID, s.StoreLocationID, s.StorageID); db.err != nil {
		return db.err
	}

	return nil
}

package models

import (
	"strings"

	"database/sql"

	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/constants"
	"github.com/tbellembois/gochimitheque/helpers"
)

// GetStoragesUnits return the units matching the search criteria
func (db *SQLiteDataStore) GetStoragesUnits(p helpers.Dbselectparam) ([]Unit, int, error) {
	var (
		units                              []Unit
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	precreq.WriteString(" SELECT count(DISTINCT unit.unit_id)")
	presreq.WriteString(" SELECT unit_id, unit_label")

	comreq.WriteString(" FROM unit")
	comreq.WriteString(" WHERE unit_label LIKE :search")
	postsreq.WriteString(" ORDER BY unit_label  " + p.GetOrder())

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
		"search": p.GetSearch(),
		"order":  p.GetOrder(),
		"limit":  p.GetLimit(),
		"offset": p.GetOffset(),
	}

	// select
	if err = snstmt.Select(&units, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	log.WithFields(log.Fields{"units": units}).Debug("GetStoragesUnits")
	return units, count, nil
}

// GetStoragesSuppliers return the suppliers matching the search criteria
func (db *SQLiteDataStore) GetStoragesSuppliers(p helpers.Dbselectparam) ([]Supplier, int, error) {
	var (
		suppliers                          []Supplier
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	precreq.WriteString(" SELECT count(DISTINCT supplier.supplier_id)")
	presreq.WriteString(" SELECT supplier_id, supplier_label")

	comreq.WriteString(" FROM supplier")
	comreq.WriteString(" WHERE supplier_label LIKE :search")
	postsreq.WriteString(" ORDER BY supplier_label  " + p.GetOrder())

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
		"search": p.GetSearch(),
		"order":  p.GetOrder(),
		"limit":  p.GetLimit(),
		"offset": p.GetOffset(),
	}

	// select
	if err = snstmt.Select(&suppliers, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for formula matching exactly the search
	s := p.GetSearch()
	s = strings.TrimPrefix(s, "%")
	s = strings.TrimSuffix(s, "%")
	var supplier Supplier

	r := db.QueryRowx(`SELECT supplier_id, supplier_label FROM supplier WHERE supplier_label == ?`, s)
	if err = r.StructScan(&supplier); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	} else {
		for i, s := range suppliers {
			if s.SupplierID == supplier.SupplierID {
				suppliers[i].C = 1
			}
		}
	}

	log.WithFields(log.Fields{"suppliers": suppliers}).Debug("GetStoragesSuppliers")
	return suppliers, count, nil
}

func (db *SQLiteDataStore) GetStorages(p helpers.DbselectparamStorage) ([]Storage, int, error) {
	var (
		storages                           []Storage
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)
	log.WithFields(log.Fields{"p": p}).Debug("GetStorages")

	// pre request: select or count
	precreq.WriteString(" SELECT count(DISTINCT storage.storage_id)")
	presreq.WriteString(` SELECT storage.storage_id,
		storage.storage_creationdate,
		storage.storage_quantity,
		storage.storage_barecode,
		storage.storage_comment,
		unit.unit_label AS "unit.unit_label",
		supplier.supplier_label AS "supplier.supplier_label",
		person.person_email AS "person.person_email", 
		product.product_id AS "product.product_id",
		name.name_label AS "product.name.name_label",	 
		storelocation.storelocation_name AS "storelocation.storelocation_name",
		storelocation.storelocation_color AS "storelocation.storelocation_color",
		storelocation.storelocation_fullpath AS "storelocation.storelocation_fullpath",
		entity.entity_id AS "storelocation.entity.entity_id"
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
	// get entity
	comreq.WriteString(" JOIN entity ON storelocation.entity = entity.entity_id")
	// get unit
	comreq.WriteString(" LEFT JOIN unit ON storage.unit = unit.unit_id")
	// get supplier
	comreq.WriteString(" LEFT JOIN supplier ON storage.supplier = supplier.supplier_id")
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
	comreq.WriteString(" WHERE (storelocation.storelocation_fullpath LIKE :search OR name.name_label LIKE :search)")
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
		"product":  p.GetProduct(),
	}

	// select
	if err = snstmt.Select(&storages, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	return storages, count, nil
}

// GetStorage returns the storage with id "id"
func (db *SQLiteDataStore) GetStorage(id int) (Storage, error) {
	var (
		storage Storage
		sqlr    string
		err     error
	)
	log.WithFields(log.Fields{"id": id}).Debug("GetStorage")

	sqlr = `SELECT storage.storage_id,
	storage.storage_creationdate,
	storage.storage_quantity,
	storage.storage_barecode,
	storage.storage_comment,
	unit.unit_label AS "unit.unit_label",
	supplier.supplier_label AS "supplier.supplier_label",
	person.person_email AS "person.person_email",
	name.name_label AS "product.name.name_label",
	casnumber.casnumber_label AS "product.casnumber.casnumber_label",
	storelocation.storelocation_id AS "storelocation.storelocation_id",
	storelocation.storelocation_name AS "storelocation.storelocation_name",
	storelocation.storelocation_color AS "storelocation.storelocation_color",
	storelocation.storelocation_fullpath AS "storelocation.storelocation_fullpath"
	FROM storage
	JOIN storelocation ON storage.storelocation = storelocation.storelocation_id
	LEFT JOIN unit ON storage.unit = unit.unit_id
	LEFT JOIN supplier ON storage.supplier = supplier.supplier_id
	JOIN person ON storage.person = person.person_id
	JOIN product ON storage.product = product.product_id
	JOIN casnumber ON product.casnumber = casnumber.casnumber_id
	JOIN name ON product.name = name.name_id
	WHERE storage.storage_id = ?`
	if err = db.Get(&storage, sqlr, id); err != nil {
		return Storage{}, err
	}
	log.WithFields(log.Fields{"ID": id, "storage": storage}).Debug("GetStorage")
	return storage, nil
}

func (db *SQLiteDataStore) DeleteStorage(id int) error {
	var (
		sqlr string
		err  error
	)
	sqlr = `DELETE FROM storage 
	WHERE storage_id = ?`
	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}
	return nil
}

func (db *SQLiteDataStore) CreateStorage(s Storage) (error, int) {

	var (
		sqlr   string
		res    sql.Result
		lastid int64
		err    error
	)
	// FIXME: use a transaction here
	sqlr = `INSERT INTO storage(storage_creationdate, storage_comment, person, product, storelocation) VALUES (?, ?, ?, ?, ?)`
	if res, err = db.Exec(sqlr, s.StorageCreationDate, s.StorageComment, s.PersonID, s.ProductID, s.StoreLocationID); err != nil {
		return err, 0
	}

	// getting the last inserted id
	if lastid, err = res.LastInsertId(); err != nil {
		return err, 0
	}

	return nil, int(lastid)
}
func (db *SQLiteDataStore) UpdateStorage(s Storage) error {

	var (
		sqlr     string
		err      error
		tx       *sql.Tx
		res      sql.Result
		lastid   int64
		sqla     []interface{}
		ubuilder sq.UpdateBuilder
	)

	// beginning transaction
	if tx, err = db.Begin(); err != nil {
		return err
	}

	// if SupplierID = -1 then it is a new supplier
	if v, err := s.Supplier.SupplierID.Value(); s.Supplier.SupplierID.Valid && err == nil && v.(int64) == -1 {
		sqlr = `INSERT INTO supplier (supplier_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, s.Supplier.SupplierLabel); err != nil {
			tx.Rollback()
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return err
		}
		// updating the storage SupplierId (SupplierLabel already set)
		s.Supplier.SupplierID = sql.NullInt64{Int64: lastid}
	}
	if err != nil {
		log.Error("supplier error - " + err.Error())
		tx.Rollback()
		return err
	}

	// finally updating the storage
	m := make(map[string]interface{})
	if s.StorageComment.Valid {
		m["storage_comment"] = s.StorageComment.String
	}
	if s.StorageQuantity.Valid {
		m["storage_quantity"] = s.StorageQuantity.Float64
	}
	if s.StorageBarecode.Valid {
		m["storage_barecode"] = s.StorageBarecode.String
	}
	m["person"] = s.PersonID
	m["storelocation"] = s.StoreLocationID
	m["unit"] = s.UnitID
	m["supplier"] = s.SupplierID

	ubuilder = sq.Update("storage").
		SetMap(m).
		Where(sq.Eq{"storage_id": s.StorageID})
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

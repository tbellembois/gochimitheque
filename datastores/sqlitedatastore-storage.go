package datastores

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"

	// register sqlite3 driver.
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/zmqclient"
)

// ToogleStorageBorrowing borrow/unborrow the storage for the connected user.
func (db *SQLiteDataStore) ToogleStorageBorrowing(s models.Storage) error {
	var (
		sqlr  string
		count int
		err   error
	)

	sqlr = `SELECT COUNT(borrowing_id) FROM borrowing WHERE storage = ?`
	if err = db.Get(&count, sqlr, s.StorageID.Int64); err != nil {
		return err
	}

	if count == 0 {
		sqlr = `INSERT into borrowing(person, storage, borrower, borrowing_comment) VALUES (?, ?, ?, ?)`
		if _, err = db.Exec(sqlr, s.Borrowing.Person.PersonID, s.StorageID.Int64, s.Borrowing.Borrower.PersonID, s.Borrowing.BorrowingComment); err != nil {
			return err
		}
	} else {
		sqlr = `DELETE from borrowing WHERE storage = ?`
		if _, err = db.Exec(sqlr, s.StorageID.Int64); err != nil {
			return err
		}
	}

	return nil
}

// GetStoragesUnits return the units.
func (db *SQLiteDataStore) GetStoragesUnits(f zmqclient.RequestFilter) ([]models.Unit, int, error) {
	var (
		units                              []models.Unit
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	// hack to bypass optionnal default on the Rust part.
	if f.Search == "" {
		f.Search = "%%"
	}

	if f.OrderBy == "" {
		f.OrderBy = "unit_id"
	}

	precreq.WriteString(" SELECT count(DISTINCT unit.unit_id)")
	presreq.WriteString(" SELECT unit_id, unit_label, unit_type")

	comreq.WriteString(" FROM unit")
	comreq.WriteString(" WHERE unit_label LIKE :search")

	if f.UnitType != "" {
		comreq.WriteString(" AND unit_type=:unit_type")
	}

	postsreq.WriteString(" ORDER BY unit.unit_type, unit_id  " + f.Order)

	// limit
	if f.Limit != ^uint64(0) {
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
		"search":    f.Search,
		"order":     f.Order,
		"limit":     f.Limit,
		"offset":    f.Offset,
		"unit_type": f.UnitType,
	}

	// Select.
	if err = snstmt.Select(&units, m); err != nil {
		return nil, 0, err
	}
	// Count.
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	logger.Log.WithFields(logrus.Fields{"units": units}).Debug("GetStoragesUnits")

	return units, count, nil
}

// GetStorages returns the storages matching the request parameters p.
// Only storages that the logged user can see are returned given his permissions
// and membership.
func (db *SQLiteDataStore) GetStorages(f zmqclient.RequestFilter, person_id int) ([]models.Storage, int, error) {
	var (
		storages                                  []models.Storage
		count                                     int
		precreq, presreq, comreq, postsreq, reqhc strings.Builder
		cnstmt                                    *sqlx.NamedStmt
		snstmt                                    *sqlx.NamedStmt
		err                                       error
		isadmin                                   bool
	)

	logger.Log.WithFields(logrus.Fields{"f": f}).Debug("GetStorages")

	// hack to bypass optionnal default on the Rust part.
	if f.Search == "" {
		f.Search = "%%"
	}

	if f.OrderBy == "" {
		f.OrderBy = "storage_id"
	} else if f.OrderBy == "product.name.name_label" {
		f.OrderBy = "name.name_label"
	} else if strings.HasPrefix(f.OrderBy, "storage_") {
		f.OrderBy = fmt.Sprintf("s.%s", f.OrderBy)
	}

	// is the user an admin?
	if isadmin, err = db.IsPersonAdmin(person_id); err != nil {
		return nil, 0, err
	}

	// pre request: select or count
	precreq.WriteString(" SELECT count(DISTINCT s.storage_id)")
	presreq.WriteString(` SELECT s.storage_id AS "storage_id",
		s.storage_entrydate,
		s.storage_exitdate,
		s.storage_openingdate,
		s.storage_expirationdate,
		s.storage_reference,
		s.storage_batchnumber,
		s.storage_todestroy,
		s.storage_creationdate,
		s.storage_modificationdate,
		s.storage_quantity,
		s.storage_barecode,
		s.storage_qrcode,
		s.storage_comment,
		s.storage_archive,
		s.storage_concentration,
		s.storage_number_of_carton,
		s.storage_number_of_bag,
		s.storage_number_of_unit,
		storage.storage_id AS "storage.storage_id",
		uq.unit_id AS "unit_quantity.unit_id",
		uq.unit_label AS "unit_quantity.unit_label",
		uc.unit_id AS "unit_concentration.unit_id",
		uc.unit_label AS "unit_concentration.unit_label",
		supplier.supplier_id AS "supplier.supplier_id",
		supplier.supplier_label AS "supplier.supplier_label",
		person.person_id AS "person.person_id", 
		person.person_email AS "person.person_email", 
		product.product_id AS "product.product_id",
		product.product_specificity AS "product.product_specificity",
		product.product_number_per_carton AS "product.product_number_per_carton",
		product.product_number_per_bag AS "product.product_number_per_bag",
        producerref.producerref_id AS "product.producerref.producerref_id",
		producerref.producerref_label AS "product.producerref.producerref_label",
		name.name_id AS "product.name.name_id",
		name.name_label AS "product.name.name_label",
		casnumber.casnumber_id AS "product.casnumber.casnumber_id",
		casnumber.casnumber_label AS "product.casnumber.casnumber_label",
		borrowing.borrowing_id AS "borrowing.borrowing_id",
		borrowing.borrowing_comment AS "borrowing.borrowing_comment",
		storelocation.storelocation_id AS "storelocation.storelocation_id",
		storelocation.storelocation_name AS "storelocation.storelocation_name",
		storelocation.storelocation_color AS "storelocation.storelocation_color",
		storelocation.storelocation_fullpath AS "storelocation.storelocation_fullpath",
		entity.entity_id AS "storelocation.entity.entity_id"
		`)

	if f.CasNumberCmr {
		presreq.WriteString(`,GROUP_CONCAT(DISTINCT hazardstatement.hazardstatement_cmr) AS "product.hazardstatement_cmr"`)
	}

	// common parts
	comreq.WriteString(" FROM storage as s")
	// get storage history parent
	comreq.WriteString(" LEFT JOIN storage ON s.storage = storage.storage_id")
	// get product
	comreq.WriteString(" JOIN product ON s.product = product.product_id")
	// CMR
	if f.CasNumberCmr {
		comreq.WriteString(" LEFT JOIN producthazardstatements ON producthazardstatements.producthazardstatements_product_id = product.product_id")
		comreq.WriteString(" LEFT JOIN hazardstatement ON producthazardstatements.producthazardstatements_hazardstatement_id = hazardstatement.hazardstatement_id")
	}
	// get producerref
	if f.ProducerRef != 0 {
		comreq.WriteString(" JOIN producerref ON product.producerref = :producerref")
	} else {
		comreq.WriteString(" LEFT JOIN producerref ON product.producerref = producerref.producerref_id")
	}
	// get name
	comreq.WriteString(" JOIN name ON product.name = name.name_id")
	// get category
	if f.Category != 0 {
		comreq.WriteString(" JOIN category ON product.category = :category")
	}
	// get signal word
	comreq.WriteString(" LEFT JOIN signalword ON product.signalword = signalword.signalword_id")
	// get person
	comreq.WriteString(" JOIN person ON s.person = person.person_id")
	// get store location
	comreq.WriteString(" JOIN storelocation ON s.storelocation = storelocation.storelocation_id")
	// get entity
	comreq.WriteString(" JOIN entity ON storelocation.entity = entity.entity_id")
	// get unit quantity
	comreq.WriteString(" LEFT JOIN unit uq ON s.unit_quantity = uq.unit_id")
	// get unit concentration
	comreq.WriteString(" LEFT JOIN unit uc ON s.unit_concentration = uc.unit_id")
	// get supplier
	comreq.WriteString(" LEFT JOIN supplier ON s.supplier = supplier.supplier_id")
	// get borrowings
	if f.Borrowing {
		comreq.WriteString(" JOIN borrowing ON borrowing.storage = s.storage_id AND borrowing.borrower = :personid")
	} else {
		comreq.WriteString(" LEFT JOIN borrowing ON s.storage_id = borrowing.storage")
	}
	// get casnumber
	comreq.WriteString(" LEFT JOIN casnumber ON product.casnumber = casnumber.casnumber_id")
	// get empirical formula
	comreq.WriteString(" LEFT JOIN empiricalformula ON product.empiricalformula = empiricalformula.empiricalformula_id")
	// get symbols
	if len(f.Symbols) != 0 {
		comreq.WriteString(" JOIN productsymbols AS ps ON ps.productsymbols_product_id = product.product_id")
	}
	// get hazardstatements
	if len(f.HazardStatements) != 0 {
		comreq.WriteString(" JOIN producthazardstatements AS phs ON phs.producthazardstatements_product_id = product.product_id")
	}
	// get precautionarystatements
	if len(f.PrecautionaryStatements) != 0 {
		comreq.WriteString(" JOIN productprecautionarystatements AS pps ON pps.productprecautionarystatements_product_id = product.product_id")
	}
	// get tags
	if len(f.Tags) != 0 {
		comreq.WriteString(" JOIN producttags AS ptags ON ptags.producttags_product_id = product.product_id")
	}
	// get bookmarks
	if f.Bookmark {
		comreq.WriteString(" JOIN bookmark AS b ON b.product = product.product_id AND b.person = :personid")
	}

	// filter by entities
	if !isadmin {
		comreq.WriteString(` JOIN personentities ON (personentities_entity_id = storelocation.entity AND personentities_person_id = :personid)`)
	}

	// filter by permissions
	comreq.WriteString(` JOIN permission AS perm, entity as e ON
	perm.person = :personid and (perm.permission_item_name in ("all", "storages")) and (perm.permission_perm_name in ("all", "r", "w")) and (perm.permission_entity_id in (-1, e.entity_id))
	`)
	comreq.WriteString(" WHERE 1")
	if len(f.Ids) > 0 {
		comreq.WriteString(" AND s.storage_id in (")

		for _, id := range f.Ids {
			comreq.WriteString(fmt.Sprintf("%d,", id))
		}
		// to complete the last comma
		comreq.WriteString("-1")
		comreq.WriteString(" )")
	}
	if f.StorageToDestroy {
		comreq.WriteString(" AND s.storage_todestroy = true")
	}
	if f.CasNumberCmr {
		comreq.WriteString(" AND (casnumber.casnumber_cmr IS NOT NULL OR (hazardstatement_cmr IS NOT NULL AND hazardstatement_cmr != ''))")
	}
	if f.Product != 0 {
		comreq.WriteString(" AND product.product_id = :product")
	}
	if f.Entity != 0 {
		comreq.WriteString(" AND entity.entity_id = :entity")
	}
	if f.Storelocation != 0 {
		comreq.WriteString(" AND storelocation.storelocation_id = :storelocation")
	}
	if f.Storage != 0 {
		if f.History {
			comreq.WriteString(" AND (s.storage = :storage OR s.storage_id = :storage)")
		} else {
			comreq.WriteString(" AND (s.storage_id = :storage")
			// getting storages with identical barecode
			comreq.WriteString(" OR (s.storage_barecode = (SELECT storage_barecode FROM storage WHERE storage_id = :storage)))")
		}
	}
	if !f.History {
		comreq.WriteString(" AND s.storage IS NULL")
	}
	if f.StorageArchive {
		comreq.WriteString(" AND s.storage_archive = true")
	} else {
		comreq.WriteString(" AND s.storage_archive = false")
	}

	// search form parameters
	if f.Name != 0 {
		comreq.WriteString(" AND name.name_id = :name")
	}
	if f.CasNumber != 0 {
		comreq.WriteString(" AND casnumber.casnumber_id = :casnumber")
	}
	if f.EmpiricalFormula != 0 {
		comreq.WriteString(" AND empiricalformula.empiricalformula_id = :empiricalformula")
	}
	if f.StorageBarecode != "" {
		comreq.WriteString(" AND s.storage_barecode LIKE :storage_barecode")
	}
	if f.StorageBatchNumber != "" {
		comreq.WriteString(" AND s.storage_batchnumber LIKE :storage_batchnumber")
	}
	if f.CustomNamePartOf != "" {
		comreq.WriteString(" AND name.name_label LIKE :custom_name_part_of")
	}
	if len(f.Symbols) != 0 {
		comreq.WriteString(" AND ps.productsymbols_symbol_id IN (")

		for _, s := range f.Symbols {
			comreq.WriteString(fmt.Sprintf("%d,", s))
		}

		// to complete the last comma
		comreq.WriteString("-1")
		comreq.WriteString(" )")
	}
	if len(f.HazardStatements) != 0 {
		comreq.WriteString(" AND phs.producthazardstatements_hazardstatement_id IN (")

		for _, s := range f.HazardStatements {
			comreq.WriteString(fmt.Sprintf("%d,", s))
		}

		// to complete the last comma
		comreq.WriteString("-1")
		comreq.WriteString(" )")
	}
	if len(f.PrecautionaryStatements) != 0 {
		comreq.WriteString(" AND pps.productprecautionarystatements_precautionarystatement_id IN (")

		for _, s := range f.PrecautionaryStatements {
			comreq.WriteString(fmt.Sprintf("%d,", s))
		}

		// to complete the last comma
		comreq.WriteString("-1")
		comreq.WriteString(" )")
	}
	if len(f.Tags) != 0 {
		comreq.WriteString(" AND ptags.producttags_tag_id IN (")

		for _, t := range f.Tags {
			comreq.WriteString(fmt.Sprintf("%d,", t))
		}

		// to complete the last comma
		comreq.WriteString("-1")
		comreq.WriteString(" )")
	}
	if f.SignalWord != 0 {
		comreq.WriteString(" AND signalword.signalword_id = :signalword")
	}

	// show bio/chem/consu
	switch {
	case !f.ShowChem && !f.ShowBio && f.ShowConsu:
		comreq.WriteString(" AND (product_number_per_carton IS NOT NULL AND product_number_per_carton != 0)")
	case !f.ShowChem && f.ShowBio && !f.ShowConsu:
		comreq.WriteString(" AND producerref IS NOT NULL")
		comreq.WriteString(" AND (product_number_per_carton IS NULL OR product_number_per_carton == 0)")
	case !f.ShowChem && f.ShowBio && f.ShowConsu:
		comreq.WriteString(" AND ((product_number_per_carton IS NOT NULL AND product_number_per_carton != 0)")
		comreq.WriteString(" OR producerref IS NOT NULL)")
	case f.ShowChem && !f.ShowBio && !f.ShowConsu:
		comreq.WriteString(" AND producerref IS NULL")
		comreq.WriteString(" AND (product_number_per_carton IS NULL OR product_number_per_carton == 0)")
	case f.ShowChem && !f.ShowBio && f.ShowConsu:
		comreq.WriteString(" AND (producerref IS NULL")
		comreq.WriteString(" OR (product_number_per_carton IS NOT NULL AND product_number_per_carton != 0))")
	case f.ShowChem && f.ShowBio && !f.ShowConsu:
		comreq.WriteString(" AND (product_number_per_carton IS NULL OR product_number_per_carton == 0)")
	}

	// post select request
	postsreq.WriteString(" GROUP BY s.storage_id")
	postsreq.WriteString(" ORDER BY " + f.OrderBy + " " + f.Order)

	// limit
	if f.Limit != ^uint64(0) {
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
		"ids":                 f.Ids,
		"search":              f.Search,
		"personid":            person_id,
		"order":               f.Order,
		"limit":               f.Limit,
		"offset":              f.Offset,
		"entity":              f.Entity,
		"product":             f.Product,
		"storelocation":       f.Storelocation,
		"storage":             f.Storage,
		"name":                f.Name,
		"casnumber":           f.CasNumber,
		"empiricalformula":    f.EmpiricalFormula,
		"storage_barecode":    f.StorageBarecode,
		"storage_batchnumber": f.StorageBatchNumber,
		"custom_name_part_of": "%" + f.CustomNamePartOf + "%",
		"signalword":          f.SignalWord,
		"producerref":         f.ProducerRef,
		"category":            f.Category,
	}

	logger.Log.Debug(presreq.String() + comreq.String() + postsreq.String())
	logger.Log.Debug(m)

	// Select.
	if err = snstmt.Select(&storages, m); err != nil {
		return nil, 0, err
	}
	// Count.
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	//
	// getting product type
	//
	for i, st := range storages {
		if st.Product.ProductNumberPerCarton.Valid {
			storages[i].Product.ProductType = "CONS"
		} else if st.Product.ProducerRef.ProducerRefID.Valid {
			storages[i].Product.ProductType = "BIO"
		} else {
			storages[i].Product.ProductType = "CHEM"
		}
	}

	//
	// getting number of history for each storage
	//
	for i, st := range storages {
		// getting the total storage count
		// logger.Log.Debug(st)
		reqhc.Reset()
		reqhc.WriteString("SELECT count(DISTINCT storage_id) from storage WHERE storage.storage = ?")
		if err = db.Get(&storages[i].StorageHC, reqhc.String(), st.StorageID); err != nil {
			return nil, 0, err
		}
	}

	//
	// getting borrower for each storage
	//
	for i, st := range storages {
		reqhc.Reset()
		reqhc.WriteString(`SELECT borrowing_id, 
		borrowing_comment, 
		person.person_email AS "borrower.person_email" 
		from borrowing 
		JOIN person 
		ON borrowing.borrower = person.person_id 
		WHERE borrowing.storage = ?`)

		var borrowing models.Borrowing

		if err = db.Get(&borrowing, reqhc.String(), st.StorageID); err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}

		storages[i].Borrowing = &borrowing
	}

	return storages, count, nil
}

// GetOtherStorages returns the entity manager(s) email of the entities
// storing the product with the id passed in the request parameters p.
func (db *SQLiteDataStore) GetOtherStorages(f zmqclient.RequestFilter, person_id int) ([]models.Entity, int, error) {
	var (
		entities                           []models.Entity
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	logger.Log.WithFields(logrus.Fields{"f": f}).Debug("GetOtherStorages")

	// hack to bypass optionnal default on the Rust part.
	if f.Search == "" {
		f.Search = "%%"
	}

	if f.OrderBy == "" {
		f.OrderBy = "storage_id"
	} else if f.OrderBy == "product.name.name_label" {
		f.OrderBy = "name.name_label"
	}

	// pre request: select or count
	precreq.WriteString(" SELECT count(DISTINCT e.entity_id)")
	presreq.WriteString(` SELECT e.entity_id AS "entity_id",
	e.entity_name AS "entity_name",
	GROUP_CONCAT(DISTINCT person.person_email) AS "entity_description"
	`)

	// common parts
	comreq.WriteString(" FROM entity as e")

	// get store location
	comreq.WriteString(" JOIN storelocation ON storelocation.entity = e.entity_id")
	// get storages
	comreq.WriteString(" JOIN storage ON storage.storelocation = storelocation.storelocation_id")

	// get managers
	comreq.WriteString(" JOIN entitypeople ON e.entity_id = entitypeople.entitypeople_entity_id")
	comreq.WriteString(" JOIN person ON entitypeople.entitypeople_person_id = person.person_id")

	comreq.WriteString(" WHERE 1")
	if f.Product != 0 {
		comreq.WriteString(" AND storage.product = :product")
	}

	// post select request
	postsreq.WriteString(" GROUP BY e.entity_id")

	// building count and select statements
	if cnstmt, err = db.PrepareNamed(precreq.String() + comreq.String()); err != nil {
		return nil, 0, err
	}
	if snstmt, err = db.PrepareNamed(presreq.String() + comreq.String() + postsreq.String()); err != nil {
		return nil, 0, err
	}

	// building argument map
	m := map[string]interface{}{
		"search":              f.Search,
		"personid":            person_id,
		"order":               f.Order,
		"limit":               f.Limit,
		"offset":              f.Offset,
		"entity":              f.Entity,
		"product":             f.Product,
		"storelocation":       f.Storelocation,
		"storage":             f.Storage,
		"name":                f.Name,
		"casnumber":           f.CasNumber,
		"empiricalformula":    f.EmpiricalFormula,
		"storage_barecode":    f.StorageBarecode,
		"custom_name_part_of": "%" + f.CustomNamePartOf + "%",
		"signalword":          f.SignalWord,
	}

	// Select.
	if err = snstmt.Select(&entities, m); err != nil {
		return nil, 0, err
	}
	// Count.
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	return entities, count, nil
}

// GetStorage returns the storage with id "id".
func (db *SQLiteDataStore) GetStorage(id int) (models.Storage, error) {
	var (
		storage models.Storage
		sqlr    string
		err     error
	)

	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetStorage")

	sqlr = `SELECT storage.storage_id,
	storage.storage_entrydate,
	storage.storage_exitdate,
	storage.storage_openingdate,
	storage.storage_expirationdate,
	storage.storage_reference,
	storage.storage_batchnumber,
	storage.storage_todestroy,
	storage.storage_creationdate,
	storage.storage_modificationdate,
	storage.storage_quantity,
	storage.storage_barecode,
	storage.storage_qrcode,
	storage.storage_comment,
	storage.storage_archive,
	storage.storage_number_of_carton,
	storage.storage_number_of_bag,
	storage.storage_number_of_unit,
	uq.unit_id AS "unit_quantity.unit_id",
	uq.unit_label AS "unit_quantity.unit_label",
	uc.unit_id AS "unit_concentration.unit_id",
	uc.unit_label AS "unit_concentration.unit_label",
	supplier.supplier_id AS "supplier.supplier_id",
	supplier.supplier_label AS "supplier.supplier_label",
	person.person_id AS "person.person_id",
	person.person_email AS "person.person_email",
	name.name_id AS "product.name.name_id",
	name.name_label AS "product.name.name_label",
	product.product_id AS "product.product_id",
	product.product_number_per_carton AS "product.product_number_per_carton",
	producerref.producerref_id AS "product.producerref.producerref_id",
	casnumber.casnumber_id AS "product.casnumber.casnumber_id",
	casnumber.casnumber_label AS "product.casnumber.casnumber_label",
	storelocation.storelocation_id AS "storelocation.storelocation_id",
	storelocation.storelocation_name AS "storelocation.storelocation_name",
	storelocation.storelocation_color AS "storelocation.storelocation_color",
	storelocation.storelocation_fullpath AS "storelocation.storelocation_fullpath",
	entity.entity_id AS "storelocation.entity.entity_id"
	FROM storage
	JOIN storelocation ON storage.storelocation = storelocation.storelocation_id
	JOIN entity ON storelocation.entity = entity.entity_id
	LEFT JOIN unit uq ON storage.unit_quantity = uq.unit_id
	LEFT JOIN unit uc ON storage.unit_concentration = uc.unit_id
	LEFT JOIN supplier ON storage.supplier = supplier.supplier_id
	JOIN person ON storage.person = person.person_id
	JOIN product ON storage.product = product.product_id
	LEFT JOIN producerref ON product.producerref = producerref.producerref_id
	LEFT JOIN casnumber ON product.casnumber = casnumber.casnumber_id
	JOIN name ON product.name = name.name_id
	WHERE storage.storage_id = ?`
	if err = db.Get(&storage, sqlr, id); err != nil {
		return models.Storage{}, err
	}

	if storage.Product.ProductNumberPerCarton.Valid {
		storage.Product.ProductType = "CONS"
	} else if storage.Product.ProducerRef.ProducerRefID.Valid {
		storage.Product.ProductType = "BIO"
	} else {
		storage.Product.ProductType = "CHEM"
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "storage": storage}).Debug("GetStorage")

	return storage, nil
}

// GetStorageEntity returns the entity of the storage with id "id".
func (db *SQLiteDataStore) GetStorageEntity(id int) (models.Entity, error) {
	var (
		entity models.Entity
		sqlr   string
		err    error
	)

	sqlr = `SELECT 
	entity.entity_id AS "entity_id",
	entity.entity_name AS "entity_name"
	FROM storage
	JOIN storelocation ON storage.storelocation = storelocation.storelocation_id
	JOIN entity ON storelocation.entity = entity.entity_id
	WHERE storage.storage_id = ?`
	if err = db.Get(&entity, sqlr, id); err != nil {
		return models.Entity{}, err
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "entity": entity}).Debug("GetStorageEntity")

	return entity, nil
}

// DeleteStorage deletes the storages with the given id.
func (db *SQLiteDataStore) DeleteStorage(id int) error {
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("DeleteStorage")

	var (
		sqlr string
		err  error
	)

	// Delete history first.
	sqlr = `DELETE FROM storage 
	WHERE storage = ?`
	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	sqlr = `DELETE FROM storage 
	WHERE storage_id = ?`
	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	return nil
}

// ArchiveStorage archives the storages with the given id.
func (db *SQLiteDataStore) ArchiveStorage(id int) error {
	var (
		sqlr string
		err  error
	)

	sqlr = `UPDATE storage SET storage_archive = true 
	WHERE storage_id = ?`

	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	sqlr = `UPDATE storage SET storage_archive = true 
	WHERE storage.storage = ?`

	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	return nil
}

// RestoreStorage restores (unarchive) the storages with the given id.
func (db *SQLiteDataStore) RestoreStorage(id int) error {
	var (
		sqlr string
		err  error
	)

	sqlr = `UPDATE storage SET storage_archive = false 
	WHERE storage_id = ?`

	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	sqlr = `UPDATE storage SET storage_archive = false 
	WHERE storage.storage = ?`

	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	return nil
}

// CreateStorage creates a new storage.
func (db *SQLiteDataStore) CreateUpdateStorage(s models.Storage, itemNumber int, update bool) (lastInsertID int64, err error) {
	var (
		v            driver.Value
		tx           *sql.Tx
		sqlr         string
		res          sql.Result
		args         []interface{}
		prefix       string
		major, minor string
	)

	// Default major.
	major = strconv.Itoa(s.ProductID)

	dialect := goqu.Dialect("sqlite3")
	tableStorage := goqu.T("storage")

	if tx, err = db.Begin(); err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			logger.Log.Error(err)
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.Log.Error(rbErr)
				err = rbErr

				return
			}

			return
		}

		err = tx.Commit()
	}()

	if update {
		// create an history of the storage
		sqlr = `INSERT into storage (storage_creationdate, 
		storage_modificationdate,
		storage_entrydate, 
		storage_exitdate, 
		storage_openingdate, 
		storage_expirationdate,
		storage_comment,
		storage_reference,
		storage_batchnumber,
		storage_quantity,
		storage_barecode,
		storage_todestroy,
		storage_archive,
		storage_concentration,
		storage_number_of_unit,
		storage_number_of_bag,
		storage_number_of_carton,
		person,
		product,
		storelocation,
		unit_quantity,
		unit_concentration,
		supplier,
		storage) select storage_creationdate, 
				storage_modificationdate,
				storage_entrydate, 
				storage_exitdate, 
				storage_openingdate, 
				storage_expirationdate,
				storage_comment,
				storage_reference,
				storage_batchnumber,
				storage_quantity,
				storage_barecode,
				storage_todestroy,
				storage_archive,
				storage_concentration,
				storage_number_of_unit,
				storage_number_of_bag,
				storage_number_of_carton,
				person,
				product,
				storelocation,
				unit_quantity,
				unit_concentration,
				supplier,
				? FROM storage WHERE storage_id = ?`
		if _, err = tx.Exec(sqlr, s.StorageID, s.StorageID); err != nil {
			return
		}
	}

	// Generating barecode if empty.
	if !update {
		if !(s.StorageBarecode.Valid) || s.StorageBarecode.String == "" {
			//
			// Getting the barecode prefix from the storelocation name.
			//
			// regex to detect store locations names starting with [_a-zA-Z] to build barecode prefixes
			prefixRegex := regexp.MustCompile(`^\[(?P<groupone>[_a-zA-Z]{1,5})\].*$`)
			groupNames := prefixRegex.SubexpNames()
			matches := prefixRegex.FindAllStringSubmatch(s.StoreLocationName.String, -1)
			// Building a map of matches.
			matchesMap := map[string]string{}

			if len(matches) != 0 {
				for i, j := range matches[0] {
					matchesMap[groupNames[i]] = j
				}
			}

			if len(matchesMap) > 0 {
				prefix = matchesMap["groupone"]
			} else {
				prefix = "_"
			}

			//
			// Getting the storage barecodes matching the regex
			// for the same product in the same entity.
			//
			sqlr := `SELECT storage_barecode FROM storage 
		JOIN storelocation on storage.storelocation = storelocation.storelocation_id 
		WHERE product = ? AND storelocation.entity = ? AND regexp('^[_a-zA-Z]{0,5}[0-9]+\.[0-9]+$', '' || storage_barecode || '') = true
		ORDER BY storage_barecode desc`

			var rows *sql.Rows

			if rows, err = tx.Query(sqlr, s.ProductID, s.EntityID); err != nil && err != sql.ErrNoRows {
				return
			}

			var (
				count    = 0
				newMinor = 0
			)

			for rows.Next() {
				var barecode string
				if err = rows.Scan(&barecode); err != nil && err != sql.ErrNoRows {
					return
				}

				majorRegex := regexp.MustCompile(`^[_a-zA-Z]{0,5}(?P<groupone>[0-9]+)\.(?P<grouptwo>[0-9]+)$`)
				groupNames = majorRegex.SubexpNames()
				matches = majorRegex.FindAllStringSubmatch(barecode, -1)
				// Building a map of matches.
				matchesMap = map[string]string{}

				if len(matches) != 0 {
					for i, j := range matches[0] {
						matchesMap[groupNames[i]] = j
					}
				}

				if count == 0 {
					// All of the major number are the same.
					// Extracting it ones.
					major = matchesMap["groupone"]
				}

				minor = matchesMap["grouptwo"]

				var iminor int

				if iminor, err = strconv.Atoi(minor); err != nil {
					return 0, err
				}

				if iminor > newMinor {
					newMinor = iminor
				}

				count++
			}

			if (!s.StorageIdenticalBarecode.Valid || !s.StorageIdenticalBarecode.Bool) || (s.StorageIdenticalBarecode.Valid && s.StorageIdenticalBarecode.Bool && itemNumber == 1) {
				newMinor++
			}

			minor = strconv.Itoa(newMinor)

			s.StorageBarecode.String = prefix + major + "." + minor
			s.StorageBarecode.Valid = true

			logger.Log.WithFields(logrus.Fields{"s.StorageBarecode.String": s.StorageBarecode.String}).Debug("CreateStorage")
		}
	}

	// if SupplierID = -1 then it is a new supplier
	if v, err = s.Supplier.SupplierID.Value(); s.Supplier.SupplierID.Valid && err == nil && v.(int64) == -1 {
		sqlr = `INSERT INTO supplier (supplier_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, s.Supplier.SupplierLabel); err != nil {
			return
		}
		// getting the last inserted id
		if lastInsertID, err = res.LastInsertId(); err != nil {
			return
		}
		// updating the storage SupplierId (SupplierLabel already set)
		s.Supplier.SupplierID = sql.NullInt64{Valid: true, Int64: lastInsertID}
	}
	if err != nil {
		logger.Log.Error("supplier error - " + err.Error())
		return
	}

	// finally updating the storage
	insertCols := goqu.Record{}
	if s.StorageComment.Valid {
		insertCols["storage_comment"] = s.StorageComment.String
	} else {
		insertCols["storage_comment"] = nil
	}

	if s.StorageQuantity.Valid {
		insertCols["storage_quantity"] = s.StorageQuantity.Float64
	} else {
		insertCols["storage_quantity"] = nil
	}

	if s.StorageBarecode.Valid {
		insertCols["storage_barecode"] = s.StorageBarecode.String
	} else {
		insertCols["storage_barecode"] = nil
	}

	if s.UnitQuantity.UnitID.Valid {
		insertCols["unit_quantity"] = s.UnitQuantity.UnitID.Int64
	} else {
		insertCols["unit_quantity"] = nil
	}

	if s.SupplierID.Valid {
		insertCols["supplier"] = s.SupplierID.Int64
	} else {
		insertCols["supplier"] = nil
	}

	if s.StorageEntryDate.Valid {
		insertCols["storage_entrydate"] = s.StorageEntryDate.Time
	} else {
		insertCols["storage_entrydate"] = nil
	}

	if s.StorageExitDate.Valid {
		insertCols["storage_exitdate"] = s.StorageExitDate.Time
	} else {
		insertCols["storage_exitdate"] = nil
	}

	if s.StorageOpeningDate.Valid {
		insertCols["storage_openingdate"] = s.StorageOpeningDate.Time
	} else {
		insertCols["storage_openingdate"] = nil
	}

	if s.StorageExpirationDate.Valid {
		insertCols["storage_expirationdate"] = s.StorageExpirationDate.Time
	} else {
		insertCols["storage_expirationdate"] = nil
	}

	if s.StorageReference.Valid {
		insertCols["storage_reference"] = s.StorageReference.String
	} else {
		insertCols["storage_reference"] = nil
	}

	if s.StorageBatchNumber.Valid {
		insertCols["storage_batchnumber"] = s.StorageBatchNumber.String
	} else {
		insertCols["storage_batchnumber"] = nil
	}

	if s.StorageToDestroy.Valid {
		insertCols["storage_todestroy"] = s.StorageToDestroy.Bool
	} else {
		insertCols["storage_todestroy"] = nil
	}

	if s.StorageConcentration.Valid {
		insertCols["storage_concentration"] = int(s.StorageConcentration.Int64)
	} else {
		insertCols["storage_concentration"] = nil
	}

	if s.StorageNumberOfBag.Valid {
		insertCols["storage_number_of_bag"] = int(s.StorageNumberOfBag.Int64)
	} else {
		insertCols["storage_number_of_bag"] = nil
	}

	if s.StorageNumberOfCarton.Valid {
		insertCols["storage_number_of_carton"] = int(s.StorageNumberOfCarton.Int64)
	} else {
		insertCols["storage_number_of_carton"] = nil
	}

	if s.StorageNumberOfUnit.Valid {
		insertCols["storage_number_of_unit"] = int(s.StorageNumberOfUnit.Int64)
	} else {
		insertCols["storage_number_of_unit"] = nil
	}

	if s.UnitConcentration.UnitID.Valid {
		insertCols["unit_concentration"] = int(s.UnitConcentration.UnitID.Int64)
	} else {
		insertCols["unit_concentration"] = nil
	}

	insertCols["person"] = s.PersonID
	insertCols["storelocation"] = s.StoreLocationID.Int64
	insertCols["product"] = s.ProductID
	insertCols["storage_creationdate"] = s.StorageCreationDate
	insertCols["storage_modificationdate"] = s.StorageModificationDate
	insertCols["storage_archive"] = false

	if update {
		iQuery := dialect.Update(tableStorage).Set(insertCols).Where(goqu.I("storage_id").Eq(s.StorageID))
		if sqlr, args, err = iQuery.ToSQL(); err != nil {
			return
		}
	} else {
		iQuery := dialect.Insert(tableStorage).Rows(insertCols)
		if sqlr, args, err = iQuery.ToSQL(); err != nil {
			return
		}
	}

	// logger.Log.Debug(sqlr)
	// logger.Log.Debug(args)

	if res, err = tx.Exec(sqlr, args...); err != nil {
		return
	}

	// getting the last inserted id
	if !update {
		if lastInsertID, err = res.LastInsertId(); err != nil {
			return
		}
	}

	//
	// qrcode
	//
	qr := strconv.FormatInt(lastInsertID, 10)
	if s.StorageQRCode, err = qrcode.Encode(qr, qrcode.Medium, 512); err != nil {
		return
	}

	sqlr = `UPDATE storage SET storage_qrcode=? WHERE storage_id=?`
	if _, err = tx.Exec(sqlr, s.StorageQRCode, lastInsertID); err != nil {
		return
	}

	s.StorageID = sql.NullInt64{Valid: true, Int64: lastInsertID}

	logger.Log.WithFields(logrus.Fields{"s": s}).Debug("CreateUpdateStorage")

	return
}

// UpdateAllQRCodes updates the storages QRCodes.
func (db *SQLiteDataStore) UpdateAllQRCodes() error {
	var (
		err  error
		tx   *sqlx.Tx
		sts  []models.Storage
		png  []byte
		sqlr string
	)

	// retrieving storages
	if err = db.Select(&sts, ` SELECT storage_id
        FROM storage`); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}

	// beginning new transaction
	if tx, err = db.Beginx(); err != nil {
		return err
	}

	for _, s := range sts {
		// generating qrcode
		newqrcode := strconv.FormatInt(s.StorageID.Int64, 10)
		logger.Log.Debug("  " + strconv.FormatInt(s.StorageID.Int64, 10) + " " + newqrcode)

		if png, err = qrcode.Encode(newqrcode, qrcode.Medium, 512); err != nil {
			return err
		}

		sqlr = `UPDATE storage
				SET storage_qrcode = ?
				WHERE storage_id = ?`

		if _, err = tx.Exec(sqlr, png, s.StorageID); err != nil {
			logger.Log.Error("error updating storage qrcode")
			if errr := tx.Rollback(); errr != nil {
				return errr
			}

			return err
		}
	}

	// committing changes
	if err = tx.Commit(); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}

	return nil
}

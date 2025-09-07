package datastores

import (
	"database/sql"
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
		s.storage_entry_date,
		s.storage_exit_date,
		s.storage_opening_date,
		s.storage_expiration_date,
		s.storage_reference,
		s.storage_batch_number,
		s.storage_to_destroy,
		s.storage_creation_date,
		s.storage_modification_date,
		s.storage_quantity,
		s.storage_barecode,
		s.storage_qrcode,
		s.storage_comment,
		s.storage_archive,
		s.storage_concentration,
		s.storage_number_of_carton,
		s.storage_number_of_bag,
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
        producer_ref.producer_ref_id AS "product.producer_ref.producer_ref_id",
		producer_ref.producer_ref_label AS "product.producer_ref.producer_ref_label",
		name.name_id AS "product.name.name_id",
		name.name_label AS "product.name.name_label",
		cas_number.cas_number_id AS "product.cas_number.cas_number_id",
		cas_number.cas_number_label AS "product.cas_number.cas_number_label",
		borrowing.borrowing_id AS "borrowing.borrowing_id",
		borrowing.borrowing_comment AS "borrowing.borrowing_comment",
		store_location.store_location_id AS "store_location.store_location_id",
		store_location.store_location_name AS "store_location.store_location_name",
		store_location.store_location_color AS "store_location.store_location_color",
		store_location.store_location_full_path AS "store_location.store_location_full_path",
		entity.entity_id AS "store_location.entity.entity_id"
		`)

	if f.CasNumberCmr {
		presreq.WriteString(`,GROUP_CONCAT(DISTINCT hazard_statement.hazard_statement_cmr) AS "product.hazard_statement_cmr"`)
	}

	// common parts
	comreq.WriteString(" FROM storage as s")
	// get storage history parent
	comreq.WriteString(" LEFT JOIN storage ON s.storage = storage.storage_id")
	// get product
	comreq.WriteString(" JOIN product ON s.product = product.product_id")
	// CMR
	if f.CasNumberCmr {
		comreq.WriteString(" LEFT JOIN producthazard_statements ON producthazard_statements.producthazard_statements_product_id = product.product_id")
		comreq.WriteString(" LEFT JOIN hazard_statement ON producthazard_statements.producthazard_statements_hazard_statement_id = hazard_statement.hazard_statement_id")
	}
	// get producer_ref
	if f.ProducerRef != 0 {
		comreq.WriteString(" JOIN producer_ref ON product.producer_ref = :producer_ref")
	} else {
		comreq.WriteString(" LEFT JOIN producer_ref ON product.producer_ref = producer_ref.producer_ref_id")
	}
	// get name
	comreq.WriteString(" JOIN name ON product.name = name.name_id")
	// get category
	if f.Category != 0 {
		comreq.WriteString(" JOIN category ON product.category = :category")
	}
	// get signal word
	comreq.WriteString(" LEFT JOIN signal_word ON product.signal_word = signal_word.signal_word_id")
	// get person
	comreq.WriteString(" JOIN person ON s.person = person.person_id")
	// get store location
	comreq.WriteString(" JOIN store_location ON s.store_location = store_location.store_location_id")
	// get entity
	comreq.WriteString(" JOIN entity ON store_location.entity = entity.entity_id")
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
	// get cas_number
	comreq.WriteString(" LEFT JOIN cas_number ON product.cas_number = cas_number.cas_number_id")
	// get empirical formula
	comreq.WriteString(" LEFT JOIN empirical_formula ON product.empirical_formula = empirical_formula.empirical_formula_id")
	// get symbols
	if len(f.Symbols) != 0 {
		comreq.WriteString(" JOIN productsymbols AS ps ON ps.productsymbols_product_id = product.product_id")
	}
	// get hazard_statements
	if len(f.HazardStatements) != 0 {
		comreq.WriteString(" JOIN producthazard_statements AS phs ON phs.producthazard_statements_product_id = product.product_id")
	}
	// get precautionary_statements
	if len(f.PrecautionaryStatements) != 0 {
		comreq.WriteString(" JOIN productprecautionary_statements AS pps ON pps.productprecautionary_statements_product_id = product.product_id")
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
		comreq.WriteString(` JOIN personentities ON (personentities_entity_id = store_location.entity AND personentities_person_id = :personid)`)
	}

	// filter by permissions
	comreq.WriteString(` JOIN permission AS perm ON
	perm.person = :personid and (perm.permission_item in ("all", "storages")) and (perm.permission_name in ("all", "r", "w")) and (perm.permission_entity in (-1, entity.entity_id))
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
		comreq.WriteString(" AND s.storage_to_destroy = true")
	}
	if f.CasNumberCmr {
		comreq.WriteString(" AND (cas_number.cas_number_cmr IS NOT NULL OR (hazard_statement_cmr IS NOT NULL AND hazard_statement_cmr != ''))")
	}
	if f.Product != 0 {
		comreq.WriteString(" AND product.product_id = :product")
	}
	if f.Entity != 0 {
		comreq.WriteString(" AND entity.entity_id = :entity")
	}
	if f.Storelocation != 0 {
		comreq.WriteString(" AND store_location.store_location_id = :store_location")
	}
	if f.Storage != 0 {
		if f.History {
			comreq.WriteString(" AND (s.storage = :storage OR s.storage_id = :storage)")
		} else {
			comreq.WriteString(" AND (s.storage_id = :storage")
			// getting storages with identical barecode
			comreq.WriteString(" OR (s.storage_barecode = (SELECT storage_barecode FROM storage WHERE storage_id = :storage)))")
		}
	} else if f.Id != 0 {
		if f.History {
			comreq.WriteString(" AND (s.storage = :id OR s.storage_id = :id)")
		} else {
			comreq.WriteString(" AND (s.storage_id = :id")
			// getting storages with identical barecode
			comreq.WriteString(" OR (s.storage_barecode = (SELECT storage_barecode FROM storage WHERE storage_id = :id)))")
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
		comreq.WriteString(" AND cas_number.cas_number_id = :cas_number")
	}
	if f.EmpiricalFormula != 0 {
		comreq.WriteString(" AND empirical_formula.empirical_formula_id = :empirical_formula")
	}
	if f.StorageBarecode != "" {
		comreq.WriteString(" AND s.storage_barecode LIKE :storage_barecode")
	}
	if f.StorageBatchNumber != "" {
		comreq.WriteString(" AND s.storage_batch_number LIKE :storage_batch_number")
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
		comreq.WriteString(" AND phs.producthazard_statements_hazard_statement_id IN (")

		for _, s := range f.HazardStatements {
			comreq.WriteString(fmt.Sprintf("%d,", s))
		}

		// to complete the last comma
		comreq.WriteString("-1")
		comreq.WriteString(" )")
	}
	if len(f.PrecautionaryStatements) != 0 {
		comreq.WriteString(" AND pps.productprecautionary_statements_precautionary_statement_id IN (")

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
		comreq.WriteString(" AND signal_word.signal_word_id = :signal_word")
	}

	// show bio/chem/consu
	switch {
	case !f.ShowChem && !f.ShowBio && f.ShowConsu:
		comreq.WriteString(" AND (product_number_per_carton IS NOT NULL AND product_number_per_carton != 0)")
	case !f.ShowChem && f.ShowBio && !f.ShowConsu:
		comreq.WriteString(" AND producer_ref IS NOT NULL")
		comreq.WriteString(" AND (product_number_per_carton IS NULL OR product_number_per_carton == 0)")
	case !f.ShowChem && f.ShowBio && f.ShowConsu:
		comreq.WriteString(" AND ((product_number_per_carton IS NOT NULL AND product_number_per_carton != 0)")
		comreq.WriteString(" OR producer_ref IS NOT NULL)")
	case f.ShowChem && !f.ShowBio && !f.ShowConsu:
		comreq.WriteString(" AND producer_ref IS NULL")
		comreq.WriteString(" AND (product_number_per_carton IS NULL OR product_number_per_carton == 0)")
	case f.ShowChem && !f.ShowBio && f.ShowConsu:
		comreq.WriteString(" AND (producer_ref IS NULL")
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
		"id":                   f.Id,
		"ids":                  f.Ids,
		"search":               f.Search,
		"personid":             person_id,
		"order":                f.Order,
		"limit":                f.Limit,
		"offset":               f.Offset,
		"entity":               f.Entity,
		"product":              f.Product,
		"store_location":       f.Storelocation,
		"storage":              f.Storage,
		"name":                 f.Name,
		"cas_number":           f.CasNumber,
		"empirical_formula":    f.EmpiricalFormula,
		"storage_barecode":     "%" + f.StorageBarecode + "%",
		"storage_batch_number": f.StorageBatchNumber,
		"custom_name_part_of":  "%" + f.CustomNamePartOf + "%",
		"signal_word":          f.SignalWord,
		"producer_ref":         f.ProducerRef,
		"category":             f.Category,
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
		if st.Product.ProductNumberPerCarton != nil {
			storages[i].Product.ProductType = "cons"
		} else if st.Product.ProducerRef.ProducerRefID != nil {
			storages[i].Product.ProductType = "bio"
		} else {
			storages[i].Product.ProductType = "chem"
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
	comreq.WriteString(" JOIN store_location ON store_location.entity = e.entity_id")
	// get storages
	comreq.WriteString(" JOIN storage ON storage.store_location = store_location.store_location_id")

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
		"store_location":      f.Storelocation,
		"storage":             f.Storage,
		"name":                f.Name,
		"cas_number":          f.CasNumber,
		"empirical_formula":   f.EmpiricalFormula,
		"storage_barecode":    f.StorageBarecode,
		"custom_name_part_of": "%" + f.CustomNamePartOf + "%",
		"signal_word":         f.SignalWord,
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
		tx           *sql.Tx
		sqlr         string
		res          sql.Result
		args         []interface{}
		prefix       string
		major, minor string
	)

	logger.Log.WithFields(logrus.Fields{"s": fmt.Sprintf("%+v", s)}).Debug("CreateUpdateStorage")

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
		sqlr = `INSERT into storage (storage_creation_date,
		storage_modification_date,
		storage_entry_date,
		storage_exit_date,
		storage_opening_date,
		storage_expiration_date,
		storage_comment,
		storage_reference,
		storage_batch_number,
		storage_quantity,
		storage_barecode,
		storage_to_destroy,
		storage_archive,
		storage_concentration,
		storage_number_of_bag,
		storage_number_of_carton,
		person,
		product,
		store_location,
		unit_quantity,
		unit_concentration,
		supplier,
		storage) select storage_creation_date,
				storage_modification_date,
				storage_entry_date,
				storage_exit_date,
				storage_opening_date,
				storage_expiration_date,
				storage_comment,
				storage_reference,
				storage_batch_number,
				storage_quantity,
				storage_barecode,
				storage_to_destroy,
				storage_archive,
				storage_concentration,
				storage_number_of_bag,
				storage_number_of_carton,
				person,
				product,
				store_location,
				unit_quantity,
				unit_concentration,
				supplier,
				? FROM storage WHERE storage_id = ?`
		if _, err = tx.Exec(sqlr, s.StorageID, s.StorageID); err != nil {
			logger.Log.Error("error creating storage history")
			return
		}
	}

	// Generating barecode if empty.
	if !update {
		if s.StorageBarecode == nil || *s.StorageBarecode == "" {
			//
			// Getting the barecode prefix from the store_location name.
			//
			// regex to detect store locations names starting with [_a-zA-Z] to build barecode prefixes
			prefixRegex := regexp.MustCompile(`^\[(?P<groupone>[_a-zA-Z]{1,5})\].*$`)
			groupNames := prefixRegex.SubexpNames()
			matches := prefixRegex.FindAllStringSubmatch(s.StoreLocationName.String, -1)
			// Building a map of matches.
			matchesMap := map[string]string{}

			logger.Log.WithFields(logrus.Fields{"s.StoreLocationName.String": s.StoreLocationName.String, "matches": matches}).Debug("CreateStorage")

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
		JOIN store_location on storage.store_location = store_location.store_location_id
		WHERE product = ? AND store_location.entity = ? AND regexp('^[_a-zA-Z]{0,5}[0-9]+\.[0-9]+$', '' || storage_barecode || '') = true
		ORDER BY storage_barecode desc`

			var rows *sql.Rows

			if rows, err = tx.Query(sqlr, s.ProductID, s.EntityID); err != nil && err != sql.ErrNoRows {
				logger.Log.Error("error getting storage barecode")
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

			if !s.StorageIdenticalBarecode || (s.StorageIdenticalBarecode && itemNumber == 1) {
				newMinor++
			}

			minor = strconv.Itoa(newMinor)

			*s.StorageBarecode = prefix + major + "." + minor

			logger.Log.WithFields(logrus.Fields{"s.StorageBarecode.String": s.StorageBarecode}).Debug("CreateStorage")
		}
	}

	// if SupplierID = -1 then it is a new supplier
	if s.Supplier.SupplierID != nil && err == nil && *s.Supplier.SupplierID == -1 {
		sqlr = `INSERT INTO supplier (supplier_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, s.Supplier.SupplierLabel); err != nil {
			return
		}
		// getting the last inserted id
		if lastInsertID, err = res.LastInsertId(); err != nil {
			return
		}
		// updating the storage SupplierId (SupplierLabel already set)
		*s.Supplier.SupplierID = lastInsertID
	}
	if err != nil {
		logger.Log.Error("supplier error - " + err.Error())
		return
	}

	// finally updating the storage
	insertCols := goqu.Record{}
	if s.StorageComment != nil {
		insertCols["storage_comment"] = s.StorageComment
	} else {
		insertCols["storage_comment"] = nil
	}

	if s.StorageQuantity != nil {
		insertCols["storage_quantity"] = s.StorageQuantity
	} else {
		insertCols["storage_quantity"] = nil
	}

	if s.StorageBarecode != nil {
		insertCols["storage_barecode"] = s.StorageBarecode
	} else {
		insertCols["storage_barecode"] = nil
	}

	if s.UnitQuantity.UnitID != nil {
		insertCols["unit_quantity"] = *s.UnitQuantity.UnitID
	} else {
		insertCols["unit_quantity"] = nil
	}

	if s.Supplier.SupplierID != nil {
		insertCols["supplier"] = *s.SupplierID
	} else {
		insertCols["supplier"] = nil
	}

	if s.StorageEntryDate != nil {
		insertCols["storage_entry_date"] = s.StorageEntryDate.Unix()
	} else {
		insertCols["storage_entry_date"] = nil
	}

	if s.StorageExitDate != nil {
		insertCols["storage_exit_date"] = s.StorageExitDate.Unix()
	} else {
		insertCols["storage_exit_date"] = nil
	}

	if s.StorageOpeningDate != nil {
		insertCols["storage_opening_date"] = s.StorageOpeningDate.Unix()
	} else {
		insertCols["storage_opening_date"] = nil
	}

	if s.StorageExpirationDate != nil {
		insertCols["storage_expiration_date"] = s.StorageExpirationDate.Unix()
	} else {
		insertCols["storage_expiration_date"] = nil
	}

	if s.StorageReference != nil {
		insertCols["storage_reference"] = s.StorageReference
	} else {
		insertCols["storage_reference"] = nil
	}

	if s.StorageBatchNumber != nil {
		insertCols["storage_batch_number"] = s.StorageBatchNumber
	} else {
		insertCols["storage_batch_number"] = nil
	}

	if s.StorageToDestroy {
		insertCols["storage_to_destroy"] = s.StorageToDestroy
	} else {
		insertCols["storage_to_destroy"] = false
	}

	if s.StorageConcentration != nil {
		insertCols["storage_concentration"] = int(*s.StorageConcentration)
	} else {
		insertCols["storage_concentration"] = nil
	}

	if s.StorageNumberOfBag != nil {
		insertCols["storage_number_of_bag"] = int(*s.StorageNumberOfBag)
	} else {
		insertCols["storage_number_of_bag"] = nil
	}

	if s.StorageNumberOfCarton != nil {
		insertCols["storage_number_of_carton"] = int(*s.StorageNumberOfCarton)
	} else {
		insertCols["storage_number_of_carton"] = nil
	}

	if s.UnitConcentration.UnitID != nil {
		insertCols["unit_concentration"] = int(*s.UnitConcentration.UnitID)
	} else {
		insertCols["unit_concentration"] = nil
	}

	insertCols["person"] = s.PersonID
	insertCols["store_location"] = s.StoreLocationID.Int64
	insertCols["product"] = s.ProductID
	insertCols["storage_creation_date"] = s.StorageCreationDate.Unix()
	insertCols["storage_modification_date"] = s.StorageModificationDate.Unix()
	insertCols["storage_archive"] = false

	if update {
		iQuery := dialect.Update(tableStorage).Set(insertCols).Where(goqu.I("storage_id").Eq(s.StorageID))
		if sqlr, args, err = iQuery.ToSQL(); err != nil {
			logger.Log.Error("error preparing update storage")
			return
		}
	} else {
		iQuery := dialect.Insert(tableStorage).Rows(insertCols)
		if sqlr, args, err = iQuery.ToSQL(); err != nil {
			logger.Log.Error("error preparing create storage")
			return
		}
	}

	// logger.Log.Debug(sqlr)
	// logger.Log.Debug(args)

	if res, err = tx.Exec(sqlr, args...); err != nil {
		logger.Log.Error("error creating/updating storage")
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

	var storage_id int64 = lastInsertID
	s.StorageID = &storage_id

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

	// beginning new transaction
	if tx, err = db.Beginx(); err != nil {
		return err
	}

	// retrieving storages
	if err = db.Select(&sts, ` SELECT storage_id
        FROM storage`); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}

	for _, s := range sts {
		// generating qrcode
		newqrcode := strconv.FormatInt(*s.StorageID, 10)
		logger.Log.Debug("  " + strconv.FormatInt(*s.StorageID, 10) + " " + newqrcode)

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

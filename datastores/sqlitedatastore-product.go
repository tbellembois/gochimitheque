package datastores

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/zmqclient"
)

// IsProductBookmark returns true if there is a bookmark for the product pr for the person pe.
func (db *SQLiteDataStore) IsProductBookmark(pr models.Product, pe models.Person) (bool, error) {
	var (
		err   error
		sqlr  string
		args  []interface{}
		count int
	)

	dialect := goqu.Dialect("sqlite3")
	tableProduct := goqu.T("bookmark")

	sQuery := dialect.From(tableProduct).Select(
		goqu.COUNT("*"),
	).Where(
		goqu.I("person").Eq(pe.PersonID),
		goqu.I("product").Eq(pr.ProductID),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return false, err
	}

	if err = db.Get(&count, sqlr, args...); err != nil {
		return false, err
	}

	return count != 0, nil
}

// CreateProductBookmark bookmarks the product pr for the person pe.
func (db *SQLiteDataStore) CreateProductBookmark(pr models.Product, pe models.Person) (err error) {
	var tx *sqlx.Tx

	logger.Log.WithFields(logrus.Fields{"s": fmt.Sprintf("pr:%+v pe:%+v", pr, pe)}).Debug("CreateProductBookmark")

	dialect := goqu.Dialect("sqlite3")
	tableBookmark := goqu.T("bookmark")

	if tx, err = db.Beginx(); err != nil {
		return err
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

	iQuery := dialect.Insert(tableBookmark)

	setClause := goqu.Record{
		"person":  pe.PersonID,
		"product": pr.ProductID,
	}

	var (
		sqlr      string
		args      []interface{}
		sqlResult sql.Result
	)

	if sqlr, args, err = iQuery.Rows(setClause).ToSQL(); err != nil {
		return
	}

	if sqlResult, err = tx.Exec(sqlr, args...); err != nil {
		return
	}

	_, err = sqlResult.LastInsertId()

	return
}

// DeleteProductBookmark remove the bookmark for the product pr and the person pe.
func (db *SQLiteDataStore) DeleteProductBookmark(pr models.Product, pe models.Person) error {
	dialect := goqu.Dialect("sqlite3")
	tableBookmark := goqu.T("bookmark")

	dQuery := dialect.From(tableBookmark).Where(
		goqu.I("person").Eq(pe.PersonID),
		goqu.I("product").Eq(pr.ProductID),
	).Delete()

	var (
		err  error
		sqlr string
		args []interface{}
	)

	if sqlr, args, err = dQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return err
	}

	if _, err = db.Exec(sqlr, args...); err != nil {
		return err
	}

	return nil
}

// GetProducts return the products matching the search criteria.
func (db *SQLiteDataStore) GetProducts(f zmqclient.Filter, person_id int, public bool) ([]models.Product, int, error) {
	// defer TimeTrack(time.Now(), "GetProducts")
	var (
		products                           []models.Product
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
		rperm                              bool
		isadmin                            bool
		wg                                 sync.WaitGroup
	)

	logger.Log.WithFields(logrus.Fields{"f": fmt.Sprintf("%+v", f)}).Debug("GetProducts")

	if f.OrderBy == "" {
		f.OrderBy = "product_id"
	}

	if !public {
		// is the user an admin?
		if isadmin, err = db.IsPersonAdmin(person_id); err != nil {
			return nil, 0, err
		}
		// has the person rproducts permission?
		if rperm, err = db.HasPersonReadRestrictedProductPermission(person_id); err != nil {
			return nil, 0, err
		}
	} else {
		isadmin = true
		rperm = true
	}

	// pre request: select or count
	precreq.WriteString(" SELECT count(DISTINCT p.product_id)")
	presreq.WriteString(` SELECT p.product_id, 
	p.product_specificity, 
	p.product_msds,
	p.product_restricted,
	p.product_radioactive,
	p.product_threedformula,
	p.product_twodformula,
	p.product_molformula,
	p.product_disposalcomment,
	p.product_remark,
	p.product_sheet,
	p.product_temperature,
	p.product_number_per_carton,
	p.product_number_per_bag,
	linearformula.linearformula_id AS "linearformula.linearformula_id",
	linearformula.linearformula_label AS "linearformula.linearformula_label",
	empiricalformula.empiricalformula_id AS "empiricalformula.empiricalformula_id",
	empiricalformula.empiricalformula_label AS "empiricalformula.empiricalformula_label",
	physicalstate.physicalstate_id AS "physicalstate.physicalstate_id",
	physicalstate.physicalstate_label AS "physicalstate.physicalstate_label",
	signalword.signalword_id AS "signalword.signalword_id",
	signalword.signalword_label AS "signalword.signalword_label",
	person.person_id AS "person.person_id",
	person.person_email AS "person.person_email",
	name.name_id AS "name.name_id",
	name.name_label AS "name.name_label",
	bookmark.bookmark_id AS "bookmark.bookmark_id",
	cenumber.cenumber_id AS "cenumber.cenumber_id",
	cenumber.cenumber_label AS "cenumber.cenumber_label",
	casnumber.casnumber_id AS "casnumber.casnumber_id",
	casnumber.casnumber_label AS "casnumber.casnumber_label",
	casnumber.casnumber_cmr AS "casnumber.casnumber_cmr",
	producerref.producerref_id AS "producerref.producerref_id",
	producerref.producerref_label AS "producerref.producerref_label",
	producer.producer_id AS "producerref.producer.producer_id",
	producer.producer_label AS "producerref.producer.producer_label",
	ut.unit_id AS "unit_temperature.unit_id",
	ut.unit_label AS "unit_temperature.unit_label",
	category.category_id AS "category.category_id",
	category.category_label AS "category.category_label"
	`)

	if !public {
		presreq.WriteString(`,GROUP_CONCAT(DISTINCT storage.storage_barecode) AS "product_sl"`)
	}

	if f.CasNumberCmr {
		presreq.WriteString(`,GROUP_CONCAT(DISTINCT hazardstatement.hazardstatement_cmr) AS "hazardstatement_cmr"`)
	}

	// common parts
	comreq.WriteString(" FROM product as p")
	// CMR
	if f.CasNumberCmr {
		comreq.WriteString(" LEFT JOIN producthazardstatements ON producthazardstatements.producthazardstatements_product_id = p.product_id")
		comreq.WriteString(" LEFT JOIN hazardstatement ON producthazardstatements.producthazardstatements_hazardstatement_id = hazardstatement.hazardstatement_id")
	}
	// get name
	comreq.WriteString(" JOIN name ON p.name = name.name_id")
	// get category
	if f.Category != 0 {
		comreq.WriteString(" JOIN category ON p.category = :category")
	} else {
		comreq.WriteString(" LEFT JOIN category ON p.category = category.category_id")
	}
	// get unit_temperature
	comreq.WriteString(" LEFT JOIN unit ut ON p.unit_temperature = ut.unit_id")
	// get producerref
	if f.ProducerRef != 0 {
		comreq.WriteString(" JOIN producerref ON p.producerref = :producerref")
	} else {
		comreq.WriteString(" LEFT JOIN producerref ON p.producerref = producerref.producerref_id")
	}
	// get producer
	comreq.WriteString(" LEFT JOIN producer ON producerref.producer = producer.producer_id")
	// get casnumber
	comreq.WriteString(" LEFT JOIN casnumber ON p.casnumber = casnumber.casnumber_id")
	// get cenumber
	comreq.WriteString(" LEFT JOIN cenumber ON p.cenumber = cenumber.cenumber_id")
	// get person
	comreq.WriteString(" JOIN person ON p.person = person.person_id")
	// get physical state
	comreq.WriteString(" LEFT JOIN physicalstate ON p.physicalstate = physicalstate.physicalstate_id")
	// get signal word
	comreq.WriteString(" LEFT JOIN signalword ON p.signalword = signalword.signalword_id")
	// get empirical formula
	comreq.WriteString(" LEFT JOIN empiricalformula ON p.empiricalformula = empiricalformula.empiricalformula_id")
	// get linear formula
	comreq.WriteString(" LEFT JOIN linearformula ON p.linearformula = linearformula.linearformula_id")
	// get bookmark
	comreq.WriteString(" LEFT JOIN bookmark ON (bookmark.product = p.product_id AND bookmark.person = :personid)")
	// get storages, store locations and entities
	comreq.WriteString(" LEFT JOIN storage ON storage.product = p.product_id")

	if f.Entity != 0 || f.Storelocation != 0 || f.StorageBarecode != "" {
		comreq.WriteString(" JOIN storelocation ON storage.storelocation = storelocation.storelocation_id")
		comreq.WriteString(" JOIN entity ON storelocation.entity = entity.entity_id")
	}
	// get borrowings
	if f.Borrowing {
		comreq.WriteString(" JOIN borrowing ON borrowing.storage = storage.storage_id AND borrowing.borrower = :personid")
	}
	// get bookmarks
	if f.Bookmark {
		comreq.WriteString(" JOIN bookmark AS b ON b.product = p.product_id AND b.person = :personid")
	}
	// get symbols
	if len(f.Symbols) != 0 {
		comreq.WriteString(" JOIN productsymbols AS ps ON ps.productsymbols_product_id = p.product_id")
	}
	// get synonyms
	if f.Name != 0 {
		comreq.WriteString(" JOIN productsynonyms AS psyn ON psyn.productsynonyms_product_id = p.product_id")
	}
	// get hazardstatements
	if len(f.HazardStatements) != 0 {
		comreq.WriteString(" JOIN producthazardstatements AS phs ON phs.producthazardstatements_product_id = p.product_id")
	}
	// get precautionarystatements
	if len(f.PrecautionaryStatements) != 0 {
		comreq.WriteString(" JOIN productprecautionarystatements AS pps ON pps.productprecautionarystatements_product_id = p.product_id")
	}
	// get tags
	if len(f.Tags) != 0 {
		comreq.WriteString(" JOIN producttags AS ptags ON ptags.producttags_product_id = p.product_id")
	}

	// filter by permissions
	if !public {
		comreq.WriteString(` JOIN permission AS perm ON
	perm.person = :personid and 
	(perm.permission_item_name in ("all", "products")) and 
	(perm.permission_perm_name in ("all", "r", "w"))
	`)
	}

	comreq.WriteString(" WHERE 1")

	if f.StorageToDestroy {
		comreq.WriteString(" AND storage.storage_todestroy = true")
	}

	if f.CasNumberCmr {
		comreq.WriteString(" AND (casnumber.casnumber_cmr IS NOT NULL OR (hazardstatement_cmr IS NOT NULL AND hazardstatement_cmr != ''))")
	}

	if f.Product != 0 {
		comreq.WriteString(" AND p.product_id = :product")
	}

	if f.Entity != 0 {
		comreq.WriteString(" AND entity.entity_id = :entity")
	}

	if f.Storelocation != 0 {
		comreq.WriteString(" AND storelocation.storelocation_id = :storelocation")
	}

	if f.ProductSpecificity != "" {
		comreq.WriteString(" AND p.product_specificity = :product_specificity")
	}

	// search form parameters
	if f.Name != 0 {
		comreq.WriteString(" AND (name.name_id = :name")
		comreq.WriteString(" OR psyn.productsynonyms_name_id = :name)")
	}

	if f.CasNumber != 0 {
		comreq.WriteString(" AND casnumber.casnumber_id = :casnumber")
	}

	if f.EmpiricalFormula != 0 {
		comreq.WriteString(" AND empiricalformula.empiricalformula_id = :empiricalformula")
	}

	if f.StorageBarecode != "" {
		comreq.WriteString(" AND storage.storage_barecode LIKE :storage_barecode")
	}

	if f.StorageBatchNumber != "" {
		comreq.WriteString(" AND storage.storage_batchnumber LIKE :storage_batchnumber")
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

	// filter restricted product
	if !rperm {
		comreq.WriteString(" AND p.product_restricted = false")
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
	postsreq.WriteString(" GROUP BY p.product_id")
	postsreq.WriteString(" ORDER BY " + f.OrderBy + " " + f.Order)

	// limit
	if !public {
		if f.Limit != ^uint64(0) {
			postsreq.WriteString(" LIMIT :limit OFFSET :offset")
		}
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
		"search":              f.Search,
		"personid":            person_id,
		"order":               f.Order,
		"limit":               f.Limit,
		"offset":              f.Offset,
		"entity":              f.Entity,
		"product":             f.Product,
		"storelocation":       f.Storelocation,
		"name":                f.Name,
		"casnumber":           f.CasNumber,
		"empiricalformula":    f.EmpiricalFormula,
		"product_specificity": f.ProductSpecificity,
		"storage_barecode":    f.StorageBarecode,
		"storage_batchnumber": f.StorageBatchNumber,
		"custom_name_part_of": "%" + f.CustomNamePartOf + "%",
		"signalword":          f.SignalWord,
		"producerref":         f.ProducerRef,
		"category":            f.Category,
	}

	// Select.
	if err = snstmt.Select(&products, m); err != nil {
		return nil, 0, err
	}

	logger.Log.WithFields(logrus.Fields{"select req": presreq.String() + comreq.String() + postsreq.String()}).Debug("GetProducts")
	logger.Log.WithFields(logrus.Fields{"m": m}).Debug("GetProducts")

	// Count.
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	wg.Add(1)

	go func() {
		for i, pr := range products {
			switch {
			case pr.ProductNumberPerCarton.Valid:
				products[i].ProductType = "CONS"
			case pr.ProducerRef.ProducerRefID.Valid:
				products[i].ProductType = "BIO"
			default:
				products[i].ProductType = "CHEM"
			}
		}

		wg.Done()
	}()

	//
	// cleaning product_sl
	//
	if !public {
		wg.Add(1)

		go func() {
			r := regexp.MustCompile(`([a-zA-Z]{1}[0-9]+)\.[0-9]+`)
			for i, pr := range products {
				// note: do not modify p but products[i] instead
				m := r.FindAllStringSubmatch(pr.ProductSL.String, -1)

				if len(m) > 0 {
					differentSL := false
					mBackup := m[0][1]

					for i := range m {
						if (m[i][1]) != mBackup {
							differentSL = true
							break
						}
					}

					if !differentSL {
						products[i].ProductSL.String = m[0][1]
					} else {
						products[i].ProductSL.String = ""
					}
				} else {
					products[i].ProductSL.String = ""
				}
			}

			wg.Done()
		}()
	}

	//
	// getting supplierref
	//
	if !public {
		wg.Add(1)

		go func() {
			var reqSupplierref strings.Builder

			for i, pr := range products {
				// note: do not modify p but products[i] instead
				reqSupplierref.Reset()
				reqSupplierref.WriteString(`SELECT supplierref_id, 
			supplierref_label,
			supplier.supplier_id AS "supplier.supplier_id",
			supplier.supplier_label AS "supplier.supplier_label"
			FROM supplierref`)
				reqSupplierref.WriteString(" JOIN productsupplierrefs ON productsupplierrefs.productsupplierrefs_supplierref_id = supplierref.supplierref_id AND productsupplierrefs.productsupplierrefs_product_id = ?")
				reqSupplierref.WriteString(" JOIN supplier ON supplierref.supplier = supplier.supplier_id")

				if err = db.Select(&products[i].SupplierRefs, reqSupplierref.String(), pr.ProductID); err != nil {
					logger.Log.WithFields(logrus.Fields{"err": err}).Error("GetProducts:goroutine:supplierref")
				}
			}

			wg.Done()
		}()
	}

	//
	// getting tags
	//
	if !public {
		wg.Add(1)

		go func() {
			var reqTags strings.Builder

			for i, pr := range products {
				// note: do not modify p but products[i] instead
				reqTags.Reset()
				reqTags.WriteString("SELECT tag_id, tag_label FROM tag")
				reqTags.WriteString(" JOIN producttags ON producttags.producttags_tag_id = tag.tag_id")
				reqTags.WriteString(" JOIN product ON producttags.producttags_product_id = product.product_id")
				reqTags.WriteString(" WHERE product.product_id = ?")

				if err = db.Select(&products[i].Tags, reqTags.String(), pr.ProductID); err != nil {
					logger.Log.WithFields(logrus.Fields{"err": err}).Error("GetProducts:goroutine:tags")
				}
			}

			wg.Done()
		}()
	}

	//
	// getting symbols
	//
	wg.Add(1)

	go func() {
		var reqSymbols strings.Builder

		for i, pr := range products {
			// note: do not modify p but products[i] instead
			reqSymbols.Reset()
			reqSymbols.WriteString("SELECT symbol_id, symbol_label, symbol_image FROM symbol")
			reqSymbols.WriteString(" JOIN productsymbols ON productsymbols.productsymbols_symbol_id = symbol.symbol_id")
			reqSymbols.WriteString(" JOIN product ON productsymbols.productsymbols_product_id = product.product_id")
			reqSymbols.WriteString(" WHERE product.product_id = ?")

			if err = db.Select(&products[i].Symbols, reqSymbols.String(), pr.ProductID); err != nil {
				logger.Log.WithFields(logrus.Fields{"err": err}).Error("GetProducts:goroutine:symbols")
			}
		}

		wg.Done()
	}()

	//
	// getting classes of compounds
	//
	wg.Add(1)

	go func() {
		var reqCoc strings.Builder

		for i, pr := range products {
			// note: do not modify p but products[i] instead
			reqCoc.Reset()
			reqCoc.WriteString("SELECT classofcompound_id, classofcompound_label FROM classofcompound")
			reqCoc.WriteString(" JOIN productclassofcompound ON productclassofcompound.productclassofcompound_classofcompound_id = classofcompound.classofcompound_id")
			reqCoc.WriteString(" JOIN product ON productclassofcompound.productclassofcompound_product_id = product.product_id")
			reqCoc.WriteString(" WHERE product.product_id = ?")

			if err = db.Select(&products[i].ClassOfCompound, reqCoc.String(), pr.ProductID); err != nil {
				logger.Log.WithFields(logrus.Fields{"err": err}).Error("GetProducts:goroutine:coc")
			}
		}

		wg.Done()
	}()

	//
	// getting synonyms
	//
	wg.Add(1)

	go func() {
		var reqSynonyms strings.Builder

		for i, pr := range products {
			// note: do not modify p but products[i] instead
			reqSynonyms.Reset()
			reqSynonyms.WriteString("SELECT name_id, name_label FROM name")
			reqSynonyms.WriteString(" JOIN productsynonyms ON productsynonyms.productsynonyms_name_id = name.name_id")
			reqSynonyms.WriteString(" JOIN product ON productsynonyms.productsynonyms_product_id = product.product_id")
			reqSynonyms.WriteString(" WHERE product.product_id = ?")

			if err = db.Select(&products[i].Synonyms, reqSynonyms.String(), pr.ProductID); err != nil {
				logger.Log.WithFields(logrus.Fields{"err": err}).Error("GetProducts:goroutine:synonyms")
			}
		}

		wg.Done()
	}()

	//
	// getting hazard statements
	//
	wg.Add(1)

	go func() {
		var reqHS strings.Builder

		for i, pr := range products {
			// note: do not modify p but products[i] instead
			reqHS.Reset()
			reqHS.WriteString("SELECT hazardstatement_id, hazardstatement_label, hazardstatement_reference, hazardstatement_cmr FROM hazardstatement")
			reqHS.WriteString(" JOIN producthazardstatements ON producthazardstatements.producthazardstatements_hazardstatement_id = hazardstatement.hazardstatement_id")
			reqHS.WriteString(" JOIN product ON producthazardstatements.producthazardstatements_product_id = product.product_id")
			reqHS.WriteString(" WHERE product.product_id = ?")

			if err = db.Select(&products[i].HazardStatements, reqHS.String(), pr.ProductID); err != nil {
				logger.Log.WithFields(logrus.Fields{"err": err}).Error("GetProducts:goroutine:hs")
			}
		}

		wg.Done()
	}()

	//
	// getting precautionary statements
	//
	wg.Add(1)

	go func() {
		var reqPS strings.Builder

		for i, pr := range products {
			// note: do not modify p but products[i] instead
			reqPS.Reset()
			reqPS.WriteString("SELECT precautionarystatement_id, precautionarystatement_label, precautionarystatement_reference FROM precautionarystatement")
			reqPS.WriteString(" JOIN productprecautionarystatements ON productprecautionarystatements.productprecautionarystatements_precautionarystatement_id = precautionarystatement.precautionarystatement_id")
			reqPS.WriteString(" JOIN product ON productprecautionarystatements.productprecautionarystatements_product_id = product.product_id")
			reqPS.WriteString(" WHERE product.product_id = ?")

			if err = db.Select(&products[i].PrecautionaryStatements, reqPS.String(), pr.ProductID); err != nil {
				logger.Log.WithFields(logrus.Fields{"err": err}).Error("GetProducts:goroutine:ps")
			}
		}

		wg.Done()
	}()

	//
	// getting number of storages for each product
	//
	if !public {
		wg.Add(1)

		go func() {
			var reqtsc, reqsc, reqasc strings.Builder

			for i, pr := range products {
				// getting the total storage count
				reqtsc.Reset()
				reqtsc.WriteString("SELECT count(DISTINCT storage_id) from storage")
				reqtsc.WriteString(" WHERE storage.product = ? AND storage.storage IS NULL AND storage.storage_archive == false")

				if isadmin {
					reqsc.Reset()
					reqsc.WriteString("SELECT count(DISTINCT storage_id) from storage")
					reqsc.WriteString(" WHERE storage.product = ? AND storage.storage IS NULL AND storage.storage_archive == false")

					reqasc.Reset()
					reqasc.WriteString("SELECT count(DISTINCT storage_id) from storage")
					reqasc.WriteString(" WHERE storage.product = ? AND storage.storage_archive == true")

					if err = db.Get(&products[i].ProductSC, reqsc.String(), pr.ProductID); err != nil {
						logger.Log.WithFields(logrus.Fields{"err": err}).Error("GetProducts:goroutine:SC")
					}

					if err = db.Get(&products[i].ProductASC, reqasc.String(), pr.ProductID); err != nil {
						logger.Log.WithFields(logrus.Fields{"err": err}).Error("GetProducts:goroutine:ASC")
					}
				} else {
					// getting the storage count of the logged user entities
					reqsc.Reset()
					reqsc.WriteString("SELECT count(DISTINCT storage_id) from storage")
					reqsc.WriteString(" JOIN storelocation ON storage.storelocation = storelocation.storelocation_id")
					reqsc.WriteString(" JOIN entity ON storelocation.entity = entity.entity_id")
					reqsc.WriteString(" JOIN personentities ON (entity.entity_id = personentities.personentities_entity_id) AND")
					reqsc.WriteString(" (personentities.personentities_person_id = ?)")
					reqsc.WriteString(" WHERE storage.product = ? AND storage.storage IS NULL AND storage.storage_archive == false")

					reqasc.Reset()
					reqasc.WriteString("SELECT count(DISTINCT storage_id) from storage")
					reqasc.WriteString(" JOIN storelocation ON storage.storelocation = storelocation.storelocation_id")
					reqasc.WriteString(" JOIN entity ON storelocation.entity = entity.entity_id")
					reqasc.WriteString(" JOIN personentities ON (entity.entity_id = personentities.personentities_entity_id) AND")
					reqasc.WriteString(" (personentities.personentities_person_id = ?)")
					reqasc.WriteString(" WHERE storage.product = ? AND storage.storage_archive == true")

					if err = db.Get(&products[i].ProductSC, reqsc.String(), person_id, pr.ProductID); err != nil {
						logger.Log.WithFields(logrus.Fields{"err": err}).Error("GetProducts:goroutine:SC")
					}

					if err = db.Get(&products[i].ProductASC, reqasc.String(), person_id, pr.ProductID); err != nil {
						logger.Log.WithFields(logrus.Fields{"err": err}).Error("GetProducts:goroutine:ASC")
					}
				}

				if err = db.Get(&products[i].ProductTSC, reqtsc.String(), pr.ProductID); err != nil {
					logger.Log.WithFields(logrus.Fields{"err": err}).Error("GetProducts:goroutine:TSC")
				}
			}

			wg.Done()
		}()
	}

	wg.Wait()

	return products, count, nil
}

// CountProducts returns the number of products.
func (db *SQLiteDataStore) CountProducts() (int, error) {
	var (
		count int
		sqlr  string
		err   error
	)

	sqlr = `SELECT count(*) FROM product`
	if err = db.Get(&count, sqlr); err != nil {
		return 0, err
	}

	logger.Log.WithFields(logrus.Fields{"count": count}).Debug("CountProducts")

	return count, nil
}

// CountProductStorages returns the number of storages for the product with the given id.
func (db *SQLiteDataStore) CountProductStorages(id int) (int, error) {
	var (
		count int
		sqlr  string
		err   error
	)

	sqlr = `SELECT count(*) FROM storage WHERE product = ?`
	if err = db.Get(&count, sqlr, id); err != nil {
		return 0, err
	}

	logger.Log.WithFields(logrus.Fields{"count": count}).Debug("CountProductStorages")

	return count, nil
}

// GetProduct returns the product with the given id.
func (db *SQLiteDataStore) GetProduct(id int) (models.Product, error) {
	var (
		product models.Product
		sqlr    string
		err     error
	)

	sqlr = `SELECT product.product_id, 
	product.product_specificity, 
	product_msds,
	product_restricted,
	product_radioactive,
	product_threedformula,
	product_twodformula,
	product_molformula,
	product_disposalcomment,
	product_remark,
	product_sheet,
	product_temperature,
	product_number_per_carton,
	product_number_per_bag,
	linearformula.linearformula_id AS "linearformula.linearformula_id",
	linearformula.linearformula_label AS "linearformula.linearformula_label",
	empiricalformula.empiricalformula_id AS "empiricalformula.empiricalformula_id",
	empiricalformula.empiricalformula_label AS "empiricalformula.empiricalformula_label",
	physicalstate.physicalstate_id AS "physicalstate.physicalstate_id",
	physicalstate.physicalstate_label AS "physicalstate.physicalstate_label",
	signalword.signalword_id AS "signalword.signalword_id",
	signalword.signalword_label AS "signalword.signalword_label",
	person.person_id AS "person.person_id",
	person.person_email AS "person.person_email",
	name.name_id AS "name.name_id",
	name.name_label AS "name.name_label",
	cenumber.cenumber_id AS "cenumber.cenumber_id",
	cenumber.cenumber_label AS "cenumber.cenumber_label",
	casnumber.casnumber_id AS "casnumber.casnumber_id",
	casnumber.casnumber_label AS "casnumber.casnumber_label",
	casnumber.casnumber_cmr AS "casnumber.casnumber_cmr",
	producerref.producerref_id AS "producerref.producerref_id",
	producerref.producerref_label AS "producerref.producerref_label",
	producer.producer_id AS "producerref.producer.producer_id",
	producer.producer_label AS "producerref.producer.producer_label",
	ut.unit_id AS "unit_temperature.unit_id",
	ut.unit_label AS "unit_temperature.unit_label",
	category.category_id AS "category.category_id",
	category.category_label AS "category.category_label"
	FROM product
	JOIN name ON product.name = name.name_id
	LEFT JOIN casnumber ON product.casnumber = casnumber.casnumber_id
	LEFT JOIN cenumber ON product.cenumber = cenumber.cenumber_id
	JOIN person ON product.person = person.person_id
	LEFT JOIN empiricalformula ON product.empiricalformula = empiricalformula.empiricalformula_id
	LEFT JOIN linearformula ON product.linearformula = linearformula.linearformula_id
	LEFT JOIN physicalstate ON product.physicalstate = physicalstate.physicalstate_id
	LEFT JOIN signalword ON product.signalword = signalword.signalword_id
	LEFT JOIN category ON product.category = category.category_id
	LEFT JOIN unit ut ON product.unit_temperature = ut.unit_id
	LEFT JOIN producerref ON product.producerref = producerref.producerref_id
	LEFT JOIN producer ON producerref.producer = producer.producer_id
	WHERE product_id = ?`
	if err = db.Get(&product, sqlr, id); err != nil {
		return models.Product{}, err
	}

	//
	// getting supplierref
	//
	sqlr = `SELECT supplierref_id, 
	supplierref_label,
	supplier.supplier_id AS "supplier.supplier_id",
	supplier.supplier_label AS "supplier.supplier_label"
	FROM supplierref
	JOIN productsupplierrefs ON productsupplierrefs.productsupplierrefs_supplierref_id = supplierref.supplierref_id AND productsupplierrefs.productsupplierrefs_product_id = ?
	JOIN supplier ON supplierref.supplier = supplier.supplier_id`
	if err = db.Select(&product.SupplierRefs, sqlr, product.ProductID); err != nil {
		return product, err
	}

	//
	// getting tags
	//
	sqlr = `SELECT tag_id, tag_label FROM tag
	JOIN producttags ON producttags.producttags_tag_id = tag.tag_id
	JOIN product ON producttags.producttags_product_id = product.product_id
	WHERE product.product_id = ?`
	if err = db.Select(&product.Tags, sqlr, product.ProductID); err != nil {
		return product, err
	}

	//
	// getting symbols
	//
	sqlr = `SELECT symbol_id, symbol_label, symbol_image FROM symbol
	JOIN productsymbols ON productsymbols.productsymbols_symbol_id = symbol.symbol_id
	JOIN product ON productsymbols.productsymbols_product_id = product.product_id
	WHERE product.product_id = ?`
	if err = db.Select(&product.Symbols, sqlr, product.ProductID); err != nil {
		return product, err
	}

	//
	// getting synonyms
	//
	sqlr = `SELECT name_id, name_label FROM name
	JOIN productsynonyms ON productsynonyms.productsynonyms_name_id = name.name_id
	JOIN product ON productsynonyms.productsynonyms_product_id = product.product_id
	WHERE product.product_id = ?`
	if err = db.Select(&product.Synonyms, sqlr, product.ProductID); err != nil {
		return product, err
	}

	//
	// getting classes of compounds
	//
	sqlr = `SELECT classofcompound_id, classofcompound_label FROM classofcompound
	JOIN productclassofcompound ON productclassofcompound.productclassofcompound_classofcompound_id = classofcompound.classofcompound_id
	JOIN product ON productclassofcompound.productclassofcompound_product_id = product.product_id
	WHERE product.product_id = ?`
	if err = db.Select(&product.ClassOfCompound, sqlr, product.ProductID); err != nil {
		return product, err
	}

	//
	// getting hazard statements
	//
	sqlr = `SELECT hazardstatement_id, hazardstatement_label, hazardstatement_reference, hazardstatement_cmr FROM hazardstatement
	JOIN producthazardstatements ON producthazardstatements.producthazardstatements_hazardstatement_id = hazardstatement.hazardstatement_id
	JOIN product ON producthazardstatements.producthazardstatements_product_id = product.product_id
	WHERE product.product_id = ?`
	if err = db.Select(&product.HazardStatements, sqlr, product.ProductID); err != nil {
		return product, err
	}

	//
	// getting precautionary statements
	//
	sqlr = `SELECT precautionarystatement_id, precautionarystatement_label, precautionarystatement_reference FROM precautionarystatement
	JOIN productprecautionarystatements ON productprecautionarystatements.productprecautionarystatements_precautionarystatement_id = precautionarystatement.precautionarystatement_id
	JOIN product ON productprecautionarystatements.productprecautionarystatements_product_id = product.product_id
	WHERE product.product_id = ?`
	if err = db.Select(&product.PrecautionaryStatements, sqlr, product.ProductID); err != nil {
		return product, err
	}

	switch {
	case product.ProductNumberPerCarton.Valid:
		product.ProductType = "CONS"
	case product.ProducerRef.ProducerRefID.Valid:
		product.ProductType = "BIO"
	default:
		product.ProductType = "CHEM"
	}

	logger.Log.WithFields(logrus.Fields{"id": id, "product": product}).Debug("GetProduct")

	return product, nil
}

// DeleteProduct deletes the product with the given id.
func (db *SQLiteDataStore) DeleteProduct(id int) error {
	var (
		sqlr string
		err  error
	)

	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("DeleteProduct")

	// deleting bookmarks
	sqlr = `DELETE FROM bookmark WHERE bookmark.product = (?)`
	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	// deleting symbols
	sqlr = `DELETE FROM productsymbols WHERE productsymbols.productsymbols_product_id = (?)`
	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	// deleting synonyms
	sqlr = `DELETE FROM productsynonyms WHERE productsynonyms.productsynonyms_product_id = (?)`
	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	// deleting classes of compounds
	sqlr = `DELETE FROM productclassofcompound WHERE productclassofcompound.productclassofcompound_product_id = (?)`
	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	// deleting hazard statements
	sqlr = `DELETE FROM producthazardstatements WHERE producthazardstatements.producthazardstatements_product_id = (?)`
	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	// deleting precautionary statements
	sqlr = `DELETE FROM productprecautionarystatements WHERE productprecautionarystatements.productprecautionarystatements_product_id = (?)`
	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	// deleting product
	sqlr = `DELETE FROM product WHERE product_id = ?`
	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}

	return nil
}

// CreateUpdateProduct insert/update the product p into the database.
func (db *SQLiteDataStore) CreateUpdateProduct(p models.Product, update bool) (lastInsertID int64, err error) {
	var (
		v    driver.Value
		sqlr string
		args []interface{}
		tx   *sql.Tx
		res  sql.Result
	)

	dialect := goqu.Dialect("sqlite3")
	tableProduct := goqu.T("product")

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

	// if CasNumberID = -1 then it is a new cas
	if v, err = p.CasNumber.CasNumberID.Value(); p.CasNumber.CasNumberID.Valid && err == nil && v.(int64) == -1 {
		logger.Log.Debug("new casnumber " + p.CasNumberLabel.String)

		sqlr = `INSERT INTO casnumber (casnumber_label) VALUES (?)`

		if res, err = tx.Exec(sqlr, p.CasNumberLabel); err != nil {
			return
		}
		// getting the last inserted id
		if lastInsertID, err = res.LastInsertId(); err != nil {
			return
		}
		// updating the product CasNumberID (CasNumberLabel already set)
		p.CasNumber.CasNumberID = sql.NullInt64{Valid: true, Int64: lastInsertID}
	}

	// if CeNumberID = -1 then it is a new ce
	if v, err = p.CeNumber.CeNumberID.Value(); p.CeNumber.CeNumberID.Valid && err == nil && v.(int64) == -1 {
		logger.Log.Debug("new cenumber " + p.CeNumberLabel.String)

		sqlr = `INSERT INTO cenumber (cenumber_label) VALUES (?)`

		if res, err = tx.Exec(sqlr, p.CeNumberLabel.String); err != nil {
			return
		}
		// getting the last inserted id
		if lastInsertID, err = res.LastInsertId(); err != nil {
			return
		}
		// updating the product CeNumberID (CeNumberLabel already set)
		p.CeNumber.CeNumberID = sql.NullInt64{Valid: true, Int64: lastInsertID}
	}

	if err != nil {
		logger.Log.Error("cenumber error - " + err.Error())
		return
	}

	// if NameID = -1 then it is a new name
	if p.Name.NameID == -1 {
		logger.Log.Debug("new name " + p.NameLabel)

		sqlr = `INSERT INTO name (name_label) VALUES (?)`

		if res, err = tx.Exec(sqlr, strings.ToUpper(p.NameLabel)); err != nil {
			return
		}
		// getting the last inserted id
		if lastInsertID, err = res.LastInsertId(); err != nil {
			return
		}
		// updating the product NameID (NameLabel already set)
		p.Name.NameID = int(lastInsertID)
	}

	// if NameID = -1 then it is a new name
	for i, syn := range p.Synonyms {
		if syn.NameID == -1 {
			logger.Log.Debug("new name(syn) " + syn.NameLabel)

			sqlr = `INSERT INTO name (name_label) VALUES (?)`

			if res, err = tx.Exec(sqlr, strings.ToUpper(syn.NameLabel)); err != nil {
				return
			}

			// getting the last inserted id
			if lastInsertID, err = res.LastInsertId(); err != nil {
				return
			}

			p.Synonyms[i].NameID = int(lastInsertID)
		}
	}

	// if ClassOfCompoundID = -1 then it is a new class of compounds
	for i, coc := range p.ClassOfCompound {
		if coc.ClassOfCompoundID == -1 {
			logger.Log.Debug("new classofcompound " + coc.ClassOfCompoundLabel)

			sqlr = `INSERT INTO classofcompound (classofcompound_label) VALUES (?)`

			if res, err = tx.Exec(sqlr, strings.ToUpper(coc.ClassOfCompoundLabel)); err != nil {
				return
			}
			// getting the last inserted id
			if lastInsertID, err = res.LastInsertId(); err != nil {
				return
			}

			p.ClassOfCompound[i].ClassOfCompoundID = int(lastInsertID)
		}
	}

	// if SupplierRefID = -1 then it is a new supplier ref
	for i, sr := range p.SupplierRefs {
		if sr.SupplierRefID == -1 {
			logger.Log.Debug("new supplierref " + sr.SupplierRefLabel)

			sqlr = `INSERT INTO supplierref (supplierref_label, supplier) VALUES (?, ?)`

			if res, err = tx.Exec(sqlr, sr.SupplierRefLabel, sr.Supplier.SupplierID); err != nil {
				return
			}
			// getting the last inserted id
			if lastInsertID, err = res.LastInsertId(); err != nil {
				return
			}

			p.SupplierRefs[i].SupplierRefID = int(lastInsertID)
		}
	}

	// if TagID = -1 then it is a new tag
	for i, tag := range p.Tags {
		if tag.TagID == -1 {
			logger.Log.Debug("new tag " + tag.TagLabel)

			sqlr = `INSERT INTO tag (tag_label) VALUES (?)`

			if res, err = tx.Exec(sqlr, tag.TagLabel); err != nil {
				return
			}
			// getting the last inserted id
			if lastInsertID, err = res.LastInsertId(); err != nil {
				return
			}

			p.Tags[i].TagID = int(lastInsertID)
		}
	}

	// if EmpiricalFormulaID = -1 then it is a new empirical formula
	if v, err = p.EmpiricalFormula.EmpiricalFormulaID.Value(); p.EmpiricalFormula.EmpiricalFormulaID.Valid && err == nil && v.(int64) == -1 {
		logger.Log.Debug("new empiricalformula " + p.EmpiricalFormulaLabel.String)

		sqlr = `INSERT INTO empiricalformula (empiricalformula_label) VALUES (?)`

		if res, err = tx.Exec(sqlr, p.EmpiricalFormulaLabel); err != nil {
			return
		}
		// getting the last inserted id
		if lastInsertID, err = res.LastInsertId(); err != nil {
			return
		}
		// updating the product EmpiricalFormulaID (EmpiricalFormulaLabel already set)
		p.EmpiricalFormula.EmpiricalFormulaID = sql.NullInt64{Valid: true, Int64: lastInsertID}
	}

	// if LinearFormulaID = -1 then it is a new linear formula
	if v, err = p.LinearFormula.LinearFormulaID.Value(); p.LinearFormula.LinearFormulaID.Valid && err == nil && v.(int64) == -1 {
		logger.Log.Debug("new linearformula " + p.LinearFormulaLabel.String)

		sqlr = `INSERT INTO linearformula (linearformula_label) VALUES (?)`

		if res, err = tx.Exec(sqlr, p.LinearFormulaLabel.String); err != nil {
			return
		}
		// getting the last inserted id
		if lastInsertID, err = res.LastInsertId(); err != nil {
			return
		}
		// updating the product LinearFormulaID (LinearFormulaLabel already set)
		p.LinearFormula.LinearFormulaID = sql.NullInt64{Valid: true, Int64: lastInsertID}
	}

	// if PhysicalStateID = -1 then it is a new physical state
	if v, err = p.PhysicalState.PhysicalStateID.Value(); p.PhysicalState.PhysicalStateID.Valid && err == nil && v.(int64) == -1 {
		logger.Log.Debug("new physicalstate " + p.PhysicalStateLabel.String)

		sqlr = `INSERT INTO physicalstate (physicalstate_label) VALUES (?)`

		if res, err = tx.Exec(sqlr, p.PhysicalStateLabel.String); err != nil {
			return
		}
		// getting the last inserted id
		if lastInsertID, err = res.LastInsertId(); err != nil {
			return
		}
		// updating the product PhysicalStateID (PhysicalStateLabel already set)
		p.PhysicalState.PhysicalStateID = sql.NullInt64{Valid: true, Int64: lastInsertID}
	}

	// if CategoryID = -1 then it is a new category
	if v, err = p.Category.CategoryID.Value(); p.Category.CategoryID.Valid && err == nil && v.(int64) == -1 {
		logger.Log.Debug("new category " + p.CategoryLabel.String)

		sqlr = `INSERT INTO category (category_label) VALUES (?)`

		if res, err = tx.Exec(sqlr, p.CategoryLabel.String); err != nil {
			return
		}
		// getting the last inserted id
		if lastInsertID, err = res.LastInsertId(); err != nil {
			return
		}
		// updating the product PhysicalStateID (PhysicalStateLabel already set)
		p.Category.CategoryID = sql.NullInt64{Valid: true, Int64: lastInsertID}
	}

	// if ProducerRefID = -1 then it is a new producer ref
	if v, err = p.ProducerRef.ProducerRefID.Value(); p.ProducerRef.ProducerRefID.Valid && err == nil && v.(int64) == -1 {
		logger.Log.Debug("new producerref " + p.ProducerRefLabel.String)

		sqlr = `INSERT INTO producerref (producerref_label, producer) VALUES (?, ?)`

		if res, err = tx.Exec(sqlr, p.ProducerRefLabel.String, p.Producer.ProducerID); err != nil {
			return
		}
		// getting the last inserted id
		if lastInsertID, err = res.LastInsertId(); err != nil {
			return
		}
		// updating the product ProducerRefID (ProducerRefLabel already set)
		p.ProducerRef.ProducerRefID = sql.NullInt64{Valid: true, Int64: lastInsertID}
	}

	// finally updating the product
	insertCols := goqu.Record{}

	if p.ProductSpecificity.Valid {
		insertCols["product_specificity"] = p.ProductSpecificity.String
	} else {
		insertCols["product_specificity"] = nil
	}

	if p.ProductMSDS.Valid {
		insertCols["product_msds"] = p.ProductMSDS.String
	} else {
		insertCols["product_msds"] = nil
	}

	if p.ProductSheet.Valid {
		insertCols["product_sheet"] = p.ProductSheet.String
	} else {
		insertCols["product_sheet"] = nil
	}

	if p.ProductTemperature.Valid {
		insertCols["product_temperature"] = int(p.ProductTemperature.Int64)
	} else {
		insertCols["product_temperature"] = nil
	}

	if p.ProductRestricted.Valid {
		insertCols["product_restricted"] = p.ProductRestricted.Bool
	} else {
		insertCols["product_restricted"] = false
	}

	if p.ProductRadioactive.Valid {
		insertCols["product_radioactive"] = p.ProductRadioactive.Bool
	} else {
		insertCols["product_radioactive"] = false
	}

	if p.Category.CategoryID.Valid {
		insertCols["category"] = int(p.Category.CategoryID.Int64)
	} else {
		insertCols["category"] = nil
	}

	if p.UnitTemperature.UnitID.Valid {
		insertCols["unit_temperature"] = int(p.UnitTemperature.UnitID.Int64)
	} else {
		insertCols["unit_temperature"] = nil
	}

	if p.ProductThreeDFormula.Valid {
		insertCols["product_threedformula"] = p.ProductThreeDFormula.String
	} else {
		insertCols["product_threedformula"] = nil
	}

	if p.ProductTwoDFormula.Valid {
		insertCols["product_twodformula"] = p.ProductTwoDFormula.String
	}
	// } else {
	// 	insertCols["product_twodformula"] = nil
	// }

	if p.ProductDisposalComment.Valid {
		insertCols["product_disposalcomment"] = p.ProductDisposalComment.String
	} else {
		insertCols["product_disposalcomment"] = nil
	}

	if p.ProductRemark.Valid {
		insertCols["product_remark"] = p.ProductRemark.String
	} else {
		insertCols["product_remark"] = nil
	}

	if p.ProductNumberPerCarton.Valid {
		insertCols["product_number_per_carton"] = p.ProductNumberPerCarton.Int64
	} else {
		insertCols["product_number_per_carton"] = nil
	}

	if p.ProductNumberPerBag.Valid {
		insertCols["product_number_per_bag"] = p.ProductNumberPerBag.Int64
	} else {
		insertCols["product_number_per_bag"] = nil
	}

	if p.EmpiricalFormulaID.Valid {
		insertCols["empiricalformula"] = int(p.EmpiricalFormulaID.Int64)
	} else {
		insertCols["empiricalformula"] = nil
	}

	if p.LinearFormulaID.Valid {
		insertCols["linearformula"] = int(p.LinearFormulaID.Int64)
	} else {
		insertCols["linearformula"] = nil
	}

	if p.PhysicalStateID.Valid {
		insertCols["physicalstate"] = int(p.PhysicalStateID.Int64)
	} else {
		insertCols["physicalstate"] = nil
	}

	if p.SignalWordID.Valid {
		insertCols["signalword"] = int(p.SignalWordID.Int64)
	} else {
		insertCols["signalword"] = nil
	}

	if p.CasNumberID.Valid {
		insertCols["casnumber"] = int(p.CasNumberID.Int64)
	} else {
		insertCols["casnumber"] = nil
	}

	if p.CeNumberID.Valid {
		insertCols["cenumber"] = int(p.CeNumberID.Int64)
	} else {
		insertCols["cenumber"] = nil
	}

	if p.ProducerRefID.Valid {
		insertCols["producerref"] = int(p.ProducerRefID.Int64)
	} else {
		insertCols["producerref"] = nil
	}

	if p.ProductMolFormula.Valid {
		insertCols["product_molformula"] = p.ProductMolFormula.String
	} else {
		insertCols["product_molformula"] = nil
	}

	insertCols["name"] = p.NameID
	insertCols["person"] = p.PersonID

	if update {
		iQuery := dialect.Update(tableProduct).Set(insertCols).Where(goqu.I("product_id").Eq(p.ProductID))
		if sqlr, args, err = iQuery.ToSQL(); err != nil {
			return
		}
	} else {
		iQuery := dialect.Insert(tableProduct).Rows(insertCols)
		if sqlr, args, err = iQuery.ToSQL(); err != nil {
			return
		}
	}

	// logger.Log.Debug(sqlr)
	// logger.Log.Debug(args)

	if res, err = tx.Exec(sqlr, args...); err != nil {
		return
	}

	if !update {
		// getting the last inserted id
		if lastInsertID, err = res.LastInsertId(); err != nil {
			return
		}

		p.ProductID = int(lastInsertID)
	}

	// adding supplierrefs
	if update {
		sqlr = `DELETE FROM productsupplierrefs WHERE productsupplierrefs.productsupplierrefs_product_id = (?)`
		if _, err = tx.Exec(sqlr, p.ProductID); err != nil {
			logger.Log.Error("error DELETE FROM productsupplierrefs")
			return
		}
	}

	for _, sr := range p.SupplierRefs {
		sqlr = `INSERT INTO productsupplierrefs (productsupplierrefs_product_id, productsupplierrefs_supplierref_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, sr.SupplierRefID); err != nil {
			logger.Log.Error("error INSERT INTO productsupplierrefs")
			return
		}
	}

	// adding tags
	if update {
		sqlr = `DELETE FROM producttags WHERE producttags.producttags_product_id = (?)`
		if _, err = tx.Exec(sqlr, p.ProductID); err != nil {
			logger.Log.Error("error DELETE FROM producttags")
			return
		}
	}

	for _, tag := range p.Tags {
		sqlr = `INSERT INTO producttags (producttags_product_id, producttags_tag_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, tag.TagID); err != nil {
			logger.Log.Error("error INSERT INTO producttags")
			return
		}
	}

	// adding symbols
	if update {
		sqlr = `DELETE FROM productsymbols WHERE productsymbols.productsymbols_product_id = (?)`
		if _, err = tx.Exec(sqlr, p.ProductID); err != nil {
			logger.Log.Error("error DELETE FROM productsymbols")
			return
		}
	}

	for _, sym := range p.Symbols {
		sqlr = `INSERT INTO productsymbols (productsymbols_product_id, productsymbols_symbol_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, sym.SymbolID); err != nil {
			logger.Log.Error("error INSERT INTO productsymbols")
			return
		}
	}

	// adding classes of compounds
	if update {
		sqlr = `DELETE FROM productclassofcompound WHERE productclassofcompound.productclassofcompound_product_id = (?)`
		if _, err = tx.Exec(sqlr, p.ProductID); err != nil {
			logger.Log.Error("error DELETE FROM productclassofcompound")
			return
		}
	}

	for _, coc := range p.ClassOfCompound {
		sqlr = `INSERT INTO productclassofcompound (productclassofcompound_product_id, productclassofcompound_classofcompound_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, coc.ClassOfCompoundID); err != nil {
			logger.Log.Error("error INSERT INTO productclassofcompound")
			return
		}
	}

	// adding hazard statements
	if update {
		sqlr = `DELETE FROM producthazardstatements WHERE producthazardstatements.producthazardstatements_product_id = (?)`
		if _, err = tx.Exec(sqlr, p.ProductID); err != nil {
			logger.Log.Error("error DELETE FROM producthazardstatements")
			return
		}
	}

	for _, hs := range p.HazardStatements {
		sqlr = `INSERT INTO producthazardstatements (producthazardstatements_product_id, producthazardstatements_hazardstatement_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, hs.HazardStatementID); err != nil {
			logger.Log.Error("error INSERT INTO producthazardstatements")
			return
		}
	}

	// adding precautionary statements
	if update {
		sqlr = `DELETE FROM productprecautionarystatements WHERE productprecautionarystatements.productprecautionarystatements_product_id = (?)`
		if _, err = tx.Exec(sqlr, p.ProductID); err != nil {
			logger.Log.Error("error DELETE FROM productprecautionarystatements")
			return
		}
	}

	for _, ps := range p.PrecautionaryStatements {
		sqlr = `INSERT INTO productprecautionarystatements (productprecautionarystatements_product_id, productprecautionarystatements_precautionarystatement_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, ps.PrecautionaryStatementID); err != nil {
			logger.Log.Error("error INSERT INTO productprecautionarystatements")
			return
		}
	}

	// adding synonyms
	if update {
		sqlr = `DELETE FROM productsynonyms WHERE productsynonyms.productsynonyms_product_id = (?)`
		if _, err = tx.Exec(sqlr, p.ProductID); err != nil {
			logger.Log.Error("error DELETE FROM productsynonyms")
			return
		}
	}

	for _, syn := range p.Synonyms {
		sqlr = `INSERT INTO productsynonyms (productsynonyms_product_id, productsynonyms_name_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, syn.NameID); err != nil {
			logger.Log.Error("error INSERT INTO productsynonyms")
			return
		}
	}

	return
}

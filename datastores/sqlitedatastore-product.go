package datastores

import (
	"database/sql"
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
func (db *SQLiteDataStore) GetProducts(f zmqclient.RequestFilter, person_id int, public bool) ([]models.Product, int, error) {
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

	// hack to bypass optionnal default on the Rust part.
	if f.Search == "" {
		f.Search = "%%"
	}

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
	p.product_inchi,
	p.product_inchikey,
	p.product_specificity, 
	p.product_canonical_smiles,
	p.product_molecular_weight,
	p.product_msds,
	p.product_restricted,
	p.product_radioactive,
	p.product_threed_formula,
	p.product_twod_formula,
	p.product_disposal_comment,
	p.product_remark,
	p.product_sheet,
	p.product_temperature,
	p.product_number_per_carton,
	p.product_number_per_bag,
	linear_formula.linear_formula_id AS "linear_formula.linear_formula_id",
	linear_formula.linear_formula_label AS "linear_formula.linear_formula_label",
	empirical_formula.empirical_formula_id AS "empirical_formula.empirical_formula_id",
	empirical_formula.empirical_formula_label AS "empirical_formula.empirical_formula_label",
	physical_state.physical_state_id AS "physical_state.physical_state_id",
	physical_state.physical_state_label AS "physical_state.physical_state_label",
	signal_word.signal_word_id AS "signal_word.signal_word_id",
	signal_word.signal_word_label AS "signal_word.signal_word_label",
	person.person_id AS "person.person_id",
	person.person_email AS "person.person_email",
	name.name_id AS "name.name_id",
	name.name_label AS "name.name_label",
	bookmark.bookmark_id AS "bookmark.bookmark_id",
	ce_number.ce_number_id AS "ce_number.ce_number_id",
	ce_number.ce_number_label AS "ce_number.ce_number_label",
	cas_number.cas_number_id AS "cas_number.cas_number_id",
	cas_number.cas_number_label AS "cas_number.cas_number_label",
	cas_number.cas_number_cmr AS "cas_number.cas_number_cmr",
	producer_ref.producer_ref_id AS "producer_ref.producer_ref_id",
	producer_ref.producer_ref_label AS "producer_ref.producer_ref_label",
	producer.producer_id AS "producer_ref.producer.producer_id",
	producer.producer_label AS "producer_ref.producer.producer_label",
	ut.unit_id AS "unit_temperature.unit_id",
	ut.unit_label AS "unit_temperature.unit_label",
	umw.unit_id AS "unit_molecularweight.unit_id",
	umw.unit_label AS "unit_molecularweight.unit_label",
	category.category_id AS "category.category_id",
	category.category_label AS "category.category_label"
	`)

	if !public {
		presreq.WriteString(`,GROUP_CONCAT(DISTINCT storage.storage_barecode) AS "product_sl"`)
	}

	if f.CasNumberCmr {
		presreq.WriteString(`,GROUP_CONCAT(DISTINCT hazard_statement.hazard_statement_cmr) AS "hazard_statement_cmr"`)
	}

	// common parts
	comreq.WriteString(" FROM product as p")
	// CMR
	if f.CasNumberCmr {
		comreq.WriteString(" LEFT JOIN producthazardstatements ON producthazardstatements.producthazardstatements_product_id = p.product_id")
		comreq.WriteString(" LEFT JOIN hazard_statement ON producthazardstatements.producthazardstatements_hazard_statement_id = hazard_statement.hazard_statement_id")
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
	// get unit_molecularweight
	comreq.WriteString(" LEFT JOIN unit umw ON p.unit_molecular_weight = umw.unit_id")
	// get producer_ref
	if f.ProducerRef != 0 {
		comreq.WriteString(" JOIN producer_ref ON p.producer_ref = :producer_ref")
	} else {
		comreq.WriteString(" LEFT JOIN producer_ref ON p.producer_ref = producer_ref.producer_ref_id")
	}
	// get producer
	comreq.WriteString(" LEFT JOIN producer ON producer_ref.producer = producer.producer_id")
	// get cas_number
	comreq.WriteString(" LEFT JOIN cas_number ON p.cas_number = cas_number.cas_number_id")
	// get ce_number
	comreq.WriteString(" LEFT JOIN ce_number ON p.ce_number = ce_number.ce_number_id")
	// get person
	comreq.WriteString(" JOIN person ON p.person = person.person_id")
	// get physical state
	comreq.WriteString(" LEFT JOIN physical_state ON p.physical_state = physical_state.physical_state_id")
	// get signal word
	comreq.WriteString(" LEFT JOIN signal_word ON p.signal_word = signal_word.signal_word_id")
	// get empirical formula
	comreq.WriteString(" LEFT JOIN empirical_formula ON p.empirical_formula = empirical_formula.empirical_formula_id")
	// get linear formula
	comreq.WriteString(" LEFT JOIN linear_formula ON p.linear_formula = linear_formula.linear_formula_id")
	// get bookmark
	comreq.WriteString(" LEFT JOIN bookmark ON (bookmark.product = p.product_id AND bookmark.person = :personid)")
	// get storages, store locations and entities
	comreq.WriteString(" LEFT JOIN storage ON storage.product = p.product_id")

	if f.Entity != 0 || f.Storelocation != 0 || f.StorageBarecode != "" {
		comreq.WriteString(" JOIN store_location ON storage.store_location = store_location.store_location_id")
		comreq.WriteString(" JOIN entity ON store_location.entity = entity.entity_id")
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
	// get hazard_statements
	if len(f.HazardStatements) != 0 {
		comreq.WriteString(" JOIN producthazardstatements AS phs ON phs.producthazardstatements_product_id = p.product_id")
	}
	// get precautionary_statements
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
		comreq.WriteString(" AND storage.storage_to_destroy = true")
	}

	if f.CasNumberCmr {
		comreq.WriteString(" AND (cas_number.cas_number_cmr IS NOT NULL OR (hazard_statement_cmr IS NOT NULL AND hazard_statement_cmr != ''))")
	}

	if f.Product != 0 {
		comreq.WriteString(" AND p.product_id = :product")
	}

	if f.Entity != 0 {
		comreq.WriteString(" AND entity.entity_id = :entity")
	}

	if f.Storelocation != 0 {
		comreq.WriteString(" AND store_location.store_location_id = :store_location")
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
		comreq.WriteString(" AND cas_number.cas_number_id = :cas_number")
	}

	if f.EmpiricalFormula != 0 {
		comreq.WriteString(" AND empirical_formula.empirical_formula_id = :empirical_formula")
	}

	if f.StorageBarecode != "" {
		comreq.WriteString(" AND storage.storage_barecode LIKE :storage_barecode")
	}

	if f.StorageBatchNumber != "" {
		comreq.WriteString(" AND storage.storage_batch_number LIKE :storage_batch_number")
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
		comreq.WriteString(" AND phs.producthazardstatements_hazard_statement_id IN (")

		for _, s := range f.HazardStatements {
			comreq.WriteString(fmt.Sprintf("%d,", s))
		}

		// to complete the last comma
		comreq.WriteString("-1")
		comreq.WriteString(" )")
	}

	if len(f.PrecautionaryStatements) != 0 {
		comreq.WriteString(" AND pps.productprecautionarystatements_precautionary_statement_id IN (")

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

	// filter restricted product
	if !rperm {
		comreq.WriteString(" AND p.product_restricted = false")
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
		"search":               f.Search,
		"personid":             person_id,
		"order":                f.Order,
		"limit":                f.Limit,
		"offset":               f.Offset,
		"entity":               f.Entity,
		"product":              f.Product,
		"store_location":       f.Storelocation,
		"name":                 f.Name,
		"cas_number":           f.CasNumber,
		"empirical_formula":    f.EmpiricalFormula,
		"product_specificity":  f.ProductSpecificity,
		"storage_barecode":     f.StorageBarecode,
		"storage_batch_number": f.StorageBatchNumber,
		"custom_name_part_of":  "%" + f.CustomNamePartOf + "%",
		"signal_word":          f.SignalWord,
		"producer_ref":         f.ProducerRef,
		"category":             f.Category,
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
			case pr.ProductNumberPerCarton != nil:
				products[i].ProductType = "CONS"
			case pr.ProducerRef.ProducerRefID != nil:
				products[i].ProductType = "BIO"
			default:
				products[i].ProductType = "CHEM"
			}
		}

		wg.Done()
	}()

	//
	// cleaning product_sl (storage barecodes concatenation)
	//
	if !public {
		wg.Add(1)

		go func() {
			r := regexp.MustCompile(`([a-zA-Z]{1}[0-9]+)\.[0-9]+`)
			for i, pr := range products {
				// note: do not modify p but products[i] instead
				m := r.FindAllStringSubmatch(*pr.ProductSL, -1)

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
						*products[i].ProductSL = m[0][1]
					} else {
						*products[i].ProductSL = ""
					}
				} else {
					*products[i].ProductSL = ""
				}
			}

			wg.Done()
		}()
	}

	//
	// getting supplier_ref
	//
	if !public {
		wg.Add(1)

		go func() {
			var reqSupplierref strings.Builder

			for i, pr := range products {
				// note: do not modify p but products[i] instead
				reqSupplierref.Reset()
				reqSupplierref.WriteString(`SELECT supplier_ref_id,
			supplier_ref_label,
			supplier.supplier_id AS "supplier.supplier_id",
			supplier.supplier_label AS "supplier.supplier_label"
			FROM supplier_ref`)
				reqSupplierref.WriteString(" JOIN productsupplierrefs ON productsupplierrefs.productsupplierrefs_supplier_ref_id = supplier_ref.supplier_ref_id AND productsupplierrefs.productsupplierrefs_product_id = ?")
				reqSupplierref.WriteString(" JOIN supplier ON supplier_ref.supplier = supplier.supplier_id")

				if err = db.Select(&products[i].SupplierRefs, reqSupplierref.String(), pr.ProductID); err != nil {
					logger.Log.WithFields(logrus.Fields{"err": err}).Error("GetProducts:goroutine:supplier_ref")
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
			reqSymbols.WriteString("SELECT symbol_id, symbol_label FROM symbol")
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
			reqCoc.WriteString("SELECT class_of_compound_id, class_of_compound_label FROM class_of_compound")
			reqCoc.WriteString(" JOIN productclassesofcompounds ON productclassesofcompounds.productclassesofcompounds_class_of_compound_id = class_of_compound.class_of_compound_id")
			reqCoc.WriteString(" JOIN product ON productclassesofcompounds.productclassesofcompounds_product_id = product.product_id")
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
			reqHS.WriteString("SELECT hazard_statement_id, hazard_statement_label, hazard_statement_reference, hazard_statement_cmr FROM hazard_statement")
			reqHS.WriteString(" JOIN producthazardstatements ON producthazardstatements.producthazardstatements_hazard_statement_id = hazard_statement.hazard_statement_id")
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
			reqPS.WriteString("SELECT precautionary_statement_id, precautionary_statement_label, precautionary_statement_reference FROM precautionary_statement")
			reqPS.WriteString(" JOIN productprecautionarystatements ON productprecautionarystatements.productprecautionarystatements_precautionary_statement_id = precautionary_statement.precautionary_statement_id")
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
					reqsc.WriteString(" JOIN store_location ON storage.store_location = store_location.store_location_id")
					reqsc.WriteString(" JOIN entity ON store_location.entity = entity.entity_id")
					reqsc.WriteString(" JOIN personentities ON (entity.entity_id = personentities.personentities_entity_id) AND")
					reqsc.WriteString(" (personentities.personentities_person_id = ?)")
					reqsc.WriteString(" WHERE storage.product = ? AND storage.storage IS NULL AND storage.storage_archive == false")

					reqasc.Reset()
					reqasc.WriteString("SELECT count(DISTINCT storage_id) from storage")
					reqasc.WriteString(" JOIN store_location ON storage.store_location = store_location.store_location_id")
					reqasc.WriteString(" JOIN entity ON store_location.entity = entity.entity_id")
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
	product.product_inchi,
	product.product_inchikey,
	product.product_specificity,
	product.product_canonical_smiles, 
	product.product_molecular_weight,
	product_msds,
	product_restricted,
	product_radioactive,
	product_threed_formula,
	product_twod_formula,
	product_disposal_comment,
	product_remark,
	product_sheet,
	product_temperature,
	product_number_per_carton,
	product_number_per_bag,
	linear_formula.linear_formula_id AS "linear_formula.linear_formula_id",
	linear_formula.linear_formula_label AS "linear_formula.linear_formula_label",
	empirical_formula.empirical_formula_id AS "empirical_formula.empirical_formula_id",
	empirical_formula.empirical_formula_label AS "empirical_formula.empirical_formula_label",
	physical_state.physical_state_id AS "physical_state.physical_state_id",
	physical_state.physical_state_label AS "physical_state.physical_state_label",
	signal_word.signal_word_id AS "signal_word.signal_word_id",
	signal_word.signal_word_label AS "signal_word.signal_word_label",
	person.person_id AS "person.person_id",
	person.person_email AS "person.person_email",
	name.name_id AS "name.name_id",
	name.name_label AS "name.name_label",
	ce_number.ce_number_id AS "ce_number.ce_number_id",
	ce_number.ce_number_label AS "ce_number.ce_number_label",
	cas_number.cas_number_id AS "cas_number.cas_number_id",
	cas_number.cas_number_label AS "cas_number.cas_number_label",
	cas_number.cas_number_cmr AS "cas_number.cas_number_cmr",
	producer_ref.producer_ref_id AS "producer_ref.producer_ref_id",
	producer_ref.producer_ref_label AS "producer_ref.producer_ref_label",
	producer.producer_id AS "producer_ref.producer.producer_id",
	producer.producer_label AS "producer_ref.producer.producer_label",
	ut.unit_id AS "unit_temperature.unit_id",
	ut.unit_label AS "unit_temperature.unit_label",
	umw.unit_id AS "unit_molecular_weight.unit_id",
	umw.unit_label AS "unit_molecular_weight.unit_label",
	category.category_id AS "category.category_id",
	category.category_label AS "category.category_label"
	FROM product
	JOIN name ON product.name = name.name_id
	LEFT JOIN cas_number ON product.cas_number = cas_number.cas_number_id
	LEFT JOIN ce_number ON product.ce_number = ce_number.ce_number_id
	JOIN person ON product.person = person.person_id
	LEFT JOIN empirical_formula ON product.empirical_formula = empirical_formula.empirical_formula_id
	LEFT JOIN linear_formula ON product.linear_formula = linear_formula.linear_formula_id
	LEFT JOIN physical_state ON product.physical_state = physical_state.physical_state_id
	LEFT JOIN signal_word ON product.signal_word = signal_word.signal_word_id
	LEFT JOIN category ON product.category = category.category_id
	LEFT JOIN unit ut ON product.unit_temperature = ut.unit_id
	LEFT JOIN unit umw ON product.unit_temperature = umw.unit_id
	LEFT JOIN producer_ref ON product.producer_ref = producer_ref.producer_ref_id
	LEFT JOIN producer ON producer_ref.producer = producer.producer_id
	WHERE product_id = ?`
	if err = db.Get(&product, sqlr, id); err != nil {
		return models.Product{}, err
	}

	//
	// getting supplier_ref
	//
	sqlr = `SELECT supplier_ref_id,
	supplier_ref_label,
	supplier.supplier_id AS "supplier.supplier_id",
	supplier.supplier_label AS "supplier.supplier_label"
	FROM supplier_ref
	JOIN productsupplierrefs ON productsupplierrefs.productsupplierrefs_supplier_ref_id = supplier_ref.supplier_ref_id AND productsupplierrefs.productsupplierrefs_product_id = ?
	JOIN supplier ON supplier_ref.supplier = supplier.supplier_id`
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
	sqlr = `SELECT symbol_id, symbol_label FROM symbol
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
	sqlr = `SELECT class_of_compound_id, class_of_compound_label FROM class_of_compound
	JOIN productclassesofcompounds ON productclassesofcompounds.productclassesofcompounds_class_of_compound_id = class_of_compound.class_of_compound_id
	JOIN product ON productclassesofcompounds.productclassesofcompounds_product_id = product.product_id
	WHERE product.product_id = ?`
	if err = db.Select(&product.ClassOfCompound, sqlr, product.ProductID); err != nil {
		return product, err
	}

	//
	// getting hazard statements
	//
	sqlr = `SELECT hazard_statement_id, hazard_statement_label, hazard_statement_reference, hazard_statement_cmr FROM hazard_statement
	JOIN producthazardstatements ON producthazardstatements.producthazardstatements_hazard_statement_id = hazard_statement.hazard_statement_id
	JOIN product ON producthazardstatements.producthazardstatements_product_id = product.product_id
	WHERE product.product_id = ?`
	if err = db.Select(&product.HazardStatements, sqlr, product.ProductID); err != nil {
		return product, err
	}

	//
	// getting precautionary statements
	//
	sqlr = `SELECT precautionary_statement_id, precautionary_statement_label, precautionary_statement_reference FROM precautionary_statement
	JOIN productprecautionarystatements ON productprecautionarystatements.productprecautionarystatements_precautionary_statement_id = precautionary_statement.precautionary_statement_id
	JOIN product ON productprecautionarystatements.productprecautionarystatements_product_id = product.product_id
	WHERE product.product_id = ?`
	if err = db.Select(&product.PrecautionaryStatements, sqlr, product.ProductID); err != nil {
		return product, err
	}

	switch {
	case product.ProductNumberPerCarton != nil:
		product.ProductType = "CONS"
	case product.ProducerRef.ProducerRefID != nil:
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
	sqlr = `DELETE FROM productclassesofcompounds WHERE productclassesofcompounds.productclassesofcompounds_product_id = (?)`
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
	if p.CasNumber.CasNumberID != nil && err == nil && *p.CasNumber.CasNumberID == -1 {
		// logger.Log.Debug("new cas_number " + p.CasNumberLabel)
		logger.Log.Debug("new cas_number " + *p.CasNumberLabel)

		sqlr = `INSERT INTO cas_number (cas_number_label) VALUES (?)`

		if res, err = tx.Exec(sqlr, p.CasNumberLabel); err != nil {
			return
		}
		// getting the last inserted id
		if lastInsertID, err = res.LastInsertId(); err != nil {
			return
		}
		// updating the product CasNumberID (CasNumberLabel already set)
		// p.CasNumber.CasNumberID = sql.NullInt64{Valid: true, Int64: lastInsertID}
		p.CasNumber.CasNumberID = &lastInsertID
	}

	// if CeNumberID = -1 then it is a new ce
	if p.CeNumber.CeNumberID != nil && err == nil && *p.CeNumber.CeNumberID == -1 {
		logger.Log.Debug("new ce_number " + *p.CeNumberLabel)

		sqlr = `INSERT INTO ce_number (ce_number_label) VALUES (?)`

		if res, err = tx.Exec(sqlr, *p.CeNumberLabel); err != nil {
			return
		}
		// getting the last inserted id
		if lastInsertID, err = res.LastInsertId(); err != nil {
			return
		}
		// updating the product CeNumberID (CeNumberLabel already set)
		var CeNumberIDPointer *int64
		CeNumberIDPointer = new(int64)
		*CeNumberIDPointer = lastInsertID
		p.CeNumber.CeNumberID = CeNumberIDPointer
	}

	if err != nil {
		logger.Log.Error("ce_number error - " + err.Error())
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
			logger.Log.Debug("new class_of_compound " + coc.ClassOfCompoundLabel)

			sqlr = `INSERT INTO class_of_compound (class_of_compound_label) VALUES (?)`

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
			logger.Log.Debug("new supplier_ref " + sr.SupplierRefLabel)

			sqlr = `INSERT INTO supplier_ref (supplier_ref_label, supplier) VALUES (?, ?)`

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
	if p.EmpiricalFormula.EmpiricalFormulaID != nil && err == nil && *p.EmpiricalFormula.EmpiricalFormulaID == -1 {
		logger.Log.Debug("new empirical_formula " + *p.EmpiricalFormulaLabel)

		sqlr = `INSERT INTO empirical_formula (empirical_formula_label) VALUES (?)`

		if res, err = tx.Exec(sqlr, p.EmpiricalFormulaLabel); err != nil {
			return
		}
		// getting the last inserted id
		if lastInsertID, err = res.LastInsertId(); err != nil {
			return
		}
		// updating the product EmpiricalFormulaIDPointer (EmpiricalFormulaLabel already set)
		var EmpiricalFormulaIDPointer *int64
		EmpiricalFormulaIDPointer = new(int64)
		*EmpiricalFormulaIDPointer = lastInsertID
		p.EmpiricalFormula.EmpiricalFormulaID = EmpiricalFormulaIDPointer
	}

	// if LinearFormulaID = -1 then it is a new linear formula
	if p.LinearFormula.LinearFormulaID != nil && err == nil && *p.LinearFormula.LinearFormulaID == -1 {
		logger.Log.Debug("new linear_formula " + *p.LinearFormulaLabel)

		sqlr = `INSERT INTO linear_formula (linear_formula_label) VALUES (?)`

		if res, err = tx.Exec(sqlr, p.LinearFormulaLabel); err != nil {
			return
		}
		// getting the last inserted id
		if lastInsertID, err = res.LastInsertId(); err != nil {
			return
		}
		// updating the product LinearFormulaID (LinearFormulaLabel already set)
		var LinearFormulaID *int64
		LinearFormulaID = new(int64)
		*LinearFormulaID = lastInsertID
		p.LinearFormula.LinearFormulaID = LinearFormulaID
	}

	// if PhysicalStateID = -1 then it is a new physical state
	if p.PhysicalState.PhysicalStateID != nil && err == nil && *p.PhysicalState.PhysicalStateID == -1 {
		logger.Log.Debug("new physical_state " + *p.PhysicalStateLabel)

		sqlr = `INSERT INTO physical_state (physical_state_label) VALUES (?)`

		if res, err = tx.Exec(sqlr, *p.PhysicalStateLabel); err != nil {
			return
		}
		// getting the last inserted id
		if lastInsertID, err = res.LastInsertId(); err != nil {
			return
		}
		// updating the product PhysicalStateID (PhysicalStateLabel already set)
		var PhysicalStateIDPointer *int64
		PhysicalStateIDPointer = new(int64)
		*PhysicalStateIDPointer = lastInsertID
		p.PhysicalState.PhysicalStateID = PhysicalStateIDPointer
	}

	// if CategoryID = -1 then it is a new category
	if p.Category.CategoryID != nil && err == nil && *p.Category.CategoryID == -1 {
		logger.Log.Debug("new category " + *p.CategoryLabel)

		sqlr = `INSERT INTO category (category_label) VALUES (?)`

		if res, err = tx.Exec(sqlr, *p.CategoryLabel); err != nil {
			return
		}
		// getting the last inserted id
		if lastInsertID, err = res.LastInsertId(); err != nil {
			return
		}
		// updating the product PhysicalStateID (PhysicalStateLabel already set)
		var CategoryIDPointer *int64
		CategoryIDPointer = new(int64)
		*CategoryIDPointer = lastInsertID
		p.Category.CategoryID = CategoryIDPointer
	}

	// if ProducerRefID = -1 then it is a new producer ref
	if p.ProducerRef.ProducerRefID != nil && err == nil && *p.ProducerRef.ProducerRefID == -1 {
		logger.Log.Debug("new producer_ref " + *p.ProducerRefLabel)

		sqlr = `INSERT INTO producer_ref (producer_ref_label, producer) VALUES (?, ?)`

		if res, err = tx.Exec(sqlr, p.ProducerRefLabel, p.Producer.ProducerID); err != nil {
			return
		}
		// getting the last inserted id
		if lastInsertID, err = res.LastInsertId(); err != nil {
			return
		}
		// updating the product ProducerRefID (ProducerRefLabel already set)
		var ProducerRefIDPointer *int64
		ProducerRefIDPointer = new(int64)
		*ProducerRefIDPointer = lastInsertID
		p.ProducerRef.ProducerRefID = ProducerRefIDPointer
	}

	// finally updating the product
	insertCols := goqu.Record{}

	if p.ProductInchi != nil {
		insertCols["product_inchi"] = *p.ProductInchi
	} else {
		insertCols["product_inchi"] = nil
	}

	if p.ProductInchikey != nil {
		insertCols["product_inchikey"] = *p.ProductInchikey
	} else {
		insertCols["product_inchikey"] = nil
	}

	if p.ProductCanonicalSmiles != nil {
		insertCols["product_canonical_smiles"] = *p.ProductCanonicalSmiles
	} else {
		insertCols["product_canonical_smiles"] = nil
	}

	if p.ProductMolecularWeight != nil {
		insertCols["product_molecular_weight"] = *p.ProductMolecularWeight
	} else {
		insertCols["product_molecular_weight"] = nil
	}

	if p.ProductSpecificity != nil {
		insertCols["product_specificity"] = *p.ProductSpecificity
	} else {
		insertCols["product_specificity"] = nil
	}

	if p.ProductMSDS != nil {
		insertCols["product_msds"] = *p.ProductMSDS
	} else {
		insertCols["product_msds"] = nil
	}

	if p.ProductSheet != nil {
		insertCols["product_sheet"] = *p.ProductSheet
	} else {
		insertCols["product_sheet"] = nil
	}

	if p.ProductTemperature != nil {
		insertCols["product_temperature"] = int(*p.ProductTemperature)
	} else {
		insertCols["product_temperature"] = nil
	}

	if p.ProductRestricted {
		insertCols["product_restricted"] = p.ProductRestricted
	} else {
		insertCols["product_restricted"] = false
	}

	if p.ProductRadioactive {
		insertCols["product_radioactive"] = p.ProductRadioactive
	} else {
		insertCols["product_radioactive"] = false
	}

	if p.Category.CategoryID != nil {
		insertCols["category"] = int(*p.Category.CategoryID)
	} else {
		insertCols["category"] = nil
	}

	if p.UnitTemperature.UnitID != nil {
		insertCols["unit_temperature"] = int(*p.UnitTemperature.UnitID)
	} else {
		insertCols["unit_temperature"] = nil
	}

	if p.UnitMolecularWeight.UnitID != nil {
		insertCols["unit_molecular_weight"] = int(*p.UnitMolecularWeight.UnitID)
	} else {
		insertCols["unit_molecular_weight"] = nil
	}

	if p.ProductThreeDFormula != nil {
		insertCols["product_threed_formula"] = *p.ProductThreeDFormula
	} else {
		insertCols["product_threed_formula"] = nil
	}

	if p.ProductTwoDFormula != nil {
		insertCols["product_twod_formula"] = *p.ProductTwoDFormula
	}
	// } else {
	// 	insertCols["product_twod_formula"] = nil
	// }

	if p.ProductDisposalComment != nil {
		insertCols["product_disposal_comment"] = *p.ProductDisposalComment
	} else {
		insertCols["product_disposal_comment"] = nil
	}

	if p.ProductRemark != nil {
		insertCols["product_remark"] = *p.ProductRemark
	} else {
		insertCols["product_remark"] = nil
	}

	if p.ProductNumberPerCarton != nil {
		insertCols["product_number_per_carton"] = *p.ProductNumberPerCarton
	} else {
		insertCols["product_number_per_carton"] = nil
	}

	if p.ProductNumberPerBag != nil {
		insertCols["product_number_per_bag"] = *p.ProductNumberPerBag
	} else {
		insertCols["product_number_per_bag"] = nil
	}

	if p.EmpiricalFormulaID != nil {
		insertCols["empirical_formula"] = *p.EmpiricalFormulaID
	} else {
		insertCols["empirical_formula"] = nil
	}

	if p.LinearFormulaID != nil {
		insertCols["linear_formula"] = *p.LinearFormulaID
	} else {
		insertCols["linear_formula"] = nil
	}

	if p.PhysicalStateID != nil {
		insertCols["physical_state"] = int(*p.PhysicalStateID)
	} else {
		insertCols["physical_state"] = nil
	}

	if p.SignalWordID != nil {
		insertCols["signal_word"] = int(*p.SignalWordID)
	} else {
		insertCols["signal_word"] = nil
	}

	// if p.CasNumberID!= nil {
	if p.CasNumberID != nil {
		// insertCols["cas_number"] = int(p.CasNumberID)
		insertCols["cas_number"] = int(*p.CasNumberID)
	} else {
		insertCols["cas_number"] = nil
	}

	if p.CeNumberID != nil {
		insertCols["ce_number"] = int(*p.CeNumberID)
	} else {
		insertCols["ce_number"] = nil
	}

	if p.ProducerRefID != nil {
		insertCols["producer_ref"] = int(*p.ProducerRefID)
	} else {
		insertCols["producer_ref"] = nil
	}

	// if p.ProductMolFormula!= nil {
	// 	insertCols["product_molformula"] = p.ProductMolFormula
	// } else {
	// 	insertCols["product_molformula"] = nil
	// }

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

	// adding supplier_refs
	if update {
		sqlr = `DELETE FROM productsupplierrefs WHERE productsupplierrefs.productsupplierrefs_product_id = (?)`
		if _, err = tx.Exec(sqlr, p.ProductID); err != nil {
			logger.Log.Error("error DELETE FROM productsupplierrefs")
			return
		}
	}

	for _, sr := range p.SupplierRefs {
		sqlr = `INSERT INTO productsupplierrefs (productsupplierrefs_product_id, productsupplierrefs_supplier_ref_id) VALUES (?,?)`
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
		sqlr = `DELETE FROM productclassesofcompounds WHERE productclassesofcompounds.productclassesofcompounds_product_id = (?)`
		if _, err = tx.Exec(sqlr, p.ProductID); err != nil {
			logger.Log.Error("error DELETE FROM productclassesofcompounds")
			return
		}
	}

	for _, coc := range p.ClassOfCompound {
		sqlr = `INSERT INTO productclassesofcompounds (productclassesofcompounds_product_id, productclassesofcompounds_class_of_compound_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, coc.ClassOfCompoundID); err != nil {
			logger.Log.Error("error INSERT INTO productclassesofcompounds")
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
		sqlr = `INSERT INTO producthazardstatements (producthazardstatements_product_id, producthazardstatements_hazard_statement_id) VALUES (?,?)`
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
		sqlr = `INSERT INTO productprecautionarystatements (productprecautionarystatements_product_id, productprecautionarystatements_precautionary_statement_id) VALUES (?,?)`
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

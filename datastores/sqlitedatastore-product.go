package datastores

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
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
		product.ProductType = "cons"
	case product.ProducerRef.ProducerRefID != nil:
		product.ProductType = "bio"
	default:
		product.ProductType = "chem"
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
	insertCols["product_type"] = p.ProductType

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

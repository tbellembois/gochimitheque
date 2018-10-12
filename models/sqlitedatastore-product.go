package models

import (
	"strings"

	"database/sql"

	"github.com/jmoiron/sqlx"
	"reflect"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/constants"
	"github.com/tbellembois/gochimitheque/helpers"
)

// GetProductsCasNumbers return the cas numbers matching the search criteria
func (db *SQLiteDataStore) GetProductsCasNumbers(p helpers.Dbselectparam) ([]CasNumber, int, error) {
	var (
		casnumbers                         []CasNumber
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	precreq.WriteString(" SELECT count(DISTINCT casnumber.casnumber_id)")
	presreq.WriteString(" SELECT casnumber_id, casnumber_label")

	comreq.WriteString(" FROM casnumber")
	comreq.WriteString(" WHERE casnumber_label LIKE :search")
	postsreq.WriteString(" ORDER BY casnumber_label  " + p.GetOrder())

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
	if err = snstmt.Select(&casnumbers, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	log.WithFields(log.Fields{"casnumbers": casnumbers}).Debug("GetProductsCasNumbers")
	return casnumbers, count, nil
}

// GetProductsCeNumbers return the cas numbers matching the search criteria
func (db *SQLiteDataStore) GetProductsCeNumbers(p helpers.Dbselectparam) ([]CeNumber, int, error) {
	var (
		cenumbers                          []CeNumber
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	precreq.WriteString(" SELECT count(DISTINCT cenumber.cenumber_id)")
	presreq.WriteString(" SELECT cenumber_id, cenumber_label")

	comreq.WriteString(" FROM cenumber")
	comreq.WriteString(" WHERE cenumber_label LIKE :search")
	postsreq.WriteString(" ORDER BY cenumber_label  " + p.GetOrder())

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
	if err = snstmt.Select(&cenumbers, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	log.WithFields(log.Fields{"cenumbers": cenumbers}).Debug("GetProductsCeNumbers")
	return cenumbers, count, nil
}

// GetProductsEmpiricalFormulas return the empirical formulas matching the search criteria
func (db *SQLiteDataStore) GetProductsEmpiricalFormulas(p helpers.Dbselectparam) ([]EmpiricalFormula, int, error) {
	var (
		eformulas                          []EmpiricalFormula
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	precreq.WriteString(" SELECT count(DISTINCT empiricalformula.empiricalformula_id)")
	presreq.WriteString(" SELECT empiricalformula_id, empiricalformula_label")

	comreq.WriteString(" FROM empiricalformula")
	comreq.WriteString(" WHERE empiricalformula_label LIKE :search")
	postsreq.WriteString(" ORDER BY empiricalformula_label  " + p.GetOrder())

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
	if err = snstmt.Select(&eformulas, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for formula matching exactly the search
	// TODO: search hardcoded
	if err = db.Select(&eformulas, `SELECT count(*) AS c, empiricalformula_id, empiricalformula_label FROM empiricalformula WHERE empiricalformula_label == ?`, "ClNa"); err != nil {
		return nil, 0, err
	}

	log.WithFields(log.Fields{"eformulas": eformulas}).Debug("GetProductsEmpiricalFormulas")
	return eformulas, count, nil
}

// GetProductsNames return the names matching the search criteria
func (db *SQLiteDataStore) GetProductsNames(p helpers.Dbselectparam) ([]Name, int, error) {
	var (
		names                              []Name
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	precreq.WriteString(" SELECT count(DISTINCT name.name_id)")
	presreq.WriteString(" SELECT name_id, name_label")

	comreq.WriteString(" FROM name")
	comreq.WriteString(" WHERE name_label LIKE :search")
	postsreq.WriteString(" ORDER BY name_label  " + p.GetOrder())

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
	if err = snstmt.Select(&names, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	log.WithFields(log.Fields{"names": names}).Debug("GetProductsNames")
	return names, count, nil
}

// GetProductsSymbols return the symbols matching the search criteria
func (db *SQLiteDataStore) GetProductsSymbols(p helpers.Dbselectparam) ([]Symbol, int, error) {
	var (
		symbols                            []Symbol
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	precreq.WriteString(" SELECT count(DISTINCT symbol.symbol_id)")
	presreq.WriteString(" SELECT symbol_id, symbol_label, symbol_image")

	comreq.WriteString(" FROM symbol")
	comreq.WriteString(" WHERE symbol_label LIKE :search")
	postsreq.WriteString(" ORDER BY symbol_label  " + p.GetOrder())

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
	if err = snstmt.Select(&symbols, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	log.WithFields(log.Fields{"symbols": symbols}).Debug("GetProductsSymbols")
	return symbols, count, nil
}

// GetProducts return the products matching the search criteria
func (db *SQLiteDataStore) GetProducts(p helpers.DbselectparamProduct) ([]Product, int, error) {
	var (
		products                                []Product
		count                                   int
		req, precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                                  *sqlx.NamedStmt
		snstmt                                  *sqlx.NamedStmt
		err                                     error
	)
	log.WithFields(log.Fields{"p": p}).Debug("GetProducts")

	// pre request: select or count
	precreq.WriteString(" SELECT count(DISTINCT product.product_id)")
	presreq.WriteString(` SELECT product.product_id, 
	product.product_specificity, 
	empiricalformula.empiricalformula_id AS "empiricalformula.empiricalformula_id",
	empiricalformula.empiricalformula_label AS "empiricalformula.empiricalformula_label",
	person.person_id AS "person.person_id",
	person.person_email AS "person.person_email",
	name.name_id AS "name.name_id",
	name.name_label AS "name.name_label",
	cenumber.cenumber_id AS "cenumber.cenumber_id",
	cenumber.cenumber_label AS "cenumber.cenumber_label",
	casnumber.casnumber_id AS "casnumber.casnumber_id",
	casnumber.casnumber_label AS "casnumber.casnumber_label"`)

	// common parts
	comreq.WriteString(" FROM product")
	// get name
	comreq.WriteString(" JOIN name ON product.name = name.name_id")
	// get casnumber
	comreq.WriteString(" JOIN casnumber ON product.casnumber = casnumber.casnumber_id")
	// get cenumber
	comreq.WriteString(" LEFT JOIN cenumber ON product.cenumber = cenumber.cenumber_id")
	// get person
	comreq.WriteString(" JOIN person ON product.person = person.person_id")
	// get empirical formula
	comreq.WriteString(" JOIN empiricalformula ON product.empiricalformula = empiricalformula.empiricalformula_id")
	// filter by permissions
	comreq.WriteString(` JOIN permission AS perm, entity as e ON
	(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
	(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "products" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
	(perm.person = :personid and perm.permission_item_name = "products" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "products" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "products" and perm.permission_perm_name = "r" and perm.permission_entity_id = e.entity_id)
	`)
	comreq.WriteString(" WHERE name.name_label LIKE :search")
	if p.GetProduct() != -1 {
		comreq.WriteString(" AND product.product_id = :product")
	}

	// post select request
	postsreq.WriteString(" GROUP BY product.product_id")
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
	if err = snstmt.Select(&products, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	//
	// getting symbols
	//
	for i, p := range products {
		// note: do not modify p but products[i] instead
		req.Reset()
		req.WriteString("SELECT symbol_id, symbol_label, symbol_image FROM symbol")
		req.WriteString(" JOIN productsymbols ON productsymbols.productsymbols_symbol_id = symbol.symbol_id")
		req.WriteString(" JOIN product ON productsymbols.productsymbols_product_id = product.product_id")
		req.WriteString(" WHERE product.product_id = ?")

		if err = db.Select(&products[i].Symbols, req.String(), p.ProductID); err != nil {
			return nil, 0, err
		}
	}

	//
	// getting synonyms
	//
	for i, p := range products {
		// note: do not modify p but products[i] instead
		req.Reset()
		req.WriteString("SELECT name_id, name_label FROM name")
		req.WriteString(" JOIN productsynonyms ON productsynonyms.productsynonyms_name_id = name.name_id")
		req.WriteString(" JOIN product ON productsynonyms.productsynonyms_product_id = product.product_id")
		req.WriteString(" WHERE product.product_id = ?")

		if err = db.Select(&products[i].Synonyms, req.String(), p.ProductID); err != nil {
			return nil, 0, err
		}
	}

	return products, count, nil
}

func (db *SQLiteDataStore) GetProduct(id int) (Product, error) {
	var (
		product Product
		sqlr    string
		err     error
	)

	sqlr = `SELECT product_id, 
	product_specificity, 
	empiricalformula.empiricalformula_id AS "empiricalformula.empiricalformula_id",
	empiricalformula.empiricalformula_label AS "empiricalformula.empiricalformula_label",
	person.person_id AS "person.person_id",
	person.person_email AS "person.person_email",
	name.name_id AS "name.name_id",
	name.name_label AS "name.name_label",
	cenumber.cenumber_id AS "cenumber.cenumber_id",
	cenumber.cenumber_label AS "cenumber.cenumber_label",
	casnumber.casnumber_id AS "casnumber.casnumber_id",
	casnumber.casnumber_label AS "casnumber.casnumber_label"
	FROM product
	JOIN name ON product.name = name.name_id
	JOIN casnumber ON product.casnumber = casnumber.casnumber_id
	LEFT JOIN cenumber ON product.cenumber = cenumber.cenumber_id
	JOIN person ON product.person = person.person_id
	JOIN empiricalformula ON product.empiricalformula = empiricalformula.empiricalformula_id
	WHERE product_id = ?`
	if err = db.Get(&product, sqlr, id); err != nil {
		return Product{}, err
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

	log.WithFields(log.Fields{"ID": id, "product": product}).Debug("GetProduct")
	return product, nil
}

func (db *SQLiteDataStore) DeleteProduct(id int) error {
	var (
		sqlr string
		err  error
	)
	// TODO: synonyms, symbols
	sqlr = `DELETE FROM product 
	WHERE product_id = ?`
	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}
	return nil
}

func (db *SQLiteDataStore) CreateProduct(p Product) (error, int) {
	var (
		lastid   int64
		tx       *sql.Tx
		sqlr     string
		res      sql.Result
		sqla     []interface{}
		ibuilder sq.InsertBuilder
		err      error
	)

	// beginning transaction
	if tx, err = db.Begin(); err != nil {
		return err, 0
	}

	// if CasNumberID = -1 then it is a new cas
	if p.CasNumber.CasNumberID == -1 {
		sqlr = `INSERT INTO casnumber (casnumber_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.CasNumberLabel); err != nil {
			tx.Rollback()
			return err, 0
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return err, 0
		}
		// updating the product CasNumberID (CasNumberLabel already set)
		p.CasNumber.CasNumberID = int(lastid)
	}
	// if CeNumberID = -1 then it is a new ce
	if v, err := p.CeNumber.CeNumberID.Value(); err == nil && v == -1 {
		sqlr = `INSERT INTO cenumber (cenumber_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.CeNumberLabel); err != nil {
			tx.Rollback()
			return err, 0
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return err, 0
		}
		// updating the product CeNumberID (CeNumberLabel already set)
		p.CeNumber.CeNumberID = sql.NullInt64{Int64: lastid}
	}
	// if NameID = -1 then it is a new name
	if p.Name.NameID == -1 {
		sqlr = `INSERT INTO name (name_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, strings.ToUpper(p.NameLabel)); err != nil {
			tx.Rollback()
			return err, 0
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return err, 0
		}
		// updating the product NameID (NameLabel already set)
		p.Name.NameID = int(lastid)
	}
	// if EmpiricalFormulaID = -1 then it is a new empirical formula
	if p.EmpiricalFormula.EmpiricalFormulaID == -1 {
		sqlr = `INSERT INTO empiricalformula (empiricalformula_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.EmpiricalFormulaLabel); err != nil {
			tx.Rollback()
			return err, 0
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return err, 0
		}
		// updating the product EmpiricalFormulaID (EmpiricalFormulaLabel already set)
		p.EmpiricalFormula.EmpiricalFormulaID = int(lastid)
	}

	// finally updating the product
	s := make(map[string]interface{})
	s["product_specificity"] = p.ProductSpecificity
	s["casnumber"] = p.CasNumberID
	s["name"] = p.NameID
	s["empiricalformula"] = p.EmpiricalFormulaID
	s["person"] = p.PersonID
	if p.CeNumber.CeNumberID.Valid {
		s["cenumber"] = int(p.CeNumberID.Int64)
	}
	// building column names/values
	col := make([]string, 0, len(s))
	val := make([]interface{}, 0, len(s))
	for k, v := range s {
		col = append(col, k)
		rt := reflect.TypeOf(v)
		rv := reflect.ValueOf(v)
		switch rt.Kind() {
		case reflect.Int:
			val = append(val, strconv.Itoa(int(rv.Int())))
		case reflect.String:
			val = append(val, rv.String())
		default:
			panic("unknown type")
		}
	}

	ibuilder = sq.Insert("product").Columns(col...).Values(val...)
	if sqlr, sqla, err = ibuilder.ToSql(); err != nil {
		tx.Rollback()
		return err, 0
	}

	if res, err = tx.Exec(sqlr, sqla...); err != nil {
		tx.Rollback()
		return err, 0
	}

	// getting the last inserted id
	if lastid, err = res.LastInsertId(); err != nil {
		tx.Rollback()
		return err, 0
	}
	p.ProductID = int(lastid)
	log.WithFields(log.Fields{"p": p}).Debug("CreateProduct")

	// adding symbols
	for _, sym := range p.Symbols {
		sqlr = `INSERT INTO productsymbols (productsymbols_product_id, productsymbols_symbol_id) VALUES (?,?)`
		if res, err = tx.Exec(sqlr, p.ProductID, sym.SymbolID); err != nil {
			tx.Rollback()
			return err, 0
		}
	}
	// adding synonyms
	for _, syn := range p.Synonyms {
		if syn.NameID == -1 {
			sqlr = `INSERT INTO name (name_label) VALUES (?)`
			if res, err = tx.Exec(sqlr, strings.ToUpper(syn.NameLabel)); err != nil {
				tx.Rollback()
				return err, 0
			}
			// getting the last inserted id
			if lastid, err = res.LastInsertId(); err != nil {
				tx.Rollback()
				return err, 0
			}
			syn.NameID = int(lastid)
		}
		sqlr = `INSERT INTO productsynonyms (productsynonyms_product_id, productsynonyms_name_id) VALUES (?,?)`
		if res, err = tx.Exec(sqlr, p.ProductID, syn.NameID); err != nil {
			tx.Rollback()
			return err, 0
		}
	}

	// committing changes
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err, 0
	}

	return nil, p.ProductID
}

func (db *SQLiteDataStore) UpdateProduct(p Product) error {
	var (
		lastid   int64
		tx       *sql.Tx
		sqlr     string
		res      sql.Result
		sqla     []interface{}
		ubuilder sq.UpdateBuilder
		err      error
	)

	// beginning transaction
	if tx, err = db.Begin(); err != nil {
		return err
	}

	// if CasNumberID = -1 then it is a new cas
	if p.CasNumber.CasNumberID == -1 {
		sqlr = `INSERT INTO casnumber (casnumber_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.CasNumberLabel); err != nil {
			tx.Rollback()
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return err
		}
		// updating the product CasNumberID (CasNumberLabel already set)
		p.CasNumber.CasNumberID = int(lastid)
	}
	// if CeNumberID = -1 then it is a new ce
	if p.CeNumber.CeNumberID.Int64 == -1 {
		sqlr = `INSERT INTO cenumber (cenumber_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.CeNumberLabel); err != nil {
			tx.Rollback()
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return err
		}
		// updating the product CeNumberID (CeNumberLabel already set)
		p.CeNumber.CeNumberID = sql.NullInt64{Int64: lastid, Valid: true}
	}
	// if NameID = -1 then it is a new name
	if p.Name.NameID == -1 {
		sqlr = `INSERT INTO name (name_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, strings.ToUpper(p.NameLabel)); err != nil {
			tx.Rollback()
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return err
		}
		// updating the product NameID (NameLabel already set)
		p.Name.NameID = int(lastid)
	}
	// if EmpiricalFormulaID = -1 then it is a new empirical formula
	if p.EmpiricalFormula.EmpiricalFormulaID == -1 {
		sqlr = `INSERT INTO empiricalformula (empiricalformula_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.EmpiricalFormulaLabel); err != nil {
			tx.Rollback()
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return err
		}
		// updating the product EmpiricalFormulaID (EmpiricalFormulaLabel already set)
		p.EmpiricalFormula.EmpiricalFormulaID = int(lastid)
	}

	// deleting symbols
	sqlr = `DELETE FROM productsymbols WHERE productsymbols.productsymbols_product_id = (?)`
	if res, err = tx.Exec(sqlr, p.ProductID); err != nil {
		tx.Rollback()
		return err
	}
	// adding new ones
	for _, sym := range p.Symbols {
		sqlr = `INSERT INTO productsymbols (productsymbols_product_id, productsymbols_symbol_id) VALUES (?,?)`
		if res, err = tx.Exec(sqlr, p.ProductID, sym.SymbolID); err != nil {
			tx.Rollback()
			return err
		}
	}

	// deleting synonyms
	sqlr = `DELETE FROM productsynonyms WHERE productsynonyms.productsynonyms_product_id = (?)`
	if res, err = tx.Exec(sqlr, p.ProductID); err != nil {
		tx.Rollback()
		return err
	}
	// adding new ones
	for _, syn := range p.Synonyms {
		if syn.NameID == -1 {
			sqlr = `INSERT INTO name (name_label) VALUES (?)`
			if res, err = tx.Exec(sqlr, strings.ToUpper(syn.NameLabel)); err != nil {
				tx.Rollback()
				return err
			}
			// getting the last inserted id
			if lastid, err = res.LastInsertId(); err != nil {
				tx.Rollback()
				return err
			}
			syn.NameID = int(lastid)
		}
		sqlr = `INSERT INTO productsynonyms (productsynonyms_product_id, productsynonyms_name_id) VALUES (?,?)`
		if res, err = tx.Exec(sqlr, p.ProductID, syn.NameID); err != nil {
			tx.Rollback()
			return err
		}
	}

	// finally updating the product
	s := make(map[string]interface{})
	s["product_specificity"] = p.ProductSpecificity
	s["casnumber"] = p.CasNumberID
	s["name"] = p.NameID
	s["empiricalformula"] = p.EmpiricalFormulaID
	s["person"] = p.PersonID
	if p.CeNumber.CeNumberID.Valid {
		s["cenumber"] = p.CeNumberID.Int64
	} else {
		s["cenumber"] = nil
	}

	ubuilder = sq.Update("product").
		SetMap(s).
		Where(sq.Eq{"product_id": p.ProductID})
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

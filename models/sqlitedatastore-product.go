package models

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/constants"
	"github.com/tbellembois/gochimitheque/helpers"
)

// IsProductBookmark returns true if there is a bookmark for the product pr for the person pe
func (db *SQLiteDataStore) IsProductBookmark(pr Product, pe Person) (bool, error) {
	var (
		sqlr string
		err  error
		i    int
	)
	sqlr = `SELECT count(*) FROM bookmark WHERE person = ? AND product = ?`
	if err = db.Get(&i, sqlr, pe.PersonID, pr.ProductID); err != nil {
		return false, err
	}
	return i != 0, err
}

// CreateProductBookmark bookmarks the product pr for the person pe
func (db *SQLiteDataStore) CreateProductBookmark(pr Product, pe Person) error {
	var (
		sqlr string
		err  error
	)
	sqlr = `INSERT into bookmark(person, product) VALUES (? , ?)`
	if _, err = db.Exec(sqlr, pe.PersonID, pr.ProductID); err != nil {
		return err
	}
	return nil
}

// DeleteProductBookmark remove the bookmark for the product pr and the person pe
func (db *SQLiteDataStore) DeleteProductBookmark(pr Product, pe Person) error {
	var (
		sqlr string
		err  error
	)
	sqlr = `DELETE from bookmark WHERE person = ? AND product = ?`
	if _, err = db.Exec(sqlr, pe.PersonID, pr.ProductID); err != nil {
		return err
	}
	return nil
}

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

	// setting the C attribute for formula matching exactly the search
	s := p.GetSearch()
	s = strings.TrimPrefix(s, "%")
	s = strings.TrimSuffix(s, "%")
	var casn CasNumber

	r := db.QueryRowx(`SELECT casnumber_id, casnumber_label FROM casnumber WHERE casnumber_label == ?`, s)
	if err = r.StructScan(&casn); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	for i, c := range casnumbers {
		if c.CasNumberID == casn.CasNumberID {
			casnumbers[i].C = 1
		}
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

	// setting the C attribute for formula matching exactly the search
	s := p.GetSearch()
	s = strings.TrimPrefix(s, "%")
	s = strings.TrimSuffix(s, "%")
	var cen CeNumber

	r := db.QueryRowx(`SELECT cenumber_id, cenumber_label FROM cenumber WHERE cenumber_label == ?`, s)
	if err = r.StructScan(&cen); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	for i, c := range cenumbers {
		if c.CeNumberID == cen.CeNumberID {
			cenumbers[i].C = 1
		}
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
	s := p.GetSearch()
	s = strings.TrimPrefix(s, "%")
	s = strings.TrimSuffix(s, "%")
	var ef EmpiricalFormula

	r := db.QueryRowx(`SELECT empiricalformula_id, empiricalformula_label FROM empiricalformula WHERE empiricalformula_label == ?`, s)
	if err = r.StructScan(&ef); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	for i, e := range eformulas {
		if e.EmpiricalFormulaID == ef.EmpiricalFormulaID {
			eformulas[i].C = 1
		}
	}

	log.WithFields(log.Fields{"eformulas": eformulas}).Debug("GetProductsEmpiricalFormulas")
	return eformulas, count, nil
}

// GetProductsLinearFormulas return the empirical formulas matching the search criteria
func (db *SQLiteDataStore) GetProductsLinearFormulas(p helpers.Dbselectparam) ([]LinearFormula, int, error) {
	var (
		lformulas                          []LinearFormula
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	precreq.WriteString(" SELECT count(DISTINCT linearformula.linearformula_id)")
	presreq.WriteString(" SELECT linearformula_id, linearformula_label")

	comreq.WriteString(" FROM linearformula")
	comreq.WriteString(" WHERE linearformula_label LIKE :search")
	postsreq.WriteString(" ORDER BY linearformula_label  " + p.GetOrder())

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
	if err = snstmt.Select(&lformulas, m); err != nil {
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
	var lf LinearFormula

	r := db.QueryRowx(`SELECT linearformula_id, linearformula_label FROM linearformula WHERE linearformula_label == ?`, s)
	if err = r.StructScan(&lf); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	for i, e := range lformulas {
		if e.LinearFormulaID == lf.LinearFormulaID {
			lformulas[i].C = 1
		}
	}

	log.WithFields(log.Fields{"lformulas": lformulas}).Debug("GetProductsLinearFormulas")
	return lformulas, count, nil
}

// GetProductsClassOfCompounds return the classe of compounds matching the search criteria
func (db *SQLiteDataStore) GetProductsClassOfCompounds(p helpers.Dbselectparam) ([]ClassOfCompound, int, error) {
	var (
		classofcompounds                   []ClassOfCompound
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	precreq.WriteString(" SELECT count(DISTINCT classofcompound.classofcompound_id)")
	presreq.WriteString(" SELECT classofcompound_id, classofcompound_label")

	comreq.WriteString(" FROM classofcompound")
	comreq.WriteString(" WHERE classofcompound_label LIKE :search")
	postsreq.WriteString(" ORDER BY classofcompound_label  " + p.GetOrder())

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
	if err = snstmt.Select(&classofcompounds, m); err != nil {
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
	var coc ClassOfCompound

	r := db.QueryRowx(`SELECT classofcompound_id, classofcompound_label FROM classofcompound WHERE classofcompound_label == ?`, s)
	if err = r.StructScan(&coc); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	for i, e := range classofcompounds {
		if e.ClassOfCompoundID == coc.ClassOfCompoundID {
			classofcompounds[i].C = 1
		}
	}

	log.WithFields(log.Fields{"classofcompounds": classofcompounds}).Debug("GetProductsClassOfCompounds")
	return classofcompounds, count, nil
}

// GetProductsName return the name matching the given id
func (db *SQLiteDataStore) GetProductsName(id int) (Name, error) {

	var (
		name Name
		sqlr string
		err  error
	)
	log.WithFields(log.Fields{"id": id}).Debug("GetProductsName")

	sqlr = `SELECT name.name_id, name.name_label
	FROM name
	WHERE name_id = ?`
	if err = db.Get(&name, sqlr, id); err != nil {
		return Name{}, err
	}
	log.WithFields(log.Fields{"ID": id, "name": name}).Debug("GetProductsName")
	return name, nil
}

// GetProductsEmpiricalFormula return the formula matching the given id
func (db *SQLiteDataStore) GetProductsEmpiricalFormula(id int) (EmpiricalFormula, error) {

	var (
		ef   EmpiricalFormula
		sqlr string
		err  error
	)
	log.WithFields(log.Fields{"id": id}).Debug("GetProductsEmpiricalFormula")

	sqlr = `SELECT empiricalformula.empiricalformula_id, empiricalformula.empiricalformula_label
	FROM empiricalformula
	WHERE empiricalformula_id = ?`
	if err = db.Get(&ef, sqlr, id); err != nil {
		return EmpiricalFormula{}, err
	}
	log.WithFields(log.Fields{"ID": id, "ef": ef}).Debug("GetProductsEmpiricalFormula")
	return ef, nil
}

// GetProductsCasNumber return the cas numbers matching the given id
func (db *SQLiteDataStore) GetProductsCasNumber(id int) (CasNumber, error) {

	var (
		cas  CasNumber
		sqlr string
		err  error
	)
	log.WithFields(log.Fields{"id": id}).Debug("GetProductsCasNumber")

	sqlr = `SELECT casnumber.casnumber_id, casnumber.casnumber_label
	FROM casnumber
	WHERE casnumber_id = ?`
	if err = db.Get(&cas, sqlr, id); err != nil {
		return CasNumber{}, err
	}
	log.WithFields(log.Fields{"ID": id, "cas": cas}).Debug("GetProductsCasNumber")
	return cas, nil
}

// GetProductsCasNumberByLabel return the cas numbers matching the given cas number
func (db *SQLiteDataStore) GetProductsCasNumberByLabel(label string) (CasNumber, error) {

	var (
		cas  CasNumber
		sqlr string
		err  error
	)
	log.WithFields(log.Fields{"label": label}).Debug("GetProductsCasNumberByLabel")

	sqlr = `SELECT casnumber.casnumber_id, casnumber.casnumber_label
	FROM casnumber
	WHERE casnumber_label = ?`
	if err = db.Get(&cas, sqlr, label); err != nil {
		return CasNumber{}, err
	}
	log.WithFields(log.Fields{"label": label, "cas": cas}).Debug("GetProductsCasNumberByLabel")
	return cas, nil
}

// GetProductsSignalWord return the signalword matching the given id
func (db *SQLiteDataStore) GetProductsSignalWord(id int) (SignalWord, error) {

	var (
		signalword SignalWord
		sqlr       string
		err        error
	)
	log.WithFields(log.Fields{"id": id}).Debug("GetProductsSignalWord")

	sqlr = `SELECT signalword.signalword_id, signalword.signalword_label
	FROM signalword
	WHERE signalword_id = ?`
	if err = db.Get(&signalword, sqlr, id); err != nil {
		return SignalWord{}, err
	}
	log.WithFields(log.Fields{"ID": id, "signalword": signalword}).Debug("GetProductsSignalWord")
	return signalword, nil
}

// GetProductsHazardStatement return the HazardStatement matching the given id
func (db *SQLiteDataStore) GetProductsHazardStatement(id int) (HazardStatement, error) {

	var (
		hs   HazardStatement
		sqlr string
		err  error
	)
	log.WithFields(log.Fields{"id": id}).Debug("GetProductsHazardStatement")

	sqlr = `SELECT hazardstatement.hazardstatement_id, hazardstatement.hazardstatement_label, hazardstatement.hazardstatement_reference
	FROM hazardstatement
	WHERE hazardstatement_id = ?`
	if err = db.Get(&hs, sqlr, id); err != nil {
		return HazardStatement{}, err
	}
	log.WithFields(log.Fields{"ID": id, "hs": hs}).Debug("GetProductsHazardStatement")
	return hs, nil
}

// GetProductsPrecautionaryStatement return the PrecautionaryStatement matching the given id
func (db *SQLiteDataStore) GetProductsPrecautionaryStatement(id int) (PrecautionaryStatement, error) {

	var (
		ps   PrecautionaryStatement
		sqlr string
		err  error
	)
	log.WithFields(log.Fields{"id": id}).Debug("GetProductsPrecautionaryStatement")

	sqlr = `SELECT precautionarystatement.precautionarystatement_id, precautionarystatement.precautionarystatement_label, precautionarystatement.precautionarystatement_reference
	FROM precautionarystatement
	WHERE precautionarystatement_id = ?`
	if err = db.Get(&ps, sqlr, id); err != nil {
		return PrecautionaryStatement{}, err
	}
	log.WithFields(log.Fields{"ID": id, "ps": ps}).Debug("GetProductsPrecautionaryStatement")
	return ps, nil
}

// GetProductsSymbol return the symbol matching the given id
func (db *SQLiteDataStore) GetProductsSymbol(id int) (Symbol, error) {

	var (
		symbol Symbol
		sqlr   string
		err    error
	)
	log.WithFields(log.Fields{"id": id}).Debug("GetProductsSymbol")

	sqlr = `SELECT symbol.symbol_id, symbol.symbol_label, symbol.symbol_image
	FROM symbol
	WHERE symbol_id = ?`
	if err = db.Get(&symbol, sqlr, id); err != nil {
		return Symbol{}, err
	}
	log.WithFields(log.Fields{"ID": id, "symbol": symbol}).Debug("GetProductsSymbol")
	return symbol, nil
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

	//select * from name WHERE name_label like "test%" OR name_label like "%test%" order by
	//case when name_label like "test%" then 0 else 1 end

	precreq.WriteString(" SELECT count(DISTINCT name.name_id)")
	presreq.WriteString(" SELECT name_id, name_label")

	comreq.WriteString(" FROM name")
	comreq.WriteString(" WHERE name_label LIKE :search OR name_label like :searchbegin")
	postsreq.WriteString(" ORDER BY case when name_label like :searchbegin then 0 else 1 end " + p.GetOrder())

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
		"search":      p.GetSearch(),
		"searchbegin": strings.TrimPrefix(p.GetSearch(), "%"),
		"order":       p.GetOrder(),
		"limit":       p.GetLimit(),
		"offset":      p.GetOffset(),
	}

	// select
	if err = snstmt.Select(&names, m); err != nil {
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
	var name Name

	r := db.QueryRowx(`SELECT name_id, name_label FROM name WHERE name_label == ?`, s)
	if err = r.StructScan(&name); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	for i, n := range names {
		if n.NameID == name.NameID {
			names[i].C = 1
		}
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

// GetProductsHazardStatementByReference return the hazard statement matching the reference
func (db *SQLiteDataStore) GetProductsHazardStatementByReference(r string) (HazardStatement, error) {
	var (
		err error
		hs  HazardStatement
	)

	sqlr := `SELECT hazardstatement_id, hazardstatement_label, hazardstatement_reference
	FROM hazardstatement
	WHERE hazardstatement_reference = ?`
	if err = db.Get(&hs, sqlr, r); err != nil {
		return HazardStatement{}, err
	}

	return hs, nil
}

// GetProductsHazardStatements return the hazard statements matching the search criteria
func (db *SQLiteDataStore) GetProductsHazardStatements(p helpers.Dbselectparam) ([]HazardStatement, int, error) {
	var (
		hazardstatements                   []HazardStatement
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	precreq.WriteString(" SELECT count(DISTINCT hazardstatement.hazardstatement_id)")
	presreq.WriteString(" SELECT hazardstatement_id, hazardstatement_label, hazardstatement_reference")

	comreq.WriteString(" FROM hazardstatement")
	comreq.WriteString(" WHERE hazardstatement_reference LIKE :search")
	postsreq.WriteString(" ORDER BY hazardstatement_label  " + p.GetOrder())

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
	if err = snstmt.Select(&hazardstatements, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	log.WithFields(log.Fields{"hazardstatements": hazardstatements}).Debug("GetProductsHazardStatements")
	return hazardstatements, count, nil
}

// GetProductsPrecautionaryStatementByReference return the precautionary statement matching the reference
func (db *SQLiteDataStore) GetProductsPrecautionaryStatementByReference(r string) (PrecautionaryStatement, error) {
	var (
		err error
		ps  PrecautionaryStatement
	)

	sqlr := `SELECT precautionarystatement_id, precautionarystatement_label, precautionarystatement_reference
	FROM precautionarystatement
	WHERE precautionarystatement_reference = ?`
	if err = db.Get(&ps, sqlr, r); err != nil {
		return PrecautionaryStatement{}, err
	}

	return ps, nil
}

// GetProductsPrecautionaryStatements return the hazard statements matching the search criteria
func (db *SQLiteDataStore) GetProductsPrecautionaryStatements(p helpers.Dbselectparam) ([]PrecautionaryStatement, int, error) {
	var (
		precautionarystatements            []PrecautionaryStatement
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	precreq.WriteString(" SELECT count(DISTINCT precautionarystatement.precautionarystatement_id)")
	presreq.WriteString(" SELECT precautionarystatement_id, precautionarystatement_label, precautionarystatement_reference")

	comreq.WriteString(" FROM precautionarystatement")
	comreq.WriteString(" WHERE precautionarystatement_reference LIKE :search")
	postsreq.WriteString(" ORDER BY precautionarystatement_label  " + p.GetOrder())

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
	if err = snstmt.Select(&precautionarystatements, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	log.WithFields(log.Fields{"precautionarystatements": precautionarystatements}).Debug("GetProductsPrecautionaryStatements")
	return precautionarystatements, count, nil
}

// GetProductsPhysicalStates return the physical states matching the search criteria
func (db *SQLiteDataStore) GetProductsPhysicalStates(p helpers.Dbselectparam) ([]PhysicalState, int, error) {
	var (
		physicalstates                     []PhysicalState
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	precreq.WriteString(" SELECT count(DISTINCT physicalstate.physicalstate_id)")
	presreq.WriteString(" SELECT physicalstate_id, physicalstate_label")

	comreq.WriteString(" FROM physicalstate")
	comreq.WriteString(" WHERE physicalstate_label LIKE :search")
	postsreq.WriteString(" ORDER BY physicalstate_label  " + p.GetOrder())

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
	if err = snstmt.Select(&physicalstates, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for physical states matching exactly the search
	s := p.GetSearch()
	s = strings.TrimPrefix(s, "%")
	s = strings.TrimSuffix(s, "%")
	var ps PhysicalState

	r := db.QueryRowx(`SELECT physicalstate_id, physicalstate_label FROM physicalstate WHERE physicalstate_label == ?`, s)
	if err = r.StructScan(&ps); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	for i, e := range physicalstates {
		if e.PhysicalStateID == ps.PhysicalStateID {
			physicalstates[i].C = 1
		}
	}

	log.WithFields(log.Fields{"physicalstates": physicalstates}).Debug("GetProductsPhysicalStates")
	return physicalstates, count, nil
}

// GetProductsSignalWords return the signal words matching the search criteria
func (db *SQLiteDataStore) GetProductsSignalWords(p helpers.Dbselectparam) ([]SignalWord, int, error) {
	var (
		signalwords                        []SignalWord
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	precreq.WriteString(" SELECT count(DISTINCT signalword.signalword_id)")
	presreq.WriteString(" SELECT signalword_id, signalword_label")

	comreq.WriteString(" FROM signalword")
	comreq.WriteString(" WHERE signalword_label LIKE :search")
	postsreq.WriteString(" ORDER BY signalword_label  " + p.GetOrder())

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
	if err = snstmt.Select(&signalwords, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	log.WithFields(log.Fields{"signalwords": signalwords}).Debug("GetProductsSignalWords")
	return signalwords, count, nil
}

// GetExposedProducts return all the products
func (db *SQLiteDataStore) GetExposedProducts() ([]Product, int, error) {
	var (
		products                                []Product
		count                                   int
		req, precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                                  *sqlx.NamedStmt
		snstmt                                  *sqlx.NamedStmt
		err                                     error
	)

	// pre request: select or count
	precreq.WriteString(" SELECT count(DISTINCT p.product_id)")
	presreq.WriteString(` SELECT p.product_id, 
	p.product_specificity, 
	p.product_msds,
	p.product_restricted,
	p.product_radioactive,
	p.product_threedformula,
	p.product_molformula,
	p.product_disposalcomment,
	p.product_remark,
	linearformula.linearformula_id AS "linearformula.linearformula_id",
	linearformula.linearformula_label AS "linearformula.linearformula_label",
	empiricalformula.empiricalformula_id AS "empiricalformula.empiricalformula_id",
	empiricalformula.empiricalformula_label AS "empiricalformula.empiricalformula_label",
	physicalstate.physicalstate_id AS "physicalstate.physicalstate_id",
	physicalstate.physicalstate_label AS "physicalstate.physicalstate_label",
	signalword.signalword_id AS "signalword.signalword_id",
	signalword.signalword_label AS "signalword.signalword_label",
	name.name_id AS "name.name_id",
	name.name_label AS "name.name_label",
	cenumber.cenumber_id AS "cenumber.cenumber_id",
	cenumber.cenumber_label AS "cenumber.cenumber_label",
	casnumber.casnumber_id AS "casnumber.casnumber_id",
	casnumber.casnumber_label AS "casnumber.casnumber_label",
	casnumber.casnumber_cmr AS "casnumber.casnumber_cmr"`)

	// common parts
	comreq.WriteString(" FROM product as p")
	// get name
	comreq.WriteString(" JOIN name ON p.name = name.name_id")
	// get casnumber
	comreq.WriteString(" JOIN casnumber ON p.casnumber = casnumber.casnumber_id")
	// get cenumber
	comreq.WriteString(" LEFT JOIN cenumber ON p.cenumber = cenumber.cenumber_id")
	// get physical state
	comreq.WriteString(" LEFT JOIN physicalstate ON p.physicalstate = physicalstate.physicalstate_id")
	// get signal word
	comreq.WriteString(" LEFT JOIN signalword ON p.signalword = signalword.signalword_id")
	// get empirical formula
	comreq.WriteString(" JOIN empiricalformula ON p.empiricalformula = empiricalformula.empiricalformula_id")
	// get linear formula
	comreq.WriteString(" LEFT JOIN linearformula ON p.linearformula = linearformula.linearformula_id")
	// get symbols
	comreq.WriteString(" JOIN productsymbols AS ps ON ps.productsymbols_product_id = p.product_id")

	// get hazardstatements
	comreq.WriteString(" JOIN producthazardstatements AS phs ON phs.producthazardstatements_product_id = p.product_id")
	// get precautionarystatements
	comreq.WriteString(" JOIN productprecautionarystatements AS pps ON pps.productprecautionarystatements_product_id = p.product_id")

	// post select request
	postsreq.WriteString(" GROUP BY p.product_id")

	// building count and select statements
	if cnstmt, err = db.PrepareNamed(precreq.String() + comreq.String()); err != nil {
		return nil, 0, err
	}
	if snstmt, err = db.PrepareNamed(presreq.String() + comreq.String() + postsreq.String()); err != nil {
		return nil, 0, err
	}

	m := map[string]interface{}{}
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
	for i, pr := range products {
		// note: do not modify p but products[i] instead
		req.Reset()
		req.WriteString("SELECT symbol_id, symbol_label, symbol_image FROM symbol")
		req.WriteString(" JOIN productsymbols ON productsymbols.productsymbols_symbol_id = symbol.symbol_id")
		req.WriteString(" JOIN product ON productsymbols.productsymbols_product_id = product.product_id")
		req.WriteString(" WHERE product.product_id = ?")

		if err = db.Select(&products[i].Symbols, req.String(), pr.ProductID); err != nil {
			return nil, 0, err
		}
	}

	//
	// getting classes of compounds
	//
	for i, pr := range products {
		// note: do not modify p but products[i] instead
		req.Reset()
		req.WriteString("SELECT classofcompound_id, classofcompound_label FROM classofcompound")
		req.WriteString(" JOIN productclassofcompound ON productclassofcompound.productclassofcompound_classofcompound_id = classofcompound.classofcompound_id")
		req.WriteString(" JOIN product ON productclassofcompound.productclassofcompound_product_id = product.product_id")
		req.WriteString(" WHERE product.product_id = ?")

		if err = db.Select(&products[i].ClassOfCompound, req.String(), pr.ProductID); err != nil {
			return nil, 0, err
		}
	}

	//
	// getting synonyms
	//
	for i, pr := range products {
		// note: do not modify p but products[i] instead
		req.Reset()
		req.WriteString("SELECT name_id, name_label FROM name")
		req.WriteString(" JOIN productsynonyms ON productsynonyms.productsynonyms_name_id = name.name_id")
		req.WriteString(" JOIN product ON productsynonyms.productsynonyms_product_id = product.product_id")
		req.WriteString(" WHERE product.product_id = ?")

		if err = db.Select(&products[i].Synonyms, req.String(), pr.ProductID); err != nil {
			return nil, 0, err
		}
	}

	//
	// getting hazard statements
	//
	for i, pr := range products {
		// note: do not modify p but products[i] instead
		req.Reset()
		req.WriteString("SELECT hazardstatement_id, hazardstatement_label, hazardstatement_reference FROM hazardstatement")
		req.WriteString(" JOIN producthazardstatements ON producthazardstatements.producthazardstatements_hazardstatement_id = hazardstatement.hazardstatement_id")
		req.WriteString(" JOIN product ON producthazardstatements.producthazardstatements_product_id = product.product_id")
		req.WriteString(" WHERE product.product_id = ?")

		if err = db.Select(&products[i].HazardStatements, req.String(), pr.ProductID); err != nil {
			return nil, 0, err
		}
	}

	//
	// getting precautionary statements
	//
	for i, pr := range products {
		// note: do not modify p but products[i] instead
		req.Reset()
		req.WriteString("SELECT precautionarystatement_id, precautionarystatement_label, precautionarystatement_reference FROM precautionarystatement")
		req.WriteString(" JOIN productprecautionarystatements ON productprecautionarystatements.productprecautionarystatements_precautionarystatement_id = precautionarystatement.precautionarystatement_id")
		req.WriteString(" JOIN product ON productprecautionarystatements.productprecautionarystatements_product_id = product.product_id")
		req.WriteString(" WHERE product.product_id = ?")

		if err = db.Select(&products[i].PrecautionaryStatements, req.String(), pr.ProductID); err != nil {
			return nil, 0, err
		}
	}

	return products, count, nil
}

// GetProducts return the products matching the search criteria
func (db *SQLiteDataStore) GetProducts(p helpers.DbselectparamProduct) ([]Product, int, error) {
	var (
		products                                               []Product
		count                                                  int
		reqsc, reqtsc, req, precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                                                 *sqlx.NamedStmt
		snstmt                                                 *sqlx.NamedStmt
		err                                                    error
		rperm                                                  bool
		isadmin                                                bool
	)
	log.WithFields(log.Fields{"p": p}).Debug("GetProducts")

	// is the user an admin?
	if isadmin, err = db.IsPersonAdmin(p.GetLoggedPersonID()); err != nil {
		return nil, 0, err
	}

	// shortcut
	if rperm, err = db.HasPersonPermission(p.GetLoggedPersonID(), "r", "rproducts", -1); err != nil {
		return nil, 0, err
	}
	log.WithFields(log.Fields{"rperm": rperm}).Debug("GetProducts")

	// pre request: select or count
	precreq.WriteString(" SELECT count(DISTINCT p.product_id)")
	presreq.WriteString(` SELECT p.product_id, 
	p.product_specificity, 
	p.product_msds,
	p.product_restricted,
	p.product_radioactive,
	p.product_threedformula,
	p.product_molformula,
	p.product_disposalcomment,
	p.product_remark,
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
	GROUP_CONCAT(DISTINCT storage.storage_barecode) AS "product_sl"`)

	// common parts
	comreq.WriteString(" FROM product as p")
	// get name
	comreq.WriteString(" JOIN name ON p.name = name.name_id")
	// get CMR
	if p.GetCasNumberCmr() {
		comreq.WriteString(" JOIN casnumber ON p.casnumber = casnumber.casnumber_id AND casnumber.casnumber_cmr IS NOT NULL")
	} else {
		// get casnumber
		comreq.WriteString(" JOIN casnumber ON p.casnumber = casnumber.casnumber_id")
	}
	// get cenumber
	comreq.WriteString(" LEFT JOIN cenumber ON p.cenumber = cenumber.cenumber_id")
	// get person
	comreq.WriteString(" JOIN person ON p.person = person.person_id")
	// get physical state
	comreq.WriteString(" LEFT JOIN physicalstate ON p.physicalstate = physicalstate.physicalstate_id")
	// get signal word
	comreq.WriteString(" LEFT JOIN signalword ON p.signalword = signalword.signalword_id")
	// get empirical formula
	comreq.WriteString(" JOIN empiricalformula ON p.empiricalformula = empiricalformula.empiricalformula_id")
	// get linear formula
	comreq.WriteString(" LEFT JOIN linearformula ON p.linearformula = linearformula.linearformula_id")
	// get bookmark
	comreq.WriteString(" LEFT JOIN bookmark ON (bookmark.product = p.product_id AND bookmark.person = :personid)")
	// get storages, store locations and entities
	comreq.WriteString(" LEFT JOIN storage ON storage.product = p.product_id")
	if p.GetEntity() != -1 || p.GetStorelocation() != -1 || p.GetStorageBarecode() != "" {
		//comreq.WriteString(" JOIN storage ON storage.product = p.product_id")
		comreq.WriteString(" JOIN storelocation ON storage.storelocation = storelocation.storelocation_id")
		comreq.WriteString(" JOIN entity ON storelocation.entity = entity.entity_id")
	}
	// get bookmarks
	if p.GetBookmark() {
		comreq.WriteString(" JOIN bookmark AS b ON b.product = p.product_id AND b.person = :personid")
	}
	// get symbols
	if len(p.GetSymbols()) != 0 {
		comreq.WriteString(" JOIN productsymbols AS ps ON ps.productsymbols_product_id = p.product_id")
	}
	// get hazardstatements
	if len(p.GetHazardStatements()) != 0 {
		comreq.WriteString(" JOIN producthazardstatements AS phs ON phs.producthazardstatements_product_id = p.product_id")
	}
	// get precautionarystatements
	if len(p.GetPrecautionaryStatements()) != 0 {
		comreq.WriteString(" JOIN productprecautionarystatements AS pps ON pps.productprecautionarystatements_product_id = p.product_id")
	}

	// filter by permissions
	// comreq.WriteString(` JOIN permission AS perm, entity as e ON
	// (perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
	// (perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
	// (perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
	// (perm.person = :personid and perm.permission_item_name = "products" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
	// (perm.person = :personid and perm.permission_item_name = "products" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
	// (perm.person = :personid and perm.permission_item_name = "products" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
	// (perm.person = :personid and perm.permission_item_name = "products" and perm.permission_perm_name = "r" and perm.permission_entity_id = e.entity_id)
	// `)
	comreq.WriteString(` JOIN permission AS perm, entity as e ON
	perm.person = :personid and (perm.permission_item_name in ("all", "products")) and (perm.permission_perm_name in ("all", "r")) and (perm.permission_entity_id in (-1, e.entity_id))
	`)
	//comreq.WriteString(" WHERE name.name_label LIKE :search")
	comreq.WriteString(" WHERE 1")
	if p.GetProduct() != -1 {
		comreq.WriteString(" AND p.product_id = :product")
	}
	if p.GetEntity() != -1 {
		comreq.WriteString(" AND entity.entity_id = :entity")
	}
	if p.GetStorelocation() != -1 {
		comreq.WriteString(" AND storelocation.storelocation_id = :storelocation")
	}
	if p.GetProductSpecificity() != "" {
		comreq.WriteString(" AND p.product_specificity = :product_specificity")
	}

	// search form parameters
	if p.GetName() != -1 {
		comreq.WriteString(" AND name.name_id = :name")
	}
	if p.GetCasNumber() != -1 {
		comreq.WriteString(" AND casnumber.casnumber_id = :casnumber")
	}
	if p.GetEmpiricalFormula() != -1 {
		comreq.WriteString(" AND empiricalformula.empiricalformula_id = :empiricalformula")
	}
	if p.GetStorageBarecode() != "" {
		comreq.WriteString(" AND storage.storage_barecode LIKE :storage_barecode")
	}
	if p.GetCustomNamePartOf() != "" {
		comreq.WriteString(" AND name.name_label LIKE :custom_name_part_of")
	}
	if len(p.GetSymbols()) != 0 {
		comreq.WriteString(" AND ps.productsymbols_symbol_id IN (")
		for _, s := range p.GetSymbols() {
			comreq.WriteString(fmt.Sprintf("%d,", s))
		}
		// to complete the last comma
		comreq.WriteString("-1")
		comreq.WriteString(" )")
	}
	if len(p.GetHazardStatements()) != 0 {
		comreq.WriteString(" AND phs.producthazardstatements_hazardstatement_id IN (")
		for _, s := range p.GetHazardStatements() {
			comreq.WriteString(fmt.Sprintf("%d,", s))
		}
		// to complete the last comma
		comreq.WriteString("-1")
		comreq.WriteString(" )")
	}
	if len(p.GetPrecautionaryStatements()) != 0 {
		comreq.WriteString(" AND pps.productprecautionarystatements_precautionarystatement_id IN (")
		for _, s := range p.GetPrecautionaryStatements() {
			comreq.WriteString(fmt.Sprintf("%d,", s))
		}
		// to complete the last comma
		comreq.WriteString("-1")
		comreq.WriteString(" )")
	}
	if p.GetSignalWord() != -1 {
		comreq.WriteString(" AND signalword.signalword_id = :signalword")
	}

	// filter restricted product
	if !rperm {
		comreq.WriteString(" AND p.product_restricted = false")
	}

	// post select request
	postsreq.WriteString(" GROUP BY p.product_id")
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
		"search":              p.GetSearch(),
		"personid":            p.GetLoggedPersonID(),
		"order":               p.GetOrder(),
		"limit":               p.GetLimit(),
		"offset":              p.GetOffset(),
		"entity":              p.GetEntity(),
		"product":             p.GetProduct(),
		"storelocation":       p.GetStorelocation(),
		"name":                p.GetName(),
		"casnumber":           p.GetCasNumber(),
		"empiricalformula":    p.GetEmpiricalFormula(),
		"product_specificity": p.GetProductSpecificity(),
		"storage_barecode":    p.GetStorageBarecode(),
		"custom_name_part_of": "%" + p.GetCustomNamePartOf() + "%",
		"signalword":          p.GetSignalWord(),
	}

	//log.Debug(presreq.String() + comreq.String() + postsreq.String())
	// log.Debug(m)

	// select
	if err = snstmt.Select(&products, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	//
	// cleaning product_sl
	//
	r := regexp.MustCompile("([a-zA-Z]{1}[0-9]+)\\.[0-9]+")
	for i, pr := range products {
		// note: do not modify p but products[i] instead
		m := r.FindAllStringSubmatch(pr.ProductSL.String, -1)
		// lazily adding only the first match
		if len(m) > 0 {
			products[i].ProductSL.String = m[0][1]
		} else {
			products[i].ProductSL.String = ""
		}
	}

	//
	// getting symbols
	//
	for i, pr := range products {
		// note: do not modify p but products[i] instead
		req.Reset()
		req.WriteString("SELECT symbol_id, symbol_label, symbol_image FROM symbol")
		req.WriteString(" JOIN productsymbols ON productsymbols.productsymbols_symbol_id = symbol.symbol_id")
		req.WriteString(" JOIN product ON productsymbols.productsymbols_product_id = product.product_id")
		req.WriteString(" WHERE product.product_id = ?")

		if err = db.Select(&products[i].Symbols, req.String(), pr.ProductID); err != nil {
			return nil, 0, err
		}
	}

	//
	// getting classes of compounds
	//
	for i, pr := range products {
		// note: do not modify p but products[i] instead
		req.Reset()
		req.WriteString("SELECT classofcompound_id, classofcompound_label FROM classofcompound")
		req.WriteString(" JOIN productclassofcompound ON productclassofcompound.productclassofcompound_classofcompound_id = classofcompound.classofcompound_id")
		req.WriteString(" JOIN product ON productclassofcompound.productclassofcompound_product_id = product.product_id")
		req.WriteString(" WHERE product.product_id = ?")

		if err = db.Select(&products[i].ClassOfCompound, req.String(), pr.ProductID); err != nil {
			return nil, 0, err
		}
	}

	//
	// getting synonyms
	//
	for i, pr := range products {
		// note: do not modify p but products[i] instead
		req.Reset()
		req.WriteString("SELECT name_id, name_label FROM name")
		req.WriteString(" JOIN productsynonyms ON productsynonyms.productsynonyms_name_id = name.name_id")
		req.WriteString(" JOIN product ON productsynonyms.productsynonyms_product_id = product.product_id")
		req.WriteString(" WHERE product.product_id = ?")

		if err = db.Select(&products[i].Synonyms, req.String(), pr.ProductID); err != nil {
			return nil, 0, err
		}
	}

	//
	// getting hazard statements
	//
	for i, pr := range products {
		// note: do not modify p but products[i] instead
		req.Reset()
		req.WriteString("SELECT hazardstatement_id, hazardstatement_label, hazardstatement_reference FROM hazardstatement")
		req.WriteString(" JOIN producthazardstatements ON producthazardstatements.producthazardstatements_hazardstatement_id = hazardstatement.hazardstatement_id")
		req.WriteString(" JOIN product ON producthazardstatements.producthazardstatements_product_id = product.product_id")
		req.WriteString(" WHERE product.product_id = ?")

		if err = db.Select(&products[i].HazardStatements, req.String(), pr.ProductID); err != nil {
			return nil, 0, err
		}
	}

	//
	// getting precautionary statements
	//
	for i, pr := range products {
		// note: do not modify p but products[i] instead
		req.Reset()
		req.WriteString("SELECT precautionarystatement_id, precautionarystatement_label, precautionarystatement_reference FROM precautionarystatement")
		req.WriteString(" JOIN productprecautionarystatements ON productprecautionarystatements.productprecautionarystatements_precautionarystatement_id = precautionarystatement.precautionarystatement_id")
		req.WriteString(" JOIN product ON productprecautionarystatements.productprecautionarystatements_product_id = product.product_id")
		req.WriteString(" WHERE product.product_id = ?")

		if err = db.Select(&products[i].PrecautionaryStatements, req.String(), pr.ProductID); err != nil {
			return nil, 0, err
		}
	}

	//
	// getting number of storages for each product
	//
	for i, pr := range products {
		// getting the total storage count
		reqtsc.Reset()
		reqtsc.WriteString("SELECT count(DISTINCT storage_id) from storage")
		reqtsc.WriteString(" JOIN product ON storage.product = ? AND storage.storage IS NULL AND storage.storage_archive == false")
		if isadmin {
			reqsc.Reset()
			reqsc.WriteString("SELECT count(DISTINCT storage_id) from storage")
			reqsc.WriteString(" JOIN product ON storage.product = ? AND storage.storage IS NULL AND storage.storage_archive == false")
		} else {
			// getting the storage count of the logged user entities
			reqsc.Reset()
			reqsc.WriteString("SELECT count(DISTINCT storage_id) from storage")
			reqsc.WriteString(" JOIN product ON storage.product = ? AND storage.storage IS NULL AND storage.storage_archive == false")
			reqsc.WriteString(" JOIN storelocation ON storage.storelocation = storelocation.storelocation_id")
			reqsc.WriteString(" JOIN entity ON storelocation.entity = entity.entity_id")
			reqsc.WriteString(" JOIN personentities ON (entity.entity_id = personentities.personentities_entity_id) AND")
			reqsc.WriteString(" (personentities.personentities_person_id = ?)")
		}
		if err = db.Get(&products[i].ProductSC, reqsc.String(), pr.ProductID, p.GetLoggedPersonID()); err != nil {
			return nil, 0, err
		}
		if err = db.Get(&products[i].ProductTSC, reqtsc.String(), pr.ProductID, p.GetLoggedPersonID()); err != nil {
			return nil, 0, err
		}
	}

	return products, count, nil
}

// CountProductStorages returns the number of storages for the product with the given id
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

	log.WithFields(log.Fields{"count": count}).Debug("CountProductStorages")
	return count, nil
}

// GetProduct returns the product with the given id
func (db *SQLiteDataStore) GetProduct(id int) (Product, error) {
	var (
		product Product
		sqlr    string
		err     error
	)

	sqlr = `SELECT product.product_id, 
	product.product_specificity, 
	product_msds,
	product_restricted,
	product_radioactive,
	product_threedformula,
	product_molformula,
	product_disposalcomment,
	product_remark,
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
	casnumber.casnumber_cmr AS "casnumber.casnumber_cmr"
	FROM product
	JOIN name ON product.name = name.name_id
	JOIN casnumber ON product.casnumber = casnumber.casnumber_id
	LEFT JOIN cenumber ON product.cenumber = cenumber.cenumber_id
	JOIN person ON product.person = person.person_id
	JOIN empiricalformula ON product.empiricalformula = empiricalformula.empiricalformula_id
	LEFT JOIN linearformula ON product.linearformula = linearformula.linearformula_id
	LEFT JOIN physicalstate ON product.physicalstate = physicalstate.physicalstate_id
	LEFT JOIN signalword ON product.signalword = signalword.signalword_id
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
	sqlr = `SELECT hazardstatement_id, hazardstatement_label, hazardstatement_reference FROM hazardstatement
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

	log.WithFields(log.Fields{"id": id, "product": product}).Debug("GetProduct")
	return product, nil
}

// DeleteProduct deletes the product with the given id
func (db *SQLiteDataStore) DeleteProduct(id int) error {
	var (
		sqlr string
		err  error
	)
	log.WithFields(log.Fields{"id": id}).Debug("DeleteProduct")
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

// CreateProduct insert the new product p into the database
func (db *SQLiteDataStore) CreateProduct(p Product) (int, error) {
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
		return 0, err
	}

	// if CasNumberID = -1 then it is a new cas
	if p.CasNumber.CasNumberID == -1 {
		sqlr = `INSERT INTO casnumber (casnumber_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.CasNumberLabel); err != nil {
			tx.Rollback()
			return 0, err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return 0, err
		}
		// updating the product CasNumberID (CasNumberLabel already set)
		p.CasNumber.CasNumberID = int(lastid)
	}
	// if CeNumberID = -1 then it is a new ce
	if v, err := p.CeNumber.CeNumberID.Value(); p.CeNumber.CeNumberID.Valid && err == nil && v.(int64) == -1 {
		sqlr = `INSERT INTO cenumber (cenumber_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.CeNumberLabel); err != nil {
			tx.Rollback()
			return 0, err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return 0, err
		}
		// updating the product CeNumberID (CeNumberLabel already set)
		p.CeNumber.CeNumberID = sql.NullInt64{Valid: true, Int64: lastid}
	}
	if err != nil {
		log.Error("cenumber error - " + err.Error())
		tx.Rollback()
		return 0, err
	}
	// if NameID = -1 then it is a new name
	if p.Name.NameID == -1 {
		sqlr = `INSERT INTO name (name_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, strings.ToUpper(p.NameLabel)); err != nil {
			tx.Rollback()
			return 0, err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return 0, err
		}
		// updating the product NameID (NameLabel already set)
		p.Name.NameID = int(lastid)
	}
	for i, syn := range p.Synonyms {
		if syn.NameID == -1 {
			sqlr = `INSERT INTO name (name_label) VALUES (?)`
			if res, err = tx.Exec(sqlr, strings.ToUpper(syn.NameLabel)); err != nil {
				tx.Rollback()
				return 0, err
			}
			// getting the last inserted id
			if lastid, err = res.LastInsertId(); err != nil {
				tx.Rollback()
				return 0, err
			}
			p.Synonyms[i].NameID = int(lastid)
		}
	}
	// if ClassOfCompoundID = -1 then it is a new class of compounds
	for i, coc := range p.ClassOfCompound {
		if coc.ClassOfCompoundID == -1 {
			sqlr = `INSERT INTO classofcompound (classofcompound_label) VALUES (?)`
			if res, err = tx.Exec(sqlr, strings.ToUpper(coc.ClassOfCompoundLabel)); err != nil {
				tx.Rollback()
				return 0, err
			}
			// getting the last inserted id
			if lastid, err = res.LastInsertId(); err != nil {
				tx.Rollback()
				return 0, err
			}
			p.ClassOfCompound[i].ClassOfCompoundID = int(lastid)
		}
	}
	// if EmpiricalFormulaID = -1 then it is a new empirical formula
	if p.EmpiricalFormula.EmpiricalFormulaID == -1 {
		sqlr = `INSERT INTO empiricalformula (empiricalformula_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.EmpiricalFormulaLabel); err != nil {
			tx.Rollback()
			return 0, err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return 0, err
		}
		// updating the product EmpiricalFormulaID (EmpiricalFormulaLabel already set)
		p.EmpiricalFormula.EmpiricalFormulaID = int(lastid)
	}
	// if LinearFormulaID = -1 then it is a new linear formula
	if v, err := p.LinearFormula.LinearFormulaID.Value(); p.LinearFormula.LinearFormulaID.Valid && err == nil && v.(int64) == -1 {
		sqlr = `INSERT INTO linearformula (linearformula_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.LinearFormulaLabel); err != nil {
			tx.Rollback()
			return 0, err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return 0, err
		}
		// updating the product LinearFormulaID (LinearFormulaLabel already set)
		p.LinearFormula.LinearFormulaID = sql.NullInt64{Valid: true, Int64: lastid}
	}
	// if PhysicalStateID = -1 then it is a new physical state
	if v, err := p.PhysicalState.PhysicalStateID.Value(); p.PhysicalState.PhysicalStateID.Valid && err == nil && v.(int64) == -1 {
		sqlr = `INSERT INTO physicalstate (physicalstate_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.PhysicalStateLabel); err != nil {
			tx.Rollback()
			return 0, err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return 0, err
		}
		// updating the product PhysicalStateID (PhysicalStateLabel already set)
		p.PhysicalState.PhysicalStateID = sql.NullInt64{Valid: true, Int64: lastid}
	}
	// if PhysicalStateID = -1 then it is a new physical state
	if v, err := p.PhysicalState.PhysicalStateID.Value(); p.PhysicalState.PhysicalStateID.Valid && err == nil && v.(int64) == -1 {
		sqlr = `INSERT INTO physicalstate (physicalstate_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.PhysicalStateLabel); err != nil {
			tx.Rollback()
			return 0, err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return 0, err
		}
		// updating the product ClassOfCompoundID (ClassOfCompoundLabel already set)
		p.PhysicalState.PhysicalStateID = sql.NullInt64{Valid: true, Int64: lastid}
	}
	if err != nil {
		log.Error("classofcompound error - " + err.Error())
		tx.Rollback()
		return 0, err
	}

	// finally updating the product
	s := make(map[string]interface{})
	if p.ProductSpecificity.Valid {
		s["product_specificity"] = p.ProductSpecificity.String
	}
	if p.ProductMSDS.Valid {
		s["product_msds"] = p.ProductMSDS.String
	}
	if p.ProductRestricted.Valid {
		s["product_restricted"] = p.ProductRestricted.Bool
	}
	if p.ProductRadioactive.Valid {
		s["product_radioactive"] = p.ProductRadioactive.Bool
	}

	if p.ProductThreeDFormula.Valid {
		s["product_threedformula"] = p.ProductThreeDFormula.String
	}
	if p.ProductDisposalComment.Valid {
		s["product_disposalcomment"] = p.ProductDisposalComment.String
	}
	if p.ProductRemark.Valid {
		s["product_remark"] = p.ProductRemark.String
	}
	if p.LinearFormulaID.Valid {
		s["linearformula"] = int(p.LinearFormulaID.Int64)
	}
	if p.PhysicalStateID.Valid {
		s["physicalstate"] = int(p.PhysicalStateID.Int64)
	}
	if p.SignalWordID.Valid {
		s["signalword"] = int(p.SignalWordID.Int64)
	}
	if p.CeNumberID.Valid {
		s["cenumber"] = int(p.CeNumberID.Int64)
	}
	if p.ProductMolFormula.Valid {
		s["product_molformula"] = p.ProductMolFormula.String
	}
	s["casnumber"] = p.CasNumberID
	s["name"] = p.NameID
	s["empiricalformula"] = p.EmpiricalFormulaID
	s["person"] = p.PersonID

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
		case reflect.Bool:
			val = append(val, rv.Bool())
		default:
			panic("unknown type:" + rt.String())
		}
	}

	ibuilder = sq.Insert("product").Columns(col...).Values(val...)
	if sqlr, sqla, err = ibuilder.ToSql(); err != nil {
		tx.Rollback()
		return 0, err
	}

	//log.Debug(ibuilder.ToSql())

	if res, err = tx.Exec(sqlr, sqla...); err != nil {
		log.Error("product error - " + err.Error())
		log.Error("sql:" + sqlr)
		tx.Rollback()
		return 0, err
	}

	// getting the last inserted id
	if lastid, err = res.LastInsertId(); err != nil {
		tx.Rollback()
		return 0, err
	}
	p.ProductID = int(lastid)
	log.WithFields(log.Fields{"p": p}).Debug("CreateProduct")

	// adding symbols
	for _, sym := range p.Symbols {
		sqlr = `INSERT INTO productsymbols (productsymbols_product_id, productsymbols_symbol_id) VALUES (?,?)`
		if res, err = tx.Exec(sqlr, p.ProductID, sym.SymbolID); err != nil {
			log.Error("productsymbols error - " + err.Error())
			tx.Rollback()
			return 0, err
		}
	}
	// adding classes of compounds
	for _, coc := range p.ClassOfCompound {
		sqlr = `INSERT INTO productclassofcompound (productclassofcompound_product_id, productclassofcompound_classofcompound_id) VALUES (?,?)`
		if res, err = tx.Exec(sqlr, p.ProductID, coc.ClassOfCompoundID); err != nil {
			log.Error("productclassofcompound error - " + err.Error())
			tx.Rollback()
			return 0, err
		}
	}
	// adding hazard statements
	for _, hs := range p.HazardStatements {
		sqlr = `INSERT INTO producthazardstatements (producthazardstatements_product_id, producthazardstatements_hazardstatement_id) VALUES (?,?)`
		if res, err = tx.Exec(sqlr, p.ProductID, hs.HazardStatementID); err != nil {
			log.Error("producthazardstatements error - " + err.Error())
			tx.Rollback()
			return 0, err
		}
	}
	// adding precautionary statements
	for _, ps := range p.PrecautionaryStatements {
		sqlr = `INSERT INTO productprecautionarystatements (productprecautionarystatements_product_id, productprecautionarystatements_precautionarystatement_id) VALUES (?,?)`
		if res, err = tx.Exec(sqlr, p.ProductID, ps.PrecautionaryStatementID); err != nil {
			log.Error("productprecautionarystatements error - " + err.Error())
			tx.Rollback()
			return 0, err
		}
	}
	// adding synonyms
	for _, syn := range p.Synonyms {
		sqlr = `INSERT INTO productsynonyms (productsynonyms_product_id, productsynonyms_name_id) VALUES (?,?)`
		if res, err = tx.Exec(sqlr, p.ProductID, syn.NameID); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// committing changes
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return 0, err
	}

	return p.ProductID, nil
}

// UpdateProduct updates the product p into the database
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
	if v, err := p.CeNumber.CeNumberID.Value(); p.CeNumber.CeNumberID.Valid && err == nil && v.(int64) == -1 {
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
		p.CeNumber.CeNumberID = sql.NullInt64{Valid: true, Int64: lastid}
	}
	if err != nil {
		log.Error("cenumber error - " + err.Error())
		tx.Rollback()
		return err
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
	for i, syn := range p.Synonyms {
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
			p.Synonyms[i].NameID = int(lastid)
		}
	}
	// if ClassOfCompoundID = -1 then it is a new class of compounds
	for i, coc := range p.ClassOfCompound {
		if coc.ClassOfCompoundID == -1 {
			sqlr = `INSERT INTO classofcompound (classofcompound_label) VALUES (?)`
			if res, err = tx.Exec(sqlr, strings.ToUpper(coc.ClassOfCompoundLabel)); err != nil {
				tx.Rollback()
				return err
			}
			// getting the last inserted id
			if lastid, err = res.LastInsertId(); err != nil {
				tx.Rollback()
				return err
			}
			p.ClassOfCompound[i].ClassOfCompoundID = int(lastid)
		}
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
	// if LinearFormulaID = -1 then it is a new linear formula
	if v, err := p.LinearFormula.LinearFormulaID.Value(); p.LinearFormula.LinearFormulaID.Valid && err == nil && v.(int64) == -1 {
		sqlr = `INSERT INTO linearformula (linearformula_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.LinearFormulaLabel); err != nil {
			tx.Rollback()
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return err
		}
		// updating the product LinearFormulaID (LinearFormulaLabel already set)
		p.LinearFormula.LinearFormulaID = sql.NullInt64{Valid: true, Int64: lastid}
	}
	// if PhysicalStateID = -1 then it is a new physical state
	if v, err := p.PhysicalState.PhysicalStateID.Value(); p.PhysicalState.PhysicalStateID.Valid && err == nil && v.(int64) == -1 {
		sqlr = `INSERT INTO physicalstate (physicalstate_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.PhysicalStateLabel); err != nil {
			tx.Rollback()
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			tx.Rollback()
			return err
		}
		// updating the product ClassOfCompoundID (ClassOfCompoundLabel already set)
		p.PhysicalState.PhysicalStateID = sql.NullInt64{Valid: true, Int64: lastid}
	}
	if err != nil {
		log.Error("classofcompound error - " + err.Error())
		tx.Rollback()
		return err
	}

	// finally updating the product
	s := make(map[string]interface{})
	if p.ProductSpecificity.Valid {
		s["product_specificity"] = p.ProductSpecificity.String
	}
	if p.ProductMSDS.Valid {
		s["product_msds"] = p.ProductMSDS.String
	}
	if p.ProductRestricted.Valid {
		s["product_restricted"] = p.ProductRestricted.Bool
	}
	if p.ProductRadioactive.Valid {
		s["product_radioactive"] = p.ProductRadioactive.Bool
	}
	if p.ProductThreeDFormula.Valid {
		s["product_threedformula"] = p.ProductThreeDFormula.String
	}
	if p.ProductDisposalComment.Valid {
		s["product_disposalcomment"] = p.ProductDisposalComment.String
	}
	if p.ProductRemark.Valid {
		s["product_remark"] = p.ProductRemark.String
	}
	if p.LinearFormulaID.Valid {
		s["linearformula"] = int(p.LinearFormulaID.Int64)
	}
	if p.PhysicalStateID.Valid {
		s["physicalstate"] = int(p.PhysicalStateID.Int64)
	}
	if p.SignalWordID.Valid {
		s["signalword"] = int(p.SignalWordID.Int64)
	}
	if p.CeNumberID.Valid {
		s["cenumber"] = int(p.CeNumberID.Int64)
	}
	if p.ProductMolFormula.Valid {
		s["product_molformula"] = p.ProductMolFormula.String
	}
	s["casnumber"] = p.CasNumberID
	s["name"] = p.NameID
	s["empiricalformula"] = p.EmpiricalFormulaID
	s["person"] = p.PersonID

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

	//log.Debug(ubuilder.ToSql())

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
		sqlr = `INSERT INTO productsynonyms (productsynonyms_product_id, productsynonyms_name_id) VALUES (?,?)`
		if res, err = tx.Exec(sqlr, p.ProductID, syn.NameID); err != nil {
			tx.Rollback()
			return err
		}
	}

	// deleting classes of compounds
	sqlr = `DELETE FROM productclassofcompound WHERE productclassofcompound.productclassofcompound_product_id = (?)`
	if res, err = tx.Exec(sqlr, p.ProductID); err != nil {
		tx.Rollback()
		return err
	}
	// adding new ones
	for _, coc := range p.ClassOfCompound {
		sqlr = `INSERT INTO productclassofcompound (productclassofcompound_product_id, productclassofcompound_classofcompound_id) VALUES (?,?)`
		if res, err = tx.Exec(sqlr, p.ProductID, coc.ClassOfCompoundID); err != nil {
			log.Error("productclassofcompound error - " + err.Error())
			tx.Rollback()
			return err
		}
	}

	// deleting hazard statements
	sqlr = `DELETE FROM producthazardstatements WHERE producthazardstatements.producthazardstatements_product_id = (?)`
	if res, err = tx.Exec(sqlr, p.ProductID); err != nil {
		tx.Rollback()
		return err
	}
	// adding new ones
	for _, hs := range p.HazardStatements {
		sqlr = `INSERT INTO producthazardstatements (producthazardstatements_product_id, producthazardstatements_hazardstatement_id) VALUES (?,?)`
		if res, err = tx.Exec(sqlr, p.ProductID, hs.HazardStatementID); err != nil {
			log.Error("producthazardstatements error - " + err.Error())
			tx.Rollback()
			return err
		}
	}

	// deleting precautionary statements
	sqlr = `DELETE FROM productprecautionarystatements WHERE productprecautionarystatements.productprecautionarystatements_product_id = (?)`
	if res, err = tx.Exec(sqlr, p.ProductID); err != nil {
		tx.Rollback()
		return err
	}
	// adding new ones
	for _, ps := range p.PrecautionaryStatements {
		sqlr = `INSERT INTO productprecautionarystatements (productprecautionarystatements_product_id, productprecautionarystatements_precautionarystatement_id) VALUES (?,?)`
		if res, err = tx.Exec(sqlr, p.ProductID, ps.PrecautionaryStatementID); err != nil {
			log.Error("productprecautionarystatements error - " + err.Error())
			tx.Rollback()
			return err
		}
	}

	// committing changes
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

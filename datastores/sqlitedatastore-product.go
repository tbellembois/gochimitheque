package datastores

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	. "github.com/tbellembois/gochimitheque/models"
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
func (db *SQLiteDataStore) GetProductsCasNumbers(p Dbselectparam) ([]CasNumber, int, error) {
	var (
		casnumbers                         []CasNumber
		count                              int
		exactSearch                        string
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	precreq.WriteString(" SELECT count(DISTINCT casnumber.casnumber_id)")
	presreq.WriteString(" SELECT casnumber_id, casnumber_label")

	comreq.WriteString(" FROM casnumber")
	comreq.WriteString(" WHERE casnumber_label LIKE :search")
	postsreq.WriteString(" ORDER BY INSTR(casnumber_label, \"" + exactSearch + "\") ASC, casnumber_label ASC")

	// limit
	if p.GetLimit() != ^uint64(0) {
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
	var casn CasNumber

	r := db.QueryRowx(`SELECT casnumber_id, casnumber_label 
	FROM casnumber 
	WHERE casnumber_label == ?`, exactSearch)
	if err = r.StructScan(&casn); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	for i, c := range casnumbers {
		if c.CasNumberID == casn.CasNumberID {
			casnumbers[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"casnumbers": casnumbers}).Debug("GetProductsCasNumbers")
	return casnumbers, count, nil
}

// GetProductsCeNumbers return the cas numbers matching the search criteria
func (db *SQLiteDataStore) GetProductsCeNumbers(p Dbselectparam) ([]CeNumber, int, error) {
	var (
		cenumbers                          []CeNumber
		count                              int
		exactSearch                        string
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	precreq.WriteString(" SELECT count(DISTINCT cenumber.cenumber_id)")
	presreq.WriteString(" SELECT cenumber_id, cenumber_label")

	comreq.WriteString(" FROM cenumber")
	comreq.WriteString(" WHERE cenumber_label LIKE :search")
	postsreq.WriteString(" ORDER BY INSTR(cenumber_label, \"" + exactSearch + "\") ASC, cenumber_label ASC")

	// limit
	if p.GetLimit() != ^uint64(0) {
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
	var cen CeNumber

	r := db.QueryRowx(`SELECT cenumber_id, cenumber_label FROM cenumber WHERE cenumber_label == ?`, exactSearch)
	if err = r.StructScan(&cen); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	for i, c := range cenumbers {
		if c.CeNumberID == cen.CeNumberID {
			cenumbers[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"cenumbers": cenumbers}).Debug("GetProductsCeNumbers")
	return cenumbers, count, nil
}

// GetProductsCeNumberByLabel return the ce numbers matching the given ce number
func (db *SQLiteDataStore) GetProductsCeNumberByLabel(label string) (CeNumber, error) {

	var (
		ce   CeNumber
		sqlr string
		err  error
	)
	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetProductsCeNumberByLabel")

	sqlr = `SELECT cenumber.cenumber_id, cenumber.cenumber_label
	FROM cenumber
	WHERE cenumber_label = ?
	ORDER BY cenumber.cenumber_label`
	if err = db.Get(&ce, sqlr, label); err != nil {
		return CeNumber{}, err
	}
	logger.Log.WithFields(logrus.Fields{"label": label, "ce": ce}).Debug("GetProductsCeNumberByLabel")
	return ce, nil
}

// GetProductsEmpiricalFormulas return the empirical formulas matching the search criteria
func (db *SQLiteDataStore) GetProductsEmpiricalFormulas(p Dbselectparam) ([]EmpiricalFormula, int, error) {
	var (
		eformulas                          []EmpiricalFormula
		count                              int
		exactSearch                        string
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	precreq.WriteString(" SELECT count(DISTINCT empiricalformula.empiricalformula_id)")
	presreq.WriteString(" SELECT empiricalformula_id, empiricalformula_label")

	comreq.WriteString(" FROM empiricalformula")
	comreq.WriteString(" WHERE empiricalformula_label LIKE :search")
	postsreq.WriteString(" ORDER BY INSTR(empiricalformula_label, \"" + exactSearch + "\") ASC, empiricalformula_label ASC")

	// limit
	if p.GetLimit() != ^uint64(0) {
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
	var ef EmpiricalFormula

	r := db.QueryRowx(`SELECT empiricalformula_id, empiricalformula_label FROM empiricalformula WHERE empiricalformula_label == ?`, exactSearch)
	if err = r.StructScan(&ef); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	for i, e := range eformulas {
		if e.EmpiricalFormulaID == ef.EmpiricalFormulaID {
			eformulas[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"eformulas": eformulas}).Debug("GetProductsEmpiricalFormulas")
	return eformulas, count, nil
}

// GetProductsEmpiricalFormulaByLabel return the empirirical formula matching the given empirical formula
func (db *SQLiteDataStore) GetProductsEmpiricalFormulaByLabel(label string) (EmpiricalFormula, error) {

	var (
		ef   EmpiricalFormula
		sqlr string
		err  error
	)
	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetProductsEmpiricalFormulaByLabel")

	sqlr = `SELECT empiricalformula.empiricalformula_id, empiricalformula.empiricalformula_label
	FROM empiricalformula
	WHERE empiricalformula.empiricalformula_label = ?
	ORDER BY empiricalformula.empiricalformula_label`
	if err = db.Get(&ef, sqlr, label); err != nil {
		return EmpiricalFormula{}, err
	}
	logger.Log.WithFields(logrus.Fields{"label": label, "ef": ef}).Debug("GetProductsEmpiricalFormulaByLabel")
	return ef, nil
}

// GetProductsLinearFormulas return the empirical formulas matching the search criteria
func (db *SQLiteDataStore) GetProductsLinearFormulas(p Dbselectparam) ([]LinearFormula, int, error) {
	var (
		lformulas                          []LinearFormula
		count                              int
		exactSearch                        string
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	precreq.WriteString(" SELECT count(DISTINCT linearformula.linearformula_id)")
	presreq.WriteString(" SELECT linearformula_id, linearformula_label")

	comreq.WriteString(" FROM linearformula")
	comreq.WriteString(" WHERE linearformula_label LIKE :search")
	postsreq.WriteString(" ORDER BY INSTR(linearformula_label, \"" + exactSearch + "\") ASC, linearformula_label ASC")

	// limit
	if p.GetLimit() != ^uint64(0) {
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
	var lf LinearFormula

	r := db.QueryRowx(`SELECT linearformula_id, linearformula_label FROM linearformula WHERE linearformula_label == ?`, exactSearch)
	if err = r.StructScan(&lf); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	for i, e := range lformulas {
		if e.LinearFormulaID == lf.LinearFormulaID {
			lformulas[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"lformulas": lformulas}).Debug("GetProductsLinearFormulas")
	return lformulas, count, nil
}

// GetProductsLinearFormulaByLabel return the linear formula matching the given linear formula
func (db *SQLiteDataStore) GetProductsLinearFormulaByLabel(label string) (LinearFormula, error) {

	var (
		lf   LinearFormula
		sqlr string
		err  error
	)
	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetProductsLinearFormulaByLabel")

	sqlr = `SELECT linearformula.linearformula_id, linearformula.linearformula_label
	FROM linearformula
	WHERE linearformula.linearformula_label = ?
	ORDER BY linearformula.linearformula_label`
	if err = db.Get(&lf, sqlr, label); err != nil {
		return LinearFormula{}, err
	}
	logger.Log.WithFields(logrus.Fields{"label": label, "lf": lf}).Debug("GetProductsLinearFormulaByLabel")
	return lf, nil
}

// GetProductsClassOfCompounds return the classe of compounds matching the search criteria
func (db *SQLiteDataStore) GetProductsClassOfCompounds(p Dbselectparam) ([]ClassOfCompound, int, error) {
	var (
		classofcompounds                   []ClassOfCompound
		count                              int
		exactSearch                        string
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	precreq.WriteString(" SELECT count(DISTINCT classofcompound.classofcompound_id)")
	presreq.WriteString(" SELECT classofcompound_id, classofcompound_label")

	comreq.WriteString(" FROM classofcompound")
	comreq.WriteString(" WHERE classofcompound_label LIKE :search")
	postsreq.WriteString(" ORDER BY INSTR(classofcompound_label, \"" + exactSearch + "\") ASC, classofcompound_label ASC")

	// limit
	if p.GetLimit() != ^uint64(0) {
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
	var coc ClassOfCompound

	r := db.QueryRowx(`SELECT classofcompound_id, classofcompound_label FROM classofcompound WHERE classofcompound_label == ?`, exactSearch)
	if err = r.StructScan(&coc); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	for i, e := range classofcompounds {
		if e.ClassOfCompoundID == coc.ClassOfCompoundID {
			classofcompounds[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"classofcompounds": classofcompounds}).Debug("GetProductsClassOfCompounds")
	return classofcompounds, count, nil
}

// GetProductsClassOfCompoundByLabel return the class of compounds matching the given label
func (db *SQLiteDataStore) GetProductsClassOfCompoundByLabel(label string) (ClassOfCompound, error) {

	var (
		coc  ClassOfCompound
		sqlr string
		err  error
	)
	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetProductsClassOfCompoundByLabel")

	sqlr = `SELECT classofcompound.classofcompound_id, classofcompound.classofcompound_label
	FROM classofcompound
	WHERE classofcompound.classofcompound_label = ?
	ORDER BY classofcompound.classofcompound_label`
	if err = db.Get(&coc, sqlr, strings.ToUpper(label)); err != nil {
		return ClassOfCompound{}, err
	}
	logger.Log.WithFields(logrus.Fields{"label": label, "coc": coc}).Debug("GetProductsClassOfCompoundByLabel")
	return coc, nil
}

// GetProductsSuppliers return the suppliers matching the search criteria
func (db *SQLiteDataStore) GetProductsSuppliers(p Dbselectparam) ([]Supplier, int, error) {
	var (
		srs                                []Supplier
		count                              int
		exactSearch                        string
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	precreq.WriteString(" SELECT count(DISTINCT supplier.supplier_id)")
	presreq.WriteString(` SELECT supplier_id, 
								 supplier_label`)

	comreq.WriteString(" FROM supplier")
	comreq.WriteString(" WHERE supplier_label LIKE :search")
	postsreq.WriteString(" ORDER BY INSTR(supplier_label, \"" + exactSearch + "\") ASC, supplier_label ASC")

	// limit
	if p.GetLimit() != ^uint64(0) {
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
		"limit":  p.GetLimit(),
		"offset": p.GetOffset(),
	}

	// select
	if err = snstmt.Select(&srs, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for formula matching exactly the search
	var sp Supplier

	r := db.QueryRowx(`SELECT supplier_id, supplier_label FROM supplier WHERE supplier_label == ?`, exactSearch)
	if err = r.StructScan(&sp); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	for i, s := range srs {
		if s.SupplierID == sp.SupplierID {
			srs[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"srs": srs}).Debug("GetProductsSuppliers")
	return srs, count, nil
}

// GetProductsProducers return the producers matching the search criteria
func (db *SQLiteDataStore) GetProductsProducers(p Dbselectparam) ([]Producer, int, error) {
	var (
		prs                                []Producer
		count                              int
		exactSearch                        string
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	precreq.WriteString(" SELECT count(DISTINCT producer.producer_id)")
	presreq.WriteString(` SELECT producer_id, 
								 producer_label`)

	comreq.WriteString(" FROM producer")
	comreq.WriteString(" WHERE producer_label LIKE :search")
	postsreq.WriteString(" ORDER BY INSTR(producer_label, \"" + exactSearch + "\") ASC, producer_label ASC")

	// limit
	if p.GetLimit() != ^uint64(0) {
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
		"limit":  p.GetLimit(),
		"offset": p.GetOffset(),
	}

	// select
	if err = snstmt.Select(&prs, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for formula matching exactly the search
	var pr Producer

	r := db.QueryRowx(`SELECT producer_id, producer_label FROM producer WHERE producer_label == ?`, exactSearch)
	if err = r.StructScan(&pr); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	for i, p := range prs {
		if p.ProducerID == pr.ProducerID {
			prs[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"prs": prs}).Debug("GetProductsProducers")
	return prs, count, nil
}

// GetProductsProducerRefs return the producerrefs matching the search criteria
func (db *SQLiteDataStore) GetProductsProducerRefs(p DbselectparamProducerRef) ([]ProducerRef, int, error) {
	var (
		prefs                              []ProducerRef
		count                              int
		exactSearch                        string
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	precreq.WriteString(" SELECT count(DISTINCT producerref.producerref_id)")
	presreq.WriteString(` SELECT producerref_id, 
								 producerref_label, 
								 producer_id AS "producer.producer_id",
								 producer_label AS "producer.producer_label"`)

	comreq.WriteString(" FROM producerref")
	comreq.WriteString(" JOIN producer ON producerref.producer = producer.producer_id")
	comreq.WriteString(" WHERE producerref.producerref_label LIKE :search")
	if p.GetProducer() != -1 {
		comreq.WriteString(" AND producerref.producer = :producer")
	}
	postsreq.WriteString(" ORDER BY INSTR(producerref_label, \"" + exactSearch + "\") ASC, producerref_label ASC")

	// limit
	if p.GetLimit() != ^uint64(0) {
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
		"limit":    p.GetLimit(),
		"offset":   p.GetOffset(),
		"producer": p.GetProducer(),
	}

	// logger.Log.Debug(presreq.String() + comreq.String() + postsreq.String())
	// logger.Log.Debug(m)

	// select
	if err = snstmt.Select(&prefs, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for formula matching exactly the search
	var pref ProducerRef

	r := db.QueryRowx(`SELECT producerref_id, producerref_label FROM producerref WHERE producerref_label == ?`, exactSearch)
	if err = r.StructScan(&pref); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	for i, p := range prefs {
		if p.ProducerRefID == pref.ProducerRefID {
			prefs[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"prefs": prefs}).Debug("GetProductsProducerRefs")
	return prefs, count, nil
}

// GetProductsSupplierRefs return the supplierrefs matching the search criteria
func (db *SQLiteDataStore) GetProductsSupplierRefs(p DbselectparamSupplierRef) ([]SupplierRef, int, error) {
	var (
		srefs                              []SupplierRef
		count                              int
		exactSearch                        string
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	precreq.WriteString(" SELECT count(DISTINCT supplierref.supplierref_id)")
	presreq.WriteString(` SELECT supplierref_id, 
								 supplierref_label, 
								 supplier_id AS "supplier.supplier_id",
								 supplier_label AS "supplier.supplier_label"`)

	comreq.WriteString(" FROM supplierref")
	comreq.WriteString(" JOIN supplier ON supplierref.supplier = supplier.supplier_id")
	comreq.WriteString(" WHERE supplierref.supplierref_label LIKE :search")
	if p.GetSupplier() != -1 {
		comreq.WriteString(" AND supplierref.supplier = :supplier")
	}
	postsreq.WriteString(" ORDER BY INSTR(supplierref_label, \"" + exactSearch + "\") ASC, supplierref_label ASC")

	// limit
	if p.GetLimit() != ^uint64(0) {
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
		"limit":    p.GetLimit(),
		"offset":   p.GetOffset(),
		"supplier": p.GetSupplier(),
	}

	// select
	if err = snstmt.Select(&srefs, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for formula matching exactly the search
	var sref SupplierRef

	r := db.QueryRowx(`SELECT supplierref_id, supplierref_label FROM supplierref WHERE supplierref_label == ?`, exactSearch)
	if err = r.StructScan(&sref); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	for i, s := range srefs {
		if s.SupplierRefID == sref.SupplierRefID {
			srefs[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"srefs": srefs}).Debug("GetProductsSupplierRefs")
	return srefs, count, nil
}

// GetProductsTags return the tags matching the search criteria
func (db *SQLiteDataStore) GetProductsTags(p Dbselectparam) ([]Tag, int, error) {
	var (
		tags                               []Tag
		count                              int
		exactSearch                        string
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	precreq.WriteString(" SELECT count(DISTINCT tag.tag_id)")
	presreq.WriteString(" SELECT tag_id, tag_label")

	comreq.WriteString(" FROM tag")
	comreq.WriteString(" WHERE tag_label LIKE :search")
	postsreq.WriteString(" ORDER BY INSTR(tag_label, \"" + exactSearch + "\") ASC, tag_label ASC")

	// limit
	if p.GetLimit() != ^uint64(0) {
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
		"limit":  p.GetLimit(),
		"offset": p.GetOffset(),
	}

	// select
	if err = snstmt.Select(&tags, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for formula matching exactly the search
	var tag Tag

	r := db.QueryRowx(`SELECT tag_id, tag_label FROM tag WHERE tag_label == ?`, exactSearch)
	if err = r.StructScan(&tag); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	for i, t := range tags {
		if t.TagID == tag.TagID {
			tags[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"tags": tags}).Debug("GetProductsTags")
	return tags, count, nil
}

// GetProductsCategories return the categories matching the search criteria
func (db *SQLiteDataStore) GetProductsCategories(p Dbselectparam) ([]Category, int, error) {
	var (
		categories                         []Category
		count                              int
		exactSearch                        string
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	precreq.WriteString(" SELECT count(DISTINCT category.category_id)")
	presreq.WriteString(" SELECT category_id, category_label")

	comreq.WriteString(" FROM category")
	comreq.WriteString(" WHERE category_label LIKE :search")
	postsreq.WriteString(" ORDER BY INSTR(category_label, \"" + exactSearch + "\") ASC, category_label ASC")

	// limit
	if p.GetLimit() != ^uint64(0) {
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
		"limit":  p.GetLimit(),
		"offset": p.GetOffset(),
	}

	// select
	if err = snstmt.Select(&categories, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for formula matching exactly the search
	var category Category

	r := db.QueryRowx(`SELECT category_id, category_label FROM category WHERE category_label == ?`, exactSearch)
	if err = r.StructScan(&category); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	for i, c := range categories {
		if c.CategoryID == category.CategoryID {
			categories[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"categories": categories}).Debug("GetProductsCategories")
	return categories, count, nil
}

// GetProductsName return the name matching the given id
func (db *SQLiteDataStore) GetProductsName(id int) (Name, error) {

	var (
		name Name
		sqlr string
		err  error
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetProductsName")

	sqlr = `SELECT name.name_id, name.name_label
	FROM name
	WHERE name_id = ?`
	if err = db.Get(&name, sqlr, id); err != nil {
		return Name{}, err
	}
	logger.Log.WithFields(logrus.Fields{"ID": id, "name": name}).Debug("GetProductsName")
	return name, nil
}

// GetProductsNameByLabel return the name matching the given label
func (db *SQLiteDataStore) GetProductsNameByLabel(label string) (Name, error) {

	var (
		name Name
		sqlr string
		err  error
	)
	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetProductsNameByLabel")

	sqlr = `SELECT name.name_id, name.name_label
	FROM name
	WHERE name.name_label = ?`
	if err = db.Get(&name, sqlr, strings.ToUpper(label)); err != nil {
		return Name{}, err
	}
	logger.Log.WithFields(logrus.Fields{"label": label, "name": name}).Debug("GetProductsNameByLabel")
	return name, nil
}

// GetProductsEmpiricalFormula return the formula matching the given id
func (db *SQLiteDataStore) GetProductsEmpiricalFormula(id int) (EmpiricalFormula, error) {

	var (
		ef   EmpiricalFormula
		sqlr string
		err  error
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetProductsEmpiricalFormula")

	sqlr = `SELECT empiricalformula.empiricalformula_id, empiricalformula.empiricalformula_label
	FROM empiricalformula
	WHERE empiricalformula_id = ?
	ORDER BY empiricalformula.empiricalformula_label`
	if err = db.Get(&ef, sqlr, id); err != nil {
		return EmpiricalFormula{}, err
	}
	logger.Log.WithFields(logrus.Fields{"ID": id, "ef": ef}).Debug("GetProductsEmpiricalFormula")
	return ef, nil
}

// GetProductsCasNumber return the cas numbers matching the given id
func (db *SQLiteDataStore) GetProductsCasNumber(id int) (CasNumber, error) {

	var (
		cas  CasNumber
		sqlr string
		err  error
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetProductsCasNumber")

	sqlr = `SELECT casnumber.casnumber_id, casnumber.casnumber_label
	FROM casnumber
	WHERE casnumber_id = ?`
	if err = db.Get(&cas, sqlr, id); err != nil {
		return CasNumber{}, err
	}
	logger.Log.WithFields(logrus.Fields{"ID": id, "cas": cas}).Debug("GetProductsCasNumber")
	return cas, nil
}

// GetProductsCasNumberByLabel return the cas numbers matching the given cas number
func (db *SQLiteDataStore) GetProductsCasNumberByLabel(label string) (CasNumber, error) {

	var (
		cas  CasNumber
		sqlr string
		err  error
	)
	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetProductsCasNumberByLabel")

	sqlr = `SELECT casnumber.casnumber_id, casnumber.casnumber_label
	FROM casnumber
	WHERE casnumber_label = ?`
	if err = db.Get(&cas, sqlr, label); err != nil {
		return CasNumber{}, err
	}
	logger.Log.WithFields(logrus.Fields{"label": label, "cas": cas}).Debug("GetProductsCasNumberByLabel")
	return cas, nil
}

// GetProductsSignalWord return the signalword matching the given id
func (db *SQLiteDataStore) GetProductsSignalWord(id int) (SignalWord, error) {

	var (
		signalword SignalWord
		sqlr       string
		err        error
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetProductsSignalWord")

	sqlr = `SELECT signalword.signalword_id, signalword.signalword_label
	FROM signalword
	WHERE signalword_id = ?`
	if err = db.Get(&signalword, sqlr, id); err != nil {
		return SignalWord{}, err
	}
	logger.Log.WithFields(logrus.Fields{"ID": id, "signalword": signalword}).Debug("GetProductsSignalWord")
	return signalword, nil
}

// GetProductsSignalWordByLabel return the signal word matching the given label
func (db *SQLiteDataStore) GetProductsSignalWordByLabel(label string) (SignalWord, error) {

	var (
		sw   SignalWord
		sqlr string
		err  error
	)
	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetProductsSignalWordByLabel")

	sqlr = `SELECT signalword.signalword_id, signalword.signalword_label
	FROM signalword
	WHERE signalword.signalword_label = ?`
	if err = db.Get(&sw, sqlr, label); err != nil {
		return SignalWord{}, err
	}
	logger.Log.WithFields(logrus.Fields{"label": label, "sw": sw}).Debug("GetProductsSignalWordByLabel")
	return sw, nil
}

// GetProductsHazardStatement return the HazardStatement matching the given id
func (db *SQLiteDataStore) GetProductsHazardStatement(id int) (HazardStatement, error) {

	var (
		hs   HazardStatement
		sqlr string
		err  error
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetProductsHazardStatement")

	sqlr = `SELECT hazardstatement.hazardstatement_id, hazardstatement.hazardstatement_label, hazardstatement.hazardstatement_reference
	FROM hazardstatement
	WHERE hazardstatement_id = ?`
	if err = db.Get(&hs, sqlr, id); err != nil {
		return HazardStatement{}, err
	}
	logger.Log.WithFields(logrus.Fields{"ID": id, "hs": hs}).Debug("GetProductsHazardStatement")
	return hs, nil
}

// GetProductsPrecautionaryStatement return the PrecautionaryStatement matching the given id
func (db *SQLiteDataStore) GetProductsPrecautionaryStatement(id int) (PrecautionaryStatement, error) {

	var (
		ps   PrecautionaryStatement
		sqlr string
		err  error
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetProductsPrecautionaryStatement")

	sqlr = `SELECT precautionarystatement.precautionarystatement_id, precautionarystatement.precautionarystatement_label, precautionarystatement.precautionarystatement_reference
	FROM precautionarystatement
	WHERE precautionarystatement_id = ?`
	if err = db.Get(&ps, sqlr, id); err != nil {
		return PrecautionaryStatement{}, err
	}
	logger.Log.WithFields(logrus.Fields{"ID": id, "ps": ps}).Debug("GetProductsPrecautionaryStatement")
	return ps, nil
}

// GetProductsSymbol return the symbol matching the given id
func (db *SQLiteDataStore) GetProductsSymbol(id int) (Symbol, error) {

	var (
		symbol Symbol
		sqlr   string
		err    error
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetProductsSymbol")

	sqlr = `SELECT symbol.symbol_id, symbol.symbol_label, symbol.symbol_image
	FROM symbol
	WHERE symbol_id = ?`
	if err = db.Get(&symbol, sqlr, id); err != nil {
		return Symbol{}, err
	}
	logger.Log.WithFields(logrus.Fields{"ID": id, "symbol": symbol}).Debug("GetProductsSymbol")
	return symbol, nil
}

// GetProductsSymbolByLabel return the symbol matching the given label
func (db *SQLiteDataStore) GetProductsSymbolByLabel(label string) (Symbol, error) {

	var (
		symbol Symbol
		sqlr   string
		err    error
	)
	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetProductsSymbolByLabel")

	sqlr = `SELECT symbol.symbol_id, symbol.symbol_label
	FROM symbol
	WHERE symbol.symbol_label = ?`
	if err = db.Get(&symbol, sqlr, label); err != nil {
		return Symbol{}, err
	}
	logger.Log.WithFields(logrus.Fields{"label": label, "symbol": symbol}).Debug("GetProductsSymbolByLabel")
	return symbol, nil
}

// GetProductsNames return the names matching the search criteria
func (db *SQLiteDataStore) GetProductsNames(p Dbselectparam) ([]Name, int, error) {
	var (
		names                              []Name
		count                              int
		exactSearch                        string
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	precreq.WriteString(" SELECT count(DISTINCT name.name_id)")
	presreq.WriteString(" SELECT name_id, name_label")

	comreq.WriteString(" FROM name")
	comreq.WriteString(" WHERE name_label LIKE upper(:search)")
	postsreq.WriteString(" ORDER BY INSTR(name_label, upper(\"" + exactSearch + "\")) ASC, name_label ASC")

	// limit
	if p.GetLimit() != ^uint64(0) {
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

	// setting the C attribute for formula matching exactly the search
	var name Name

	r := db.QueryRowx(`SELECT name_id, name_label FROM name WHERE name_label == ?`, exactSearch)
	if err = r.StructScan(&name); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	for i, n := range names {
		if n.NameID == name.NameID {
			names[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"names": names}).Debug("GetProductsNames")
	return names, count, nil
}

// GetProductsSymbols return the symbols matching the search criteria
func (db *SQLiteDataStore) GetProductsSymbols(p Dbselectparam) ([]Symbol, int, error) {
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
	postsreq.WriteString(" ORDER BY symbol_label " + p.GetOrder())

	// limit
	if p.GetLimit() != ^uint64(0) {
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

	logger.Log.WithFields(logrus.Fields{"symbols": symbols}).Debug("GetProductsSymbols")
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
	WHERE hazardstatement_reference = ?
	ORDER BY hazardstatement_reference`
	if err = db.Get(&hs, sqlr, r); err != nil {
		return HazardStatement{}, err
	}

	return hs, nil
}

// GetProductsHazardStatements return the hazard statements matching the search criteria
func (db *SQLiteDataStore) GetProductsHazardStatements(p Dbselectparam) ([]HazardStatement, int, error) {
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
	postsreq.WriteString(" ORDER BY hazardstatement_reference  " + p.GetOrder())

	// limit
	if p.GetLimit() != ^uint64(0) {
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

	logger.Log.WithFields(logrus.Fields{"hazardstatements": hazardstatements}).Debug("GetProductsHazardStatements")
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
	WHERE precautionarystatement_reference = ?
	ORDER BY precautionarystatement_reference`
	if err = db.Get(&ps, sqlr, r); err != nil {
		return PrecautionaryStatement{}, err
	}

	return ps, nil
}

// GetProductsPrecautionaryStatements return the hazard statements matching the search criteria
func (db *SQLiteDataStore) GetProductsPrecautionaryStatements(p Dbselectparam) ([]PrecautionaryStatement, int, error) {
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
	postsreq.WriteString(" ORDER BY precautionarystatement_reference  " + p.GetOrder())

	// limit
	if p.GetLimit() != ^uint64(0) {
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

	logger.Log.WithFields(logrus.Fields{"precautionarystatements": precautionarystatements}).Debug("GetProductsPrecautionaryStatements")
	return precautionarystatements, count, nil
}

// GetProductsPhysicalStates return the physical states matching the search criteria
func (db *SQLiteDataStore) GetProductsPhysicalStates(p Dbselectparam) ([]PhysicalState, int, error) {
	var (
		physicalstates                     []PhysicalState
		count                              int
		exactSearch                        string
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
	)

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	precreq.WriteString(" SELECT count(DISTINCT physicalstate.physicalstate_id)")
	presreq.WriteString(" SELECT physicalstate_id, physicalstate_label")

	comreq.WriteString(" FROM physicalstate")
	comreq.WriteString(" WHERE physicalstate_label LIKE :search")
	postsreq.WriteString(" ORDER BY INSTR(physicalstate_label, \"" + exactSearch + ")\") ASC, physicalstate_label ASC")

	// limit
	if p.GetLimit() != ^uint64(0) {
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
	var ps PhysicalState

	r := db.QueryRowx(`SELECT physicalstate_id, physicalstate_label FROM physicalstate WHERE physicalstate_label == ?`, exactSearch)
	if err = r.StructScan(&ps); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}
	for i, e := range physicalstates {
		if e.PhysicalStateID == ps.PhysicalStateID {
			physicalstates[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"physicalstates": physicalstates}).Debug("GetProductsPhysicalStates")
	return physicalstates, count, nil
}

// GetProductsPhysicalStateByLabel return the physical state matching the given label
func (db *SQLiteDataStore) GetProductsPhysicalStateByLabel(label string) (PhysicalState, error) {

	var (
		ps   PhysicalState
		sqlr string
		err  error
	)
	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetProductsPhysicalStateByLabel")

	sqlr = `SELECT physicalstate.physicalstate_id, physicalstate.physicalstate_label
	FROM physicalstate
	WHERE physicalstate.physicalstate_label = ?
	ORDER BY physicalstate.physicalstate_label`
	if err = db.Get(&ps, sqlr, label); err != nil {
		return PhysicalState{}, err
	}
	logger.Log.WithFields(logrus.Fields{"label": label, "ps": ps}).Debug("GetProductsPhysicalStateByLabel")
	return ps, nil
}

// GetProductsSignalWords return the signal words matching the search criteria
func (db *SQLiteDataStore) GetProductsSignalWords(p Dbselectparam) ([]SignalWord, int, error) {
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
	if p.GetLimit() != ^uint64(0) {
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

	logger.Log.WithFields(logrus.Fields{"signalwords": signalwords}).Debug("GetProductsSignalWords")
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
	p.product_twodformula,
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
func (db *SQLiteDataStore) GetProducts(p DbselectparamProduct) ([]Product, int, error) {
	//defer TimeTrack(time.Now(), "GetProducts")

	var (
		products                           []Product
		count                              int
		precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                             *sqlx.NamedStmt
		snstmt                             *sqlx.NamedStmt
		err                                error
		rperm                              bool
		isadmin                            bool
		wg                                 sync.WaitGroup
	)
	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("GetProducts")

	// is the user an admin?
	if isadmin, err = db.IsPersonAdmin(p.GetLoggedPersonID()); err != nil {
		return nil, 0, err
	}
	// has the person rproducts permission?
	if rperm, err = db.HasPersonReadRestrictedProductPermission(p.GetLoggedPersonID()); err != nil {
		return nil, 0, err
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
	category.category_label AS "category.category_label",
	GROUP_CONCAT(DISTINCT storage.storage_barecode) AS "product_sl"
	`)

	if p.GetCasNumberCmr() {
		presreq.WriteString(`,GROUP_CONCAT(DISTINCT hazardstatement.hazardstatement_cmr) AS "hazardstatement_cmr"`)
	}

	// common parts
	comreq.WriteString(" FROM product as p")
	// CMR
	if p.GetCasNumberCmr() {
		comreq.WriteString(" LEFT JOIN producthazardstatements ON producthazardstatements.producthazardstatements_product_id = p.product_id")
		comreq.WriteString(" LEFT JOIN hazardstatement ON producthazardstatements.producthazardstatements_hazardstatement_id = hazardstatement.hazardstatement_id")
	}
	// get name
	comreq.WriteString(" JOIN name ON p.name = name.name_id")
	// get category
	if p.GetCategory() != -1 {
		comreq.WriteString(" JOIN category ON p.category = :category")
	} else {
		comreq.WriteString(" LEFT JOIN category ON p.category = category.category_id")
	}
	// get unit_temperature
	comreq.WriteString(" LEFT JOIN unit ut ON p.unit_temperature = ut.unit_id")
	// get producerref
	if p.GetProducerRef() != -1 {
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
	if p.GetEntity() != -1 || p.GetStorelocation() != -1 || p.GetStorageBarecode() != "" {
		comreq.WriteString(" JOIN storelocation ON storage.storelocation = storelocation.storelocation_id")
		comreq.WriteString(" JOIN entity ON storelocation.entity = entity.entity_id")
	}
	// get borrowings
	if p.GetBorrowing() {
		comreq.WriteString(" JOIN borrowing ON borrowing.storage = storage.storage_id AND borrowing.borrower = :personid")
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
	// get tags
	if len(p.GetTags()) != 0 {
		comreq.WriteString(" JOIN producttags AS ptags ON ptags.producttags_product_id = p.product_id")
	}

	// filter by permissions
	comreq.WriteString(` JOIN permission AS perm ON
	perm.person = :personid and 
	(perm.permission_item_name in ("all", "products")) and 
	(perm.permission_perm_name in ("all", "r", "w"))
	`)
	comreq.WriteString(" WHERE 1")
	if p.GetStorageToDestroy() {
		comreq.WriteString(" AND storage.storage_todestroy = true")
	}
	if p.GetCasNumberCmr() {
		comreq.WriteString(" AND (casnumber.casnumber_cmr IS NOT NULL OR (hazardstatement_cmr IS NOT NULL AND hazardstatement_cmr != ''))")
	}
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
	if p.GetStorageBatchNumber() != "" {
		comreq.WriteString(" AND storage.storage_batchnumber LIKE :storage_batchnumber")
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
	if len(p.GetTags()) != 0 {
		comreq.WriteString(" AND ptags.producttags_tag_id IN (")
		for _, t := range p.GetTags() {
			comreq.WriteString(fmt.Sprintf("%d,", t))
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

	// show bio/chem/consu
	if !p.GetShowChem() && !p.GetShowBio() && p.GetShowConsu() {
		comreq.WriteString(" AND (product_number_per_carton IS NOT NULL AND product_number_per_carton != 0)")
	} else if !p.GetShowChem() && p.GetShowBio() && !p.GetShowConsu() {
		comreq.WriteString(" AND producerref IS NOT NULL")
		comreq.WriteString(" AND (product_number_per_carton IS NULL OR product_number_per_carton == 0)")
	} else if !p.GetShowChem() && p.GetShowBio() && p.GetShowConsu() {
		comreq.WriteString(" AND ((product_number_per_carton IS NOT NULL AND product_number_per_carton != 0)")
		comreq.WriteString(" OR producerref IS NOT NULL)")
	} else if p.GetShowChem() && !p.GetShowBio() && !p.GetShowConsu() {
		comreq.WriteString(" AND producerref IS NULL")
		comreq.WriteString(" AND (product_number_per_carton IS NULL OR product_number_per_carton == 0)")
	} else if p.GetShowChem() && !p.GetShowBio() && p.GetShowConsu() {
		comreq.WriteString(" AND (producerref IS NULL")
		comreq.WriteString(" OR (product_number_per_carton IS NOT NULL AND product_number_per_carton != 0))")
	} else if p.GetShowChem() && p.GetShowBio() && !p.GetShowConsu() {
		comreq.WriteString(" AND (product_number_per_carton IS NULL OR product_number_per_carton == 0)")
	}

	// post select request
	postsreq.WriteString(" GROUP BY p.product_id")
	postsreq.WriteString(" ORDER BY " + p.GetOrderBy() + " " + p.GetOrder())

	// limit
	if p.GetLimit() != ^uint64(0) {
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
		"storage_batchnumber": p.GetStorageBatchNumber(),
		"custom_name_part_of": "%" + p.GetCustomNamePartOf() + "%",
		"signalword":          p.GetSignalWord(),
		"producerref":         p.GetProducerRef(),
		"category":            p.GetCategory(),
	}

	logger.Log.Debug(presreq.String() + comreq.String() + postsreq.String())
	logger.Log.Debug(m)

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

	//
	// getting supplierref
	//
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

	//
	// getting tags
	//
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
	wg.Add(1)
	go func() {

		var (
			reqtsc, reqsc, reqasc strings.Builder
		)

		for i, pr := range products {
			// getting the total storage count
			reqtsc.Reset()
			reqtsc.WriteString("SELECT count(DISTINCT storage_id) from storage")
			reqtsc.WriteString(" JOIN product ON storage.product = ? AND storage.storage IS NULL AND storage.storage_archive == false")

			if isadmin {
				reqsc.Reset()
				reqsc.WriteString("SELECT count(DISTINCT storage_id) from storage")
				reqsc.WriteString(" JOIN product ON storage.product = ? AND storage.storage IS NULL AND storage.storage_archive == false")

				reqasc.Reset()
				reqasc.WriteString("SELECT count(DISTINCT storage_id) from storage")
				reqasc.WriteString(" JOIN product ON storage.product = ? AND storage.storage_archive == true")
			} else {
				// getting the storage count of the logged user entities
				reqsc.Reset()
				reqsc.WriteString("SELECT count(DISTINCT storage_id) from storage")
				reqsc.WriteString(" JOIN product ON storage.product = ? AND storage.storage IS NULL AND storage.storage_archive == false")
				reqsc.WriteString(" JOIN storelocation ON storage.storelocation = storelocation.storelocation_id")
				reqsc.WriteString(" JOIN entity ON storelocation.entity = entity.entity_id")
				reqsc.WriteString(" JOIN personentities ON (entity.entity_id = personentities.personentities_entity_id) AND")
				reqsc.WriteString(" (personentities.personentities_person_id = ?)")

				reqasc.Reset()
				reqasc.WriteString("SELECT count(DISTINCT storage_id) from storage")
				reqasc.WriteString(" JOIN product ON storage.product = ? AND storage.storage_archive == true")
				reqasc.WriteString(" JOIN storelocation ON storage.storelocation = storelocation.storelocation_id")
				reqasc.WriteString(" JOIN entity ON storelocation.entity = entity.entity_id")
				reqasc.WriteString(" JOIN personentities ON (entity.entity_id = personentities.personentities_entity_id) AND")
				reqasc.WriteString(" (personentities.personentities_person_id = ?)")
			}
			if err = db.Get(&products[i].ProductSC, reqsc.String(), pr.ProductID, p.GetLoggedPersonID()); err != nil {
				logger.Log.WithFields(logrus.Fields{"err": err}).Error("GetProducts:goroutine:SC")
			}
			if err = db.Get(&products[i].ProductASC, reqasc.String(), pr.ProductID, p.GetLoggedPersonID()); err != nil {
				logger.Log.WithFields(logrus.Fields{"err": err}).Error("GetProducts:goroutine:ASC")
			}
			if err = db.Get(&products[i].ProductTSC, reqtsc.String(), pr.ProductID, p.GetLoggedPersonID()); err != nil {
				logger.Log.WithFields(logrus.Fields{"err": err}).Error("GetProducts:goroutine:TSC")
			}
		}
		wg.Done()
	}()

	wg.Wait()

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

	logger.Log.WithFields(logrus.Fields{"count": count}).Debug("CountProductStorages")
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
		return Product{}, err
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

	logger.Log.WithFields(logrus.Fields{"id": id, "product": product}).Debug("GetProduct")
	return product, nil
}

// DeleteProduct deletes the product with the given id
func (db *SQLiteDataStore) DeleteProduct(id int) error {
	var (
		sqlr string
		err  error
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("DeleteProduct")
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
	if v, err := p.CasNumber.CasNumberID.Value(); p.CasNumber.CasNumberID.Valid && err == nil && v.(int64) == -1 {
		logger.Log.Debug("new casnumber " + p.CasNumberLabel.String)
		sqlr = `INSERT INTO casnumber (casnumber_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.CasNumberLabel); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
		// updating the product CasNumberID (CasNumberLabel already set)
		p.CasNumber.CasNumberID = sql.NullInt64{Valid: true, Int64: int64(lastid)}
	}
	// if CeNumberID = -1 then it is a new ce
	if v, err := p.CeNumber.CeNumberID.Value(); p.CeNumber.CeNumberID.Valid && err == nil && v.(int64) == -1 {
		logger.Log.Debug("new cenumber " + p.CeNumberLabel.String)
		sqlr = `INSERT INTO cenumber (cenumber_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.CeNumberLabel.String); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
		// updating the product CeNumberID (CeNumberLabel already set)
		p.CeNumber.CeNumberID = sql.NullInt64{Valid: true, Int64: lastid}
	}
	if err != nil {
		logger.Log.Error("cenumber error - " + err.Error())
		if errr := tx.Rollback(); errr != nil {
			return 0, errr
		}
		return 0, err
	}
	// if NameID = -1 then it is a new name
	if p.Name.NameID == -1 {
		logger.Log.Debug("new name " + p.NameLabel)
		sqlr = `INSERT INTO name (name_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, strings.ToUpper(p.NameLabel)); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
		// updating the product NameID (NameLabel already set)
		p.Name.NameID = int(lastid)
	}
	for i, syn := range p.Synonyms {
		if syn.NameID == -1 {
			logger.Log.Debug("new name(syn) " + syn.NameLabel)
			sqlr = `INSERT INTO name (name_label) VALUES (?)`
			if res, err = tx.Exec(sqlr, strings.ToUpper(syn.NameLabel)); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return 0, errr
				}
				return 0, err
			}
			// getting the last inserted id
			if lastid, err = res.LastInsertId(); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return 0, errr
				}
				return 0, err
			}
			p.Synonyms[i].NameID = int(lastid)
		}
	}
	// if ClassOfCompoundID = -1 then it is a new class of compounds
	for i, coc := range p.ClassOfCompound {
		if coc.ClassOfCompoundID == -1 {
			logger.Log.Debug("new classofcompound " + coc.ClassOfCompoundLabel)
			sqlr = `INSERT INTO classofcompound (classofcompound_label) VALUES (?)`
			if res, err = tx.Exec(sqlr, strings.ToUpper(coc.ClassOfCompoundLabel)); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return 0, errr
				}
				return 0, err
			}
			// getting the last inserted id
			if lastid, err = res.LastInsertId(); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return 0, errr
				}
				return 0, err
			}
			p.ClassOfCompound[i].ClassOfCompoundID = int(lastid)
		}
	}
	// if SupplierRefID = -1 then it is a new supplier ref
	for i, sr := range p.SupplierRefs {
		if sr.SupplierRefID == -1 {
			logger.Log.Debug("new supplierref " + sr.SupplierRefLabel)
			sqlr = `INSERT INTO supplierref (supplierref_label, supplier) VALUES (?, ?)`
			if res, err = tx.Exec(sqlr, sr.SupplierRefLabel, sr.Supplier.SupplierID); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return 0, errr
				}
				return 0, err
			}
			// getting the last inserted id
			if lastid, err = res.LastInsertId(); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return 0, errr
				}
				return 0, err
			}
			p.SupplierRefs[i].SupplierRefID = int(lastid)
		}
	}
	// if TagID = -1 then it is a new tag
	for i, tag := range p.Tags {
		if tag.TagID == -1 {
			logger.Log.Debug("new tag " + tag.TagLabel)
			sqlr = `INSERT INTO tag (tag_label) VALUES (?)`
			if res, err = tx.Exec(sqlr, tag.TagLabel); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return 0, errr
				}
				return 0, err
			}
			// getting the last inserted id
			if lastid, err = res.LastInsertId(); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return 0, errr
				}
				return 0, err
			}
			p.Tags[i].TagID = int(lastid)
		}
	}
	// if EmpiricalFormulaID = -1 then it is a new empirical formula
	if v, err := p.EmpiricalFormula.EmpiricalFormulaID.Value(); p.EmpiricalFormula.EmpiricalFormulaID.Valid && err == nil && v.(int64) == -1 {
		logger.Log.Debug("new empiricalformula " + p.EmpiricalFormulaLabel.String)
		sqlr = `INSERT INTO empiricalformula (empiricalformula_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.EmpiricalFormulaLabel); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
		// updating the product EmpiricalFormulaID (EmpiricalFormulaLabel already set)
		p.EmpiricalFormula.EmpiricalFormulaID = sql.NullInt64{Valid: true, Int64: int64(lastid)}
	}
	// if LinearFormulaID = -1 then it is a new linear formula
	if v, err := p.LinearFormula.LinearFormulaID.Value(); p.LinearFormula.LinearFormulaID.Valid && err == nil && v.(int64) == -1 {
		logger.Log.Debug("new linearformula " + p.LinearFormulaLabel.String)
		sqlr = `INSERT INTO linearformula (linearformula_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.LinearFormulaLabel.String); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
		// updating the product LinearFormulaID (LinearFormulaLabel already set)
		p.LinearFormula.LinearFormulaID = sql.NullInt64{Valid: true, Int64: lastid}
	}
	// if PhysicalStateID = -1 then it is a new physical state
	if v, err := p.PhysicalState.PhysicalStateID.Value(); p.PhysicalState.PhysicalStateID.Valid && err == nil && v.(int64) == -1 {
		logger.Log.Debug("new physicalstate " + p.PhysicalStateLabel.String)
		sqlr = `INSERT INTO physicalstate (physicalstate_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.PhysicalStateLabel.String); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
		// updating the product PhysicalStateID (PhysicalStateLabel already set)
		p.PhysicalState.PhysicalStateID = sql.NullInt64{Valid: true, Int64: lastid}
	}
	// if CategoryID = -1 then it is a new category
	if v, err := p.Category.CategoryID.Value(); p.Category.CategoryID.Valid && err == nil && v.(int64) == -1 {
		logger.Log.Debug("new category " + p.CategoryLabel.String)
		sqlr = `INSERT INTO category (category_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.CategoryLabel.String); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
		// updating the product PhysicalStateID (PhysicalStateLabel already set)
		p.Category.CategoryID = sql.NullInt64{Valid: true, Int64: lastid}
	}
	// if ProducerRefID = -1 then it is a new producer ref
	if v, err := p.ProducerRef.ProducerRefID.Value(); p.ProducerRef.ProducerRefID.Valid && err == nil && v.(int64) == -1 {
		logger.Log.Debug("new producerref " + p.ProducerRefLabel.String)
		sqlr = `INSERT INTO producerref (producerref_label, producer) VALUES (?, ?)`
		if res, err = tx.Exec(sqlr, p.ProducerRefLabel.String, p.Producer.ProducerID); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
		// updating the product ProducerRefID (ProducerRefLabel already set)
		p.ProducerRef.ProducerRefID = sql.NullInt64{Valid: true, Int64: lastid}
	}

	// finally updating the product
	s := make(map[string]interface{})
	if p.ProductSpecificity.Valid {
		s["product_specificity"] = p.ProductSpecificity.String
	}
	if p.ProductMSDS.Valid {
		s["product_msds"] = p.ProductMSDS.String
	}
	if p.ProductSheet.Valid {
		s["product_sheet"] = p.ProductSheet.String
	}
	if p.ProductTemperature.Valid {
		s["product_temperature"] = int(p.ProductTemperature.Int64)
	}
	if p.ProductRestricted.Valid {
		s["product_restricted"] = p.ProductRestricted.Bool
	}
	if p.ProductRadioactive.Valid {
		s["product_radioactive"] = p.ProductRadioactive.Bool
	}

	if p.Category.CategoryID.Valid {
		s["category"] = int(p.Category.CategoryID.Int64)
	}
	if p.UnitTemperature.UnitID.Valid {
		s["unit_temperature"] = int(p.UnitTemperature.UnitID.Int64)
	}
	if p.ProductThreeDFormula.Valid {
		s["product_threedformula"] = p.ProductThreeDFormula.String
	}
	if p.ProductTwoDFormula.Valid {
		s["product_twodformula"] = p.ProductTwoDFormula.String
	}
	if p.ProductDisposalComment.Valid {
		s["product_disposalcomment"] = p.ProductDisposalComment.String
	}
	if p.ProductRemark.Valid {
		s["product_remark"] = p.ProductRemark.String
	}
	if p.ProductNumberPerCarton.Valid {
		s["product_number_per_carton"] = p.ProductNumberPerCarton.Int64
	}
	if p.ProductNumberPerBag.Valid {
		s["product_number_per_bag"] = p.ProductNumberPerBag.Int64
	}
	if p.EmpiricalFormulaID.Valid {
		s["empiricalformula"] = int(p.EmpiricalFormulaID.Int64)
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
	if p.CasNumberID.Valid {
		s["casnumber"] = int(p.CasNumberID.Int64)
	}
	if p.CeNumberID.Valid {
		s["cenumber"] = int(p.CeNumberID.Int64)
	}
	if p.ProducerRefID.Valid {
		s["producerref"] = int(p.ProducerRefID.Int64)
	}
	if p.ProductMolFormula.Valid {
		s["product_molformula"] = p.ProductMolFormula.String
	}

	s["name"] = p.NameID
	s["person"] = p.PersonID

	// building column names/values
	col := make([]string, 0, len(s))
	val := make([]interface{}, 0, len(s))
	for k, v := range s {
		col = append(col, k)

		switch v.(type) {
		case int:
			val = append(val, v.(int))
		case int64:
			val = append(val, v.(int64))
		case string:
			val = append(val, v.(string))
		case bool:
			val = append(val, v.(bool))
		default:
			val = append(val, v)
		}
	}

	ibuilder = sq.Insert("product").Columns(col...).Values(val...)
	if sqlr, sqla, err = ibuilder.ToSql(); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return 0, errr
		}
		return 0, err
	}

	//logger.Log.Debug(sqlr)
	//logger.Log.Debug(sqla)

	if res, err = tx.Exec(sqlr, sqla...); err != nil {
		logger.Log.Error("product error - " + err.Error())
		logger.Log.Error("sql:" + sqlr)
		if errr := tx.Rollback(); errr != nil {
			return 0, errr
		}
		return 0, err
	}

	// getting the last inserted id
	if lastid, err = res.LastInsertId(); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return 0, errr
		}
		return 0, err
	}
	p.ProductID = int(lastid)
	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("CreateProduct")

	// adding supplierrefs
	for _, sr := range p.SupplierRefs {
		sqlr = `INSERT INTO productsupplierrefs (productsupplierrefs_product_id, productsupplierrefs_supplierref_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, sr.SupplierRefID); err != nil {
			logger.Log.Error("productsupplierrefs error - " + err.Error())
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
	}

	// adding tags
	for _, tag := range p.Tags {
		sqlr = `INSERT INTO producttags (producttags_product_id, producttags_tag_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, tag.TagID); err != nil {
			logger.Log.Error("producttags error - " + err.Error())
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
	}

	// adding symbols
	for _, sym := range p.Symbols {
		sqlr = `INSERT INTO productsymbols (productsymbols_product_id, productsymbols_symbol_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, sym.SymbolID); err != nil {
			logger.Log.Error("productsymbols error - " + err.Error())
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
	}

	// adding classes of compounds
	for _, coc := range p.ClassOfCompound {
		sqlr = `INSERT INTO productclassofcompound (productclassofcompound_product_id, productclassofcompound_classofcompound_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, coc.ClassOfCompoundID); err != nil {
			logger.Log.Error("productclassofcompound error - " + err.Error())
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
	}
	// adding hazard statements
	for _, hs := range p.HazardStatements {
		sqlr = `INSERT INTO producthazardstatements (producthazardstatements_product_id, producthazardstatements_hazardstatement_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, hs.HazardStatementID); err != nil {
			logger.Log.Error("producthazardstatements error - " + err.Error())
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
	}
	// adding precautionary statements
	for _, ps := range p.PrecautionaryStatements {
		sqlr = `INSERT INTO productprecautionarystatements (productprecautionarystatements_product_id, productprecautionarystatements_precautionarystatement_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, ps.PrecautionaryStatementID); err != nil {
			logger.Log.Error("productprecautionarystatements error - " + err.Error())
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
	}
	// adding synonyms
	for _, syn := range p.Synonyms {
		sqlr = `INSERT INTO productsynonyms (productsynonyms_product_id, productsynonyms_name_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, syn.NameID); err != nil {
			logger.Log.Error("productsynonyms error - " + err.Error())
			if errr := tx.Rollback(); errr != nil {
				return 0, errr
			}
			return 0, err
		}
	}
	// committing changes
	if err = tx.Commit(); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return 0, errr
		}
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
	if v, err := p.CasNumber.CasNumberID.Value(); p.CasNumber.CasNumberID.Valid && err == nil && v.(int64) == -1 {
		logger.Log.Debug("new casnumber " + p.CasNumberLabel.String)
		sqlr = `INSERT INTO casnumber (casnumber_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.CasNumberLabel); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// updating the product CasNumberID (CasNumberLabel already set)
		p.CasNumber.CasNumberID = sql.NullInt64{Valid: true, Int64: int64(lastid)}
	}
	// if CeNumberID = -1 then it is a new ce
	if v, err := p.CeNumber.CeNumberID.Value(); p.CeNumber.CeNumberID.Valid && err == nil && v.(int64) == -1 {
		logger.Log.Debug("new cenumber " + p.CeNumberLabel.String)
		sqlr = `INSERT INTO cenumber (cenumber_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.CeNumberLabel.String); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// updating the product CeNumberID (CeNumberLabel already set)
		p.CeNumber.CeNumberID = sql.NullInt64{Valid: true, Int64: lastid}
	}
	if err != nil {
		logger.Log.Error("cenumber error - " + err.Error())
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}
	// if NameID = -1 then it is a new name
	if p.Name.NameID == -1 {
		logger.Log.Debug("new name " + p.NameLabel)
		sqlr = `INSERT INTO name (name_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, strings.ToUpper(p.NameLabel)); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// updating the product NameID (NameLabel already set)
		p.Name.NameID = int(lastid)
	}
	for i, syn := range p.Synonyms {
		if syn.NameID == -1 {
			logger.Log.Debug("new name(syn) " + syn.NameLabel)
			sqlr = `INSERT INTO name (name_label) VALUES (?)`
			if res, err = tx.Exec(sqlr, strings.ToUpper(syn.NameLabel)); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			// getting the last inserted id
			if lastid, err = res.LastInsertId(); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			p.Synonyms[i].NameID = int(lastid)
		}
	}
	// if ClassOfCompoundID = -1 then it is a new class of compounds
	for i, coc := range p.ClassOfCompound {
		if coc.ClassOfCompoundID == -1 {
			logger.Log.Debug("new classofcompound " + coc.ClassOfCompoundLabel)
			sqlr = `INSERT INTO classofcompound (classofcompound_label) VALUES (?)`
			if res, err = tx.Exec(sqlr, strings.ToUpper(coc.ClassOfCompoundLabel)); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			// getting the last inserted id
			if lastid, err = res.LastInsertId(); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			p.ClassOfCompound[i].ClassOfCompoundID = int(lastid)
		}
	}
	// if SupplierRefID = -1 then it is a new supplier ref
	for i, sr := range p.SupplierRefs {
		if sr.SupplierRefID == -1 {
			logger.Log.Debug("new supplierref " + sr.SupplierRefLabel)
			sqlr = `INSERT INTO supplierref (supplierref_label, supplier) VALUES (?, ?)`
			if res, err = tx.Exec(sqlr, sr.SupplierRefLabel, sr.Supplier.SupplierID); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			// getting the last inserted id
			if lastid, err = res.LastInsertId(); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			p.SupplierRefs[i].SupplierRefID = int(lastid)
		}
	}
	// if TagID = -1 then it is a new tag
	for i, tag := range p.Tags {
		if tag.TagID == -1 {
			logger.Log.Debug("new tag " + tag.TagLabel)
			sqlr = `INSERT INTO tag (tag_label) VALUES (?)`
			if res, err = tx.Exec(sqlr, tag.TagLabel); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			// getting the last inserted id
			if lastid, err = res.LastInsertId(); err != nil {
				if errr := tx.Rollback(); errr != nil {
					return errr
				}
				return err
			}
			p.Tags[i].TagID = int(lastid)
		}
	}
	// if EmpiricalFormulaID = -1 then it is a new empirical formula
	if v, err := p.EmpiricalFormula.EmpiricalFormulaID.Value(); p.EmpiricalFormula.EmpiricalFormulaID.Valid && err == nil && v.(int64) == -1 {
		logger.Log.Debug("new empiricalformula " + p.EmpiricalFormulaLabel.String)
		sqlr = `INSERT INTO empiricalformula (p.EmpiricalFormulaLabel) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.EmpiricalFormulaLabel); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// updating the product EmpiricalFormulaID (EmpiricalFormulaLabel already set)
		p.EmpiricalFormula.EmpiricalFormulaID = sql.NullInt64{Valid: true, Int64: int64(lastid)}
	}
	// if LinearFormulaID = -1 then it is a new linear formula
	if v, err := p.LinearFormula.LinearFormulaID.Value(); p.LinearFormula.LinearFormulaID.Valid && err == nil && v.(int64) == -1 {
		logger.Log.Debug("new linearformula " + p.LinearFormulaLabel.String)
		sqlr = `INSERT INTO linearformula (linearformula_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.LinearFormulaLabel.String); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// updating the product LinearFormulaID (LinearFormulaLabel already set)
		p.LinearFormula.LinearFormulaID = sql.NullInt64{Valid: true, Int64: lastid}
	}
	// if PhysicalStateID = -1 then it is a new physical state
	if v, err := p.PhysicalState.PhysicalStateID.Value(); p.PhysicalState.PhysicalStateID.Valid && err == nil && v.(int64) == -1 {
		logger.Log.Debug("new physicalstate " + p.PhysicalStateLabel.String)
		sqlr = `INSERT INTO physicalstate (physicalstate_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.PhysicalStateLabel.String); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// updating the product PhysicalStateID (PhysicalStateLabel already set)
		p.PhysicalState.PhysicalStateID = sql.NullInt64{Valid: true, Int64: lastid}
	}
	// if CategoryID = -1 then it is a new category
	if v, err := p.Category.CategoryID.Value(); p.Category.CategoryID.Valid && err == nil && v.(int64) == -1 {
		logger.Log.Debug("new category " + p.CategoryLabel.String)
		sqlr = `INSERT INTO category (category_label) VALUES (?)`
		if res, err = tx.Exec(sqlr, p.CategoryLabel.String); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// updating the product PhysicalStateID (PhysicalStateLabel already set)
		p.Category.CategoryID = sql.NullInt64{Valid: true, Int64: lastid}
	}
	// if ProducerRefID = -1 then it is a new physical state
	if v, err := p.ProducerRef.ProducerRefID.Value(); p.ProducerRef.ProducerRefID.Valid && err == nil && v.(int64) == -1 {
		logger.Log.Debug("new producerref " + p.ProducerRefLabel.String)
		sqlr = `INSERT INTO producerref (producerref_label, producer) VALUES (?, ?)`
		if res, err = tx.Exec(sqlr, p.ProducerRefLabel.String, p.Producer.ProducerID); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// getting the last inserted id
		if lastid, err = res.LastInsertId(); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
		// updating the product ProducerRefID (PhysicalStateLabel already set)
		p.ProducerRef.ProducerRefID = sql.NullInt64{Valid: true, Int64: lastid}
	}

	// finally updating the product
	s := make(map[string]interface{})
	if p.ProductSpecificity.Valid {
		s["product_specificity"] = p.ProductSpecificity.String
	}
	if p.ProductMSDS.Valid {
		s["product_msds"] = p.ProductMSDS.String
	}
	if p.ProductSheet.Valid {
		s["product_sheet"] = p.ProductSheet.String
	}
	if p.ProductTemperature.Valid {
		s["product_temperature"] = int(p.ProductTemperature.Int64)
	}
	if p.ProductRestricted.Valid {
		s["product_restricted"] = p.ProductRestricted.Bool
	}
	if p.ProductRadioactive.Valid {
		s["product_radioactive"] = p.ProductRadioactive.Bool
	}

	if p.Category.CategoryID.Valid {
		s["category"] = int(p.Category.CategoryID.Int64)
	}
	if p.UnitTemperature.UnitID.Valid {
		s["unit_temperature"] = int(p.UnitTemperature.UnitID.Int64)
	}
	if p.ProductThreeDFormula.Valid {
		s["product_threedformula"] = p.ProductThreeDFormula.String
	}
	if p.ProductTwoDFormula.Valid {
		s["product_twodformula"] = p.ProductTwoDFormula.String
	}
	if p.ProductDisposalComment.Valid {
		s["product_disposalcomment"] = p.ProductDisposalComment.String
	}
	if p.ProductRemark.Valid {
		s["product_remark"] = p.ProductRemark.String
	}
	if p.ProductNumberPerCarton.Valid {
		s["product_number_per_carton"] = p.ProductNumberPerCarton.Int64
	}
	if p.ProductNumberPerBag.Valid {
		s["product_number_per_bag"] = p.ProductNumberPerBag.Int64
	}
	if p.EmpiricalFormulaID.Valid {
		s["empiricalformula"] = int(p.EmpiricalFormulaID.Int64)
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
	if p.CasNumberID.Valid {
		s["casnumber"] = int(p.CasNumberID.Int64)
	}
	if p.CeNumberID.Valid {
		s["cenumber"] = int(p.CeNumberID.Int64)
	}
	if p.ProducerRefID.Valid {
		s["producerref"] = int(p.ProducerRefID.Int64)
	}
	if p.ProductMolFormula.Valid {
		s["product_molformula"] = p.ProductMolFormula.String
	}
	s["name"] = p.NameID
	s["person"] = p.PersonID

	ubuilder = sq.Update("product").
		SetMap(s).
		Where(sq.Eq{"product_id": p.ProductID})
	if sqlr, sqla, err = ubuilder.ToSql(); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}
	if _, err = tx.Exec(sqlr, sqla...); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}

	// deleting supplierrefs
	sqlr = `DELETE FROM productsupplierrefs WHERE productsupplierrefs.productsupplierrefs_product_id = (?)`
	if _, err = tx.Exec(sqlr, p.ProductID); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}
	// adding new ones
	for _, sr := range p.SupplierRefs {
		sqlr = `INSERT INTO productsupplierrefs (productsupplierrefs_product_id, productsupplierrefs_supplierref_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, sr.SupplierRefID); err != nil {
			logger.Log.Error("productsupplierrefs error - " + err.Error())
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
	}

	// deleting tags
	sqlr = `DELETE FROM producttags WHERE producttags.producttags_product_id = (?)`
	if _, err = tx.Exec(sqlr, p.ProductID); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}
	// adding new ones
	for _, tag := range p.Tags {
		sqlr = `INSERT INTO producttags (producttags_product_id, producttags_tag_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, tag.TagID); err != nil {
			logger.Log.Error("producttags error - " + err.Error())
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
	}

	// deleting symbols
	sqlr = `DELETE FROM productsymbols WHERE productsymbols.productsymbols_product_id = (?)`
	if _, err = tx.Exec(sqlr, p.ProductID); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}
	// adding new ones
	for _, sym := range p.Symbols {
		sqlr = `INSERT INTO productsymbols (productsymbols_product_id, productsymbols_symbol_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, sym.SymbolID); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
	}

	// deleting synonyms
	sqlr = `DELETE FROM productsynonyms WHERE productsynonyms.productsynonyms_product_id = (?)`
	if _, err = tx.Exec(sqlr, p.ProductID); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}
	// adding new ones
	for _, syn := range p.Synonyms {
		sqlr = `INSERT INTO productsynonyms (productsynonyms_product_id, productsynonyms_name_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, syn.NameID); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
	}

	// deleting classes of compounds
	sqlr = `DELETE FROM productclassofcompound WHERE productclassofcompound.productclassofcompound_product_id = (?)`
	if _, err = tx.Exec(sqlr, p.ProductID); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}
	// adding new ones
	for _, coc := range p.ClassOfCompound {
		sqlr = `INSERT INTO productclassofcompound (productclassofcompound_product_id, productclassofcompound_classofcompound_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, coc.ClassOfCompoundID); err != nil {
			logger.Log.Error("productclassofcompound error - " + err.Error())
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
	}

	// deleting hazard statements
	sqlr = `DELETE FROM producthazardstatements WHERE producthazardstatements.producthazardstatements_product_id = (?)`
	if _, err = tx.Exec(sqlr, p.ProductID); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}
	// adding new ones
	for _, hs := range p.HazardStatements {
		sqlr = `INSERT INTO producthazardstatements (producthazardstatements_product_id, producthazardstatements_hazardstatement_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, hs.HazardStatementID); err != nil {
			logger.Log.Error("producthazardstatements error - " + err.Error())
			if errr := tx.Rollback(); errr != nil {
				return errr
			}
			return err
		}
	}

	// deleting precautionary statements
	sqlr = `DELETE FROM productprecautionarystatements WHERE productprecautionarystatements.productprecautionarystatements_product_id = (?)`
	if _, err = tx.Exec(sqlr, p.ProductID); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}
	// adding new ones
	for _, ps := range p.PrecautionaryStatements {
		sqlr = `INSERT INTO productprecautionarystatements (productprecautionarystatements_product_id, productprecautionarystatements_precautionarystatement_id) VALUES (?,?)`
		if _, err = tx.Exec(sqlr, p.ProductID, ps.PrecautionaryStatementID); err != nil {
			logger.Log.Error("productprecautionarystatements error - " + err.Error())
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

// CreateProducer create a new producer in the db
func (db *SQLiteDataStore) CreateProducer(p Producer) (int, error) {
	var (
		sqlr   string
		tx     *sql.Tx
		res    sql.Result
		lastid int64
		err    error
	)
	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("CreateProducer")

	if !p.ProducerLabel.Valid {
		return 0, errors.New("empty string")
	}

	// beginning transaction
	if tx, err = db.Begin(); err != nil {
		return 0, err
	}

	sqlr = `INSERT INTO producer(producer_label) VALUES (?)`
	if res, err = tx.Exec(sqlr, p.ProducerLabel); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return 0, errr
		}
		return 0, err
	}

	// getting the last inserted id
	if lastid, err = res.LastInsertId(); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return 0, errr
		}
		return 0, err
	}

	// committing changes
	if err = tx.Commit(); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return 0, errr
		}
		return 0, err
	}

	return int(lastid), nil
}

// CreateSupplier create a new supplier in the db
func (db *SQLiteDataStore) CreateSupplier(s Supplier) (int, error) {
	var (
		sqlr   string
		tx     *sql.Tx
		res    sql.Result
		lastid int64
		err    error
	)
	logger.Log.WithFields(logrus.Fields{"s": s}).Debug("CreateSupplier")

	if !s.SupplierLabel.Valid {
		return 0, errors.New("empty string")
	}

	// beginning transaction
	if tx, err = db.Begin(); err != nil {
		return 0, err
	}

	sqlr = `INSERT INTO supplier(supplier_label) VALUES (?)`
	if res, err = tx.Exec(sqlr, s.SupplierLabel); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return 0, errr
		}
		return 0, err
	}

	// getting the last inserted id
	if lastid, err = res.LastInsertId(); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return 0, errr
		}
		return 0, err
	}

	// committing changes
	if err = tx.Commit(); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return 0, errr
		}
		return 0, err
	}

	return int(lastid), nil
}

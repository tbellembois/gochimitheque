package datastores

import (
	"database/sql"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"

	. "github.com/tbellembois/gochimitheque/models"
)

// GetSuppliers return the suppliers matching the search criteria
func (db *SQLiteDataStore) GetSuppliers(p Dbselectparam) ([]Supplier, int, error) {

	var (
		err                              error
		suppliers                        []Supplier
		count                            int
		exactSearch, countSql, selectSql string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("GetSuppliers")

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")
	supplierTable := goqu.T("supplier")

	// Join, where.
	joinClause := dialect.From(
		supplierTable,
	).Where(
		goqu.I("supplier_label").Like(p.GetSearch()),
	)

	if countSql, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("supplier_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}
	if selectSql, selectArgs, err = joinClause.Select(
		goqu.I("supplier_id"),
		goqu.I("supplier_label"),
	).Order(
		goqu.L("INSTR(supplier_label, \"?\")", exactSearch).Asc(),
		goqu.C("supplier_label").Asc(),
	).Limit(uint(p.GetLimit())).Offset(uint(p.GetOffset())).ToSQL(); err != nil {
		return nil, 0, err
	}

	// select
	if err = db.Select(&suppliers, selectSql, selectArgs...); err != nil {
		return nil, 0, err
	}
	// count
	if err = db.Get(&count, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for formula matching exactly the search
	sQuery := dialect.From(supplierTable).Where(
		goqu.I("supplier_label").Eq(exactSearch),
	).Select(
		"supplier_id",
		"supplier_label",
	)

	var (
		sqlr     string
		args     []interface{}
		supplier Supplier
	)
	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return nil, 0, err
	}

	if err = db.Select(&supplier, sqlr, args...); err != nil {
		return nil, 0, err
	}

	for i, e := range suppliers {
		if e.SupplierID == supplier.SupplierID {
			suppliers[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"suppliers": suppliers}).Debug("GetSuppliers")

	return suppliers, count, nil

}

// GetSupplier return the formula matching the given id
func (db *SQLiteDataStore) GetSupplier(id int) (Supplier, error) {

	var (
		err      error
		sqlr     string
		args     []interface{}
		supplier Supplier
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetSupplier")

	dialect := goqu.Dialect("sqlite3")
	supplierTable := goqu.T("supplier")

	sQuery := dialect.From(supplierTable).Where(
		goqu.I("supplier_id").Eq(id),
	).Select(
		goqu.I("supplier_id"),
		goqu.I("supplier_label"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return Supplier{}, err
	}

	if err = db.Get(&supplier, sqlr, args...); err != nil {
		return Supplier{}, err
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "supplier": supplier}).Debug("GetSupplier")

	return supplier, nil

}

// GetSupplierByLabel return the supplier matching the given supplier
func (db *SQLiteDataStore) GetSupplierByLabel(label string) (Supplier, error) {

	var (
		err      error
		sqlr     string
		args     []interface{}
		supplier Supplier
	)
	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetSupplierByLabel")

	dialect := goqu.Dialect("sqlite3")
	supplierTable := goqu.T("supplier")

	sQuery := dialect.From(supplierTable).Where(
		goqu.I("supplier_label").Eq(label),
	).Select(
		goqu.I("supplier_id"),
		goqu.I("supplier_label"),
	).Order(goqu.I("supplier_label").Asc())

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return Supplier{}, err
	}

	if err = db.Get(&supplier, sqlr, args...); err != nil {
		return Supplier{}, err
	}

	logger.Log.WithFields(logrus.Fields{"label": label, "supplier": supplier}).Debug("GetSupplierByLabel")

	return supplier, nil

}

// CreateSupplier create a new supplier in the db
func (db *SQLiteDataStore) CreateSupplier(s Supplier) (lastInsertId int64, err error) {

	var (
		sqlr string
		args []interface{}
		tx   *sql.Tx
		res  sql.Result
	)

	logger.Log.WithFields(logrus.Fields{"s": s}).Debug("CreateSupplier")

	dialect := goqu.Dialect("sqlite3")
	tableSupplier := goqu.T("supplier")

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

	iQuery := dialect.Insert(tableSupplier).Rows(
		goqu.Record{
			"supplier_label": s.SupplierLabel,
		},
	)

	if sqlr, args, err = iQuery.ToSQL(); err != nil {
		return
	}

	if res, err = tx.Exec(sqlr, args...); err != nil {
		return
	}

	if lastInsertId, err = res.LastInsertId(); err != nil {
		return
	}

	return

}

// GetSupplierRefs return the supplierrefs matching the search criteria
func (db *SQLiteDataStore) GetSupplierRefs(p DbselectparamSupplierRef) ([]SupplierRef, int, error) {

	var (
		err                              error
		supplierRefs                     []SupplierRef
		count                            int
		exactSearch, countSql, selectSql string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("GetSupplierRefs")

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")
	supplierrefTable := goqu.T("supplierref")

	// Join, where.
	whereAnd := []goqu.Expression{
		goqu.I("supplierref.supplierref_label").Like(p.GetSearch()),
	}
	if p.GetSupplier() != -1 {
		whereAnd = append(whereAnd, goqu.I("supplierref.supplier").Eq(p.GetSupplier()))
	}

	joinClause := dialect.From(
		supplierrefTable,
	).Join(
		goqu.T("supplier"),
		goqu.On(
			goqu.Ex{
				"supplierref.supplier": goqu.I("supplier.supplier_id"),
			},
		),
	).Where(
		whereAnd...,
	)

	if countSql, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("supplierref.supplierref_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}
	if selectSql, selectArgs, err = joinClause.Select(
		goqu.I("supplierref_id"),
		goqu.I("supplierref_label"),
		goqu.I("supplier_id").As("supplier.supplier_id"),
		goqu.I("supplier_label").As("supplier.supplier_label"),
	).Order(
		goqu.L("INSTR(supplierref_label, \"?\")", exactSearch).Asc(),
		goqu.C("supplierref_label").Asc(),
	).Limit(uint(p.GetLimit())).Offset(uint(p.GetOffset())).ToSQL(); err != nil {
		return nil, 0, err
	}

	// select
	if err = db.Select(&supplierRefs, selectSql, selectArgs...); err != nil {
		return nil, 0, err
	}
	// count
	if err = db.Get(&count, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for formula matching exactly the search
	sQuery := dialect.From(supplierrefTable).Where(
		goqu.I("supplierref_label").Eq(exactSearch),
	).Select(
		"supplierref_id",
		"supplierref_label",
	)

	var (
		sqlr string
		args []interface{}
		pref SupplierRef
	)
	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return nil, 0, err
	}

	if err = db.Get(&pref, sqlr, args...); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	for i, p := range supplierRefs {
		if p.SupplierRefID == pref.SupplierRefID {
			supplierRefs[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"supplierRefs": supplierRefs}).Debug("GetSupplierRefs")

	return supplierRefs, count, nil

}

package datastores

import (
	"database/sql"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/zmqclient"
)

func (db *SQLiteDataStore) GetSuppliers(f zmqclient.RequestFilter) ([]models.Supplier, int, error) {
	var (
		err                              error
		suppliers                        []models.Supplier
		count                            int
		exactSearch, countSQL, selectSQL string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"f": f}).Debug("GetSuppliers")

	exactSearch = f.Search
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")
	supplierTable := goqu.T("supplier")

	// Join, where.
	joinClause := dialect.From(
		supplierTable,
	).Where(
		goqu.I("supplier_label").Like(f.Search),
	)

	if countSQL, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("supplier_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}
	if selectSQL, selectArgs, err = joinClause.Select(
		goqu.I("supplier_id"),
		goqu.I("supplier_label"),
	).Order(
		goqu.L("INSTR(supplier_label, ?)", exactSearch).Asc(),
		goqu.C("supplier_label").Asc(),
	).Limit(uint(f.Limit)).Offset(uint(f.Offset)).ToSQL(); err != nil {
		return nil, 0, err
	}

	// Select.
	if err = db.Select(&suppliers, selectSQL, selectArgs...); err != nil {
		return nil, 0, err
	}
	// Count.
	if err = db.Get(&count, countSQL, countArgs...); err != nil {
		return nil, 0, err
	}

	// Setting the C attribute for formula matching exactly the search.
	sQuery := dialect.From(supplierTable).Where(
		goqu.I("supplier_label").Eq(exactSearch),
	).Select(
		"supplier_id",
		"supplier_label",
	)

	var (
		sqlr     string
		args     []interface{}
		supplier models.Supplier
	)
	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return nil, 0, err
	}

	if err = db.Get(&supplier, sqlr, args...); err != nil && err != sql.ErrNoRows {
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

func (db *SQLiteDataStore) GetSupplier(id int) (models.Supplier, error) {
	var (
		err      error
		sqlr     string
		args     []interface{}
		supplier models.Supplier
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
		return models.Supplier{}, err
	}

	if err = db.Get(&supplier, sqlr, args...); err != nil {
		return models.Supplier{}, err
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "supplier": supplier}).Debug("GetSupplier")

	return supplier, nil
}

func (db *SQLiteDataStore) GetSupplierByLabel(label string) (models.Supplier, error) {
	var (
		err      error
		sqlr     string
		args     []interface{}
		supplier models.Supplier
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
		return models.Supplier{}, err
	}

	if err = db.Get(&supplier, sqlr, args...); err != nil {
		return models.Supplier{}, err
	}

	logger.Log.WithFields(logrus.Fields{"label": label, "supplier": supplier}).Debug("GetSupplierByLabel")

	return supplier, nil
}

func (db *SQLiteDataStore) CreateSupplier(s models.Supplier) (lastInsertID int64, err error) {
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

	if lastInsertID, err = res.LastInsertId(); err != nil {
		return
	}

	return
}

func (db *SQLiteDataStore) GetSupplierRefs(f zmqclient.RequestFilter) ([]models.SupplierRef, int, error) {
	var (
		err                              error
		supplierRefs                     []models.SupplierRef
		count                            int
		exactSearch, countSQL, selectSQL string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"f": f}).Debug("GetSupplierRefs")

	if f.OrderBy == "" {
		f.OrderBy = "supplierref_id"
	}

	exactSearch = f.Search
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")
	supplierrefTable := goqu.T("supplierref")

	// Join, where.
	whereAnd := []goqu.Expression{
		goqu.I("supplierref.supplierref_label").Like(f.Search),
	}
	if f.Supplier != 0 {
		whereAnd = append(whereAnd, goqu.I("supplierref.supplier").Eq(f.Supplier))
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

	if countSQL, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("supplierref.supplierref_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}
	if selectSQL, selectArgs, err = joinClause.Select(
		goqu.I("supplierref_id"),
		goqu.I("supplierref_label"),
		goqu.I("supplier_id").As(goqu.C("supplier.supplier_id")),
		goqu.I("supplier_label").As(goqu.C("supplier.supplier_label")),
	).Order(
		goqu.L("INSTR(supplierref_label, ?)", exactSearch).Asc(),
		goqu.C("supplierref_label").Asc(),
	).Limit(uint(f.Limit)).Offset(uint(f.Offset)).ToSQL(); err != nil {
		return nil, 0, err
	}

	// Select.
	if err = db.Select(&supplierRefs, selectSQL, selectArgs...); err != nil {
		return nil, 0, err
	}
	// Count.
	if err = db.Get(&count, countSQL, countArgs...); err != nil {
		return nil, 0, err
	}

	// Setting the C attribute for formula matching exactly the search.
	sQuery := dialect.From(supplierrefTable).Where(
		goqu.I("supplierref_label").Eq(exactSearch),
	).Select(
		"supplierref_id",
		"supplierref_label",
	)

	var (
		sqlr string
		args []interface{}
		pref models.SupplierRef
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

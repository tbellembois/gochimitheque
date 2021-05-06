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
	. "github.com/tbellembois/gochimitheque/models"
)

// Return the store location full path.
// The caller is responsible of opening and commiting the tx transaction.
func (db *SQLiteDataStore) buildFullPath(s StoreLocation, tx *sqlx.Tx) string {

	var (
		err    error
		sqlr   string
		args   []interface{}
		parent StoreLocation
	)

	logger.Log.WithFields(logrus.Fields{"s": s}).Debug("buildFullPath")

	// Recursively getting the parents.
	if s.StoreLocation != nil && s.StoreLocation.StoreLocationID.Valid {

		dialect := goqu.Dialect("sqlite3")
		tableStorelocation := goqu.T("storelocation")

		sQuery := dialect.From(tableStorelocation.As("s")).Select(
			goqu.I("s.storelocation_id"),
			goqu.I("s.storelocation_name"),
			goqu.I("storelocation.storelocation_id").As(goqu.C("storelocation.storelocation_id")),
			goqu.I("storelocation.storelocation_name").As(goqu.C("storelocation.storelocation_name")),
		).LeftJoin(
			goqu.T("storelocation"),
			goqu.On(goqu.Ex{
				"s.storelocation": goqu.I("storelocation.storelocation_id"),
			}),
		).Where(
			goqu.I("s.storelocation_id").Eq(s.StoreLocation.StoreLocationID.Int64),
		)

		if sqlr, args, err = sQuery.ToSQL(); err != nil {
			logger.Log.Error(err)
			return ""
		}

		if err = tx.Get(&parent, sqlr, args...); err != nil {
			logger.Log.Error(err)
			return ""
		}

		return db.buildFullPath(parent, tx) + "/" + s.StoreLocationName.String

	}

	return s.StoreLocationName.String

}

// GetStoreLocations select the store locations matching p
// and visible by the connected user.
func (db *SQLiteDataStore) GetStoreLocations(p DbselectparamStoreLocation) ([]StoreLocation, int, error) {

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("GetStoreLocations")

	var err error

	dialect := goqu.Dialect("sqlite3")
	tableStorelocation := goqu.T("storelocation")

	// Map orderby clause.
	orderByClause := p.GetOrderBy()
	if orderByClause == "storelocation" {
		orderByClause = "storelocation.storelocation_id"
	}

	// Build orderby/order clause.
	orderClause := goqu.I(orderByClause).Asc()
	if strings.ToLower(p.GetOrder()) == "desc" {
		orderClause = goqu.I(orderByClause).Desc()
	}

	// Build join clause.
	joinClause := dialect.From(tableStorelocation.As("s")).Join(
		goqu.T("entity"),
		goqu.On(goqu.Ex{"s.entity": goqu.I("entity.entity_id")}),
	).LeftJoin(
		goqu.T("storelocation"),
		goqu.On(goqu.Ex{"s.storelocation": goqu.I("storelocation.storelocation_id")}),
	).Join(
		goqu.T("permission").As("perm"),
		goqu.On(
			goqu.Ex{
				"perm.person":               p.GetLoggedPersonID(),
				"perm.permission_item_name": []string{"all", "storages"},
				"perm.permission_perm_name": []string{"r", "w", "all"},
				"perm.permission_entity_id": []interface{}{-1, goqu.I("entity.entity_id")},
			},
		),
	)

	// Build where AND expression.
	whereAnd := []goqu.Expression{
		goqu.I("s.storelocation_name").Like(p.GetSearch()),
	}
	if p.GetEntity() != -1 {
		whereAnd = append(whereAnd, goqu.I("s.entity").Eq(p.GetEntity()))
	}
	if p.GetStoreLocationCanStore() {
		whereAnd = append(whereAnd, goqu.I("s.storelocation_canstore").Eq(p.GetStoreLocationCanStore()))
	}

	joinClause = joinClause.Where(goqu.And(whereAnd...))

	// Building final count.
	var (
		countSql  string
		countArgs []interface{}
	)
	if countSql, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("s.storelocation_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}

	// Building final select.
	var (
		selectSql  string
		selectArgs []interface{}
	)
	if selectSql, selectArgs, err = joinClause.Select(
		goqu.I("s.storelocation_id").As("storelocation_id"),
		goqu.I("s.storelocation_canstore").As("storelocation_canstore"),
		goqu.I("s.storelocation_color").As("storelocation_color"),
		goqu.I("s.storelocation_id").As("storelocation_id"),
		goqu.I("s.storelocation_name").As("storelocation_name"),
		goqu.I("s.storelocation_fullpath").As("storelocation_fullpath"),
		goqu.I("storelocation.storelocation_id").As(goqu.C("storelocation.storelocation_id")),
		goqu.I("storelocation.storelocation_name").As(goqu.C("storelocation.storelocation_name")),
		goqu.I("entity.entity_id").As(goqu.C("entity.entity_id")),
		goqu.I("entity.entity_name").As(goqu.C("entity.entity_name")),
	).GroupBy(goqu.I("s.storelocation_id")).Order(orderClause).Limit(uint(p.GetLimit())).Offset(uint(p.GetOffset())).ToSQL(); err != nil {
		return nil, 0, err
	}

	// logger.Log.Debug(selectSql)
	// logger.Log.Debug(selectArgs)
	// logger.Log.Debug(countSql)
	// logger.Log.Debug(countArgs)

	var (
		storelocations []StoreLocation
		count          int
	)

	if err = db.Select(&storelocations, selectSql, selectArgs...); err != nil {
		return nil, 0, err
	}

	if err = db.Get(&count, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	return storelocations, count, nil

}

// GetStoreLocation select the store location by id.
func (db *SQLiteDataStore) GetStoreLocation(id int) (StoreLocation, error) {

	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetStoreLocation")

	dialect := goqu.Dialect("sqlite3")
	tableStorelocation := goqu.T("storelocation")

	sQuery := dialect.From(tableStorelocation.As("s")).Join(
		goqu.T("entity"),
		goqu.On(goqu.Ex{"s.entity": goqu.I("entity.entity_id")}),
	).LeftJoin(
		goqu.T("storelocation"),
		goqu.On(goqu.Ex{"s.storelocation": goqu.I("storelocation.storelocation_id")}),
	).Where(
		goqu.I("s.storelocation_id").Eq(id),
	).Select(
		goqu.I("s.storelocation_id"),
		goqu.I("s.storelocation_name"),
		goqu.I("s.storelocation_canstore"),
		goqu.I("s.storelocation_color"),
		goqu.I("s.storelocation_fullpath"),
		goqu.I("storelocation.storelocation_id").As(goqu.C("storelocation.storelocation_id")),
		goqu.I("storelocation.storelocation_name").As(goqu.C("storelocation.storelocation_name")),
		goqu.I("entity.entity_id").As(goqu.C("entity.entity_id")),
		goqu.I("entity.entity_name").As(goqu.C("entity.entity_name")),
	)

	var (
		err           error
		sqlr          string
		args          []interface{}
		storelocation StoreLocation
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return StoreLocation{}, err
	}

	if err = db.Get(&storelocation, sqlr, args...); err != nil {
		return StoreLocation{}, err
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "storelocation": storelocation}).Debug("GetStoreLocation")

	return storelocation, nil

}

// GetStoreLocationChildren select the children store locations of parent id.
func (db *SQLiteDataStore) GetStoreLocationChildren(id int) ([]StoreLocation, error) {

	dialect := goqu.Dialect("sqlite3")
	tableStorelocation := goqu.T("storelocation")

	// Select
	sQuery := dialect.From(tableStorelocation.As("s")).Select(
		goqu.I("s.storelocation_id"),
		goqu.I("s.storelocation_name"),
		goqu.I("s.storelocation_canstore"),
		goqu.I("s.storelocation_color"),
		goqu.I("s.storelocation_fullpath"),
		goqu.I("storelocation.storelocation_id").As(goqu.C("storelocation.storelocation_id")),
		goqu.I("storelocation.storelocation_name").As(goqu.C("storelocation.storelocation_name")),
		goqu.I("entity.entity_id").As(goqu.C("entity.entity_id")),
		goqu.I("entity.entity_name").As(goqu.C("entity.entity_name")),
	).Join(
		goqu.T("entity"),
		goqu.On(goqu.Ex{"s.entity": goqu.I("entity.entity_id")}),
	).LeftJoin(
		goqu.T("storelocation"),
		goqu.On(goqu.Ex{"s.storelocation": goqu.I("storelocation.storelocation_id")}),
	).Where(
		goqu.I("s.storelocation").Eq(id),
	)

	var (
		err            error
		sqlr           string
		args           []interface{}
		storelocations []StoreLocation
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return nil, err
	}

	if err = db.Select(&storelocations, sqlr, args...); err != nil {
		return nil, err
	}

	return storelocations, nil

}

// DeleteStoreLocation delete the store location by id.
func (db *SQLiteDataStore) DeleteStoreLocation(id int) error {

	dialect := goqu.Dialect("sqlite3")
	tableStorelocation := goqu.T("storelocation")

	dQuery := dialect.From(tableStorelocation).Where(
		goqu.I("storelocation_id").Eq(id),
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

// CreateStoreLocation insert s.
func (db *SQLiteDataStore) CreateStoreLocation(s StoreLocation) (lastInsertId int64, err error) {

	var (
		tx *sqlx.Tx
	)

	logger.Log.WithFields(logrus.Fields{"s": fmt.Sprintf("%+v", s)}).Debug("CreateStoreLocation")

	dialect := goqu.Dialect("sqlite3")
	tableStorelocation := goqu.T("storelocation")

	if tx, err = db.Beginx(); err != nil {
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

	s.StoreLocationFullPath = db.buildFullPath(s, tx)

	iQuery := dialect.Insert(tableStorelocation)

	setClause := goqu.Record{
		"storelocation_name":     s.StoreLocationName.String,
		"entity":                 s.EntityID,
		"storelocation_fullpath": s.StoreLocationFullPath,
	}

	if s.StoreLocationCanStore.Valid {
		setClause["storelocation_canstore"] = s.StoreLocationCanStore.Bool
	}
	if s.StoreLocationColor.Valid {
		setClause["storelocation_color"] = s.StoreLocationColor.String
	}
	if s.StoreLocation != nil {
		setClause["storelocation"] = s.StoreLocation.StoreLocationID.Int64
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

	return sqlResult.LastInsertId()

}

// UpdateStoreLocation updates s.
func (db *SQLiteDataStore) UpdateStoreLocation(s StoreLocation) (err error) {

	var (
		tx *sqlx.Tx
	)

	dialect := goqu.Dialect("sqlite3")
	tableStorelocation := goqu.T("storelocation")

	if tx, err = db.Beginx(); err != nil {
		return
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

	s.StoreLocationFullPath = db.buildFullPath(s, tx)

	uQuery := dialect.Update(tableStorelocation)

	setClause := goqu.Record{
		"storelocation_name":     s.StoreLocationName.String,
		"entity":                 s.EntityID,
		"storelocation_fullpath": s.StoreLocationFullPath,
	}

	if s.StoreLocationCanStore.Valid {
		setClause["storelocation_canstore"] = s.StoreLocationCanStore.Bool
	}
	if s.StoreLocationColor.Valid {
		setClause["storelocation_color"] = s.StoreLocationColor.String
	}
	if s.StoreLocation != nil {
		setClause["storelocation"] = s.StoreLocation.StoreLocationID.Int64
	}

	var (
		sqlr string
		args []interface{}
	)

	if sqlr, args, err = uQuery.Set(
		setClause,
	).Where(
		goqu.I("storelocation_id").Eq(s.StoreLocationID),
	).ToSQL(); err != nil {
		logger.Log.Error(err)
		return
	}

	if _, err = tx.Exec(sqlr, args...); err != nil {
		return
	}

	return nil

}

// HasStorelocationStorage returns true is the store location is empty.
func (db *SQLiteDataStore) HasStorelocationStorage(id int) (bool, error) {

	var (
		err   error
		sqlr  string
		args  []interface{}
		count int
	)

	dialect := goqu.Dialect("sqlite3")
	tableStorage := goqu.T("storage")

	sQuery := dialect.From(tableStorage).Select(
		goqu.COUNT("*"),
	).Where(
		goqu.I("storelocation").Eq(id),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return false, err
	}

	if err = db.Get(&count, sqlr, args...); err != nil {
		return false, err
	}

	return count == 0, nil

}

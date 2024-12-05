package datastores

import (
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
)

func (db *SQLiteDataStore) GetStoreLocationChildren(id int) ([]models.StoreLocation, error) {
	dialect := goqu.Dialect("sqlite3")
	tableStorelocation := goqu.T("store_location")

	// Select
	sQuery := dialect.From(tableStorelocation.As("s")).Select(
		goqu.I("s.store_location_id"),
		goqu.I("s.store_location_name"),
		goqu.I("s.store_location_can_store"),
		goqu.I("s.store_location_color"),
		goqu.I("s.store_location_full_path"),
		goqu.I("store_location.store_location_id").As(goqu.C("store_location.store_location_id")),
		goqu.I("store_location.store_location_name").As(goqu.C("store_location.store_location_name")),
		goqu.I("entity.entity_id").As(goqu.C("entity.entity_id")),
		goqu.I("entity.entity_name").As(goqu.C("entity.entity_name")),
	).Join(
		goqu.T("entity"),
		goqu.On(goqu.Ex{"s.entity": goqu.I("entity.entity_id")}),
	).LeftJoin(
		goqu.T("store_location"),
		goqu.On(goqu.Ex{"s.store_location": goqu.I("store_location.store_location_id")}),
	).Where(
		goqu.I("s.store_location").Eq(id),
	)

	var (
		err             error
		sqlr            string
		args            []interface{}
		store_locations []models.StoreLocation
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return nil, err
	}

	if err = db.Select(&store_locations, sqlr, args...); err != nil {
		return nil, err
	}

	return store_locations, nil
}

func (db *SQLiteDataStore) DeleteStoreLocation(id int) error {
	dialect := goqu.Dialect("sqlite3")
	tableStorelocation := goqu.T("store_location")

	dQuery := dialect.From(tableStorelocation).Where(
		goqu.I("store_location_id").Eq(id),
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
		logger.Log.Error(err)
		return err
	}

	return nil
}

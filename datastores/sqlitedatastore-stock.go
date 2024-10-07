package datastores

import (
	"database/sql"
	"net/http"
	"sync"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/request"
	"github.com/tbellembois/gochimitheque/zmqclient"
)

type SyncStoreLocation struct {
	mu            sync.Mutex
	Storelocation *models.StoreLocation
}

// computeStockStorelocationConsumable returns the number of units of product p in the store location s.
func (db *SQLiteDataStore) computeStockStorelocationConsumable(p models.Product, s *SyncStoreLocation, mu *sync.Mutex) float64 {
	var (
		err                   error
		currentStock          float64
		totalStock            float64
		store_locationChildren []models.StoreLocation
		sqlr                  string
		args                  []interface{}
	)

	dialect := goqu.Dialect("sqlite3")
	t := goqu.T("storage")

	// Getting the store location current stock.
	sQuery := dialect.From(t).Join(
		goqu.T("product"),
		goqu.On(goqu.Ex{"storage.product": goqu.I("product.product_id")}),
	).Where(
		goqu.I("storage.storage_archive").IsFalse(),
		goqu.I("storage.storage").IsNull(),
		goqu.I("storage.store_location").Eq(s.Storelocation.StoreLocationID.Int64),
		goqu.I("storage.product").Eq(p.ProductID),
	).Select(
		goqu.SUM(goqu.L("product.product_number_per_bag * storage.storage_number_of_bag")).As("bag"),
		goqu.SUM(goqu.L("product.product_number_per_carton * storage.storage_number_of_carton")).As("carton"),
		goqu.SUM(goqu.L("storage.storage_number_of_unit")).As("unit"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return 0
	}

	type Result struct {
		Bag    sql.NullInt64 `db:"bag"`
		Carton sql.NullInt64 `db:"carton"`
		Unit   sql.NullInt64 `db:"unit"`
	}

	var result Result

	mu.Lock()
	if err = db.Get(&result, sqlr, args...); err != nil && err != sql.ErrNoRows {
		logger.Log.Error(err)
		return 0
	}
	mu.Unlock()

	var stock int64
	if result.Bag.Valid && result.Bag.Int64 > 0 {
		stock = result.Bag.Int64
	}
	if result.Carton.Valid && result.Carton.Int64 > 0 {
		stock += result.Carton.Int64
	}
	if result.Unit.Valid && result.Unit.Int64 > 0 {
		stock += result.Unit.Int64
	}

	// totalStock is initialized with currentStock
	// and increased later while processing the children.
	currentStock = float64(stock)
	totalStock = float64(stock)

	logger.Log.WithFields(logrus.Fields{
		"p.ProductID":         p.ProductID,
		"s.StoreLocationName": s.Storelocation.StoreLocationName,
		"currentStock":        currentStock,
	}).Debug("computeStockStorelocationConsumable")

	// Getting the children store locations.
	mu.Lock()
	if store_locationChildren, err = db.GetStoreLocationChildren(int(s.Storelocation.StoreLocationID.Int64)); err != nil {
		logger.Log.Error(err)
		return 0
	}
	mu.Unlock()

	for i := range store_locationChildren {
		s.Storelocation.Children = append(s.Storelocation.Children, &store_locationChildren[i])

		totalStock += db.computeStockStorelocationConsumable(p, &SyncStoreLocation{
			Storelocation: &store_locationChildren[i],
		}, mu)
	}

	s.mu.Lock()
	s.Storelocation.Stocks = append(s.Storelocation.Stocks, models.Stock{Total: totalStock, Current: currentStock})
	s.mu.Unlock()

	return currentStock
}

// computeStockStorelocation returns the quantity of product p in the store location s for the unit u.
func (db *SQLiteDataStore) computeStockStorelocation(p models.Product, s *SyncStoreLocation, u models.Unit, mu *sync.Mutex) float64 {
	var (
		err                   error
		currentStock          float64
		totalStock            float64
		store_locationChildren []models.StoreLocation
		sqlr                  string
		args                  []interface{}
	)

	dialect := goqu.Dialect("sqlite3")
	t := goqu.T("storage")

	// Getting the store location current stock.
	sQuery := dialect.From(t).Join(
		goqu.T("unit"),
		goqu.On(goqu.Ex{"storage.unit_quantity": goqu.I("unit.unit_id")}),
	).Where(
		goqu.I("storage.store_location").Eq(s.Storelocation.StoreLocationID.Int64),
		goqu.I("storage.storage").IsNull(),
		goqu.I("storage.storage_quantity").IsNotNull(),
		goqu.I("storage.storage_archive").IsFalse(),
		goqu.I("storage.product").Eq(p.ProductID),
		goqu.Or(
			goqu.I("storage.unit_quantity").Eq(*u.UnitID),
			goqu.I("unit.unit").Eq(*u.UnitID),
		),
	).Select(
		goqu.L("SUM(storage.storage_quantity * unit_multiplier)"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return 0
	}

	var nullableFloat64 sql.NullFloat64

	mu.Lock()
	if err = db.Get(&nullableFloat64, sqlr, args...); err != nil && err != sql.ErrNoRows {
		logger.Log.Error(err)
		return 0
	}
	mu.Unlock()

	// totalStock is initialized with currentStock
	// and increased later while processing the children.
	if nullableFloat64.Valid {
		currentStock = nullableFloat64.Float64
		totalStock = nullableFloat64.Float64
	}

	logger.Log.WithFields(logrus.Fields{
		"p.ProductID":         p.ProductID,
		"s.StoreLocationName": s.Storelocation.StoreLocationName,
		"currentStock":        currentStock,
	}).Debug("computeStockStorelocation")

	mu.Lock()
	if store_locationChildren, err = db.GetStoreLocationChildren(int(s.Storelocation.StoreLocationID.Int64)); err != nil {
		logger.Log.Error(err)
		return 0
	}
	mu.Unlock()

	for i := range store_locationChildren {
		s.Storelocation.Children = append(s.Storelocation.Children, &store_locationChildren[i])

		totalStock += db.computeStockStorelocation(p, &SyncStoreLocation{
			Storelocation: &store_locationChildren[i],
		}, u, mu)
	}

	s.mu.Lock()
	s.Storelocation.Stocks = append(s.Storelocation.Stocks, models.Stock{Total: totalStock, Current: currentStock, Unit: u})
	s.mu.Unlock()

	return currentStock
}

// computeStockStorelocationNoUnit returns the quantity of product p with no unit in the store location s.
func (db *SQLiteDataStore) computeStockStorelocationNoUnit(p models.Product, s *SyncStoreLocation, mu *sync.Mutex) float64 {
	var (
		currentStock          float64
		totalStock            float64
		store_locationChildren []models.StoreLocation
		err                   error
		sqlrNotNull, sqlrNull string
		argsNotNull, argsNull []interface{}
	)

	dialect := goqu.Dialect("sqlite3")
	t := goqu.T("storage")

	// Getting the store location current stock.
	sQueryNotNull := dialect.From(t).Where(
		goqu.I("storage.store_location").Eq(s.Storelocation.StoreLocationID.Int64),
		goqu.I("storage.storage").IsNull(),
		goqu.And(
			goqu.I("storage.storage_quantity").IsNotNull(),
			goqu.I("storage.storage_quantity").Neq(0),
		),
		goqu.I("storage.storage_archive").IsFalse(),
		goqu.I("storage.product").Eq(p.ProductID),
		goqu.I("storage.unit_quantity").IsNull(),
	).Select(
		goqu.SUM(goqu.I("storage.storage_quantity")),
	)

	if sqlrNotNull, argsNotNull, err = sQueryNotNull.ToSQL(); err != nil {
		logger.Log.Error(err)
		return 0
	}

	sQueryNull := dialect.From(t).Where(
		goqu.I("storage.store_location").Eq(s.Storelocation.StoreLocationID.Int64),
		goqu.I("storage.storage").IsNull(),
		goqu.Or(
			goqu.I("storage.storage_quantity").IsNull(),
			goqu.I("storage.storage_quantity").Eq(0),
		),
		goqu.I("storage.storage_archive").IsFalse(),
		goqu.I("storage.product").Eq(p.ProductID),
		goqu.I("storage.unit_quantity").IsNull(),
	).Select(
		goqu.COUNT(goqu.I("storage.storage_id").Distinct()),
	)

	if sqlrNull, argsNull, err = sQueryNull.ToSQL(); err != nil {
		logger.Log.Error(err)
		return 0
	}

	mu.Lock()

	var (
		resultNull, resultNotNull, nullableFloat64 sql.NullFloat64
	)

	if err = db.Get(&resultNotNull, sqlrNotNull, argsNotNull...); err != nil && err != sql.ErrNoRows {
		logger.Log.Error(err)
		return 0
	}

	if err = db.Get(&resultNull, sqlrNull, argsNull...); err != nil && err != sql.ErrNoRows {
		logger.Log.Error(err)
		return 0
	}

	nullableFloat64 = sql.NullFloat64{Valid: true, Float64: resultNotNull.Float64 + resultNull.Float64}

	mu.Unlock()

	// totalStock is initialized with currentStock
	// and increased later while processing the children.
	if nullableFloat64.Valid {
		currentStock = nullableFloat64.Float64
		totalStock = nullableFloat64.Float64
	}

	logger.Log.WithFields(logrus.Fields{
		"p.ProductID":         p.ProductID,
		"s.StoreLocationName": s.Storelocation.StoreLocationName,
		"currentStock":        currentStock,
	}).Debug("computeStockStorelocationNoUnit")

	// Getting the children store locations.
	mu.Lock()
	if store_locationChildren, err = db.GetStoreLocationChildren(int(s.Storelocation.StoreLocationID.Int64)); err != nil {
		logger.Log.Error(err)
		return 0
	}
	mu.Unlock()

	for i := range store_locationChildren {
		s.Storelocation.Children = append(s.Storelocation.Children, &store_locationChildren[i])

		totalStock += db.computeStockStorelocationNoUnit(p, &SyncStoreLocation{
			Storelocation: &store_locationChildren[i],
		}, mu)
	}

	s.mu.Lock()
	s.Storelocation.Stocks = append(s.Storelocation.Stocks, models.Stock{Total: totalStock, Current: currentStock, Unit: models.Unit{}})
	s.mu.Unlock()

	return currentStock
}

// ComputeStockEntity returns the root store locations of the entity(ies) of the loggued user.
// Each store location has a Stocks []models.Stock field containing the stocks of the product p for each unit.
func (db *SQLiteDataStore) ComputeStockEntity(p models.Product, r *http.Request) []models.StoreLocation {
	var (
		units              []models.Unit // reference units
		syncstore_locations []SyncStoreLocation
		entities           []models.Entity
		eids               []int
		err                error
		sqlr               string
		args               []interface{}
	)

	// Getting the entities (GetEntities returns only entities the connected user can see).
	var (
		filter zmqclient.RequestFilter
	)

	c := request.ContainerFromRequestContext(r)

	// if filter, aerr = request.NewFilter(r); aerr != nil {
	// 	logger.Log.Error(aerr.Error())
	// 	return []models.StoreLocation{}
	// }
	if filter, err = zmqclient.RequestFilterFromRawString("http://localhost/?" + r.URL.RawQuery); err != nil {
		logger.Log.Error(err)
		return []models.StoreLocation{}
	}

	if entities, _, err = db.GetEntities(filter, c.PersonID); err != nil {
		logger.Log.Error(err)
		return []models.StoreLocation{}
	}

	for _, e := range entities {
		eids = append(eids, e.EntityID)
	}

	dialect := goqu.Dialect("sqlite3")

	// Getting the reference units.
	t := goqu.T("unit")
	if sqlr, args, err = dialect.From(t).Where(
		goqu.I("unit.unit").IsNull(),
		goqu.I("unit.unit_type").Eq("quantity"),
	).Select(
		goqu.I("unit.unit_id"),
		goqu.I("unit.unit_label"),
	).ToSQL(); err != nil {
		logger.Log.Error(err)
		return []models.StoreLocation{}
	}

	if err = db.Select(&units, sqlr, args...); err != nil {
		logger.Log.Error(err)
		return []models.StoreLocation{}
	}

	// Getting the root store locations.
	t = goqu.T("store_location")
	sQuery := dialect.From(t).Where(
		goqu.I("store_location.store_location").IsNull(),
		goqu.I("store_location.entity").In(eids),
	).Select(
		goqu.I("store_location.store_location_id"),
		goqu.I("store_location.store_location_name"),
		goqu.I("store_location.store_location_color"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return []models.StoreLocation{}
	}

	var rootStoreLocations []models.StoreLocation
	if err = db.Select(&rootStoreLocations, sqlr, args...); err != nil {
		logger.Log.Error(err)
		return []models.StoreLocation{}
	}

	for i := range rootStoreLocations {
		syncstore_locations = append(syncstore_locations, SyncStoreLocation{
			Storelocation: &rootStoreLocations[i],
		})
	}

	var wg sync.WaitGroup

	mu := &sync.Mutex{}

	// Computing stocks for storages with units.
	for i := range syncstore_locations {
		for j := range units {
			wg.Add(1)

			go func(u models.Unit, sl *SyncStoreLocation) {
				db.computeStockStorelocation(p, sl, u, mu)
				wg.Done()
			}(units[j], &syncstore_locations[i])
		}
	}
	// Computing stocks for storages without units.
	for i := range syncstore_locations {
		wg.Add(1)

		go func(sl *SyncStoreLocation) {
			db.computeStockStorelocationNoUnit(p, sl, mu)
			wg.Done()
		}(&syncstore_locations[i])
	}
	// Computing stocks for consumables storages.
	for i := range syncstore_locations {
		wg.Add(1)

		go func(sl *SyncStoreLocation) {
			db.computeStockStorelocationConsumable(p, sl, mu)
			wg.Done()
		}(&syncstore_locations[i])
	}

	wg.Wait()

	var result []models.StoreLocation
	for i := range syncstore_locations {
		result = append(result, *syncstore_locations[i].Storelocation)
	}

	return result
}

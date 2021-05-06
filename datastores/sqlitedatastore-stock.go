package datastores

import (
	"database/sql"
	"net/http"
	"sync"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	. "github.com/tbellembois/gochimitheque/models"
)

type SyncStoreLocation struct {
	mu            sync.Mutex
	Storelocation *StoreLocation
}

// computeStockStorelocationConsumable returns the number of units of product p in the store location s.
func (db *SQLiteDataStore) computeStockStorelocationConsumable(p Product, s *SyncStoreLocation, mu *sync.Mutex) float64 {

	var (
		err                   error
		currentStock          float64
		totalStock            float64
		storelocationChildren []StoreLocation
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
		goqu.I("storage.storelocation").Eq(s.Storelocation.StoreLocationID.Int64),
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
		stock = stock + result.Carton.Int64
	}
	if result.Unit.Valid && result.Unit.Int64 > 0 {
		stock = stock + result.Unit.Int64
	}

	// totalStock is initialized with currentStock
	// and increased later while processing the children.
	currentStock = float64(stock)
	totalStock = float64(stock)

	logger.Log.WithFields(logrus.Fields{
		"p.ProductID":         p.ProductID,
		"s.StoreLocationName": s.Storelocation.StoreLocationName,
		"currentStock":        currentStock}).Debug("ComputeStockStorelocation")

	// Getting the children store locations.
	mu.Lock()
	if storelocationChildren, err = db.GetStoreLocationChildren(int(s.Storelocation.StoreLocationID.Int64)); err != nil {
		logger.Log.Error(err)
		return 0
	}
	mu.Unlock()

	for i := range storelocationChildren {

		s.Storelocation.Children = append(s.Storelocation.Children, &storelocationChildren[i])

		totalStock += db.computeStockStorelocationConsumable(p, &SyncStoreLocation{
			Storelocation: &storelocationChildren[i],
		}, mu)

	}

	s.mu.Lock()
	(*s).Storelocation.Stocks = append((*s).Storelocation.Stocks, Stock{Total: totalStock, Current: currentStock})
	s.mu.Unlock()

	return currentStock

}

// computeStockStorelocation returns the quantity of product p in the store location s for the unit u.
func (db *SQLiteDataStore) computeStockStorelocation(p Product, s *SyncStoreLocation, u Unit, mu *sync.Mutex) float64 {

	var (
		err                   error
		currentStock          float64
		totalStock            float64
		storelocationChildren []StoreLocation
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
		goqu.I("storage.storelocation").Eq(s.Storelocation.StoreLocationID.Int64),
		goqu.I("storage.storage").IsNull(),
		goqu.I("storage.storage_quantity").IsNotNull(),
		goqu.I("storage.storage_archive").IsFalse(),
		goqu.I("storage.product").Eq(p.ProductID),
		goqu.Or(
			goqu.I("storage.unit_quantity").Eq(u.UnitID.Int64),
			goqu.I("storage.unit_quantity").In(dialect.From("unit").Select("unit_id").Where(goqu.I("unit.unit").Eq(u.UnitID.Int64))),
		),
	).Select(
		goqu.SUM(goqu.L("storage.storage_quantity * unit_multiplier")),
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

	mu.Lock()
	if storelocationChildren, err = db.GetStoreLocationChildren(int(s.Storelocation.StoreLocationID.Int64)); err != nil {
		logger.Log.Error(err)
		return 0
	}
	mu.Unlock()

	for i := range storelocationChildren {

		s.Storelocation.Children = append(s.Storelocation.Children, &storelocationChildren[i])

		totalStock += db.computeStockStorelocation(p, &SyncStoreLocation{
			Storelocation: &storelocationChildren[i],
		}, u, mu)

	}

	s.mu.Lock()
	(*s).Storelocation.Stocks = append((*s).Storelocation.Stocks, Stock{Total: totalStock, Current: currentStock, Unit: u})
	s.mu.Unlock()

	return currentStock

}

// computeStockStorelocationNoUnit returns the quantity of product p with no unit in the store location s.
func (db *SQLiteDataStore) computeStockStorelocationNoUnit(p Product, s *SyncStoreLocation, mu *sync.Mutex) float64 {

	var (
		currentStock          float64
		totalStock            float64
		storelocationChildren []StoreLocation
		err                   error
		sqlr                  string
		args                  []interface{}
	)

	dialect := goqu.Dialect("sqlite3")
	t := goqu.T("storage")

	// Getting the store location current stock.
	sQuery := dialect.From(t).LeftJoin(
		goqu.T("unit"),
		goqu.On(goqu.Ex{"storage.unit_quantity": goqu.I("unit.unit_id")}),
	).Where(
		goqu.I("storage.storelocation").Eq(s.Storelocation.StoreLocationID.Int64),
		goqu.I("storage.storage").IsNull(),
		goqu.I("storage.storage_quantity").IsNotNull(),
		goqu.I("storage.storage_archive").IsFalse(),
		goqu.I("storage.product").Eq(p.ProductID),
		goqu.I("storage.unit_quantity").IsNull(),
	).Select(
		goqu.SUM(goqu.I("storage.storage_quantity")),
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
		"currentStock":        currentStock}).Debug("ComputeStockStorelocation")

	// Getting the children store locations.
	mu.Lock()
	if storelocationChildren, err = db.GetStoreLocationChildren(int(s.Storelocation.StoreLocationID.Int64)); err != nil {
		logger.Log.Error(err)
		return 0
	}
	mu.Unlock()

	for i := range storelocationChildren {

		s.Storelocation.Children = append(s.Storelocation.Children, &storelocationChildren[i])

		totalStock += db.computeStockStorelocationNoUnit(p, &SyncStoreLocation{
			Storelocation: &storelocationChildren[i],
		}, mu)

	}

	s.mu.Lock()
	(*s).Storelocation.Stocks = append((*s).Storelocation.Stocks, Stock{Total: totalStock, Current: currentStock, Unit: Unit{}})
	s.mu.Unlock()

	return currentStock

}

// ComputeStockEntity returns the root store locations of the entity(ies) of the loggued user.
// Each store location has a Stocks []Stock field containing the stocks of the product p for each unit.
func (db *SQLiteDataStore) ComputeStockEntity(p Product, r *http.Request) []StoreLocation {

	var (
		units              []Unit // reference units
		syncstorelocations []SyncStoreLocation
		entities           []Entity
		eids               []int
		err                error
		sqlr               string
		args               []interface{}
	)

	// Getting the entities (GetEntities returns only entities the connected user can see).
	h, _ := NewdbselectparamEntity(r, nil)
	if entities, _, err = db.GetEntities(h); err != nil {
		logger.Log.Error(err)
		return []StoreLocation{}
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
		return []StoreLocation{}
	}

	if err = db.Select(&units, sqlr, args...); err != nil {
		logger.Log.Error(err)
		return []StoreLocation{}
	}

	// Getting the root store locations.
	t = goqu.T("storelocation")
	sQuery := dialect.From(t).Where(
		goqu.I("storelocation.storelocation").IsNull(),
		goqu.I("storelocation.entity").In(eids),
	).Select(
		goqu.I("storelocation.storelocation_id"),
		goqu.I("storelocation.storelocation_name"),
		goqu.I("storelocation.storelocation_color"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return []StoreLocation{}
	}

	var rootStoreLocations []StoreLocation
	if err = db.Select(&rootStoreLocations, sqlr, args...); err != nil {
		logger.Log.Error(err)
		return []StoreLocation{}
	}

	for i := range rootStoreLocations {
		syncstorelocations = append(syncstorelocations, SyncStoreLocation{
			Storelocation: &rootStoreLocations[i],
		})
	}

	var (
		wg sync.WaitGroup
	)
	mu := &sync.Mutex{}
	// Computing stocks for storages with units.
	for i := range syncstorelocations {
		for j := range units {
			wg.Add(1)
			go func(u Unit, sl *SyncStoreLocation) {
				db.computeStockStorelocation(p, sl, u, mu)
				wg.Done()
			}(units[j], &syncstorelocations[i])
		}
	}
	// Computing stocks for storages without units.
	for i := range syncstorelocations {
		wg.Add(1)
		go func(sl *SyncStoreLocation) {
			db.computeStockStorelocationNoUnit(p, sl, mu)
			wg.Done()
		}(&syncstorelocations[i])
	}
	// Computing stocks for consumables storages.
	for i := range syncstorelocations {
		wg.Add(1)
		go func(sl *SyncStoreLocation) {
			db.computeStockStorelocationConsumable(p, sl, mu)
			wg.Done()
		}(&syncstorelocations[i])
	}

	wg.Wait()

	var result []StoreLocation
	for i := range syncstorelocations {
		result = append(result, *syncstorelocations[i].Storelocation)
	}

	return result

}

package models

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql" // register mysql driver
	log "github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/constants"
	"github.com/tbellembois/gochimitheque/helpers"
)

func (db *MySQLDataStore) ComputeStockStorelocation(p Product, s *StoreLocation, u Unit) float64 {

	var (
		c           float64 // current s stock for p
		nullc       sql.NullFloat64
		sdbchildren []StoreLocation
		t           float64 // total s stock for p
		err         error
	)

	sqlr := `SELECT SUM(storage.storage_quantity * unit_multiplier) FROM storage
	JOIN unit on storage.unit = unit.unit_id
	WHERE storage.storelocation = ? AND
	storage.storage_quantity IS NOT NULL AND
	storage.product = ? AND
	(storage.unit = ? OR storage.unit IN (select unit_id FROM unit WHERE unit.unit = ?))`

	// getting current s stock
	if err = db.Get(&nullc, sqlr, s.StoreLocationID.Int64, p.ProductID, u.UnitID.Int64, u.UnitID.Int64); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("ComputeStockStorelocation")
		return 0
	}
	if nullc.Valid {
		c = nullc.Float64
		t = nullc.Float64
	}

	// getting s children
	if sdbchildren, err = db.GetStoreLocationChildren(int(s.StoreLocationID.Int64)); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("ComputeStockStorelocation")
		return 0
	}

	// retrieving or appending children to s and computing their stocks
	for _, sdbchild := range sdbchildren {
		var (
			child      *StoreLocation
			childfound bool
		)
		childfound = false
		for i, schild := range (*s).Children {
			if schild.StoreLocationID == sdbchild.StoreLocationID {
				// child found
				child = (*s).Children[i]
				childfound = true
				break
			}
		}
		if !childfound {
			// child not found
			child = &StoreLocation{
				StoreLocationID:   sdbchild.StoreLocationID,
				StoreLocationName: sdbchild.StoreLocationName,
				Entity:            sdbchild.Entity,
			}
			(*s).Children = append((*s).Children, child)
		}

		t += db.ComputeStockStorelocation(p, child, u)
	}

	(*s).Stocks = append((*s).Stocks, Stock{Total: t, Current: c, Unit: u})

	return c
}

func (db *MySQLDataStore) ComputeStockEntity(p Product, r *http.Request) []StoreLocation {

	var (
		units          []Unit          // reference units
		storelocations []StoreLocation // e root storelocations
		entities       []Entity        // entities
		eids           []int           // entities ids
		err            error
	)

	// getting the entities (GetEntities returns only entities the connected user can see)
	h, _ := helpers.NewdbselectparamEntity(r, nil)
	if entities, _, err = db.GetEntities(h); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("ComputeStockEntity")
		return []StoreLocation{}
	}
	for _, e := range entities {
		eids = append(eids, e.EntityID)
	}

	// getting the reference units
	sqlr := `SELECT unit.unit_id, unit.unit_label FROM unit
	WHERE unit.unit IS NULL`
	if err = db.Select(&units, sqlr); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("ComputeStockEntity")
		return []StoreLocation{}
	}

	// getting the root store locations
	q, args, err := sqlx.In(`SELECT storelocation.storelocation_id, storelocation.storelocation_name, storelocation.storelocation_color
	FROM storelocation
	WHERE storelocation.storelocation IS NULL AND storelocation.entity IN (?)`, eids)
	if err = db.Select(&storelocations, q, args...); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("ComputeStockEntity")
		return []StoreLocation{}
	}

	// computing stocks
	for i := range storelocations {
		for _, u := range units {
			db.ComputeStockStorelocation(p, &storelocations[i], u)
		}
	}

	return storelocations
}

// type stockMapKey struct {
// 	sid int64 // store location id
// 	uid int64 // init id
// }

// type stockMapValue struct {
// 	t float64 // total
// 	c float64 // current
// }

// type StockMap map[stockMapKey]stockMapValue

// func (db *MySQLDataStore) ComputeStockStorelocation(p Product, s StoreLocation, u Unit, m *StockMap) float64 {

// 	var (
// 		c     float64 // current s stock for p
// 		nullc sql.NullFloat64
// 		t     float64 // total s stock for p
// 		err   error
// 		sc    []StoreLocation // s children
// 	)

// 	sqlr := `SELECT SUM(storage.storage_quantity * unit_multiplier) FROM storage
// 	JOIN unit on storage.unit = unit.unit_id
// 	WHERE storage.storelocation = ? AND
// 	storage.storage_quantity IS NOT NULL AND
// 	storage.product = ? AND
// 	(storage.unit = ? OR storage.unit IN (select unit_id FROM unit WHERE unit.unit = ?))`

// 	// getting current s stock
// 	if err = db.Get(&nullc, sqlr, s.StoreLocationID.Int64, p.ProductID, u.UnitID.Int64, u.UnitID.Int64); err != nil {
// 		log.WithFields(log.Fields{"err": err.Error()}).Error("ComputeStockStorelocation")
// 		return 0
// 	}
// 	if nullc.Valid {
// 		c = nullc.Float64
// 	}

// 	// getting s children
// 	if sc, err = db.GetStoreLocationChildren(int(s.StoreLocationID.Int64)); err != nil {
// 		log.WithFields(log.Fields{"err": err.Error()}).Error("ComputeStockStorelocation")
// 		return 0
// 	}

// 	// parsing the children
// 	for _, sci := range sc {
// 		t += db.ComputeStockStorelocation(p, sci, u, m)
// 	}

// 	k := stockMapKey{sid: s.StoreLocationID.Int64, uid: u.UnitID.Int64}
// 	v := stockMapValue{t: t, c: c}
// 	(*m)[k] = v

// 	return c
// }

// func (db *MySQLDataStore) ComputeStockEntity(p Product, e Entity) StockMap {

// 	var (
// 		m              StockMap
// 		units          []Unit          // reference units
// 		storelocations []StoreLocation // e root storelocations
// 		err            error
// 	)

// 	m = make(StockMap)

// 	// getting the reference units
// 	sqlr := `SELECT unit.unit_id FROM unit
// 	WHERE unit.unit = 1`
// 	if err = db.Select(&units, sqlr); err != nil {
// 		log.WithFields(log.Fields{"err": err.Error()}).Error("ComputeStockEntity")
// 		return StockMap{}
// 	}

// 	// getting e root store locations
// 	sqlr = `SELECT storelocation.storelocation_id
// 	FROM storelocation
// 	WHERE storelocation.storelocation IS NULL AND storelocation.entity = ?`
// 	if err = db.Select(&storelocations, sqlr, e.EntityID); err != nil {
// 		log.WithFields(log.Fields{"err": err.Error()}).Error("ComputeStockEntity")
// 		return StockMap{}
// 	}

// 	// computing stocks
// 	for _, sl := range storelocations {
// 		for _, u := range units {
// 			db.ComputeStockStorelocation(p, sl, u, &m)
// 		}
// 	}

// 	return m
// }

// GetEntities returns the entities matching the search criteria
// order, offset and limit are passed to the sql request
func (db *MySQLDataStore) GetEntities(p helpers.DbselectparamEntity) ([]Entity, int, error) {
	var (
		entities                                []Entity
		count                                   int
		req, precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                                  *sqlx.NamedStmt
		snstmt                                  *sqlx.NamedStmt
		err                                     error
	)
	log.WithFields(log.Fields{"p": p}).Debug("GetEntities")

	precreq.WriteString(" SELECT count(DISTINCT e.entity_id)")
	presreq.WriteString(" SELECT e.entity_id, e.entity_name, e.entity_description")
	comreq.WriteString(" FROM entity AS e, person as p")
	// filter by permissions
	comreq.WriteString(` JOIN permission AS perm ON
	(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
	(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
	(perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
	(perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "r" and perm.permission_entity_id = e.entity_id)
	`)
	comreq.WriteString(" WHERE e.entity_name LIKE :search")
	postsreq.WriteString(" GROUP BY e.entity_id")
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
	}

	// select
	if err = snstmt.Select(&entities, m); err != nil {
		return nil, 0, err
	}
	// count
	if err = cnstmt.Get(&count, m); err != nil {
		return nil, 0, err
	}

	//
	// getting managers
	//
	for i, e := range entities {
		// note: do not modify e but entities[i] instead
		req.Reset()
		req.WriteString("SELECT person_id, person_email FROM person")
		req.WriteString(" JOIN entitypeople ON entitypeople.entitypeople_person_id = person.person_id")
		req.WriteString(" JOIN entity ON entitypeople.entitypeople_entity_id = entity.entity_id")
		req.WriteString(" WHERE entity.entity_id = ?")

		if err = db.Select(&entities[i].Managers, req.String(), e.EntityID); err != nil {
			return nil, 0, err
		}
	}

	log.WithFields(log.Fields{"entities": entities, "count": count}).Debug("GetEntities")
	return entities, count, nil
}

// GetEntity returns the entity with id "id"
func (db *MySQLDataStore) GetEntity(id int) (Entity, error) {
	var (
		entity Entity
		sqlr   string
		err    error
	)
	log.WithFields(log.Fields{"id": id}).Debug("GetEntity")

	sqlr = `SELECT e.entity_id, e.entity_name, e.entity_description
	FROM entity AS e
	WHERE e.entity_id = ?`
	if err = db.Get(&entity, sqlr, id); err != nil {
		return Entity{}, err
	}
	log.WithFields(log.Fields{"ID": id, "entity": entity}).Debug("GetEntity")
	return entity, nil
}

// GetEntityPeople returns the entity (with id "id") managers
func (db *MySQLDataStore) GetEntityPeople(id int) ([]Person, error) {
	var (
		people []Person
		sqlr   string
		err    error
	)

	sqlr = `SELECT p.person_id, p.person_email
	FROM person AS p, entitypeople
	WHERE entitypeople.entitypeople_person_id == p.person_id AND entitypeople.entitypeople_entity_id = ?`
	if err = db.Select(&people, sqlr, id); err != nil {
		return []Person{}, err
	}
	log.WithFields(log.Fields{"ID": id, "people": people}).Debug("GetEntityPeople")
	return people, nil
}

// DeleteEntity deletes the entity with id "id"
func (db *MySQLDataStore) DeleteEntity(id int) error {
	var (
		sqlr string
		err  error
	)
	sqlr = `DELETE FROM entity 
	WHERE entity_id = ?`
	if _, err = db.Exec(sqlr, id); err != nil {
		return err
	}
	return nil
}

// CreateEntity creates the given entity
func (db *MySQLDataStore) CreateEntity(e Entity) (error, int) {
	var (
		sqlr   string
		res    sql.Result
		lastid int64
		err    error
	)
	// FIXME: use a transaction here
	sqlr = `INSERT INTO entity(entity_name, entity_description) VALUES (?, ?)`
	if res, err = db.Exec(sqlr, e.EntityName, e.EntityDescription); err != nil {
		return err, 0
	}

	// getting the last inserted id
	if lastid, err = res.LastInsertId(); err != nil {
		return err, 0
	}
	e.EntityID = int(lastid)

	// adding the new managers
	for _, m := range e.Managers {
		sqlr = `INSERT INTO entitypeople (entitypeople_entity_id, entitypeople_person_id) values (?, ?)`
		if _, err = db.Exec(sqlr, e.EntityID, m.PersonID); err != nil {
			return err, 0
		}

		// setting the manager in the entity
		sqlr = `INSERT OR IGNORE INTO personentities(personentities_person_id, personentities_entity_id) 
			VALUES (?, ?)`
		if _, err = db.Exec(sqlr, m.PersonID, e.EntityID); err != nil {
			return err, 0
		}

		// setting the manager permissions in the entity
		// 1. lazily deleting former permissions
		sqlr = `DELETE FROM permission 
			WHERE person = ? and permission_entity_id = ?`
		if _, err = db.Exec(sqlr, m.PersonID, e.EntityID); err != nil {
			return err, 0
		}
		// 2. inserting manager permissions
		sqlr = `INSERT INTO permission(person, permission_perm_name, permission_item_name, permission_entity_id) 
			VALUES (?, ?, ?, ?)`
		if _, err = db.Exec(sqlr, m.PersonID, "all", "all", e.EntityID); err != nil {
			return err, 0
		}
	}

	return nil, e.EntityID
}

// UpdateEntity updates the given entity
func (db *MySQLDataStore) UpdateEntity(e Entity) error {
	var (
		sqlr     string
		sqla     []interface{}
		sbuilder sq.DeleteBuilder
		err      error
	)
	log.WithFields(log.Fields{"e": e}).Debug("UpdateEntity")

	// updating the entity
	// FIXME: use a transaction here
	sqlr = `UPDATE entity SET entity_name = ?, entity_description = ?
	WHERE entity_id = ?`
	if _, err = db.Exec(sqlr, e.EntityName, e.EntityDescription, e.EntityID); err != nil {
		return err
	}

	if len(e.Managers) != 0 {
		// removing former managers
		notin := sq.Or{}
		// ex: AND (entitypeople_person_id <> ? OR entitypeople_person_id <> ?)
		for _, m := range e.Managers {
			notin = append(notin, sq.NotEq{"entitypeople_person_id": m.PersonID})
		}
		// ex: DELETE FROM entitypeople WHERE (entitypeople_entity_id = ? AND (entitypeople_person_id <> ? OR entitypeople_person_id <> ?)
		sbuilder = sq.Delete(`entitypeople`).Where(
			sq.And{
				sq.Eq{`entitypeople_entity_id`: e.EntityID},
				notin})
	} else {
		sbuilder = sq.Delete(`entitypeople`).Where(
			sq.Eq{`entitypeople_entity_id`: e.EntityID})
	}
	sqlr, sqla, err = sbuilder.ToSql()
	if err != nil {
		return err
	}
	if _, err = db.Exec(sqlr, sqla...); err != nil {
		return err
	}

	// TODO: removing former managers permissions

	// adding the new ones
	for _, m := range e.Managers {
		// adding the manager
		sqlr = `INSERT OR IGNORE INTO entitypeople (entitypeople_entity_id, entitypeople_person_id) VALUES (?, ?)`
		if _, err = db.Exec(sqlr, e.EntityID, m.PersonID); err != nil {
			return err
		}

		for _, man := range e.Managers {
			// setting the manager in the entity
			sqlr = `INSERT OR IGNORE INTO personentities(personentities_person_id, personentities_entity_id) 
			VALUES (?, ?)`
			if _, err = db.Exec(sqlr, man.PersonID, e.EntityID); err != nil {
				return err
			}

			// setting the manager permissions in the entity
			// 1. lazily deleting former permissions
			sqlr = `DELETE FROM permission 
			WHERE person = ? and permission_entity_id = ?`
			if _, err = db.Exec(sqlr, man.PersonID, e.EntityID); err != nil {
				return err
			}
			// 2. inserting manager permissions
			sqlr = `INSERT INTO permission(person, permission_perm_name, permission_item_name, permission_entity_id) 
			VALUES (?, ?, ?, ?)`
			if _, err = db.Exec(sqlr, man.PersonID, "all", "all", e.EntityID); err != nil {
				return err
			}

		}
	}

	return nil
}

// IsEntityEmpty returns true is the entity is empty
func (db *MySQLDataStore) IsEntityEmpty(id int) (bool, error) {
	var (
		res   bool
		count int
		sqlr  string
		err   error
	)

	sqlr = "SELECT count(*) from personentities WHERE personentities.personentities_entity_id = ?"
	if err = db.Get(&count, sqlr, id); err != nil {
		return false, err
	}
	log.WithFields(log.Fields{"id": id, "count": count}).Debug("IsEntityEmpty")
	if count == 0 {
		res = true
	} else {
		res = false
	}
	return res, nil
}

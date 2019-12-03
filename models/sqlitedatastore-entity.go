package models

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver

	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/constants"
	"github.com/tbellembois/gochimitheque/global"
	"github.com/tbellembois/gochimitheque/helpers"
)

// ComputeStockStorelocation returns the quantity of product p in the store location s for the unit u
func (db *SQLiteDataStore) ComputeStockStorelocation(p Product, s *StoreLocation, u Unit) float64 {

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
		global.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("ComputeStockStorelocation")
		return 0
	}
	if nullc.Valid {
		c = nullc.Float64
		t = nullc.Float64
	}
	global.Log.WithFields(logrus.Fields{"p": p, "s": s, "u": u, "c": c}).Debug("ComputeStockStorelocation")

	// getting s children
	if sdbchildren, err = db.GetStoreLocationChildren(int(s.StoreLocationID.Int64)); err != nil {
		global.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("ComputeStockStorelocation")
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

// ComputeStockStorelocationNoUnit returns the quantity of product p with no unit in the store location s
func (db *SQLiteDataStore) ComputeStockStorelocationNoUnit(p Product, s *StoreLocation) float64 {

	var (
		c           float64 // current s stock for p
		nullc       sql.NullFloat64
		sdbchildren []StoreLocation
		t           float64 // total s stock for p
		err         error
	)

	sqlr := `SELECT count(*) FROM storage
	LEFT JOIN unit on storage.unit = unit.unit_id
	WHERE storage.storelocation = ? AND
	storage.storage_quantity IS NOT NULL AND
	storage.product = ? AND
	storage.unit IS NULL`

	// getting current s stock
	if err = db.Get(&nullc, sqlr, s.StoreLocationID.Int64, p.ProductID); err != nil {
		global.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("ComputeStockStorelocationNoUnit")
		return 0
	}
	if nullc.Valid {
		c = nullc.Float64
		t = nullc.Float64
	}
	global.Log.WithFields(logrus.Fields{"p": p, "s": s, "c": c}).Debug("ComputeStockStorelocationNoUnit")

	// getting s children
	if sdbchildren, err = db.GetStoreLocationChildren(int(s.StoreLocationID.Int64)); err != nil {
		global.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("ComputeStockStorelocationNoUnit")
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

		t += db.ComputeStockStorelocationNoUnit(p, child)
	}

	(*s).Stocks = append((*s).Stocks, Stock{Total: t, Current: c, Unit: Unit{}})

	return c
}

// ComputeStockEntity returns the root store locations of the entity(ies) of the loggued user.
// Each store location has a Stocks []Stock field containing the stocks of the product p for each unit
func (db *SQLiteDataStore) ComputeStockEntity(p Product, r *http.Request) []StoreLocation {

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
		global.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("ComputeStockEntity")
		return []StoreLocation{}
	}
	for _, e := range entities {
		eids = append(eids, e.EntityID)
	}

	// getting the reference units
	sqlr := `SELECT unit.unit_id, unit.unit_label FROM unit
	WHERE unit.unit IS NULL`
	if err = db.Select(&units, sqlr); err != nil {
		global.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("ComputeStockEntity")
		return []StoreLocation{}
	}

	// getting the root store locations
	q, args, err := sqlx.In(`SELECT storelocation.storelocation_id, storelocation.storelocation_name, storelocation.storelocation_color
	FROM storelocation
	WHERE storelocation.storelocation IS NULL AND storelocation.entity IN (?)`, eids)
	if err = db.Select(&storelocations, q, args...); err != nil {
		global.Log.WithFields(logrus.Fields{"err": err.Error()}).Error("ComputeStockEntity")
		return []StoreLocation{}
	}

	// computing stocks for storages with units
	for i := range storelocations {
		for _, u := range units {
			db.ComputeStockStorelocation(p, &storelocations[i], u)
		}
	}
	// computing stocks for storages with units
	for i := range storelocations {
		db.ComputeStockStorelocationNoUnit(p, &storelocations[i])
	}

	return storelocations
}

// GetEntities returns the entities matching the search criteria
// order, offset and limit are passed to the sql request
func (db *SQLiteDataStore) GetEntities(p helpers.DbselectparamEntity) ([]Entity, int, error) {
	var (
		entities                                []Entity
		count                                   int
		req, precreq, presreq, comreq, postsreq strings.Builder
		cnstmt                                  *sqlx.NamedStmt
		snstmt                                  *sqlx.NamedStmt
		err                                     error
	)
	global.Log.WithFields(logrus.Fields{"p": p}).Debug("GetEntities")

	precreq.WriteString(" SELECT count(DISTINCT e.entity_id)")
	presreq.WriteString(" SELECT e.entity_id, e.entity_name, e.entity_description")
	comreq.WriteString(" FROM entity AS e, person as p")
	// filter by permissions
	// comreq.WriteString(` JOIN permission AS perm ON
	// (perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
	// (perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
	// (perm.person = :personid and perm.permission_item_name = "all" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
	// (perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "all" and perm.permission_entity_id = e.entity_id) OR
	// (perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "all" and perm.permission_entity_id = -1) OR
	// (perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "r" and perm.permission_entity_id = -1) OR
	// (perm.person = :personid and perm.permission_item_name = "entities" and perm.permission_perm_name = "r" and perm.permission_entity_id = e.entity_id)
	// `)
	comreq.WriteString(` JOIN permission AS perm ON
	perm.person = :personid and (perm.permission_item_name in ("all", "entities")) and (perm.permission_perm_name in ("all", "r")) and (perm.permission_entity_id in (-1, e.entity_id))
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

	//
	// getting number of store locations for each entity
	//
	for i, ent := range entities {
		// getting the total store location count
		req.Reset()
		req.WriteString("SELECT count(storelocation_id) from storelocation")
		req.WriteString(" WHERE entity = ?")
		if err = db.Get(&entities[i].EntitySLC, req.String(), ent.EntityID); err != nil {
			return nil, 0, err
		}
	}

	//
	// getting number of person for each entity
	//
	for i, ent := range entities {
		// getting the total person count
		req.Reset()
		req.WriteString("SELECT count(personentities_person_id) from personentities")
		req.WriteString(" WHERE personentities_entity_id = ?")
		if err = db.Get(&entities[i].EntityPC, req.String(), ent.EntityID); err != nil {
			return nil, 0, err
		}
	}

	global.Log.WithFields(logrus.Fields{"entities": entities, "count": count}).Debug("GetEntities")
	return entities, count, nil
}

// GetEntity returns the entity with id "id"
func (db *SQLiteDataStore) GetEntity(id int) (Entity, error) {
	var (
		entity Entity
		sqlr   string
		err    error
	)
	global.Log.WithFields(logrus.Fields{"id": id}).Debug("GetEntity")

	sqlr = `SELECT e.entity_id, e.entity_name, e.entity_description
	FROM entity AS e
	WHERE e.entity_id = ?`
	if err = db.Get(&entity, sqlr, id); err != nil {
		return Entity{}, err
	}
	global.Log.WithFields(logrus.Fields{"ID": id, "entity": entity}).Debug("GetEntity")
	return entity, nil
}

// GetEntityPeople returns the entity (with id "id") managers
func (db *SQLiteDataStore) GetEntityPeople(id int) ([]Person, error) {
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
	global.Log.WithFields(logrus.Fields{"ID": id, "people": people}).Debug("GetEntityPeople")
	return people, nil
}

// DeleteEntity deletes the entity with id "id"
func (db *SQLiteDataStore) DeleteEntity(id int) error {
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
func (db *SQLiteDataStore) CreateEntity(e Entity) (int, error) {
	var (
		sqlr   string
		tx     *sql.Tx
		res    sql.Result
		lastid int64
		err    error
	)

	// beginning transaction
	if tx, err = db.Begin(); err != nil {
		return 0, err
	}

	sqlr = `INSERT INTO entity(entity_name, entity_description) VALUES (?, ?)`
	if res, err = tx.Exec(sqlr, e.EntityName, e.EntityDescription); err != nil {
		tx.Rollback()
		return 0, err
	}

	// getting the last inserted id
	if lastid, err = res.LastInsertId(); err != nil {
		tx.Rollback()
		return 0, err
	}
	e.EntityID = int(lastid)

	// adding the new managers
	for _, m := range e.Managers {
		sqlr = `INSERT INTO entitypeople (entitypeople_entity_id, entitypeople_person_id) values (?, ?)`
		if _, err = tx.Exec(sqlr, e.EntityID, m.PersonID); err != nil {
			tx.Rollback()
			return 0, err
		}

		// setting the manager in the entity
		sqlr = `INSERT OR IGNORE INTO personentities(personentities_person_id, personentities_entity_id) 
			VALUES (?, ?)`
		if _, err = tx.Exec(sqlr, m.PersonID, e.EntityID); err != nil {
			tx.Rollback()
			return 0, err
		}

		// setting the manager permissions in the entity
		// 1. lazily deleting former permissions
		sqlr = `DELETE FROM permission 
			WHERE person = ? and permission_entity_id = ?`
		if _, err = tx.Exec(sqlr, m.PersonID, e.EntityID); err != nil {
			tx.Rollback()
			return 0, err
		}
		// 2. inserting manager permissions
		sqlr = `INSERT INTO permission(person, permission_perm_name, permission_item_name, permission_entity_id) 
			VALUES (?, ?, ?, ?)`
		if _, err = tx.Exec(sqlr, m.PersonID, "all", "all", e.EntityID); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// committing changes
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return 0, err
	}

	return e.EntityID, nil
}

// UpdateEntity updates the given entity
func (db *SQLiteDataStore) UpdateEntity(e Entity) error {
	var (
		sqlr     string
		sqla     []interface{}
		sbuilder sq.DeleteBuilder
		tx       *sql.Tx
		err      error
	)
	global.Log.WithFields(logrus.Fields{"e": e}).Debug("UpdateEntity")

	// beginning transaction
	if tx, err = db.Begin(); err != nil {
		return err
	}

	// updating the entity
	sqlr = `UPDATE entity SET entity_name = ?, entity_description = ?
	WHERE entity_id = ?`
	if _, err = tx.Exec(sqlr, e.EntityName, e.EntityDescription, e.EntityID); err != nil {
		tx.Rollback()
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
	if _, err = tx.Exec(sqlr, sqla...); err != nil {
		return err
	}

	// adding the new ones
	for _, m := range e.Managers {
		// adding the manager
		sqlr = `INSERT OR IGNORE INTO entitypeople (entitypeople_entity_id, entitypeople_person_id) VALUES (?, ?)`
		if _, err = tx.Exec(sqlr, e.EntityID, m.PersonID); err != nil {
			tx.Rollback()
			return err
		}

		for _, man := range e.Managers {
			// setting the manager in the entity
			sqlr = `INSERT OR IGNORE INTO personentities(personentities_person_id, personentities_entity_id) 
			VALUES (?, ?)`
			if _, err = tx.Exec(sqlr, man.PersonID, e.EntityID); err != nil {
				tx.Rollback()
				return err
			}

			// setting the manager permissions in the entity
			// 1. lazily deleting former permissions
			sqlr = `DELETE FROM permission 
			WHERE person = ? and permission_entity_id = ?`
			if _, err = tx.Exec(sqlr, man.PersonID, e.EntityID); err != nil {
				tx.Rollback()
				return err
			}
			// 2. inserting manager permissions
			sqlr = `INSERT INTO permission(person, permission_perm_name, permission_item_name, permission_entity_id) 
			VALUES (?, ?, ?, ?)`
			if _, err = tx.Exec(sqlr, man.PersonID, "all", "all", e.EntityID); err != nil {
				tx.Rollback()
				return err
			}

		}
	}

	// committing changes
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// IsEntityEmpty returns true is the entity is empty
func (db *SQLiteDataStore) IsEntityEmpty(id int) (bool, error) {
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
	global.Log.WithFields(logrus.Fields{"id": id, "count": count}).Debug("IsEntityEmpty")
	if count == 0 {
		res = true
	} else {
		res = false
	}
	return res, nil
}

// HasEntityNoStorelocation returns true is the entity has no store location
func (db *SQLiteDataStore) HasEntityNoStorelocation(id int) (bool, error) {
	var (
		res   bool
		count int
		sqlr  string
		err   error
	)

	sqlr = "SELECT count(*) from storelocation WHERE storelocation.entity = ?"
	if err = db.Get(&count, sqlr, id); err != nil {
		return false, err
	}
	global.Log.WithFields(logrus.Fields{"id": id, "count": count}).Debug("HasEntityNoStorelocation")
	if count == 0 {
		res = true
	} else {
		res = false
	}
	return res, nil
}

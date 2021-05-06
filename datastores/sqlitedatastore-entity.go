package datastores

import (
	"database/sql"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	. "github.com/tbellembois/gochimitheque/models"
)

// GetEntities select the entities matching p
// and visible by the connected user.
func (db *SQLiteDataStore) GetEntities(p DbselectparamEntity) ([]Entity, int, error) {

	var (
		err                   error
		entities              []Entity
		count                 int
		countSql, selectSql   string
		countArgs, selectArgs []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("GetEntities")

	dialect := goqu.Dialect("sqlite3")
	entityTable := goqu.T("entity")
	personTable := goqu.T("person")
	storelocationTable := goqu.T("storelocation")
	personentitiesTable := goqu.T("personentities")

	// Prepare orderby/order clause.
	orderByClause := p.GetOrderBy()
	orderClause := goqu.I(orderByClause).Asc()
	if strings.ToLower(p.GetOrder()) == "desc" {
		orderClause = goqu.I(orderByClause).Desc()
	}

	// Join, where.
	joinClause := dialect.From(
		entityTable.As("e"),
		personTable.As("p"),
	).Join(
		goqu.T("permission").As("perm"),
		goqu.On(
			goqu.Ex{
				"perm.person":               p.GetLoggedPersonID(),
				"perm.permission_item_name": []string{"all", "entities"},
				"perm.permission_perm_name": []string{"all", "r", "w"},
				"perm.permission_entity_id": []interface{}{-1, goqu.I("e.entity_id")},
			},
		),
	).Where(
		goqu.I("e.entity_name").Like(p.GetSearch()),
	)

	if countSql, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("e.entity_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}
	if selectSql, selectArgs, err = joinClause.Select(
		goqu.I("e.entity_id"),
		goqu.I("e.entity_name"),
		goqu.I("e.entity_description"),
	).GroupBy(goqu.I("e.entity_id")).Order(orderClause).Limit(uint(p.GetLimit())).Offset(uint(p.GetOffset())).ToSQL(); err != nil {
		return nil, 0, err
	}

	// select
	if err = db.Select(&entities, selectSql, selectArgs...); err != nil {
		return nil, 0, err
	}
	// count
	if err = db.Get(&count, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	//
	// Getting the entity managers.
	//
	for i, entity := range entities {

		sQuery := dialect.From(personTable).Join(
			goqu.T("entitypeople"),
			goqu.On(
				goqu.Ex{
					"entitypeople.entitypeople_person_id": goqu.I("person.person_id"),
				},
			),
		).Join(
			goqu.T("entity"),
			goqu.On(
				goqu.Ex{
					"entitypeople.entitypeople_entity_id": goqu.I("entity.entity_id"),
				},
			),
		).Where(
			goqu.I("entity.entity_id").Eq(entity.EntityID),
		).Select(
			"person_id",
			"person_email",
		)

		var (
			sqlr string
			args []interface{}
		)
		if sqlr, args, err = sQuery.ToSQL(); err != nil {
			logger.Log.Error(err)
			return nil, 0, err
		}

		if err = db.Select(&entities[i].Managers, sqlr, args...); err != nil {
			return nil, 0, err
		}

	}

	//
	// Getting entities number of store locations.
	//
	for i, entity := range entities {

		sQuery := dialect.From(storelocationTable).Where(
			goqu.I("entity").Eq(entity.EntityID),
		).Select(
			goqu.COUNT(goqu.I("storelocation_id")),
		)

		var (
			sqlr string
			args []interface{}
		)
		if sqlr, args, err = sQuery.ToSQL(); err != nil {
			logger.Log.Error(err)
			return nil, 0, err
		}

		if err = db.Get(&entities[i].EntitySLC, sqlr, args...); err != nil {
			return nil, 0, err
		}

	}

	//
	// Getting entities number of members.
	//
	for i, entity := range entities {

		sQuery := dialect.From(personentitiesTable).Where(
			goqu.I("personentities_entity_id").Eq(entity.EntityID),
		).Select(
			goqu.COUNT(goqu.I("personentities_person_id")),
		)

		var (
			sqlr string
			args []interface{}
		)
		if sqlr, args, err = sQuery.ToSQL(); err != nil {
			logger.Log.Error(err)
			return nil, 0, err
		}

		if err = db.Get(&entities[i].EntityPC, sqlr, args...); err != nil {
			return nil, 0, err
		}

	}

	logger.Log.WithFields(logrus.Fields{"entities": entities, "count": count}).Debug("GetEntities")
	return entities, count, nil

}

// GetEntity select the entity by id.
func (db *SQLiteDataStore) GetEntity(id int) (Entity, error) {

	var (
		err    error
		sqlr   string
		args   []interface{}
		entity Entity
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetEntity")

	dialect := goqu.Dialect("sqlite3")
	tableEntity := goqu.T("entity")
	tablePerson := goqu.T("person")

	sQuery := dialect.From(tableEntity.As("e")).Where(
		goqu.I("e.entity_id").Eq(id),
	).Select(
		goqu.I("e.entity_id"),
		goqu.I("e.entity_name"),
		goqu.I("e.entity_description"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return Entity{}, err
	}

	if err = db.Get(&entity, sqlr, args...); err != nil {
		return Entity{}, err
	}

	// Managers.
	sQuery = dialect.From(tablePerson).Join(
		goqu.T("entitypeople"),
		goqu.On(goqu.Ex{"entitypeople.entitypeople_person_id": goqu.I("person.person_id")}),
	).Join(
		goqu.T("entity"),
		goqu.On(goqu.Ex{"entitypeople.entitypeople_entity_id": goqu.I("entity.entity_id")}),
	).Where(
		goqu.I("entity.entity_id").Eq(id),
	).Select(
		goqu.I("person_id"),
		goqu.I("person_email"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return Entity{}, err
	}

	if err = db.Select(&entity.Managers, sqlr, args...); err != nil {
		return Entity{}, err
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "entity": entity}).Debug("GetEntity")

	return entity, nil

}

// GetEntityManager select the entity managers.
func (db *SQLiteDataStore) GetEntityManager(id int) ([]Person, error) {

	var (
		err    error
		sqlr   string
		args   []interface{}
		people []Person
	)

	dialect := goqu.Dialect("sqlite3")
	tablePerson := goqu.T("person")
	tableEntitypeople := goqu.T("entitypeople")

	sQuery := dialect.From(tablePerson.As("p"), tableEntitypeople).Where(
		goqu.Ex{
			"entitypeople.entitypeople_person_id": goqu.I("p.person_id"),
			"entitypeople.entitypeople_entity_id": id,
		},
	).Select(
		goqu.I("p.person_id"),
		goqu.I("p.person_email"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return []Person{}, err
	}

	if err = db.Select(&people, sqlr, args...); err != nil {
		return []Person{}, err
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "people": people}).Debug("GetEntityPeople")
	return people, nil

}

// DeleteEntity delete the entity by id.
func (db *SQLiteDataStore) DeleteEntity(id int) error {

	var (
		err  error
		sqlr string
		args []interface{}
	)

	dialect := goqu.Dialect("sqlite3")
	tableEntity := goqu.T("entity")

	sQuery := dialect.From(tableEntity).Where(
		goqu.I("entity_id").Eq(id),
	).Delete()

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return err
	}

	if _, err = db.Exec(sqlr, args...); err != nil {
		return err
	}

	return nil

}

// CreateEntity insert e.
func (db *SQLiteDataStore) CreateEntity(e Entity) (lastInsertId int64, err error) {

	var (
		sqlr string
		args []interface{}
		tx   *sql.Tx
		res  sql.Result
	)

	logger.Log.WithFields(logrus.Fields{"e": e}).Debug("CreateEntity")

	dialect := goqu.Dialect("sqlite3")
	tableEntity := goqu.T("entity")

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

	iQuery := dialect.Insert(tableEntity).Rows(
		goqu.Record{
			"entity_name":        e.EntityName,
			"entity_description": e.EntityDescription,
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
	e.EntityID = int(lastInsertId)

	// Setting up the managers.
	for _, m := range e.Managers {

		logger.Log.WithFields(logrus.Fields{"m": m}).Debug("CreateEntity")

		// Adding the managers.
		if sqlr, args, err = dialect.Insert(goqu.T("entitypeople")).Rows(
			goqu.Record{
				"entitypeople_entity_id": e.EntityID,
				"entitypeople_person_id": m.PersonID,
			},
		).ToSQL(); err != nil {
			return
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			return
		}

		// Adding the managers as members of the entity.
		if sqlr, args, err = dialect.Insert(goqu.T("personentities")).Rows(
			goqu.Record{
				"personentities_person_id": m.PersonID,
				"personentities_entity_id": e.EntityID,
			},
		).OnConflict(goqu.DoNothing()).ToSQL(); err != nil {
			return
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			return
		}

		// Setting the manager permissions.
		// 1. lazily deleting former permissions
		if sqlr, args, err = dialect.From(goqu.T("permission")).Where(
			goqu.Ex{
				"person":               m.PersonID,
				"permission_entity_id": e.EntityID,
			},
		).Delete().ToSQL(); err != nil {
			return
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			return
		}

		// 2. inserting new permissions
		if sqlr, args, err = dialect.From(goqu.T("permission")).Prepared(true).Insert().Rows(
			goqu.Record{
				"person":               m.PersonID,
				"permission_perm_name": "all",
				"permission_item_name": "all",
				"permission_entity_id": e.EntityID,
			},
		).ToSQL(); err != nil {
			return
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			return
		}

		args = []interface{}{m.PersonID, "w", "products", e.EntityID}
		if _, err = tx.Exec(sqlr, args...); err != nil {
			return
		}

		args = []interface{}{m.PersonID, "w", "rproducts", e.EntityID}
		if _, err = tx.Exec(sqlr, args...); err != nil {
			return
		}

	}

	return

}

// UpdateEntity update e.
func (db *SQLiteDataStore) UpdateEntity(e Entity) (err error) {

	var (
		tx *sql.Tx
	)

	dialect := goqu.Dialect("sqlite3")
	tableEntity := goqu.T("entity")

	if tx, err = db.Begin(); err != nil {
		return err
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

	var (
		sqlr string
		args []interface{}
	)

	if sqlr, args, err = dialect.Update(tableEntity).Set(
		goqu.Record{
			"entity_name":        e.EntityName,
			"entity_description": e.EntityDescription,
		},
	).Where(
		goqu.I("entity_id").Eq(e.EntityID),
	).ToSQL(); err != nil {
		logger.Log.Error(err)
		return
	}

	if _, err = tx.Exec(sqlr, args...); err != nil {
		return
	}

	// Removing former managers.
	whereAnd := []goqu.Expression{
		goqu.I("entitypeople_entity_id").Eq(e.EntityID),
	}

	if len(e.Managers) != 0 {

		// Except those not removed.
		var notIn []int
		for _, manager := range e.Managers {
			notIn = append(notIn, manager.PersonID)
		}

		whereAnd = append(whereAnd, goqu.I("entitypeople_person_id").NotIn(notIn))

	}

	dQuery := dialect.From(goqu.I("entitypeople")).Where(
		whereAnd...,
	).Delete()

	if sqlr, args, err = dQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return err
	}

	if _, err = tx.Exec(sqlr, args...); err != nil {
		return err
	}

	// Adding the new managers.
	for _, manager := range e.Managers {

		// Adding the manager.
		if sqlr, args, err = dialect.Insert(goqu.T("entitypeople")).Rows(
			goqu.Record{
				"entitypeople_entity_id": e.EntityID,
				"entitypeople_person_id": manager.PersonID,
			},
		).OnConflict(goqu.DoNothing()).ToSQL(); err != nil {
			return
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			return
		}

		// Putting the manager in his entity.
		if sqlr, args, err = dialect.Insert(goqu.T("personentities")).Rows(
			goqu.Record{
				"personentities_person_id": manager.PersonID,
				"personentities_entity_id": e.EntityID,
			},
		).OnConflict(goqu.DoNothing()).ToSQL(); err != nil {
			return
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			return
		}

		// Setting the manager permissions.
		// 1. lazily deleting former permissions
		if sqlr, args, err = dialect.From(goqu.T("permission")).Where(
			goqu.Ex{
				"person":               manager.PersonID,
				"permission_entity_id": e.EntityID,
			},
		).ToSQL(); err != nil {
			return
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			return
		}

		// 2. inserting manager permissions
		// added OR IGNORE bacause w/(r)products/-1 can already exists for man.PersonID
		if sqlr, args, err = dialect.From(goqu.T("permission")).Prepared(true).Insert().Rows(
			goqu.Record{
				"person":               manager.PersonID,
				"permission_perm_name": "all",
				"permission_item_name": "all",
				"permission_entity_id": e.EntityID,
			},
		).OnConflict(goqu.DoNothing()).ToSQL(); err != nil {
			return
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			return
		}

		args = []interface{}{manager.PersonID, "w", "products", "-1"}
		if _, err = tx.Exec(sqlr, args...); err != nil {
			return
		}

		args = []interface{}{manager.PersonID, "w", "rproducts", "-1"}
		if _, err = tx.Exec(sqlr, args...); err != nil {
			return
		}

	}

	return

}

// HasEntityMember returns true is the entity has members.
func (db *SQLiteDataStore) HasEntityMember(id int) (bool, error) {

	var (
		err   error
		sqlr  string
		args  []interface{}
		count int
	)

	dialect := goqu.Dialect("sqlite3")
	tablePersonentities := goqu.T("personentities")

	sQuery := dialect.From(tablePersonentities).Select(
		goqu.COUNT("*"),
	).Where(
		goqu.I("personentities.personentities_entity_id").Eq(id),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return false, err
	}

	if err = db.Get(&count, sqlr, args...); err != nil {
		return false, err
	}

	return count != 0, nil

}

// HasEntityStorelocation returns true is the entity has no store locations.
func (db *SQLiteDataStore) HasEntityStorelocation(id int) (bool, error) {

	var (
		err   error
		sqlr  string
		args  []interface{}
		count int
	)

	dialect := goqu.Dialect("sqlite3")
	tableStorelocation := goqu.T("storelocation")

	sQuery := dialect.From(tableStorelocation).Select(
		goqu.COUNT("*"),
	).Where(
		goqu.I("storelocation.entity").Eq(id),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return false, err
	}

	if err = db.Get(&count, sqlr, args...); err != nil {
		return false, err
	}

	return count != 0, nil

}

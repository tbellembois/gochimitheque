package datastores

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
)

func (db *SQLiteDataStore) DeleteEntity(id int) error {
	var (
		err  error
		sqlr string
		args []interface{}
	)

	dialect := goqu.Dialect("sqlite3")
	tableEntity := goqu.T("entity")
	tableEntityPeople := goqu.T("entitypeople")

	// Managers.
	sQuery := dialect.From(tableEntityPeople).Where(
		goqu.I("entitypeople_entity_id").Eq(id),
	).Delete()

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return err
	}

	if _, err = db.Exec(sqlr, args...); err != nil {
		return err
	}

	// Entity.
	sQuery = dialect.From(tableEntity).Where(
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

func (db *SQLiteDataStore) CreateEntity(e models.Entity) (lastInsertID int64, err error) {
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

	if lastInsertID, err = res.LastInsertId(); err != nil {
		return
	}

	e.EntityID = int(lastInsertID)

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
				"permission_entity": e.EntityID,
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
				"permission_name": "all",
				"permission_item": "all",
				"permission_entity": e.EntityID,
			},
		).ToSQL(); err != nil {
			return
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			return
		}

		if sqlr, args, err = dialect.From(goqu.T("permission")).Insert().Rows(
			goqu.Record{
				"person":               m.PersonID,
				"permission_name": "w",
				"permission_item": "products",
				"permission_entity": e.EntityID,
			},
		).OnConflict(goqu.DoNothing()).ToSQL(); err != nil {
			logger.Log.Errorf("error preparing inserting manager new permissions w products -1: %v", err)
			return
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			logger.Log.Errorf("error inserting manager new permissions w products -1: %v", err)
			return
		}

		if sqlr, args, err = dialect.From(goqu.T("permission")).Insert().Rows(
			goqu.Record{
				"person":               m.PersonID,
				"permission_name": "w",
				"permission_item": "rproducts",
				"permission_entity": e.EntityID,
			},
		).OnConflict(goqu.DoNothing()).ToSQL(); err != nil {
			logger.Log.Errorf("error preparing inserting manager new permissions w rproducts -1: %v", err)
			return
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			logger.Log.Errorf("error inserting manager new permissions w products -1: %v", err)
			return
		}
	}

	return
}

func (db *SQLiteDataStore) UpdateEntity(e models.Entity) (err error) {
	var tx *sql.Tx

	dialect := goqu.Dialect("sqlite3")
	tableEntity := goqu.T("entity")
	tableEntityPeople := goqu.T("entitypeople")

	if tx, err = db.Begin(); err != nil {
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
		logger.Log.Errorf("error preparing updating entity: %v", err)
		return
	}

	if _, err = tx.Exec(sqlr, args...); err != nil {
		logger.Log.Errorf("error updating entity: %v", err)
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

	dQuery := dialect.From(tableEntityPeople).Where(
		whereAnd...,
	).Delete()

	if sqlr, args, err = dQuery.ToSQL(); err != nil {
		logger.Log.Errorf("error preparing removing former managers: %v", err)
		return
	}

	if _, err = tx.Exec(sqlr, args...); err != nil {
		logger.Log.Errorf("error removing former managers: %v", err)
		return
	}

	// Adding the new managers.
	for _, manager := range e.Managers {
		// Adding the manager.
		if sqlr, args, err = dialect.Insert(tableEntityPeople).Rows(
			goqu.Record{
				"entitypeople_entity_id": e.EntityID,
				"entitypeople_person_id": manager.PersonID,
			},
		).OnConflict(goqu.DoNothing()).ToSQL(); err != nil {
			logger.Log.Errorf("error preparing inserting new managers: %v", err)
			return
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			logger.Log.Errorf("error inserting new managers: %v", err)
			return
		}

		// Putting the manager in his entity.
		if sqlr, args, err = dialect.Insert(goqu.T("personentities")).Rows(
			goqu.Record{
				"personentities_person_id": manager.PersonID,
				"personentities_entity_id": e.EntityID,
			},
		).OnConflict(goqu.DoNothing()).ToSQL(); err != nil {
			logger.Log.Errorf("error preparing putting manager in its entity: %v", err)
			return
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			logger.Log.Errorf("error putting manager in its entity: %v", err)
			return
		}

		// Setting the manager permissions.
		// 1. lazily deleting former permissions
		if sqlr, args, err = dialect.From(goqu.T("permission")).Where(
			goqu.Ex{
				"person":               manager.PersonID,
				"permission_entity": e.EntityID,
			},
		).ToSQL(); err != nil {
			logger.Log.Errorf("error preparing deleting manager permissions: %v", err)
			return
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			logger.Log.Errorf("error deleting manager permissions: %v", err)
			return
		}

		// 2. inserting manager permissions
		// added OR IGNORE bacause w/(r)products/-1 can already exists for man.PersonID
		if sqlr, args, err = dialect.From(goqu.T("permission")).Insert().Rows(
			goqu.Record{
				"person":               manager.PersonID,
				"permission_name": "all",
				"permission_item": "all",
				"permission_entity": e.EntityID,
			},
		).OnConflict(goqu.DoNothing()).ToSQL(); err != nil {
			logger.Log.Errorf("error preparing inserting manager new permissions: %v", err)
			return
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			logger.Log.Errorf("error inserting manager new permissions: %v", err)
			return
		}

		if sqlr, args, err = dialect.From(goqu.T("permission")).Insert().Rows(
			goqu.Record{
				"person":               manager.PersonID,
				"permission_name": "w",
				"permission_item": "products",
				"permission_entity": "-1",
			},
		).OnConflict(goqu.DoNothing()).ToSQL(); err != nil {
			logger.Log.Errorf("error preparing inserting manager new permissions w products -1: %v", err)
			return
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			logger.Log.Errorf("error inserting manager new permissions w products -1: %v", err)
			return
		}

		if sqlr, args, err = dialect.From(goqu.T("permission")).Insert().Rows(
			goqu.Record{
				"person":               manager.PersonID,
				"permission_name": "w",
				"permission_item": "rproducts",
				"permission_entity": "-1",
			},
		).OnConflict(goqu.DoNothing()).ToSQL(); err != nil {
			logger.Log.Errorf("error preparing inserting manager new permissions w rproducts -1: %v", err)
			return
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			logger.Log.Errorf("error inserting manager new permissions w products -1: %v", err)
			return
		}
	}

	return
}

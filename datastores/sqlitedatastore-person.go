package datastores

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	"github.com/jmoiron/sqlx"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/zmqclient"
)

// DeletePerson deletes the person with id "id".
func (db *SQLiteDataStore) DeletePerson(id int) (err error) {
	var (
		sqlr string
		args []interface{}
		tx   *sqlx.Tx
	)

	dialect := goqu.Dialect("sqlite3")
	tablePerson := goqu.T("person")
	tableStorage := goqu.T("storage")
	tableProduct := goqu.T("product")

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

	// Getting the admin.
	// TODO: remove 1 by connected user id.
	var (
		jsonRawMessage json.RawMessage
		admin          *models.Person
	)

	if jsonRawMessage, err = zmqclient.DBGetPeople("http://localhost/?search=admin@chimitheque.fr", 1); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "zmqclient.DBGetPeople",
			Code:          http.StatusInternalServerError,
		}
	}

	if admin, err = zmqclient.ConvertDBJSONToPerson(jsonRawMessage); err != nil {
		return &models.AppError{
			OriginalError: err,
			Message:       "ConvertDBJSONToPerson",
			Code:          http.StatusInternalServerError,
		}
	}

	// if admin, err = db.GetPersonByEmail("admin@chimitheque.fr"); err != nil {
	// 	return err
	// }

	// Updating storage ownership to admin.
	if sqlr, args, err = dialect.Update(tableStorage).Set(
		goqu.Record{
			"person": admin.PersonID,
		},
	).Where(
		goqu.I("person").Eq(id),
	).ToSQL(); err != nil {
		logger.Log.Errorf("prepare update storage ownership: %s", err)
		return
	}

	if _, err = tx.Exec(sqlr, args...); err != nil {
		logger.Log.Errorf("update storage ownership: %s", err)
		return
	}

	// Updating product ownership to admin.
	if sqlr, args, err = dialect.Update(tableProduct).Set(
		goqu.Record{
			"person": admin.PersonID,
		},
	).Where(
		goqu.I("person").Eq(id),
	).ToSQL(); err != nil {
		logger.Log.Errorf("prepare update product ownership: %s", err)
		return
	}

	if _, err = tx.Exec(sqlr, args...); err != nil {
		logger.Log.Errorf("update product ownership: %s", err)
		return
	}

	// Deleting entity membership.
	if sqlr, args, err = dialect.From(goqu.T("personentities")).Where(
		goqu.I("personentities_person_id").Eq(id),
	).Delete().ToSQL(); err != nil {
		logger.Log.Errorf("prepare delete entity membership: %s", err)
		return
	}

	if _, err = tx.Exec(sqlr, args...); err != nil {
		logger.Log.Errorf("delete entity membership: %s", err)
		return
	}

	// Remove manager.
	// Should not be used as we can not delete a manager.
	if sqlr, args, err = dialect.From(goqu.T("entitypeople")).Where(
		goqu.I("entitypeople_person_id").Eq(id),
	).Delete().ToSQL(); err != nil {
		logger.Log.Errorf("prepare remove manager: %s", err)
		return
	}

	if _, err = tx.Exec(sqlr, args...); err != nil {
		logger.Log.Errorf("remove manager: %s", err)
		return
	}

	// Remove permissions.
	if sqlr, args, err = dialect.From(goqu.T("permission")).Where(
		goqu.I("person").Eq(id),
	).Delete().ToSQL(); err != nil {
		logger.Log.Errorf("prepare remove permissions: %s", err)
		return
	}

	if _, err = tx.Exec(sqlr, args...); err != nil {
		logger.Log.Errorf("remove permissions: %s", err)
		return
	}

	// Remove borrowings.
	if sqlr, args, err = dialect.From(goqu.T("borrowing")).Where(
		goqu.I("borrower").Eq(id),
	).Delete().ToSQL(); err != nil {
		logger.Log.Errorf("prepare remove borrowings: %s", err)
		return
	}

	if _, err = tx.Exec(sqlr, args...); err != nil {
		logger.Log.Errorf("remove borrowings: %s", err)
		return
	}

	// Remove bookmarks.
	if sqlr, args, err = dialect.From(goqu.T("bookmark")).Where(
		goqu.I("person").Eq(id),
	).Delete().ToSQL(); err != nil {
		logger.Log.Errorf("prepare remove bookmarks: %s", err)
		return
	}

	if _, err = tx.Exec(sqlr, args...); err != nil {
		logger.Log.Errorf("remove bookmarks: %s", err)
		return
	}

	// Remove person.
	if sqlr, args, err = dialect.From(tablePerson).Where(
		goqu.I("person_id").Eq(id),
	).Delete().ToSQL(); err != nil {
		logger.Log.Errorf("prepare delete person: %s", err)
		return
	}

	if _, err = tx.Exec(sqlr, args...); err != nil {
		logger.Log.Errorf("delete person: %s", err)
		return
	}

	return
}

// CreatePerson creates the given person.
func (db *SQLiteDataStore) CreatePerson(p models.Person) (lastInsertID int64, err error) {
	var (
		sqlr string
		args []interface{}
		res  sql.Result
		tx   *sqlx.Tx
	)

	dialect := goqu.Dialect("sqlite3")
	tablePerson := goqu.T("person")

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

	iQuery := dialect.Insert(tablePerson).Rows(
		goqu.Record{
			"person_email": strings.ToLower(p.PersonEmail),
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

	p.PersonID = int(lastInsertID)

	// Inserting entity membership.
	for _, entity := range p.Entities {
		if sqlr, args, err = dialect.Insert(goqu.T("personentities")).Rows(
			goqu.Record{
				"personentities_person_id": p.PersonID,
				"personentities_entity_id": entity.EntityID,
			},
		).ToSQL(); err != nil {
			return
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			return
		}

		if sqlr, args, err = dialect.Insert(goqu.T("permission")).Rows(
			goqu.Record{
				"person":            p.PersonID,
				"permission_name":   "r",
				"permission_item":   "entities",
				"permission_entity": entity.EntityID,
			},
		).ToSQL(); err != nil {
			return
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			return
		}
	}

	// Inserting permissions.
	if err = db.insertPermissions(p, tx); err != nil {
		return
	}

	return
}

// UpdatePerson updates the given person.
// The password is not updated.
func (db *SQLiteDataStore) UpdatePerson(p models.Person) (err error) {
	var (
		sqlr string
		args []interface{}
		tx   *sqlx.Tx
	)

	dialect := goqu.Dialect("sqlite3")
	tablePerson := goqu.T("person")

	if tx, err = db.Beginx(); err != nil {
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

	if sqlr, args, err = dialect.Update(tablePerson).Set(
		goqu.Record{
			"person_email": strings.ToLower(p.PersonEmail),
		},
	).Where(
		goqu.I("person_id").Eq(p.PersonID),
	).ToSQL(); err != nil {
		logger.Log.Error(err)
		return
	}

	if _, err = tx.Exec(sqlr, args...); err != nil {
		return
	}

	// Lazily deleting former entities.
	if sqlr, args, err = dialect.From(goqu.T("personentities")).Where(
		goqu.I("personentities_person_id").Eq(p.PersonID),
	).Delete().ToSQL(); err != nil {
		logger.Log.Error(err)
		return
	}

	if _, err = tx.Exec(sqlr, args...); err != nil {
		return
	}

	// Lazily deleting former permissions.
	if sqlr, args, err = dialect.From(goqu.T("permission")).Where(
		goqu.I("person").Eq(p.PersonID),
	).Delete().ToSQL(); err != nil {
		logger.Log.Error(err)
		return
	}

	if _, err = tx.Exec(sqlr, args...); err != nil {
		return
	}

	// Updating person entities.
	for _, entity := range p.Entities {
		if sqlr, args, err = dialect.Insert(goqu.T("personentities")).Rows(
			goqu.Record{
				"personentities_person_id": p.PersonID,
				"personentities_entity_id": entity.EntityID,
			},
		).ToSQL(); err != nil {
			return
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			return
		}

		if sqlr, args, err = dialect.Insert(goqu.T("permission")).Rows(
			goqu.Record{
				"person":            p.PersonID,
				"permission_name":   "r",
				"permission_item":   "entities",
				"permission_entity": entity.EntityID,
			},
		).ToSQL(); err != nil {
			return
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			return
		}
	}

	// Inserting permissions.
	if err = db.insertPermissions(p, tx); err != nil {
		return
	}

	return
}

// GetAdmins returns the administrators.
func (db *SQLiteDataStore) GetAdmins() ([]models.Person, error) {
	dialect := goqu.Dialect("sqlite3")
	tablePerson := goqu.T("person")

	sQuery := dialect.From(tablePerson).Join(
		goqu.T("permission"),
		goqu.On(goqu.Ex{
			"permission.person":          goqu.I("person_id"),
			"permission.permission_name": "all",
			"permission.permission_item": "all",
			"permission_entity":          -1,
		},
		),
	).Where(
		goqu.I("person_email").Neq("admin@chimitheque.fr"),
	).Select(
		goqu.I("person_id"),
		goqu.I("person_email"),
	)

	var (
		err    error
		sqlr   string
		args   []interface{}
		people []models.Person
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return nil, err
	}

	if err = db.Select(&people, sqlr, args...); err != nil {
		return nil, err
	}

	return people, nil
}

// HasPersonReadRestrictedProductPermission returns true if the person
// can read restricted products.
func (db *SQLiteDataStore) HasPersonReadRestrictedProductPermission(id int) (bool, error) {
	var (
		err   error
		sqlr  string
		args  []interface{}
		count int
	)

	dialect := goqu.Dialect("sqlite3")
	tablePermission := goqu.T("permission")

	sQuery := dialect.From(tablePermission).Select(
		goqu.COUNT("*"),
	).Where(
		goqu.And(
			goqu.I("permission.person").Eq(id),
			goqu.Or(
				goqu.And(
					goqu.I("permission.permission_name").Eq("all"),
					goqu.I("permission.permission_item").Eq("all"),
					goqu.I("permission_entity").Eq(-1),
				),
				goqu.And(
					goqu.I("permission.permission_name").Neq("n"),
					goqu.I("permission.permission_item").Eq("rproducts"),
				),
			),
		),
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

// IsPersonAdmin returns true is the person is an admin.
func (db *SQLiteDataStore) IsPersonAdmin(id int) (bool, error) {
	var (
		err   error
		sqlr  string
		args  []interface{}
		count int
	)

	dialect := goqu.Dialect("sqlite3")
	tablePermission := goqu.T("permission")

	sQuery := dialect.From(tablePermission).Select(
		goqu.COUNT("*"),
	).Where(
		goqu.And(
			goqu.And(
				goqu.I("permission.person").Eq(id),
				goqu.I("permission.permission_name").Eq("all"),
				goqu.I("permission.permission_item").Eq("all"),
				goqu.I("permission_entity").Eq(-1),
			),
		),
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

// UnsetPersonAdmin unset the person admin permissions.
func (db *SQLiteDataStore) UnsetPersonAdmin(id int) error {
	dialect := goqu.Dialect("sqlite3")
	tablePermission := goqu.T("permission")

	dQuery := dialect.From(tablePermission).Where(
		goqu.And(
			goqu.I("person").Eq(id),
			goqu.I("permission_name").Eq("all"),
			goqu.I("permission_item").Eq("all"),
			goqu.I("permission_entity").Eq(-1),
		),
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

// SetPersonAdmin set the person an admin.
func (db *SQLiteDataStore) SetPersonAdmin(id int) error {
	var (
		err  error
		sqlr string
		args []interface{}
	)

	dialect := goqu.Dialect("sqlite3")
	tablePermission := goqu.T("permission")

	if sqlr, args, err = dialect.Insert(tablePermission).Rows(
		goqu.Record{
			"person":            id,
			"permission_name":   "all",
			"permission_item":   "all",
			"permission_entity": -1,
		},
	).OnConflict(goqu.DoNothing()).ToSQL(); err != nil {
		return nil
	}

	if _, err = db.Exec(sqlr, args...); err != nil {
		return err
	}

	return nil
}

func (db *SQLiteDataStore) insertPermissions(p models.Person, tx *sqlx.Tx) error {
	var (
		sqlr string
		args []interface{}
		err  error
	)

	dialect := goqu.Dialect("sqlite3")
	tablePermission := goqu.T("permission")

	if len(p.Permissions) == 0 {
		// Inserting default permission.
		p.Permissions = append(p.Permissions, &models.Permission{
			PermissionName:   "r",
			PermissionItem:   "products",
			PermissionEntity: -1,
			Person: models.Person{
				PersonID: p.PersonID,
			},
		})
	}

	for _, perm := range p.Permissions {
		iQuery := dialect.Insert(tablePermission).Rows(
			goqu.Record{
				"person":            p.PersonID,
				"permission_name":   perm.PermissionName,
				"permission_item":   perm.PermissionItem,
				"permission_entity": perm.PermissionEntity,
			},
		)

		if sqlr, args, err = iQuery.ToSQL(); err != nil {
			return err
		}

		if _, err = tx.Exec(sqlr, args...); err != nil {
			return err
		}
	}

	return nil
}

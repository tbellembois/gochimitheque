package datastores

import (
	"database/sql"
	"strings"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/zmqclient"
)

// GetPeople select the people matching p
// and visible by the connected user.
func (db *SQLiteDataStore) GetPeople(f zmqclient.RequestFilter, person_id int) ([]models.Person, int, error) {
	var err error

	logger.Log.WithFields(logrus.Fields{"f": f}).Debug("GetPeople")

	if f.OrderBy == "" {
		f.OrderBy = "person_id"
	}

	dialect := goqu.Dialect("sqlite3")
	tablePerson := goqu.T("person")
	tableEntity := goqu.T("entity")

	// Build orderby/order clause.
	orderByClause := f.OrderBy
	orderClause := goqu.I(orderByClause).Asc()

	if strings.ToLower(f.Order) == "desc" {
		orderClause = goqu.I(orderByClause).Desc()
	}

	// Is the logged user an admin?
	// We need to handle admins to see people with no entities.
	var isadmin bool

	if isadmin, err = db.IsPersonAdmin(person_id); err != nil {
		return nil, 0, err
	}

	// Build join clause.
	var joinClause *goqu.SelectDataset

	switch {
	case f.Entity != 0:
		joinClause = dialect.From(tablePerson.As("p"), tableEntity.As("e")).Join(
			goqu.T("personentities"),
			goqu.On(
				goqu.Ex{
					"personentities.personentities_person_id": goqu.I("p.person_id"),
				},
			),
		).Join(
			goqu.T("entity"),
			goqu.On(
				goqu.Ex{
					"personentities.personentities_entity_id": f.Entity,
				},
			),
		)
	case !isadmin:
		joinClause = dialect.From(tablePerson.As("p"), tableEntity.As("e")).Join(
			goqu.T("personentities"),
			goqu.On(
				goqu.Ex{
					"personentities.personentities_person_id": goqu.I("p.person_id"),
				},
			),
		).Join(
			goqu.T("entity"),
			goqu.On(
				goqu.Ex{
					"personentities.personentities_entity_id": goqu.I("e.entity_id"),
				},
			),
		).Join(
			goqu.T("permission").As("perm"),
			goqu.On(
				goqu.Ex{
					"perm.person":               person_id,
					"perm.permission_item_name": []string{"all", "people"},
					"perm.permission_perm_name": []string{"all", "r", "w"},
					"perm.permission_entity_id": []interface{}{-1, goqu.I("e.entity_id")},
				},
			),
		)
	default:
		joinClause = dialect.From(tablePerson.As("p"), tableEntity.As("e"))
	}

	joinClause = joinClause.Where(
		goqu.I("p.person_email").Like(f.Search),
	)

	// Building final count.
	var (
		countSQL  string
		countArgs []interface{}
	)

	if countSQL, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("p.person_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}

	// Building final select.
	var (
		selectSQL  string
		selectArgs []interface{}
	)

	if selectSQL, selectArgs, err = joinClause.Select(
		goqu.I("p.person_id"),
		goqu.I("p.person_email"),
	).GroupBy(goqu.I("p.person_id")).Order(orderClause).Limit(uint(f.Limit)).Offset(uint(f.Offset)).ToSQL(); err != nil {
		return nil, 0, err
	}

	var (
		people []models.Person
		count  int
	)

	if err = db.Select(&people, selectSQL, selectArgs...); err != nil {
		return nil, 0, err
	}

	if err = db.Get(&count, countSQL, countArgs...); err != nil {
		return nil, 0, err
	}

	return people, count, nil
}

// GetPerson select the person by id.
func (db *SQLiteDataStore) GetPerson(id int) (models.Person, error) {
	dialect := goqu.Dialect("sqlite3")
	tablePerson := goqu.T("person")

	sQuery := dialect.From(tablePerson).Where(
		goqu.I("person_id").Eq(id),
	).Select(
		goqu.I("person_id"),
		goqu.I("person_email"),
	)

	var (
		err    error
		sqlr   string
		args   []interface{}
		person models.Person
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return models.Person{}, err
	}

	if err = db.Get(&person, sqlr, args...); err != nil {
		return models.Person{}, err
	}

	return person, nil
}

// GetPersonByEmail select the person by email.
func (db *SQLiteDataStore) GetPersonByEmail(email string) (models.Person, error) {
	email = strings.ToLower(email)

	dialect := goqu.Dialect("sqlite3")
	tablePerson := goqu.T("person")

	sQuery := dialect.From(tablePerson).Where(
		goqu.I("person_email").Eq(email),
	).Select(
		goqu.I("person_id"),
		goqu.I("person_email"),
	)

	var (
		err    error
		sqlr   string
		args   []interface{}
		person models.Person
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return models.Person{}, err
	}

	if err = db.Get(&person, sqlr, args...); err != nil {
		return models.Person{}, err
	}

	return person, nil
}

// GetPersonPermissions return person permissions.
func (db *SQLiteDataStore) GetPersonPermissions(id int) ([]models.Permission, error) {
	dialect := goqu.Dialect("sqlite3")
	tablePermission := goqu.T("permission")

	sQuery := dialect.From(tablePermission).Where(
		goqu.I("person").Eq(id),
	).Select(
		goqu.I("permission_id"),
		goqu.I("permission_perm_name"),
		goqu.I("permission_item_name"),
		goqu.I("permission_entity_id"),
	)

	var (
		err         error
		sqlr        string
		args        []interface{}
		permissions []models.Permission
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return nil, err
	}

	if err = db.Select(&permissions, sqlr, args...); err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetPersonManageEntities returns the entities the person if manager of.
func (db *SQLiteDataStore) GetPersonManageEntities(id int) ([]models.Entity, error) {
	dialect := goqu.Dialect("sqlite3")
	tableEntity := goqu.T("entity")

	sQuery := dialect.From(tableEntity).LeftJoin(
		goqu.T("entitypeople"),
		goqu.On(goqu.Ex{"entitypeople.entitypeople_entity_id": goqu.I("entity.entity_id")}),
	).Where(
		goqu.I("entitypeople.entitypeople_person_id").Eq(id),
	).Select(
		goqu.I("entity_id"),
		goqu.I("entity_name"),
		goqu.I("entity_description"),
	)

	var (
		err      error
		sqlr     string
		args     []interface{}
		entities []models.Entity
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return nil, err
	}

	if err = db.Select(&entities, sqlr, args...); err != nil {
		return nil, err
	}

	return entities, nil
}

// GetPeople select the person entities
// and visible by the connected user.
func (db *SQLiteDataStore) GetPersonEntities(loggedPersonID int, personID int) ([]models.Entity, error) {
	var err error

	dialect := goqu.Dialect("sqlite3")
	tablePerson := goqu.T("person")
	tableEntity := goqu.T("entity")
	tablePersonentities := goqu.T("personentities")

	// Is the logged user an admin?
	var isadmin bool

	if isadmin, err = db.IsPersonAdmin(loggedPersonID); err != nil {
		return nil, err
	}

	// Build join clause.
	var joinClause *goqu.SelectDataset
	if !isadmin {
		joinClause = dialect.From(
			tableEntity.As("e"),
			tablePerson.As("p"),
			tablePersonentities.As("pe"),
		).Join(
			goqu.T("permission").As("perm"),
			goqu.On(
				goqu.Or(
					goqu.And(
						goqu.I("perm.person").Eq(personID),
						goqu.I("perm.permission_item_name").Eq("all"),
						goqu.I("perm.permission_perm_name").Eq("all"),
						goqu.I("perm.permission_entity_id").Eq(goqu.I("e.entity_id")),
					),
					goqu.And(
						goqu.I("perm.person").Eq(personID),
						goqu.I("perm.permission_item_name").Eq("all"),
						goqu.I("perm.permission_perm_name").Eq("all"),
						goqu.I("perm.permission_entity_id").Eq(-1),
					),
					goqu.And(
						goqu.I("perm.person").Eq(personID),
						goqu.I("perm.permission_item_name").Eq("all"),
						goqu.I("perm.permission_perm_name").Eq("r"),
						goqu.I("perm.permission_entity_id").Eq(-1),
					),
					goqu.And(
						goqu.I("perm.person").Eq(personID),
						goqu.I("perm.permission_item_name").Eq("entities"),
						goqu.I("perm.permission_perm_name").Eq("all"),
						goqu.I("perm.permission_entity_id").Eq(goqu.I("e.entity_id")),
					),
					goqu.And(
						goqu.I("perm.person").Eq(personID),
						goqu.I("perm.permission_item_name").Eq("entities"),
						goqu.I("perm.permission_perm_name").Eq("all"),
						goqu.I("perm.permission_entity_id").Eq(-1),
					),
					goqu.And(
						goqu.I("perm.person").Eq(personID),
						goqu.I("perm.permission_item_name").Eq("entities"),
						goqu.I("perm.permission_perm_name").Eq("r"),
						goqu.I("perm.permission_entity_id").Eq(-1),
					),
					goqu.And(
						goqu.I("perm.person").Eq(personID),
						goqu.I("perm.permission_item_name").Eq("entities"),
						goqu.I("perm.permission_perm_name").Eq("r"),
						goqu.I("perm.permission_entity_id").Eq(goqu.I("e.entity_id")),
					),
				),
			),
		)
	} else {
		joinClause = dialect.From(
			tableEntity.As("e"),
			tablePerson.As("p"),
			tablePersonentities.As("pe"),
		)
	}

	joinClause = joinClause.Where(
		goqu.Ex{
			"pe.personentities_person_id": personID,
			"e.entity_id":                 goqu.I("pe.personentities_entity_id"),
		},
	).GroupBy(goqu.I("e.entity_id")).Order(goqu.I("e.entity_name").Asc())

	var (
		sqlr     string
		args     []interface{}
		entities []models.Entity
	)

	if sqlr, args, err = joinClause.Select(
		goqu.I("e.entity_id"),
		goqu.I("e.entity_name"),
		goqu.I("e.entity_description"),
	).ToSQL(); err != nil {
		logger.Log.Error(err)
		return nil, err
	}

	if err = db.Select(&entities, sqlr, args...); err != nil {
		return nil, err
	}

	return entities, nil
}

// DoesPersonBelongsTo returns true if the person belongs to the entities.
func (db *SQLiteDataStore) DoesPersonBelongsTo(id int, entities []models.Entity) (bool, error) {
	var (
		err   error
		sqlr  string
		args  []interface{}
		count int
	)

	dialect := goqu.Dialect("sqlite3")
	tablePersonentities := goqu.T("personentities")

	var entityIds []int
	for _, i := range entities {
		entityIds = append(entityIds, i.EntityID)
	}

	sQuery := dialect.From(tablePersonentities).Select(
		goqu.COUNT("*"),
	).Where(
		goqu.Ex{
			"personentities_person_id": id,
			"personentities_entity_id": entityIds,
		},
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
	var admin models.Person

	if admin, err = db.GetPersonByEmail("admin@chimitheque.fr"); err != nil {
		return err
	}

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
				"person":               p.PersonID,
				"permission_perm_name": "r",
				"permission_item_name": "entities",
				"permission_entity_id": entity.EntityID,
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
				"person":               p.PersonID,
				"permission_perm_name": "r",
				"permission_item_name": "entities",
				"permission_entity_id": entity.EntityID,
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
			"permission.person":               goqu.I("person_id"),
			"permission.permission_perm_name": "all",
			"permission.permission_item_name": "all",
			"permission_entity_id":            -1,
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
					goqu.I("permission.permission_perm_name").Eq("all"),
					goqu.I("permission.permission_item_name").Eq("all"),
					goqu.I("permission_entity_id").Eq(-1),
				),
				goqu.And(
					goqu.I("permission.permission_perm_name").Neq("n"),
					goqu.I("permission.permission_item_name").Eq("rproducts"),
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
				goqu.I("permission.permission_perm_name").Eq("all"),
				goqu.I("permission.permission_item_name").Eq("all"),
				goqu.I("permission_entity_id").Eq(-1),
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
			goqu.I("permission_perm_name").Eq("all"),
			goqu.I("permission_item_name").Eq("all"),
			goqu.I("permission_entity_id").Eq(-1),
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
			"person":               id,
			"permission_perm_name": "all",
			"permission_item_name": "all",
			"permission_entity_id": -1,
		},
	).OnConflict(goqu.DoNothing()).ToSQL(); err != nil {
		return nil
	}

	if _, err = db.Exec(sqlr, args...); err != nil {
		return err
	}

	return nil
}

// IsPersonManager returns true is the person is a manager.
func (db *SQLiteDataStore) IsPersonManager(id int) (bool, error) {
	var (
		err   error
		sqlr  string
		args  []interface{}
		count int
	)

	dialect := goqu.Dialect("sqlite3")
	tableEntitypeople := goqu.T("entitypeople")

	sQuery := dialect.From(tableEntitypeople).Select(
		goqu.COUNT("*"),
	).Where(
		goqu.I("entitypeople.entitypeople_person_id").Eq(id),
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
			PermissionPermName: "r",
			PermissionItemName: "products",
			PermissionEntityID: -1,
			Person: models.Person{
				PersonID: p.PersonID,
			},
		})
	}

	for _, perm := range p.Permissions {
		iQuery := dialect.Insert(tablePermission).Rows(
			goqu.Record{
				"person":               p.PersonID,
				"permission_perm_name": perm.PermissionPermName,
				"permission_item_name": perm.PermissionItemName,
				"permission_entity_id": perm.PermissionEntityID,
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

package models

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver
	log "github.com/sirupsen/logrus"
)

const (
	dbdriver = "sqlite3"
)

// SQLiteDataStore implements the Datastore interface
// to store data in SQLite3
type SQLiteDataStore struct {
	*sqlx.DB
	err error
}

// NewDBstore returns a database connection to the given dataSourceName
// ie. a path to the sqlite database file
func NewDBstore(dataSourceName string) (*SQLiteDataStore, error) {
	var (
		db  *sqlx.DB
		err error
	)

	log.WithFields(log.Fields{"dbdriver": dbdriver, "dataSourceName": dataSourceName}).Debug("NewDBstore")
	if db, err = sqlx.Connect(dbdriver, dataSourceName); err != nil {
		return &SQLiteDataStore{}, err
	}
	return &SQLiteDataStore{db, nil}, nil
}

// FlushErrors returns the last DB errors and flushes it.
func (db *SQLiteDataStore) FlushErrors() error {
	// saving the last thrown error
	lastError := db.err
	// resetting the error
	db.err = nil
	// returning the last error
	return lastError
}

// CreateDatabase creates the database tables
func (db *SQLiteDataStore) CreateDatabase() error {
	// activate the foreign keys feature
	if _, db.err = db.Exec("PRAGMA foreign_keys = ON"); db.err != nil {
		return db.err
	}

	// schema definition
	schema := `CREATE TABLE IF NOT EXISTS person(
		person_id integer PRIMARY KEY,
		person_email string NOT NULL,
		person_password string NOT NULL);
	CREATE TABLE IF NOT EXISTS entity (
		entity_id integer PRIMARY KEY,
		entity_name string NOT NULL,
		entity_description string,
		entity_person_id integer,
		FOREIGN KEY (entity_person_id) references person(person_id));
	CREATE TABLE IF NOT EXISTS permission (
		permission_id integer PRIMARY KEY,
		permission_person_id integer NOT NULL,
		permission_perm_name string NOT NULL,
		permission_item_name string NOT NULL,
		permission_itemid integer,
		FOREIGN KEY (permission_person_id) references person(person_id));
	CREATE TABLE IF NOT EXISTS personentities (
		personentities_person_id integer NOT NULL,
		personentities_entity_id integer NOT NULL,
		FOREIGN KEY (personentities_person_id) references person(person_id),
		FOREIGN KEY (personentities_entity_id) references entity(entity_id));`

	// tables creation
	if _, db.err = db.Exec(schema); db.err != nil {
		return db.err
	}

	// inserting sample values if tables are empty
	var c int
	_ = db.Get(&c, `SELECT count(*) FROM person`)
	log.WithFields(log.Fields{"c": c}).Debug("CreateDatabase")
	if c == 0 {
		log.Debug("populating database")
		// preparing requests
		people := `INSERT INTO person (person_email, person_password) VALUES (?, ?)`
		entities := `INSERT INTO entity (entity_name, entity_description, entity_person_id) VALUES (?, ?, ?)`
		permissions := `INSERT INTO permission (permission_person_id, permission_perm_name, permission_item_name, permission_itemid) VALUES (?, ?, ?, ?)`
		personentities := `INSERT INTO personentities (personentities_person_id, personentities_entity_id) VALUES (? ,?)`
		// inserting people
		res1 := db.MustExec(people, "john.doe@foo.com", "johndoe")
		res2 := db.MustExec(people, "mickey.mouse@foo.com", "mickeymouse")
		res3 := db.MustExec(people, "obione.kenobi@foo.com", "obionekenobi")
		res4 := db.MustExec(people, "dark.vader@foo.com", "darkvader")
		// getting last inserted ids
		johnid, _ := res1.LastInsertId()
		mickeyid, _ := res2.LastInsertId()
		obioneid, _ := res3.LastInsertId()
		darkid, _ := res4.LastInsertId()
		// inserted entities and permissions
		res5 := db.MustExec(entities, "entity1", "sample entity one", johnid)
		res6 := db.MustExec(entities, "entity2", "sample entity two", mickeyid)
		res7 := db.MustExec(entities, "entity3", "sample entity three", obioneid)
		// getting last inserted ids
		entity1id, _ := res5.LastInsertId()
		entity2id, _ := res6.LastInsertId()
		entity3id, _ := res7.LastInsertId()
		db.MustExec(permissions, johnid, "r", "entities", nil)
		db.MustExec(permissions, johnid, "r", "entity", entity1id)
		db.MustExec(permissions, johnid, "w", "entity", entity1id)
		db.MustExec(permissions, mickeyid, "r", "entities", nil)
		db.MustExec(permissions, mickeyid, "r", "entity", entity2id)
		db.MustExec(permissions, mickeyid, "w", "entity", entity2id)
		db.MustExec(permissions, obioneid, "all", "all", nil)
		db.MustExec(permissions, darkid, "r", "entities", nil)
		// then people entities
		db.MustExec(personentities, johnid, entity1id)
		db.MustExec(personentities, mickeyid, entity2id)
		db.MustExec(personentities, obioneid, entity1id)
		db.MustExec(personentities, obioneid, entity2id)
		db.MustExec(personentities, obioneid, entity3id)
		db.MustExec(personentities, darkid, entity1id)
		db.MustExec(personentities, darkid, entity2id)
		db.MustExec(personentities, darkid, entity3id)
	}
	return nil
}

func (db *SQLiteDataStore) DeleteEntity(id int) error {
	var (
		sqlr string
	)
	sqlr = `DELETE FROM entity 
	WHERE entity_id = ?`
	if _, db.err = db.Exec(sqlr, id); db.err != nil {
		return db.err
	}
	return nil
}

func (db *SQLiteDataStore) UpdateEntity(e Entity) error {
	var (
		sqlr string
	)
	sqlr = `UPDATE entity SET entity_name = ?, entity_description = ?, entity_person_id = ?
	WHERE entity_id = ?`
	if _, db.err = db.Exec(sqlr, e.EntityName, e.EntityDescription, e.PersonID, e.EntityID); db.err != nil {
		return db.err
	}
	return nil
}

func (db *SQLiteDataStore) CreateEntity(e Entity) error {
	var (
		sqlr string
	)
	sqlr = `INSERT INTO entity(entity_name, entity_description, entity_person_id) VALUES (?, ?, ?)`
	if _, db.err = db.Exec(sqlr, e.EntityName, e.EntityDescription, e.Person.PersonID); db.err != nil {
		return db.err
	}
	return nil
}

func (db *SQLiteDataStore) GetEntities(search string, order string, offset uint64, limit uint64) ([]Entity, error) {
	var (
		entities []Entity
		sqlr     string
		sqla     []interface{}
	)
	log.WithFields(log.Fields{"search": search, "order": order, "offset": offset, "limit": limit}).Debug("GetEntities")

	sqlr, sqla, db.err = sq.Select(`e.entity_id, 
		e.entity_name, 
		e.entity_description, 
		p.person_id, 
		p.person_email, 
		p.person_password`).
		From("entity AS e, person AS p").
		Where("e.entity_person_id = p.person_id AND e.entity_name LIKE ?", fmt.Sprint("%", search, "%")).
		OrderBy(fmt.Sprintf("entity_name %s", order)).
		Offset(offset).
		Limit(limit).ToSql()

	if db.err != nil {
		return nil, db.err
	}

	if db.err = db.Select(&entities, sqlr, sqla...); db.err != nil {
		return nil, db.err
	}
	return entities, nil
}

func (db *SQLiteDataStore) GetEntity(ID int) (Entity, error) {
	var (
		entity Entity
		sqlr   string
	)

	sqlr = "SELECT e.entity_id, e.entity_name, e.entity_description, p.person_id, p.person_email, p.person_password FROM entity AS e, person AS p WHERE e.entity_person_id = p.person_id AND e.entity_id = ?"
	if db.err = db.Get(&entity, sqlr, ID); db.err != nil {
		return Entity{}, db.err
	}
	log.WithFields(log.Fields{"ID": ID, "entity": entity}).Debug("GetEntity")
	return entity, nil
}

func (db *SQLiteDataStore) HasEntityWithNameExcept(name string, except ...string) (bool, error) {
	var (
		res   bool
		count int
		sqlr  sq.SelectBuilder
		w     sq.And
	)

	w = append(w, sq.Eq{"entity.entity_name": name})
	for _, e := range except {
		w = append(w, sq.NotEq{"entity.entity_name": e})
	}

	sqlr = sq.Select("count(*)").From("entity").Where(w)
	sql, args, _ := sqlr.ToSql()
	log.WithFields(log.Fields{"sql": sql, "args": args}).Debug("HasEntityWithNameExcept")

	if db.err = db.Get(&count, sql, args...); db.err != nil {
		return false, db.err
	}
	log.WithFields(log.Fields{"name": name, "count": count}).Debug("HasEntityWithNameExcept")
	if count == 0 {
		res = false
	} else {
		res = true
	}
	return res, nil
}

func (db *SQLiteDataStore) HasEntityWithName(name string) (bool, error) {
	var (
		res   bool
		count int
		sqlr  string
	)

	sqlr = "SELECT count(*) from entity WHERE entity.entity_name = ?"
	if db.err = db.Get(&count, sqlr, name); db.err != nil {
		return false, db.err
	}
	log.WithFields(log.Fields{"name": name, "count": count}).Debug("HasEntityWithName")
	if count == 0 {
		res = false
	} else {
		res = true
	}
	return res, nil
}

func (db *SQLiteDataStore) HasPermission(id int, perm string, item string, itemid int) (bool, error) {
	var (
		res   bool
		count int
		sqlr  string
	)

	log.WithFields(log.Fields{
		"id":     id,
		"perm":   perm,
		"item":   item,
		"itemid": itemid}).Debug("HasPermission")

	// then counting the permissions matching the parameters
	if itemid == -1 {
		sqlr = `SELECT count(*) FROM permission WHERE 
		permission_person_id = ? AND permission_perm_name = ? AND permission_item_name = ? OR 
		permission_person_id = ? AND permission_perm_name = "all" AND permission_item_name = ? OR
		permission_person_id = ? AND permission_perm_name = ? AND permission_item_name = "all" OR
		permission_person_id = ? AND permission_perm_name = "all" AND permission_item_name = "all"`
		if db.err = db.Get(&count, sqlr, id, perm, item, id, item, id, perm, id); db.err != nil {
			switch {
			case db.err == sql.ErrNoRows:
				return false, nil
			default:
				return false, db.err
			}
		}
	} else {
		sqlr = `SELECT count(*) FROM permission WHERE 
		permission_person_id = ? AND permission_perm_name = ? AND permission_item_name = ? AND permission_itemid = ? OR 
		permission_person_id = ? AND permission_perm_name = ? AND permission_item_name = "all" AND permission_itemid = ? OR
		permission_person_id = ? AND permission_perm_name = "all" AND permission_item_name = ? AND permission_itemid = ? OR
		permission_person_id = ? AND permission_perm_name = "all" AND permission_item_name = "all"`
		if db.err = db.Get(&count, sqlr, id, perm, item, itemid, id, perm, itemid, id, item, itemid, id); db.err != nil {
			switch {
			case db.err == sql.ErrNoRows:
				return false, nil
			default:
				return false, db.err
			}
		}
	}
	log.WithFields(log.Fields{"count": count}).Debug("HasPermission")

	if count == 0 {
		res = false
	} else {
		res = true
	}
	return res, nil
}

func (db *SQLiteDataStore) GetPeople() ([]Person, error) {
	var (
		people []Person
		sqlr   string
	)

	sqlr = "SELECT person_id, person_email FROM person"
	if db.err = db.Select(&people, sqlr); db.err != nil {
		return nil, db.err
	}
	return people, nil
}

func (db *SQLiteDataStore) GetPerson(personID int) (Person, error) {
	var (
		person Person
		sqlr   string
	)

	sqlr = "SELECT person_id, person_email FROM person WHERE person_id = ?"
	if db.err = db.Get(&person, sqlr, personID); db.err != nil {
		return Person{}, db.err
	}
	return person, nil
}

func (db *SQLiteDataStore) GetPersonByEmail(personEmail string) (Person, error) {
	var (
		person Person
		sqlr   string
	)

	sqlr = "SELECT person_id, person_email FROM person WHERE person_email = ?"
	if db.err = db.Get(&person, sqlr, personEmail); db.err != nil {
		return Person{}, db.err
	}
	return person, nil
}

func (db *SQLiteDataStore) GetPersonPermissions(personID int) ([]Permission, error) {
	var (
		ps   []Permission
		sqlr string
	)

	sqlr = `SELECT permission_id, permission_perm_name, permission_item_name, permission_itemid 
	FROM permission
	WHERE permission_person_id = ?`
	if db.err = db.Select(&ps, sqlr, personID); db.err != nil {
		return nil, db.err
	}
	log.WithFields(log.Fields{"personID": personID, "ps": ps}).Debug("GetPersonPermissions")
	return ps, nil
}

func (db *SQLiteDataStore) GetPersonEntities(personID int) ([]Entity, error) {
	var (
		es   []Entity
		sqlr string
	)

	sqlr = `SELECT entity_id, entity_name, entity_description 
	FROM entity
	INNER JOIN personentities ON personentities.personentities_entity_id = entity.entity_id
	WHERE personentities.personentities_person_id = ?`
	if db.err = db.Select(&es, sqlr, personID); db.err != nil {
		return nil, db.err
	}
	log.WithFields(log.Fields{"personID": personID, "es": es}).Debug("GetPersonEntities")
	return es, nil
}

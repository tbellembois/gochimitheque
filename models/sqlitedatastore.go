package models

import (
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
		db.MustExec(permissions, johnid, "r", "entity", entity1id)
		db.MustExec(permissions, johnid, "w", "entity", entity1id)
		db.MustExec(permissions, mickeyid, "r", "entity", -1)
		db.MustExec(permissions, mickeyid, "r", "entity", entity2id)
		db.MustExec(permissions, mickeyid, "w", "entity", entity2id)
		db.MustExec(permissions, obioneid, "all", "all", -1)
		db.MustExec(permissions, darkid, "r", "entity", -1)
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

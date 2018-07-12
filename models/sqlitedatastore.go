package models

import (
	"fmt"
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

// buildPermissionFilter return the sql join to return only items of tableName that the person permission_person_id can permission_perm_name
func buildPermissionFilter(tableName string, tableAlias string, tableJoinField string, permName string) string {
	return fmt.Sprintf(`permission AS perm on perm.permission_person_id = ? and (
		(perm.permission_item_name = "all" and perm.permission_perm_name = "all") or
		(perm.permission_item_name == "all" and perm.permission_perm_name == "%s" and perm.permission_entityid == -1) or
		(perm.permission_item_name == "%s" and perm.permission_perm_name == "all" and perm.permission_entityid == %s.%s) or
		(perm.permission_item_name == "%s" and perm.permission_perm_name == "all" and perm.permission_entityid == -1) or
		(perm.permission_item_name == "%s" and perm.permission_perm_name == "%s" and perm.permission_entityid == -1) or
		(perm.permission_item_name == "%s" and perm.permission_perm_name == "%s" and perm.permission_entityid == %s.%s)
		)`, permName, tableName, tableAlias, tableJoinField, tableName, tableName, permName, tableName, permName, tableAlias, tableJoinField)
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
		entity_description string);
	CREATE TABLE IF NOT EXISTS permission (
		permission_id integer PRIMARY KEY,
		permission_person_id integer NOT NULL,
		permission_perm_name string NOT NULL,
		permission_item_name string NOT NULL,
		permission_entityid integer,
		FOREIGN KEY (permission_person_id) references person(person_id));
	-- entities people belongs to
	CREATE TABLE IF NOT EXISTS personentities (
		personentities_person_id integer NOT NULL,
		personentities_entity_id integer NOT NULL,
		FOREIGN KEY (personentities_person_id) references person(person_id),
		FOREIGN KEY (personentities_entity_id) references entity(entity_id));
	-- entities managers	
	CREATE TABLE IF NOT EXISTS entitypeople (
		entitypeople_entity_id integer NOT NULL,
		entitypeople_person_id integer NOT NULL,
		PRIMARY KEY (entitypeople_entity_id, entitypeople_person_id),
		FOREIGN KEY (entitypeople_person_id) references person(person_id),
		FOREIGN KEY (entitypeople_entity_id) references entity(entity_id));
	`

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
		entities := `INSERT INTO entity (entity_name, entity_description) VALUES (?, ?)`
		permissions := `INSERT INTO permission (permission_person_id, permission_perm_name, permission_item_name, permission_entityid) VALUES (?, ?, ?, ?)`
		personentities := `INSERT INTO personentities (personentities_person_id, personentities_entity_id) VALUES (? ,?)`
		entitypeople := `INSERT INTO entitypeople (entitypeople_entity_id, entitypeople_person_id) VALUES (? ,?)`
		// inserting people
		user1, _ := db.MustExec(people, "user1@entity1.com", "user").LastInsertId()
		user11, _ := db.MustExec(people, "user11@entity1.com", "user").LastInsertId()
		userm1, _ := db.MustExec(people, "manager1@entity1.com", "user").LastInsertId()

		user2, _ := db.MustExec(people, "user2@entity2.com", "user").LastInsertId()
		user22, _ := db.MustExec(people, "user22@entity2.com", "user").LastInsertId()
		userm2, _ := db.MustExec(people, "manager2@entity2.com", "user").LastInsertId()

		user3, _ := db.MustExec(people, "user3@entity3.com", "user").LastInsertId()
		user33, _ := db.MustExec(people, "user33@entity3.com", "user").LastInsertId()
		userm3, _ := db.MustExec(people, "manager3@entity3.com", "user").LastInsertId()
		userm33, _ := db.MustExec(people, "manager33@entity3.com", "user").LastInsertId()

		usersuper, _ := db.MustExec(people, "user@super.com", "user").LastInsertId()
		// inserting entities
		entity1id, _ := db.MustExec(entities, "entity1", "sample entity one").LastInsertId()
		entity2id, _ := db.MustExec(entities, "entity2", "sample entity two").LastInsertId()
		entity3id, _ := db.MustExec(entities, "entity3", "sample entity three").LastInsertId()

		// setting up permissions
		// entity1 users
		db.MustExec(permissions, userm1, "r", "entity", entity1id)
		db.MustExec(permissions, userm1, "w", "entity", entity1id)
		db.MustExec(permissions, userm1, "r", "person", entity1id)
		db.MustExec(permissions, userm1, "w", "person", entity1id)
		db.MustExec(permissions, userm1, "all", "product", -1)
		db.MustExec(permissions, userm1, "all", "storage", entity1id)

		db.MustExec(permissions, user1, "r", "product", -1)
		db.MustExec(permissions, user1, "r", "entity", entity1id)
		db.MustExec(permissions, user1, "r", "storage", entity1id)

		db.MustExec(permissions, user11, "r", "product", -1)
		db.MustExec(permissions, user11, "r", "storage", entity1id)
		db.MustExec(permissions, user11, "w", "storage", entity1id)

		// entity2 users
		db.MustExec(permissions, userm2, "r", "entity", entity2id)
		db.MustExec(permissions, userm2, "w", "entity", entity2id)
		db.MustExec(permissions, userm2, "r", "person", entity2id)
		db.MustExec(permissions, userm2, "w", "person", entity2id)
		db.MustExec(permissions, userm2, "all", "product", -1)
		db.MustExec(permissions, userm2, "all", "storage", entity2id)

		db.MustExec(permissions, user2, "r", "product", -1)
		db.MustExec(permissions, user2, "r", "entity", entity2id)
		db.MustExec(permissions, user2, "r", "storage", entity2id)

		db.MustExec(permissions, user22, "r", "product", -1)
		db.MustExec(permissions, user22, "r", "storage", entity2id)
		db.MustExec(permissions, user22, "w", "storage", entity2id)

		// entity3 users
		db.MustExec(permissions, userm3, "r", "entity", entity3id)
		db.MustExec(permissions, userm3, "w", "entity", entity3id)
		db.MustExec(permissions, userm3, "r", "person", entity3id)
		db.MustExec(permissions, userm3, "w", "person", entity3id)
		db.MustExec(permissions, userm3, "all", "product", -1)
		db.MustExec(permissions, userm3, "all", "storage", entity3id)

		db.MustExec(permissions, userm33, "r", "entity", entity3id)
		db.MustExec(permissions, userm33, "w", "entity", entity3id)
		db.MustExec(permissions, userm33, "r", "person", entity3id)
		db.MustExec(permissions, userm33, "w", "person", entity3id)
		db.MustExec(permissions, userm33, "all", "product", -1)
		db.MustExec(permissions, userm33, "all", "storage", entity3id)

		db.MustExec(permissions, user3, "r", "product", -1)
		db.MustExec(permissions, user3, "r", "storage", entity3id)

		db.MustExec(permissions, user33, "r", "product", -1)
		db.MustExec(permissions, user33, "r", "storage", entity3id)
		db.MustExec(permissions, user33, "w", "storage", entity3id)

		// super admin
		db.MustExec(permissions, usersuper, "all", "all", -1)

		// then people entities
		db.MustExec(personentities, user1, entity1id)
		db.MustExec(personentities, user11, entity1id)
		db.MustExec(personentities, userm1, entity1id)
		db.MustExec(personentities, user2, entity2id)
		db.MustExec(personentities, user22, entity2id)
		db.MustExec(personentities, userm2, entity2id)
		db.MustExec(personentities, user3, entity3id)
		db.MustExec(personentities, user33, entity3id)
		db.MustExec(personentities, userm3, entity3id)
		db.MustExec(personentities, userm33, entity3id)

		// then entities managers
		db.MustExec(entitypeople, entity1id, userm1)
		db.MustExec(entitypeople, entity2id, userm2)
		db.MustExec(entitypeople, entity3id, userm3)
		db.MustExec(entitypeople, entity3id, userm33)
	}
	return nil
}

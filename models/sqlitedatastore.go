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
		entity_description string);
	CREATE TABLE IF NOT EXISTS storelocation (
		storelocation_id integer PRIMARY KEY,
		storelocation_name string NOT NULL,
		storelocation_entity_id integer NOT NULL,
		FOREIGN KEY(storelocation_entity_id) references entity(entity_id));
	CREATE TABLE IF NOT EXISTS permission (
		permission_id integer PRIMARY KEY,
		permission_person_id integer NOT NULL,
		permission_perm_name string NOT NULL,
		permission_item_name string NOT NULL,
		permission_entity_id integer,
		FOREIGN KEY(permission_person_id) references person(person_id));
	-- entities people belongs to
	CREATE TABLE IF NOT EXISTS personentities (
		personentities_person_id integer NOT NULL,
		personentities_entity_id integer NOT NULL,
		PRIMARY KEY(personentities_person_id, personentities_entity_id),
		FOREIGN KEY(personentities_person_id) references person(person_id),
		FOREIGN KEY(personentities_entity_id) references entity(entity_id));
	-- entities managers	
	CREATE TABLE IF NOT EXISTS entitypeople (
		entitypeople_entity_id integer NOT NULL,
		entitypeople_person_id integer NOT NULL,
		PRIMARY KEY(entitypeople_entity_id, entitypeople_person_id),
		FOREIGN KEY(entitypeople_person_id) references person(person_id),
		FOREIGN KEY(entitypeople_entity_id) references entity(entity_id));
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

		m1 := Person{PersonEmail: "manager@lab-one.com"}
		m2 := Person{PersonEmail: "manager@lab-two.com"}
		m3 := Person{PersonEmail: "manager@lab-three.com"}

		_, m1.PersonID = db.CreatePerson(m1)
		_, m2.PersonID = db.CreatePerson(m1)
		_, m3.PersonID = db.CreatePerson(m1)

		e1 := Entity{EntityName: "lab one", Managers: []Person{m1}}
		e2 := Entity{EntityName: "lab two", Managers: []Person{m2}}
		e3 := Entity{EntityName: "lab three", Managers: []Person{m3}}

		_, e1.EntityID = db.CreateEntity(e1)
		_, e2.EntityID = db.CreateEntity(e2)
		_, e3.EntityID = db.CreateEntity(e3)

		sl1 := StoreLocation{StoreLocationName: "fridgeA1", Entity: e1}
		sl2 := StoreLocation{StoreLocationName: "fridgeB1", Entity: e1}
		sl3 := StoreLocation{StoreLocationName: "fridgeA2", Entity: e2}
		sl4 := StoreLocation{StoreLocationName: "fridgeB2", Entity: e2}
		sl5 := StoreLocation{StoreLocationName: "fridgeA3", Entity: e3}
		sl6 := StoreLocation{StoreLocationName: "fridgeB3", Entity: e3}

		db.CreateStoreLocation(sl1)
		db.CreateStoreLocation(sl2)
		db.CreateStoreLocation(sl3)
		db.CreateStoreLocation(sl4)
		db.CreateStoreLocation(sl5)
		db.CreateStoreLocation(sl6)

		m1.Entities = []Entity{e1}
		m2.Entities = []Entity{e2}
		m3.Entities = []Entity{e3}
		m1.Permissions = []Permission{Permission{PermissionPermName: "all", PermissionItemName: "all", PermissionEntityID: e1.EntityID}}
		m2.Permissions = []Permission{Permission{PermissionPermName: "all", PermissionItemName: "all", PermissionEntityID: e2.EntityID}}
		m3.Permissions = []Permission{Permission{PermissionPermName: "all", PermissionItemName: "all", PermissionEntityID: e3.EntityID}}

		db.UpdatePerson(m1)
		db.UpdatePerson(m2)
		db.UpdatePerson(m3)

		p0 := Person{PersonEmail: "user@super.com", Permissions: []Permission{Permission{PermissionPermName: "all", PermissionItemName: "all", PermissionEntityID: -1}}}
		p1 := Person{PersonEmail: "john@lab-one.com", Entities: []Entity{e1}}
		p2 := Person{PersonEmail: "mickey@lab-one.com", Entities: []Entity{e1}}
		p3 := Person{PersonEmail: "donald@lab-one.com", Entities: []Entity{e1}}
		p4 := Person{PersonEmail: "tom@lab-two.com", Entities: []Entity{e2}}
		p5 := Person{PersonEmail: "mike@lab-two.com", Entities: []Entity{e2}}
		p6 := Person{PersonEmail: "ralf@lab-two.com", Entities: []Entity{e2}}
		p7 := Person{PersonEmail: "john@lab-three.com", Entities: []Entity{e3}}
		p8 := Person{PersonEmail: "rob@lab-three.com", Entities: []Entity{e3}}
		p9 := Person{PersonEmail: "harrison@lab-three.com", Entities: []Entity{e3}}
		p10 := Person{PersonEmail: "alone@no-entity.com"}

		db.CreatePerson(p0)
		db.CreatePerson(p1)
		db.CreatePerson(p2)
		db.CreatePerson(p3)
		db.CreatePerson(p4)
		db.CreatePerson(p5)
		db.CreatePerson(p6)
		db.CreatePerson(p7)
		db.CreatePerson(p8)
		db.CreatePerson(p9)
		db.CreatePerson(p10)
	}
	return nil
}

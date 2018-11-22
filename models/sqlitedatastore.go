package models

import (
	"fmt"

	"bufio"
	"database/sql"
	"encoding/csv"
	"os"
	"strings"
	"time"

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
	return &SQLiteDataStore{db}, nil
}

// InsertSamples insert sample values in the database
func (db *SQLiteDataStore) InsertSamples() error {
	var (
		c   int
		err error
	)
	_ = db.Get(&c, `SELECT count(*) FROM person`)
	if c == 0 {
		// inserting sample values
		// FIXME: remove this before release
		scas, _ := os.Open("sample_cas.txt")
		sname, _ := os.Open("sample_name.txt")
		sempiricalformula, _ := os.Open("sample_empiricalformula.txt")
		defer scas.Close()
		defer sname.Close()
		defer sempiricalformula.Close()

		scanner := bufio.NewScanner(scas)
		scanner.Split(bufio.ScanLines)
		log.Debug("- creating sample cas")
		for scanner.Scan() {
			if _, err = db.Exec(`INSERT OR IGNORE INTO casnumber ("casnumber_label") VALUES ("` + scanner.Text() + `");`); err != nil {
				return err
			}
		}

		scanner = bufio.NewScanner(sname)
		scanner.Split(bufio.ScanLines)
		log.Debug("- creating sample names")
		i := 0
		for scanner.Scan() && i < 50 {
			if _, err = db.Exec(`INSERT OR IGNORE INTO name ("name_label") VALUES ("` + scanner.Text() + `");`); err != nil {
				return err
			}
			i++
		}

		scanner = bufio.NewScanner(sempiricalformula)
		scanner.Split(bufio.ScanLines)
		log.Debug("- creating sample empirical formulas")
		i = 0
		for scanner.Scan() && i < 50 {
			if _, err = db.Exec(`INSERT OR IGNORE INTO empiricalformula ("empiricalformula_label") VALUES ("` + scanner.Text() + `");`); err != nil {
				return err
			}
			i++
		}

		m1 := Person{PersonEmail: "manager@lab-one.com"}
		m2 := Person{PersonEmail: "manager@lab-two.com"}
		m3 := Person{PersonEmail: "manager@lab-three.com"}
		m4 := Person{PersonEmail: "delphine.pitrat@ens-lyon.fr"}

		log.Debug("- creating 4 sample managers")
		_, m1.PersonID = db.CreatePerson(m1)
		_, m2.PersonID = db.CreatePerson(m2)
		_, m3.PersonID = db.CreatePerson(m3)
		_, m4.PersonID = db.CreatePerson(m4)

		e1 := Entity{EntityName: "lab one", EntityDescription: "the lab one", Managers: []Person{m1}}
		e2 := Entity{EntityName: "lab two", EntityDescription: "the lab two", Managers: []Person{m2}}
		e3 := Entity{EntityName: "lab three", EntityDescription: "the lab three", Managers: []Person{m3}}
		e4 := Entity{EntityName: "laboratoire de chimie", EntityDescription: "laboratoire de chimie de l'ENS de Lyon", Managers: []Person{m4}}

		log.Debug("- creating 4 sample entities")
		_, e1.EntityID = db.CreateEntity(e1)
		_, e2.EntityID = db.CreateEntity(e2)
		_, e3.EntityID = db.CreateEntity(e3)
		_, e4.EntityID = db.CreateEntity(e4)

		sl1 := StoreLocation{StoreLocationColor: sql.NullString{Valid: true, String: "rgb(255, 38, 38)"}, StoreLocationName: sql.NullString{Valid: true, String: "fridgeE1-A"}, Entity: e1, StoreLocationCanStore: sql.NullBool{Valid: true, Bool: true}}
		sl2 := StoreLocation{StoreLocationColor: sql.NullString{Valid: true, String: "rgb(255, 129, 129)"}, StoreLocationName: sql.NullString{Valid: true, String: "fridgeE1-B"}, Entity: e1, StoreLocationCanStore: sql.NullBool{Valid: true, Bool: true}}
		sl3 := StoreLocation{StoreLocationColor: sql.NullString{Valid: true, String: "rgb(33, 185, 102)"}, StoreLocationName: sql.NullString{Valid: true, String: "fridgeE2-A"}, Entity: e2, StoreLocationCanStore: sql.NullBool{Valid: true, Bool: true}}
		sl4 := StoreLocation{StoreLocationColor: sql.NullString{Valid: true, String: "rgb(99, 232, 159)"}, StoreLocationName: sql.NullString{Valid: true, String: "fridgeE2-B"}, Entity: e2, StoreLocationCanStore: sql.NullBool{Valid: true, Bool: true}}
		sl5 := StoreLocation{StoreLocationColor: sql.NullString{Valid: true, String: "rgb(32, 103, 208)"}, StoreLocationName: sql.NullString{Valid: true, String: "fridgeE3-A"}, Entity: e3, StoreLocationCanStore: sql.NullBool{Valid: true, Bool: true}}
		sl6 := StoreLocation{StoreLocationColor: sql.NullString{Valid: true, String: "rgb(255, 38, 38)"}, StoreLocationName: sql.NullString{Valid: true, String: "roomE3-B"}, Entity: e3, StoreLocationCanStore: sql.NullBool{Valid: true, Bool: false}}

		log.Debug("- creating 5 sample storelocations")
		db.CreateStoreLocation(sl1)
		db.CreateStoreLocation(sl2)
		db.CreateStoreLocation(sl3)
		db.CreateStoreLocation(sl4)
		db.CreateStoreLocation(sl5)

		log.Debug("- creating laboratoire de chimie sample storelocations")
		var lastid int
		slch1 := StoreLocation{
			StoreLocationColor:    sql.NullString{Valid: true, String: "rgb(0, 139, 139)"},
			StoreLocationName:     sql.NullString{Valid: true, String: "[P]M6"},
			Entity:                e4,
			StoreLocationCanStore: sql.NullBool{Valid: true, Bool: false},
		}
		_, lastid = db.CreateStoreLocation(slch1)
		slch1.StoreLocationID = sql.NullInt64{Valid: true, Int64: int64(lastid)}
		slch2 := StoreLocation{
			StoreLocationColor:    sql.NullString{Valid: true, String: "rgb(0, 206, 209)"},
			StoreLocationName:     sql.NullString{Valid: true, String: "[I]Inflammable"},
			Entity:                e4,
			StoreLocation:         &slch1,
			StoreLocationCanStore: sql.NullBool{Valid: true, Bool: true},
		}
		slch3 := StoreLocation{
			StoreLocationColor:    sql.NullString{Valid: true, String: "rgb(32, 178, 170)"},
			StoreLocationName:     sql.NullString{Valid: true, String: "Labo central"},
			Entity:                e4,
			StoreLocation:         &slch1,
			StoreLocationCanStore: sql.NullBool{Valid: true, Bool: false},
		}
		_, lastid = db.CreateStoreLocation(slch3)
		slch3.StoreLocationID = sql.NullInt64{Valid: true, Int64: int64(lastid)}
		slch4 := StoreLocation{
			StoreLocationColor:    sql.NullString{Valid: true, String: "rgb(72, 209, 204)"},
			StoreLocationName:     sql.NullString{Valid: true, String: "[A]Acides"},
			Entity:                e4,
			StoreLocation:         &slch3,
			StoreLocationCanStore: sql.NullBool{Valid: true, Bool: true},
		}
		slch5 := StoreLocation{
			StoreLocationColor:    sql.NullString{Valid: true, String: "rgb(72, 209, 204)"},
			StoreLocationName:     sql.NullString{Valid: true, String: "[C]Congélateur"},
			Entity:                e4,
			StoreLocation:         &slch3,
			StoreLocationCanStore: sql.NullBool{Valid: true, Bool: true},
		}
		slch6 := StoreLocation{
			StoreLocationColor:    sql.NullString{Valid: true, String: "rgb(72, 209, 204)"},
			StoreLocationName:     sql.NullString{Valid: true, String: "[D]Dessicateur"},
			Entity:                e4,
			StoreLocation:         &slch3,
			StoreLocationCanStore: sql.NullBool{Valid: true, Bool: true},
		}
		slch7 := StoreLocation{
			StoreLocationColor:    sql.NullString{Valid: true, String: "rgb(72, 209, 204)"},
			StoreLocationName:     sql.NullString{Valid: true, String: "[F]Frigo"},
			Entity:                e4,
			StoreLocation:         &slch3,
			StoreLocationCanStore: sql.NullBool{Valid: true, Bool: true},
		}
		slch8 := StoreLocation{
			StoreLocationColor:    sql.NullString{Valid: true, String: "rgb(72, 209, 204)"},
			StoreLocationName:     sql.NullString{Valid: true, String: "[P]Placard"},
			Entity:                e4,
			StoreLocation:         &slch3,
			StoreLocationCanStore: sql.NullBool{Valid: true, Bool: true},
		}
		slch9 := StoreLocation{
			StoreLocationColor:    sql.NullString{Valid: true, String: "rgb(72, 209, 204)"},
			StoreLocationName:     sql.NullString{Valid: true, String: "[S]Placard sels et solides"},
			Entity:                e4,
			StoreLocation:         &slch3,
			StoreLocationCanStore: sql.NullBool{Valid: true, Bool: true},
		}
		slch10 := StoreLocation{
			StoreLocationColor:    sql.NullString{Valid: true, String: "rgb(255, 0, 255)"},
			StoreLocationName:     sql.NullString{Valid: true, String: "[P]M6.072"},
			Entity:                e4,
			StoreLocation:         &slch1,
			StoreLocationCanStore: sql.NullBool{Valid: true, Bool: true},
		}
		slch11 := StoreLocation{
			StoreLocationColor:    sql.NullString{Valid: true, String: "rgb(255, 0, 255)"},
			StoreLocationName:     sql.NullString{Valid: true, String: "[P]M6.121"},
			Entity:                e4,
			StoreLocation:         &slch1,
			StoreLocationCanStore: sql.NullBool{Valid: true, Bool: true},
		}
		slch12 := StoreLocation{
			StoreLocationColor:    sql.NullString{Valid: true, String: "rgb(255, 0, 255)"},
			StoreLocationName:     sql.NullString{Valid: true, String: "[P]M6.156"},
			Entity:                e4,
			StoreLocation:         &slch1,
			StoreLocationCanStore: sql.NullBool{Valid: true, Bool: true},
		}
		slch13 := StoreLocation{
			StoreLocationColor:    sql.NullString{Valid: true, String: "rgb(139, 0, 139)"},
			StoreLocationName:     sql.NullString{Valid: true, String: "Soute - local déchets"},
			Entity:                e4,
			StoreLocationCanStore: sql.NullBool{Valid: true, Bool: true},
		}
		slch14 := StoreLocation{
			StoreLocationColor:    sql.NullString{Valid: true, String: "rgb(139, 0, 139)"},
			StoreLocationName:     sql.NullString{Valid: true, String: "[T]Frigo CMR/Toxiques"},
			Entity:                e4,
			StoreLocationCanStore: sql.NullBool{Valid: true, Bool: true},
		}

		db.CreateStoreLocation(slch2)
		db.CreateStoreLocation(slch4)
		db.CreateStoreLocation(slch5)
		db.CreateStoreLocation(slch6)
		db.CreateStoreLocation(slch7)
		db.CreateStoreLocation(slch8)
		db.CreateStoreLocation(slch9)
		db.CreateStoreLocation(slch10)
		db.CreateStoreLocation(slch11)
		db.CreateStoreLocation(slch12)
		db.CreateStoreLocation(slch13)
		db.CreateStoreLocation(slch14)

		m1.Entities = []Entity{e1}
		m2.Entities = []Entity{e2}
		m3.Entities = []Entity{e3}
		m1.Permissions = []Permission{Permission{PermissionPermName: "all", PermissionItemName: "all", PermissionEntityID: e1.EntityID}}
		m2.Permissions = []Permission{Permission{PermissionPermName: "all", PermissionItemName: "all", PermissionEntityID: e2.EntityID}}
		m3.Permissions = []Permission{Permission{PermissionPermName: "all", PermissionItemName: "all", PermissionEntityID: e3.EntityID}}

		log.Debug("- updating the 3 managers")
		db.UpdatePerson(m1)
		db.UpdatePerson(m2)
		db.UpdatePerson(m3)

		p0 := Person{PersonEmail: "user@super.com", Permissions: []Permission{Permission{PermissionPermName: "all", PermissionItemName: "all", PermissionEntityID: -1}}}
		p1 := Person{PersonEmail: "john@lab-one.com", Entities: []Entity{e1}, Permissions: []Permission{Permission{PermissionPermName: "r", PermissionItemName: "products", PermissionEntityID: -1}}}
		p2 := Person{PersonEmail: "mickey@lab-one.com", Entities: []Entity{e1}}
		p3 := Person{PersonEmail: "donald@lab-one.com", Entities: []Entity{e1}}
		p4 := Person{PersonEmail: "tom@lab-two.com", Entities: []Entity{e2}}
		p5 := Person{PersonEmail: "mike@lab-two.com", Entities: []Entity{e2}}
		p6 := Person{PersonEmail: "ralf@lab-two.com", Entities: []Entity{e2}}
		p7 := Person{PersonEmail: "john@lab-three.com", Entities: []Entity{e3}}
		p8 := Person{PersonEmail: "rob@lab-three.com", Entities: []Entity{e3}}
		p9 := Person{PersonEmail: "harrison@lab-three.com", Entities: []Entity{e3}}
		p10 := Person{PersonEmail: "alone@no-entity.com"}

		log.Debug("- creating 11 sample users")
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

		// inserting sample products
		// attention: values are wrongs, just for devel purposes
		log.Debug("- creating sample products")
		for i := 1; i <= 50; i++ {
			ins := fmt.Sprintf("(\"spec%d\", \"%d\", \"%d\", 1, \"%d\")", i, i, i, i)
			if _, err = db.Exec(`INSERT INTO product ("product_specificity", "casnumber", "name", "person", "empiricalformula") VALUES ` + ins + `;`); err != nil {
				return err
			}
			if _, err = db.Exec(`INSERT INTO productsymbols ("productsymbols_product_id", "productsymbols_symbol_id") VALUES 
			(?, ?), (?, ?), (?, ?), (?, ?);`, i, (i%9)+1, i, ((i+1)%9)+1, i, ((i+2)%9)+1, i, ((i+3)%9)+1); err != nil {
				return err
			}
		}

		// inserting sample storages
		// attention: values are wrongs, just for devel purposes
		log.Debug("- creating sample storages")
		for i := 1; i <= 300; i++ {
			comment := fmt.Sprintf("(\"comment%d\", \"%d\", \"%d\")", i, i, i)
			datetime := time.Now()
			person := i%10 + 1
			product := i%19 + 1
			storelocation := i%18 + 1
			unit := i%6 + 1
			quantity := i
			if storelocation == int(sl6.StoreLocationID.Int64) {
				storelocation = int(sl6.StoreLocationID.Int64) + 1
			}
			if _, err = db.Exec(`INSERT INTO storage ("storage_creationdate", "storage_modificationdate", "storage_comment", "person", "product", "storelocation", "storage_quantity", "unit") VALUES (?,?,?,?,?,?,?,?);`, datetime, datetime, comment, person, product, storelocation, quantity, unit); err != nil {
				return err
			}
		}

	}
	log.Debug("done")
	return nil
}

// CreateDatabase creates the database tables
func (db *SQLiteDataStore) CreateDatabase() error {
	var (
		err     error
		c       int
		r       *csv.Reader
		records [][]string
	)

	// activate the foreign keys feature
	log.Debug("enabling foreign keys")
	if _, err = db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return err
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
		storelocation_color string,
		storelocation_canstore boolean default 0,
		storelocation_fullpath string,
		entity integer NOT NULL,
		storelocation integer,
		FOREIGN KEY(storelocation) references storelocation(storelocation_id),
		FOREIGN KEY(entity) references entity(entity_id));
	CREATE TABLE IF NOT EXISTS supplier (
		supplier_id integer PRIMARY KEY,
		supplier_label string NOT NULL);
	CREATE TABLE IF NOT EXISTS unit (
		unit_id integer PRIMARY KEY,
		unit_label string NOT NULL,
		unit_multiplier integer NOT NULL default 1,
		unit integer,
		FOREIGN KEY(unit) references unit(unit_id));
	CREATE TABLE IF NOT EXISTS storage (
		storage_id integer PRIMARY KEY,
		storage_creationdate datetime NOT NULL,
		storage_modificationdate datetime NOT NULL,
		storage_entrydate datetime,
		storage_exitdate datetime,
		storage_openingdate datetime,
		storage_expirationdate datetime,
		storage_quantity integer,
		storage_barecode string,
		storage_comment string,
		storage_reference string,
		storage_batchnumber string,
		storage_todestroy boolean default 0,
		storage_archive boolean default 0,
		person integer NOT NULL,
		product integer NOT NULL,
		storelocation integer NOT NULL,
		unit integer,
		supplier integer,
		storage integer,
		FOREIGN KEY(storage) references storage(storage_id),
		FOREIGN KEY(unit) references unit(unit_id),
		FOREIGN KEY(supplier) references supplier(supplier_id),
		FOREIGN KEY(person) references person(person_id),
		FOREIGN KEY(product) references product(product_id),
		FOREIGN KEY(storelocation) references storelocation(storelocation_id));
	CREATE TABLE IF NOT EXISTS permission (
		permission_id integer PRIMARY KEY,
		person integer NOT NULL,
		permission_perm_name string NOT NULL,
		permission_item_name string NOT NULL,
		permission_entity_id integer,
		FOREIGN KEY(person) references person(person_id));
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
	-- products symbols
	CREATE TABLE IF NOT EXISTS symbol (
		symbol_id integer PRIMARY KEY,
		symbol_label string NOT NULL,
		symbol_image string);
	-- products names
	CREATE TABLE IF NOT EXISTS name (
		name_id integer PRIMARY KEY,
		name_label string NOT NULL UNIQUE);
	-- products cas numbers
	CREATE TABLE IF NOT EXISTS casnumber (
		casnumber_id integer PRIMARY KEY,
		casnumber_label string NOT NULL UNIQUE);
	-- products ce numbers
	CREATE TABLE IF NOT EXISTS cenumber (
		cenumber_id integer PRIMARY KEY,
		cenumber_label string NOT NULL UNIQUE);
	-- products empirical formulas
	CREATE TABLE IF NOT EXISTS empiricalformula (
		empiricalformula_id integer PRIMARY KEY,
		empiricalformula_label string NOT NULL UNIQUE);
	-- products physical states
	CREATE TABLE IF NOT EXISTS physicalstate (
		physicalstate_id integer PRIMARY KEY,
		physicalstate_label string NOT NULL UNIQUE);
	-- products signal words
	CREATE TABLE IF NOT EXISTS signalword (
		signalword_id integer PRIMARY KEY,
		signalword_label string NOT NULL UNIQUE);
	-- products classes of compound
	CREATE TABLE IF NOT EXISTS classofcompound (
		classofcompound_id integer PRIMARY KEY,
		classofcompound_label string NOT NULL UNIQUE);
	-- products hazard statements
	CREATE TABLE IF NOT EXISTS hazardstatement (
		hazardstatement_id integer PRIMARY KEY,
		hazardstatement_label string NOT NULL,
		hazardstatement_reference string NOT NULL);
	-- products precautionary statements
	CREATE TABLE IF NOT EXISTS precautionarystatement (
		precautionarystatement_id integer PRIMARY KEY,
		precautionarystatement_label string NOT NULL,
		precautionarystatement_reference string NOT NULL);
	-- products
	CREATE TABLE IF NOT EXISTS product (
		product_id integer PRIMARY KEY,
		product_specificity string,
		product_msds string,
		product_restricted boolean default 0,
		product_radioactive boolean default 0,
		product_linearformula string,
		product_threedformula string,
		product_disposalcomment string,
		product_remark string,
		casnumber integer,
		cenumber integer,
		person integer NOT NULL,
		empiricalformula integer NOT NULL,
		physicalstate integer,
		signalword integer,
		classofcompound integer,
		name integer NOT NULL,
		FOREIGN KEY(casnumber) references casnumber(casnumber_id),
		FOREIGN KEY(cenumber) references cenumber(cenumber_id),
		FOREIGN KEY(person) references person(person_id),
		FOREIGN KEY(empiricalformula) references empiricalformula(empiricalformula_id),
		FOREIGN KEY(physicalstate) references physicalstate(physicalstate_id),
		FOREIGN KEY(signalword) references signalword(signalword_id),
		FOREIGN KEY(classofcompound) references classofcompound(classofcompound_id),
		FOREIGN KEY(name) references name(name_id));
	CREATE TABLE IF NOT EXISTS productsymbols (
		productsymbols_product_id integer NOT NULL,
		productsymbols_symbol_id integer NOT NULL,
		PRIMARY KEY(productsymbols_product_id, productsymbols_symbol_id),
		FOREIGN KEY(productsymbols_product_id) references product(product_id),
		FOREIGN KEY(productsymbols_symbol_id) references symbol(symbol_id));
	CREATE TABLE IF NOT EXISTS productsynonyms (
		productsynonyms_product_id integer NOT NULL,
		productsynonyms_name_id integer NOT NULL,
		PRIMARY KEY(productsynonyms_product_id, productsynonyms_name_id),
		FOREIGN KEY(productsynonyms_product_id) references product(product_id),
		FOREIGN KEY(productsynonyms_name_id) references name(name_id));
	CREATE TABLE IF NOT EXISTS producthazardstatements (
		producthazardstatements_product_id integer NOT NULL,
		producthazardstatements_hazardstatement_id integer NOT NULL,
		PRIMARY KEY(producthazardstatements_product_id, producthazardstatements_hazardstatement_id),
		FOREIGN KEY(producthazardstatements_product_id) references product(product_id),
		FOREIGN KEY(producthazardstatements_hazardstatement_id) references hazardstatement(hazardstatement_id));	
	CREATE TABLE IF NOT EXISTS productprecautionarystatements (
		productprecautionarystatements_product_id integer NOT NULL,
		productprecautionarystatements_precautionarystatement_id integer NOT NULL,
		PRIMARY KEY(productprecautionarystatements_product_id, productprecautionarystatements_precautionarystatement_id),
		FOREIGN KEY(productprecautionarystatements_product_id) references product(product_id),
		FOREIGN KEY(productprecautionarystatements_precautionarystatement_id) references precautionarystatement(precautionarystatement_id));
	CREATE TABLE IF NOT EXISTS bookmark (
		bookmark_id integer PRIMARY KEY,
		person integer NOT NULL,
		product integer NOT NULL,
		FOREIGN KEY(person) references person(person_id),
		FOREIGN KEY(product) references product(product_id));	
		`

	// values definition
	inssymbol := `INSERT INTO symbol (symbol_label, symbol_image) VALUES 
	("sgh01", "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAInSURBVFiFzdi9b45RGMfxz11KNalS0napgVaE0HhLxEtpVEIQxKAbiUhMJCoisTyriYTBLCZGsRCrP8JkkLB4iVlyDL1EU9refZ7raZ3kJPfrub75/a7zWpVSpJSqaoBSGintlVJarzQKJWojo81sqDS4XKUSlcu3LwkuFyoRLh8qCa49UAlw6VDoRUercOlK4T7WtapcNlQHHmfYmgm1CVO4MOPZBow31V4SVDde4i22xrMeXEfVVK5myB4gh/AQVwPqSljbEe/XYVfd9lOgIvAt3AnAS1iBczgV168wVTdOClSAPcMwzmIg4EbRP+u7behZKF6r9q3BTTzFC1wLO49iD/owHioex2nswGpsnC9uU1BYhUE8R8EH3As1DuIYtmAnDsT9SZwPJScxMp8o9RKRtQHSFUk8jBHcxpPIr95QqC+svIxHGKiVDrM4VqpRSik/qqoaxTecwSe8CUWO4Dve4W6o9xFf8Bl9VVV1RgfoDLXfl1J+LhR0bp+nVRjGZoxhLw7jRNhzIwAKXmMCD/AVDVxsRq3ayY/1GEK/6RF+u+k5cTAUGJoxVk1ionaPnjf568HtD6h9GJunY3RjN7qahfobrEYP9Xv0brUuaoCt+VO7oeYGaydcS5N4u+BSlj3ZcKkLxSy4tiytW4Vr62ak2SBLsn1bbLAl3fDWDbosRwQLBV/WQ5W5IP6LY6h/w6VA5YAl2jez1lrBLlhKaaiqP9cJ5Rf+De5Q3HyidwAAAABJRU5ErkJggg=="),
	("sgh02", "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAIvSURBVFiFzdgxbI5BGMDx36uNJsJApFtFIh2QEIkBtYiYJFKsrDaJqYOofhKRMFhsFhNNMBgkTAaD0KUiBomN1EpYBHWGnvi8vu9r7/2eti65vHfv3T3P/57nuXvfuyqlJCRVVQuk1AqRl1LqP9NKpJxbETKjocLgYi0VaLl49wXBxUIFwsVDBcGFQWEAg1FwYZbCGMajLBfmPkzgUZRbw2IKFzGPrRFw/bpvD/bn8jUkXM719f3A9eu+k3iXA/92Bnub2yYx1NgDfbrvXIYZx8dcThjBExxvPOmGltqLIzmuEt63QSVczc+z/2whSw2ThpbajS+4UgOq59O4gYFSuGaByWb8zKvwN8RXXKiBPc7PLaWx3ARqY37O1CBe5/cvO1huVy+ZnfSX7y9MYxRTNeX32lZj+/sXWNfVnV3g1tT/aJeQ5vAGp3L9eXbjTFv7NzzM9VncSSnNF2lp4MqjNYvcxwEcy+0HcQg32/q8Kndl+YrcgM9Z4YdsrZ21PtvxHT9yv1vNgr8cbiIrnMUmbKu177PwVZjLgKPNt4sCOKzF0ww32aF9CA+yxSZKoTqDlVnucI6lMxhpg76OuxhrKr8oIENyXx/xxQKTE/hUkIdLJ1tlRd3TwtF/KtcuSalVVdUwdvQe+Fd6ljhfl9NzRKT5I8cvq/B+xi3vzFfk+FaqbEUPvEtVuipXBIspX9VLlW4Q/8U1VGe4EKgYsED3tefBgt271y7dUlV/ygHpF8bRglXiwx7BAAAAAElFTkSuQmCC"),
	("sgh03", "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAJhSURBVFiFzdjNq41RFMfxz/FSXK4iUV4GUjJRMpG8lQnJUF0mTAy8xkApimOEMDA2oSgD5R9gcE0obxO6JkooBroJhRttg7Nu94lzzz37OdvLrt16zrPXs9d3/9Ze+5zzNFJKirRGowlSahaZL6XUe6eZSNGbJeYsDVUMrqxSBZUrn75CcMWgsARTSsEVUwq7sLWUcmXS1wK7iOvd+pcDa5++qVgU1/fxKq4n9wrXa/qW4Bb6MIKE2diPmb3A1VYq7MqAuRQ2YTdeY6CXtNZVaj4uYG0FaLR/D3sc0+vC1d3okwJgsAI0iB+Vz5dxe1TdXLg6UHPCvg2AT2E34VobBaflzD8+2AQPRYqu4kUEPh1KzcKOuPck7CMcQF92nOyVtCquqsg8PI5C2IyHWBFjn8NuzM5Mdu7ZGcGO4k1U5EgF9CNO4QuG4t6x7ALLrhY2RLB9uBMAJ7Ea63A+CuMVlobvidzqzz9fmFtR5jvWtPHZHj4Xww5MNO+vHJNktpTSezxAP26klO618bkZah4JRe/mxslOZSiyLZQ43MHnTPicy1Wr1uavBH6Hsx3Gr+ADZudC1TouKoFv4CX624wtwDBO1oH6HSwDDsvxTetrZ2Hl/jKtg3UYs+pAtQfLg1uldcqPhH2qVanPsL4u1Phg3ayIPdiLg3hu7IAdwqEY24vFtRZdew/wtQLTqW+ptYc7gnWYLPbS8i76jFyo7sBqTFri+T86eS/P/dmV/5W/b7nB/uof3m6D/pNXBBMF/6cvVcaD+C9eQ7WHKwJVBqxg+qp9SvYvy3YtpaZGY+y6QPsJlPiFVobY9AkAAAAASUVORK5CYII="),
	("sgh04", "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAFtSURBVFiFzdixLgRBGADgb4UoBYXKA3gAHYU3EKWH0FKuAk8gGm+hvU4oJAqlaDVqiUSCUbjEJe5ys7v/7NnkT66Ynf/bmZ3bf6ZKKQm5qqoGKdUh/aWUugd1Ig2jjugzGhWGix2pwJGLn74gXCwqEBePCsKVQQXgyqE63lcW1eH+8qiW/fSDatFff6iG/faLatB//6jMPKEoLGEX53jEHTaxj0sMcvN1QmEBWzjGLT6QxsQX7nGUO3KNUdjAAa7wOgEyGjdYazqt+auEQzxnQEbjBdtt3rn5nCq3qqp1nGJuStM3XGMwjIc0fKrGV9YKYQXv/o7Ip58X/AQ7WIxaofnLl7Mh5gkX2MNyqb+NrEYjuNXOkMx8jRr3hRoP6wPX6pNUGtfpI14KF1L2RONCC8UoXJHSuiuu6GakbZJetm9Nk/W64c1NOpMjgmnJZ3qoMgnxL46hxuNCUDGwwOkbjawKNqParFXV7++A6xtDLLIHRMAuWAAAAABJRU5ErkJggg=="),
	("sgh05", "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAI9SURBVFiFzdhPiI1RGMfxz2FmgYVo1CyYhVJISuPPDmFGNwtW2EzKn4TEELKaa2NNdiyUlWxtUGYptiytrRSNIk3GsbinceO+7vvOPe8dp56677k9v+d7nnPOczonxBhlaSE0QYzNLHoxxt6NZiQma+bQzA2VDS5vpjJmLv/0ZYLLC5URLj9UJrh6oDLA1QfVo1+9UD341w+1QJ2exLAG+3AKq9DAGUxgJ1YsFG6gy9k3lb5uibEZQhjHNnxFxEe8w1ucwxfcxxw24GAIYTkCBmOMTSFIulNCKDxbO4N1gEq/P6SsfMdw6hvFDB5hGS5iLYZSRr/hE6bRAikDV2X6sC4B7MZQ13XSAlufII5iBGOllknFNTWGq6W3PJO4ngYygRu4WSoJJTO1N62TXWnqtmK0BNgR3E4DmsQdjJRJRtlMNXAX53EBz7G6BNgAjmEH9hT6dIhfagtjCV5gC67hUgmoA3iNl3iDe7iSQAe7wRWWixBCAyfauj5jo9bOPBRCeNz239kY40yb72lsx5MEth9PY4zvQwgr8aMo7nwrTCWXtWpVGRv+I1vH8SxlaRxLsblKIQ9J6K/aFXiITV1H1mrTMcbZ9o4QwmGc1CrGr/BTa83N4kGMca5T3PmattAjI4uVKhf9hqtUYPsFV6YS9OJcF9S/weqAq6CXVSynTi2iOfxrFe/Fr96R9+X6VjVYXy+8ZYMuyhNBt+CL+qhSBPFfPEN1hssClQcs4/S1W/GFt0r7fVcsvMBWbb8AgnCJLinP5ycAAAAASUVORK5CYII="),
	("sgh06", "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAK6SURBVFiFzdhPiFdVFMDxzxsnpTIFNQpJZjEoVEwQyoRYRP82JUQySAYuZITEDBeCuvPnInXhzGAaRTMtDNyMBElUi9m4ECVGI6FVMOuilUG4COO0eBfmj2/m997vd0e7cOG9y7vnfN8555577i0iQpZWFC0Q0coiLyK677SCSL2VQ2ZuqGxweS2V0XL53ZcJLi9URrj8UJnglgcqA1zXUHgc7+EsvsQxbOwWrluoAXyDHXgkjfXhC7zbDVw3UIO4gFWVgjmIjzqF6xTqhWSVNxcVzAq0cKgTuI4CF6cxgr4lhTOG850siMZQSeEr+B29NcBONZVfDba0pXrmPH+Ll7EX7+OxNP56GtuBy1iXxot5P9Lu5xtADWICn6AHn2M3Av9gc/puOo19jXNp7FlM4kJtfQ3cdx39eAPncS4p3IPXUrD34xkcwUZ8hW24hJXJtUO1Flhtn/M9Pk6JdBc+wx8pLfRgP+5gLV7ET7iRIA+keXtxpk74FFEu6RNl2eikRSrQoijeSop/w3Cyzq94CbeTpQaSu99JMXUN2zGOX7ATP0TE3xUK5nH0VEFUtYiYiojpiPgrIkZxFE/jJp7HFA4rF8R3CXImIvZFxPWIuBsRk5VQiyhsnpVL1w0pk+wVPKeMvQ8T7EVlDI5hSzt51a4sFd1nyoUuLYriOJ5Srrgn8LPSnffQi1X4E5twFa9iBquxRunaf/FjREzNEVytt9YKKS00jrH0vi8J2zDnm0ct2AmU8TehjL+VuIUPmqWL9nCFsqxZj5G27pmdt1a52W/FcGcJtj3cIEbxZF2wNO9tfIrVtdNTk4DM0rvaxJcLLkvZkxsua6GYC25ZSutu4Zb1MNKpkgdyfGuq7IEeeOsqfShXBO2UP9RLlcUg/hfXUNVwWaDygGV039zeW10+NmwRLUUx+5yh/QdzLVcJBJ5ddQAAAABJRU5ErkJggg=="),
	("sgh07", "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAF/SURBVFiFzdghT8NAGMbx/y1g0HMLQWDxJDgCkk+AQhE8BlkDGL4JwczyAXBkCY4vQAhBgiI7xJpwK1fW9nm6ccklNe/7/u7eW9c2xBixjBAKAGIsLPlijPqEIkIsZ+HI6UbZcN6dMu6cv30mnBdlxPlRJlw/KAOuP5QYJycHBsARMAYmwCUwVHHyioELIFbmvbpzchuAmwzsSW2rfDaA0wzsTs4rrwz2M7BruRPyWYBRBnYin10FVcIC8FGB7TWJ/fPXrqAS3KQCGzaOr6kro0rYbYJ6bxufq+/5w4WrBPbQJUfVMZAfgWfjuea6+zC1cht4Y7Zjx55WGg5/idsEDhyoeZiIA7aAXQfqN6wjDjgHvoApi+76Det0CsrAXpi/j+0oqHpYSxzwmKCmwEjNK213AjsEXoFP4MyyWPUsJLh1YMOBagbrkNQR32tyJa7flS/l9a1tsaW+8DYtupJPBIuKr/SjSh3iX3yGyuMsKA/M2L50rpmeNgtC+Lk2jG/Rx4o589viKwAAAABJRU5ErkJggg=="),
	("sgh08", "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAKJSURBVFiFzdixix1FGADw38QDUbQQqyNewiEYCxNN4ECjkuIaFYLNQQpTCBLkzD8gJsVLOqPGJiBYeGB1SeFhYXEhmCIpTgIpRDAEosWhckXU6AknCUyKN49bn+/tze6OnAsDuzA789v5vvnevg0xRkWOEHogxl6R8WKM3Ru9SEytV2LM0qhiuLIrVXDlyoevEK4sqiCuGApTeBo7SuA6o7AX3yKmdh3PdsV1RU3hzwpq0NbxZBdc1/AtjEAN2lddcq41KsG+qYH9jYfb4lqjEuxWDSxissl4ebCMQbBSg/qp7bjjYbk3c7oGttjpoduiEuz5GthbndKkLaqCu5Qgd/FdOv8FD3XaWF1QCfZawqzhWDo/mnv/2FLUEfUA3kuYL/FyOn87GzZm/rbhm8JH+LmSU6t4pHL9A87gmTa4Hf9+px19hBAmQgjzIYQVnMIGHqt0+QKvVq6ncRAzIYSlEMLZEMKu3PmyQ4kTldWI+lX/MD7GZ5jHJH7ESbyJs/qbYnDP79idF8rc7Tu6mK7hHexJfQLewE7jfxWO5Cd/zvblypiJ/sBzQ32vjun7T1hWudjqCdiH20OTrOJz7MJTqd8reAk3RqCu4fGsCDWqLezHr2mSdcxgFss4nkK4iPcxh8v6myTiazyandONCx8H8Ck+wM3KapxP+Ta43tDfqXM4hweb1MzGha+CuzMUpoWU+MPh+yR3g+XD6nEvpMQfAJaxNIT6EKEpKg9Wj3vR5jv/b/r1a4B6t81KNYPV4w7hL3yfcu+e6ivPf/pnZGvcLC7iAl7vimoOq8dN44kSqHawnEm35RPBVpNv60eVcYj/xWeo0bgiqDKwguGrtgkljhh7Qtg8L3DcB497IINNg8B2AAAAAElFTkSuQmCC"),
	("sgh09", "image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAI8SURBVFiFzdi7a1VBEMDhb0EjYhFfaCEICopBi4haWUiqqGgQW1OKaBpNIQEVvREi2lhaWCiC/4CSImgv2IgQELSSdKKlMUKEtTgr3jxuPI+9iQMD5+y9u/NjdmbO7IYYoywSQgvE2MqyXoyxudKKxKStHGvmhsoGl9dTGT2Xf/syweWFygiXHyoTXHegMsDVgprjLk5hS4wRenLD1fYUhjGans9jF/rRlwOuMhS2Yz924F0CO4EBjOTyXK0YwUVcwAcEHMJr9OaKudqLYBA/cQx7cT+NB1xrCld5cvLOGA7gqSIRricv7sED9Df1XB1PbcVOXMZnzOA7pvEcg3UTqjNY1QDlKmLSsRX+14NNlXamLlQyeDJBfenw+xmcxQ3cxunSCVYpINmYatUfHUlg3xaN9+MmjqR5uxPkaNltrZbChcFYUi8tmjuOqbIxt25pT7uifMThtvcBPMS84kvwKY2fizE+hhDCBkXGzuBZaUsN4qtPUVQj3uOJIhmGFSXlOO4loEmsr5KhtYIfvRjCUTxS1K8reKnIwAkLt/UHhqqUjUblIkFuxi18TRCv8KsNah5v8UbqRqqVi5JwOKhI/ym8wKzlg3+xzmGyrJ2QjC2U4ox4J72NS2fFEMI27CsdwEtlNsY43Wn9BVIlILNoo494t+CytD254bI2irngutJaN4Xr6mGkrpFVOb5VNbaqB96yRtfkiuBfxtf0UqUTxH9xDbU8XBaoPGAZt69dq3awy0uMLSH8fc4gvwFyuYuihNiCxwAAAABJRU5ErkJggg==");`
	insphysicalstate := `INSERT INTO physicalstate (physicalstate_label) VALUES ("gaz"), ("liquid"), ("solid")`
	inssignalword := `INSERT INTO signalword (signalword_label) VALUES ("danger"), ("warning")`
	insunit := `INSERT INTO unit (unit_label, unit_multiplier, unit) VALUES 
	("l", 1, NULL), ("ml", 0.001, 1), ("µl", 0.00001, 1),
	("kg", 1000, 2), ("g", 1, NULL), ("mg", 0.001, 2), ("µg", 0.00001, 2),
	("m", 1, NULL), ("dm", 0.1, 3), ("cm", 0.01, 3)`

	// tables creation
	log.Debug("creating tables")
	if _, err = db.Exec(schema); err != nil {
		return err
	}

	// symbols
	if err = db.Get(&c, `SELECT count(*) FROM symbol`); err != nil {
		return err
	}
	if c == 0 {
		if _, err = db.Exec(inssymbol); err != nil {
			return err
		}
	}

	// signal words
	if err = db.Get(&c, `SELECT count(*) FROM signalword`); err != nil {
		return err
	}
	if c == 0 {
		if _, err = db.Exec(inssignalword); err != nil {
			return err
		}
	}

	// signal units
	if err = db.Get(&c, `SELECT count(*) FROM unit`); err != nil {
		return err
	}
	if c == 0 {
		if _, err = db.Exec(insunit); err != nil {
			return err
		}
	}

	// physical states
	if err = db.Get(&c, `SELECT count(*) FROM physicalstate`); err != nil {
		return err
	}
	if c == 0 {
		if _, err = db.Exec(insphysicalstate); err != nil {
			return err
		}
	}

	// hazard statements
	if err = db.Get(&c, `SELECT count(*) FROM hazardstatement`); err != nil {
		return err
	}
	if c == 0 {
		r = csv.NewReader(strings.NewReader(HAZARDSTATEMENT))
		r.Comma = '\t'
		if records, err = r.ReadAll(); err != nil {
			return err
		}
		for _, record := range records {
			if _, err = db.Exec(`INSERT INTO hazardstatement (hazardstatement_label, hazardstatement_reference) VALUES (?, ?)`, record[0], record[1]); err != nil {
				return err
			}
		}
	}

	// precautionary statements
	if err = db.Get(&c, `SELECT count(*) FROM precautionarystatement`); err != nil {
		return err
	}
	if c == 0 {
		r = csv.NewReader(strings.NewReader(PRECAUTIONARYSTATEMENT))
		r.Comma = '\t'
		if records, err = r.ReadAll(); err != nil {
			return err
		}
		for _, record := range records {
			if _, err = db.Exec(`INSERT INTO precautionarystatement (precautionarystatement_label, precautionarystatement_reference) VALUES (?, ?)`, record[0], record[1]); err != nil {
				return err
			}
		}
	}
	return nil
}

package datastores

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
)

type CasbinJSON struct {
	PType string
	V0    string
	V1    string
	V2    string
	V3    string
}

// SQLiteDataStore implements the Datastore interface
// to store data in SQLite3.
type SQLiteDataStore struct {
	*sqlx.DB
}

var regex = func(re, s string) bool {
	var (
		m   bool
		err error
	)

	if m, err = regexp.MatchString(re, s); err != nil {
		return false
	}

	return m
}

func init() {
	sql.Register("sqlite3_with_go_func",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				return conn.RegisterFunc("regexp", regex, true)
			},
		})
}

// GetWelcomeAnnounce returns the welcome announce.
func (db *SQLiteDataStore) GetWelcomeAnnounce() (models.WelcomeAnnounce, error) {
	var (
		wa   models.WelcomeAnnounce
		sqlr string
		err  error
	)

	sqlr = `SELECT welcome_announce.welcome_announce_id, welcome_announce.welcome_announce_text
	FROM welcome_announce LIMIT 1`
	if err = db.Get(&wa, sqlr); err != nil {
		return models.WelcomeAnnounce{}, err
	}

	logger.Log.WithFields(logrus.Fields{"wa": wa}).Debug("GetWelcomeAnnounce")

	return wa, nil
}

// UpdateWelcomeAnnounce updates the main page announce.
func (db *SQLiteDataStore) UpdateWelcomeAnnounce(w models.WelcomeAnnounce) error {
	var (
		sqlr string
		tx   *sqlx.Tx
		err  error
	)

	// beginning new transaction
	if tx, err = db.Beginx(); err != nil {
		return err
	}

	// updating person
	sqlr = `UPDATE welcome_announce SET welcome_announce_text = ?
	WHERE welcome_announce_id = (SELECT welcome_announce_id FROM welcome_announce LIMIT 1)`
	if _, err = tx.Exec(sqlr, w.WelcomeAnnounceText); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}

	// committing changes
	if err = tx.Commit(); err != nil {
		if errr := tx.Rollback(); errr != nil {
			return errr
		}
	}

	return nil
}

func (db *SQLiteDataStore) GetDB() *sqlx.DB {
	return db.DB
}

// NewSQLiteDBstore returns a database connection to the given dataSourceName
// ie. a path to the sqlite database file.
func NewSQLiteDBstore(dataSourceName string) (*SQLiteDataStore, error) {
	var (
		db  *sqlx.DB
		err error
	)

	logger.Log.WithFields(logrus.Fields{"dbdriver": "sqlite3", "dataSourceName": dataSourceName}).Debug("NewDBstore")

	if db, err = sqlx.Connect("sqlite3_with_go_func", dataSourceName+"?_journal=wal&_fk=1"); err != nil {
		return &SQLiteDataStore{}, err
	}

	return &SQLiteDataStore{db}, nil
}

// ToCasbinJSONAdapter returns a JSON as a slice of bytes
// following the format: https://github.com/casbin/json-adapter#policy-json
func (db *SQLiteDataStore) ToCasbinJSONAdapter() ([]byte, error) {
	var (
		ps   []models.Permission
		err  error
		res  []byte
		sqlr string
	)

	sqlr = `SELECT person AS "person.person_id", permission_perm_name, permission_item_name, permission_entity_id 
	FROM permission`
	if err = db.Select(&ps, sqlr); err != nil {
		return nil, err
	}

	js := make([]CasbinJSON, 0, len(ps))

	for _, p := range ps {
		js = append(js, CasbinJSON{
			PType: "p",
			V0:    strconv.Itoa(p.Person.PersonID),
			V1:    p.PermissionPermName,
			V2:    p.PermissionItemName,
			V3:    strconv.Itoa(p.PermissionEntityID),
		})
	}

	if res, err = json.Marshal(js); err != nil {
		return nil, err
	}

	return res, nil
}

// CreateDatabase creates the database tables.
func (db *SQLiteDataStore) CreateDatabase() error {
	var (
		err         error
		c           int
		userVersion int
		r           *csv.Reader
		records     [][]string
	)

	// tables creation
	var count_entity_table int
	if err = db.Get(&count_entity_table, `SELECT count(*) FROM sqlite_schema WHERE type='table' AND name='entity'`); err != nil {
		logger.Log.Fatal(err)
	}

	if count_entity_table == 0 {
		logger.Log.Info("  creating sqlite tables")
		if _, err = db.Exec(schema); err != nil {
			return err
		}
	}

	// shema migration
	logger.Log.Info("  running migrations")
	if err = db.Get(&userVersion, `PRAGMA user_version`); err != nil {
		return err
	}

	logger.Log.Info(fmt.Sprintf("  user_version:%d", userVersion))

	nextVersion := userVersion + 1
	for _, version := range versionToMigration[userVersion:] {
		logger.Log.Infof("  upgrading version to %d ", nextVersion)

		if _, err = db.Exec(version); err != nil {
			return err
		}
		nextVersion++
		logger.Log.Infof("  done")
	}

	// welcome announce
	if err = db.Get(&c, `SELECT count(*) FROM welcome_announce`); err != nil {
		return err
	}

	if c == 0 {
		logger.Log.Info("  inserting welcome announce")

		if _, err = db.Exec(inswelcomeannounce); err != nil {
			return err
		}
	}

	// symbols
	if err = db.Get(&c, `SELECT count(*) FROM symbol`); err != nil {
		return err
	}

	if c == 0 {
		logger.Log.Info("  inserting symbols")

		if _, err = db.Exec(inssymbol); err != nil {
			return err
		}
	}

	// physical states
	if err = db.Get(&c, `SELECT count(*) FROM physical_state`); err != nil {
		return err
	}

	if c == 0 {
		logger.Log.Info("  inserting physical states")

		if _, err = db.Exec(insphysicalstate); err != nil {
			return err
		}
	}

	// signal words
	if err = db.Get(&c, `SELECT count(*) FROM signal_word`); err != nil {
		return err
	}

	if c == 0 {
		logger.Log.Info("  inserting signal words")

		if _, err = db.Exec(inssignalword); err != nil {
			return err
		}
	}

	// cas numbers
	if err = db.Get(&c, `SELECT count(*) FROM cas_number`); err != nil {
		return err
	}

	if c == 0 {
		logger.Log.Info("  inserting CMRs")

		r = csv.NewReader(strings.NewReader(CMR_CAS))
		r.Comma = ','

		if records, err = r.ReadAll(); err != nil {
			return err
		}

		for _, record := range records {
			if _, err = db.Exec(`INSERT INTO cas_number (cas_number_label, cas_number_cmr) VALUES (?, ?)`,
				record[0],
				record[1]); err != nil {
				return err
			}
		}
	}

	// tags
	if err = db.Get(&c, `SELECT count(*) FROM tag`); err != nil {
		return err
	}

	if c == 0 {
		logger.Log.Info("  inserting tags")

		r = csv.NewReader(strings.NewReader(TAG))
		r.Comma = ','

		if records, err = r.ReadAll(); err != nil {
			return err
		}

		for _, record := range records {
			if _, err = db.Exec(`INSERT INTO tag (tag_label) VALUES (?)`, record[0]); err != nil {
				return err
			}
		}
	}

	// categories
	if err = db.Get(&c, `SELECT count(*) FROM category`); err != nil {
		return err
	}

	if c == 0 {
		logger.Log.Info("  inserting categories")

		r = csv.NewReader(strings.NewReader(CATEGORY))
		r.Comma = ';'

		if records, err = r.ReadAll(); err != nil {
			return err
		}

		for _, record := range records {
			if _, err = db.Exec(`INSERT INTO category (category_label) VALUES (?)`, record[0]); err != nil {
				return err
			}
		}
	}

	// suppliers
	if err = db.Get(&c, `SELECT count(*) FROM supplier`); err != nil {
		return err
	}

	if c == 0 {
		logger.Log.Info("  inserting suppliers")

		r = csv.NewReader(strings.NewReader(SUPPLIER))
		r.Comma = ','

		if records, err = r.ReadAll(); err != nil {
			return err
		}

		for _, record := range records {
			if _, err = db.Exec(`INSERT INTO supplier (supplier_label) VALUES (?)`, record[0]); err != nil {
				return err
			}
		}
	}

	// producers
	if err = db.Get(&c, `SELECT count(*) FROM producer`); err != nil {
		return err
	}

	if c == 0 {
		logger.Log.Info("  inserting producers")

		r = csv.NewReader(strings.NewReader(PRODUCER))
		r.Comma = ','

		if records, err = r.ReadAll(); err != nil {
			return err
		}

		for _, record := range records {
			if _, err = db.Exec(`INSERT INTO producer (producer_label) VALUES (?)`, record[0]); err != nil {
				return err
			}
		}
	}

	// classes of compounds
	if err = db.Get(&c, `SELECT count(*) FROM class_of_compound`); err != nil {
		return err
	}

	if c == 0 {
		logger.Log.Info("  inserting classes of compounds")

		r = csv.NewReader(strings.NewReader(CLASSOFCOMPOUND))
		r.Comma = ','

		if records, err = r.ReadAll(); err != nil {
			return err
		}

		for _, record := range records {
			if _, err = db.Exec(`INSERT INTO class_of_compound (class_of_compound_label) VALUES (?)`, record[0]); err != nil {
				return err
			}
		}
	}

	// hazard statements
	if err = db.Get(&c, `SELECT count(*) FROM hazard_statement`); err != nil {
		return err
	}

	if c == 0 {
		logger.Log.Info("  inserting hazard statements")

		r = csv.NewReader(strings.NewReader(HAZARDSTATEMENT))
		r.Comma = '\t'

		if records, err = r.ReadAll(); err != nil {
			return err
		}

		for _, record := range records {
			if _, err = db.Exec(`INSERT INTO hazard_statement (
				hazard_statement_label,
				hazard_statement_reference,
				hazard_statement_cmr) VALUES (?, ?, ?)`,
				record[0], record[1], record[2]); err != nil {
				return err
			}
		}
	}

	// precautionary statements
	if err = db.Get(&c, `SELECT count(*) FROM precautionary_statement`); err != nil {
		return err
	}

	if c == 0 {
		logger.Log.Info("  inserting precautionary statements")

		r = csv.NewReader(strings.NewReader(PRECAUTIONARYSTATEMENT))
		r.Comma = '\t'

		if records, err = r.ReadAll(); err != nil {
			return err
		}

		for _, record := range records {
			if _, err = db.Exec(`INSERT INTO precautionary_statement (
				precautionary_statement_label,
				precautionary_statement_reference) VALUES (?, ?)`,
				record[0], record[1]); err != nil {
				return err
			}
		}
	}

	// inserting default admin
	var admin *models.Person

	if err = db.Get(&c, `SELECT count(*) FROM person`); err != nil {
		return err
	}

	if c == 0 {
		logger.Log.Info("  inserting admin user")

		admin = &models.Person{
			PersonEmail: "admin@chimitheque.fr",
			Permissions: []*models.Permission{{PermissionPermName: "all", PermissionItemName: "all", PermissionEntityID: -1}},
		}

		if _, err = db.CreatePerson(*admin); err != nil {
			return err
		}

	}

	// tables creation
	logger.Log.Info("  vacuuming database")

	if _, err = db.Exec("VACUUM;"); err != nil {
		return err
	}

	return nil
}

func (db *SQLiteDataStore) Maintenance() {
	var (
		err  error
		sqlr string
		tx   *sql.Tx
	)

	//
	// Cleaning up casnumber labels duplicates.
	//
	if tx, err = db.Begin(); err != nil {
		logger.Log.Error(err)

		return
	}

	var casNumbers []models.CasNumber

	sqlr = `SELECT cas_number_id, cas_number_label FROM cas_number;`

	if err = db.Select(&casNumbers, sqlr); err != nil {
		logger.Log.Error(err)

		return
	}

	for _, casNumber := range casNumbers {
		// if strings.HasPrefix(casNumber.CasNumberLabel.String, " ") || strings.HasSuffix(casNumber.CasNumberLabel.String, " ") {
		if strings.HasPrefix(*casNumber.CasNumberLabel, " ") || strings.HasSuffix(*casNumber.CasNumberLabel, " ") {
			logger.Log.Infof("casnumber %s contains spaces", *casNumber.CasNumberLabel)

			trimmedLabel := strings.Trim(*casNumber.CasNumberLabel, " ")

			// Checking if the trimmed label already exists.
			var existCasNumber models.CasNumber

			sqlr = `SELECT cas_number_id, cas_number_label FROM cas_number WHERE cas_number_label=?;`

			if err = db.Get(&existCasNumber, sqlr, trimmedLabel); err != nil {
				switch err {
				case sql.ErrNoRows:
					// Just fixing the label.
					logger.Log.Info("  - fixing it")

					sqlr = `UPDATE cas_number SET cas_number_label=? WHERE cas_number_id=?;`

					if _, err = tx.Exec(sqlr, trimmedLabel, casNumber.CasNumberID); err != nil {
						logger.Log.Error(err)

						if errr := tx.Rollback(); errr != nil {
							logger.Log.Error(err)

							return
						}

						return
					}

					continue
				default:
					logger.Log.Error(err)

					return
				}
			}

			// Updating products with the found casnumber.
			// logger.Log.Infof("  - correct cas number found, replacing it: %d -> %d",
			// 	existCasNumber.CasNumberID.Int64,
			// 	casNumber.CasNumberID.Int64)

			sqlr = `UPDATE product SET cas_number=? WHERE cas_number=?;`

			if _, err = tx.Exec(sqlr, existCasNumber.CasNumberID, casNumber.CasNumberID); err != nil {
				logger.Log.Error(err)

				if errr := tx.Rollback(); errr != nil {
					logger.Log.Error(err)

					return
				}

				return
			}

			// Deleting the wrong cas number.
			logger.Log.Info("  - deleting it")

			sqlr = `DELETE FROM cas_number WHERE cas_number_id=?;`

			if _, err = tx.Exec(sqlr, casNumber.CasNumberID); err != nil {
				logger.Log.Error(err)

				if errr := tx.Rollback(); errr != nil {
					logger.Log.Error(err)

					return
				}

				return
			}
		}
	}

	if err = tx.Commit(); err != nil {
		logger.Log.Error(err)

		if errr := tx.Rollback(); errr != nil {
			logger.Log.Error(errr)

			return
		}
	}
}

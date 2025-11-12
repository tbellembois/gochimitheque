package datastores

import (
	"database/sql"
	"encoding/json"
	"regexp"
	"strconv"

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

	sqlr = `SELECT person AS "person.person_id", permission_name, permission_item, permission_entity
	FROM permission`
	if err = db.Select(&ps, sqlr); err != nil {
		return nil, err
	}

	js := make([]CasbinJSON, 0, len(ps))

	for _, p := range ps {
		js = append(js, CasbinJSON{
			PType: "p",
			V0:    strconv.Itoa(int(*p.Person.PersonID)),
			V1:    p.PermissionName,
			V2:    p.PermissionItem,
			V3:    strconv.Itoa(int(p.PermissionEntity)),
		})
	}

	if res, err = json.Marshal(js); err != nil {
		return nil, err
	}

	return res, nil
}

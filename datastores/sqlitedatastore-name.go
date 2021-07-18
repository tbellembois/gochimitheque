package datastores

import (
	"database/sql"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"

	. "github.com/tbellembois/gochimitheque/models"
)

// GetNames return the names matching the search criteria
func (db *SQLiteDataStore) GetNames(p Dbselectparam) ([]Name, int, error) {

	var (
		err                              error
		names                            []Name
		count                            int
		exactSearch, countSql, selectSql string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("GetNames")

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")
	nameTable := goqu.T("name")

	// Join, where.
	joinClause := dialect.From(
		nameTable,
	).Where(
		goqu.I("name_label").Like(p.GetSearch()),
	)

	if countSql, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("name_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}
	if selectSql, selectArgs, err = joinClause.Select(
		goqu.I("name_id"),
		goqu.I("name_label"),
	).Order(
		goqu.L("INSTR(name_label, \"?\")", exactSearch).Asc(),
		goqu.C("name_label").Asc(),
	).Limit(uint(p.GetLimit())).Offset(uint(p.GetOffset())).ToSQL(); err != nil {
		return nil, 0, err
	}

	// select
	if err = db.Select(&names, selectSql, selectArgs...); err != nil {
		return nil, 0, err
	}
	// count
	if err = db.Get(&count, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for formula matching exactly the search
	sQuery := dialect.From(nameTable).Where(
		goqu.I("name_label").Eq(exactSearch),
	).Select(
		"name_id",
		"name_label",
	)

	var (
		sqlr string
		args []interface{}
		name Name
	)
	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return nil, 0, err
	}

	if err = db.Get(&name, sqlr, args...); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	for i, e := range names {
		if e.NameID == name.NameID {
			names[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"names": names}).Debug("GetNames")

	return names, count, nil

}

// GetName return the formula matching the given id
func (db *SQLiteDataStore) GetName(id int) (Name, error) {

	var (
		err  error
		sqlr string
		args []interface{}
		name Name
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetName")

	dialect := goqu.Dialect("sqlite3")
	nameTable := goqu.T("name")

	sQuery := dialect.From(nameTable).Where(
		goqu.I("name_id").Eq(id),
	).Select(
		goqu.I("name_id"),
		goqu.I("name_label"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return Name{}, err
	}

	if err = db.Get(&name, sqlr, args...); err != nil {
		return Name{}, err
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "name": name}).Debug("GetName")

	return name, nil

}

// GetNameByLabel return the name matching the given name
func (db *SQLiteDataStore) GetNameByLabel(label string) (Name, error) {

	var (
		err  error
		sqlr string
		args []interface{}
		name Name
	)
	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetNameByLabel")

	dialect := goqu.Dialect("sqlite3")
	nameTable := goqu.T("name")

	sQuery := dialect.From(nameTable).Where(
		goqu.I("name_label").Eq(label),
	).Select(
		goqu.I("name_id"),
		goqu.I("name_label"),
	).Order(goqu.I("name_label").Asc())

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return Name{}, err
	}

	if err = db.Get(&name, sqlr, args...); err != nil {
		return Name{}, err
	}

	logger.Log.WithFields(logrus.Fields{"label": label, "name": name}).Debug("GetNameByLabel")

	return name, nil

}

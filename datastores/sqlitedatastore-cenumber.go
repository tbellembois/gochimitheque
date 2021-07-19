package datastores

import (
	"database/sql"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"

	. "github.com/tbellembois/gochimitheque/models"
)

// GetCeNumbers return the cas numbers matching the search criteria
func (db *SQLiteDataStore) GetCeNumbers(p SelectFilter) ([]CeNumber, int, error) {

	var (
		err                              error
		ceNumbers                        []CeNumber
		count                            int
		exactSearch, countSql, selectSql string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("GetProductsCeNumbers")

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")
	cenumberTable := goqu.T("cenumber")

	// Join, where.
	joinClause := dialect.From(
		cenumberTable,
	).Where(
		goqu.I("cenumber_label").Like(p.GetSearch()),
	)

	if countSql, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("cesnumber_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}
	if selectSql, selectArgs, err = joinClause.Select(
		goqu.I("cenumber_id"),
		goqu.I("cenumber_label"),
	).Order(
		goqu.L("INSTR(cenumber_label, \"?\")", exactSearch).Asc(),
		goqu.C("cenumber_label").Asc(),
	).Limit(uint(p.GetLimit())).Offset(uint(p.GetOffset())).ToSQL(); err != nil {
		return nil, 0, err
	}

	// select
	if err = db.Select(&ceNumbers, selectSql, selectArgs...); err != nil {
		return nil, 0, err
	}
	// count
	if err = db.Get(&count, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for formula matching exactly the search
	sQuery := dialect.From(cenumberTable).Where(
		goqu.I("cenumber_label").Eq(exactSearch),
	).Select(
		"cenumber_id",
		"cenumber_label",
	)

	var (
		sqlr string
		args []interface{}
		cen  CeNumber
	)
	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return nil, 0, err
	}

	if err = db.Get(&cen, sqlr, args...); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	for i, c := range ceNumbers {
		if c.CeNumberID == cen.CeNumberID {
			ceNumbers[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"ceNumbers": ceNumbers}).Debug("GetProductsCeNumbers")

	return ceNumbers, count, nil

}

// GetGetCeNumber return the cas numbers matching the given id
func (db *SQLiteDataStore) GetCeNumber(id int) (CeNumber, error) {

	var (
		err  error
		sqlr string
		args []interface{}
		ce   CeNumber
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetCeNumber")

	dialect := goqu.Dialect("sqlite3")
	cenumberTable := goqu.T("cenumber")

	sQuery := dialect.From(cenumberTable).Where(
		goqu.I("cenumber_id").Eq(id),
	).Select(
		goqu.I("cenumber_id"),
		goqu.I("cenumber_label"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return CeNumber{}, err
	}

	if err = db.Get(&ce, sqlr, args...); err != nil {
		return CeNumber{}, err
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "ce": ce}).Debug("GetCeNumber")

	return ce, nil

}

// GetCeNumberByLabel return the ce numbers matching the given ce number
func (db *SQLiteDataStore) GetCeNumberByLabel(label string) (CeNumber, error) {

	var (
		err  error
		sqlr string
		args []interface{}
		ce   CeNumber
	)
	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetProductsCeNumberByLabel")

	dialect := goqu.Dialect("sqlite3")
	cenumberTable := goqu.T("cenumber")

	sQuery := dialect.From(cenumberTable).Where(
		goqu.I("cenumber_label").Eq(label),
	).Select(
		goqu.I("cenumber_id"),
		goqu.I("cenumber_label"),
	).Order(goqu.I("cenumber_labe").Asc())

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return CeNumber{}, err
	}

	if err = db.Get(&ce, sqlr, args...); err != nil {
		return CeNumber{}, err
	}

	logger.Log.WithFields(logrus.Fields{"label": label, "ce": ce}).Debug("GetProductsCeNumberByLabel")

	return ce, nil

}

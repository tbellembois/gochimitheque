package datastores

import (
	"database/sql"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"

	. "github.com/tbellembois/gochimitheque/models"
)

// GetCasNumbers return the cas numbers matching the search criteria
func (db *SQLiteDataStore) GetCasNumbers(p SelectFilter) ([]CasNumber, int, error) {

	var (
		err                              error
		casNumbers                       []CasNumber
		count                            int
		exactSearch, countSql, selectSql string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("GetProductsCasNumbers")

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")
	casnumberTable := goqu.T("casnumber")

	// Join, where.
	joinClause := dialect.From(
		casnumberTable,
	).Where(
		goqu.I("casnumber_label").Like(p.GetSearch()),
	)

	if countSql, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("casnumber_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}
	if selectSql, selectArgs, err = joinClause.Select(
		goqu.I("casnumber_id"),
		goqu.I("casnumber_label"),
	).Order(
		goqu.L("INSTR(casnumber_label, \"?\")", exactSearch).Asc(),
		goqu.C("casnumber_label").Asc(),
	).Limit(uint(p.GetLimit())).Offset(uint(p.GetOffset())).ToSQL(); err != nil {
		return nil, 0, err
	}

	// select
	if err = db.Select(&casNumbers, selectSql, selectArgs...); err != nil {
		return nil, 0, err
	}
	// count
	if err = db.Get(&count, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for formula matching exactly the search
	sQuery := dialect.From(casnumberTable).Where(
		goqu.I("casnumber_label").Eq(exactSearch),
	).Select(
		"casnumber_id",
		"casnumber_label",
	)

	var (
		sqlr string
		args []interface{}
		casn CasNumber
	)
	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return nil, 0, err
	}

	if err = db.Get(&casn, sqlr, args...); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	for i, c := range casNumbers {
		if c.CasNumberID == casn.CasNumberID {
			casNumbers[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"casNumbers": casNumbers}).Debug("GetProductsCasNumbers")

	return casNumbers, count, nil

}

// GetCasNumber return the cas numbers matching the given id
func (db *SQLiteDataStore) GetCasNumber(id int) (CasNumber, error) {

	var (
		err  error
		sqlr string
		args []interface{}
		cas  CasNumber
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetCasNumber")

	dialect := goqu.Dialect("sqlite3")
	casnumberTable := goqu.T("casnumber")

	sQuery := dialect.From(casnumberTable).Where(
		goqu.I("casnumber_id").Eq(id),
	).Select(
		goqu.I("casnumber_id"),
		goqu.I("casnumber_label"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return CasNumber{}, err
	}

	if err = db.Get(&cas, sqlr, args...); err != nil {
		return CasNumber{}, err
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "cas": cas}).Debug("GetCasNumber")

	return cas, nil

}

// GetCasNumberByLabel return the cas numbers matching the given cas number
func (db *SQLiteDataStore) GetCasNumberByLabel(label string) (CasNumber, error) {

	var (
		err  error
		sqlr string
		args []interface{}
		cas  CasNumber
	)
	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetCasNumberByLabel")

	dialect := goqu.Dialect("sqlite3")
	casnumberTable := goqu.T("casnumber")

	sQuery := dialect.From(casnumberTable).Where(
		goqu.I("casnumber_label").Eq(label),
	).Select(
		goqu.I("casnumber_id"),
		goqu.I("casnumber_label"),
	).Order(goqu.I("casnumber_label").Asc())

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return CasNumber{}, err
	}

	if err = db.Get(&cas, sqlr, args...); err != nil {
		return CasNumber{}, err
	}

	logger.Log.WithFields(logrus.Fields{"label": label, "ca": cas}).Debug("GetCasNumberByLabel")

	return cas, nil

}

package datastores

import (
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"

	. "github.com/tbellembois/gochimitheque/models"
)

// GetPrecautionaryStatements return the precautionary statements matching the search criteria
func (db *SQLiteDataStore) GetPrecautionaryStatements(p Dbselectparam) ([]PrecautionaryStatement, int, error) {

	var (
		err                              error
		precautionaryStatements          []PrecautionaryStatement
		count                            int
		exactSearch, countSql, selectSql string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("GetPrecautionaryStatements")

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")
	precautionarystatementTable := goqu.T("precautionarystatement")

	// Join, where.
	joinClause := dialect.From(
		precautionarystatementTable,
	).Where(
		goqu.I("precautionarystatement_label").Like(p.GetSearch()),
	)

	if countSql, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("precautionarystatement_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}
	if selectSql, selectArgs, err = joinClause.Select(
		goqu.I("precautionarystatement_id"),
		goqu.I("precautionarystatement_label"),
		goqu.I("precautionarystatement_reference"),
	).Order(
		goqu.L("INSTR(precautionarystatement_label, \"?\")", exactSearch).Asc(),
		goqu.C("precautionarystatement_label").Asc(),
	).Limit(uint(p.GetLimit())).Offset(uint(p.GetOffset())).ToSQL(); err != nil {
		return nil, 0, err
	}

	// select
	if err = db.Select(&precautionaryStatements, selectSql, selectArgs...); err != nil {
		return nil, 0, err
	}
	// count
	if err = db.Get(&count, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	logger.Log.WithFields(logrus.Fields{"precautionaryStatements": precautionaryStatements}).Debug("GetPrecautionaryStatements")

	return precautionaryStatements, count, nil

}

// GetPrecautionaryStatement return the formula matching the given id
func (db *SQLiteDataStore) GetPrecautionaryStatement(id int) (PrecautionaryStatement, error) {

	var (
		err                    error
		sqlr                   string
		args                   []interface{}
		precautionaryStatement PrecautionaryStatement
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetPrecautionaryStatement")

	dialect := goqu.Dialect("sqlite3")
	precautionarystatementTable := goqu.T("precautionarystatement")

	sQuery := dialect.From(precautionarystatementTable).Where(
		goqu.I("precautionarystatement_id").Eq(id),
	).Select(
		goqu.I("precautionarystatement_id"),
		goqu.I("precautionarystatement_label"),
		goqu.I("precautionarystatement_reference"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return PrecautionaryStatement{}, err
	}

	if err = db.Get(&precautionaryStatement, sqlr, args...); err != nil {
		return PrecautionaryStatement{}, err
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "precautionaryStatement": precautionaryStatement}).Debug("GetPrecautionaryStatement")

	return precautionaryStatement, nil

}

// GetPrecautionaryStatementByReference return the empirirical formula matching the given precautionary statement
func (db *SQLiteDataStore) GetPrecautionaryStatementByReference(reference string) (PrecautionaryStatement, error) {

	var (
		err                    error
		sqlr                   string
		args                   []interface{}
		precautionaryStatement PrecautionaryStatement
	)
	logger.Log.WithFields(logrus.Fields{"reference": reference}).Debug("GetPrecautionaryStatementByLabel")

	dialect := goqu.Dialect("sqlite3")
	precautionarystatementTable := goqu.T("precautionarystatement")

	sQuery := dialect.From(precautionarystatementTable).Where(
		goqu.I("precautionarystatement_reference").Eq(reference),
	).Select(
		goqu.I("precautionarystatement_id"),
		goqu.I("precautionarystatement_label"),
		goqu.I("precautionarystatement_reference"),
	).Order(goqu.I("precautionarystatement_label").Asc())

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return PrecautionaryStatement{}, err
	}

	if err = db.Get(&precautionaryStatement, sqlr, args...); err != nil {
		return PrecautionaryStatement{}, err
	}

	logger.Log.WithFields(logrus.Fields{"reference": reference, "precautionaryStatement": precautionaryStatement}).Debug("GetPrecautionaryStatementByLabel")

	return precautionaryStatement, nil

}

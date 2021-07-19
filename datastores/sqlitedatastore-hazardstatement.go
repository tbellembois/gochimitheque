package datastores

import (
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"

	. "github.com/tbellembois/gochimitheque/models"
)

// GetHazardStatements return the hazard statements matching the search criteria
func (db *SQLiteDataStore) GetHazardStatements(p SelectFilter) ([]HazardStatement, int, error) {

	var (
		err                              error
		hazardStatements                 []HazardStatement
		count                            int
		exactSearch, countSql, selectSql string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("GetHazardStatements")

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")
	hazardstatementTable := goqu.T("hazardstatement")

	// Join, where.
	joinClause := dialect.From(
		hazardstatementTable,
	).Where(
		goqu.I("hazardstatement_label").Like(p.GetSearch()),
	)

	if countSql, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("hazardstatement_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}
	if selectSql, selectArgs, err = joinClause.Select(
		goqu.I("hazardstatement_id"),
		goqu.I("hazardstatement_label"),
		goqu.I("hazardstatement_reference"),
	).Order(
		goqu.L("INSTR(hazardstatement_label, \"?\")", exactSearch).Asc(),
		goqu.C("hazardstatement_label").Asc(),
	).Limit(uint(p.GetLimit())).Offset(uint(p.GetOffset())).ToSQL(); err != nil {
		return nil, 0, err
	}

	// select
	if err = db.Select(&hazardStatements, selectSql, selectArgs...); err != nil {
		return nil, 0, err
	}
	// count
	if err = db.Get(&count, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	logger.Log.WithFields(logrus.Fields{"hazardStatements": hazardStatements}).Debug("GetHazardStatements")

	return hazardStatements, count, nil

}

// GetHazardStatement return the formula matching the given id
func (db *SQLiteDataStore) GetHazardStatement(id int) (HazardStatement, error) {

	var (
		err             error
		sqlr            string
		args            []interface{}
		hazardStatement HazardStatement
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetHazardStatement")

	dialect := goqu.Dialect("sqlite3")
	hazardstatementTable := goqu.T("hazardstatement")

	sQuery := dialect.From(hazardstatementTable).Where(
		goqu.I("hazardstatement_id").Eq(id),
	).Select(
		goqu.I("hazardstatement_id"),
		goqu.I("hazardstatement_label"),
		goqu.I("hazardstatement_reference"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return HazardStatement{}, err
	}

	if err = db.Get(&hazardStatement, sqlr, args...); err != nil {
		return HazardStatement{}, err
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "hazardStatement": hazardStatement}).Debug("GetHazardStatement")

	return hazardStatement, nil

}

// GetHazardStatementByReference return the empirirical formula matching the given hazard statement
func (db *SQLiteDataStore) GetHazardStatementByReference(reference string) (HazardStatement, error) {

	var (
		err             error
		sqlr            string
		args            []interface{}
		hazardStatement HazardStatement
	)
	logger.Log.WithFields(logrus.Fields{"reference": reference}).Debug("GetHazardStatementByLabel")

	dialect := goqu.Dialect("sqlite3")
	hazardstatementTable := goqu.T("hazardstatement")

	sQuery := dialect.From(hazardstatementTable).Where(
		goqu.I("hazardstatement_reference").Eq(reference),
	).Select(
		goqu.I("hazardstatement_id"),
		goqu.I("hazardstatement_label"),
		goqu.I("hazardstatement_reference"),
	).Order(goqu.I("hazardstatement_label").Asc())

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return HazardStatement{}, err
	}

	if err = db.Get(&hazardStatement, sqlr, args...); err != nil {
		return HazardStatement{}, err
	}

	logger.Log.WithFields(logrus.Fields{"reference": reference, "hazardStatement": hazardStatement}).Debug("GetHazardStatementByLabel")

	return hazardStatement, nil

}

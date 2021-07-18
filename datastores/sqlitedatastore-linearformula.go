package datastores

import (
	"database/sql"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"

	. "github.com/tbellembois/gochimitheque/models"
)

// GetLinearFormulas return the linear formulas matching the search criteria
func (db *SQLiteDataStore) GetLinearFormulas(p Dbselectparam) ([]LinearFormula, int, error) {

	var (
		err                              error
		linearFormulas                   []LinearFormula
		count                            int
		exactSearch, countSql, selectSql string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("GetLinearFormulas")

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")
	linearformulaTable := goqu.T("linearformula")

	// Join, where.
	joinClause := dialect.From(
		linearformulaTable,
	).Where(
		goqu.I("linearformula_label").Like(p.GetSearch()),
	)

	if countSql, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("linearformula_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}
	if selectSql, selectArgs, err = joinClause.Select(
		goqu.I("linearformula_id"),
		goqu.I("linearformula_label"),
	).Order(
		goqu.L("INSTR(linearformula_label, \"?\")", exactSearch).Asc(),
		goqu.C("linearformula_label").Asc(),
	).Limit(uint(p.GetLimit())).Offset(uint(p.GetOffset())).ToSQL(); err != nil {
		return nil, 0, err
	}

	// select
	if err = db.Select(&linearFormulas, selectSql, selectArgs...); err != nil {
		return nil, 0, err
	}
	// count
	if err = db.Get(&count, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for formula matching exactly the search
	sQuery := dialect.From(linearformulaTable).Where(
		goqu.I("linearformula_label").Eq(exactSearch),
	).Select(
		"linearformula_id",
		"linearformula_label",
	)

	var (
		sqlr string
		args []interface{}
		lf   LinearFormula
	)
	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return nil, 0, err
	}

	if err = db.Get(&lf, sqlr, args...); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	for i, e := range linearFormulas {
		if e.LinearFormulaID == lf.LinearFormulaID {
			linearFormulas[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"linearFormulas": linearFormulas}).Debug("GetLinearFormulas")

	return linearFormulas, count, nil

}

// GetLinearFormula return the formula matching the given id
func (db *SQLiteDataStore) GetLinearFormula(id int) (LinearFormula, error) {

	var (
		err           error
		sqlr          string
		args          []interface{}
		linearFormula LinearFormula
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetLinearFormula")

	dialect := goqu.Dialect("sqlite3")
	linearformulaTable := goqu.T("linearformula")

	sQuery := dialect.From(linearformulaTable).Where(
		goqu.I("linearformula_id").Eq(id),
	).Select(
		goqu.I("linearformula_id"),
		goqu.I("linearformula_label"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return LinearFormula{}, err
	}

	if err = db.Get(&linearFormula, sqlr, args...); err != nil {
		return LinearFormula{}, err
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "linearFormula": linearFormula}).Debug("GetLinearFormula")

	return linearFormula, nil

}

// GetLinearFormulaByLabel return the empirirical formula matching the given linear formula
func (db *SQLiteDataStore) GetLinearFormulaByLabel(label string) (LinearFormula, error) {

	var (
		err           error
		sqlr          string
		args          []interface{}
		linearFormula LinearFormula
	)
	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetLinearFormulaByLabel")

	dialect := goqu.Dialect("sqlite3")
	linearformulaTable := goqu.T("linearformula")

	sQuery := dialect.From(linearformulaTable).Where(
		goqu.I("linearformula_label").Eq(label),
	).Select(
		goqu.I("linearformula_id"),
		goqu.I("linearformula_label"),
	).Order(goqu.I("linearformula_label").Asc())

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return LinearFormula{}, err
	}

	if err = db.Get(&linearFormula, sqlr, args...); err != nil {
		return LinearFormula{}, err
	}

	logger.Log.WithFields(logrus.Fields{"label": label, "linearFormula": linearFormula}).Debug("GetLinearFormulaByLabel")

	return linearFormula, nil

}

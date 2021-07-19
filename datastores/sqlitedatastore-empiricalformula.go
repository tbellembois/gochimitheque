package datastores

import (
	"database/sql"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"

	. "github.com/tbellembois/gochimitheque/models"
)

// GetEmpiricalFormulas return the empirical formulas matching the search criteria
func (db *SQLiteDataStore) GetEmpiricalFormulas(p SelectFilter) ([]EmpiricalFormula, int, error) {

	var (
		err                              error
		empiricalFormulas                []EmpiricalFormula
		count                            int
		exactSearch, countSql, selectSql string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("GetEmpiricalFormulas")

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")
	empiricalformulaTable := goqu.T("empiricalformula")

	// Join, where.
	joinClause := dialect.From(
		empiricalformulaTable,
	).Where(
		goqu.I("empiricalformula_label").Like(p.GetSearch()),
	)

	if countSql, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("empiricalformula_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}
	if selectSql, selectArgs, err = joinClause.Select(
		goqu.I("empiricalformula_id"),
		goqu.I("empiricalformula_label"),
	).Order(
		goqu.L("INSTR(empiricalformula_label, \"?\")", exactSearch).Asc(),
		goqu.C("empiricalformula_label").Asc(),
	).Limit(uint(p.GetLimit())).Offset(uint(p.GetOffset())).ToSQL(); err != nil {
		return nil, 0, err
	}

	// select
	if err = db.Select(&empiricalFormulas, selectSql, selectArgs...); err != nil {
		return nil, 0, err
	}
	// count
	if err = db.Get(&count, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for formula matching exactly the search
	sQuery := dialect.From(empiricalformulaTable).Where(
		goqu.I("empiricalformula_label").Eq(exactSearch),
	).Select(
		"empiricalformula_id",
		"empiricalformula_label",
	)

	var (
		sqlr string
		args []interface{}
		ef   EmpiricalFormula
	)
	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return nil, 0, err
	}

	if err = db.Get(&ef, sqlr, args...); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	for i, e := range empiricalFormulas {
		if e.EmpiricalFormulaID == ef.EmpiricalFormulaID {
			empiricalFormulas[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"empiricalFormulas": empiricalFormulas}).Debug("GetEmpiricalFormulas")

	return empiricalFormulas, count, nil

}

// GetEmpiricalFormula return the formula matching the given id
func (db *SQLiteDataStore) GetEmpiricalFormula(id int) (EmpiricalFormula, error) {

	var (
		err              error
		sqlr             string
		args             []interface{}
		empiricalFormula EmpiricalFormula
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetEmpiricalFormula")

	dialect := goqu.Dialect("sqlite3")
	empiricalformulaTable := goqu.T("empiricalformula")

	sQuery := dialect.From(empiricalformulaTable).Where(
		goqu.I("empiricalformula_id").Eq(id),
	).Select(
		goqu.I("empiricalformula_id"),
		goqu.I("empiricalformula_label"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return EmpiricalFormula{}, err
	}

	if err = db.Get(&empiricalFormula, sqlr, args...); err != nil {
		return EmpiricalFormula{}, err
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "empiricalFormula": empiricalFormula}).Debug("GetEmpiricalFormula")

	return empiricalFormula, nil

}

// GetEmpiricalFormulaByLabel return the empirirical formula matching the given empirical formula
func (db *SQLiteDataStore) GetEmpiricalFormulaByLabel(label string) (EmpiricalFormula, error) {

	var (
		err              error
		sqlr             string
		args             []interface{}
		empiricalFormula EmpiricalFormula
	)
	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetEmpiricalFormulaByLabel")

	dialect := goqu.Dialect("sqlite3")
	empiricalformulaTable := goqu.T("empiricalformula")

	sQuery := dialect.From(empiricalformulaTable).Where(
		goqu.I("empiricalformula_label").Eq(label),
	).Select(
		goqu.I("empiricalformula_id"),
		goqu.I("empiricalformula_label"),
	).Order(goqu.I("empiricalformula_label").Asc())

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return EmpiricalFormula{}, err
	}

	if err = db.Get(&empiricalFormula, sqlr, args...); err != nil {
		return EmpiricalFormula{}, err
	}

	logger.Log.WithFields(logrus.Fields{"label": label, "empiricalFormula": empiricalFormula}).Debug("GetEmpiricalFormulaByLabel")

	return empiricalFormula, nil

}

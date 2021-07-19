package datastores

import (
	"database/sql"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"

	. "github.com/tbellembois/gochimitheque/models"
)

// GetClassesOfCompound return the class of compounds matching the search criteria
func (db *SQLiteDataStore) GetClassesOfCompound(p SelectFilter) ([]ClassOfCompound, int, error) {

	var (
		err                              error
		classOfCompounds                 []ClassOfCompound
		count                            int
		exactSearch, countSql, selectSql string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("GetClassOfCompounds")

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")
	classofcompoundTable := goqu.T("classofcompound")

	// Join, where.
	joinClause := dialect.From(
		classofcompoundTable,
	).Where(
		goqu.I("classofcompound_label").Like(p.GetSearch()),
	)

	if countSql, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("classofcompound_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}
	if selectSql, selectArgs, err = joinClause.Select(
		goqu.I("classofcompound_id"),
		goqu.I("classofcompound_label"),
	).Order(
		goqu.L("INSTR(classofcompound_label, \"?\")", exactSearch).Asc(),
		goqu.C("classofcompound_label").Asc(),
	).Limit(uint(p.GetLimit())).Offset(uint(p.GetOffset())).ToSQL(); err != nil {
		return nil, 0, err
	}

	// select
	if err = db.Select(&classOfCompounds, selectSql, selectArgs...); err != nil {
		return nil, 0, err
	}
	// count
	if err = db.Get(&count, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for formula matching exactly the search
	sQuery := dialect.From(classofcompoundTable).Where(
		goqu.I("classofcompound_label").Eq(exactSearch),
	).Select(
		"classofcompound_id",
		"classofcompound_label",
	)

	var (
		sqlr string
		args []interface{}
		coc  ClassOfCompound
	)
	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return nil, 0, err
	}

	if err = db.Get(&coc, sqlr, args...); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	for i, c := range classOfCompounds {
		if c.ClassOfCompoundID == coc.ClassOfCompoundID {
			classOfCompounds[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"classOfCompounds": classOfCompounds}).Debug("GetClassOfCompounds")

	return classOfCompounds, count, nil

}

// GetClassOfCompound return the formula matching the given id
func (db *SQLiteDataStore) GetClassOfCompound(id int) (ClassOfCompound, error) {

	var (
		err             error
		sqlr            string
		args            []interface{}
		classOfCompound ClassOfCompound
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetClassOfCompound")

	dialect := goqu.Dialect("sqlite3")
	classofcompoundTable := goqu.T("classofcompound")

	sQuery := dialect.From(classofcompoundTable).Where(
		goqu.I("classofcompound_id").Eq(id),
	).Select(
		goqu.I("classofcompound_id"),
		goqu.I("classofcompound_label"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return ClassOfCompound{}, err
	}

	if err = db.Get(&classOfCompound, sqlr, args...); err != nil {
		return ClassOfCompound{}, err
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "classOfCompound": classOfCompound}).Debug("GetClassOfCompound")

	return classOfCompound, nil

}

// GetClassOfCompoundByLabel return the empirirical formula matching the given class of compound
func (db *SQLiteDataStore) GetClassOfCompoundByLabel(label string) (ClassOfCompound, error) {

	var (
		err             error
		sqlr            string
		args            []interface{}
		classOfCompound ClassOfCompound
	)
	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetClassOfCompoundByLabel")

	dialect := goqu.Dialect("sqlite3")
	classofcompoundTable := goqu.T("classofcompound")

	sQuery := dialect.From(classofcompoundTable).Where(
		goqu.I("classofcompound_label").Eq(label),
	).Select(
		goqu.I("classofcompound_id"),
		goqu.I("classofcompound_label"),
	).Order(goqu.I("classofcompound_label").Asc())

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return ClassOfCompound{}, err
	}

	if err = db.Get(&classOfCompound, sqlr, args...); err != nil {
		return ClassOfCompound{}, err
	}

	logger.Log.WithFields(logrus.Fields{"label": label, "classOfCompound": classOfCompound}).Debug("GetClassOfCompoundByLabel")

	return classOfCompound, nil

}

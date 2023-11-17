package datastores

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/zmqclient"
)

func GetByMany[T models.Searchable](searchable T, db *sqlx.DB, filter zmqclient.RequestFilter) (ts []T, count int, err error) {
	var (
		exactSearch, countSQL, selectSQL string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"filter": filter}).Debug("GetByMany")

	exactSearch = filter.Search
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")

	// Join, where.
	joinClause := dialect.From(
		searchable.GetTableName(),
	).Where(
		goqu.I(searchable.GetTextFieldName()).Like(filter.Search),
	)

	if countSQL, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I(searchable.GetIDFieldName()).Distinct()),
	).ToSQL(); err != nil {
		return
	}

	if selectSQL, selectArgs, err = joinClause.Select(
		goqu.I("*"),
	).Order(
		goqu.L(fmt.Sprintf("INSTR(%s, ?)", searchable.GetTextFieldName()), exactSearch).Asc(),
		goqu.C(searchable.GetTextFieldName()).Asc(),
	).Limit(
		uint(filter.Limit),
	).Offset(
		uint(filter.Offset),
	).ToSQL(); err != nil {
		return
	}

	// Select.
	if err = db.Select(&ts, selectSQL, selectArgs...); err != nil {
		return
	}
	// Count.
	if err = db.Get(&count, countSQL, countArgs...); err != nil {
		return
	}

	// Setting the C attribute for formula matching exactly the search.
	sQuery := dialect.From(searchable.GetTableName()).Where(
		goqu.I(searchable.GetTextFieldName()).Eq(exactSearch),
	).Select("*")

	var (
		sqlr string
		args []interface{}
		t    T
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return
	}

	if err = db.Get(&t, sqlr, args...); err != nil && err != sql.ErrNoRows {
		logger.Log.Error(err)
		return
	}

	err = nil

	for i, c := range ts {
		if c.GetID() == t.GetID() {
			ts[i] = (ts[i].SetC(1)).(T)
		}
	}

	logger.Log.WithFields(logrus.Fields{"ts": ts}).Debug("GetByMany")

	return
}

func GetByID[T models.Searchable](searchable T, db *sqlx.DB, id int) (t T, err error) {
	var (
		sqlr string
		args []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"ID": id}).Debug("GetByID")

	dialect := goqu.Dialect("sqlite3")

	sQuery := dialect.From(
		searchable.GetTableName(),
	).Where(
		goqu.I(searchable.GetIDFieldName()).Eq(id),
	).Select(
		goqu.I("*"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return
	}

	if err = db.Get(&t, sqlr, args...); err != nil && err != sql.ErrNoRows {
		logger.Log.Error(err)
		return
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "t": t}).Debug("GetByID")

	return
}

func GetByText[T models.Searchable](searchable T, db *sqlx.DB, text string) (t T, err error) {
	var (
		sqlr string
		args []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"text": text}).Debug("GetByText")

	dialect := goqu.Dialect("sqlite3")

	sQuery := dialect.From(
		searchable.GetTableName(),
	).Where(
		goqu.I(searchable.GetTextFieldName()).Eq(text),
	).Select(
		goqu.I("*"),
	).Order(goqu.I(searchable.GetTextFieldName()).Asc())

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return
	}

	logger.Log.WithFields(logrus.Fields{"sqlr": sqlr, "args": args}).Debug("GetByText")

	if err = db.Get(&t, sqlr, args...); err != nil && err != sql.ErrNoRows {
		logger.Log.Error(err)
		return
	}

	logger.Log.WithFields(logrus.Fields{"text": text, "t": t}).Debug("GetByText")

	return
}

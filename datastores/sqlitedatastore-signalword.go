package datastores

import (
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"

	. "github.com/tbellembois/gochimitheque/models"
)

// GetSignalWords return the signal words matching the search criteria
func (db *SQLiteDataStore) GetSignalWords(p Dbselectparam) ([]SignalWord, int, error) {

	var (
		err                              error
		signalWords                      []SignalWord
		count                            int
		exactSearch, countSql, selectSql string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("GetSignalWords")

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")
	signalwordTable := goqu.T("signalword")

	// Join, where.
	joinClause := dialect.From(
		signalwordTable,
	).Where(
		goqu.I("signalword_label").Like(p.GetSearch()),
	)

	if countSql, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("signalword_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}
	if selectSql, selectArgs, err = joinClause.Select(
		goqu.I("signalword_id"),
		goqu.I("signalword_label"),
	).Order(
		goqu.L("INSTR(signalword_label, \"?\")", exactSearch).Asc(),
		goqu.C("signalword_label").Asc(),
	).Limit(uint(p.GetLimit())).Offset(uint(p.GetOffset())).ToSQL(); err != nil {
		return nil, 0, err
	}

	// select
	if err = db.Select(&signalWords, selectSql, selectArgs...); err != nil {
		return nil, 0, err
	}
	// count
	if err = db.Get(&count, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	logger.Log.WithFields(logrus.Fields{"signalWords": signalWords}).Debug("GetSignalWords")

	return signalWords, count, nil

}

// GetSignalWord return the formula matching the given id
func (db *SQLiteDataStore) GetSignalWord(id int) (SignalWord, error) {

	var (
		err        error
		sqlr       string
		args       []interface{}
		signalWord SignalWord
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetSignalWord")

	dialect := goqu.Dialect("sqlite3")
	signalwordTable := goqu.T("signalword")

	sQuery := dialect.From(signalwordTable).Where(
		goqu.I("signalword_id").Eq(id),
	).Select(
		goqu.I("signalword_id"),
		goqu.I("signalword_label"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return SignalWord{}, err
	}

	if err = db.Get(&signalWord, sqlr, args...); err != nil {
		return SignalWord{}, err
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "signalWord": signalWord}).Debug("GetSignalWord")

	return signalWord, nil

}

// GetSignalWordByLabel return the empirirical formula matching the given signal word
func (db *SQLiteDataStore) GetSignalWordByLabel(label string) (SignalWord, error) {

	var (
		err        error
		sqlr       string
		args       []interface{}
		signalWord SignalWord
	)
	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetSignalWordByLabel")

	dialect := goqu.Dialect("sqlite3")
	signalwordTable := goqu.T("signalword")

	sQuery := dialect.From(signalwordTable).Where(
		goqu.I("signalword_label").Eq(label),
	).Select(
		goqu.I("signalword_id"),
		goqu.I("signalword_label"),
	).Order(goqu.I("signalword_label").Asc())

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return SignalWord{}, err
	}

	if err = db.Get(&signalWord, sqlr, args...); err != nil {
		return SignalWord{}, err
	}

	logger.Log.WithFields(logrus.Fields{"label": label, "signalWord": signalWord}).Debug("GetSignalWordByLabel")

	return signalWord, nil

}

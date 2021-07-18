package datastores

import (
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"

	. "github.com/tbellembois/gochimitheque/models"
)

// GetSymbols return the symbols matching the search criteria
func (db *SQLiteDataStore) GetSymbols(p Dbselectparam) ([]Symbol, int, error) {

	var (
		err                              error
		symbols                          []Symbol
		count                            int
		exactSearch, countSql, selectSql string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("GetSymbols")

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")
	symbolTable := goqu.T("symbol")

	// Join, where.
	joinClause := dialect.From(
		symbolTable,
	).Where(
		goqu.I("symbol_label").Like(p.GetSearch()),
	)

	if countSql, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("symbol_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}
	if selectSql, selectArgs, err = joinClause.Select(
		goqu.I("symbol_id"),
		goqu.I("symbol_label"),
		goqu.I("symbol_image"),
	).Order(
		goqu.L("INSTR(symbol_label, \"?\")", exactSearch).Asc(),
		goqu.C("symbol_label").Asc(),
	).Limit(uint(p.GetLimit())).Offset(uint(p.GetOffset())).ToSQL(); err != nil {
		return nil, 0, err
	}

	// select
	if err = db.Select(&symbols, selectSql, selectArgs...); err != nil {
		return nil, 0, err
	}
	// count
	if err = db.Get(&count, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	logger.Log.WithFields(logrus.Fields{"symbols": symbols}).Debug("GetSymbols")

	return symbols, count, nil

}

// GetSymbol return the formula matching the given id
func (db *SQLiteDataStore) GetSymbol(id int) (Symbol, error) {

	var (
		err    error
		sqlr   string
		args   []interface{}
		symbol Symbol
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetSymbol")

	dialect := goqu.Dialect("sqlite3")
	symbolTable := goqu.T("symbol")

	sQuery := dialect.From(symbolTable).Where(
		goqu.I("symbol_id").Eq(id),
	).Select(
		goqu.I("symbol_id"),
		goqu.I("symbol_label"),
		goqu.I("symbol_image"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return Symbol{}, err
	}

	if err = db.Get(&symbol, sqlr, args...); err != nil {
		return Symbol{}, err
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "symbol": symbol}).Debug("GetSymbol")

	return symbol, nil

}

// GetSymbolByLabel return the symbol matching the given symbol
func (db *SQLiteDataStore) GetSymbolByLabel(label string) (Symbol, error) {

	var (
		err    error
		sqlr   string
		args   []interface{}
		symbol Symbol
	)
	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetSymbolByLabel")

	dialect := goqu.Dialect("sqlite3")
	symbolTable := goqu.T("symbol")

	sQuery := dialect.From(symbolTable).Where(
		goqu.I("symbol_label").Eq(label),
	).Select(
		goqu.I("symbol_id"),
		goqu.I("symbol_label"),
		goqu.I("symbol_image"),
	).Order(goqu.I("symbol_label").Asc())

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return Symbol{}, err
	}

	if err = db.Get(&symbol, sqlr, args...); err != nil {
		return Symbol{}, err
	}

	logger.Log.WithFields(logrus.Fields{"label": label, "symbol": symbol}).Debug("GetSymbolByLabel")

	return symbol, nil

}

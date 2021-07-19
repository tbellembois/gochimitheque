package datastores

import (
	"database/sql"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"

	. "github.com/tbellembois/gochimitheque/models"
)

// GetCategories return the categories matching the search criteria
func (db *SQLiteDataStore) GetCategories(p SelectFilter) ([]Category, int, error) {

	var (
		err                              error
		categories                       []Category
		count                            int
		exactSearch, countSql, selectSql string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("GetCategories")

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")
	categoryTable := goqu.T("category")

	// Join, where.
	joinClause := dialect.From(
		categoryTable,
	).Where(
		goqu.I("category_label").Like(p.GetSearch()),
	)

	if countSql, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("category_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}
	if selectSql, selectArgs, err = joinClause.Select(
		goqu.I("category_id"),
		goqu.I("category_label"),
	).Order(
		goqu.L("INSTR(category_label, \"?\")", exactSearch).Asc(),
		goqu.C("category_label").Asc(),
	).Limit(uint(p.GetLimit())).Offset(uint(p.GetOffset())).ToSQL(); err != nil {
		return nil, 0, err
	}

	// select
	if err = db.Select(&categories, selectSql, selectArgs...); err != nil {
		return nil, 0, err
	}
	// count
	if err = db.Get(&count, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for formula matching exactly the search
	sQuery := dialect.From(categoryTable).Where(
		goqu.I("category_label").Eq(exactSearch),
	).Select(
		"category_id",
		"category_label",
	)

	var (
		sqlr     string
		args     []interface{}
		category Category
	)
	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return nil, 0, err
	}

	if err = db.Get(&category, sqlr, args...); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	for i, e := range categories {
		if e.CategoryID == category.CategoryID {
			categories[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"categories": categories}).Debug("GetCategories")

	return categories, count, nil

}

// GetCategory return the formula matching the given id
func (db *SQLiteDataStore) GetCategory(id int) (Category, error) {

	var (
		err      error
		sqlr     string
		args     []interface{}
		category Category
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetCategory")

	dialect := goqu.Dialect("sqlite3")
	categoryTable := goqu.T("category")

	sQuery := dialect.From(categoryTable).Where(
		goqu.I("category_id").Eq(id),
	).Select(
		goqu.I("category_id"),
		goqu.I("category_label"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return Category{}, err
	}

	if err = db.Get(&category, sqlr, args...); err != nil {
		return Category{}, err
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "category": category}).Debug("GetCategory")

	return category, nil

}

// GetCategoryByLabel return the category matching the given category
func (db *SQLiteDataStore) GetCategoryByLabel(label string) (Category, error) {

	var (
		err      error
		sqlr     string
		args     []interface{}
		category Category
	)
	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetCategoryByLabel")

	dialect := goqu.Dialect("sqlite3")
	categoryTable := goqu.T("category")

	sQuery := dialect.From(categoryTable).Where(
		goqu.I("category_label").Eq(label),
	).Select(
		goqu.I("category_id"),
		goqu.I("category_label"),
	).Order(goqu.I("category_label").Asc())

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return Category{}, err
	}

	if err = db.Get(&category, sqlr, args...); err != nil {
		return Category{}, err
	}

	logger.Log.WithFields(logrus.Fields{"label": label, "category": category}).Debug("GetCategoryByLabel")

	return category, nil

}

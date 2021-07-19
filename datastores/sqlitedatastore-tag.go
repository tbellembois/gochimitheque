package datastores

import (
	"database/sql"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"

	. "github.com/tbellembois/gochimitheque/models"
)

// GetTags return the tags matching the search criteria
func (db *SQLiteDataStore) GetTags(p SelectFilter) ([]Tag, int, error) {

	var (
		err                              error
		tags                             []Tag
		count                            int
		exactSearch, countSql, selectSql string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("GetTags")

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")
	tagTable := goqu.T("tag")

	// Join, where.
	joinClause := dialect.From(
		tagTable,
	).Where(
		goqu.I("tag_label").Like(p.GetSearch()),
	)

	if countSql, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("tag_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}
	if selectSql, selectArgs, err = joinClause.Select(
		goqu.I("tag_id"),
		goqu.I("tag_label"),
	).Order(
		goqu.L("INSTR(tag_label, \"?\")", exactSearch).Asc(),
		goqu.C("tag_label").Asc(),
	).Limit(uint(p.GetLimit())).Offset(uint(p.GetOffset())).ToSQL(); err != nil {
		return nil, 0, err
	}

	// select
	if err = db.Select(&tags, selectSql, selectArgs...); err != nil {
		return nil, 0, err
	}
	// count
	if err = db.Get(&count, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for formula matching exactly the search
	sQuery := dialect.From(tagTable).Where(
		goqu.I("tag_label").Eq(exactSearch),
	).Select(
		"tag_id",
		"tag_label",
	)

	var (
		sqlr string
		args []interface{}
		tag  Tag
	)
	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return nil, 0, err
	}

	if err = db.Get(&tag, sqlr, args...); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	for i, e := range tags {
		if e.TagID == tag.TagID {
			tags[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"tags": tags}).Debug("GetTags")

	return tags, count, nil

}

// GetTag return the formula matching the given id
func (db *SQLiteDataStore) GetTag(id int) (Tag, error) {

	var (
		err  error
		sqlr string
		args []interface{}
		tag  Tag
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetTag")

	dialect := goqu.Dialect("sqlite3")
	tagTable := goqu.T("tag")

	sQuery := dialect.From(tagTable).Where(
		goqu.I("tag_id").Eq(id),
	).Select(
		goqu.I("tag_id"),
		goqu.I("tag_label"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return Tag{}, err
	}

	if err = db.Get(&tag, sqlr, args...); err != nil {
		return Tag{}, err
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "tag": tag}).Debug("GetTag")

	return tag, nil

}

// GetTagByLabel return the tag matching the given tag
func (db *SQLiteDataStore) GetTagByLabel(label string) (Tag, error) {

	var (
		err  error
		sqlr string
		args []interface{}
		tag  Tag
	)
	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetTagByLabel")

	dialect := goqu.Dialect("sqlite3")
	tagTable := goqu.T("tag")

	sQuery := dialect.From(tagTable).Where(
		goqu.I("tag_label").Eq(label),
	).Select(
		goqu.I("tag_id"),
		goqu.I("tag_label"),
	).Order(goqu.I("tag_label").Asc())

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return Tag{}, err
	}

	if err = db.Get(&tag, sqlr, args...); err != nil {
		return Tag{}, err
	}

	logger.Log.WithFields(logrus.Fields{"label": label, "tag": tag}).Debug("GetTagByLabel")

	return tag, nil

}

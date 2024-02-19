package datastores

import (
	"database/sql"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/zmqclient"
)

func (db *SQLiteDataStore) GetProducers(f zmqclient.RequestFilter) ([]models.Producer, int, error) {
	var (
		err                              error
		producers                        []models.Producer
		count                            int
		exactSearch, countSQL, selectSQL string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"f": f}).Debug("GetProducers")

	// hack to bypass optionnal default on the Rust part.
	if f.Search == "" {
		f.Search = "%%"
	}
	exactSearch = f.Search
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")
	producerTable := goqu.T("producer")

	// Join, where.
	joinClause := dialect.From(
		producerTable,
	).Where(
		goqu.I("producer_label").Like(f.Search),
	)

	if countSQL, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("producer_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}

	if selectSQL, selectArgs, err = joinClause.Select(
		goqu.I("producer_id"),
		goqu.I("producer_label"),
	).Order(
		goqu.L("INSTR(producer_label, ?)", exactSearch).Asc(),
		goqu.C("producer_label").Asc(),
	).Limit(uint(f.Limit)).Offset(uint(f.Offset)).ToSQL(); err != nil {
		return nil, 0, err
	}

	// Select.
	if err = db.Select(&producers, selectSQL, selectArgs...); err != nil {
		return nil, 0, err
	}
	// Count.
	if err = db.Get(&count, countSQL, countArgs...); err != nil {
		return nil, 0, err
	}

	// Setting the C attribute for formula matching exactly the search.
	sQuery := dialect.From(producerTable).Where(
		goqu.I("producer_label").Eq(exactSearch),
	).Select(
		"producer_id",
		"producer_label",
	)

	var (
		sqlr     string
		args     []interface{}
		producer models.Producer
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return nil, 0, err
	}

	if err = db.Get(&producer, sqlr, args...); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	for i, e := range producers {
		if e.ProducerID == producer.ProducerID {
			producers[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"producers": producers}).Debug("GetProducers")

	return producers, count, nil
}

func (db *SQLiteDataStore) GetProducer(id int) (models.Producer, error) {
	var (
		err      error
		sqlr     string
		args     []interface{}
		producer models.Producer
	)

	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetProducer")

	dialect := goqu.Dialect("sqlite3")
	producerTable := goqu.T("producer")

	sQuery := dialect.From(producerTable).Where(
		goqu.I("producer_id").Eq(id),
	).Select(
		goqu.I("producer_id"),
		goqu.I("producer_label"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return models.Producer{}, err
	}

	if err = db.Get(&producer, sqlr, args...); err != nil {
		return models.Producer{}, err
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "producer": producer}).Debug("GetProducer")

	return producer, nil
}

func (db *SQLiteDataStore) GetProducerByLabel(label string) (models.Producer, error) {
	var (
		err      error
		sqlr     string
		args     []interface{}
		producer models.Producer
	)

	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetProducerByLabel")

	dialect := goqu.Dialect("sqlite3")
	producerTable := goqu.T("producer")

	sQuery := dialect.From(producerTable).Where(
		goqu.I("producer_label").Eq(label),
	).Select(
		goqu.I("producer_id"),
		goqu.I("producer_label"),
	).Order(goqu.I("producer_label").Asc())

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return models.Producer{}, err
	}

	if err = db.Get(&producer, sqlr, args...); err != nil {
		return models.Producer{}, err
	}

	logger.Log.WithFields(logrus.Fields{"label": label, "producer": producer}).Debug("GetProducerByLabel")

	return producer, nil
}

func (db *SQLiteDataStore) CreateProducer(p models.Producer) (lastInsertID int64, err error) {
	var (
		sqlr string
		args []interface{}
		tx   *sql.Tx
		res  sql.Result
	)

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("CreateProducer")

	dialect := goqu.Dialect("sqlite3")
	tableProducer := goqu.T("producer")

	if tx, err = db.Begin(); err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			logger.Log.Error(err)

			if rbErr := tx.Rollback(); rbErr != nil {
				logger.Log.Error(rbErr)
				err = rbErr

				return
			}

			return
		}

		err = tx.Commit()
	}()

	iQuery := dialect.Insert(tableProducer).Rows(
		goqu.Record{
			"producer_label": p.ProducerLabel,
		},
	)

	if sqlr, args, err = iQuery.ToSQL(); err != nil {
		return
	}

	if res, err = tx.Exec(sqlr, args...); err != nil {
		return
	}

	if lastInsertID, err = res.LastInsertId(); err != nil {
		return
	}

	return
}

func (db *SQLiteDataStore) GetProducerRefs(f zmqclient.RequestFilter) ([]models.ProducerRef, int, error) {
	var (
		err                              error
		producerRefs                     []models.ProducerRef
		count                            int
		exactSearch, countSQL, selectSQL string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"f": f}).Debug("GetProducerRefs")

	if f.OrderBy == "" {
		f.OrderBy = "producerref_id"
	}

	// hack to bypass optionnal default on the Rust part.
	if f.Search == "" {
		f.Search = "%%"
	}
	exactSearch = f.Search
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")
	producerrefTable := goqu.T("producerref")

	// Join, where.
	whereAnd := []goqu.Expression{
		goqu.I("producerref.producerref_label").Like(f.Search),
	}
	if f.Producer != 0 {
		whereAnd = append(whereAnd, goqu.I("producerref.producer").Eq(f.Producer))
	}

	joinClause := dialect.From(
		producerrefTable,
	).Join(
		goqu.T("producer"),
		goqu.On(
			goqu.Ex{
				"producerref.producer": goqu.I("producer.producer_id"),
			},
		),
	).Where(
		whereAnd...,
	)

	if countSQL, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("producerref.producerref_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}

	if selectSQL, selectArgs, err = joinClause.Select(
		goqu.I("producerref_id"),
		goqu.I("producerref_label"),
		goqu.I("producer_id").As(goqu.C("producer.producer_id")),
		goqu.I("producer_label").As(goqu.C("producer.producer_label")),
	).Order(
		goqu.L("INSTR(producerref_label, ?)", exactSearch).Asc(),
		goqu.C("producerref_label").Asc(),
	).Limit(uint(f.Limit)).Offset(uint(f.Offset)).ToSQL(); err != nil {
		return nil, 0, err
	}

	// Select.
	if err = db.Select(&producerRefs, selectSQL, selectArgs...); err != nil {
		return nil, 0, err
	}
	// Count.
	if err = db.Get(&count, countSQL, countArgs...); err != nil {
		return nil, 0, err
	}

	// Setting the C attribute for formula matching exactly the search.
	sQuery := dialect.From(producerrefTable).Where(
		goqu.I("producerref_label").Eq(exactSearch),
	).Select(
		"producerref_id",
		"producerref_label",
	)

	var (
		sqlr string
		args []interface{}
		pref models.ProducerRef
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return nil, 0, err
	}

	if err = db.Get(&pref, sqlr, args...); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	for i, p := range producerRefs {
		if p.ProducerRefID == pref.ProducerRefID {
			producerRefs[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"producerRefs": producerRefs}).Debug("GetProducerRefs")

	return producerRefs, count, nil
}

package datastores

import (
	"database/sql"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"

	. "github.com/tbellembois/gochimitheque/models"
)

// GetProducers return the producers matching the search criteria
func (db *SQLiteDataStore) GetProducers(p SelectFilter) ([]Producer, int, error) {

	var (
		err                              error
		producers                        []Producer
		count                            int
		exactSearch, countSql, selectSql string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("GetProducers")

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")
	producerTable := goqu.T("producer")

	// Join, where.
	joinClause := dialect.From(
		producerTable,
	).Where(
		goqu.I("producer_label").Like(p.GetSearch()),
	)

	if countSql, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("producer_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}
	if selectSql, selectArgs, err = joinClause.Select(
		goqu.I("producer_id"),
		goqu.I("producer_label"),
	).Order(
		goqu.L("INSTR(producer_label, \"?\")", exactSearch).Asc(),
		goqu.C("producer_label").Asc(),
	).Limit(uint(p.GetLimit())).Offset(uint(p.GetOffset())).ToSQL(); err != nil {
		return nil, 0, err
	}

	// select
	if err = db.Select(&producers, selectSql, selectArgs...); err != nil {
		return nil, 0, err
	}
	// count
	if err = db.Get(&count, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for formula matching exactly the search
	sQuery := dialect.From(producerTable).Where(
		goqu.I("producer_label").Eq(exactSearch),
	).Select(
		"producer_id",
		"producer_label",
	)

	var (
		sqlr     string
		args     []interface{}
		producer Producer
	)
	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return nil, 0, err
	}

	if err = db.Select(&producer, sqlr, args...); err != nil {
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

// GetProducer return the formula matching the given id
func (db *SQLiteDataStore) GetProducer(id int) (Producer, error) {

	var (
		err      error
		sqlr     string
		args     []interface{}
		producer Producer
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
		return Producer{}, err
	}

	if err = db.Get(&producer, sqlr, args...); err != nil {
		return Producer{}, err
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "producer": producer}).Debug("GetProducer")

	return producer, nil

}

// GetProducerByLabel return the producer matching the given producer
func (db *SQLiteDataStore) GetProducerByLabel(label string) (Producer, error) {

	var (
		err      error
		sqlr     string
		args     []interface{}
		producer Producer
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
		return Producer{}, err
	}

	if err = db.Get(&producer, sqlr, args...); err != nil {
		return Producer{}, err
	}

	logger.Log.WithFields(logrus.Fields{"label": label, "producer": producer}).Debug("GetProducerByLabel")

	return producer, nil

}

// CreateProducer create a new producer in the db
func (db *SQLiteDataStore) CreateProducer(p Producer) (lastInsertId int64, err error) {

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

	if lastInsertId, err = res.LastInsertId(); err != nil {
		return
	}

	return

}

// GetProducerRefs return the producerrefs matching the search criteria
func (db *SQLiteDataStore) GetProducerRefs(p SelectFilterProducerRef) ([]ProducerRef, int, error) {

	var (
		err                              error
		producerRefs                     []ProducerRef
		count                            int
		exactSearch, countSql, selectSql string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("GetProducerRefs")

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")
	producerrefTable := goqu.T("producerref")

	// Join, where.
	whereAnd := []goqu.Expression{
		goqu.I("producerref.producerref_label").Like(p.GetSearch()),
	}
	if p.GetProducer() != -1 {
		whereAnd = append(whereAnd, goqu.I("producerref.producer").Eq(p.GetProducer()))
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

	if countSql, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("producerref.producerref_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}
	if selectSql, selectArgs, err = joinClause.Select(
		goqu.I("producerref_id"),
		goqu.I("producerref_label"),
		goqu.I("producer_id").As("producer.producer_id"),
		goqu.I("producer_label").As("producer.producer_label"),
	).Order(
		goqu.L("INSTR(producerref_label, \"?\")", exactSearch).Asc(),
		goqu.C("producerref_label").Asc(),
	).Limit(uint(p.GetLimit())).Offset(uint(p.GetOffset())).ToSQL(); err != nil {
		return nil, 0, err
	}

	// select
	if err = db.Select(&producerRefs, selectSql, selectArgs...); err != nil {
		return nil, 0, err
	}
	// count
	if err = db.Get(&count, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for formula matching exactly the search
	sQuery := dialect.From(producerrefTable).Where(
		goqu.I("producerref_label").Eq(exactSearch),
	).Select(
		"producerref_id",
		"producerref_label",
	)

	var (
		sqlr string
		args []interface{}
		pref ProducerRef
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

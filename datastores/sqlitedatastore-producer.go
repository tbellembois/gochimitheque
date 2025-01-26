package datastores

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
)

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

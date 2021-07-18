package datastores

import (
	"database/sql"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/logger"

	. "github.com/tbellembois/gochimitheque/models"
)

// GetPhysicalStates return the physical states matching the search criteria
func (db *SQLiteDataStore) GetPhysicalStates(p Dbselectparam) ([]PhysicalState, int, error) {

	var (
		err                              error
		physicalStates                   []PhysicalState
		count                            int
		exactSearch, countSql, selectSql string
		countArgs, selectArgs            []interface{}
	)

	logger.Log.WithFields(logrus.Fields{"p": p}).Debug("GetPhysicalStates")

	exactSearch = p.GetSearch()
	exactSearch = strings.TrimPrefix(exactSearch, "%")
	exactSearch = strings.TrimSuffix(exactSearch, "%")

	dialect := goqu.Dialect("sqlite3")
	physicalstateTable := goqu.T("physicalstate")

	// Join, where.
	joinClause := dialect.From(
		physicalstateTable,
	).Where(
		goqu.I("physicalstate_label").Like(p.GetSearch()),
	)

	if countSql, countArgs, err = joinClause.Select(
		goqu.COUNT(goqu.I("physicalstate_id").Distinct()),
	).ToSQL(); err != nil {
		return nil, 0, err
	}
	if selectSql, selectArgs, err = joinClause.Select(
		goqu.I("physicalstate_id"),
		goqu.I("physicalstate_label"),
	).Order(
		goqu.L("INSTR(physicalstate_label, \"?\")", exactSearch).Asc(),
		goqu.C("physicalstate_label").Asc(),
	).Limit(uint(p.GetLimit())).Offset(uint(p.GetOffset())).ToSQL(); err != nil {
		return nil, 0, err
	}

	// select
	if err = db.Select(&physicalStates, selectSql, selectArgs...); err != nil {
		return nil, 0, err
	}
	// count
	if err = db.Get(&count, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	// setting the C attribute for formula matching exactly the search
	sQuery := dialect.From(physicalstateTable).Where(
		goqu.I("physicalstate_label").Eq(exactSearch),
	).Select(
		"physicalstate_id",
		"physicalstate_label",
	)

	var (
		sqlr string
		args []interface{}
		ps   PhysicalState
	)
	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return nil, 0, err
	}

	if err = db.Get(&ps, sqlr, args...); err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	for i, e := range physicalStates {
		if e.PhysicalStateID == ps.PhysicalStateID {
			physicalStates[i].C = 1
		}
	}

	logger.Log.WithFields(logrus.Fields{"physicalStates": physicalStates}).Debug("GetPhysicalStates")

	return physicalStates, count, nil

}

// GetPhysicalState return the formula matching the given id
func (db *SQLiteDataStore) GetPhysicalState(id int) (PhysicalState, error) {

	var (
		err           error
		sqlr          string
		args          []interface{}
		physicalState PhysicalState
	)
	logger.Log.WithFields(logrus.Fields{"id": id}).Debug("GetPhysicalState")

	dialect := goqu.Dialect("sqlite3")
	physicalstateTable := goqu.T("physicalstate")

	sQuery := dialect.From(physicalstateTable).Where(
		goqu.I("physicalstate_id").Eq(id),
	).Select(
		goqu.I("physicalstate_id"),
		goqu.I("physicalstate_label"),
	)

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return PhysicalState{}, err
	}

	if err = db.Get(&physicalState, sqlr, args...); err != nil {
		return PhysicalState{}, err
	}

	logger.Log.WithFields(logrus.Fields{"ID": id, "physicalState": physicalState}).Debug("GetPhysicalState")

	return physicalState, nil

}

// GetPhysicalStateByLabel return the empirirical formula matching the given physical state
func (db *SQLiteDataStore) GetPhysicalStateByLabel(label string) (PhysicalState, error) {

	var (
		err           error
		sqlr          string
		args          []interface{}
		physicalState PhysicalState
	)
	logger.Log.WithFields(logrus.Fields{"label": label}).Debug("GetPhysicalStateByLabel")

	dialect := goqu.Dialect("sqlite3")
	physicalstateTable := goqu.T("physicalstate")

	sQuery := dialect.From(physicalstateTable).Where(
		goqu.I("physicalstate_label").Eq(label),
	).Select(
		goqu.I("physicalstate_id"),
		goqu.I("physicalstate_label"),
	).Order(goqu.I("physicalstate_label").Asc())

	if sqlr, args, err = sQuery.ToSQL(); err != nil {
		logger.Log.Error(err)
		return PhysicalState{}, err
	}

	if err = db.Get(&physicalState, sqlr, args...); err != nil {
		return PhysicalState{}, err
	}

	logger.Log.WithFields(logrus.Fields{"label": label, "physicalState": physicalState}).Debug("GetPhysicalStateByLabel")

	return physicalState, nil

}

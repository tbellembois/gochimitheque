package models

import "database/sql"

// Unit is a volume or weight unit.
type Unit struct {
	UnitID         sql.NullInt64  `db:"unit_id" json:"unit_id" schema:"unit_id" `
	UnitLabel      sql.NullString `db:"unit_label" json:"unit_label" schema:"unit_label" `
	UnitType       sql.NullString `db:"unit_type" json:"unit_type" schema:"unit_type" `
	Unit           *Unit          `db:"unit" json:"unit" schema:"unit"` // reference unit
	UnitMultiplier int            `db:"unit_multiplier" json:"-" schema:"-"`
}

type UnitsResp struct {
	Rows  []Unit `json:"rows"`
	Total int    `json:"total"`
}

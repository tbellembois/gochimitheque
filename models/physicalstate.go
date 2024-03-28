package models

import "database/sql"

// PhysicalState is a product physical state.
type PhysicalState struct {
	// nullable values to handle optional Product foreign key (gorilla shema nil values)
	MatchExactSearch   bool           `db:"match_exact_case" json:"match_exact_case"` // not stored in db but db:"c" set for sqlx
	PhysicalStateID    sql.NullInt64  `db:"physicalstate_id" json:"physicalstate_id" schema:"physicalstate_id" `
	PhysicalStateLabel sql.NullString `db:"physicalstate_label" json:"physicalstate_label" schema:"physicalstate_label" `
}

func (physicalstate PhysicalState) SetC(count int) Searchable {
	if count > 1 {
		physicalstate.MatchExactSearch = true
	} else {
		physicalstate.MatchExactSearch = false
	}
	return physicalstate
}

func (physicalstate PhysicalState) GetTableName() string {
	return ("physicalstate")
}

func (physicalstate PhysicalState) GetIDFieldName() string {
	return ("physicalstate_id")
}

func (physicalstate PhysicalState) GetTextFieldName() string {
	return ("physicalstate_label")
}

func (physicalstate PhysicalState) GetID() int64 {
	return physicalstate.PhysicalStateID.Int64
}

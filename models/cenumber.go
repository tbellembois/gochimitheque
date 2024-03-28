package models

import "database/sql"

// CeNumber is a product CE number.
type CeNumber struct {
	MatchExactSearch bool           `db:"match_exact_case" json:"match_exact_case"` // not stored in db but db:"c" set for sqlx
	CeNumberID       sql.NullInt64  `db:"cenumber_id" json:"cenumber_id" schema:"cenumber_id" `
	CeNumberLabel    sql.NullString `db:"cenumber_label" json:"cenumber_label" schema:"cenumber_label" `
}

func (cenumber CeNumber) SetC(count int) Searchable {
	if count > 1 {
		cenumber.MatchExactSearch = true
	} else {
		cenumber.MatchExactSearch = false
	}

	return cenumber
}

func (cenumber CeNumber) GetTableName() string {
	return ("cenumber")
}

func (cenumber CeNumber) GetIDFieldName() string {
	return ("cenumber_id")
}

func (cenumber CeNumber) GetTextFieldName() string {
	return ("cenumber_label")
}

func (cenumber CeNumber) GetID() int64 {
	return cenumber.CeNumberID.Int64
}

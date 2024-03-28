package models

import "database/sql"

// CasNumber is a product CAS number.
type CasNumber struct {
	MatchExactSearch bool           `db:"match_exact_case" json:"match_exact_case"` // not stored in db but db:"c" set for sqlx
	CasNumberID      sql.NullInt64  `db:"casnumber_id" json:"casnumber_id" schema:"casnumber_id" `
	CasNumberLabel   sql.NullString `db:"casnumber_label" json:"casnumber_label" schema:"casnumber_label" `
	CasNumberCMR     sql.NullString `db:"casnumber_cmr" json:"casnumber_cmr" schema:"casnumber_cmr" `
}

func (casnumber CasNumber) SetC(count int) Searchable {
	if count > 1 {
		casnumber.MatchExactSearch = true
	} else {
		casnumber.MatchExactSearch = false
	}

	return casnumber
}

func (casnumber CasNumber) GetTableName() string {
	return ("casnumber")
}

func (casnumber CasNumber) GetIDFieldName() string {
	return ("casnumber_id")
}

func (casnumber CasNumber) GetTextFieldName() string {
	return ("casnumber_label")
}

func (casnumber CasNumber) GetID() int64 {
	return casnumber.CasNumberID.Int64
}

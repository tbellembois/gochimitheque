package models

// CasNumber is a product CAS number.
type CasNumber struct {
	MatchExactSearch bool `db:"match_exact_search" json:"match_exact_search"` // not stored in db but db:"c" set for sqlx
	// CasNumberID      sql.NullInt64  `db:"cas_number_id" json:"cas_number_id" schema:"cas_number_id" `
	CasNumberID    *int64  `db:"cas_number_id" json:"cas_number_id" schema:"cas_number_id" `
	CasNumberLabel *string `db:"cas_number_label" json:"cas_number_label" schema:"cas_number_label" `
	CasNumberCMR   *string `db:"cas_number_cmr" json:"cas_number_cmr" schema:"cas_number_cmr" `
}

func (cas_number CasNumber) SetC(count int) Searchable {
	if count > 1 {
		cas_number.MatchExactSearch = true
	} else {
		cas_number.MatchExactSearch = false
	}

	return cas_number
}

func (cas_number CasNumber) GetTableName() string {
	return ("cas_number")
}

func (cas_number CasNumber) GetIDFieldName() string {
	return ("cas_number_id")
}

func (cas_number CasNumber) GetTextFieldName() string {
	return ("cas_number_label")
}

func (cas_number CasNumber) GetID() int64 {
	// return cas_number.CasNumberID.Int64
	return *cas_number.CasNumberID
}

package models

// CeNumber is a product CE number.
type CeNumber struct {
	MatchExactSearch bool    `db:"match_exact_case" json:"match_exact_case"` // not stored in db but db:"c" set for sqlx
	CeNumberID       *int64  `db:"ce_number_id" json:"ce_number_id" schema:"ce_number_id" `
	CeNumberLabel    *string `db:"ce_number_label" json:"ce_number_label" schema:"ce_number_label" `
}

func (ce_number CeNumber) SetC(count int) Searchable {
	if count > 1 {
		ce_number.MatchExactSearch = true
	} else {
		ce_number.MatchExactSearch = false
	}

	return ce_number
}

func (ce_number CeNumber) GetTableName() string {
	return ("ce_number")
}

func (ce_number CeNumber) GetIDFieldName() string {
	return ("ce_number_id")
}

func (ce_number CeNumber) GetTextFieldName() string {
	return ("ce_number_label")
}

func (ce_number CeNumber) GetID() int64 {
	return *ce_number.CeNumberID
}

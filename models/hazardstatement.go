package models

// HazardStatement is a product hazard statement.
type HazardStatement struct {
	MatchExactSearch bool `db:"match_exact_search" json:"match_exact_search"` // not stored in db but db:"c" set for sqlx

	HazardStatementID        *int64  `db:"hazard_statement_id" json:"hazard_statement_id" schema:"hazard_statement_id"`
	HazardStatementLabel     string  `db:"hazard_statement_label" json:"hazard_statement_label" schema:"hazard_statement_label"`
	HazardStatementReference string  `db:"hazard_statement_reference" json:"hazard_statement_reference" schema:"hazard_statement_reference"`
	HazardStatementCMR       *string `db:"hazard_statement_cmr" json:"hazard_statement_cmr" schema:"hazard_statement_cmr" `
}

func (hs HazardStatement) SetC(count int) Searchable {
	return hs
}

func (hs HazardStatement) GetTableName() string {
	return ("hazard_statement")
}

func (hs HazardStatement) GetIDFieldName() string {
	return ("hazard_statement_id")
}

func (hs HazardStatement) GetTextFieldName() string {
	return ("hazard_statement_reference")
}

func (hs HazardStatement) GetID() int64 {
	return int64(*hs.HazardStatementID)
}

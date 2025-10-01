package models

// PrecautionaryStatement is a product precautionary statement.
type PrecautionaryStatement struct {
	MatchExactSearch bool `db:"match_exact_search" json:"match_exact_search"` // not stored in db but db:"c" set for sqlx

	PrecautionaryStatementID        *int64 `db:"precautionary_statement_id" json:"precautionary_statement_id" schema:"precautionary_statement_id"`
	PrecautionaryStatementLabel     string `db:"precautionary_statement_label" json:"precautionary_statement_label" schema:"precautionary_statement_label"`
	PrecautionaryStatementReference string `db:"precautionary_statement_reference" json:"precautionary_statement_reference" schema:"precautionary_statement_reference"`
}

func (ps PrecautionaryStatement) SetC(count int) Searchable {
	return ps
}

func (ps PrecautionaryStatement) GetTableName() string {
	return ("precautionary_statement")
}

func (ps PrecautionaryStatement) GetIDFieldName() string {
	return ("precautionary_statement_id")
}

func (ps PrecautionaryStatement) GetTextFieldName() string {
	return ("precautionary_statement_reference")
}

func (ps PrecautionaryStatement) GetID() int64 {
	return *ps.PrecautionaryStatementID
}

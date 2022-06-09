package models

// PrecautionaryStatement is a product precautionary statement.
type PrecautionaryStatement struct {
	PrecautionaryStatementID        int    `db:"precautionarystatement_id" json:"precautionarystatement_id" schema:"precautionarystatement_id"`
	PrecautionaryStatementLabel     string `db:"precautionarystatement_label" json:"precautionarystatement_label" schema:"precautionarystatement_label"`
	PrecautionaryStatementReference string `db:"precautionarystatement_reference" json:"precautionarystatement_reference" schema:"precautionarystatement_reference"`
}

func (ps PrecautionaryStatement) SetC(count int) Searchable {
	return ps
}

func (ps PrecautionaryStatement) GetTableName() string {
	return ("precautionarystatement")
}

func (ps PrecautionaryStatement) GetIDFieldName() string {
	return ("precautionarystatement_id")
}

func (ps PrecautionaryStatement) GetTextFieldName() string {
	return ("precautionarystatement_reference")
}

func (ps PrecautionaryStatement) GetID() int64 {
	return int64(ps.PrecautionaryStatementID)
}

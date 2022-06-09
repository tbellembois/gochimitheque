package models

import "database/sql"

// HazardStatement is a product hazard statement.
type HazardStatement struct {
	HazardStatementID        int            `db:"hazardstatement_id" json:"hazardstatement_id" schema:"hazardstatement_id"`
	HazardStatementLabel     string         `db:"hazardstatement_label" json:"hazardstatement_label" schema:"hazardstatement_label"`
	HazardStatementReference string         `db:"hazardstatement_reference" json:"hazardstatement_reference" schema:"hazardstatement_reference"`
	HazardStatementCMR       sql.NullString `db:"hazardstatement_cmr" json:"hazardstatement_cmr" schema:"hazardstatement_cmr" `
}

func (hs HazardStatement) SetC(count int) Searchable {
	return hs
}

func (hs HazardStatement) GetTableName() string {
	return ("hazardstatement")
}

func (hs HazardStatement) GetIDFieldName() string {
	return ("hazardstatement_id")
}

func (hs HazardStatement) GetTextFieldName() string {
	return ("hazardstatement_reference")
}

func (hs HazardStatement) GetID() int64 {
	return int64(hs.HazardStatementID)
}

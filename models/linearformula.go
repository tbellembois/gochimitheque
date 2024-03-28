package models

import "database/sql"

// LinearFormula is a product linear formula.
type LinearFormula struct {
	MatchExactSearch   bool           `db:"match_exact_case" json:"match_exact_case"` // not stored in db but db:"c" set for sqlx
	LinearFormulaID    sql.NullInt64  `db:"linearformula_id" json:"linearformula_id" schema:"linearformula_id" `
	LinearFormulaLabel sql.NullString `db:"linearformula_label" json:"linearformula_label" schema:"linearformula_label" `
}

func (linearformula LinearFormula) SetC(count int) Searchable {
	if count > 1 {
		linearformula.MatchExactSearch = true
	} else {
		linearformula.MatchExactSearch = false
	}

	return linearformula
}

func (linearformula LinearFormula) GetTableName() string {
	return ("linearformula")
}

func (linearformula LinearFormula) GetIDFieldName() string {
	return ("linearformula_id")
}

func (linearformula LinearFormula) GetTextFieldName() string {
	return ("linearformula_label")
}

func (linearformula LinearFormula) GetID() int64 {
	return linearformula.LinearFormulaID.Int64
}

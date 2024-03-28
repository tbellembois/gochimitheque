package models

import "database/sql"

// EmpiricalFormula is a product empirical formula.
type EmpiricalFormula struct {
	MatchExactSearch      bool           `db:"match_exact_case" json:"match_exact_case"` // not stored in db but db:"c" set for sqlx
	EmpiricalFormulaID    sql.NullInt64  `db:"empiricalformula_id" json:"empiricalformula_id" schema:"empiricalformula_id" `
	EmpiricalFormulaLabel sql.NullString `db:"empiricalformula_label" json:"empiricalformula_label" schema:"empiricalformula_label" `
}

func (empiricalformula EmpiricalFormula) SetC(count int) Searchable {
	if count > 1 {
		empiricalformula.MatchExactSearch = true
	} else {
		empiricalformula.MatchExactSearch = false
	}

	return empiricalformula
}

func (empiricalformula EmpiricalFormula) GetTableName() string {
	return ("empiricalformula")
}

func (empiricalformula EmpiricalFormula) GetIDFieldName() string {
	return ("empiricalformula_id")
}

func (empiricalformula EmpiricalFormula) GetTextFieldName() string {
	return ("empiricalformula_label")
}

func (empiricalformula EmpiricalFormula) GetID() int64 {
	return empiricalformula.EmpiricalFormulaID.Int64
}

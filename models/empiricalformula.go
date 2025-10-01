package models

// EmpiricalFormula is a product empirical formula.
type EmpiricalFormula struct {
	MatchExactSearch      bool    `db:"match_exact_search" json:"match_exact_search"` // not stored in db but db:"c" set for sqlx
	EmpiricalFormulaID    *int64  `db:"empirical_formula_id" json:"empirical_formula_id" schema:"empirical_formula_id" `
	EmpiricalFormulaLabel *string `db:"empirical_formula_label" json:"empirical_formula_label" schema:"empirical_formula_label" `
}

func (empirical_formula EmpiricalFormula) SetC(count int) Searchable {
	if count > 1 {
		empirical_formula.MatchExactSearch = true
	} else {
		empirical_formula.MatchExactSearch = false
	}

	return empirical_formula
}

func (empirical_formula EmpiricalFormula) GetTableName() string {
	return ("empirical_formula")
}

func (empirical_formula EmpiricalFormula) GetIDFieldName() string {
	return ("empirical_formula_id")
}

func (empirical_formula EmpiricalFormula) GetTextFieldName() string {
	return ("empirical_formula_label")
}

func (empirical_formula EmpiricalFormula) GetID() int64 {
	return *empirical_formula.EmpiricalFormulaID
}

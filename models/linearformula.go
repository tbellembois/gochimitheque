package models

// LinearFormula is a product linear formula.
type LinearFormula struct {
	MatchExactSearch   bool    `db:"match_exact_search" json:"match_exact_search"` // not stored in db but db:"c" set for sqlx
	LinearFormulaID    *int64  `db:"linear_formula_id" json:"linear_formula_id" schema:"linear_formula_id" `
	LinearFormulaLabel *string `db:"linear_formula_label" json:"linear_formula_label" schema:"linear_formula_label" `
}

func (linear_formula LinearFormula) SetC(count int) Searchable {
	if count > 1 {
		linear_formula.MatchExactSearch = true
	} else {
		linear_formula.MatchExactSearch = false
	}

	return linear_formula
}

func (linear_formula LinearFormula) GetTableName() string {
	return ("linear_formula")
}

func (linear_formula LinearFormula) GetIDFieldName() string {
	return ("linear_formula_id")
}

func (linear_formula LinearFormula) GetTextFieldName() string {
	return ("linear_formula_label")
}

func (linear_formula LinearFormula) GetID() int64 {
	return *linear_formula.LinearFormulaID
}

package models

import "database/sql"

// LinearFormula is a product linear formula.
type LinearFormula struct {
	C                  int            `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	LinearFormulaID    sql.NullInt64  `db:"linearformula_id" json:"linearformula_id" schema:"linearformula_id" `
	LinearFormulaLabel sql.NullString `db:"linearformula_label" json:"linearformula_label" schema:"linearformula_label" `
}

func (linearformula LinearFormula) SetC(count int) Searchable {
	linearformula.C = count

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

package models

import "database/sql"

// Category is a product category.
type Category struct {
	MatchExactSearch bool           `db:"match_exact_case" json:"match_exact_case"` // not stored in db but db:"c" set for sqlx
	CategoryID       sql.NullInt64  `db:"category_id" json:"category_id" schema:"category_id" `
	CategoryLabel    sql.NullString `db:"category_label" json:"category_label" schema:"category_label" `
}

func (category Category) SetC(count int) Searchable {
	if count > 1 {
		category.MatchExactSearch = true
	} else {
		category.MatchExactSearch = false
	}

	return category
}

func (category Category) GetTableName() string {
	return ("category")
}

func (category Category) GetIDFieldName() string {
	return ("category_id")
}

func (category Category) GetTextFieldName() string {
	return ("category_label")
}

func (category Category) GetID() int64 {
	return category.CategoryID.Int64
}

package models

// ClassOfCompound is a product class of compound.
type ClassOfCompound struct {
	// nullable values to handle optional Product foreign key (gorilla shema nil values)
	MatchExactSearch     bool   `db:"match_exact_case" json:"match_exact_case"` // not stored in db but db:"c" set for sqlx
	ClassOfCompoundID    int    `db:"class_of_compound_id" json:"class_of_compound_id" schema:"class_of_compound_id" `
	ClassOfCompoundLabel string `db:"class_of_compound_label" json:"class_of_compound_label" schema:"class_of_compound_label" `
}

func (coc ClassOfCompound) SetC(count int) Searchable {
	if count > 1 {
		coc.MatchExactSearch = true
	} else {
		coc.MatchExactSearch = false
	}

	return coc
}

func (coc ClassOfCompound) GetTableName() string {
	return ("class_of_compound")
}

func (coc ClassOfCompound) GetIDFieldName() string {
	return ("class_of_compound_id")
}

func (coc ClassOfCompound) GetTextFieldName() string {
	return ("class_of_compound_label")
}

func (coc ClassOfCompound) GetID() int64 {
	return int64(coc.ClassOfCompoundID)
}

package models

// ClassOfCompound is a product class of compound.
type ClassOfCompound struct {
	// nullable values to handle optional Product foreign key (gorilla shema nil values)
	MatchExactSearch     bool   `db:"match_exact_case" json:"match_exact_case"` // not stored in db but db:"c" set for sqlx
	ClassOfCompoundID    int    `db:"classofcompound_id" json:"classofcompound_id" schema:"classofcompound_id" `
	ClassOfCompoundLabel string `db:"classofcompound_label" json:"classofcompound_label" schema:"classofcompound_label" `
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
	return ("classofcompound")
}

func (coc ClassOfCompound) GetIDFieldName() string {
	return ("classofcompound_id")
}

func (coc ClassOfCompound) GetTextFieldName() string {
	return ("classofcompound_label")
}

func (coc ClassOfCompound) GetID() int64 {
	return int64(coc.ClassOfCompoundID)
}

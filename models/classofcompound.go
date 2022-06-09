package models

// ClassOfCompound is a product class of compound.
type ClassOfCompound struct {
	// nullable values to handle optional Product foreign key (gorilla shema nil values)
	C                    int    `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	ClassOfCompoundID    int    `db:"classofcompound_id" json:"classofcompound_id" schema:"classofcompound_id" `
	ClassOfCompoundLabel string `db:"classofcompound_label" json:"classofcompound_label" schema:"classofcompound_label" `
}

func (coc ClassOfCompound) SetC(count int) Searchable {
	coc.C = count

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

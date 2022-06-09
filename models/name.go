package models

// Name is a product name.
type Name struct {
	C         int    `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	NameID    int    `db:"name_id" json:"name_id" schema:"name_id"`
	NameLabel string `db:"name_label" json:"name_label" schema:"name_label"`
}

func (name Name) SetC(count int) Searchable {
	name.C = count

	return name
}

func (name Name) GetTableName() string {
	return ("name")
}

func (name Name) GetIDFieldName() string {
	return ("name_id")
}

func (name Name) GetTextFieldName() string {
	return ("name_label")
}

func (name Name) GetID() int64 {
	return int64(name.NameID)
}

package models

// Tag is a product tag.
type Tag struct {
	C        int    `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	TagID    int    `db:"tag_id" json:"tag_id" schema:"tag_id"`
	TagLabel string `db:"tag_label" json:"tag_label" schema:"tag_label"`
}

func (tag Tag) SetC(count int) Searchable {
	tag.C = count

	return tag
}

func (tag Tag) GetTableName() string {
	return ("tag")
}

func (tag Tag) GetIDFieldName() string {
	return ("tag_id")
}

func (tag Tag) GetTextFieldName() string {
	return ("tag_label")
}

func (tag Tag) GetID() int64 {
	return int64(tag.TagID)
}

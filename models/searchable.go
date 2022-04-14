package models

type Searchable interface {
	SetC(int)
	GetTableName() string
	GetIDFieldName() string
	GetTextFieldName() string
	GetID() int64
}

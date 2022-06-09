package models

type Searchable interface {
	SetC(int) Searchable
	GetTableName() string
	GetIDFieldName() string
	GetTextFieldName() string
	GetID() int64
}

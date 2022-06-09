package models

import "database/sql"

// SignalWord is a product signal word.
type SignalWord struct {
	// nullable values to handle optional Product foreign key (gorilla shema nil values)
	SignalWordID    sql.NullInt64  `db:"signalword_id" json:"signalword_id" schema:"signalword_id" `
	SignalWordLabel sql.NullString `db:"signalword_label" json:"signalword_label" schema:"signalword_label" `
}

func (signalword SignalWord) SetC(count int) Searchable {
	return signalword
}

func (signalword SignalWord) GetTableName() string {
	return ("signalword")
}

func (signalword SignalWord) GetIDFieldName() string {
	return ("signalword_id")
}

func (signalword SignalWord) GetTextFieldName() string {
	return ("signalword_label")
}

func (signalword SignalWord) GetID() int64 {
	return signalword.SignalWordID.Int64
}

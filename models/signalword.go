package models

// SignalWord is a product signal word.
type SignalWord struct {
	// nullable values to handle optional Product foreign key (gorilla shema nil values)
	MatchExactSearch bool    `db:"match_exact_search" json:"match_exact_search"` // not stored in db but db:"c" set for sqlx
	SignalWordID     *int64  `db:"signal_word_id" json:"signal_word_id" schema:"signal_word_id" `
	SignalWordLabel  *string `db:"signal_word_label" json:"signal_word_label" schema:"signal_word_label" `
}

func (signal_word SignalWord) SetC(count int) Searchable {
	return signal_word
}

func (signal_word SignalWord) GetTableName() string {
	return ("signal_word")
}

func (signal_word SignalWord) GetIDFieldName() string {
	return ("signal_word_id")
}

func (signal_word SignalWord) GetTextFieldName() string {
	return ("signal_word_label")
}

func (signal_word SignalWord) GetID() int64 {
	return *signal_word.SignalWordID
}

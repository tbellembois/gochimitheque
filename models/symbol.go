package models

// Symbol is a product symbol.
type Symbol struct {
	MatchExactSearch bool `db:"match_exact_search" json:"match_exact_search"` // not stored in db but db:"c" set for sqlx

	SymbolID    *int64 `db:"symbol_id" json:"symbol_id" schema:"symbol_id"`
	SymbolLabel string `db:"symbol_label" json:"symbol_label" schema:"symbol_label"`
}

func (symbol Symbol) SetC(count int) Searchable {
	return symbol
}

func (symbol Symbol) GetTableName() string {
	return ("symbol")
}

func (symbol Symbol) GetIDFieldName() string {
	return ("symbol_id")
}

func (symbol Symbol) GetTextFieldName() string {
	return ("symbol_label")
}

func (symbol Symbol) GetID() int64 {
	return *symbol.SymbolID
}

package models

// PhysicalState is a product physical state.
type PhysicalState struct {
	// nullable values to handle optional Product foreign key (gorilla shema nil values)
	MatchExactSearch   bool    `db:"match_exact_search" json:"match_exact_search"` // not stored in db but db:"c" set for sqlx
	PhysicalStateID    *int64  `db:"physical_state_id" json:"physical_state_id" schema:"physical_state_id" `
	PhysicalStateLabel *string `db:"physical_state_label" json:"physical_state_label" schema:"physical_state_label" `
}

func (physical_state PhysicalState) SetC(count int) Searchable {
	if count > 1 {
		physical_state.MatchExactSearch = true
	} else {
		physical_state.MatchExactSearch = false
	}
	return physical_state
}

func (physical_state PhysicalState) GetTableName() string {
	return ("physical_state")
}

func (physical_state PhysicalState) GetIDFieldName() string {
	return ("physical_state_id")
}

func (physical_state PhysicalState) GetTextFieldName() string {
	return ("physical_state_label")
}

func (physical_state PhysicalState) GetID() int64 {
	return *physical_state.PhysicalStateID
}

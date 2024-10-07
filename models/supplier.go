package models

// Supplier is a product supplier.
type Supplier struct {
	MatchExactSearch bool    `db:"match_exact_search" json:"match_exact_search"` // not stored in db but db:"c" set for sqlx
	SupplierID       *int64  `db:"supplier_id" json:"supplier_id" schema:"supplier_id" `
	SupplierLabel    *string `db:"supplier_label" json:"supplier_label" schema:"supplier_label" `
}

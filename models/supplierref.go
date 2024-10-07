package models

// SupplierRef is a product supplier reference.
type SupplierRef struct {
	MatchExactSearch bool      `db:"match_exact_case" json:"match_exact_case"` // not stored in db but db:"c" set for sqlx
	SupplierRefID    int       `db:"supplier_ref_id" json:"supplier_ref_id" schema:"supplier_ref_id"`
	SupplierRefLabel string    `db:"supplier_ref_label" json:"supplier_ref_label" schema:"supplier_ref_label"`
	Supplier         *Supplier `db:"supplier" json:"supplier" schema:"supplier"`
}

package models

import "database/sql"

// Supplier is a product supplier.
type Supplier struct {
	C             int            `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	SupplierID    sql.NullInt64  `db:"supplier_id" json:"supplier_id" schema:"supplier_id" `
	SupplierLabel sql.NullString `db:"supplier_label" json:"supplier_label" schema:"supplier_label" `
}

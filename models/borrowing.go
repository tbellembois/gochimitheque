package models

import "database/sql"

// Borrowing represent a storage borrowing.
type Borrowing struct {
	BorrowingID      sql.NullInt64  `db:"borrowing_id" json:"borrowing_id" schema:"borrowing_id" `
	BorrowingComment sql.NullString `db:"borrowing_comment" json:"borrowing_comment" schema:"borrowing_comment" `
	Person           *Person        `db:"person" json:"person" schema:"person"` // logged person
	// Storage          `db:"storage" json:"storage" schema:"storage"`
	Borrower *Person `db:"borrower" json:"borrower" schema:"borrower"` // logged person
}

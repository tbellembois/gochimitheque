package models

// Bookmark is a product person bookmark.
type Bookmark struct {
	BookmarkID *int64 `db:"bookmark_id" json:"bookmark_id" schema:"bookmark_id" `
	Person     `db:"person" json:"person" schema:"person"`
	Product    `db:"product" json:"product" schema:"product"`
}

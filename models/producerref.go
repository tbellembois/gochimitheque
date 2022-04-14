package models

import "database/sql"

// ProducerRef is a product producer reference.
type ProducerRef struct {
	C                int            `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	ProducerRefID    sql.NullInt64  `db:"producerref_id" json:"producerref_id" schema:"producerref_id" `
	ProducerRefLabel sql.NullString `db:"producerref_label" json:"producerref_label" schema:"producerref_label" `
	Producer         *Producer      `db:"producer" json:"producer" schema:"producer"`
}

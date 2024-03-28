package models

import "database/sql"

// ProducerRef is a product producer reference.
type ProducerRef struct {
	MatchExactSearch bool           `db:"match_exact_case" json:"match_exact_case"` // not stored in db but db:"c" set for sqlx
	ProducerRefID    sql.NullInt64  `db:"producerref_id" json:"producerref_id" schema:"producerref_id" `
	ProducerRefLabel sql.NullString `db:"producerref_label" json:"producerref_label" schema:"producerref_label" `
	Producer         *Producer      `db:"producer" json:"producer" schema:"producer"`
}

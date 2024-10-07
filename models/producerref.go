package models

// ProducerRef is a product producer reference.
type ProducerRef struct {
	MatchExactSearch bool      `db:"match_exact_case" json:"match_exact_case"` // not stored in db but db:"c" set for sqlx
	ProducerRefID    *int64    `db:"producer_ref_id" json:"producer_ref_id" schema:"producer_ref_id" `
	ProducerRefLabel *string   `db:"producer_ref_label" json:"producer_ref_label" schema:"producer_ref_label" `
	Producer         *Producer `db:"producer" json:"producer" schema:"producer"`
}

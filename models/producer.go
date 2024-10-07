package models

// Producer is a product producer.
type Producer struct {
	MatchExactSearch bool    `db:"match_exact_case" json:"match_exact_case"` // not stored in db but db:"c" set for sqlx
	ProducerID       *int64  `db:"producer_id" json:"producer_id" schema:"producer_id" `
	ProducerLabel    *string `db:"producer_label" json:"producer_label" schema:"producer_label" `
}

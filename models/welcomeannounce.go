package models

// WelcomeAnnounce is the custom welcome page message.
type WelcomeAnnounce struct {
	WelcomeAnnounceID   int    `db:"welcome_announce_id" json:"welcome_announce_id" schema:"welcome_announce_id"`
	WelcomeAnnounceText string `db:"welcome_announce_text" json:"welcome_announce_text" schema:"welcome_announce_text"`
	WelcomeAnnounceHTML string `db:"welcome_announce_html" json:"welcome_announce_html" schema:"welcome_announce_html"`
}

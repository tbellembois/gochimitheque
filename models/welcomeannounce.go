package models

// WelcomeAnnounce is the custom welcome page message.
type WelcomeAnnounce struct {
	WelcomeAnnounceID   int    `db:"welcomeannounce_id" json:"welcomeannounce_id" schema:"welcomeannounce_id"`
	WelcomeAnnounceText string `db:"welcomeannounce_text" json:"welcomeannounce_text" schema:"welcomeannounce_text"`
	WelcomeAnnounceHTML string `db:"welcomeannounce_html" json:"welcomeannounce_html" schema:"welcomeannounce_html"`
}

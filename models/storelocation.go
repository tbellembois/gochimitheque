package models

import "database/sql"

// StoreLocation is where products are stored in entities.
type StoreLocation struct {
	// nullable values to handle optional StoreLocation foreign key (gorilla shema nil values)
	StoreLocationID       sql.NullInt64  `db:"storelocation_id" json:"storelocation_id" schema:"storelocation_id" `
	StoreLocationName     sql.NullString `db:"storelocation_name" json:"storelocation_name" schema:"storelocation_name" `
	StoreLocationCanStore sql.NullBool   `db:"storelocation_canstore" json:"storelocation_canstore" schema:"storelocation_canstore" `
	StoreLocationColor    sql.NullString `db:"storelocation_color" json:"storelocation_color" schema:"storelocation_color" `
	Entity                `db:"entity" json:"entity" schema:"entity"`
	StoreLocation         *StoreLocation `db:"storelocation" json:"storelocation" schema:"storelocation"`
	StoreLocationFullPath string         `db:"storelocation_fullpath" json:"storelocation_fullpath" schema:"storelocation_fullpath"`

	Children []*StoreLocation `db:"-" json:"children" schema:"-"`
	Stocks   []Stock          `db:"-" json:"stock" schema:"-"`
}

type StoreLocationsResp struct {
	Rows  []StoreLocation `json:"rows"`
	Total int             `json:"total"`
}

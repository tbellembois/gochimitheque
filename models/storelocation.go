package models

// StoreLocation is where products are stored in entities.
type StoreLocation struct {
	// nullable values to handle optional StoreLocation foreign key (gorilla shema nil values)
	StoreLocationID       NullInt64  `db:"store_location_id" json:"store_location_id" schema:"store_location_id" `
	StoreLocationName     NullString `db:"store_location_name" json:"store_location_name" schema:"store_location_name" `
	StoreLocationCanStore NullBool   `db:"store_location_can_store" json:"store_location_can_store" schema:"store_location_can_store" `
	StoreLocationColor    NullString `db:"store_location_color" json:"store_location_color" schema:"store_location_color" `
	Entity                `db:"entity" json:"entity" schema:"entity"`
	StoreLocation         *StoreLocation `db:"store_location" json:"store_location" schema:"store_location"`
	StoreLocationFullPath string         `db:"store_location_full_path" json:"store_location_full_path" schema:"store_location_full_path"`

	Children []*StoreLocation `db:"-" json:"children" schema:"-"`
	Stocks   []Stock          `db:"-" json:"stock" schema:"-"`
}

type StoreLocationsResp struct {
	Rows  []StoreLocation `json:"rows"`
	Total int             `json:"total"`
}

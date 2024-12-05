package models

// Stock is a store location stock for a given product.
type Stock struct {
	StoreLocation StoreLocation `json:"store_location"`
	Product       Product       `json:"product"`
	Unit          *Unit         `json:"unit"`
	Quantity      float64       `json:"quantity"`
}

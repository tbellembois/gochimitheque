package models

// Stock is a store location stock for a given product.
type Stock struct {
	Total   float64 `json:"total"`
	Current float64 `json:"current"`
	Unit    Unit    `json:"unit"`
}

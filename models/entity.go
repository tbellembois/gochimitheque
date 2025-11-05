package models

// Entity represent a department, a laboratory...
type Entity struct {
	EntityID          *int64    `db:"entity_id" json:"entity_id" schema:"entity_id"`
	EntityName        string    `db:"entity_name" json:"entity_name" schema:"entity_name"`
	EntityDescription string    `db:"entity_description" json:"entity_description" schema:"entity_description"`
	Managers          *[]Person `db:"-" json:"managers" schema:"managers"`

	// total store location count
	EntityNbStoreLocations *int64 `db:"entity_nb_store_locations" json:"entity_nb_store_locations" schema:"entity_nb_store_locations"` // not in db but sqlx requires the "db" entry
	// total person count
	EntityNbPeople *int64 `db:"entity_nb_people" json:"entity_nb_people" schema:"entity_nb_people"` // not in db but sqlx requires the "db" entry
}

// Equal tests the entity equality.
func (e1 Entity) Equal(e2 Entity) bool {
	return e1.EntityID == e2.EntityID
}

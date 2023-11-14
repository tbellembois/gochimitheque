package models

// Entity represent a department, a laboratory...
type Entity struct {
	EntityID          int       `db:"entity_id" json:"entity_id" schema:"entity_id"`
	EntityName        string    `db:"entity_name" json:"entity_name" schema:"entity_name"`
	EntityDescription string    `db:"entity_description" json:"entity_description" schema:"entity_description"`
	Managers          []*Person `db:"-" json:"managers" schema:"managers"`

	// total store location count
	EntitySLC int `db:"entity_slc" json:"entity_slc" schema:"entity_slc"` // not in db but sqlx requires the "db" entry
	// total person count
	EntityPC int `db:"entity_pc" json:"entity_pc" schema:"entity_pc"` // not in db but sqlx requires the "db" entry
}

// Equal tests the entity equality.
func (e1 Entity) Equal(e2 Entity) bool {
	return e1.EntityID == e2.EntityID
}

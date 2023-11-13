package models

// Person represent a person.
type Person struct {
	PersonID    int           `db:"person_id" json:"person_id" schema:"person_id"`
	PersonEmail string        `db:"person_email" json:"person_email" schema:"person_email"`
	Permissions []*Permission `db:"-" json:"permissions" schema:"permissions"`
	Entities    []*Entity     `db:"-" json:"entities" schema:"entities"`
}

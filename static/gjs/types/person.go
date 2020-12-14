package types

import "github.com/gopherjs/gopherjs/js"

type People struct {
	*js.Object
	Rows  []*Person `js:"rows" json:"rows"`
	Total int       `js:"total" json:"total"`
}
type Person struct {
	*js.Object
	PersonId    int    `js:"person_id" json:"person_id"`
	PersonEmail string `js:"person_email" json:"person_email"`
}

func NewPerson() *Person {

	return &Person{Object: js.Global.Get("Object").New()}

}

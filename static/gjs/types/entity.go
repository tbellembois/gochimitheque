package types

import (
	"github.com/gopherjs/gopherjs/js"
)

type Entity struct {
	*js.Object
	EntityID          int       `js:"entity_id" json:"entity_id"`
	EntityName        string    `js:"entity_name" json:"entity_name"`
	EntityDescription string    `js:"entity_description" json:"entity_description"`
	Managers          []*Person `js:"managers" json:"managers"`
	EntitySLC         int       `js:"entity_slc" json:"entity_slc"`
	EntityPC          int       `js:"entity_pc" json:"entity_pc"`
}

func NewEntity() *Entity {

	e := &Entity{Object: js.Global.Get("Object").New()}
	e.Managers = make([]*Person, 0)
	return e

}

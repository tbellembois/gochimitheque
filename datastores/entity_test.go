package datastores_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	. "github.com/tbellembois/gochimitheque/models"
	. "github.com/tbellembois/gochimitheque/test"
)

var (
	err error
	m   *Person // entity manager
	p   Person  // entity person
)

func TestMain(m *testing.M) {
	TestInit()
	r := m.Run()
	TestOut()
	os.Exit(r)
}

func TestEntityCreateDuplicateName(t *testing.T) {

	e := Entity{EntityName: "EntityTestDuplicate"}
	if _, err = TDatastore.CreateEntity(e); err != nil {
		t.Fatal(err.Error())
	}

	if _, err = TDatastore.CreateEntity(e); err == nil {
		t.Errorf("could create duplicate entity")
	}

}

func TestEntityMembership(t *testing.T) {

	var (
		e   *Entity
		es  []Entity     // a person entities
		ps  []Permission // a person permissions
		err error
	)

	// creating entity manager
	m = &Person{PersonEmail: "m@test.com"}
	if m.PersonID, err = TDatastore.CreatePerson(*m); err != nil {
		t.Fatal(err.Error())
	}

	// creating entity
	e = &Entity{EntityName: "EntityTestMembership", Managers: []*Person{m}}
	if e.EntityID, err = TDatastore.CreateEntity(*e); err != nil {
		t.Fatal(err.Error())
	}

	// creating entity member
	p = Person{PersonEmail: "p@test.com", Entities: []*Entity{e}}
	if p.PersonID, err = TDatastore.CreatePerson(p); err != nil {
		t.Fatal(err.Error())
	}

	// testing people membership
	if es, err = TDatastore.GetPersonEntities(Admin.PersonID, m.PersonID); err != nil {
		t.Fatal(err.Error())
	}
	if !cmp.Equal(*e, es[0]) {
		t.Error("manager not a member of his entity")
	}
	if es, err = TDatastore.GetPersonEntities(Admin.PersonID, p.PersonID); err != nil {
		t.Fatal(err.Error())
	}
	if !cmp.Equal(*e, es[0]) {
		t.Error("person not a member of his entity")
	}

	// testing people permissions
	if ps, err = TDatastore.GetPersonPermissions(m.PersonID); err != nil {
		t.Fatal(err.Error())
	}
	expected := Permission{PermissionPermName: "all", PermissionItemName: "all", PermissionEntityID: e.EntityID}
	if !cmp.Equal(expected, ps[0]) {
		t.Errorf("entity manager has wrong permission: %v expected: %v", ps[0], expected)
	}
	if ps, err = TDatastore.GetPersonPermissions(p.PersonID); err != nil {
		t.Fatal(err.Error())
	}
	expected = Permission{PermissionPermName: "r", PermissionItemName: "entities", PermissionEntityID: e.EntityID}
	if !cmp.Equal(expected, ps[0]) {
		t.Errorf("entity manager has wrong permission: %v expected: %v", ps[0], expected)
	}

}

func TestEntityCreateEmptyName(t *testing.T) {

	var (
		e Entity
	)

	// testing entity creation with an empty name
	e.EntityName = ""
	if _, err = TDatastore.CreateEntity(e); err == nil {
		t.Error("error: could create entity with an empty name")
	}

}

func TestEntityCreate(t *testing.T) {

	var (
		e Entity
	)

	e = Entity{EntityName: "EntityTestCreate"}
	if e.EntityID, err = TDatastore.CreateEntity(e); err != nil {
		t.Fatal(err.Error())
	}
	if e, err = TDatastore.GetEntity(e.EntityID); err != nil {
		t.Fatal(err.Error())
	}
	if e.EntityName != "EntityTestCreate" {
		t.Errorf("EntityTestA expected - got %v", e)
	}

}

func TestEntityUpdate(t *testing.T) {

	var (
		e Entity
	)

	e = Entity{EntityName: "EntityTestUpdate"}
	if e.EntityID, err = TDatastore.CreateEntity(e); err != nil {
		t.Fatal(err.Error())
	}
	e.EntityName = "EntityTestUpdated"
	if err = TDatastore.UpdateEntity(e); err != nil {
		t.Fatal(err.Error())
	}
	if e, err = TDatastore.GetEntity(e.EntityID); err != nil {
		t.Fatal(err.Error())
	}
	if e.EntityName != "EntityTestUpdated" {
		t.Errorf("EntityTestUpdated expected - got %v", e.EntityName)
	}

}

func TestEntityUpdateDuplicateName(t *testing.T) {

	var (
		e1, e2 Entity
	)

	e1 = Entity{EntityName: "EntityTestUpdateDuplicateName1"}
	if e1.EntityID, err = TDatastore.CreateEntity(e1); err != nil {
		t.Fatal(err.Error())
	}
	e2 = Entity{EntityName: "EntityTestUpdateDuplicateName2"}
	if e2.EntityID, err = TDatastore.CreateEntity(e2); err != nil {
		t.Fatal(err.Error())
	}
	e2.EntityName = "EntityTestUpdateDuplicateName1"
	if err = TDatastore.UpdateEntity(e2); err == nil {
		t.Error("could update duplicate entity")
	}

}

func TestEntityUpdateEmptyName(t *testing.T) {

	e := Entity{EntityName: "EntityTestUpdateNameEmpty"}
	if e.EntityID, err = TDatastore.CreateEntity(e); err != nil {
		t.Fatal(err.Error())
	}
	e.EntityName = ""
	if err = TDatastore.UpdateEntity(e); err == nil {
		t.Error("could update an entity with an empty name")
	}

}

func TestEntityDelete(t *testing.T) {

	var (
		e Entity
	)

	e = Entity{EntityName: "EntityTest"}
	if e.EntityID, err = TDatastore.CreateEntity(e); err != nil {
		t.Fatal(err.Error())
	}
	if err = TDatastore.DeleteEntity(e.EntityID); err != nil {
		t.Fatal(err.Error())
	}
	if e, err = TDatastore.GetEntity(e.EntityID); err == nil {
		t.Error("could get deleted entity")
	}

}

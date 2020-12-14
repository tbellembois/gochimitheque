package globals_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/tbellembois/gochimitheque/globals"
	. "github.com/tbellembois/gochimitheque/test"
)

func TestMain(m *testing.M) {
	TestInit()
	r := m.Run()
	TestOut()
	os.Exit(r)
}

type EnforcerReq struct {
	personId    string
	item        string
	action      string
	itemId      string
	expectedRes bool
}

func Enforcer(t *testing.T, e EnforcerReq) {

	var (
		err error
		ok  bool
	)

	if ok, err = globals.Enforcer.Enforce(e.personId, e.action, e.item, e.itemId); err != nil {
		t.Fatalf("enforce error: %v", err)
	}
	if ok != e.expectedRes {
		t.Errorf("enforce error for person %s action %s item %s itemId %s: expected: %t got %t",
			e.personId,
			e.action,
			e.item,
			e.itemId,
			e.expectedRes,
			ok)
	}

}

func TestEnforcer(t *testing.T) {

	// 1	admin@chimitheque.fr
	// 2	manager1@test.com
	// 3	manager2@test.com
	// 4	person1a@test.com
	// 5	person1b@test.com
	// 6	person2a@test.com
	// 7	person2b@test.com

	// 1	sample entity
	// 2	entity1
	// 3	entity2

	// 1	entity1_sl1
	// 2	entity1_sl2
	// 3	entity2_sl1
	// 4	entity2_sl2

	reqs := []EnforcerReq{
		{
			personId:    "1",
			item:        "entities",
			action:      "w",
			itemId:      "-1",
			expectedRes: true,
		},
		{
			personId:    "1",
			item:        "products",
			action:      "r",
			itemId:      "",
			expectedRes: true,
		},
		{
			personId:    "1",
			item:        "storages",
			action:      "w",
			itemId:      "-2",
			expectedRes: true,
		},
		{
			personId:    "1",
			item:        "products",
			action:      "w",
			itemId:      "-2",
			expectedRes: true,
		},
		{
			personId:    strconv.Itoa(Man1.PersonID),
			item:        "entities",
			action:      "r",
			itemId:      strconv.Itoa(E1.EntityID),
			expectedRes: true,
		},
		{
			personId:    strconv.Itoa(Man1.PersonID),
			item:        "entities",
			action:      "w",
			itemId:      strconv.Itoa(E1.EntityID),
			expectedRes: true,
		},
		{
			personId:    strconv.Itoa(Man1.PersonID),
			item:        "entities",
			action:      "r",
			itemId:      strconv.Itoa(E2.EntityID),
			expectedRes: false,
		},
		{
			personId:    strconv.Itoa(Mickey1.PersonID),
			item:        "storages",
			action:      "r",
			itemId:      strconv.Itoa(int(S1a.StorageID.Int64)),
			expectedRes: false,
		},
		{
			personId:    strconv.Itoa(Pluto1.PersonID),
			item:        "storages",
			action:      "w",
			itemId:      strconv.Itoa(int(S1a.StorageID.Int64)),
			expectedRes: true,
		},
		{
			personId:    strconv.Itoa(Man1.PersonID),
			item:        "storages",
			action:      "w",
			itemId:      strconv.Itoa(int(S1a.StorageID.Int64)),
			expectedRes: true,
		},
		{
			personId:    strconv.Itoa(Mickey1.PersonID),
			item:        "storages",
			action:      "w",
			itemId:      strconv.Itoa(int(S1a.StorageID.Int64)),
			expectedRes: false,
		},
		{
			personId:    strconv.Itoa(Mickey1.PersonID),
			item:        "storages",
			action:      "r",
			itemId:      strconv.Itoa(int(S2a.StorageID.Int64)),
			expectedRes: false,
		},
		{
			personId:    strconv.Itoa(Mickey1.PersonID),
			item:        "peoplepass",
			action:      "r",
			itemId:      "-1",
			expectedRes: true,
		},
		{
			personId:    strconv.Itoa(Mickey1.PersonID),
			item:        "storelocations",
			action:      "r",
			itemId:      strconv.Itoa(int(SL1a.StoreLocationID.Int64)),
			expectedRes: false,
		},
		{
			personId:    strconv.Itoa(Pluto1.PersonID),
			item:        "storelocations",
			action:      "r",
			itemId:      strconv.Itoa(int(SL1a.StoreLocationID.Int64)),
			expectedRes: true,
		},
		{
			personId:    strconv.Itoa(Pluto1.PersonID),
			item:        "storelocations",
			action:      "r",
			itemId:      strconv.Itoa(int(SL2a.StoreLocationID.Int64)),
			expectedRes: false,
		},
		{
			personId:    strconv.Itoa(Mickey1.PersonID),
			item:        "people",
			action:      "r",
			itemId:      strconv.Itoa(Pluto1.PersonID),
			expectedRes: false,
		},
		{
			personId:    strconv.Itoa(Pluto1.PersonID),
			item:        "people",
			action:      "r",
			itemId:      strconv.Itoa(Mickey1.PersonID),
			expectedRes: true,
		},
		{
			personId:    strconv.Itoa(Pluto1.PersonID),
			item:        "people",
			action:      "r",
			itemId:      strconv.Itoa(int(Pluto2.PersonID)),
			expectedRes: false,
		},
		{
			personId:    strconv.Itoa(Man1.PersonID),
			item:        "storelocations",
			action:      "w",
			itemId:      strconv.Itoa(int(SL1a.StoreLocationID.Int64)),
			expectedRes: true,
		},
	}

	for _, req := range reqs {
		Enforcer(t, req)
	}

}

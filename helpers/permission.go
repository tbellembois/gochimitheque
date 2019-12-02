package helpers

type PermKey struct {
	View string
	Item string
	Verb string
	Id   string
}
type PermValue struct {
	Type string
	Item string
	Id   string
}

var PermMatrix = map[PermKey]PermValue{
	// application root
	PermKey{View: "", Item: "", Verb: "GET", Id: ""}: PermValue{Type: "r", Item: "products", Id: "-2"},

	// products
	PermKey{View: "v", Item: "products", Verb: "GET", Id: "id"}:  PermValue{Type: "r"},
	PermKey{View: "vc", Item: "products", Verb: "GET", Id: "id"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "products", Verb: "GET", Id: "id"}:   PermValue{Type: "r"},

	PermKey{View: "v", Item: "products", Verb: "GET", Id: "-1"}:  PermValue{Type: "r"},
	PermKey{View: "vc", Item: "products", Verb: "GET", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "products", Verb: "GET", Id: "-1"}:   PermValue{Type: "r"},

	PermKey{View: "v", Item: "products", Verb: "GET", Id: "-2"}:  PermValue{Type: "r"},
	PermKey{View: "vc", Item: "products", Verb: "GET", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "products", Verb: "GET", Id: "-2"}:   PermValue{Type: "r"},

	PermKey{View: "v", Item: "products", Verb: "GET", Id: ""}:  PermValue{Type: "rall"},
	PermKey{View: "vc", Item: "products", Verb: "GET", Id: ""}: PermValue{Type: "wall"},
	PermKey{View: "", Item: "products", Verb: "GET", Id: ""}:   PermValue{Type: "rall"},

	PermKey{View: "", Item: "products", Verb: "PUT", Id: "id"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "products", Verb: "PUT", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "products", Verb: "PUT", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "products", Verb: "PUT", Id: ""}:   PermValue{Type: "wall"},

	PermKey{View: "", Item: "products", Verb: "POST", Id: "id"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "products", Verb: "POST", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "products", Verb: "POST", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "products", Verb: "POST", Id: ""}:   PermValue{Type: "wall"},

	PermKey{View: "", Item: "products", Verb: "DELETE", Id: "id"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "products", Verb: "DELETE", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "products", Verb: "DELETE", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "products", Verb: "DELETE", Id: ""}:   PermValue{Type: "wall"},

	// rproducts
	PermKey{View: "v", Item: "rproducts", Verb: "GET", Id: "id"}:  PermValue{Type: "r"},
	PermKey{View: "vc", Item: "rproducts", Verb: "GET", Id: "id"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "rproducts", Verb: "GET", Id: "id"}:   PermValue{Type: "r"},

	PermKey{View: "v", Item: "rproducts", Verb: "GET", Id: "-1"}:  PermValue{Type: "r"},
	PermKey{View: "vc", Item: "rproducts", Verb: "GET", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "rproducts", Verb: "GET", Id: "-1"}:   PermValue{Type: "r"},

	PermKey{View: "v", Item: "rproducts", Verb: "GET", Id: "-2"}:  PermValue{Type: "r"},
	PermKey{View: "vc", Item: "rproducts", Verb: "GET", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "rproducts", Verb: "GET", Id: "-2"}:   PermValue{Type: "r"},

	PermKey{View: "v", Item: "rproducts", Verb: "GET", Id: ""}:  PermValue{Type: "rall"},
	PermKey{View: "vc", Item: "rproducts", Verb: "GET", Id: ""}: PermValue{Type: "wall"},
	PermKey{View: "", Item: "rproducts", Verb: "GET", Id: ""}:   PermValue{Type: "rall"},

	PermKey{View: "", Item: "rproducts", Verb: "PUT", Id: "id"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "rproducts", Verb: "PUT", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "rproducts", Verb: "PUT", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "rproducts", Verb: "PUT", Id: ""}:   PermValue{Type: "wall"},

	PermKey{View: "", Item: "rproducts", Verb: "POST", Id: "id"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "rproducts", Verb: "POST", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "rproducts", Verb: "POST", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "rproducts", Verb: "POST", Id: ""}:   PermValue{Type: "wall"},

	PermKey{View: "", Item: "rproducts", Verb: "DELETE", Id: "id"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "rproducts", Verb: "DELETE", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "rproducts", Verb: "DELETE", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "rproducts", Verb: "DELETE", Id: ""}:   PermValue{Type: "wall"},

	// storages
	PermKey{View: "v", Item: "storages", Verb: "GET", Id: "id"}:  PermValue{Type: "rent"},
	PermKey{View: "vc", Item: "storages", Verb: "GET", Id: "id"}: PermValue{Type: "wrent"},
	PermKey{View: "", Item: "storages", Verb: "GET", Id: "id"}:   PermValue{Type: "rent"},

	PermKey{View: "v", Item: "storages", Verb: "GET", Id: "-1"}:  PermValue{Type: "r"},
	PermKey{View: "vc", Item: "storages", Verb: "GET", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "storages", Verb: "GET", Id: "-1"}:   PermValue{Type: "r"},

	PermKey{View: "v", Item: "storages", Verb: "GET", Id: "-2"}:  PermValue{Type: "r"},
	PermKey{View: "vc", Item: "storages", Verb: "GET", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "storages", Verb: "GET", Id: "-2"}:   PermValue{Type: "r"},

	PermKey{View: "v", Item: "storages", Verb: "GET", Id: ""}:  PermValue{Type: "rany"},
	PermKey{View: "vc", Item: "storages", Verb: "GET", Id: ""}: PermValue{Type: "wany"},
	PermKey{View: "", Item: "storages", Verb: "GET", Id: ""}:   PermValue{Type: "rany"},

	PermKey{View: "", Item: "storages", Verb: "PUT", Id: "id"}: PermValue{Type: "went"},
	PermKey{View: "", Item: "storages", Verb: "PUT", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "storages", Verb: "PUT", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "storages", Verb: "PUT", Id: ""}:   PermValue{Type: "wany"},

	PermKey{View: "", Item: "storages", Verb: "POST", Id: "id"}: PermValue{Type: "went"},
	PermKey{View: "", Item: "storages", Verb: "POST", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "storages", Verb: "POST", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "storages", Verb: "POST", Id: ""}:   PermValue{Type: "wany"},

	PermKey{View: "", Item: "storages", Verb: "DELETE", Id: "id"}: PermValue{Type: "went"},
	PermKey{View: "", Item: "storages", Verb: "DELETE", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "storages", Verb: "DELETE", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "storages", Verb: "DELETE", Id: ""}:   PermValue{Type: "wany"},

	// borrowings
	PermKey{View: "", Item: "borrowings", Verb: "PUT", Id: "id"}: PermValue{Type: "went", Item: "storages"},

	// stocks
	PermKey{View: "", Item: "stocks", Verb: "GET", Id: "id"}: PermValue{Type: "r", Item: "storages", Id: "-2"},

	// storelocations
	PermKey{View: "v", Item: "storelocations", Verb: "GET", Id: "id"}:  PermValue{Type: "rent"},
	PermKey{View: "vc", Item: "storelocations", Verb: "GET", Id: "id"}: PermValue{Type: "wrent"},
	PermKey{View: "", Item: "storelocations", Verb: "GET", Id: "id"}:   PermValue{Type: "rent"},

	PermKey{View: "v", Item: "storelocations", Verb: "GET", Id: "-1"}:  PermValue{Type: "r"},
	PermKey{View: "vc", Item: "storelocations", Verb: "GET", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "storelocations", Verb: "GET", Id: "-1"}:   PermValue{Type: "r"},

	PermKey{View: "v", Item: "storelocations", Verb: "GET", Id: "-2"}:  PermValue{Type: "r"},
	PermKey{View: "vc", Item: "storelocations", Verb: "GET", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "storelocations", Verb: "GET", Id: "-2"}:   PermValue{Type: "r"},

	PermKey{View: "v", Item: "storelocations", Verb: "GET", Id: ""}:  PermValue{Type: "rany"},
	PermKey{View: "vc", Item: "storelocations", Verb: "GET", Id: ""}: PermValue{Type: "wany"},
	PermKey{View: "", Item: "storelocations", Verb: "GET", Id: ""}:   PermValue{Type: "rany"},

	PermKey{View: "", Item: "storelocations", Verb: "PUT", Id: "id"}: PermValue{Type: "went"},
	PermKey{View: "", Item: "storelocations", Verb: "PUT", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "storelocations", Verb: "PUT", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "storelocations", Verb: "PUT", Id: ""}:   PermValue{Type: "wany"},

	PermKey{View: "", Item: "storelocations", Verb: "POST", Id: "id"}: PermValue{Type: "went"},
	PermKey{View: "", Item: "storelocations", Verb: "POST", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "storelocations", Verb: "POST", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "storelocations", Verb: "POST", Id: ""}:   PermValue{Type: "wany"},

	PermKey{View: "", Item: "storelocations", Verb: "DELETE", Id: "id"}: PermValue{Type: "went"},
	PermKey{View: "", Item: "storelocations", Verb: "DELETE", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "storelocations", Verb: "DELETE", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "storelocations", Verb: "DELETE", Id: ""}:   PermValue{Type: "wany"},

	// entities
	PermKey{View: "v", Item: "entities", Verb: "GET", Id: "id"}:  PermValue{Type: "r"},
	PermKey{View: "vc", Item: "entities", Verb: "GET", Id: "id"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "entities", Verb: "GET", Id: "id"}:   PermValue{Type: "r"},

	PermKey{View: "v", Item: "entities", Verb: "GET", Id: "-1"}:  PermValue{Type: "r"},
	PermKey{View: "vc", Item: "entities", Verb: "GET", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "entities", Verb: "GET", Id: "-1"}:   PermValue{Type: "r"},

	PermKey{View: "v", Item: "entities", Verb: "GET", Id: "-2"}:  PermValue{Type: "r"},
	PermKey{View: "vc", Item: "entities", Verb: "GET", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "entities", Verb: "GET", Id: "-2"}:   PermValue{Type: "r"},

	PermKey{View: "v", Item: "entities", Verb: "GET", Id: ""}:  PermValue{Type: "rany"},
	PermKey{View: "vc", Item: "entities", Verb: "GET", Id: ""}: PermValue{Type: "wany"},
	PermKey{View: "", Item: "entities", Verb: "GET", Id: ""}:   PermValue{Type: "rany"},

	PermKey{View: "", Item: "entities", Verb: "PUT", Id: "id"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "entities", Verb: "PUT", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "entities", Verb: "PUT", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "entities", Verb: "PUT", Id: ""}:   PermValue{Type: "all", Item: "all"},

	PermKey{View: "", Item: "entities", Verb: "POST", Id: "id"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "entities", Verb: "POST", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "entities", Verb: "POST", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "entities", Verb: "POST", Id: ""}:   PermValue{Type: "all", Item: "all"},

	PermKey{View: "", Item: "entities", Verb: "DELETE", Id: "id"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "entities", Verb: "DELETE", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "entities", Verb: "DELETE", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "entities", Verb: "DELETE", Id: ""}:   PermValue{Type: "all", Item: "all"},

	// welcomeannounce
	PermKey{View: "v", Item: "welcomeannounce", Verb: "GET", Id: "id"}: PermValue{Type: "all", Item: "all"},
	PermKey{View: "", Item: "welcomeannounce", Verb: "GET", Id: "id"}:  PermValue{Type: "all", Item: "all"},

	PermKey{View: "v", Item: "welcomeannounce", Verb: "GET", Id: "-1"}: PermValue{Type: "all", Item: "all"},
	PermKey{View: "", Item: "welcomeannounce", Verb: "GET", Id: "-1"}:  PermValue{Type: "all", Item: "all"},

	PermKey{View: "v", Item: "welcomeannounce", Verb: "GET", Id: "-2"}: PermValue{Type: "all", Item: "all"},
	PermKey{View: "", Item: "welcomeannounce", Verb: "GET", Id: "-2"}:  PermValue{Type: "all", Item: "all"},

	PermKey{View: "v", Item: "welcomeannounce", Verb: "GET", Id: ""}: PermValue{Type: "all", Item: "all"},
	PermKey{View: "", Item: "welcomeannounce", Verb: "GET", Id: ""}:  PermValue{Type: "all", Item: "all"},

	// people
	PermKey{View: "v", Item: "people", Verb: "GET", Id: "id"}:  PermValue{Type: "rent"},
	PermKey{View: "vc", Item: "people", Verb: "GET", Id: "id"}: PermValue{Type: "wrent"},
	PermKey{View: "", Item: "people", Verb: "GET", Id: "id"}:   PermValue{Type: "rent"},

	PermKey{View: "v", Item: "people", Verb: "GET", Id: "-1"}:  PermValue{Type: "r"},
	PermKey{View: "vc", Item: "people", Verb: "GET", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "people", Verb: "GET", Id: "-1"}:   PermValue{Type: "r"},

	PermKey{View: "v", Item: "people", Verb: "GET", Id: "-2"}:  PermValue{Type: "r"},
	PermKey{View: "vc", Item: "people", Verb: "GET", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "people", Verb: "GET", Id: "-2"}:   PermValue{Type: "r"},

	PermKey{View: "v", Item: "people", Verb: "GET", Id: ""}:  PermValue{Type: "rany"},
	PermKey{View: "vc", Item: "people", Verb: "GET", Id: ""}: PermValue{Type: "wany"},
	PermKey{View: "", Item: "people", Verb: "GET", Id: ""}:   PermValue{Type: "rany"},

	PermKey{View: "", Item: "people", Verb: "PUT", Id: "id"}: PermValue{Type: "went"},
	PermKey{View: "", Item: "people", Verb: "PUT", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "people", Verb: "PUT", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "people", Verb: "PUT", Id: ""}:   PermValue{Type: "wany"},

	PermKey{View: "", Item: "people", Verb: "POST", Id: "id"}: PermValue{Type: "went"},
	PermKey{View: "", Item: "people", Verb: "POST", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "people", Verb: "POST", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "people", Verb: "POST", Id: ""}:   PermValue{Type: "wany"},

	PermKey{View: "", Item: "people", Verb: "DELETE", Id: "id"}: PermValue{Type: "went"},
	PermKey{View: "", Item: "people", Verb: "DELETE", Id: "-1"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "people", Verb: "DELETE", Id: "-2"}: PermValue{Type: "w"},
	PermKey{View: "", Item: "people", Verb: "DELETE", Id: ""}:   PermValue{Type: "wany"},
}

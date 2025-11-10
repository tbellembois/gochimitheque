package datastores

// Datastore is an interface to be implemented
// to store data.
type Datastore interface {
	ToCasbinJSONAdapter() ([]byte, error)
}

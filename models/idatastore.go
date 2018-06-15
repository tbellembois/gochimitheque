package models

// Datastore is an interface to be implemented
// to store data
type Datastore interface {
	FlushErrors() error
	CreateDatabase() error

	GetEntities(search string, order string, offset uint64, limit uint64) ([]Entity, error)
	GetEntity(id int) (Entity, error)
	DeleteEntity(id int) error
	CreateEntity(e Entity) error
	UpdateEntity(e Entity) error
	IsEntityWithName(name string) (bool, error)
	IsEntityWithNameExcept(string, ...string) (bool, error)

	GetPeople() ([]Person, error)
	GetPerson(id int) (Person, error)
	GetPersonByEmail(email string) (Person, error)
	GetPersonPermissions(id int) ([]Permission, error)
	GetPersonEntities(id int) ([]Entity, error)
	HasPersonPermission(id int, perm string, item string, itemid int) (bool, error)
}

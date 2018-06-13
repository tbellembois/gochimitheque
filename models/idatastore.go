package models

// Datastore is an interface to be implemented
// to store data
type Datastore interface {
	FlushErrors() error
	CreateDatabase() error

	GetEntities(search string, order string, offset uint64, limit uint64) ([]Entity, error)
	GetEntity(int) (Entity, error)
	DeleteEntity(int) error
	CreateEntity(Entity) error
	UpdateEntity(Entity) error
	HasEntityWithName(string) (bool, error)
	HasEntityWithNameExcept(string, ...string) (bool, error)

	GetPeople() ([]Person, error)
	GetPerson(int) (Person, error)
	GetPersonByEmail(string) (Person, error)
	GetPersonPermissions(PersonID int) ([]Permission, error)
	GetPersonEntities(PersonID int) ([]Entity, error)

	HasPermission(int, string, string, int) (bool, error)
}

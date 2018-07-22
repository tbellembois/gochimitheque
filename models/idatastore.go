package models

// Datastore is an interface to be implemented
// to store data
type Datastore interface {
	FlushErrors() error
	CreateDatabase() error

	// entities
	GetEntities(personID int, search string, order string, offset uint64, limit uint64) ([]Entity, int, error)
	GetEntity(id int) (Entity, error)
	GetEntityPeople(id int) ([]Person, error)
	DeleteEntity(id int) error
	CreateEntity(e Entity) (error, int)
	UpdateEntity(e Entity) error
	IsEntityEmpty(id int) (bool, error)
	IsEntityWithName(name string) (bool, error)
	IsEntityWithNameExcept(string, ...string) (bool, error)

	// people
	GetPeople(personID int, search string, order string, offset uint64, limit uint64) ([]Person, int, error)
	GetPerson(id int) (Person, error)
	GetPersonByEmail(email string) (Person, error)
	GetPersonPermissions(id int) ([]Permission, error)
	GetPersonEntities(personID int, id int) ([]Entity, error)
	GetPersonManageEntities(id int) ([]Entity, error)
	DoesPersonBelongsTo(id int, entities []Entity) (bool, error)
	HasPersonPermission(id int, perm string, item string, itemid int) (bool, error)
	CreatePerson(p Person) (error, int)
	UpdatePerson(p Person) error
	DeletePerson(id int) error
	IsPersonManager(id int) (bool, error)
	IsPersonWithEmail(email string) (bool, error)
	IsPersonWithEmailExcept(string, ...string) (bool, error)
}

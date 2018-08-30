package models

// GetCommonParameters contains the common parameters
// passed to the database Get* functions returning multiple values
// such as GetStoreLocations, GetEntities...
type GetCommonParameters struct {
	LoggedPersonID int
	Search         string
	Order          string
	Offset         uint64
	Limit          uint64
}

// GetPeopleParameters contains the parameters of the GetPeople function
type GetPeopleParameters struct {
	GetCommonParameters
	EntityID int
}

// Datastore is an interface to be implemented
// to store data
type Datastore interface {
	FlushErrors() error
	CreateDatabase() error

	// store locations
	GetStoreLocations(loggedpersonID int, search string, order string, offset uint64, limit uint64) ([]StoreLocation, int, error)
	GetStoreLocation(id int) (StoreLocation, error)
	DeleteStoreLocation(id int) error
	CreateStoreLocation(s StoreLocation) (error, int)
	UpdateStoreLocation(s StoreLocation) error

	// entities
	GetEntities(loggedpersonID int, search string, order string, offset uint64, limit uint64) ([]Entity, int, error)
	GetEntity(id int) (Entity, error)
	GetEntityPeople(id int) ([]Person, error)
	DeleteEntity(id int) error
	CreateEntity(e Entity) (error, int)
	UpdateEntity(e Entity) error
	IsEntityEmpty(id int) (bool, error)
	IsEntityWithName(name string) (bool, error)
	IsEntityWithNameExcept(string, ...string) (bool, error)

	// people
	GetPeople(GetPeopleParameters) ([]Person, int, error)
	GetPerson(id int) (Person, error)
	GetPersonByEmail(email string) (Person, error)
	GetPersonPermissions(id int) ([]Permission, error)
	GetPersonEntities(loggedpersonID int, id int) ([]Entity, error)
	GetPersonManageEntities(id int) ([]Entity, error)
	DoesPersonBelongsTo(id int, entities []Entity) (bool, error)
	HasPersonPermission(id int, perm string, item string, itemid int) (bool, error)
	CreatePerson(p Person) (error, int)
	UpdatePerson(p Person) error
	DeletePerson(id int) error
	IsPersonAdmin(id int) (bool, error)
	IsPersonManager(id int) (bool, error)
	IsPersonWithEmail(email string) (bool, error)
	IsPersonWithEmailExcept(string, ...string) (bool, error)
}

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

// GetEntitiesParameters contains the parameters of the GetEntities function
type GetEntitiesParameters struct {
	GetCommonParameters
}

// GetStoreLocationsParameters contains the parameters of the GetStoreLocations function
type GetStoreLocationsParameters struct {
	GetCommonParameters
	EntityID int
}

// GetProductsParameters contains the parameters of the GetProducts function
type GetProductsParameters struct {
	GetCommonParameters
	EntityID int
}

// Datastore is an interface to be implemented
// to store data
type Datastore interface {
	FlushErrors() error
	CreateDatabase() error

	// products
	GetProducts(GetProductsParameters) ([]Product, int, error)
	GetProductsCasNumbers(GetCommonParameters) ([]CasNumber, int, error)
	GetProductsNames(GetCommonParameters) ([]Name, int, error)
	GetProductsSymbols(GetCommonParameters) ([]Symbol, int, error)
	GetProduct(id int) (Product, error)
	DeleteProduct(id int) error
	CreateProduct(p Product) (error, int)
	UpdateProduct(p Product) error
	IsProductWithName(name string) (bool, error)

	// store locations
	GetStoreLocations(GetStoreLocationsParameters) ([]StoreLocation, int, error)
	GetStoreLocation(id int) (StoreLocation, error)
	GetStoreLocationEntity(id int) (Entity, error)
	DeleteStoreLocation(id int) error
	CreateStoreLocation(s StoreLocation) (error, int)
	UpdateStoreLocation(s StoreLocation) error

	// entities
	GetEntities(GetEntitiesParameters) ([]Entity, int, error)
	GetEntity(id int) (Entity, error)
	GetEntityPeople(id int) ([]Person, error)
	DeleteEntity(id int) error
	CreateEntity(e Entity) (error, int)
	UpdateEntity(e Entity) error
	IsEntityEmpty(id int) (bool, error)

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

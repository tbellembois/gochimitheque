package models

import (
	"net/http"

	"github.com/tbellembois/gochimitheque/helpers"
)

// Datastore is an interface to be implemented
// to store data
type Datastore interface {
	CreateDatabase() error

	// products
	GetProducts(helpers.DbselectparamProduct) ([]Product, int, error)
	GetProductsCasNumbers(helpers.Dbselectparam) ([]CasNumber, int, error)
	GetProductsCasNumber(id int) (CasNumber, error)
	GetProductsCeNumbers(helpers.Dbselectparam) ([]CeNumber, int, error)
	GetProductsNames(helpers.Dbselectparam) ([]Name, int, error)
	GetProductsName(id int) (Name, error)
	GetProductsSymbols(helpers.Dbselectparam) ([]Symbol, int, error)
	GetProductsEmpiricalFormulas(helpers.Dbselectparam) ([]EmpiricalFormula, int, error)
	GetProductsEmpiricalFormula(id int) (EmpiricalFormula, error)
	GetProductsPhysicalStates(helpers.Dbselectparam) ([]PhysicalState, int, error)
	GetProductsSignalWords(helpers.Dbselectparam) ([]SignalWord, int, error)
	GetProductsClassOfCompounds(helpers.Dbselectparam) ([]ClassOfCompound, int, error)
	GetProductsHazardStatements(helpers.Dbselectparam) ([]HazardStatement, int, error)
	GetProductsPrecautionaryStatements(helpers.Dbselectparam) ([]PrecautionaryStatement, int, error)
	GetProduct(id int) (Product, error)
	DeleteProduct(id int) error
	CreateProduct(p Product) (error, int)
	UpdateProduct(p Product) error
	CreateProductBookmark(pr Product, pe Person) error
	DeleteProductBookmark(pr Product, pe Person) error
	IsProductBookmark(pr Product, pe Person) (bool, error)

	// storages
	GetStorages(helpers.DbselectparamStorage) ([]Storage, int, error)
	GetStorage(id int) (Storage, error)
	GetStoragesUnits(helpers.Dbselectparam) ([]Unit, int, error)
	GetStoragesSuppliers(helpers.Dbselectparam) ([]Supplier, int, error)
	DeleteStorage(id int) error
	ArchiveStorage(id int) error
	RestoreStorage(id int) error
	CreateStorage(s Storage) (error, int)
	UpdateStorage(s Storage) error
	GenerateAndUpdateStorageBarecode(s *Storage) error

	// store locations
	GetStoreLocations(helpers.DbselectparamStoreLocation) ([]StoreLocation, int, error)
	GetStoreLocation(id int) (StoreLocation, error)
	GetStoreLocationChildren(id int) ([]StoreLocation, error)
	GetStoreLocationEntity(id int) (Entity, error)
	DeleteStoreLocation(id int) error
	CreateStoreLocation(s StoreLocation) (error, int)
	UpdateStoreLocation(s StoreLocation) error

	// entities
	ComputeStockEntity(p Product, r *http.Request) []StoreLocation

	GetEntities(helpers.DbselectparamEntity) ([]Entity, int, error)
	GetEntity(id int) (Entity, error)
	GetEntityPeople(id int) ([]Person, error)
	DeleteEntity(id int) error
	CreateEntity(e Entity) (error, int)
	UpdateEntity(e Entity) error
	IsEntityEmpty(id int) (bool, error)

	// people
	GetPeople(helpers.DbselectparamPerson) ([]Person, int, error)
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
}

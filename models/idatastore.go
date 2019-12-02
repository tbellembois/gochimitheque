package models

import (
	"net/http"

	"github.com/steambap/captcha"
	"github.com/tbellembois/gochimitheque/helpers"
)

// Datastore is an interface to be implemented
// to store data
type Datastore interface {
	CreateDatabase() error
	ImportV1(dir string) error
	Import(url string) error

	// welcome announce
	GetWelcomeAnnounce() (WelcomeAnnounce, error)
	UpdateWelcomeAnnounce(w WelcomeAnnounce) error

	// products
	GetExposedProducts() ([]Product, int, error)
	GetProducts(helpers.DbselectparamProduct) ([]Product, int, error)

	GetProductsCasNumbers(helpers.Dbselectparam) ([]CasNumber, int, error)
	GetProductsCasNumber(id int) (CasNumber, error)
	GetProductsCasNumberByLabel(label string) (CasNumber, error)

	GetProductsCeNumbers(helpers.Dbselectparam) ([]CeNumber, int, error)
	GetProductsCeNumberByLabel(label string) (CeNumber, error)

	GetProductsNames(helpers.Dbselectparam) ([]Name, int, error)
	GetProductsName(id int) (Name, error)
	GetProductsNameByLabel(label string) (Name, error)

	GetProductsSymbols(helpers.Dbselectparam) ([]Symbol, int, error)
	GetProductsSymbol(id int) (Symbol, error)
	GetProductsSymbolByLabel(label string) (Symbol, error)

	GetProductsEmpiricalFormulas(helpers.Dbselectparam) ([]EmpiricalFormula, int, error)
	GetProductsEmpiricalFormula(id int) (EmpiricalFormula, error)
	GetProductsEmpiricalFormulaByLabel(label string) (EmpiricalFormula, error)

	GetProductsLinearFormulas(helpers.Dbselectparam) ([]LinearFormula, int, error)
	GetProductsLinearFormulaByLabel(label string) (LinearFormula, error)

	GetProductsPhysicalStates(helpers.Dbselectparam) ([]PhysicalState, int, error)
	GetProductsPhysicalStateByLabel(label string) (PhysicalState, error)

	GetProductsSignalWords(helpers.Dbselectparam) ([]SignalWord, int, error)
	GetProductsSignalWord(id int) (SignalWord, error)
	GetProductsSignalWordByLabel(label string) (SignalWord, error)

	GetProductsClassOfCompounds(helpers.Dbselectparam) ([]ClassOfCompound, int, error)
	GetProductsClassOfCompoundByLabel(label string) (ClassOfCompound, error)

	GetProductsHazardStatementByReference(string) (HazardStatement, error)
	GetProductsHazardStatements(helpers.Dbselectparam) ([]HazardStatement, int, error)
	GetProductsHazardStatement(id int) (HazardStatement, error)

	GetProductsPrecautionaryStatementByReference(string) (PrecautionaryStatement, error)
	GetProductsPrecautionaryStatements(helpers.Dbselectparam) ([]PrecautionaryStatement, int, error)
	GetProductsPrecautionaryStatement(id int) (PrecautionaryStatement, error)

	GetProduct(id int) (Product, error)
	CountProductStorages(id int) (int, error)
	DeleteProduct(id int) error
	CreateProduct(p Product) (int, error)
	UpdateProduct(p Product) error
	CreateProductBookmark(pr Product, pe Person) error
	DeleteProductBookmark(pr Product, pe Person) error
	IsProductBookmark(pr Product, pe Person) (bool, error)

	// storages
	GetStorages(helpers.DbselectparamStorage) ([]Storage, int, error)
	GetOtherStorages(helpers.DbselectparamStorage) ([]Entity, int, error)
	GetStorage(id int) (Storage, error)
	GetStoragesUnits(helpers.Dbselectparam) ([]Unit, int, error)
	GetStoragesSuppliers(helpers.Dbselectparam) ([]Supplier, int, error)
	GetStorageEntity(id int) (Entity, error)
	DeleteStorage(id int) error
	ArchiveStorage(id int) error
	RestoreStorage(id int) error
	CreateStorage(s Storage) (int, error)
	UpdateStorage(s Storage) error
	GenerateAndUpdateStorageBarecode(s *Storage) error
	IsStorageBorrowing(b Borrowing) (bool, error)
	CreateStorageBorrowing(b Borrowing) error
	DeleteStorageBorrowing(b Borrowing) error
	UpdateAllQRCodes() error

	// store locations
	GetStoreLocations(helpers.DbselectparamStoreLocation) ([]StoreLocation, int, error)
	GetStoreLocation(id int) (StoreLocation, error)
	GetStoreLocationChildren(id int) ([]StoreLocation, error)
	GetStoreLocationEntity(id int) (Entity, error)
	DeleteStoreLocation(id int) error
	CreateStoreLocation(s StoreLocation) (int, error)
	UpdateStoreLocation(s StoreLocation) error
	IsStoreLocationEmpty(id int) (bool, error)
	ComputeStockStorelocation(p Product, s *StoreLocation, u Unit) float64

	// entities
	ComputeStockEntity(p Product, r *http.Request) []StoreLocation

	GetEntities(helpers.DbselectparamEntity) ([]Entity, int, error)
	GetEntity(id int) (Entity, error)
	GetEntityPeople(id int) ([]Person, error)
	DeleteEntity(id int) error
	CreateEntity(e Entity) (int, error)
	UpdateEntity(e Entity) error
	IsEntityEmpty(id int) (bool, error)
	HasEntityNoStorelocation(id int) (bool, error)

	// people
	GetPeople(helpers.DbselectparamPerson) ([]Person, int, error)
	GetPerson(id int) (Person, error)
	GetPersonByEmail(email string) (Person, error)
	GetPersonPermissions(id int) ([]Permission, error)
	GetPersonEntities(loggedpersonID int, id int) ([]Entity, error)
	GetPersonManageEntities(id int) ([]Entity, error)
	DoesPersonBelongsTo(id int, entities []Entity) (bool, error)
	HasPersonPermission(id int, perm string, item string, eids []int) (bool, error)
	CreatePerson(p Person) (int, error)
	UpdatePerson(p Person) error
	UpdatePersonPassword(p Person) error
	DeletePerson(id int) error
	GetAdmins() ([]Person, error)
	IsPersonAdmin(id int) (bool, error)
	UnsetPersonAdmin(id int) error
	SetPersonAdmin(id int) error
	IsPersonManager(id int) (bool, error)

	// captcha
	InsertCaptcha(*captcha.Data) (string, error)
	ValidateCaptcha(token string, text string) (bool, error)
}

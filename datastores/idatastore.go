package datastores

import (
	"net/http"

	"github.com/steambap/captcha"
	. "github.com/tbellembois/gochimitheque/models"
)

// Datastore is an interface to be implemented
// to store data
type Datastore interface {
	Maintenance()

	CreateDatabase() error
	ImportV1(dir string) error
	Import(url string) error
	ToCasbinJSONAdapter() ([]byte, error)

	// welcome announce
	GetWelcomeAnnounce() (WelcomeAnnounce, error)
	UpdateWelcomeAnnounce(w WelcomeAnnounce) error

	// products
	GetExposedProducts() ([]Product, int, error)
	GetProducts(DbselectparamProduct) ([]Product, int, error)

	GetProductsCasNumbers(Dbselectparam) ([]CasNumber, int, error)
	GetProductsCasNumber(id int) (CasNumber, error)
	GetProductsCasNumberByLabel(label string) (CasNumber, error)

	GetProductsCeNumbers(Dbselectparam) ([]CeNumber, int, error)
	GetProductsCeNumberByLabel(label string) (CeNumber, error)

	GetProductsNames(Dbselectparam) ([]Name, int, error)
	GetProductsName(id int) (Name, error)
	GetProductsNameByLabel(label string) (Name, error)

	GetProductsSymbols(Dbselectparam) ([]Symbol, int, error)
	GetProductsSymbol(id int) (Symbol, error)
	GetProductsSymbolByLabel(label string) (Symbol, error)

	GetProductsEmpiricalFormulas(Dbselectparam) ([]EmpiricalFormula, int, error)
	GetProductsEmpiricalFormula(id int) (EmpiricalFormula, error)
	GetProductsEmpiricalFormulaByLabel(label string) (EmpiricalFormula, error)

	GetProductsLinearFormulas(Dbselectparam) ([]LinearFormula, int, error)
	GetProductsLinearFormulaByLabel(label string) (LinearFormula, error)

	GetProductsPhysicalStates(Dbselectparam) ([]PhysicalState, int, error)
	GetProductsPhysicalStateByLabel(label string) (PhysicalState, error)

	GetProductsSignalWords(Dbselectparam) ([]SignalWord, int, error)
	GetProductsSignalWord(id int) (SignalWord, error)
	GetProductsSignalWordByLabel(label string) (SignalWord, error)

	GetProductsClassOfCompounds(Dbselectparam) ([]ClassOfCompound, int, error)
	GetProductsClassOfCompoundByLabel(label string) (ClassOfCompound, error)

	GetProductsHazardStatementByReference(string) (HazardStatement, error)
	GetProductsHazardStatements(Dbselectparam) ([]HazardStatement, int, error)
	GetProductsHazardStatement(id int) (HazardStatement, error)

	GetProductsPrecautionaryStatementByReference(string) (PrecautionaryStatement, error)
	GetProductsPrecautionaryStatements(Dbselectparam) ([]PrecautionaryStatement, int, error)
	GetProductsPrecautionaryStatement(id int) (PrecautionaryStatement, error)

	GetProductsProducerRefs(DbselectparamProducerRef) ([]ProducerRef, int, error)
	GetProductsProducers(Dbselectparam) ([]Producer, int, error)
	GetProductsCategories(Dbselectparam) ([]Category, int, error)
	GetProductsTags(Dbselectparam) ([]Tag, int, error)
	GetProductsSuppliers(Dbselectparam) ([]Supplier, int, error)
	GetProductsSupplierRefs(DbselectparamSupplierRef) ([]SupplierRef, int, error)

	GetProduct(id int) (Product, error)
	CountProductStorages(id int) (int, error)
	DeleteProduct(id int) error
	CreateProduct(p Product) (int, error)
	UpdateProduct(p Product) error
	CreateProductBookmark(pr Product, pe Person) error
	DeleteProductBookmark(pr Product, pe Person) error
	IsProductBookmark(pr Product, pe Person) (bool, error)

	CreateProducer(p Producer) (int, error)
	CreateSupplier(s Supplier) (int, error)

	// storages
	GetStorages(DbselectparamStorage) ([]Storage, int, error)
	GetOtherStorages(DbselectparamStorage) ([]Entity, int, error)
	GetStorage(id int) (Storage, error)
	GetStoragesUnits(DbselectparamUnit) ([]Unit, int, error)
	GetStoragesSuppliers(Dbselectparam) ([]Supplier, int, error)
	GetStorageEntity(id int) (Entity, error)
	DeleteStorage(id int) error
	ArchiveStorage(id int) error
	RestoreStorage(id int) error
	CreateStorage(s Storage, itemNumber int) (int, error)
	UpdateStorage(s Storage) error
	ToogleStorageBorrowing(s Storage) error
	UpdateAllQRCodes() error

	// store locations
	GetStoreLocations(DbselectparamStoreLocation) ([]StoreLocation, int, error)
	GetStoreLocation(id int) (StoreLocation, error)
	GetStoreLocationChildren(id int) ([]StoreLocation, error)
	DeleteStoreLocation(id int) error
	CreateStoreLocation(s StoreLocation) (int64, error)
	UpdateStoreLocation(s StoreLocation) error
	HasStorelocationStorage(id int) (bool, error)

	// entities
	ComputeStockEntity(p Product, r *http.Request) []StoreLocation

	GetEntities(DbselectparamEntity) ([]Entity, int, error)
	GetEntity(id int) (Entity, error)
	GetEntityManager(id int) ([]Person, error)
	DeleteEntity(id int) error
	CreateEntity(e Entity) (int64, error)
	UpdateEntity(e Entity) error
	HasEntityMember(id int) (bool, error)
	HasEntityStorelocation(id int) (bool, error)

	// people
	GetPeople(DbselectparamPerson) ([]Person, int, error)
	GetPerson(id int) (Person, error)
	GetPersonByEmail(email string) (Person, error)
	GetPersonPermissions(id int) ([]Permission, error)
	GetPersonEntities(loggedpersonID int, id int) ([]Entity, error)
	GetPersonManageEntities(id int) ([]Entity, error)
	DoesPersonBelongsTo(id int, entities []Entity) (bool, error)
	CreatePerson(p Person) (int64, error)
	UpdatePerson(p Person) error
	UpdatePersonPassword(p Person) error
	DeletePerson(id int) error
	GetAdmins() ([]Person, error)
	IsPersonAdmin(id int) (bool, error)
	UnsetPersonAdmin(id int) error
	SetPersonAdmin(id int) error
	IsPersonManager(id int) (bool, error)
	HasPersonReadRestrictedProductPermission(id int) (bool, error)

	// captcha
	InsertCaptcha(string, *captcha.Data) error
	ValidateCaptcha(token string, text string) (bool, error)
}

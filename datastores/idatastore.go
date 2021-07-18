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

	GetCasNumbers(Dbselectparam) ([]CasNumber, int, error)
	GetCasNumber(id int) (CasNumber, error)
	GetCasNumberByLabel(label string) (CasNumber, error)

	GetCeNumbers(Dbselectparam) ([]CeNumber, int, error)
	GetCeNumber(id int) (CeNumber, error)
	GetCeNumberByLabel(label string) (CeNumber, error)

	GetNames(Dbselectparam) ([]Name, int, error)
	GetName(id int) (Name, error)
	GetNameByLabel(label string) (Name, error)

	GetSymbols(Dbselectparam) ([]Symbol, int, error)
	GetSymbol(id int) (Symbol, error)
	GetSymbolByLabel(label string) (Symbol, error)

	GetEmpiricalFormulas(Dbselectparam) ([]EmpiricalFormula, int, error)
	GetEmpiricalFormula(id int) (EmpiricalFormula, error)
	GetEmpiricalFormulaByLabel(label string) (EmpiricalFormula, error)

	GetLinearFormulas(Dbselectparam) ([]LinearFormula, int, error)
	GetLinearFormula(id int) (LinearFormula, error)
	GetLinearFormulaByLabel(label string) (LinearFormula, error)

	GetPhysicalStates(Dbselectparam) ([]PhysicalState, int, error)
	GetPhysicalState(id int) (PhysicalState, error)
	GetPhysicalStateByLabel(label string) (PhysicalState, error)

	GetSignalWords(Dbselectparam) ([]SignalWord, int, error)
	GetSignalWord(id int) (SignalWord, error)
	GetSignalWordByLabel(label string) (SignalWord, error)

	GetClassesOfCompound(Dbselectparam) ([]ClassOfCompound, int, error)
	GetClassOfCompound(id int) (ClassOfCompound, error)
	GetClassOfCompoundByLabel(label string) (ClassOfCompound, error)

	GetHazardStatementByReference(string) (HazardStatement, error)
	GetHazardStatements(Dbselectparam) ([]HazardStatement, int, error)
	GetHazardStatement(id int) (HazardStatement, error)

	GetPrecautionaryStatementByReference(string) (PrecautionaryStatement, error)
	GetPrecautionaryStatements(Dbselectparam) ([]PrecautionaryStatement, int, error)
	GetPrecautionaryStatement(id int) (PrecautionaryStatement, error)

	GetTags(Dbselectparam) ([]Tag, int, error)
	GetTag(id int) (Tag, error)
	GetTagByLabel(label string) (Tag, error)

	GetCategories(Dbselectparam) ([]Category, int, error)
	GetCategory(id int) (Category, error)
	GetCategoryByLabel(label string) (Category, error)

	GetProducers(Dbselectparam) ([]Producer, int, error)
	GetProducer(id int) (Producer, error)
	GetProducerByLabel(label string) (Producer, error)
	CreateProducer(p Producer) (int64, error)

	GetSuppliers(Dbselectparam) ([]Supplier, int, error)
	GetSupplier(id int) (Supplier, error)
	GetSupplierByLabel(label string) (Supplier, error)
	CreateSupplier(s Supplier) (int64, error)

	GetProducerRefs(DbselectparamProducerRef) ([]ProducerRef, int, error)
	GetSupplierRefs(DbselectparamSupplierRef) ([]SupplierRef, int, error)

	GetProduct(id int) (Product, error)
	CountProductStorages(id int) (int, error)
	DeleteProduct(id int) error
	CreateProduct(p Product) (int, error)
	UpdateProduct(p Product) error
	CreateProductBookmark(pr Product, pe Person) error
	DeleteProductBookmark(pr Product, pe Person) error
	IsProductBookmark(pr Product, pe Person) (bool, error)

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

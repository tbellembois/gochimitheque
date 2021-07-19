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

	GetWelcomeAnnounce() (WelcomeAnnounce, error)
	UpdateWelcomeAnnounce(w WelcomeAnnounce) error

	GetProducts(SelectFilterProduct, bool) ([]Product, int, error)
	GetProduct(id int) (Product, error)
	CountProductStorages(id int) (int, error)
	DeleteProduct(id int) error
	CreateUpdateProduct(p Product, update bool) (int64, error)
	CreateProductBookmark(pr Product, pe Person) error
	DeleteProductBookmark(pr Product, pe Person) error
	IsProductBookmark(pr Product, pe Person) (bool, error)

	GetCasNumbers(SelectFilter) ([]CasNumber, int, error)
	GetCasNumber(id int) (CasNumber, error)
	GetCasNumberByLabel(label string) (CasNumber, error)

	GetCeNumbers(SelectFilter) ([]CeNumber, int, error)
	GetCeNumber(id int) (CeNumber, error)
	GetCeNumberByLabel(label string) (CeNumber, error)

	GetNames(SelectFilter) ([]Name, int, error)
	GetName(id int) (Name, error)
	GetNameByLabel(label string) (Name, error)

	GetSymbols(SelectFilter) ([]Symbol, int, error)
	GetSymbol(id int) (Symbol, error)
	GetSymbolByLabel(label string) (Symbol, error)

	GetEmpiricalFormulas(SelectFilter) ([]EmpiricalFormula, int, error)
	GetEmpiricalFormula(id int) (EmpiricalFormula, error)
	GetEmpiricalFormulaByLabel(label string) (EmpiricalFormula, error)

	GetLinearFormulas(SelectFilter) ([]LinearFormula, int, error)
	GetLinearFormula(id int) (LinearFormula, error)
	GetLinearFormulaByLabel(label string) (LinearFormula, error)

	GetPhysicalStates(SelectFilter) ([]PhysicalState, int, error)
	GetPhysicalState(id int) (PhysicalState, error)
	GetPhysicalStateByLabel(label string) (PhysicalState, error)

	GetSignalWords(SelectFilter) ([]SignalWord, int, error)
	GetSignalWord(id int) (SignalWord, error)
	GetSignalWordByLabel(label string) (SignalWord, error)

	GetClassesOfCompound(SelectFilter) ([]ClassOfCompound, int, error)
	GetClassOfCompound(id int) (ClassOfCompound, error)
	GetClassOfCompoundByLabel(label string) (ClassOfCompound, error)

	GetHazardStatementByReference(string) (HazardStatement, error)
	GetHazardStatements(SelectFilter) ([]HazardStatement, int, error)
	GetHazardStatement(id int) (HazardStatement, error)

	GetPrecautionaryStatementByReference(string) (PrecautionaryStatement, error)
	GetPrecautionaryStatements(SelectFilter) ([]PrecautionaryStatement, int, error)
	GetPrecautionaryStatement(id int) (PrecautionaryStatement, error)

	GetTags(SelectFilter) ([]Tag, int, error)
	GetTag(id int) (Tag, error)
	GetTagByLabel(label string) (Tag, error)

	GetCategories(SelectFilter) ([]Category, int, error)
	GetCategory(id int) (Category, error)
	GetCategoryByLabel(label string) (Category, error)

	GetProducers(SelectFilter) ([]Producer, int, error)
	GetProducer(id int) (Producer, error)
	GetProducerByLabel(label string) (Producer, error)
	CreateProducer(p Producer) (int64, error)

	GetSuppliers(SelectFilter) ([]Supplier, int, error)
	GetSupplier(id int) (Supplier, error)
	GetSupplierByLabel(label string) (Supplier, error)
	CreateSupplier(s Supplier) (int64, error)

	GetProducerRefs(SelectFilterProducerRef) ([]ProducerRef, int, error)
	GetSupplierRefs(SelectFilterSupplierRef) ([]SupplierRef, int, error)

	// storages
	GetStorages(SelectFilterStorage) ([]Storage, int, error)
	GetOtherStorages(SelectFilterStorage) ([]Entity, int, error)
	GetStorage(id int) (Storage, error)
	GetStoragesUnits(SelectFilterUnit) ([]Unit, int, error)
	GetStorageEntity(id int) (Entity, error)
	DeleteStorage(id int) error
	ArchiveStorage(id int) error
	RestoreStorage(id int) error
	CreateUpdateStorage(s Storage, itemNumber int, update bool) (int64, error)
	ToogleStorageBorrowing(s Storage) error
	UpdateAllQRCodes() error

	// store locations
	GetStoreLocations(SelectFilterStoreLocation) ([]StoreLocation, int, error)
	GetStoreLocation(id int) (StoreLocation, error)
	GetStoreLocationChildren(id int) ([]StoreLocation, error)
	DeleteStoreLocation(id int) error
	CreateStoreLocation(s StoreLocation) (int64, error)
	UpdateStoreLocation(s StoreLocation) error
	HasStorelocationStorage(id int) (bool, error)

	// entities
	ComputeStockEntity(p Product, r *http.Request) []StoreLocation

	GetEntities(SelectFilterEntity) ([]Entity, int, error)
	GetEntity(id int) (Entity, error)
	GetEntityManager(id int) ([]Person, error)
	DeleteEntity(id int) error
	CreateEntity(e Entity) (int64, error)
	UpdateEntity(e Entity) error
	HasEntityMember(id int) (bool, error)
	HasEntityStorelocation(id int) (bool, error)

	// people
	GetPeople(SelectFilterPerson) ([]Person, int, error)
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

package datastores

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/steambap/captcha"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/request"
)

// Datastore is an interface to be implemented
// to store data.
type Datastore interface {
	GetDB() *sqlx.DB

	Maintenance()

	CreateDatabase() error
	Import(url string) error
	ToCasbinJSONAdapter() ([]byte, error)

	GetWelcomeAnnounce() (models.WelcomeAnnounce, error)
	UpdateWelcomeAnnounce(w models.WelcomeAnnounce) error

	GetProducts(request.Filter, bool) ([]models.Product, int, error)
	GetProduct(id int) (models.Product, error)
	CountProductStorages(id int) (int, error)
	DeleteProduct(id int) error
	CreateUpdateProduct(p models.Product, update bool) (int64, error)
	CreateProductBookmark(pr models.Product, pe models.Person) error
	DeleteProductBookmark(pr models.Product, pe models.Person) error
	IsProductBookmark(pr models.Product, pe models.Person) (bool, error)

	// GetCasNumbers(request.Filter) ([]models.CasNumber, int, error)
	// GetCasNumber(id int) (models.CasNumber, error)
	// GetCasNumberByLabel(label string) (models.CasNumber, error)

	// GetCeNumbers(request.Filter) ([]models.CeNumber, int, error)
	// GetCeNumber(id int) (models.CeNumber, error)
	// GetCeNumberByLabel(label string) (models.CeNumber, error)

	// GetNames(request.Filter) ([]models.Name, int, error)
	// GetName(id int) (models.Name, error)
	// GetNameByLabel(label string) (models.Name, error)

	// GetSymbols(request.Filter) ([]models.Symbol, int, error)
	// GetSymbol(id int) (models.Symbol, error)
	// GetSymbolByLabel(label string) (models.Symbol, error)

	// GetEmpiricalFormulas(request.Filter) ([]models.EmpiricalFormula, int, error)
	// GetEmpiricalFormula(id int) (models.EmpiricalFormula, error)
	// GetEmpiricalFormulaByLabel(label string) (models.EmpiricalFormula, error)

	// GetLinearFormulas(request.Filter) ([]models.LinearFormula, int, error)
	// GetLinearFormula(id int) (models.LinearFormula, error)
	// GetLinearFormulaByLabel(label string) (models.LinearFormula, error)

	// GetPhysicalStates(request.Filter) ([]models.PhysicalState, int, error)
	// GetPhysicalState(id int) (models.PhysicalState, error)
	// GetPhysicalStateByLabel(label string) (models.PhysicalState, error)

	// GetSignalWords(request.Filter) ([]models.SignalWord, int, error)
	// GetSignalWord(id int) (models.SignalWord, error)
	// GetSignalWordByLabel(label string) (models.SignalWord, error)

	// GetClassesOfCompound(request.Filter) ([]models.ClassOfCompound, int, error)
	// GetClassOfCompound(id int) (models.ClassOfCompound, error)
	// GetClassOfCompoundByLabel(label string) (models.ClassOfCompound, error)

	// GetHazardStatementByReference(string) (models.HazardStatement, error)
	// GetHazardStatements(request.Filter) ([]models.HazardStatement, int, error)
	// GetHazardStatement(id int) (models.HazardStatement, error)

	// GetPrecautionaryStatementByReference(string) (models.PrecautionaryStatement, error)
	// GetPrecautionaryStatements(request.Filter) ([]models.PrecautionaryStatement, int, error)
	// GetPrecautionaryStatement(id int) (models.PrecautionaryStatement, error)

	// GetTags(request.Filter) ([]models.Tag, int, error)
	// GetTag(id int) (models.Tag, error)
	// GetTagByLabel(label string) (models.Tag, error)

	// GetCategories(request.Filter) ([]models.Category, int, error)
	// GetCategory(id int) (models.Category, error)
	// GetCategoryByLabel(label string) (models.Category, error)

	GetProducers(request.Filter) ([]models.Producer, int, error)
	GetProducer(id int) (models.Producer, error)
	GetProducerByLabel(label string) (models.Producer, error)
	CreateProducer(p models.Producer) (int64, error)

	GetSuppliers(request.Filter) ([]models.Supplier, int, error)
	GetSupplier(id int) (models.Supplier, error)
	GetSupplierByLabel(label string) (models.Supplier, error)
	CreateSupplier(s models.Supplier) (int64, error)

	GetProducerRefs(request.Filter) ([]models.ProducerRef, int, error)
	GetSupplierRefs(request.Filter) ([]models.SupplierRef, int, error)

	// storages
	GetStorages(request.Filter) ([]models.Storage, int, error)
	GetOtherStorages(request.Filter) ([]models.Entity, int, error)
	GetStorage(id int) (models.Storage, error)
	GetStoragesUnits(request.Filter) ([]models.Unit, int, error)
	GetStorageEntity(id int) (models.Entity, error)
	DeleteStorage(id int) error
	ArchiveStorage(id int) error
	RestoreStorage(id int) error
	CreateUpdateStorage(s models.Storage, itemNumber int, update bool) (int64, error)
	ToogleStorageBorrowing(s models.Storage) error
	UpdateAllQRCodes() error

	// store locations
	GetStoreLocations(request.Filter) ([]models.StoreLocation, int, error)
	GetStoreLocation(id int) (models.StoreLocation, error)
	GetStoreLocationChildren(id int) ([]models.StoreLocation, error)
	DeleteStoreLocation(id int) error
	CreateStoreLocation(s models.StoreLocation) (int64, error)
	UpdateStoreLocation(s models.StoreLocation) error
	HasStorelocationStorage(id int) (bool, error)

	// entities
	ComputeStockEntity(p models.Product, r *http.Request) []models.StoreLocation

	GetEntities(request.Filter) ([]models.Entity, int, error)
	GetEntity(id int) (models.Entity, error)
	GetEntityManager(id int) ([]models.Person, error)
	DeleteEntity(id int) error
	CreateEntity(e models.Entity) (int64, error)
	UpdateEntity(e models.Entity) error
	HasEntityMember(id int) (bool, error)
	HasEntityStorelocation(id int) (bool, error)

	// people
	GetPeople(request.Filter) ([]models.Person, int, error)
	GetPerson(id int) (models.Person, error)
	GetPersonByEmail(email string) (models.Person, error)
	GetPersonPermissions(id int) ([]models.Permission, error)
	GetPersonEntities(loggedpersonID int, id int) ([]models.Entity, error)
	GetPersonManageEntities(id int) ([]models.Entity, error)
	DoesPersonBelongsTo(id int, entities []models.Entity) (bool, error)
	CreatePerson(p models.Person) (int64, error)
	UpdatePerson(p models.Person) error
	UpdatePersonPassword(p models.Person) error
	UpdatePersonAESKey(p models.Person) error
	DeletePerson(id int) error
	GetAdmins() ([]models.Person, error)
	IsPersonAdmin(id int) (bool, error)
	UnsetPersonAdmin(id int) error
	SetPersonAdmin(id int) error
	IsPersonManager(id int) (bool, error)
	HasPersonReadRestrictedProductPermission(id int) (bool, error)

	// captcha
	InsertCaptcha(string, *captcha.Data) error
	ValidateCaptcha(token string, text string) (bool, error)
}

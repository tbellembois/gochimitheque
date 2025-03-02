package datastores

import (
	"github.com/jmoiron/sqlx"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/zmqclient"
)

// Datastore is an interface to be implemented
// to store data.
type Datastore interface {
	GetDB() *sqlx.DB

	Maintenance()

	CreateDatabase() error
	ToCasbinJSONAdapter() ([]byte, error)

	GetWelcomeAnnounce() (models.WelcomeAnnounce, error)
	UpdateWelcomeAnnounce(w models.WelcomeAnnounce) error

	// GetProducts(zmqclient.RequestFilter, int, bool) ([]models.Product, int, error)
	// GetProduct(id int) (models.Product, error)
	CountProductStorages(id int) (int, error)
	CountProducts() (int, error)
	DeleteProduct(id int) error
	CreateUpdateProduct(p models.Product, update bool) (int64, error)
	CreateProductBookmark(pr models.Product, pe models.Person) error
	DeleteProductBookmark(pr models.Product, pe models.Person) error
	IsProductBookmark(pr models.Product, pe models.Person) (bool, error)

	// GetProducers(zmqclient.RequestFilter) ([]models.Producer, int, error)
	// GetProducer(id int) (models.Producer, error)
	// GetProducerByLabel(label string) (models.Producer, error)
	// CreateProducer(p models.Producer) (int64, error)

	// GetSuppliers(zmqclient.RequestFilter) ([]models.Supplier, int, error)
	// GetSupplier(id int) (models.Supplier, error)
	// GetSupplierByLabel(label string) (models.Supplier, error)
	// CreateSupplier(s models.Supplier) (int64, error)

	// GetProducerRefs(zmqclient.RequestFilter) ([]models.ProducerRef, int, error)
	// GetSupplierRefs(zmqclient.RequestFilter) ([]models.SupplierRef, int, error)

	// storages
	GetStorages(zmqclient.RequestFilter, int) ([]models.Storage, int, error)
	GetOtherStorages(zmqclient.RequestFilter, int) ([]models.Entity, int, error)
	GetStorage(id int) (models.Storage, error)
	// GetStoragesUnits(zmqclient.RequestFilter) ([]models.Unit, int, error)
	GetStorageEntity(id int) (models.Entity, error)
	DeleteStorage(id int) error
	ArchiveStorage(id int) error
	RestoreStorage(id int) error
	CreateUpdateStorage(s models.Storage, itemNumber int, update bool) (int64, error)
	ToogleStorageBorrowing(s models.Storage) error
	UpdateAllQRCodes() error

	// store locations
	// GetStoreLocations(zmqclient.RequestFilter, int) ([]models.StoreLocation, int, error)
	// GetStoreLocation(id int) (models.StoreLocation, error)
	// GetStoreLocationChildren(id int) ([]models.StoreLocation, error)
	// DeleteStoreLocation(id int) error
	// CreateStoreLocation(s models.StoreLocation) (int64, error)
	// UpdateStoreLocation(s models.StoreLocation) error
	// HasStorelocationStorage(id int) (bool, error)

	// entities
	// ComputeStockEntity(p models.Product, r *http.Request) []models.StoreLocation

	// GetEntities(zmqclient.RequestFilter, int) ([]models.Entity, int, error)
	// GetEntity(id int) (models.Entity, error)
	// GetEntityManager(id int) ([]models.Person, error)
	DeleteEntity(id int) error
	CreateEntity(e models.Entity) (int64, error)
	UpdateEntity(e models.Entity) error
	// HasEntityMember(id int) (bool, error)
	// HasEntityStorelocation(id int) (bool, error)

	// people
	// GetPeople(zmqclient.RequestFilter, int) ([]models.Person, int, error)
	// GetOrphanPeople() ([]models.Person, error)
	// IsOrphanPerson(id int) (bool, error)
	// GetPerson(id int) (models.Person, error)
	// GetPersonByEmail(email string) (models.Person, error)
	// GetPersonPermissions(id int) ([]models.Permission, error)
	// GetPersonEntities(loggedpersonID int, id int) ([]models.Entity, error)
	// GetPersonManageEntities(id int) ([]models.Entity, error)
	// DoesPersonBelongsTo(id int, entities []models.Entity) (bool, error)
	CreatePerson(p models.Person) (int64, error)
	UpdatePerson(p models.Person) error
	DeletePerson(id int) error
	GetAdmins() ([]models.Person, error)
	IsPersonAdmin(id int) (bool, error)
	UnsetPersonAdmin(id int) error
	SetPersonAdmin(id int) error
	// IsPersonManager(id int) (bool, error)
	HasPersonReadRestrictedProductPermission(id int) (bool, error)
}

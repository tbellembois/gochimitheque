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
	CloseDB() error

	Maintenance()

	CreateDatabase() error
	ToCasbinJSONAdapter() ([]byte, error)

	GetWelcomeAnnounce() (models.WelcomeAnnounce, error)
	UpdateWelcomeAnnounce(w models.WelcomeAnnounce) error

	DeleteProduct(id int) error
	CreateUpdateProduct(p models.Product, update bool) (int64, error)

	GetStorages(zmqclient.RequestFilter, int) ([]models.Storage, int, error)
	GetOtherStorages(zmqclient.RequestFilter, int) ([]models.Entity, int, error)
	GetStorage(id int) (models.Storage, error)
	GetStorageEntity(id int) (models.Entity, error)
	DeleteStorage(id int) error
	ArchiveStorage(id int) error
	RestoreStorage(id int) error
	CreateUpdateStorage(s models.Storage, itemNumber int, update bool) (int64, error)
	UpdateAllQRCodes() error

	DeleteEntity(id int) error
	CreateEntity(e models.Entity) (int64, error)
	UpdateEntity(e models.Entity) error

	CreatePerson(p models.Person) (int64, error)
	UpdatePerson(p models.Person) error
	DeletePerson(id int) error
	GetAdmins() ([]models.Person, error)
	IsPersonAdmin(id int) (bool, error)
	UnsetPersonAdmin(id int) error
	SetPersonAdmin(id int) error
	HasPersonReadRestrictedProductPermission(id int) (bool, error)
}

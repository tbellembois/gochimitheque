package datastores

import (
	"github.com/jmoiron/sqlx"
	"github.com/tbellembois/gochimitheque/models"
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

	DeleteStorage(id int) error
	ArchiveStorage(id int) error
	RestoreStorage(id int) error
	UpdateAllQRCodes() error

	DeleteEntity(id int) error
	CreateEntity(e models.Entity) (int64, error)
	UpdateEntity(e models.Entity) error

	CreatePerson(p models.Person) (int64, error)
	UpdatePerson(p models.Person) error
	DeletePerson(id int) error
	GetAdmins() ([]models.Person, error)
	IsPersonAdmin(id int) (bool, error)
	UnsetPersonAdmin(id int64) error
	SetPersonAdmin(id int64) error
}

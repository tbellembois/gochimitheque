package datastores

import (
	"github.com/tbellembois/gochimitheque/models"
)

// Datastore is an interface to be implemented
// to store data.
type Datastore interface {
	ToCasbinJSONAdapter() ([]byte, error)

	GetWelcomeAnnounce() (models.WelcomeAnnounce, error)
	UpdateWelcomeAnnounce(w models.WelcomeAnnounce) error
}

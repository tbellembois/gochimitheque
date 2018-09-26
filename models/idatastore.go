package models

import (
	"net/http"
	"strconv"

	"github.com/tbellembois/gochimitheque/constants"
	"github.com/tbellembois/gopicluster/models"
)

// selectParameters contains the common parameters
// passed to the database Get* functions returning multiple values
// such as GetStoreLocations, GetEntities...
type selectParameters struct {
	LoggedPersonID int
	Search         string
	Order          string
	Offset         uint64
	Limit          uint64
}

// NewSelectParameters returns a getCommonParameters struct
// with default values
func NewSelectParameters() selectParameters {
	return selectParameters{
		LoggedPersonID: 0,
		Search:         "%%",
		Order:          "asc",
		Offset:         0,
		Limit:          constants.MaxUint64,
	}
}

// NewSelectParametersFromRequest returns a getCommonParameters struct
// with values in the request r except LoggedPersonID
func NewSelectParametersFromRequest(r *http.Request) (selectParameters, *models.AppError) {
	var err error

	cp := NewSelectParameters()

	if s, ok := r.URL.Query()["search"]; ok {
		cp.Search = "%" + s[0] + "%"
	}
	if o, ok := r.URL.Query()["order"]; ok {
		cp.Order = o[0]
	}
	if o, ok := r.URL.Query()["offset"]; ok {
		var of int
		if of, err = strconv.Atoi(o[0]); err != nil {
			return selectParameters{}, &models.AppError{
				Error:   err,
				Code:    http.StatusInternalServerError,
				Message: "offset atoi conversion",
			}
		}
		cp.Offset = uint64(of)
	}
	if l, ok := r.URL.Query()["limit"]; ok {
		var lm int
		if lm, err = strconv.Atoi(l[0]); err != nil {
			return selectParameters{}, &models.AppError{
				Error:   err,
				Code:    http.StatusInternalServerError,
				Message: "limit atoi conversion",
			}
		}
		cp.Limit = uint64(lm)
	}
	return cp, nil
}

// GetPeopleParameters contains the parameters of the GetPeople function
type GetPeopleParameters struct {
	CP       selectParameters
	EntityID int
}

// GetEntitiesParameters contains the parameters of the GetEntities function
type GetEntitiesParameters struct {
	CP selectParameters
}

// GetStoreLocationsParameters contains the parameters of the GetStoreLocations function
type GetStoreLocationsParameters struct {
	CP       selectParameters
	EntityID int
}

// GetProductsParameters contains the parameters of the GetProducts function
type GetProductsParameters struct {
	CP       selectParameters
	EntityID int
}

// Datastore is an interface to be implemented
// to store data
type Datastore interface {
	FlushErrors() error
	CreateDatabase() error

	// products
	GetProducts(GetProductsParameters) ([]Product, int, error)
	GetProductsCasNumbers(selectParameters) ([]CasNumber, int, error)
	GetProductsNames(selectParameters) ([]Name, int, error)
	GetProductsSymbols(selectParameters) ([]Symbol, int, error)
	GetProduct(id int) (Product, error)
	DeleteProduct(id int) error
	CreateProduct(p Product) (error, int)
	UpdateProduct(p Product) error

	// store locations
	GetStoreLocations(GetStoreLocationsParameters) ([]StoreLocation, int, error)
	GetStoreLocation(id int) (StoreLocation, error)
	GetStoreLocationEntity(id int) (Entity, error)
	DeleteStoreLocation(id int) error
	CreateStoreLocation(s StoreLocation) (error, int)
	UpdateStoreLocation(s StoreLocation) error

	// entities
	GetEntities(GetEntitiesParameters) ([]Entity, int, error)
	GetEntity(id int) (Entity, error)
	GetEntityPeople(id int) ([]Person, error)
	DeleteEntity(id int) error
	CreateEntity(e Entity) (error, int)
	UpdateEntity(e Entity) error
	IsEntityEmpty(id int) (bool, error)

	// people
	GetPeople(GetPeopleParameters) ([]Person, int, error)
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

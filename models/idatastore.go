package models

import (
	"net/http"
	"strconv"

	"github.com/tbellembois/gochimitheque/constants"
	"github.com/tbellembois/gopicluster/models"
)

// getCommonParameters contains the common parameters
// passed to the database Get* functions returning multiple values
// such as GetStoreLocations, GetEntities...
type getCommonParameters struct {
	LoggedPersonID int
	Search         string
	Order          string
	Offset         uint64
	Limit          uint64
}

// NewGetCommonParameters returns a getCommonParameters struct
// with default values
func NewGetCommonParameters() getCommonParameters {
	return getCommonParameters{
		LoggedPersonID: 0,
		Search:         "%%",
		Order:          "asc",
		Offset:         0,
		Limit:          constants.MaxUint64,
	}
}

// NewGetCommonParametersFromRequest returns a getCommonParameters struct
// with values in the request r except LoggedPersonID
func NewGetCommonParametersFromRequest(r *http.Request) (getCommonParameters, *models.AppError) {
	var err error

	cp := getCommonParameters{}

	if s, ok := r.URL.Query()["search"]; ok {
		cp.Search = "%" + s[0] + "%"
	}
	if o, ok := r.URL.Query()["order"]; ok {
		cp.Order = o[0]
	}
	if o, ok := r.URL.Query()["offset"]; ok {
		var of int
		if of, err = strconv.Atoi(o[0]); err != nil {
			return getCommonParameters{}, &models.AppError{
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
			return getCommonParameters{}, &models.AppError{
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
	CP       getCommonParameters
	EntityID int
}

// GetEntitiesParameters contains the parameters of the GetEntities function
type GetEntitiesParameters struct {
	CP getCommonParameters
}

// GetStoreLocationsParameters contains the parameters of the GetStoreLocations function
type GetStoreLocationsParameters struct {
	CP       getCommonParameters
	EntityID int
}

// GetProductsParameters contains the parameters of the GetProducts function
type GetProductsParameters struct {
	CP       getCommonParameters
	EntityID int
}

// Datastore is an interface to be implemented
// to store data
type Datastore interface {
	FlushErrors() error
	CreateDatabase() error

	// products
	GetProducts(GetProductsParameters) ([]Product, int, error)
	GetProductsCasNumbers(getCommonParameters) ([]CasNumber, int, error)
	GetProductsNames(getCommonParameters) ([]Name, int, error)
	GetProductsSymbols(getCommonParameters) ([]Symbol, int, error)
	GetProduct(id int) (Product, error)
	DeleteProduct(id int) error
	CreateProduct(p Product) (error, int)
	UpdateProduct(p Product) error
	IsProductWithName(name string) (bool, error)

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
	IsPersonWithEmail(email string) (bool, error)
	IsPersonWithEmailExcept(string, ...string) (bool, error)
}

package models

import (
	"net/http"
	"strconv"

	"github.com/tbellembois/gochimitheque/constants"
)

// dbselectparam contains the common parameters
// of the db select requests
// such as in GetStoreLocations, GetEntities...
type Dbselectparam interface{}
type dbselectparam struct {
	LoggedPersonID int // logged person, used to filter results
	Search         string
	Order          string
	Offset         uint64
	Limit          uint64
}

// dbselectparamProduct contains the parameters of the GetProducts function
type DbselectparamProduct interface {
	Dbselectparam
}
type dbselectparamProduct struct {
	dbselectparam
	Entity    int // id
	CasNumber int // id
}

// dbselectparamPerson contains the parameters of the GetPeople function
type DbselectparamPerson interface {
	Dbselectparam
}
type dbselectparamPerson struct {
	dbselectparam
	Entity int
}

// dbselectparamEntities contains the parameters of the GetEntities function
type DbselectparamEntity interface {
	Dbselectparam
}
type dbselectparamEntity struct {
	dbselectparam
}

// dbselectparamStoreLocation contains the parameters of the GetStoreLocations function
type DbselectparamStoreLocation interface {
	Dbselectparam
}
type dbselectparamStoreLocation struct {
	dbselectparam
	Entity int
}

// Newdbselectparam returns a dbselectparam struct
// with values populated from the request parameters
func Newdbselectparam(r *http.Request) (dbselectparam, *AppError) {

	var err error

	// initializing default values
	dsp := dbselectparam{
		LoggedPersonID: 0,
		Search:         "%%",
		Order:          "asc",
		Offset:         0,
		Limit:          constants.MaxUint64,
	}

	// returning default values
	if r == nil {
		return dsp, nil
	}

	// populating with request values
	if s, ok := r.URL.Query()["search"]; ok {
		dsp.Search = "%" + s[0] + "%"
	}
	if o, ok := r.URL.Query()["order"]; ok {
		dsp.Order = o[0]
	}
	if o, ok := r.URL.Query()["offset"]; ok {
		var of int
		if of, err = strconv.Atoi(o[0]); err != nil {
			return dbselectparam{}, &AppError{
				Error:   err,
				Code:    http.StatusInternalServerError,
				Message: "offset atoi conversion",
			}
		}
		dsp.Offset = uint64(of)
	}
	if l, ok := r.URL.Query()["limit"]; ok {
		var lm int
		if lm, err = strconv.Atoi(l[0]); err != nil {
			return dbselectparam{}, &AppError{
				Error:   err,
				Code:    http.StatusInternalServerError,
				Message: "limit atoi conversion",
			}
		}
		dsp.Limit = uint64(lm)
	}
	return dsp, nil

}

// NewdbselectparamProduct returns a dbselectparamProduct struct
// with values populated from the request parameters
func NewdbselectparamProduct(r *http.Request) (dbselectparamProduct, *AppError) {

	var (
		err  error
		aerr *AppError
		dsp  dbselectparam
		dspp dbselectparamProduct
	)

	// returning default values if no request
	if r == nil {
		if dsp, aerr = Newdbselectparam(nil); aerr != nil {
			return dbselectparamProduct{}, aerr
		}
		dspp.dbselectparam = dsp
		return dspp, nil
	}

	// or populating with request values
	if dsp, aerr = Newdbselectparam(r); err != nil {
		return dbselectparamProduct{}, aerr
	}
	dspp.dbselectparam = dsp

	if entityid, ok := r.URL.Query()["entity"]; ok {
		var eid int
		if eid, err = strconv.Atoi(entityid[0]); err != nil {
			return dbselectparamProduct{}, &AppError{
				Error:   err,
				Code:    http.StatusInternalServerError,
				Message: "limit atoi conversion",
			}
		}
		dspp.Entity = eid
	}
	return dspp, nil

}

// NewdbselectparamStoreLocation returns a dbselectparamStoreLocation struct
// with values populated from the request parameters
func NewdbselectparamStoreLocation(r *http.Request) (dbselectparamStoreLocation, *AppError) {

	var (
		err   error
		aerr  *AppError
		dsp   dbselectparam
		dspsl dbselectparamStoreLocation
	)

	// returning default values if no request
	if r == nil {
		if dsp, aerr = Newdbselectparam(nil); aerr != nil {
			return dbselectparamStoreLocation{}, aerr
		}
		dspsl.dbselectparam = dsp
		return dspsl, nil
	}

	// or populating with request values
	if dsp, aerr = Newdbselectparam(r); err != nil {
		return dbselectparamStoreLocation{}, aerr
	}
	dspsl.dbselectparam = dsp

	if entityid, ok := r.URL.Query()["entity"]; ok {
		var eid int
		if eid, err = strconv.Atoi(entityid[0]); err != nil {
			return dbselectparamStoreLocation{}, &AppError{
				Error:   err,
				Code:    http.StatusInternalServerError,
				Message: "limit atoi conversion",
			}
		}
		dspsl.Entity = eid
	}
	return dspsl, nil

}

// NewdbselectparamEntity returns a dbselectparamEntity struct
// with values populated from the request parameters
func NewdbselectparamEntity(r *http.Request) (dbselectparamEntity, *AppError) {

	var (
		err  error
		aerr *AppError
		dsp  dbselectparam
		dspe dbselectparamEntity
	)

	// returning default values if no request
	if r == nil {
		if dsp, aerr = Newdbselectparam(nil); aerr != nil {
			return dbselectparamEntity{}, aerr
		}
		dspe.dbselectparam = dsp
		return dspe, nil
	}

	// or populating with request values
	if dsp, aerr = Newdbselectparam(r); err != nil {
		return dbselectparamEntity{}, aerr
	}
	dspe.dbselectparam = dsp
	return dspe, nil

}

// Datastore is an interface to be implemented
// to store data
type Datastore interface {
	FlushErrors() error
	CreateDatabase() error

	// products
	GetProducts(dbselectparamProduct) ([]Product, int, error)
	GetProductsCasNumbers(dbselectparam) ([]CasNumber, int, error)
	GetProductsNames(dbselectparam) ([]Name, int, error)
	GetProductsSymbols(dbselectparam) ([]Symbol, int, error)
	GetProduct(id int) (Product, error)
	DeleteProduct(id int) error
	CreateProduct(p Product) (error, int)
	UpdateProduct(p Product) error

	// store locations
	GetStoreLocations(dbselectparamStoreLocation) ([]StoreLocation, int, error)
	GetStoreLocation(id int) (StoreLocation, error)
	GetStoreLocationEntity(id int) (Entity, error)
	DeleteStoreLocation(id int) error
	CreateStoreLocation(s StoreLocation) (error, int)
	UpdateStoreLocation(s StoreLocation) error

	// entities
	GetEntities(dbselectparamEntity) ([]Entity, int, error)
	GetEntity(id int) (Entity, error)
	GetEntityPeople(id int) ([]Person, error)
	DeleteEntity(id int) error
	CreateEntity(e Entity) (error, int)
	UpdateEntity(e Entity) error
	IsEntityEmpty(id int) (bool, error)

	// people
	GetPeople(dbselectparamPerson) ([]Person, int, error)
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

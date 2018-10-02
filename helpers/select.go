package helpers

import (
	"github.com/tbellembois/gochimitheque/constants"
	"net/http"
	"strconv"
)

// dbselectparam contains the common parameters
// of the db select requests
// such as in GetStoreLocations, GetEntities...
type Dbselectparam interface {
	SetLoggedPersonID(int)
	SetSearch(string)

	GetLoggedPersonID() int
	GetSearch() string
	GetOrder() string
	GetOrderBy() string
	GetOffset() uint64
	GetLimit() uint64
}
type dbselectparam struct {
	LoggedPersonID int // logged person, used to filter results
	Search         string
	OrderBy        string
	Order          string
	Offset         uint64
	Limit          uint64
}

// dbselectparamProduct contains the parameters of the GetProducts function
type DbselectparamProduct interface {
	Dbselectparam
	SetEntity(int)
	SetProduct(int)

	GetEntity() int
	GetProduct() int
}
type dbselectparamProduct struct {
	dbselectparam
	Entity    int // id
	Product   int // id
	CasNumber int // id
}

// dbselectparamStorage contains the parameters of the GetStorages function
type DbselectparamStorage interface {
	Dbselectparam
	SetEntity(int)
	SetProduct(int)

	GetEntity() int
	GetProduct() int
}
type dbselectparamStorage struct {
	dbselectparam
	Entity  int // id
	Product int // id
}

// dbselectparamPerson contains the parameters of the GetPeople function
type DbselectparamPerson interface {
	Dbselectparam
	SetEntity(int)

	GetEntity() int
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
	SetEntity(int)

	GetEntity() int
}
type dbselectparamStoreLocation struct {
	dbselectparam
	Entity int
}

//
// dbselectparam functions
//
func (d *dbselectparam) SetLoggedPersonID(i int) {
	d.LoggedPersonID = i
}

func (d *dbselectparam) SetSearch(s string) {
	d.Search = s
}

func (d dbselectparam) GetLoggedPersonID() int {
	return d.LoggedPersonID
}

func (d dbselectparam) GetSearch() string {
	return d.Search
}

func (d dbselectparam) GetOrder() string {
	return d.Order
}

func (d dbselectparam) GetOrderBy() string {
	return d.OrderBy
}

func (d dbselectparam) GetOffset() uint64 {
	return d.Offset
}

func (d dbselectparam) GetLimit() uint64 {
	return d.Limit
}

//
// dbselectparamPerson functions
//
func (d *dbselectparamPerson) SetEntity(i int) {
	d.Entity = i
}

func (d dbselectparamPerson) GetEntity() int {
	return d.Entity
}

//
// dbselectparamStoreLocation functions
//
func (d *dbselectparamStoreLocation) SetEntity(i int) {
	d.Entity = i
}

func (d *dbselectparamStoreLocation) GetEntity() int {
	return d.Entity
}

//
// dbselectparamProduct functions
//
func (d *dbselectparamProduct) SetEntity(i int) {
	d.Entity = i
}

func (d dbselectparamProduct) GetEntity() int {
	return d.Entity
}

func (d *dbselectparamProduct) SetProduct(i int) {
	d.Product = i
}

func (d dbselectparamProduct) GetProduct() int {
	return d.Product
}

//
// dbselectparamStorage functions
//
func (d *dbselectparamStorage) SetEntity(i int) {
	d.Entity = i
}

func (d dbselectparamStorage) GetEntity() int {
	return d.Entity
}

func (d *dbselectparamStorage) SetProduct(i int) {
	d.Product = i
}

func (d dbselectparamStorage) GetProduct() int {
	return d.Product
}

// Newdbselectparam returns a dbselectparam struct
// with values populated from the request parameters
func Newdbselectparam(r *http.Request) (*dbselectparam, *AppError) {

	var err error

	// initializing default values
	dsp := dbselectparam{
		LoggedPersonID: 0,
		Search:         "%%",
		OrderBy:        "",
		Order:          "asc",
		Offset:         0,
		Limit:          constants.MaxUint64,
	}

	// returning default values
	if r == nil {
		return &dsp, nil
	}

	// retrieving the logged user id from request context
	c := ContainerFromRequestContext(r)
	dsp.LoggedPersonID = c.PersonID

	// populating with request values
	if s, ok := r.URL.Query()["search"]; ok {
		dsp.Search = "%" + s[0] + "%"
	}
	if o, ok := r.URL.Query()["order"]; ok {
		dsp.Order = o[0]
	}
	if o, ok := r.URL.Query()["sort"]; ok {
		dsp.OrderBy = o[0]
	}
	if o, ok := r.URL.Query()["offset"]; ok {
		var of int
		if of, err = strconv.Atoi(o[0]); err != nil {
			return nil, &AppError{
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
			return nil, &AppError{
				Error:   err,
				Code:    http.StatusInternalServerError,
				Message: "limit atoi conversion",
			}
		}
		dsp.Limit = uint64(lm)
	}
	return &dsp, nil

}

// NewdbselectparamProduct returns a dbselectparamProduct struct
// with values populated from the request parameters
func NewdbselectparamProduct(r *http.Request) (*dbselectparamProduct, *AppError) {

	var (
		err  error
		aerr *AppError
		dsp  *dbselectparam
		dspp dbselectparamProduct
	)

	// init defaults
	dspp.Entity = -1
	dspp.Product = -1
	if dsp, aerr = Newdbselectparam(r); aerr != nil {
		return nil, aerr
	}
	dspp.dbselectparam = *dsp

	if r != nil {
		if o, ok := r.URL.Query()["sort"]; ok {
			dspp.OrderBy = o[0]
		} else {
			dspp.OrderBy = "product_id"
		}
		if entityid, ok := r.URL.Query()["entity"]; ok {
			var eid int
			if eid, err = strconv.Atoi(entityid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "limit atoi conversion",
				}
			}
			dspp.Entity = eid
		}
		if productid, ok := r.URL.Query()["product"]; ok {
			var pid int
			if pid, err = strconv.Atoi(productid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "limit atoi conversion",
				}
			}
			dspp.Product = pid
		}
	}
	return &dspp, nil

}

// NewdbselectparamStorage returns a dbselectparamStorage struct
// with values populated from the request parameters
func NewdbselectparamStorage(r *http.Request) (*dbselectparamStorage, *AppError) {

	var (
		err  error
		aerr *AppError
		dsp  *dbselectparam
		dsps dbselectparamStorage
	)

	// init defaults
	dsps.Entity = -1
	dsps.Product = -1
	if dsp, aerr = Newdbselectparam(r); aerr != nil {
		return nil, aerr
	}
	dsps.dbselectparam = *dsp

	if r != nil {
		if o, ok := r.URL.Query()["sort"]; ok {
			switch o[0] {
			case "product.name.name_label":
				dsps.OrderBy = "name.name_label"
			default:
				dsps.OrderBy = o[0]
			}
		} else {
			dsps.OrderBy = "storage_id"
		}
		if entityid, ok := r.URL.Query()["entity"]; ok {
			var eid int
			if eid, err = strconv.Atoi(entityid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "limit atoi conversion",
				}
			}
			dsps.Entity = eid
		}
		if productid, ok := r.URL.Query()["product"]; ok {
			var pid int
			if pid, err = strconv.Atoi(productid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "limit atoi conversion",
				}
			}
			dsps.Product = pid
		}
	}
	return &dsps, nil

}

// NewdbselectparamStoreLocation returns a dbselectparamStoreLocation struct
// with values populated from the request parameters
func NewdbselectparamStoreLocation(r *http.Request) (*dbselectparamStoreLocation, *AppError) {

	var (
		err   error
		aerr  *AppError
		dsp   *dbselectparam
		dspsl dbselectparamStoreLocation
	)

	// init defaults
	dspsl.Entity = -1
	if dsp, aerr = Newdbselectparam(r); aerr != nil {
		return nil, aerr
	}
	dspsl.dbselectparam = *dsp

	if r != nil {
		if o, ok := r.URL.Query()["sort"]; ok {
			dspsl.OrderBy = o[0]
		} else {
			dspsl.OrderBy = "storelocation_id"
		}
		if entityid, ok := r.URL.Query()["entity"]; ok {
			var eid int
			if eid, err = strconv.Atoi(entityid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "limit atoi conversion",
				}
			}
			dspsl.Entity = eid
		}
	}
	return &dspsl, nil

}

// NewdbselectparamPerson returns a dbselectparamStorePerson struct
// with values populated from the request parameters
func NewdbselectparamPerson(r *http.Request) (*dbselectparamPerson, *AppError) {

	var (
		err  error
		aerr *AppError
		dsp  *dbselectparam
		dspp dbselectparamPerson
	)

	// init defaults
	dspp.Entity = -1
	if dsp, aerr = Newdbselectparam(r); aerr != nil {
		return nil, aerr
	}
	dspp.dbselectparam = *dsp

	if r != nil {
		if o, ok := r.URL.Query()["sort"]; ok {
			dspp.OrderBy = o[0]
		} else {
			dspp.OrderBy = "person_id"
		}
		if entityid, ok := r.URL.Query()["entity"]; ok {
			var eid int
			if eid, err = strconv.Atoi(entityid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "limit atoi conversion",
				}
			}
			dspp.Entity = eid
		}
	}
	return &dspp, nil

}

// NewdbselectparamEntity returns a dbselectparamEntity struct
// with values populated from the request parameters
func NewdbselectparamEntity(r *http.Request) (*dbselectparamEntity, *AppError) {

	var (
		aerr *AppError
		dsp  *dbselectparam
		dspe dbselectparamEntity
	)

	// init defaults
	if dsp, aerr = Newdbselectparam(r); aerr != nil {
		return nil, aerr
	}
	dspe.dbselectparam = *dsp
	dspe.OrderBy = "entity_id"

	return &dspe, nil

}

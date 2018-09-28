package helpers

import (
	"net/http"
	"strconv"

	"github.com/tbellembois/gochimitheque/constants"
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
	GetOffset() uint64
	GetLimit() uint64
}
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
	SetEntity(int)

	GetEntity() int
}
type dbselectparamProduct struct {
	dbselectparam
	Entity    int // id
	CasNumber int // id
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

// Newdbselectparam returns a dbselectparam struct
// with values populated from the request parameters
func Newdbselectparam(r *http.Request) (*dbselectparam, *AppError) {

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

	// returning default values if no request
	if r == nil {
		if dsp, aerr = Newdbselectparam(nil); aerr != nil {
			return nil, aerr
		}
		dspp.dbselectparam = *dsp
		dspp.Entity = -1
		return &dspp, nil
	}

	// or populating with request values
	if dsp, aerr = Newdbselectparam(r); err != nil {
		return nil, aerr
	}
	dspp.dbselectparam = *dsp

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
	return &dspp, nil

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

	// returning default values if no request
	if r == nil {
		if dsp, aerr = Newdbselectparam(nil); aerr != nil {
			return nil, aerr
		}
		dspsl.dbselectparam = *dsp
		dspsl.Entity = -1
		return &dspsl, nil
	}

	// or populating with request values
	if dsp, aerr = Newdbselectparam(r); err != nil {
		return nil, aerr
	}
	dspsl.dbselectparam = *dsp

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

	// returning default values if no request
	if r == nil {
		if dsp, aerr = Newdbselectparam(nil); aerr != nil {
			return nil, aerr
		}
		dspp.dbselectparam = *dsp
		dspp.Entity = -1
		return &dspp, nil
	}

	// or populating with request values
	if dsp, aerr = Newdbselectparam(r); err != nil {
		return nil, aerr
	}
	dspp.dbselectparam = *dsp

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
	return &dspp, nil

}

// NewdbselectparamEntity returns a dbselectparamEntity struct
// with values populated from the request parameters
func NewdbselectparamEntity(r *http.Request) (*dbselectparamEntity, *AppError) {

	var (
		err  error
		aerr *AppError
		dsp  *dbselectparam
		dspe dbselectparamEntity
	)

	// returning default values if no request
	if r == nil {
		if dsp, aerr = Newdbselectparam(nil); aerr != nil {
			return nil, aerr
		}
		dspe.dbselectparam = *dsp
		return &dspe, nil
	}

	// or populating with request values
	if dsp, aerr = Newdbselectparam(r); err != nil {
		return nil, aerr
	}
	dspe.dbselectparam = *dsp
	return &dspe, nil

}

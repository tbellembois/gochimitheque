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
	GetOrderBy() string
	GetOffset() uint64
	GetLimit() uint64

	SetOrderBy(string)
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
	SetStorelocation(int)
	SetBookmark(bool)
	SetName(int)
	SetEmpiricalFormula(int)
	SetCasNumber(int)
	SetStorageBarecode(string)

	SetCustomNamePartOf(string)

	GetEntity() int
	GetProduct() int
	GetStorelocation() int
	GetBookmark() bool
	GetName() int
	GetEmpiricalFormula() int
	GetCasNumber() int
	GetStorageBarecode() string

	GetCustomNamePartOf() string
}
type dbselectparamProduct struct {
	dbselectparam
	Entity           int // id
	Product          int // id
	Storelocation    int // id
	Bookmark         bool
	Name             int // id
	EmpiricalFormula int // id
	CasNumber        int // id
	StorageBarecode  string

	// custom search
	CustomNamePartOf string
}

// dbselectparamStorage contains the parameters of the GetStorages function
type DbselectparamStorage interface {
	Dbselectparam
	SetEntity(int)
	SetProduct(int)
	SetStorelocation(int)
	SetStorage(int)
	SetHistory(bool)
	SetStorageArchive(bool)

	GetEntity() int
	GetProduct() int
	GetStorelocation() int
	GetStorage() int
	GetHistory() bool
	GetStorageArchive() bool
}
type dbselectparamStorage struct {
	dbselectparam
	Entity         int // id
	Product        int // id
	Storelocation  int // id
	Storage        int // id
	History        bool
	StorageArchive bool
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
	SetStoreLocationCanStore(bool)

	GetEntity() int
	GetStoreLocationCanStore() bool
}
type dbselectparamStoreLocation struct {
	dbselectparam
	Entity                int
	StoreLocationCanStore bool
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

func (d *dbselectparam) SetOrderBy(o string) {
	d.OrderBy = o
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

func (d *dbselectparamStoreLocation) GetStoreLocationCanStore() bool {
	return d.StoreLocationCanStore
}

func (d *dbselectparamStoreLocation) SetStoreLocationCanStore(b bool) {
	d.StoreLocationCanStore = b
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

func (d *dbselectparamProduct) SetStorelocation(i int) {
	d.Storelocation = i
}

func (d dbselectparamProduct) GetStorelocation() int {
	return d.Storelocation
}

func (d *dbselectparamProduct) SetBookmark(b bool) {
	d.Bookmark = b
}

func (d dbselectparamProduct) GetBookmark() bool {
	return d.Bookmark
}

func (d *dbselectparamProduct) SetName(n int) {
	d.Name = n
}

func (d dbselectparamProduct) GetName() int {
	return d.Name
}

func (d *dbselectparamProduct) SetEmpiricalFormula(n int) {
	d.Name = n
}

func (d dbselectparamProduct) GetEmpiricalFormula() int {
	return d.EmpiricalFormula
}

func (d *dbselectparamProduct) SetCasNumber(n int) {
	d.CasNumber = n
}

func (d dbselectparamProduct) GetCasNumber() int {
	return d.CasNumber
}

func (d *dbselectparamProduct) SetStorageBarecode(n string) {
	d.StorageBarecode = n
}

func (d dbselectparamProduct) GetStorageBarecode() string {
	return d.StorageBarecode
}

func (d *dbselectparamProduct) SetCustomNamePartOf(n string) {
	d.CustomNamePartOf = n
}

func (d dbselectparamProduct) GetCustomNamePartOf() string {
	return d.CustomNamePartOf
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

func (d *dbselectparamStorage) SetStorelocation(i int) {
	d.Storelocation = i
}

func (d dbselectparamStorage) GetStorelocation() int {
	return d.Storelocation
}

func (d *dbselectparamStorage) SetStorage(i int) {
	d.Storage = i
}

func (d dbselectparamStorage) GetStorage() int {
	return d.Storage
}

func (d dbselectparamStorage) GetHistory() bool {
	return d.History
}

func (d *dbselectparamStorage) SetHistory(b bool) {
	d.History = b
}

func (d dbselectparamStorage) GetStorageArchive() bool {
	return d.StorageArchive
}

func (d *dbselectparamStorage) SetStorageArchive(b bool) {
	d.StorageArchive = b
}

// Newdbselectparam returns a dbselectparam struct
// with values populated from the request parameters
func Newdbselectparam(r *http.Request, f func(string) (string, error)) (*dbselectparam, *AppError) {

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
		if f != nil {
			fs, _ := f(s[0])
			dsp.Search = "%" + fs + "%"
		} else {
			dsp.Search = "%" + s[0] + "%"
		}
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
func NewdbselectparamProduct(r *http.Request, f func(string) (string, error)) (*dbselectparamProduct, *AppError) {

	var (
		err  error
		aerr *AppError
		dsp  *dbselectparam
		dspp dbselectparamProduct
	)

	// init defaults
	dspp.Entity = -1
	dspp.Product = -1
	dspp.Storelocation = -1
	dspp.Bookmark = false
	dspp.Name = -1
	dspp.CasNumber = -1
	dspp.EmpiricalFormula = -1
	dspp.StorageBarecode = ""
	dspp.CustomNamePartOf = ""
	if dsp, aerr = Newdbselectparam(r, f); aerr != nil {
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
			if dspp.Entity, err = strconv.Atoi(entityid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "entity atoi conversion",
				}
			}
		}
		if productid, ok := r.URL.Query()["product"]; ok {
			if dspp.Product, err = strconv.Atoi(productid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "product atoi conversion",
				}
			}
		}
		if storelocationid, ok := r.URL.Query()["storelocation"]; ok {
			if dspp.Storelocation, err = strconv.Atoi(storelocationid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "storelocation atoi conversion",
				}
			}
		}
		if bookmark, ok := r.URL.Query()["bookmark"]; ok {
			if dspp.Bookmark, err = strconv.ParseBool(bookmark[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "bookmark bool conversion",
				}
			}
		}
		if nameid, ok := r.URL.Query()["name"]; ok {
			if dspp.Name, err = strconv.Atoi(nameid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "name atoi conversion",
				}
			}
		}
		if casnumberid, ok := r.URL.Query()["casnumber"]; ok {
			if dspp.CasNumber, err = strconv.Atoi(casnumberid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "casnumber atoi conversion",
				}
			}
		}
		if empiricalformulaid, ok := r.URL.Query()["empiricalformula"]; ok {
			if dspp.EmpiricalFormula, err = strconv.Atoi(empiricalformulaid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "empiricalformula atoi conversion",
				}
			}
		}
		if storage_barecode, ok := r.URL.Query()["storage_barecode"]; ok {
			dspp.StorageBarecode = storage_barecode[0]
		}
		if custom_name_part_of, ok := r.URL.Query()["custom_name_part_of"]; ok {
			dspp.CustomNamePartOf = custom_name_part_of[0]
		}
	}
	return &dspp, nil

}

// NewdbselectparamStorage returns a dbselectparamStorage struct
// with values populated from the request parameters
func NewdbselectparamStorage(r *http.Request, f func(string) (string, error)) (*dbselectparamStorage, *AppError) {

	var (
		err  error
		aerr *AppError
		dsp  *dbselectparam
		dsps dbselectparamStorage
	)

	// init defaults
	dsps.Entity = -1
	dsps.Product = -1
	dsps.Storelocation = -1
	dsps.Storage = -1
	dsps.History = false
	dsps.StorageArchive = false
	if dsp, aerr = Newdbselectparam(r, f); aerr != nil {
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
			if dsps.Entity, err = strconv.Atoi(entityid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "entity atoi conversion",
				}
			}
		}
		if productid, ok := r.URL.Query()["product"]; ok {
			if dsps.Product, err = strconv.Atoi(productid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "product atoi conversion",
				}
			}
		}
		if storelocationid, ok := r.URL.Query()["storelocation"]; ok {
			if dsps.Storelocation, err = strconv.Atoi(storelocationid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "storelocation atoi conversion",
				}
			}
		}
		if storageid, ok := r.URL.Query()["storage"]; ok {
			if dsps.Storage, err = strconv.Atoi(storageid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "storage atoi conversion",
				}
			}
		}
		if history, ok := r.URL.Query()["history"]; ok {
			if dsps.History, err = strconv.ParseBool(history[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "history bool conversion",
				}
			}
		}
		if storage_archive, ok := r.URL.Query()["storage_archive"]; ok {
			if dsps.StorageArchive, err = strconv.ParseBool(storage_archive[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "storage_archive bool conversion",
				}
			}
		}
	}
	return &dsps, nil

}

// NewdbselectparamStoreLocation returns a dbselectparamStoreLocation struct
// with values populated from the request parameters
func NewdbselectparamStoreLocation(r *http.Request, f func(string) (string, error)) (*dbselectparamStoreLocation, *AppError) {

	var (
		err   error
		aerr  *AppError
		dsp   *dbselectparam
		dspsl dbselectparamStoreLocation
	)

	// init defaults
	dspsl.Entity = -1
	if dsp, aerr = Newdbselectparam(r, f); aerr != nil {
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
			if dspsl.Entity, err = strconv.Atoi(entityid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "limit atoi conversion",
				}
			}
		}
		if c, ok := r.URL.Query()["storelocation_canstore"]; ok {
			if dspsl.StoreLocationCanStore, err = strconv.ParseBool(c[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "storelocation_canstore bool conversion",
				}
			}
		}
	}
	return &dspsl, nil

}

// NewdbselectparamPerson returns a dbselectparamStorePerson struct
// with values populated from the request parameters
func NewdbselectparamPerson(r *http.Request, f func(string) (string, error)) (*dbselectparamPerson, *AppError) {

	var (
		err  error
		aerr *AppError
		dsp  *dbselectparam
		dspp dbselectparamPerson
	)

	// init defaults
	dspp.Entity = -1
	if dsp, aerr = Newdbselectparam(r, f); aerr != nil {
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
			if dspp.Entity, err = strconv.Atoi(entityid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "limit atoi conversion",
				}
			}
		}
	}
	return &dspp, nil

}

// NewdbselectparamEntity returns a dbselectparamEntity struct
// with values populated from the request parameters
func NewdbselectparamEntity(r *http.Request, f func(string) (string, error)) (*dbselectparamEntity, *AppError) {

	var (
		aerr *AppError
		dsp  *dbselectparam
		dspe dbselectparamEntity
	)

	// init defaults
	if dsp, aerr = Newdbselectparam(r, f); aerr != nil {
		return nil, aerr
	}
	dspe.dbselectparam = *dsp
	dspe.OrderBy = "entity_id"

	if r != nil {
		if o, ok := r.URL.Query()["sort"]; ok {
			dspe.OrderBy = o[0]
		} else {
			dspe.OrderBy = "entity_id"
		}
	}

	return &dspe, nil

}

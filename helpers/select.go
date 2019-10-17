package helpers

import (
	"net/http"
	"strconv"

	"github.com/tbellembois/gochimitheque/constants"
)

// Dbselectparam contains the common parameters
// of the db select requests
// such as in GetStoreLocations, GetEntities...
type Dbselectparam interface {
	SetLoggedPersonID(int)
	SetSearch(string)
	SetLimit(uint64)

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
	SetProductSpecificity(string)

	SetCustomNamePartOf(string)
	SetName(int)
	SetEmpiricalFormula(int)
	SetCasNumber(int)
	SetStorageBarecode(string)
	SetSymbols([]int)
	SetHazardStatements([]int)
	SetPrecautionaryStatements([]int)
	SetSignalWord(int)
	SetCasNumberCmr(bool)

	GetEntity() int
	GetProduct() int
	GetStorelocation() int
	GetBookmark() bool
	GetProductSpecificity() string

	GetCustomNamePartOf() string
	GetName() int
	GetEmpiricalFormula() int
	GetCasNumber() int
	GetStorageBarecode() string
	GetSymbols() []int
	GetHazardStatements() []int
	GetPrecautionaryStatements() []int
	GetSignalWord() int
	GetCasNumberCmr() bool
}
type dbselectparamProduct struct {
	dbselectparam
	Entity        int // id
	Product       int // id
	Storelocation int // id
	Bookmark      bool

	CustomNamePartOf        string
	Name                    int // id
	EmpiricalFormula        int // id
	CasNumber               int // id
	StorageBarecode         string
	Symbols                 []int // ids
	HazardStatements        []int //ids
	PrecautionaryStatements []int //ids
	SignalWord              int   // id
	CasNumberCmr            bool
	ProductSpecificity      string
}

// DbselectparamStorage contains the parameters of the GetStorages function
type DbselectparamStorage interface {
	Dbselectparam
	SetEntity(int)
	SetProduct(int)
	SetStorelocation(int)
	SetBookmark(bool)
	SetStorage(int)
	SetHistory(bool)
	SetStorageArchive(bool)

	SetCustomNamePartOf(string)
	SetName(int)
	SetEmpiricalFormula(int)
	SetCasNumber(int)
	SetStorageBarecode(string)
	SetSymbols([]int)
	SetHazardStatements([]int)
	SetPrecautionaryStatements([]int)
	SetSignalWord(int)
	SetCasNumberCmr(bool)

	GetEntity() int
	GetProduct() int
	GetStorelocation() int
	GetBookmark() bool
	GetStorage() int
	GetHistory() bool
	GetStorageArchive() bool

	GetCustomNamePartOf() string
	GetName() int
	GetEmpiricalFormula() int
	GetCasNumber() int
	GetStorageBarecode() string
	GetSymbols() []int
	GetHazardStatements() []int
	GetPrecautionaryStatements() []int
	GetSignalWord() int
	GetCasNumberCmr() bool
}
type dbselectparamStorage struct {
	dbselectparam
	Entity         int // id
	Product        int // id
	Storelocation  int // id
	Storage        int // id
	Bookmark       bool
	History        bool
	StorageArchive bool

	CustomNamePartOf        string
	Name                    int // id
	EmpiricalFormula        int // id
	CasNumber               int // id
	StorageBarecode         string
	Symbols                 []int // ids
	HazardStatements        []int //ids
	PrecautionaryStatements []int //ids
	SignalWord              int   // id
	CasNumberCmr            bool
}

// DbselectparamPerson contains the parameters of the GetPeople function
type DbselectparamPerson interface {
	Dbselectparam
	SetEntity(int)

	GetEntity() int
}
type dbselectparamPerson struct {
	dbselectparam
	Entity int
}

// DbselectparamEntity contains the parameters of the GetEntities function
type DbselectparamEntity interface {
	Dbselectparam
}
type dbselectparamEntity struct {
	dbselectparam
}

// DbselectparamStoreLocation contains the parameters of the GetStoreLocations function
type DbselectparamStoreLocation interface {
	Dbselectparam
	SetEntity(int)
	SetStoreLocationCanStore(bool)

	SetPermission(string)

	GetEntity() int
	GetStoreLocationCanStore() bool

	GetPermission() string
}
type dbselectparamStoreLocation struct {
	dbselectparam
	Entity                int
	StoreLocationCanStore bool

	Permission string
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

func (d *dbselectparam) SetLimit(l uint64) {
	d.Limit = l
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

func (d *dbselectparamStoreLocation) GetPermission() string {
	return d.Permission
}

func (d *dbselectparamStoreLocation) SetPermission(p string) {
	d.Permission = p
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

func (d *dbselectparamProduct) SetSymbols(n []int) {
	d.Symbols = n
}

func (d dbselectparamProduct) GetSymbols() []int {
	return d.Symbols
}

func (d *dbselectparamProduct) SetCustomNamePartOf(n string) {
	d.CustomNamePartOf = n
}

func (d dbselectparamProduct) GetCustomNamePartOf() string {
	return d.CustomNamePartOf
}

func (d *dbselectparamProduct) SetHazardStatements(n []int) {
	d.HazardStatements = n
}

func (d dbselectparamProduct) GetHazardStatements() []int {
	return d.HazardStatements
}

func (d *dbselectparamProduct) SetPrecautionaryStatements(n []int) {
	d.PrecautionaryStatements = n
}

func (d dbselectparamProduct) GetPrecautionaryStatements() []int {
	return d.PrecautionaryStatements
}

func (d *dbselectparamProduct) SetSignalWord(n int) {
	d.SignalWord = n
}

func (d dbselectparamProduct) GetSignalWord() int {
	return d.SignalWord
}

func (d *dbselectparamProduct) SetCasNumberCmr(n bool) {
	d.CasNumberCmr = n
}

func (d dbselectparamProduct) GetCasNumberCmr() bool {
	return d.CasNumberCmr
}

func (d *dbselectparamProduct) SetProductSpecificity(s string) {
	d.ProductSpecificity = s
}

func (d *dbselectparamProduct) GetProductSpecificity() string {
	return d.ProductSpecificity
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

func (d *dbselectparamStorage) SetBookmark(b bool) {
	d.Bookmark = b
}

func (d dbselectparamStorage) GetBookmark() bool {
	return d.Bookmark
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

func (d *dbselectparamStorage) SetName(n int) {
	d.Name = n
}

func (d dbselectparamStorage) GetName() int {
	return d.Name
}

func (d *dbselectparamStorage) SetEmpiricalFormula(n int) {
	d.Name = n
}

func (d dbselectparamStorage) GetEmpiricalFormula() int {
	return d.EmpiricalFormula
}

func (d *dbselectparamStorage) SetCasNumber(n int) {
	d.CasNumber = n
}

func (d dbselectparamStorage) GetCasNumber() int {
	return d.CasNumber
}

func (d *dbselectparamStorage) SetStorageBarecode(n string) {
	d.StorageBarecode = n
}

func (d dbselectparamStorage) GetStorageBarecode() string {
	return d.StorageBarecode
}

func (d *dbselectparamStorage) SetSymbols(n []int) {
	d.Symbols = n
}

func (d dbselectparamStorage) GetSymbols() []int {
	return d.Symbols
}

func (d *dbselectparamStorage) SetCustomNamePartOf(n string) {
	d.CustomNamePartOf = n
}

func (d dbselectparamStorage) GetCustomNamePartOf() string {
	return d.CustomNamePartOf
}

func (d *dbselectparamStorage) SetHazardStatements(n []int) {
	d.HazardStatements = n
}

func (d dbselectparamStorage) GetHazardStatements() []int {
	return d.HazardStatements
}

func (d *dbselectparamStorage) SetPrecautionaryStatements(n []int) {
	d.PrecautionaryStatements = n
}

func (d dbselectparamStorage) GetPrecautionaryStatements() []int {
	return d.PrecautionaryStatements
}

func (d *dbselectparamStorage) SetSignalWord(n int) {
	d.SignalWord = n
}

func (d dbselectparamStorage) GetSignalWord() int {
	return d.SignalWord
}

func (d *dbselectparamStorage) SetCasNumberCmr(n bool) {
	d.CasNumberCmr = n
}

func (d dbselectparamStorage) GetCasNumberCmr() bool {
	return d.CasNumberCmr
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
		if f != nil && s[0] != "" {
			fs, err := f(s[0])
			if err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "error calling f",
				}
			}
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
	// no limit on export
	if _, ok := r.URL.Query()["export"]; !ok {
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
	dspp.SignalWord = -1
	dspp.CasNumberCmr = false
	dspp.ProductSpecificity = ""
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
		if symbolsids, ok := r.URL.Query()["symbols[]"]; ok {
			var sint int
			for _, s := range symbolsids {
				if sint, err = strconv.Atoi(s); err != nil {
					return nil, &AppError{
						Error:   err,
						Code:    http.StatusInternalServerError,
						Message: "symbol atoi conversion",
					}
				}
				dspp.Symbols = append(dspp.Symbols, sint)
			}
		}
		if hsids, ok := r.URL.Query()["hazardstatements[]"]; ok {
			var sint int
			for _, s := range hsids {
				if sint, err = strconv.Atoi(s); err != nil {
					return nil, &AppError{
						Error:   err,
						Code:    http.StatusInternalServerError,
						Message: "hazardstatement atoi conversion",
					}
				}
				dspp.HazardStatements = append(dspp.HazardStatements, sint)
			}
		}
		if psids, ok := r.URL.Query()["precautionarystatements[]"]; ok {
			var sint int
			for _, s := range psids {
				if sint, err = strconv.Atoi(s); err != nil {
					return nil, &AppError{
						Error:   err,
						Code:    http.StatusInternalServerError,
						Message: "precautionarystatement atoi conversion",
					}
				}
				dspp.PrecautionaryStatements = append(dspp.PrecautionaryStatements, sint)
			}
		}
		if signalwordid, ok := r.URL.Query()["signalword"]; ok {
			if dspp.SignalWord, err = strconv.Atoi(signalwordid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "signalword atoi conversion",
				}
			}
		}
		if storage_barecode, ok := r.URL.Query()["storage_barecode"]; ok {
			dspp.StorageBarecode = "%" + storage_barecode[0] + "%"
		}
		if custom_name_part_of, ok := r.URL.Query()["custom_name_part_of"]; ok {
			dspp.CustomNamePartOf = "%" + custom_name_part_of[0] + "%"
		}
		if product_specificity, ok := r.URL.Query()["product_specificity"]; ok {
			dspp.ProductSpecificity = product_specificity[0]
		}
		if casnumber_cmr, ok := r.URL.Query()["casnumber_cmr"]; ok {
			if dspp.CasNumberCmr, err = strconv.ParseBool(casnumber_cmr[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "casnumber_cmr bool conversion",
				}
			}
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
	dsps.Bookmark = false
	dsps.History = false
	dsps.StorageArchive = false
	dsps.Name = -1
	dsps.CasNumber = -1
	dsps.EmpiricalFormula = -1
	dsps.StorageBarecode = ""
	dsps.CustomNamePartOf = ""
	dsps.SignalWord = -1
	dsps.CasNumberCmr = false
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
		if bookmark, ok := r.URL.Query()["bookmark"]; ok {
			if dsps.Bookmark, err = strconv.ParseBool(bookmark[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "bookmark bool conversion",
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
		if nameid, ok := r.URL.Query()["name"]; ok {
			if dsps.Name, err = strconv.Atoi(nameid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "name atoi conversion",
				}
			}
		}
		if casnumberid, ok := r.URL.Query()["casnumber"]; ok {
			if dsps.CasNumber, err = strconv.Atoi(casnumberid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "casnumber atoi conversion",
				}
			}
		}
		if empiricalformulaid, ok := r.URL.Query()["empiricalformula"]; ok {
			if dsps.EmpiricalFormula, err = strconv.Atoi(empiricalformulaid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "empiricalformula atoi conversion",
				}
			}
		}
		if symbolsids, ok := r.URL.Query()["symbols[]"]; ok {
			var sint int
			for _, s := range symbolsids {
				if sint, err = strconv.Atoi(s); err != nil {
					return nil, &AppError{
						Error:   err,
						Code:    http.StatusInternalServerError,
						Message: "symbol atoi conversion",
					}
				}
				dsps.Symbols = append(dsps.Symbols, sint)
			}
		}
		if hsids, ok := r.URL.Query()["hazardstatements[]"]; ok {
			var sint int
			for _, s := range hsids {
				if sint, err = strconv.Atoi(s); err != nil {
					return nil, &AppError{
						Error:   err,
						Code:    http.StatusInternalServerError,
						Message: "hazardstatement atoi conversion",
					}
				}
				dsps.HazardStatements = append(dsps.HazardStatements, sint)
			}
		}
		if psids, ok := r.URL.Query()["precautionarystatements[]"]; ok {
			var sint int
			for _, s := range psids {
				if sint, err = strconv.Atoi(s); err != nil {
					return nil, &AppError{
						Error:   err,
						Code:    http.StatusInternalServerError,
						Message: "precautionarystatement atoi conversion",
					}
				}
				dsps.PrecautionaryStatements = append(dsps.PrecautionaryStatements, sint)
			}
		}
		if signalwordid, ok := r.URL.Query()["signalword"]; ok {
			if dsps.SignalWord, err = strconv.Atoi(signalwordid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "signalword atoi conversion",
				}
			}
		}
		if storage_barecode, ok := r.URL.Query()["storage_barecode"]; ok {
			dsps.StorageBarecode = "%" + storage_barecode[0] + "%"
		}
		if custom_name_part_of, ok := r.URL.Query()["custom_name_part_of"]; ok {
			dsps.CustomNamePartOf = "%" + custom_name_part_of[0] + "%"
		}
		if casnumber_cmr, ok := r.URL.Query()["casnumber_cmr"]; ok {
			if dsps.CasNumberCmr, err = strconv.ParseBool(casnumber_cmr[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "casnumber_cmr bool conversion",
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
	dspsl.Permission = "r"
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
		if p, ok := r.URL.Query()["permission"]; ok {
			dspsl.Permission = p[0]
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

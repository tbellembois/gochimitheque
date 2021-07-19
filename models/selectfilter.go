package models

import (
	"net/http"
	"strconv"
)

// SelectFilter contains the common parameters
// of the db select requests
// such as in GetStoreLocations, GetEntities...
type SelectFilter interface {
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
type selectFilter struct {
	LoggedPersonID int // logged person, used to filter results
	Search         string
	OrderBy        string
	Order          string
	Offset         uint64
	Limit          uint64
}

type SelectFilterProduct struct {
	selectFilter
	Entity        int // id
	Product       int // id
	Storelocation int // id
	ProducerRef   int // id
	Bookmark      bool

	CustomNamePartOf        string
	Name                    int // id
	EmpiricalFormula        int // id
	CasNumber               int // id
	StorageBarecode         string
	StorageBatchNumber      string
	Symbols                 []int // ids
	HazardStatements        []int //ids
	PrecautionaryStatements []int //ids
	SignalWord              int   // id
	CasNumberCmr            bool
	ProductSpecificity      string
	Borrowing               bool
	StorageToDestroy        bool
	Tags                    []int
	Category                int

	ShowBio   bool
	ShowChem  bool
	ShowConsu bool
}

type SelectFilterStorage struct {
	selectFilter
	Ids            []int
	Entity         int // id
	Product        int // id
	Storelocation  int // id
	Storage        int // id
	ProducerRef    int // id
	Bookmark       bool
	History        bool
	StorageArchive bool

	CustomNamePartOf        string
	Name                    int // id
	EmpiricalFormula        int // id
	CasNumber               int // id
	StorageBarecode         string
	StorageBatchNumber      string
	Symbols                 []int // ids
	HazardStatements        []int //ids
	PrecautionaryStatements []int //ids
	SignalWord              int   // id
	CasNumberCmr            bool
	Borrowing               bool
	StorageToDestroy        bool
	Tags                    []int
	Category                int

	ShowBio   bool
	ShowChem  bool
	ShowConsu bool
}

type SelectFilterPerson struct {
	selectFilter
	Entity int
}

type SelectFilterEntity struct {
	selectFilter
}

type SelectFilterStoreLocation struct {
	selectFilter
	Entity                int
	StoreLocationCanStore bool

	Permission string
}

type SelectFilterUnit struct {
	selectFilter
	UnitType string
}

type SelectFilterProducerRef struct {
	selectFilter
	Producer int
}

type SelectFilterSupplierRef struct {
	selectFilter
	Supplier int
}

//
// dbselectparam functions
//
func (d *selectFilter) SetLoggedPersonID(i int) {
	d.LoggedPersonID = i
}

func (d *selectFilter) SetSearch(s string) {
	d.Search = s
}

func (d selectFilter) GetLoggedPersonID() int {
	return d.LoggedPersonID
}

func (d selectFilter) GetSearch() string {
	return d.Search
}

func (d selectFilter) GetOrder() string {
	return d.Order
}

func (d selectFilter) GetOrderBy() string {
	return d.OrderBy
}

func (d *selectFilter) SetOrderBy(o string) {
	d.OrderBy = o
}

func (d selectFilter) GetOffset() uint64 {
	return d.Offset
}

func (d selectFilter) GetLimit() uint64 {
	return d.Limit
}

func (d *selectFilter) SetLimit(l uint64) {
	d.Limit = l
}

//
// dbselectparamUnit functions
//
func (d *SelectFilterUnit) SetUnitType(s string) {
	d.UnitType = s
}

func (d SelectFilterUnit) GetUnitType() string {
	return d.UnitType
}

//
// dbselectparamProducerRef functions
//
func (d *SelectFilterProducerRef) SetProducer(i int) {
	d.Producer = i
}

func (d SelectFilterProducerRef) GetProducer() int {
	return d.Producer
}

//
// dbselectparamSupplierRef functions
//
func (d *SelectFilterSupplierRef) SetSupplier(i int) {
	d.Supplier = i
}

func (d SelectFilterSupplierRef) GetSupplier() int {
	return d.Supplier
}

//
// dbselectparamPerson functions
//
func (d *SelectFilterPerson) SetEntity(i int) {
	d.Entity = i
}

func (d SelectFilterPerson) GetEntity() int {
	return d.Entity
}

//
// dbselectparamStoreLocation functions
//
func (d *SelectFilterStoreLocation) SetEntity(i int) {
	d.Entity = i
}

func (d *SelectFilterStoreLocation) GetEntity() int {
	return d.Entity
}

func (d *SelectFilterStoreLocation) GetStoreLocationCanStore() bool {
	return d.StoreLocationCanStore
}

func (d *SelectFilterStoreLocation) SetStoreLocationCanStore(b bool) {
	d.StoreLocationCanStore = b
}

func (d *SelectFilterStoreLocation) GetPermission() string {
	return d.Permission
}

func (d *SelectFilterStoreLocation) SetPermission(p string) {
	d.Permission = p
}

//
// dbselectparamProduct functions
//
func (d *SelectFilterProduct) SetTags(t []int) {
	d.Tags = t
}

func (d *SelectFilterProduct) GetTags() []int {
	return d.Tags
}

func (d *SelectFilterProduct) SetCategory(c int) {
	d.Category = c
}

func (d *SelectFilterProduct) GetCategory() int {
	return d.Category
}

func (d *SelectFilterProduct) SetShowBio(b bool) {
	d.ShowBio = b
}

func (d *SelectFilterProduct) GetShowBio() bool {
	return d.ShowBio
}

func (d *SelectFilterProduct) SetShowChem(b bool) {
	d.ShowChem = b
}

func (d *SelectFilterProduct) GetShowChem() bool {
	return d.ShowChem
}

func (d *SelectFilterProduct) SetShowConsu(b bool) {
	d.ShowConsu = b
}

func (d *SelectFilterProduct) GetShowConsu() bool {
	return d.ShowConsu
}

func (d *SelectFilterProduct) SetProducerRef(i int) {
	d.ProducerRef = i
}

func (d SelectFilterProduct) GetProducerRef() int {
	return d.ProducerRef
}

func (d *SelectFilterProduct) SetEntity(i int) {
	d.Entity = i
}

func (d SelectFilterProduct) GetEntity() int {
	return d.Entity
}

func (d *SelectFilterProduct) SetProduct(i int) {
	d.Product = i
}

func (d SelectFilterProduct) GetProduct() int {
	return d.Product
}

func (d *SelectFilterProduct) SetStorelocation(i int) {
	d.Storelocation = i
}

func (d SelectFilterProduct) GetStorelocation() int {
	return d.Storelocation
}

func (d *SelectFilterProduct) SetBookmark(b bool) {
	d.Bookmark = b
}

func (d SelectFilterProduct) GetBookmark() bool {
	return d.Bookmark
}

func (d *SelectFilterProduct) SetName(n int) {
	d.Name = n
}

func (d SelectFilterProduct) GetName() int {
	return d.Name
}

func (d *SelectFilterProduct) SetEmpiricalFormula(n int) {
	d.Name = n
}

func (d SelectFilterProduct) GetEmpiricalFormula() int {
	return d.EmpiricalFormula
}

func (d *SelectFilterProduct) SetCasNumber(n int) {
	d.CasNumber = n
}

func (d SelectFilterProduct) GetCasNumber() int {
	return d.CasNumber
}

func (d SelectFilterProduct) GetStorageBatchNumber() string {
	return d.StorageBatchNumber
}

func (d *SelectFilterProduct) SetStorageBatchNumber(n string) {
	d.StorageBatchNumber = n
}

func (d *SelectFilterProduct) SetStorageBarecode(n string) {
	d.StorageBarecode = n
}

func (d SelectFilterProduct) GetStorageBarecode() string {
	return d.StorageBarecode
}

func (d *SelectFilterProduct) SetSymbols(n []int) {
	d.Symbols = n
}

func (d SelectFilterProduct) GetSymbols() []int {
	return d.Symbols
}

func (d *SelectFilterProduct) SetCustomNamePartOf(n string) {
	d.CustomNamePartOf = n
}

func (d SelectFilterProduct) GetCustomNamePartOf() string {
	return d.CustomNamePartOf
}

func (d *SelectFilterProduct) SetHazardStatements(n []int) {
	d.HazardStatements = n
}

func (d SelectFilterProduct) GetHazardStatements() []int {
	return d.HazardStatements
}

func (d *SelectFilterProduct) SetPrecautionaryStatements(n []int) {
	d.PrecautionaryStatements = n
}

func (d SelectFilterProduct) GetPrecautionaryStatements() []int {
	return d.PrecautionaryStatements
}

func (d *SelectFilterProduct) SetSignalWord(n int) {
	d.SignalWord = n
}

func (d SelectFilterProduct) GetSignalWord() int {
	return d.SignalWord
}

func (d *SelectFilterProduct) SetCasNumberCmr(n bool) {
	d.CasNumberCmr = n
}

func (d SelectFilterProduct) GetCasNumberCmr() bool {
	return d.CasNumberCmr
}

func (d *SelectFilterProduct) SetProductSpecificity(s string) {
	d.ProductSpecificity = s
}

func (d *SelectFilterProduct) GetProductSpecificity() string {
	return d.ProductSpecificity
}

func (d *SelectFilterProduct) SetBorrowing(b bool) {
	d.Borrowing = b
}

func (d *SelectFilterProduct) GetBorrowing() bool {
	return d.Borrowing
}

func (d *SelectFilterProduct) SetStorageToDestroy(b bool) {
	d.StorageToDestroy = b
}

func (d *SelectFilterProduct) GetStorageToDestroy() bool {
	return d.StorageToDestroy
}

//
// dbselectparamStorage functions
//
func (d *SelectFilterStorage) SetTags(t []int) {
	d.Tags = t
}

func (d *SelectFilterStorage) GetTags() []int {
	return d.Tags
}

func (d *SelectFilterStorage) SetCategory(c int) {
	d.Category = c
}

func (d *SelectFilterStorage) GetCategory() int {
	return d.Category
}

func (d *SelectFilterStorage) SetShowBio(b bool) {
	d.ShowBio = b
}

func (d *SelectFilterStorage) GetShowBio() bool {
	return d.ShowBio
}

func (d *SelectFilterStorage) SetShowChem(b bool) {
	d.ShowChem = b
}

func (d *SelectFilterStorage) GetShowChem() bool {
	return d.ShowChem
}

func (d *SelectFilterStorage) SetShowConsu(b bool) {
	d.ShowConsu = b
}

func (d *SelectFilterStorage) GetShowConsu() bool {
	return d.ShowConsu
}

func (d *SelectFilterStorage) SetIds(i []int) {
	d.Ids = i
}

func (d *SelectFilterStorage) GetIds() []int {
	return d.Ids
}

func (d *SelectFilterStorage) SetProducerRef(i int) {
	d.ProducerRef = i
}

func (d SelectFilterStorage) GetProducerRef() int {
	return d.ProducerRef
}

func (d *SelectFilterStorage) SetEntity(i int) {
	d.Entity = i
}

func (d SelectFilterStorage) GetEntity() int {
	return d.Entity
}

func (d *SelectFilterStorage) SetProduct(i int) {
	d.Product = i
}

func (d SelectFilterStorage) GetProduct() int {
	return d.Product
}

func (d *SelectFilterStorage) SetStorelocation(i int) {
	d.Storelocation = i
}

func (d SelectFilterStorage) GetStorelocation() int {
	return d.Storelocation
}

func (d *SelectFilterStorage) SetStorage(i int) {
	d.Storage = i
}

func (d SelectFilterStorage) GetStorage() int {
	return d.Storage
}

func (d *SelectFilterStorage) SetBookmark(b bool) {
	d.Bookmark = b
}

func (d SelectFilterStorage) GetBookmark() bool {
	return d.Bookmark
}

func (d SelectFilterStorage) GetHistory() bool {
	return d.History
}

func (d *SelectFilterStorage) SetHistory(b bool) {
	d.History = b
}

func (d SelectFilterStorage) GetStorageArchive() bool {
	return d.StorageArchive
}

func (d *SelectFilterStorage) SetStorageArchive(b bool) {
	d.StorageArchive = b
}

func (d *SelectFilterStorage) SetName(n int) {
	d.Name = n
}

func (d SelectFilterStorage) GetName() int {
	return d.Name
}

func (d *SelectFilterStorage) SetEmpiricalFormula(n int) {
	d.Name = n
}

func (d SelectFilterStorage) GetEmpiricalFormula() int {
	return d.EmpiricalFormula
}

func (d *SelectFilterStorage) SetCasNumber(n int) {
	d.CasNumber = n
}

func (d SelectFilterStorage) GetCasNumber() int {
	return d.CasNumber
}

func (d *SelectFilterStorage) SetStorageBatchNumber(n string) {
	d.StorageBatchNumber = n
}

func (d SelectFilterStorage) GetStorageBatchNumber() string {
	return d.StorageBatchNumber
}

func (d *SelectFilterStorage) SetStorageBarecode(n string) {
	d.StorageBarecode = n
}

func (d SelectFilterStorage) GetStorageBarecode() string {
	return d.StorageBarecode
}

func (d *SelectFilterStorage) SetSymbols(n []int) {
	d.Symbols = n
}

func (d SelectFilterStorage) GetSymbols() []int {
	return d.Symbols
}

func (d *SelectFilterStorage) SetCustomNamePartOf(n string) {
	d.CustomNamePartOf = n
}

func (d SelectFilterStorage) GetCustomNamePartOf() string {
	return d.CustomNamePartOf
}

func (d *SelectFilterStorage) SetHazardStatements(n []int) {
	d.HazardStatements = n
}

func (d SelectFilterStorage) GetHazardStatements() []int {
	return d.HazardStatements
}

func (d *SelectFilterStorage) SetPrecautionaryStatements(n []int) {
	d.PrecautionaryStatements = n
}

func (d SelectFilterStorage) GetPrecautionaryStatements() []int {
	return d.PrecautionaryStatements
}

func (d *SelectFilterStorage) SetSignalWord(n int) {
	d.SignalWord = n
}

func (d SelectFilterStorage) GetSignalWord() int {
	return d.SignalWord
}

func (d *SelectFilterStorage) SetCasNumberCmr(n bool) {
	d.CasNumberCmr = n
}

func (d SelectFilterStorage) GetCasNumberCmr() bool {
	return d.CasNumberCmr
}

func (d *SelectFilterStorage) SetBorrowing(b bool) {
	d.Borrowing = b
}

func (d *SelectFilterStorage) GetBorrowing() bool {
	return d.Borrowing
}

func (d *SelectFilterStorage) SetStorageToDestroy(b bool) {
	d.StorageToDestroy = b
}

func (d *SelectFilterStorage) GetStorageToDestroy() bool {
	return d.StorageToDestroy
}

// Newdbselectparam returns a dbselectparam struct
// with values populated from the request parameters
func Newdbselectparam(r *http.Request, f func(string) (string, error)) (*selectFilter, *AppError) {

	var err error

	// initializing default values
	dsp := selectFilter{
		LoggedPersonID: 0,
		Search:         "%%",
		OrderBy:        "",
		Order:          "asc",
		Offset:         0,
		Limit:          ^uint64(0),
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
					Code:    http.StatusBadRequest,
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
func NewdbselectparamProduct(r *http.Request, f func(string) (string, error)) (*SelectFilterProduct, *AppError) {

	var (
		err  error
		aerr *AppError
		dsp  *selectFilter
		dspp SelectFilterProduct
	)

	// init defaults
	dspp.Category = -1
	dspp.ProducerRef = -1
	dspp.Entity = -1
	dspp.Product = -1
	dspp.Storelocation = -1
	dspp.Bookmark = false
	dspp.Name = -1
	dspp.CasNumber = -1
	dspp.EmpiricalFormula = -1
	dspp.StorageBarecode = ""
	dspp.StorageBatchNumber = ""
	dspp.CustomNamePartOf = ""
	dspp.SignalWord = -1
	dspp.CasNumberCmr = false
	dspp.ProductSpecificity = ""
	dspp.Borrowing = false
	dspp.StorageToDestroy = false
	if dsp, aerr = Newdbselectparam(r, f); aerr != nil {
		return nil, aerr
	}
	dspp.selectFilter = *dsp

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
		if producerrefid, ok := r.URL.Query()["producerref"]; ok {
			if dspp.ProducerRef, err = strconv.Atoi(producerrefid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "producerref atoi conversion",
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
		if storage_batchnumber, ok := r.URL.Query()["storage_batchnumber"]; ok {
			dspp.StorageBatchNumber = "%" + storage_batchnumber[0] + "%"
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
		if borrowing, ok := r.URL.Query()["borrowing"]; ok {
			if dspp.Borrowing, err = strconv.ParseBool(borrowing[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "borrowing bool conversion",
				}
			}
		}
		if storage_to_destroy, ok := r.URL.Query()["storage_to_destroy"]; ok {
			if dspp.StorageToDestroy, err = strconv.ParseBool(storage_to_destroy[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "storage_to_destroy bool conversion",
				}
			}
		}
		if showbio, ok := r.URL.Query()["showbio"]; ok {
			if dspp.ShowBio, err = strconv.ParseBool(showbio[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "showbio bool conversion",
				}
			}
		}
		if showchem, ok := r.URL.Query()["showchem"]; ok {
			if dspp.ShowChem, err = strconv.ParseBool(showchem[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "showchem bool conversion",
				}
			}
		}
		if showconsu, ok := r.URL.Query()["showconsu"]; ok {
			if dspp.ShowConsu, err = strconv.ParseBool(showconsu[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "showconsu bool conversion",
				}
			}
		}
		if categoryid, ok := r.URL.Query()["category"]; ok {
			if dspp.Category, err = strconv.Atoi(categoryid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "category atoi conversion",
				}
			}
		}
		if tagsids, ok := r.URL.Query()["tags[]"]; ok {
			var tint int
			for _, t := range tagsids {
				if tint, err = strconv.Atoi(t); err != nil {
					return nil, &AppError{
						Error:   err,
						Code:    http.StatusInternalServerError,
						Message: "tag atoi conversion",
					}
				}
				dspp.Tags = append(dspp.Tags, tint)
			}
		}
	}
	return &dspp, nil

}

// NewdbselectparamStorage returns a dbselectparamStorage struct
// with values populated from the request parameters
func NewdbselectparamStorage(r *http.Request, f func(string) (string, error)) (*SelectFilterStorage, *AppError) {

	var (
		err  error
		aerr *AppError
		dsp  *selectFilter
		dsps SelectFilterStorage
	)

	// init defaults
	dsps.Category = -1
	dsps.ProducerRef = -1
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
	dsps.StorageBatchNumber = ""
	dsps.CustomNamePartOf = ""
	dsps.SignalWord = -1
	dsps.CasNumberCmr = false
	dsps.Borrowing = false
	dsps.StorageToDestroy = false
	if dsp, aerr = Newdbselectparam(r, f); aerr != nil {
		return nil, aerr
	}
	dsps.selectFilter = *dsp

	if r != nil {
		if ids, ok := r.URL.Query()["ids"]; ok {
			for _, id := range ids {
				idInt, err := strconv.Atoi(id)
				if err != nil {
					return nil, &AppError{
						Error:   err,
						Code:    http.StatusInternalServerError,
						Message: "ids atoi conversion",
					}
				}
				dsps.Ids = append(dsps.Ids, idInt)
			}
		}
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
		if producerrefid, ok := r.URL.Query()["producerref"]; ok {
			if dsps.ProducerRef, err = strconv.Atoi(producerrefid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "producerref atoi conversion",
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
		if storage_batchnumber, ok := r.URL.Query()["storage_batchnumber"]; ok {
			dsps.StorageBatchNumber = "%" + storage_batchnumber[0] + "%"
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
		if borrowing, ok := r.URL.Query()["borrowing"]; ok {
			if dsps.Borrowing, err = strconv.ParseBool(borrowing[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "borrowing bool conversion",
				}
			}
		}
		if storage_to_destroy, ok := r.URL.Query()["storage_to_destroy"]; ok {
			if dsps.StorageToDestroy, err = strconv.ParseBool(storage_to_destroy[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "storage_to_destroy bool conversion",
				}
			}
		}
		if showbio, ok := r.URL.Query()["showbio"]; ok {
			if dsps.ShowBio, err = strconv.ParseBool(showbio[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "showbio bool conversion",
				}
			}
		}
		if showchem, ok := r.URL.Query()["showchem"]; ok {
			if dsps.ShowChem, err = strconv.ParseBool(showchem[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "showchem bool conversion",
				}
			}
		}
		if showconsu, ok := r.URL.Query()["showconsu"]; ok {
			if dsps.ShowConsu, err = strconv.ParseBool(showconsu[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "showconsu bool conversion",
				}
			}
		}
		if categoryid, ok := r.URL.Query()["category"]; ok {
			if dsps.Category, err = strconv.Atoi(categoryid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "category atoi conversion",
				}
			}
		}
		if tagsids, ok := r.URL.Query()["tags[]"]; ok {
			var tint int
			for _, t := range tagsids {
				if tint, err = strconv.Atoi(t); err != nil {
					return nil, &AppError{
						Error:   err,
						Code:    http.StatusInternalServerError,
						Message: "tag atoi conversion",
					}
				}
				dsps.Tags = append(dsps.Tags, tint)
			}
		}
	}
	return &dsps, nil

}

// NewdbselectparamStoreLocation returns a dbselectparamStoreLocation struct
// with values populated from the request parameters
func NewdbselectparamStoreLocation(r *http.Request, f func(string) (string, error)) (*SelectFilterStoreLocation, *AppError) {

	var (
		err   error
		aerr  *AppError
		dsp   *selectFilter
		dspsl SelectFilterStoreLocation
	)

	// init defaults
	dspsl.Entity = -1
	dspsl.Permission = "r"
	if dsp, aerr = Newdbselectparam(r, f); aerr != nil {
		return nil, aerr
	}
	dspsl.selectFilter = *dsp

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
					Message: "entity atoi conversion",
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

// NewdbselectparamUnit returns a dbselectparamUnit struct
// with values populated from the request parameters
func NewdbselectparamUnit(r *http.Request, f func(string) (string, error)) (*SelectFilterUnit, *AppError) {

	var (
		aerr *AppError
		dsp  *selectFilter
		dspu SelectFilterUnit
	)

	// init defaults
	dspu.UnitType = ""
	if dsp, aerr = Newdbselectparam(r, f); aerr != nil {
		return nil, aerr
	}
	dspu.selectFilter = *dsp

	if r != nil {
		if o, ok := r.URL.Query()["sort"]; ok {
			dspu.OrderBy = o[0]
		} else {
			dspu.OrderBy = "unit_id"
		}
		if unitType, ok := r.URL.Query()["unit_type"]; ok {
			dspu.UnitType = unitType[0]
		}
	}
	return &dspu, nil
}

// NewdbselectparamProducerRef returns a dbselectparamProducerRef struct
// with values populated from the request parameters
func NewdbselectparamProducerRef(r *http.Request, f func(string) (string, error)) (*SelectFilterProducerRef, *AppError) {

	var (
		err  error
		aerr *AppError
		dsp  *selectFilter
		dspp SelectFilterProducerRef
	)

	// init defaults
	dspp.Producer = -1
	if dsp, aerr = Newdbselectparam(r, f); aerr != nil {
		return nil, aerr
	}
	dspp.selectFilter = *dsp

	if r != nil {
		if o, ok := r.URL.Query()["sort"]; ok {
			dspp.OrderBy = o[0]
		} else {
			dspp.OrderBy = "producerref_id"
		}
		if producerid, ok := r.URL.Query()["producer"]; ok {
			if dspp.Producer, err = strconv.Atoi(producerid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "producer atoi conversion",
				}
			}
		}
	}
	return &dspp, nil
}

// NewdbselectparamSupplierRef returns a dbselectparamSupplierRef struct
// with values populated from the request parameters
func NewdbselectparamSupplierRef(r *http.Request, f func(string) (string, error)) (*SelectFilterSupplierRef, *AppError) {

	var (
		err  error
		aerr *AppError
		dsp  *selectFilter
		dspp SelectFilterSupplierRef
	)

	// init defaults
	dspp.Supplier = -1
	if dsp, aerr = Newdbselectparam(r, f); aerr != nil {
		return nil, aerr
	}
	dspp.selectFilter = *dsp

	if r != nil {
		if o, ok := r.URL.Query()["sort"]; ok {
			dspp.OrderBy = o[0]
		} else {
			dspp.OrderBy = "supplierref_id"
		}
		if supplierid, ok := r.URL.Query()["supplier"]; ok {
			if dspp.Supplier, err = strconv.Atoi(supplierid[0]); err != nil {
				return nil, &AppError{
					Error:   err,
					Code:    http.StatusInternalServerError,
					Message: "supplier atoi conversion",
				}
			}
		}
	}
	return &dspp, nil
}

// NewdbselectparamPerson returns a dbselectparamStorePerson struct
// with values populated from the request parameters
func NewdbselectparamPerson(r *http.Request, f func(string) (string, error)) (*SelectFilterPerson, *AppError) {

	var (
		err  error
		aerr *AppError
		dsp  *selectFilter
		dspp SelectFilterPerson
	)

	// init defaults
	dspp.Entity = -1
	if dsp, aerr = Newdbselectparam(r, f); aerr != nil {
		return nil, aerr
	}
	dspp.selectFilter = *dsp

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
					Message: "entity atoi conversion",
				}
			}
		}
	}
	return &dspp, nil

}

// NewdbselectparamEntity returns a dbselectparamEntity struct
// with values populated from the request parameters
func NewdbselectparamEntity(r *http.Request, f func(string) (string, error)) (*SelectFilterEntity, *AppError) {

	var (
		aerr *AppError
		dsp  *selectFilter
		dspe SelectFilterEntity
	)

	// init defaults
	if dsp, aerr = Newdbselectparam(r, f); aerr != nil {
		return nil, aerr
	}
	dspe.selectFilter = *dsp
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

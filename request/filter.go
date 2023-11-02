package request

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/tbellembois/gochimitheque/models"
)

// type paramType int

// const (
// 	Int paramType = iota
// 	String
// 	SliceOfInt
// 	None
// 	Bool
// )

type Filter struct {
	LoggedPersonID int // logged person, used to filter results
	Search         string
	OrderBy        string
	Order          string
	Offset         uint64
	Limit          uint64

	Bookmark                bool
	Borrowing               bool
	CasNumber               int // id
	CasNumberCmr            bool
	Category                int
	CustomNamePartOf        string
	EmpiricalFormula        int   // id
	Entity                  int   // id
	HazardStatements        []int // ids
	History                 bool
	Ids                     []int // FIXME: Storage_id
	Name                    int   // id
	Permission              string
	PrecautionaryStatements []int // ids
	Producer                int
	ProducerRef             int // id
	Product                 int // id
	ProductSpecificity      string
	ShowBio                 bool
	ShowChem                bool
	ShowConsu               bool
	SignalWord              int // id
	Storage                 int // id
	StorageArchive          bool
	StorageBarecode         string
	StorageBatchNumber      string
	StorageToDestroy        bool
	Storelocation           int // id
	StoreLocationCanStore   bool
	Supplier                int
	Symbols                 []int // ids
	Tags                    []int
	UnitType                string
}

// var filterMap map[string]paramType

// func init() {

// 	filterMap = make(map[string]paramType)
// 	filterMap["search"] = String
// 	filterMap["order"] = String
// 	// TODO: adapt to type
// 	filterMap["sort"] = String
// 	filterMap["offset"] = Int
// 	filterMap["limit"] = Int
// 	filterMap["export"] = None

// 	filterMap["bookmark"] = Bool
// 	filterMap["borrowing"] = Bool
// 	filterMap["casnumber_cmr"] = Bool
// 	filterMap["casnumber"] = Int
// 	filterMap["category"] = Int
// 	filterMap["custom_name_part_of"] = String
// 	filterMap["empiricalformula"] = Int
// 	filterMap["entity"] = Int
// 	filterMap["hazardstatements[]"] = SliceOfInt
// 	filterMap["history"] = Bool
// 	// FIXME: storage_id[]
// 	filterMap["ids"] = SliceOfInt
// 	filterMap["name"] = Int
// 	filterMap["permission"] = None
// 	filterMap["precautionarystatements[]"] = SliceOfInt
// 	filterMap["producer"] = Int
// 	filterMap["producerref"] = Int
// 	filterMap["product_specificity"] = String
// 	filterMap["product"] = Int
// 	filterMap["showbio"] = Bool
// 	filterMap["showchem"] = Bool
// 	filterMap["showconsu"] = Bool
// 	filterMap["signalword"] = Int
// 	filterMap["storage"] = Int
// 	filterMap["storage_barecode"] = String
// 	filterMap["storage_batchnumber"] = String
// 	filterMap["storage_to_destroy"] = Bool
// 	filterMap["storage_archive"] = Bool
// 	filterMap["storelocation"] = Int
// 	filterMap["storelocation_canstore"] = Bool
// 	filterMap["supplier"] = Int
// 	filterMap["symbols[]"] = SliceOfInt
// 	filterMap["tags[]"] = SliceOfInt
// 	filterMap["unit_type"] = String

// }

func NewFilter(r *http.Request) (filter *Filter, apperr *models.AppError) {
	var err error

	// Init defaults.
	filter = &Filter{
		LoggedPersonID: 0,
		Search:         "%%",
		OrderBy:        "",
		Order:          "asc",
		Offset:         0,
		Limit:          ^uint64(0),

		CasNumber:        -1,
		Category:         -1,
		EmpiricalFormula: -1,
		Entity:           -1,
		Name:             -1,
		Permission:       "r",
		Producer:         -1,
		ProducerRef:      -1,
		Product:          -1,
		SignalWord:       -1,
		Storage:          -1,
		Storelocation:    -1,
		Supplier:         -1,
	}

	if r == nil {
		return
	}

	c := ContainerFromRequestContext(r)
	filter.LoggedPersonID = c.PersonID

	// if s, ok := r.URL.Query()["search"]; ok {
	// 	if f != nil && s[0] != "" {
	// 		fs, err := f(s[0])
	// 		if err != nil {
	// 			return nil, &models.AppError{
	// 				OriginalError: err,
	// 				Code:          http.StatusBadRequest,
	// 				Message:       "error calling f",
	// 			}
	// 		}

	// 		filter.Search = "%" + fs + "%"
	// 	} else {
	// 		filter.Search = "%" + s[0] + "%"
	// 	}
	// }

	if s, ok := r.URL.Query()["search"]; ok {
		filter.Search = "%" + s[0] + "%"
	}

	if o, ok := r.URL.Query()["order"]; ok {
		filter.Order = o[0]
	}

	if o, ok := r.URL.Query()["sort"]; ok {
		filter.OrderBy = o[0]
	}

	if o, ok := r.URL.Query()["offset"]; ok {
		var of int

		if of, err = strconv.Atoi(o[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "offset atoi conversion",
			}
		}

		filter.Offset = uint64(of)
	}

	// No limit on export.
	if _, ok := r.URL.Query()["export"]; !ok {
		if l, ok := r.URL.Query()["limit"]; ok {
			var lm int

			if lm, err = strconv.Atoi(l[0]); err != nil {
				return nil, &models.AppError{
					OriginalError: err,
					Code:          http.StatusInternalServerError,
					Message:       "limit atoi conversion",
				}
			}

			filter.Limit = uint64(lm)
		}
	}

	if entityid, ok := r.URL.Query()["entity"]; ok {
		if filter.Entity, err = strconv.Atoi(entityid[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "entity atoi conversion",
			}
		}
	}

	if producerrefid, ok := r.URL.Query()["producerref"]; ok {
		if filter.ProducerRef, err = strconv.Atoi(producerrefid[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "producerref atoi conversion",
			}
		}
	}

	if productid, ok := r.URL.Query()["product"]; ok {
		if filter.Product, err = strconv.Atoi(productid[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "product atoi conversion",
			}
		}
	}

	if storelocationid, ok := r.URL.Query()["storelocation"]; ok {
		if filter.Storelocation, err = strconv.Atoi(storelocationid[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "storelocation atoi conversion",
			}
		}
	}

	if storageid, ok := r.URL.Query()["storage"]; ok {
		if filter.Storage, err = strconv.Atoi(storageid[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "storage atoi conversion",
			}
		}
	}

	if bookmark, ok := r.URL.Query()["bookmark"]; ok {
		if filter.Bookmark, err = strconv.ParseBool(bookmark[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "bookmark bool conversion",
			}
		}
	}

	if nameid, ok := r.URL.Query()["name"]; ok {
		if filter.Name, err = strconv.Atoi(nameid[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "name atoi conversion",
			}
		}
	}

	if casnumberid, ok := r.URL.Query()["casnumber"]; ok {
		if filter.CasNumber, err = strconv.Atoi(casnumberid[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "casnumber atoi conversion",
			}
		}
	}

	if empiricalformulaid, ok := r.URL.Query()["empiricalformula"]; ok {
		if filter.EmpiricalFormula, err = strconv.Atoi(empiricalformulaid[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "empiricalformula atoi conversion",
			}
		}
	}

	if symbolsids, ok := r.URL.Query()["symbols[]"]; ok {
		var sint int

		for _, s := range symbolsids {
			if sint, err = strconv.Atoi(s); err != nil {
				return nil, &models.AppError{
					OriginalError: err,
					Code:          http.StatusInternalServerError,
					Message:       "symbol atoi conversion",
				}
			}

			filter.Symbols = append(filter.Symbols, sint)
		}
	}

	if hsids, ok := r.URL.Query()["hazardstatements[]"]; ok {
		var sint int

		for _, s := range hsids {
			if sint, err = strconv.Atoi(s); err != nil {
				return nil, &models.AppError{
					OriginalError: err,
					Code:          http.StatusInternalServerError,
					Message:       "hazardstatement atoi conversion",
				}
			}

			filter.HazardStatements = append(filter.HazardStatements, sint)
		}
	}

	if psids, ok := r.URL.Query()["precautionarystatements[]"]; ok {
		var sint int

		for _, s := range psids {
			if sint, err = strconv.Atoi(s); err != nil {
				return nil, &models.AppError{
					OriginalError: err,
					Code:          http.StatusInternalServerError,
					Message:       "precautionarystatement atoi conversion",
				}
			}

			filter.PrecautionaryStatements = append(filter.PrecautionaryStatements, sint)
		}
	}

	if signalwordid, ok := r.URL.Query()["signalword"]; ok {
		if filter.SignalWord, err = strconv.Atoi(signalwordid[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "signalword atoi conversion",
			}
		}
	}

	if storage_barecode, ok := r.URL.Query()["storage_barecode"]; ok {
		filter.StorageBarecode = "%" + storage_barecode[0] + "%"
	}

	if storage_batchnumber, ok := r.URL.Query()["storage_batchnumber"]; ok {
		filter.StorageBatchNumber = "%" + storage_batchnumber[0] + "%"
	}

	if custom_name_part_of, ok := r.URL.Query()["custom_name_part_of"]; ok {
		filter.CustomNamePartOf = "%" + strings.ToUpper(custom_name_part_of[0]) + "%"
	}

	if product_specificity, ok := r.URL.Query()["product_specificity"]; ok {
		filter.ProductSpecificity = product_specificity[0]
	}

	if casnumber_cmr, ok := r.URL.Query()["casnumber_cmr"]; ok {
		if filter.CasNumberCmr, err = strconv.ParseBool(casnumber_cmr[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "casnumber_cmr bool conversion",
			}
		}
	}

	if borrowing, ok := r.URL.Query()["borrowing"]; ok {
		if filter.Borrowing, err = strconv.ParseBool(borrowing[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "borrowing bool conversion",
			}
		}
	}

	if storage_to_destroy, ok := r.URL.Query()["storage_to_destroy"]; ok {
		if filter.StorageToDestroy, err = strconv.ParseBool(storage_to_destroy[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "storage_to_destroy bool conversion",
			}
		}
	}

	if storage_archive, ok := r.URL.Query()["storage_archive"]; ok {
		if filter.StorageArchive, err = strconv.ParseBool(storage_archive[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "storage_archive bool conversion",
			}
		}
	}

	if showbio, ok := r.URL.Query()["showbio"]; ok {
		if filter.ShowBio, err = strconv.ParseBool(showbio[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "showbio bool conversion",
			}
		}
	}

	if showchem, ok := r.URL.Query()["showchem"]; ok {
		if filter.ShowChem, err = strconv.ParseBool(showchem[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "showchem bool conversion",
			}
		}
	}

	if showconsu, ok := r.URL.Query()["showconsu"]; ok {
		if filter.ShowConsu, err = strconv.ParseBool(showconsu[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "showconsu bool conversion",
			}
		}
	}

	if categoryid, ok := r.URL.Query()["category"]; ok {
		if filter.Category, err = strconv.Atoi(categoryid[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "category atoi conversion",
			}
		}
	}

	if tagsids, ok := r.URL.Query()["tags[]"]; ok {
		var tint int

		for _, t := range tagsids {
			if tint, err = strconv.Atoi(t); err != nil {
				return nil, &models.AppError{
					OriginalError: err,
					Code:          http.StatusInternalServerError,
					Message:       "tag atoi conversion",
				}
			}

			filter.Tags = append(filter.Tags, tint)
		}
	}

	// FIXME: storage_id
	if ids, ok := r.URL.Query()["ids"]; ok {
		for _, id := range ids {
			idInt, err := strconv.Atoi(id)
			if err != nil {
				return nil, &models.AppError{
					OriginalError: err,
					Code:          http.StatusInternalServerError,
					Message:       "ids atoi conversion",
				}
			}

			filter.Ids = append(filter.Ids, idInt)
		}
	}

	if history, ok := r.URL.Query()["history"]; ok {
		if filter.History, err = strconv.ParseBool(history[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "history bool conversion",
			}
		}
	}

	if c, ok := r.URL.Query()["storelocation_canstore"]; ok {
		if filter.StoreLocationCanStore, err = strconv.ParseBool(c[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "storelocation_canstore bool conversion",
			}
		}
	}

	if p, ok := r.URL.Query()["permission"]; ok {
		filter.Permission = p[0]
	}

	if unitType, ok := r.URL.Query()["unit_type"]; ok {
		filter.UnitType = unitType[0]
	}

	if producerid, ok := r.URL.Query()["producer"]; ok {
		if filter.Producer, err = strconv.Atoi(producerid[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "producer atoi conversion",
			}
		}
	}

	if supplierid, ok := r.URL.Query()["supplier"]; ok {
		if filter.Supplier, err = strconv.Atoi(supplierid[0]); err != nil {
			return nil, &models.AppError{
				OriginalError: err,
				Code:          http.StatusInternalServerError,
				Message:       "supplier atoi conversion",
			}
		}
	}

	return
}

package models

import (
	"database/sql"
	"encoding/csv"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tbellembois/gochimitheque/logger"
)

// ChimithequeContextKey is the Go request context
// used in each request
type ChimithequeContextKey string

// AppError is the error type returned by the custom handlers
type AppError struct {
	Error   error
	Message string
	Code    int
}

// AppHandlerFunc is an HandlerFunc returning an AppError
type AppHandlerFunc func(http.ResponseWriter, *http.Request) *AppError

// Stock is a store location stock for a given product
type Stock struct {
	Total   float64 `json:"total"`
	Current float64 `json:"current"`
	Unit    Unit    `json:"unit"`
}

// WelcomeAnnounce is the custom welcome page message
type WelcomeAnnounce struct {
	WelcomeAnnounceID   int    `db:"welcomeannounce_id" json:"welcomeannounce_id" schema:"welcomeannounce_id"`
	WelcomeAnnounceText string `db:"welcomeannounce_text" json:"welcomeannounce_text" schema:"welcomeannounce_text"`
	WelcomeAnnounceHTML string `db:"welcomeannounce_html" json:"welcomeannounce_html" schema:"welcomeannounce_html"`
}

// StoreLocation is where products are stored in entities
type StoreLocation struct {
	// nullable values to handle optional StoreLocation foreign key (gorilla shema nil values)
	StoreLocationID       sql.NullInt64  `db:"storelocation_id" json:"storelocation_id" schema:"storelocation_id"`
	StoreLocationName     sql.NullString `db:"storelocation_name" json:"storelocation_name" schema:"storelocation_name"`
	StoreLocationCanStore sql.NullBool   `db:"storelocation_canstore" json:"storelocation_canstore" schema:"storelocation_canstore"`
	StoreLocationColor    sql.NullString `db:"storelocation_color" json:"storelocation_color" schema:"storelocation_color"`
	Entity                `db:"entity" json:"entity" schema:"entity"`
	StoreLocation         *StoreLocation `db:"storelocation" json:"storelocation" schema:"storelocation"`
	StoreLocationFullPath string         `db:"storelocation_fullpath" json:"storelocation_fullpath" schema:"storelocation_fullpath"`

	Children []*StoreLocation `db:"-" json:"children" schema:"-"`
	Stocks   []Stock          `db:"-" json:"stock" schema:"-"`
}

// Entity represent a department, a laboratory...
type Entity struct {
	EntityID          int       `db:"entity_id" json:"entity_id" schema:"entity_id"`
	EntityName        string    `db:"entity_name" json:"entity_name" schema:"entity_name"`
	EntityDescription string    `db:"entity_description" json:"entity_description" schema:"entity_description"`
	Managers          []*Person `db:"-" json:"managers" schema:"managers"`

	// total store location count
	EntitySLC int `db:"entity_slc" json:"entity_slc" schema:"entity_slc"` // not in db but sqlx requires the "db" entry
	// total person count
	EntityPC int `db:"entity_pc" json:"entity_pc" schema:"entity_pc"` // not in db but sqlx requires the "db" entry
}

// Equal tests the entity equality
func (e1 Entity) Equal(e2 Entity) bool {
	return e1.EntityID == e2.EntityID
}

// Person represent a person
type Person struct {
	PersonID       int           `db:"person_id" json:"person_id" schema:"person_id"`
	PersonEmail    string        `db:"person_email" json:"person_email" schema:"person_email"`
	PersonPassword string        `db:"person_password" json:"person_password" schema:"person_password"`
	Permissions    []*Permission `db:"-" schema:"permissions"`
	Entities       []*Entity     `db:"-" schema:"entities"`
	CaptchaText    string        `db:"-" schema:"captcha_text" json:"captcha_text"`
	CaptchaUID     string        `db:"-" schema:"captcha_uid" json:"captcha_uid"`
}

// Unit is a volume or weight unit
type Unit struct {
	UnitID         sql.NullInt64  `db:"unit_id" json:"unit_id" schema:"unit_id"`
	UnitLabel      sql.NullString `db:"unit_label" json:"unit_label" schema:"unit_label"`
	UnitType       sql.NullString `db:"unit_type" json:"unit_type" schema:"unit_type"`
	Unit           *Unit          `db:"unit" json:"unit" schema:"unit"` // reference
	UnitMultiplier int            `db:"unit_multiplier" json:"-" schema:"-"`
}

// Supplier is a product supplier
type Supplier struct {
	C             int            `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	SupplierID    sql.NullInt64  `db:"supplier_id" json:"supplier_id" schema:"supplier_id"`
	SupplierLabel sql.NullString `db:"supplier_label" json:"supplier_label" schema:"supplier_label"`
}

// SupplierRef is a product supplier reference
type SupplierRef struct {
	C                int       `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	SupplierRefID    int       `db:"supplierref_id" json:"supplierref_id" schema:"supplierref_id"`
	SupplierRefLabel string    `db:"supplierref_label" json:"supplierref_label" schema:"supplierref_label"`
	Supplier         *Supplier `db:"supplier" json:"supplier" schema:"supplier"`
}

// Producer is a product producer
type Producer struct {
	C             int            `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	ProducerID    sql.NullInt64  `db:"producer_id" json:"producer_id" schema:"producer_id"`
	ProducerLabel sql.NullString `db:"producer_label" json:"producer_label" schema:"producer_label"`
}

// ProducerRef is a product producer reference
type ProducerRef struct {
	C                int            `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	ProducerRefID    sql.NullInt64  `db:"producerref_id" json:"producerref_id" schema:"producerref_id"`
	ProducerRefLabel sql.NullString `db:"producerref_label" json:"producerref_label" schema:"producerref_label"`
	Producer         *Producer      `db:"producer" json:"producer" schema:"producer"`
	// ProducerID       sql.NullInt64  `json:"producer_id" schema:"producer_id"`       // for Gorilla only
	// ProducerLabel    sql.NullString `json:"producer_label" schema:"producer_label"` // for Gorilla only
}

// Category is a product category
type Category struct {
	C             int            `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	CategoryID    sql.NullInt64  `db:"category_id" json:"category_id" schema:"category_id"`
	CategoryLabel sql.NullString `db:"category_label" json:"category_label" schema:"category_label"`
}

// Tag is a product tag
type Tag struct {
	C        int    `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	TagID    int    `db:"tag_id" json:"tag_id" schema:"tag_id"`
	TagLabel string `db:"tag_label" json:"tag_label" schema:"tag_label"`
}

// Storage is a product storage in a store location
type Storage struct {
	StorageID                sql.NullInt64   `db:"storage_id" json:"storage_id" schema:"storage_id"`
	StorageCreationDate      time.Time       `db:"storage_creationdate" json:"storage_creationdate" schema:"storage_creationdate"`
	StorageModificationDate  time.Time       `db:"storage_modificationdate" json:"storage_modificationdate" schema:"storage_modificationdate"`
	StorageEntryDate         sql.NullTime    `db:"storage_entrydate" json:"storage_entrydate" schema:"storage_entrydate"`
	StorageExitDate          sql.NullTime    `db:"storage_exitdate" json:"storage_exitdate" schema:"storage_exitdate"`
	StorageOpeningDate       sql.NullTime    `db:"storage_openingdate" json:"storage_openingdate" schema:"storage_openingdate"`
	StorageExpirationDate    sql.NullTime    `db:"storage_expirationdate" json:"storage_expirationdate" schema:"storage_expirationdate"`
	StorageComment           sql.NullString  `db:"storage_comment" json:"storage_comment" schema:"storage_comment"`
	StorageReference         sql.NullString  `db:"storage_reference" json:"storage_reference" schema:"storage_reference"`
	StorageBatchNumber       sql.NullString  `db:"storage_batchnumber" json:"storage_batchnumber" schema:"storage_batchnumber"`
	StorageQuantity          sql.NullFloat64 `db:"storage_quantity" json:"storage_quantity" schema:"storage_quantity"`
	StorageNbItem            int             `db:"-" json:"storage_nbitem" schema:"storage_nbitem"`
	StorageIdenticalBarecode sql.NullBool    `db:"-" json:"storage_identicalbarecode" schema:"storage_identicalbarecode"`
	StorageBarecode          sql.NullString  `db:"storage_barecode" json:"storage_barecode" schema:"storage_barecode"`
	StorageQRCode            []byte          `db:"storage_qrcode" json:"storage_qrcode" schema:"storage_qrcode"`
	StorageToDestroy         sql.NullBool    `db:"storage_todestroy" json:"storage_todestroy" schema:"storage_todestroy"`
	StorageArchive           sql.NullBool    `db:"storage_archive" json:"storage_archive" schema:"storage_archive"`
	StorageConcentration     sql.NullInt64   `db:"storage_concentration" json:"storage_concentration" schema:"storage_concentration"`
	StorageNumberOfUnit      sql.NullInt64   `db:"storage_number_of_unit" json:"storage_number_of_unit" schema:"storage_number_of_unit"`
	StorageNumberOfBag       sql.NullInt64   `db:"storage_number_of_bag" json:"storage_number_of_bag" schema:"storage_number_of_bag"`
	StorageNumberOfCarton    sql.NullInt64   `db:"storage_number_of_carton" json:"storage_number_of_carton" schema:"storage_number_of_carton"`
	Person                   `db:"person" json:"person" schema:"person"`
	Product                  `db:"product" json:"product" schema:"product"`
	StoreLocation            `db:"storelocation" json:"storelocation" schema:"storelocation"`
	UnitQuantity             Unit `db:"unit_quantity" json:"unit_quantity" schema:"unit_quantity"`
	UnitConcentration        Unit `db:"unit_concentration" json:"unit_concentration" schema:"unit_concentration"`
	Supplier                 `db:"supplier" json:"supplier" schema:"supplier"`
	Storage                  *Storage   `db:"storage" json:"storage" schema:"storage"`       // history reference storage
	Borrowing                *Borrowing `db:"borrowing" json:"borrowing" schema:"borrowing"` // not un db but sqlx requires the "db" entry

	// storage history count
	StorageHC int `db:"storage_hc" json:"storage_hc" schema:"storage_hc"` // not in db but sqlx requires the "db" entry
}

// Borrowing represent a storage borrowing
type Borrowing struct {
	BorrowingID      sql.NullInt64  `db:"borrowing_id" json:"borrowing_id" schema:"borrowing_id"`
	BorrowingComment sql.NullString `db:"borrowing_comment" json:"borrowing_comment" schema:"borrowing_comment"`
	Person           *Person        `db:"person" json:"person" schema:"person"` // logged person
	//Storage          `db:"storage" json:"storage" schema:"storage"`
	Borrower *Person `db:"borrower" json:"borrower" schema:"borrower"` // logged person
}

// Permission represent who is able to do what on something
type Permission struct {
	PermissionID       int    `db:"permission_id" json:"permission_id"`
	PermissionPermName string `db:"permission_perm_name" json:"permission_perm_name" schema:"permission_perm_name"` // ex: r
	PermissionItemName string `db:"permission_item_name" json:"permission_item_name" schema:"permission_item_name"` // ex: entity
	PermissionEntityID int    `db:"permission_entity_id" json:"permission_entity_id" schema:"permission_entity_id"` // ex: 8
	Person             `db:"person" json:"person"`
}

// Equal tests the permission equality
func (p1 Permission) Equal(p2 Permission) bool {
	return (p1.PermissionPermName == p2.PermissionPermName &&
		p1.PermissionItemName == p2.PermissionItemName &&
		p1.PermissionEntityID == p2.PermissionEntityID)
}

// Symbol is a product symbol
type Symbol struct {
	SymbolID    int    `db:"symbol_id" json:"symbol_id" schema:"symbol_id"`
	SymbolLabel string `db:"symbol_label" json:"symbol_label" schema:"symbol_label"`
	SymbolImage string `db:"symbol_image" json:"symbol_image" schema:"symbol_image"`
}

// Name is a product name
type Name struct {
	C         int    `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	NameID    int    `db:"name_id" json:"name_id" schema:"name_id"`
	NameLabel string `db:"name_label" json:"name_label" schema:"name_label"`
}

// CasNumber is a product CAS number
type CasNumber struct {
	C              int            `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	CasNumberID    sql.NullInt64  `db:"casnumber_id" json:"casnumber_id" schema:"casnumber_id"`
	CasNumberLabel sql.NullString `db:"casnumber_label" json:"casnumber_label" schema:"casnumber_label"`
	CasNumberCMR   sql.NullString `db:"casnumber_cmr" json:"casnumber_cmr" schema:"casnumber_cmr"`
}

// CeNumber is a product CE number
type CeNumber struct {
	C             int            `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	CeNumberID    sql.NullInt64  `db:"cenumber_id" json:"cenumber_id" schema:"cenumber_id"`
	CeNumberLabel sql.NullString `db:"cenumber_label" json:"cenumber_label" schema:"cenumber_label"`
}

// EmpiricalFormula is a product empirical formula
type EmpiricalFormula struct {
	C                     int            `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	EmpiricalFormulaID    sql.NullInt64  `db:"empiricalformula_id" json:"empiricalformula_id" schema:"empiricalformula_id"`
	EmpiricalFormulaLabel sql.NullString `db:"empiricalformula_label" json:"empiricalformula_label" schema:"empiricalformula_label"`
}

// LinearFormula is a product linear formula
type LinearFormula struct {
	C                  int            `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	LinearFormulaID    sql.NullInt64  `db:"linearformula_id" json:"linearformula_id" schema:"linearformula_id"`
	LinearFormulaLabel sql.NullString `db:"linearformula_label" json:"linearformula_label" schema:"linearformula_label"`
}

// PhysicalState is a product physical state
type PhysicalState struct {
	// nullable values to handle optional Product foreign key (gorilla shema nil values)
	C                  int            `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	PhysicalStateID    sql.NullInt64  `db:"physicalstate_id" json:"physicalstate_id" schema:"physicalstate_id"`
	PhysicalStateLabel sql.NullString `db:"physicalstate_label" json:"physicalstate_label" schema:"physicalstate_label"`
}

// ClassOfCompound is a product class of compound
type ClassOfCompound struct {
	// nullable values to handle optional Product foreign key (gorilla shema nil values)
	C                    int    `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	ClassOfCompoundID    int    `db:"classofcompound_id" json:"classofcompound_id" schema:"classofcompound_id"`
	ClassOfCompoundLabel string `db:"classofcompound_label" json:"classofcompound_label" schema:"classofcompound_label"`
}

// SignalWord is a product signal word
type SignalWord struct {
	// nullable values to handle optional Product foreign key (gorilla shema nil values)
	SignalWordID    sql.NullInt64  `db:"signalword_id" json:"signalword_id" schema:"signalword_id"`
	SignalWordLabel sql.NullString `db:"signalword_label" json:"signalword_label" schema:"signalword_label"`
}

// HazardStatement is a product hazard statement
type HazardStatement struct {
	HazardStatementID        int            `db:"hazardstatement_id" json:"hazardstatement_id" schema:"hazardstatement_id"`
	HazardStatementLabel     string         `db:"hazardstatement_label" json:"hazardstatement_label" schema:"hazardstatement_label"`
	HazardStatementReference string         `db:"hazardstatement_reference" json:"hazardstatement_reference" schema:"hazardstatement_reference"`
	HazardStatementCMR       sql.NullString `db:"hazardstatement_cmr" json:"hazardstatement_cmr" schema:"hazardstatement_cmr"`
}

// PrecautionaryStatement is a product precautionary statement
type PrecautionaryStatement struct {
	PrecautionaryStatementID        int    `db:"precautionarystatement_id" json:"precautionarystatement_id" schema:"precautionarystatement_id"`
	PrecautionaryStatementLabel     string `db:"precautionarystatement_label" json:"precautionarystatement_label" schema:"precautionarystatement_label"`
	PrecautionaryStatementReference string `db:"precautionarystatement_reference" json:"precautionarystatement_reference" schema:"precautionarystatement_reference"`
}

// Product is a chemical product card
type Product struct {
	ProductID              int            `db:"product_id" json:"product_id" schema:"product_id"`
	ProductSpecificity     sql.NullString `db:"product_specificity" json:"product_specificity" schema:"product_specificity"`
	ProductMSDS            sql.NullString `db:"product_msds" json:"product_msds" schema:"product_msds"`
	ProductRestricted      sql.NullBool   `db:"product_restricted" json:"product_restricted" schema:"product_restricted"`
	ProductRadioactive     sql.NullBool   `db:"product_radioactive" json:"product_radioactive" schema:"product_radioactive"`
	ProductThreeDFormula   sql.NullString `db:"product_threedformula" json:"product_threedformula" schema:"product_threedformula"`
	ProductTwoDFormula     sql.NullString `db:"product_twodformula" json:"product_twodformula" schema:"product_twodformula"`
	ProductMolFormula      sql.NullString `db:"product_molformula" json:"product_molformula" schema:"product_molformula"`
	ProductDisposalComment sql.NullString `db:"product_disposalcomment" json:"product_disposalcomment" schema:"product_disposalcomment"`
	ProductRemark          sql.NullString `db:"product_remark" json:"product_remark" schema:"product_remark"`
	ProductTemperature     sql.NullInt64  `db:"product_temperature" json:"product_temperature" schema:"product_temperature"`
	ProductSheet           sql.NullString `db:"product_sheet" json:"product_sheet" schema:"product_sheet"`
	ProductNumberPerCarton sql.NullInt64  `db:"product_number_per_carton" json:"product_number_per_carton" schema:"product_number_per_carton"`
	ProductNumberPerBag    sql.NullInt64  `db:"product_number_per_bag" json:"product_number_per_bag" schema:"product_number_per_bag"`
	EmpiricalFormula       `db:"empiricalformula" json:"empiricalformula" schema:"empiricalformula"`
	LinearFormula          `db:"linearformula" json:"linearformula" schema:"linearformula"`
	PhysicalState          `db:"physicalstate" json:"physicalstate" schema:"physicalstate"`
	SignalWord             `db:"signalword" json:"signalword" schema:"signalword"`
	Person                 `db:"person" json:"person" schema:"person"`
	CasNumber              `db:"casnumber" json:"casnumber" schema:"casnumber"`
	CeNumber               `db:"cenumber" json:"cenumber" schema:"cenumber"`
	Name                   `db:"name" json:"name" schema:"name"`
	ProducerRef            `db:"producerref" json:"producerref" schema:"producerref"`
	Category               `db:"category" json:"category" schema:"category"`
	UnitTemperature        Unit `db:"unit_temperature" json:"unit_temperature" schema:"unit_temperature"`

	ClassOfCompound         []ClassOfCompound        `db:"-" schema:"classofcompound" json:"classofcompound"`
	Synonyms                []Name                   `db:"-" schema:"synonyms" json:"synonyms"`
	Symbols                 []Symbol                 `db:"-" schema:"symbols" json:"symbols"`
	HazardStatements        []HazardStatement        `db:"-" schema:"hazardstatements" json:"hazardstatements"`
	PrecautionaryStatements []PrecautionaryStatement `db:"-" schema:"precautionarystatements" json:"precautionarystatements"`
	SupplierRefs            []SupplierRef            `db:"-" json:"supplierrefs" schema:"supplierrefs"`
	Tags                    []Tag                    `db:"-" json:"tags" schema:"tags"`

	Bookmark *Bookmark `db:"bookmark" json:"bookmark" schema:"bookmark"` // not in db but sqlx requires the "db" entry

	// archived storage count in the logged user entity(ies)
	ProductASC int `db:"product_asc" json:"product_asc" schema:"product_asc"` // not in db but sqlx requires the "db" entry
	// total storage count
	ProductTSC int `db:"product_tsc" json:"product_tsc" schema:"product_tsc"` // not in db but sqlx requires the "db" entry
	// storage count in the logged user entity(ies)
	ProductSC int `db:"product_sc" json:"product_sc" schema:"product_sc"` // not in db but sqlx requires the "db" entry
	// storage barecode concatenation
	ProductSL sql.NullString `db:"product_sl" json:"product_sl" schema:"product_sl"` // not in db but sqlx requires the "db" entry
	// hazard statement CMR concatenation
	HazardStatementCMR sql.NullString `db:"hazardstatement_cmr" json:"hazardstatement_cmr" schema:"hazardstatement_cmr"` // not in db but sqlx requires the "db" entry
}

// Bookmark is a product person bookmark
type Bookmark struct {
	BookmarkID sql.NullInt64 `db:"bookmark_id" json:"bookmark_id" schema:"bookmark_id"`
	Person     `db:"person" json:"person" schema:"person"`
	Product    `db:"product" json:"product" schema:"product"`
}

func (p Product) ProductToStringSlice() []string {
	ret := make([]string, 0)

	ret = append(ret, strconv.Itoa(p.ProductID))

	ret = append(ret, p.NameLabel)
	syn := ""
	for _, s := range p.Synonyms {
		syn += "|" + s.NameLabel
	}
	ret = append(ret, syn)

	ret = append(ret, p.CasNumberLabel.String)
	ret = append(ret, p.CeNumberLabel.String)

	ret = append(ret, p.ProductSpecificity.String)
	ret = append(ret, p.EmpiricalFormulaLabel.String)
	ret = append(ret, p.LinearFormulaLabel.String)
	ret = append(ret, p.ProductThreeDFormula.String)

	ret = append(ret, p.ProductMSDS.String)

	ret = append(ret, p.PhysicalStateLabel.String)

	ret = append(ret, p.SignalWordLabel.String)

	coc := ""
	for _, c := range p.ClassOfCompound {
		coc += "|" + c.ClassOfCompoundLabel
	}
	ret = append(ret, coc)
	sym := ""
	for _, s := range p.Symbols {
		sym += "|" + s.SymbolLabel
	}
	ret = append(ret, sym)
	hs := ""
	for _, h := range p.HazardStatements {
		hs += "|" + h.HazardStatementReference
	}
	ret = append(ret, hs)
	ps := ""
	for _, p := range p.PrecautionaryStatements {
		ps += "|" + p.PrecautionaryStatementReference
	}
	ret = append(ret, ps)

	ret = append(ret, p.ProductRemark.String)
	ret = append(ret, p.ProductDisposalComment.String)

	ret = append(ret, strconv.FormatBool(p.ProductRestricted.Bool))
	ret = append(ret, strconv.FormatBool(p.ProductRadioactive.Bool))

	return ret
}

func (s Storage) StorageToStringSlice() []string {
	ret := make([]string, 0)

	ret = append(ret, strconv.FormatInt(s.StorageID.Int64, 10))
	ret = append(ret, s.Product.Name.NameLabel)
	ret = append(ret, s.Product.CasNumber.CasNumberLabel.String)
	ret = append(ret, s.Product.ProductSpecificity.String)

	ret = append(ret, s.StoreLocation.StoreLocationFullPath)

	ret = append(ret, strconv.FormatFloat(s.StorageQuantity.Float64, 'E', -1, 64))
	ret = append(ret, s.UnitQuantity.UnitLabel.String)
	ret = append(ret, s.StorageBarecode.String)
	ret = append(ret, s.Supplier.SupplierLabel.String)

	ret = append(ret, s.StorageCreationDate.Format("2006-01-02"))
	ret = append(ret, s.StorageModificationDate.Format("2006-01-02"))
	ret = append(ret, s.StorageEntryDate.Time.Format("2006-01-02"))
	ret = append(ret, s.StorageExitDate.Time.Format("2006-01-02"))
	ret = append(ret, s.StorageOpeningDate.Time.Format("2006-01-02"))
	ret = append(ret, s.StorageExpirationDate.Time.Format("2006-01-02"))

	ret = append(ret, s.StorageComment.String)
	ret = append(ret, s.StorageReference.String)
	ret = append(ret, s.StorageBatchNumber.String)

	ret = append(ret, strconv.FormatBool(s.StorageToDestroy.Bool))
	ret = append(ret, strconv.FormatBool(s.StorageArchive.Bool))

	return ret
}

// ProductsToCSV returns a file name of the products prs
// exported into CSV
func ProductsToCSV(prs []Product) string {

	header := []string{"product_id",
		"product_name",
		"product_synonyms",
		"product_cas",
		"product_ce",
		"product_specificity",
		"empirical_formula",
		"linear_formula",
		"3D_formula",
		"MSDS",
		"class_of_compounds",
		"physical_state",
		"signal_word",
		"symbols",
		"hazard_statements",
		"precautionary_statements",
		"remark",
		"disposal_comment",
		"restricted?",
		"radioactive?"}

	// create a temp file
	tmpFile, err := ioutil.TempFile(os.TempDir(), "chimitheque-")
	if err != nil {
		logger.Log.Error("cannot create temporary file", err)
	}
	// creates a csv writer that uses the io buffer
	csvwr := csv.NewWriter(tmpFile)
	// write the header
	_ = csvwr.Write(header)
	for _, p := range prs {
		_ = csvwr.Write(p.ProductToStringSlice())
	}

	csvwr.Flush()
	return strings.Split(tmpFile.Name(), "chimitheque-")[1]
}

// StoragesToCSV returns a file name of the products prs
// exported into CSV
func StoragesToCSV(sts []Storage) (string, error) {

	var (
		err     error
		tmpFile *os.File
	)

	header := []string{"storage_id",
		"product_name",
		"product_casnumber",
		"product_specificity",
		"storelocation",
		"quantity",
		"unit",
		"barecode",
		"supplier",
		"creation_date",
		"modification_date",
		"entry_date",
		"exit_date",
		"opening_date",
		"expiration_date",
		"comment",
		"reference",
		"batch_number",
		"to_destroy?",
		"archive?"}

	// create a temp file
	if tmpFile, err = ioutil.TempFile(os.TempDir(), "chimitheque-"); err != nil {
		logger.Log.Error("cannot create temporary file", err)
		return "", err
	}
	// creates a csv writer that uses the io buffer
	csvwr := csv.NewWriter(tmpFile)
	// write the header
	if err = csvwr.Write(header); err != nil {
		logger.Log.Error("cannot write header", err)
		return "", err
	}

	for _, s := range sts {
		if err = csvwr.Write(s.StorageToStringSlice()); err != nil {
			logger.Log.Error("cannot write entry", err)
			return "", err
		}
	}

	csvwr.Flush()

	return strings.Split(tmpFile.Name(), "chimitheque-")[1], nil
}

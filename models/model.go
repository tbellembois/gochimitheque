package models

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/tbellembois/gochimitheque/helpers"
)

// AppHandlerFunc is an HandlerFunc returning an AppError
type AppHandlerFunc func(http.ResponseWriter, *http.Request) *helpers.AppError

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
}

// Entity represent a department, a laboratory...
type Entity struct {
	EntityID          int      `db:"entity_id" json:"entity_id" schema:"entity_id"`
	EntityName        string   `db:"entity_name" json:"entity_name" schema:"entity_name"`
	EntityDescription string   `db:"entity_description" json:"entity_description" schema:"entity_description"`
	Managers          []Person `db:"-" json:"managers" schema:"managers"`
}

// Person represent a person
type Person struct {
	PersonID       int          `db:"person_id" json:"person_id" schema:"person_id"`
	PersonEmail    string       `db:"person_email" json:"person_email" schema:"person_email"`
	PersonPassword string       `db:"person_password" json:"person_password" schema:"person_password"`
	Permissions    []Permission `db:"-" schema:"permissions"`
	Entities       []Entity     `db:"-" schema:"entities"`
}

// Storage is a product storage in a store location
type Storage struct {
	StorageID           int       `db:"storage_id" json:"storage_id" schema:"storage_id"`
	StorageCreationDate time.Time `db:"storage_creationdate" json:"storage_creationdate" schema:"storage_creationdate"`
	StorageComment      string    `db:"storage_comment" json:"storage_comment" schema:"storage_comment"`
	Person              `db:"person" json:"person" schema:"person"`
	Product             `db:"product" json:"product" schema:"product"`
	StoreLocation       `db:"storelocation" json:"storelocation" schema:"storelocation"`
}

// Permission represent who is able to do what on something
type Permission struct {
	PermissionID       int    `db:"permission_id" json:"permission_id"`
	PermissionPermName string `db:"permission_perm_name" json:"permission_perm_name" schema:"permission_perm_name"` // ex: r
	PermissionItemName string `db:"permission_item_name" json:"permission_item_name" schema:"permission_item_name"` // ex: entity
	PermissionEntityID int    `db:"permission_entity_id" json:"permission_entity_id" schema:"permission_entity_id"` // ex: 8
	Person             `db:"person" json:"person"`
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
	C              int    `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	CasNumberID    int    `db:"casnumber_id" json:"casnumber_id" schema:"casnumber_id"`
	CasNumberLabel string `db:"casnumber_label" json:"casnumber_label" schema:"casnumber_label"`
}

// CeNumber is a product CE number
type CeNumber struct {
	C             int            `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	CeNumberID    sql.NullInt64  `db:"cenumber_id" json:"cenumber_id" schema:"cenumber_id"`
	CeNumberLabel sql.NullString `db:"cenumber_label" json:"cenumber_label" schema:"cenumber_label"`
}

// EmpiricalFormula is a product empirical formula
type EmpiricalFormula struct {
	C                     int    `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	EmpiricalFormulaID    int    `db:"empiricalformula_id" json:"empiricalformula_id" schema:"empiricalformula_id"`
	EmpiricalFormulaLabel string `db:"empiricalformula_label" json:"empiricalformula_label" schema:"empiricalformula_label"`
}

// PhysicalState is a product physical state
type PhysicalState struct {
	// nullable values to handle optional Product foreign key (gorilla shema nil values)
	PhysicalStateID    sql.NullInt64  `db:"physicalstate_id" json:"physicalstate_id" schema:"physicalstate_id"`
	PhysicalStateLabel sql.NullString `db:"physicalstate_label" json:"physicalstate_label" schema:"physicalstate_label"`
}

// ClassOfCompound is a product class of compound
type ClassOfCompound struct {
	// nullable values to handle optional Product foreign key (gorilla shema nil values)
	C                    int            `db:"c" json:"c"` // not stored in db but db:"c" set for sqlx
	ClassOfCompoundID    sql.NullInt64  `db:"classofcompound_id" json:"classofcompound_id" schema:"classofcompound_id"`
	ClassOfCompoundLabel sql.NullString `db:"classofcompound_label" json:"classofcompound_label" schema:"classofcompound_label"`
}

// SignalWord is a product signal word
type SignalWord struct {
	// nullable values to handle optional Product foreign key (gorilla shema nil values)
	SignalWordID    sql.NullInt64  `db:"signalword_id" json:"signalword_id" schema:"signalword_id"`
	SignalWordLabel sql.NullString `db:"signalword_label" json:"signalword_label" schema:"signalword_label"`
}

// HazardStatement is a product hazard statement
type HazardStatement struct {
	HazardStatementID        int    `db:"hazardstatement_id" json:"hazardstatement_id" schema:"hazardstatement_id"`
	HazardStatementLabel     string `db:"hazardstatement_label" json:"hazardstatement_label" schema:"hazardstatement_label"`
	HazardStatementReference string `db:"hazardstatement_reference" json:"hazardstatement_reference" schema:"hazardstatement_reference"`
}

// PrecautionaryStatement is a product precautionary statement
type PrecautionaryStatement struct {
	PrecautionaryStatementID        int    `db:"precautionarystatement_id" json:"precautionarystatement_id" schema:"precautionarystatement_id"`
	PrecautionaryStatementLabel     string `db:"precautionarystatement_label" json:"precautionarystatement_label" schema:"precautionarystatement_label"`
	PrecautionaryStatementReference string `db:"precautionarystatement_reference" json:"precautionarystatement_reference" schema:"precautionarystatement_reference"`
}

// Product is a chemical product card
type Product struct {
	ProductID               int            `db:"product_id" json:"product_id" schema:"product_id"`
	ProductSpecificity      sql.NullString `db:"product_specificity" json:"product_specificity" schema:"product_specificity"`
	ProductMSDS             sql.NullString `db:"product_msds" json:"product_msds" schema:"product_msds"`
	ProductRestricted       sql.NullBool   `db:"product_restricted" json:"product_restricted" schema:"product_restricted"`
	ProductRadioactive      sql.NullBool   `db:"product_radioactive" json:"product_radioactive" schema:"product_radioactive"`
	ProductLinearFormula    sql.NullString `db:"product_linearformula" json:"product_linearformula" schema:"product_linearformula"`
	ProductThreeDFormula    sql.NullString `db:"product_threedformula" json:"product_threedformula" schema:"product_threedformula"`
	ProductDisposalComment  sql.NullString `db:"product_disposalcomment" json:"product_disposalcomment" schema:"product_disposalcomment"`
	ProductRemark           sql.NullString `db:"product_remark" json:"product_remark" schema:"product_remark"`
	EmpiricalFormula        `db:"empiricalformula" json:"empiricalformula" schema:"empiricalformula"`
	PhysicalState           `db:"physicalstate" json:"physicalstate" schema:"physicalstate"`
	SignalWord              `db:"signalword" json:"signalword" schema:"signalword"`
	ClassOfCompound         `db:"classofcompound" json:"classofcompound" schema:"classofcompound"`
	Person                  `db:"person" json:"person" schema:"person"`
	CasNumber               `db:"casnumber" json:"casnumber" schema:"casnumber"`
	CeNumber                `db:"cenumber" json:"cenumber" schema:"cenumber"`
	Name                    `db:"name" json:"name" schema:"name"`
	Synonyms                []Name                   `db:"-" schema:"synonyms" json:"synonyms"`
	Symbols                 []Symbol                 `db:"-" schema:"symbols" json:"symbols"`
	HazardStatements        []HazardStatement        `db:"-" schema:"hazardstatements" json:"hazardstatements"`
	PrecautionaryStatements []PrecautionaryStatement `db:"-" schema:"precautionarystatements" json:"precautionarystatements"`
}

func (p Product) String() string {
	return fmt.Sprintf(`ProductID:%d | 
	ProductSpecificity:%s | 
	EmpiricalFormula:%+v | 
	Person:%+v | 
	CasNumber:%s | 
	CeNumber:%s | 
	Name:%s | 
	Synonyms:%+v | 
	Symbols:%+v |
	DisposalComment:%+v |
	Remark:%+v |`, p.ProductID, p.ProductSpecificity, p.EmpiricalFormula, p.Person, p.CasNumber, p.CeNumber, p.Name, p.Synonyms, p.Symbols, p.ProductDisposalComment, p.ProductRemark)
}

func (p Person) String() string {
	return fmt.Sprintf("PersonEmail: %s", p.PersonEmail)
}

func (s StoreLocation) String() string {
	return fmt.Sprintf("StoreLocationName: %s | StoreLocationCanStore: Valid:%t Bool:%t | StoreLocationColor: %s | Entity: %d | StoreLocation: %v", s.StoreLocationName, s.StoreLocationCanStore.Valid, s.StoreLocationCanStore.Bool, s.StoreLocationColor, s.EntityID, s.StoreLocation)
}

func (p Permission) String() string {
	return fmt.Sprintf("PermissionPermName: %s | PermissionItemName: %s | PermissionEntityID: %d", p.PermissionPermName, p.PermissionItemName, p.PermissionEntityID)
}

func (s Symbol) String() string {
	return fmt.Sprintf("SymbolLabel: %s", s.SymbolLabel)
}

func (c CasNumber) String() string {
	return fmt.Sprintf("CasNumberID: %d | CasNumberLabel: %s", c.CasNumberID, c.CasNumberLabel)
}

func (c CeNumber) String() string {
	i, _ := c.CeNumberID.Value()
	v, _ := c.CeNumberLabel.Value()
	b := c.CeNumberID.Valid
	return fmt.Sprintf("CeNumberID: %d | CeNumberValid: %t | CeNumberLabel: %s", i, b, v)
}

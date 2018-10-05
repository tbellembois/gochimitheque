package models

import (
	"database/sql"
	"fmt"
	"github.com/tbellembois/gochimitheque/helpers"
	"net/http"
	"time"
)

// AppHandlerFunc is an HandlerFunc returning an AppError
type AppHandlerFunc func(http.ResponseWriter, *http.Request) *helpers.AppError

type Model interface {
	GetID() int
}

// StoreLocation is where products are stored in entities
type StoreLocation struct {
	StoreLocationID   int    `db:"storelocation_id" json:"storelocation_id" schema:"storelocation_id"`
	StoreLocationName string `db:"storelocation_name" json:"storelocation_name" schema:"storelocation_name"`
	Entity            `db:"entity" json:"entity" schema:"entity"`
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

func (p Person) GetID() int {
	return p.PersonID
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
	NameID    int    `db:"name_id" json:"name_id" schema:"name_id"`
	NameLabel string `db:"name_label" json:"name_label" schema:"name_label"`
}

func (n Name) GetID() int {
	return n.NameID
}

// CasNumber is a product CAS number
type CasNumber struct {
	CasNumberID    int    `db:"casnumber_id" json:"casnumber_id" schema:"casnumber_id"`
	CasNumberLabel string `db:"casnumber_label" json:"casnumber_label" schema:"casnumber_label"`
}

func (c CasNumber) GetID() int {
	return c.CasNumberID
}

// CeNumber is a product CE number
type CeNumber struct {
	CeNumberID    sql.NullInt64  `db:"cenumber_id" json:"cenumber_id" schema:"cenumber_id"`
	CeNumberLabel sql.NullString `db:"cenumber_label" json:"cenumber_label" schema:"cenumber_label"`
}

func (c CeNumber) GetID() int {
	if c.CeNumberID.Valid {
		return int(c.CeNumberID.Int64)
	}
	return -1
}

// EmpiricalFormula is a product empirical formula
type EmpiricalFormula struct {
	EmpiricalFormulaID    int    `db:"empiricalformula_id" json:"empiricalformula_id" schema:"empiricalformula_id"`
	EmpiricalFormulaLabel string `db:"empiricalformula_label" json:"empiricalformula_label" schema:"empiricalformula_label"`
}

func (e EmpiricalFormula) GetID() int {
	return e.EmpiricalFormulaID
}

// Product is a chemical product card
type Product struct {
	ProductID          int    `db:"product_id" json:"product_id" schema:"product_id"`
	ProductSpecificity string `db:"product_specificity" json:"product_specificity" schema:"product_specificity"`
	EmpiricalFormula   `db:"empiricalformula" json:"empiricalformula" schema:"empiricalformula"`
	Person             `db:"person" json:"person" schema:"person"`
	CasNumber          `db:"casnumber" json:"casnumber" schema:"casnumber"`
	CeNumber           `db:"cenumber" json:"cenumber" schema:"cenumber"`
	Name               `db:"name" json:"name" schema:"name"`
	Synonyms           []Name   `db:"-" schema:"synonyms" json:"synonyms"`
	Symbols            []Symbol `db:"-" schema:"symbols" json:"symbols"`
}

func (p Product) String() string {
	return fmt.Sprintf("ProductID:%d | ProductSpecificity:%s | EmpiricalFormula:%+v | Person:%+v | CasNumber:%s | CeNumber:%s | Name:%s | Synonyms:%+v | Symbols:%+v", p.ProductID, p.ProductSpecificity, p.EmpiricalFormula, p.Person, p.CasNumber, p.CeNumber, p.Name, p.Synonyms, p.Symbols)
}

func (p Person) String() string {
	return fmt.Sprintf("PersonEmail: %s", p.PersonEmail)
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

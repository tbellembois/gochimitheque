package models

import (
	"github.com/tbellembois/gochimitheque/helpers"
	"net/http"
)

// AppHandlerFunc is an HandlerFunc returning an AppError
type AppHandlerFunc func(http.ResponseWriter, *http.Request) *helpers.AppError

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

// CasNumber is a product CAS number
type CasNumber struct {
	CasNumberID    int    `db:"casnumber_id" json:"casnumber_id" schema:"casnumber_id"`
	CasNumberLabel string `db:"casnumber_label" json:"casnumber_label" schema:"casnumber_label"`
}

// Product is a chemical product card
type Product struct {
	ProductID          int    `db:"product_id" json:"product_id" schema:"product_id"`
	ProductSpecificity string `db:"product_specificity" json:"product_specificity" schema:"product_specificity"`
	CasNumber          `db:"casnumber" json:"casnumber" schema:"casnumber"`
	Name               `db:"name" json:"name" schema:"name"`
	Symbols            []Symbol `db:"-" schema:"symbols" json:"symbols" schema:"symbols"`
}

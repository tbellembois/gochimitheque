package models

import (
	"net/http"
	"net/url"
)

// AppHandlerFunc is an HandlerFunc returning an AppError
type AppHandlerFunc func(http.ResponseWriter, *http.Request) *AppError

// AppError is the error type returned by the custom handlers
type AppError struct {
	Error   error
	Message string
	Code    int
}

// ViewContainer is a struct passed to the view
type ViewContainer struct {
	PersonEmail string
	PersonID    int
	URLValues   url.Values
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

// Permission represent who is able to do what on something
type Permission struct {
	PermissionID       int    `db:"permission_id" json:"permission_id"`
	PermissionPermName string `db:"permission_perm_name" json:"permission_perm_name" schema:"permission_perm_name"` // ex: r
	PermissionItemName string `db:"permission_item_name" json:"permission_item_name" schema:"permission_item_name"` // ex: entity
	PermissionEntityID int    `db:"permission_entity_id" json:"permission_entity_id" schema:"permission_entity_id"` // ex: 8
	Person             `db:"person" json:"person"`
}

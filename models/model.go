package models

import (
	"database/sql"
	"net/http"
)

// AppHandlerFunc is an HandlerFunc returning an AppError
type AppHandlerFunc func(http.ResponseWriter, *http.Request) *AppError

// AppError is the error type returned by the custom handlers
type AppError struct {
	Error   error
	Message string
	Code    int
}

// ViewContainer is a struct passed to the view containing the logged user email
// and his permissions
type ViewContainer struct {
	PersonEmail string
	PersonID    int
	//Permissions []Permission
}

// Entity represent a department, a laboratory...
type Entity struct {
	EntityID          int    `db:"entity_id" json:"entity_id" schema:"entity_id"`
	EntityName        string `db:"entity_name" json:"entity_name" schema:"entity_name"`
	EntityDescription string `db:"entity_description" json:"entity_description" schema:"entity_description"`
	// manager
	Person `json:"entity_person_id" schema:"entity_person_id"`
}

// Person represent a person
type Person struct {
	PersonID       int    `db:"person_id" json:"person_id" schema:"person_id"`
	PersonEmail    string `db:"person_email" json:"person_email" schema:"person_email"`
	PersonPassword string `db:"person_password" json:"person_password" schema:"person_password"`
}

// PersonEntities stores the person entities
type PersonEntities struct {
	Person
	Entity
}

// Permission represent who is able to do what on something
type Permission struct {
	PermissionID       int           `db:"permission_id" json:"permission_id" schema:"permission_id"`
	PermissionPermName string        `db:"permission_perm_name" json:"permission_perm_name" schema:"permission_perm_name"` // ex: r
	PermissionItemName string        `db:"permission_item_name" json:"permission_item_name" schema:"permission_item_name"` // ex: entity
	PermissionItemID   sql.NullInt64 `db:"permission_itemid" json:"permission_itemid" schema:"permission_itemid"`          // ex: 8
	Person             `json:"permission_person_id" schema:"permission_person_id"`
}

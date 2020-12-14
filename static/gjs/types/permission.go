package types

import "github.com/gopherjs/gopherjs/js"

type PersonPermission struct {
	*js.Object
	PermissionID       int    `json:"permission_id" js:"permission_id"`
	PermissionPermName string `json:"permission_perm_name" js:"permission_perm_name"` // ex: r
	PermissionItemName string `json:"permission_item_name" js:"permission_item_name"` // ex: entity
	PermissionEntityID int    `json:"permission_entity_id" js:"permission_entity_id"` // ex: 8
}

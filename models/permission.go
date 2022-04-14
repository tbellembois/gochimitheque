package models

// Permission represent who is able to do what on something.
type Permission struct {
	PermissionID       int    `db:"permission_id" json:"permission_id"`
	PermissionPermName string `db:"permission_perm_name" json:"permission_perm_name" schema:"permission_perm_name"` // ex: r
	PermissionItemName string `db:"permission_item_name" json:"permission_item_name" schema:"permission_item_name"` // ex: entity
	PermissionEntityID int    `db:"permission_entity_id" json:"permission_entity_id" schema:"permission_entity_id"` // ex: 8
	Person             `db:"person" json:"person"`
}

// Equal tests the permission equality.
func (p1 Permission) Equal(p2 Permission) bool {
	return (p1.PermissionPermName == p2.PermissionPermName &&
		p1.PermissionItemName == p2.PermissionItemName &&
		p1.PermissionEntityID == p2.PermissionEntityID)
}

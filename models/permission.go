package models

// Permission represent who is able to do what on something.
type Permission struct {
	PermissionID     int    `db:"permission_id" json:"permission_id"`
	PermissionName   string `db:"permission_name" json:"permission_name" schema:"permission_name"`       // ex: r
	PermissionItem   string `db:"permission_item" json:"permission_item" schema:"permission_item"`       // ex: entity
	PermissionEntity int64  `db:"permission_entity" json:"permission_entity" schema:"permission_entity"` // ex: 8
	Person           `db:"person" json:"person"`
}

// Equal tests the permission equality.
func (p1 Permission) Equal(p2 Permission) bool {
	return (p1.PermissionName == p2.PermissionName &&
		p1.PermissionItem == p2.PermissionItem &&
		p1.PermissionEntity == p2.PermissionEntity)
}

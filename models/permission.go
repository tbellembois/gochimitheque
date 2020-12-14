package models

type PermKey struct {
	View string
	Item string
	Verb string
	Id   string
}
type PermValue struct {
	Type string
	Item string
	Id   string
}

type CasbinJSON struct {
	PType string
	V0    string
	V1    string
	V2    string
	V3    string
}

package helpers

const (
	PermTypeRead  = "r"
	PermTypeWrite = "w"
	PermTypeAll   = "all"
	PermTypeNA    = "na"

	PermOnAll    = "all"
	PermOnAny    = "any"
	PermOnEntity = "entity"

	ViewRead   = "v"
	ViewCreate = "vc"

	VerbGet    = "GET"
	VerbPut    = "PUT"
	VerbPost   = "POST"
	VerbDelete = "DELETE"

	ItemId   = "id"
	ItemAll  = "-1"
	ItemAny  = "-2"
	ItemNone = ""
)

type PermMatrixRow struct {
	View string
	Item string
}
type PermMatrixCol struct {
	Verb   string
	ItemId string
}

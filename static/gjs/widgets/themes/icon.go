package themes

import "fmt"

type Icon interface {
	ToString() string
}

type MDIcon interface {
	Icon
}

// Material Design icons
// https://materialdesignicons.com/
type mdiIcon struct {
	Face IconFace
	Size IconSize
}

type IconClass string
type IconFace string
type IconSize string

const (
	MDI IconClass = "mdi"

	MDI_STORELOCATION IconFace = "mdi-docker"
	MDI_PEOPLE        IconFace = "mdi-account-group"
	MDI_EDIT          IconFace = "mdi-border-color"
	MDI_DELETE        IconFace = "mdi-delete"
	MDI_NONE          IconFace = "mdi-border-none-variant"
	MDI_ERROR         IconFace = "mdi-alert-circle-outline"

	MDI_24PX IconSize = "mdi-24px"
)

func NewMdiIcon(face IconFace, size IconSize) Icon {

	if face == "" {
		face = MDI_NONE
	}
	if size == "" {
		size = MDI_24PX
	}

	return mdiIcon{Face: face, Size: size}

}

func (i mdiIcon) ToString() string {
	return fmt.Sprintf("%s %s %s", MDI, i.Face, i.Size)
}

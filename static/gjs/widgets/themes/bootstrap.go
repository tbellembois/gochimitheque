package themes

import "fmt"

type BSClass string

const (
	BS_BTN BSClass = "btn"

	BS_BNT_LINK BSClass = "btn-link"

	BS_BNT_SM BSClass = "btn-sm"
)

func (b BSClass) ToString() string {
	return fmt.Sprintf("%s", b)
}

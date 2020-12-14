package widgets

import (
	. "github.com/tbellembois/gochimitheque/static/gjs/globals"
	"github.com/tbellembois/gochimitheque/static/gjs/widgets/themes"

	"honnef.co/go/js/dom"
)

type icon struct {
	Widget
}

type IconAttributes struct {
	BaseAttributes
	Text string
	Icon themes.MDIcon
}

func NewIcon(args IconAttributes) *icon {

	i := &icon{}
	i.HTMLElement = Doc.CreateElement("span").(*dom.HTMLSpanElement)

	i.SetInnerHTML(args.Text)
	// Appending mateial design icon to classes.
	args.BaseAttributes.Classes = append(args.BaseAttributes.Classes, args.Icon.ToString())
	args.BaseAttributes.Visible = true

	i.SetBaseAttributes(args.BaseAttributes)

	return i

}

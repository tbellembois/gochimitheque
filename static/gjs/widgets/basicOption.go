package widgets

import (
	. "github.com/tbellembois/gochimitheque/static/gjs/globals"
	"honnef.co/go/js/dom"
)

type option struct {
	Widget
}

type OptionAttributes struct {
	BaseAttributes
	Text            string
	Value           string
	DefaultSelected bool
	Selected        bool
}

func NewOption(args OptionAttributes) *option {

	o := &option{}

	htmlElement := Doc.CreateElement("option").(*dom.HTMLOptionElement)
	htmlElement.SetTextContent(args.Text)
	htmlElement.SetAttribute("value", args.Value)
	htmlElement.Selected = args.Selected
	htmlElement.DefaultSelected = args.DefaultSelected

	o.HTMLElement = htmlElement

	o.SetBaseAttributes(args.BaseAttributes)

	return o

}

package widgets

import (
	. "github.com/tbellembois/gochimitheque/static/gjs/globals"

	"honnef.co/go/js/dom"
)

type button struct {
	Widget
}

type ButtonAttributes struct {
	BaseAttributes
	Label string
	Title string
}

func NewButton(args ButtonAttributes) *button {

	b := &button{}
	b.HTMLElement = Doc.CreateElement("button").(*dom.HTMLButtonElement)

	b.SetTitle(args.Title)
	if args.Label != "" {
		b.SetAttribute("label", args.Label)
	}

	b.SetBaseAttributes(args.BaseAttributes)

	return b

}

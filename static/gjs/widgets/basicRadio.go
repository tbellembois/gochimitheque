package widgets

import (
	. "github.com/tbellembois/gochimitheque/static/gjs/globals"
	"honnef.co/go/js/dom"
)

type radio struct {
	Widget
}

type RadioAttributes struct {
	BaseAttributes
	Checked bool
}

func NewRadio(args RadioAttributes) *radio {

	r := &radio{}

	htmlElement := Doc.CreateElement("input").(*dom.HTMLInputElement)
	htmlElement.SetAttribute("type", "radio")
	htmlElement.Checked = args.Checked

	r.HTMLElement = htmlElement

	r.SetBaseAttributes(args.BaseAttributes)

	return r

}

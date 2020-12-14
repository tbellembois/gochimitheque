package widgets

import (
	. "github.com/tbellembois/gochimitheque/static/gjs/globals"
	"honnef.co/go/js/dom"
)

type div struct {
	Widget
}

type DivAttributes struct {
	BaseAttributes
}

func NewDiv(args DivAttributes) *div {

	d := &div{}
	d.HTMLElement = Doc.CreateElement("div").(*dom.HTMLDivElement)

	d.SetBaseAttributes(args.BaseAttributes)

	return d

}

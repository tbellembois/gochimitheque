package widgets

import (
	. "github.com/tbellembois/gochimitheque/static/gjs/globals"
	"honnef.co/go/js/dom"
)

type ul struct {
	Widget
}

type UlAttributes struct {
	BaseAttributes
}

func NewUl(args UlAttributes) *ul {

	o := &ul{}
	o.HTMLElement = Doc.CreateElement("ul").(*dom.HTMLUListElement)

	o.SetBaseAttributes(args.BaseAttributes)

	return o

}

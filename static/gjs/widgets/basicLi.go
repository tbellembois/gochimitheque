package widgets

import (
	. "github.com/tbellembois/gochimitheque/static/gjs/globals"
	"honnef.co/go/js/dom"
)

type li struct {
	Widget
}

type LiAttributes struct {
	BaseAttributes
	Text string
}

func NewLi(args LiAttributes) *li {

	o := &li{}
	o.HTMLElement = Doc.CreateElement("li").(*dom.HTMLLIElement)

	o.SetBaseAttributes(args.BaseAttributes)

	return o

}

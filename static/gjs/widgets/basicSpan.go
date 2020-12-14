package widgets

import (
	. "github.com/tbellembois/gochimitheque/static/gjs/globals"

	"honnef.co/go/js/dom"
)

type span struct {
	Widget
}

type SpanAttributes struct {
	BaseAttributes
	Text string
}

func NewSpan(args SpanAttributes) *span {

	s := &span{}
	s.HTMLElement = Doc.CreateElement("span").(*dom.HTMLSpanElement)

	s.SetInnerHTML(args.Text)

	s.SetBaseAttributes(args.BaseAttributes)

	return s

}

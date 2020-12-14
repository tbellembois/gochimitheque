package widgets

import (
	. "github.com/tbellembois/gochimitheque/static/gjs/globals"
	"honnef.co/go/js/dom"
)

// Title returns a page title with a customized icon
func Title(msgText string, msgType string) *dom.HTMLDivElement {
	t := Doc.CreateElement("div").(*dom.HTMLDivElement)
	s := Doc.CreateElement("span").(*dom.HTMLSpanElement)
	icon := Doc.CreateElement("span").(*dom.HTMLSpanElement)
	t.Class().SetString("mt-md-3 mb-md-3 row text-right")
	s.Class().SetString("col-sm-11 align-bottom")
	s.SetTextContent(msgText)

	switch msgType {
	case "history":
		icon.Class().SetString("mdi mdi-24px mdi-alarm")
	case "bookmark":
		icon.Class().SetString("mdi mdi-24px mdi-bookmark")
	case "entity":
		icon.Class().SetString("mdi mdi-24px mdi-store")
	case "storelocation":
		icon.Class().SetString("mdi mdi-24px mdi-Docker")
	case "product":
		icon.Class().SetString("mdi mdi-24px mdi-tag")
	case "storage":
		icon.Class().SetString("mdi mdi-24px mdi-cube-unfolded")
	default:
		icon.Class().SetString("mdi mdi-24px mdi-menu-right-outline")
	}

	t.AppendChild(s)
	t.AppendChild(icon)

	return t
}

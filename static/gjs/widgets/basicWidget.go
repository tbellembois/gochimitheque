package widgets

import (
	"strings"

	"honnef.co/go/js/dom"
)

type Widget struct {
	dom.HTMLElement
}

type BaseAttributes struct {
	Id         string
	Classes    []string
	Attributes map[string]string
	Visible    bool
}

func (w Widget) SetBaseAttributes(a BaseAttributes) {

	if a.Id != "" {
		w.SetID(a.Id)
	}
	if len(a.Classes) > 0 {
		w.Class().SetString(strings.Join(a.Classes, " "))
	}
	for attributeName, attributeValue := range a.Attributes {
		w.SetAttribute(attributeName, attributeValue)
	}
	if !a.Visible {
		w.SetAttribute("style", "display: none;")
	}

}

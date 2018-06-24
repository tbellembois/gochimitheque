package main

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

var (
	window   dom.Window
	document dom.Document

	// database tables
	tableitems = [6]string{
		"product",
		"rproduct",
		"storage",
		"astorage",
		"classofcompounds",
		"supplier"}
)

func init() {
	window = dom.GetWindow()
	document = window.Document()
}

// Permission represent who is able to do what on something
type Permission struct {
	PermissionID       int    `db:"permission_id" json:"permission_id" schema:"permission_id"`
	PermissionPermName string `db:"permission_perm_name" json:"permission_perm_name" schema:"permission_perm_name"` // ex: r
	PermissionItemName string `db:"permission_item_name" json:"permission_item_name" schema:"permission_item_name"` // ex: entity
	PermissionItemID   int    `db:"permission_itemid" json:"permission_itemid" schema:"permission_itemid"`          // ex: 8
}

func Test(params []interface{}) {
	println(params)
	for _, p := range params {
		println(p)
	}
}

// BuildInlineRadioElement return a radio inline block such as:
//
// <div class="form-check form-check-inline">
//   <input class="form-check-input" type="radio" name="inlineRadioOptions" id="inlineRadio1" value="option1">
//   <label class="form-check-label" for="inlineRadio1">1</label>
// </div>
//
// inputattr much contain at least id, name, label and value
func BuildInlineRadioElement(inputattr map[string]string) *dom.HTMLDivElement {
	// main div
	maindiv := document.CreateElement("div").(*dom.HTMLDivElement)
	maindiv.SetClass("form-check form-check-inline")
	// input
	input := document.CreateElement("input").(*dom.HTMLInputElement)
	input.SetAttribute("type", "radio")
	input.SetClass("form-check-input")
	// ...setting up additional attributes
	for a := range inputattr {
		input.SetAttribute(a, inputattr[a])
	}
	// label
	label := document.CreateElement("label").(*dom.HTMLLabelElement)
	label.SetClass("form-check-label")
	label.SetAttribute("for", inputattr["id"])
	label.SetInnerHTML(inputattr["label"])

	// building the final result
	maindiv.AppendChild(input)
	maindiv.AppendChild(label)
	return maindiv
}

// BuildPermissionWidget return a widget to setup people permissions
func BuildPermissionWidget(persID int) *dom.HTMLDivElement {

	var widgetdiv *dom.HTMLDivElement
	// create main widget div
	widgetdiv = document.CreateElement("div").(*dom.HTMLDivElement)
	widgetdiv.SetID(fmt.Sprintf("perm%d", persID))

	for _, i := range tableitems {
		// building main row
		mainrowdiv := document.CreateElement("div").(*dom.HTMLDivElement)
		mainrowdiv.SetClass("row")
		// building first col for table name
		label := document.CreateElement("div").(*dom.HTMLDivElement)
		label.SetClass("alert alert-primary")
		label.SetInnerHTML(i)
		firstcoldiv := document.CreateElement("div").(*dom.HTMLDivElement)
		firstcoldiv.SetClass("col-sm-6")
		firstcoldiv.AppendChild(label)
		// building second col for radios
		noneradioattrs := map[string]string{"id": i, "name": fmt.Sprintf("perm%s", i), "value": "none", "label": "no permission"}
		readradioattrs := map[string]string{"id": i, "name": fmt.Sprintf("perm%s", i), "value": "r", "label": "read"}
		writeradioattrs := map[string]string{"id": i, "name": fmt.Sprintf("perm%s", i), "value": "w", "label": "read/write"}
		secondcoldiv := document.CreateElement("div").(*dom.HTMLDivElement)
		secondcoldiv.SetClass("col-sm-6")
		secondcoldiv.AppendChild(BuildInlineRadioElement(noneradioattrs))
		secondcoldiv.AppendChild(BuildInlineRadioElement(readradioattrs))
		secondcoldiv.AppendChild(BuildInlineRadioElement(writeradioattrs))

		// appending to final div
		mainrowdiv.AppendChild(firstcoldiv)
		mainrowdiv.AppendChild(secondcoldiv)
		widgetdiv.AppendChild(mainrowdiv)
	}

	return widgetdiv
}

func main() {

	// exporting functions to be called from other JS files
	js.Global.Set("global", map[string]interface{}{
		"buildPermissionWidget": BuildPermissionWidget,
		"test":                  Test,
	})

}

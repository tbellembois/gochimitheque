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
	permitems = [2]string{
		"products",
		"storages"}
)

func init() {
	window = dom.GetWindow()
	document = window.Document()
}

func PopulatePermissionWidget(params []*js.Object) {
	// unchecking all permissions
	for _, e := range document.GetElementsByClassName("perm") {
		e.(*dom.HTMLInputElement).RemoveAttribute("checked")
	}

	// setting all permissions at none by defaut
	for _, e := range document.GetElementsByClassName("permn") {
		e.(*dom.HTMLInputElement).Checked = true
	}

	// then setting up new permissions
	for _, p := range params {
		pitemname := p.Get("permission_item_name").String()
		ppermname := p.Get("permission_perm_name").String()
		pentityid := p.Get("permission_entity_id").String()

		// println("---")
		// println(pitemname)
		// println(ppermname)
		// println(pentityid)

		switch pitemname {
		case "products":
			switch ppermname {
			case "w", "all":
				if pentityid == "-1" {
					for _, e := range document.GetElementsByClassName("permwproducts") {
						e.(*dom.HTMLInputElement).Checked = true
					}
				} else {
					document.GetElementByID("permw" + pitemname + pentityid).(*dom.HTMLInputElement).Checked = true
				}
			case "r":
				if pentityid == "-1" {
					for _, e := range document.GetElementsByClassName("permrproducts") {
						e.(*dom.HTMLInputElement).Checked = true
					}
				} else {
					// should never happen
					// permissions on products are not related to an entity id
					document.GetElementByID("permr" + pitemname + pentityid).(*dom.HTMLInputElement).Checked = true
				}
			}
		case "storages":
			switch ppermname {
			case "w", "all":
				if pentityid == "-1" {
					for _, e := range document.GetElementsByClassName("permwstorages") {
						e.(*dom.HTMLInputElement).Checked = true
					}
				} else {
					document.GetElementByID("permw" + pitemname + pentityid).(*dom.HTMLInputElement).Checked = true
				}
			case "r":
				if pentityid == "-1" {
					for _, e := range document.GetElementsByClassName("permrstorages") {
						e.(*dom.HTMLInputElement).Checked = true
					}
				} else {
					document.GetElementByID("permr" + pitemname + pentityid).(*dom.HTMLInputElement).Checked = true
				}
			}
		case "all":
			switch ppermname {
			case "w", "all":
				if pentityid == "-1" {
					// super admin (if "all")
					for _, e := range document.GetElementsByClassName("permw") {
						e.(*dom.HTMLInputElement).Checked = true
					}
				} else {
					// manager (if "all")
					for _, e := range document.GetElementsByClassName("permw") {
						e.(*dom.HTMLInputElement).Checked = true
					}
				}
			case "r":
				for _, e := range document.GetElementsByClassName("permr") {
					e.(*dom.HTMLInputElement).Checked = true
				}
			}
		}
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
func BuildPermissionWidget(entityID int, entityName string) *dom.HTMLDivElement {

	var widgetdiv *dom.HTMLDivElement
	// create main widget div
	widgetdiv = document.CreateElement("div").(*dom.HTMLDivElement)
	widgetdiv.SetID(fmt.Sprintf("perm%d", entityID))
	widgetdiv.SetClass("col-sm-12")
	title := document.CreateElement("label").(*dom.HTMLLabelElement)
	title.SetInnerHTML(entityName)

	widgetdiv.AppendChild(title)
	for _, i := range permitems {
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
		noneradioattrs := map[string]string{
			"id":        fmt.Sprintf("permn%s%d", i, entityID),
			"name":      fmt.Sprintf("perm%s%d", i, entityID),
			"value":     "none",
			"label":     "_",
			"perm_name": "n",
			"item_name": fmt.Sprintf("%s", i),
			"entity_id": fmt.Sprintf("%d", entityID),
			"class":     fmt.Sprintf("perm permn permn%s", i)}
		readradioattrs := map[string]string{
			"id":        fmt.Sprintf("permr%s%d", i, entityID),
			"name":      fmt.Sprintf("perm%s%d", i, entityID),
			"value":     "r",
			"label":     "r",
			"perm_name": "r",
			"item_name": fmt.Sprintf("%s", i),
			"entity_id": fmt.Sprintf("%d", entityID),
			"class":     fmt.Sprintf("perm permr permr%s", i)}
		writeradioattrs := map[string]string{
			"id":        fmt.Sprintf("permw%s%d", i, entityID),
			"name":      fmt.Sprintf("perm%s%d", i, entityID),
			"value":     "w",
			"label":     "rw",
			"perm_name": "w",
			"item_name": fmt.Sprintf("%s", i),
			"entity_id": fmt.Sprintf("%d", entityID),
			"class":     fmt.Sprintf("perm permw permw%s", i)}
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

func main2() {

	// exporting functions to be called from other JS files
	js.Global.Set("global", map[string]interface{}{
		"buildPermissionWidget":    BuildPermissionWidget,
		"populatePermissionWidget": PopulatePermissionWidget,
	})

}

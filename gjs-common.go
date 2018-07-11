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
	tableitems = [2]string{
		"product",
		"storage"}
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
	for _, e := range tableitems {
		//document.GetElementByID("perm"+e).SetAttribute("checked", "checked")
		document.GetElementByID("perm" + e).(*dom.HTMLInputElement).Checked = true
	}

	// then setting up new permissions
	for _, p := range params {
		pitemname := p.Get("permission_item_name").String()
		ppermname := p.Get("permission_perm_name").String()
		//pitemid := p.Get("permission_entityid").Get("Int64").Int64()

		switch pitemname {
		case "product", "storage":
			switch ppermname {
			case "w", "all":
				//document.GetElementByID("perm"+pitemname+"rw").SetAttribute("checked", "checked")
				document.GetElementByID("perm" + pitemname + "rw").(*dom.HTMLInputElement).Checked = true
			case "r":
				//document.GetElementByID("perm"+pitemname+"r").SetAttribute("checked", "checked")
				document.GetElementByID("perm" + pitemname + "r").(*dom.HTMLInputElement).Checked = true
			}
		case "all":
			switch ppermname {
			case "w", "all":
				for _, e := range tableitems {
					//document.GetElementByID("perm"+e+"rw").SetAttribute("checked", "checked")
					document.GetElementByID("perm" + e + "rw").(*dom.HTMLInputElement).Checked = true
				}
			case "r":
				for _, e := range tableitems {
					//document.GetElementByID("perm"+e+"r").SetAttribute("checked", "checked")
					document.GetElementByID("perm" + e + "r").(*dom.HTMLInputElement).Checked = true
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

// func main() {

// 	// exporting functions to be called from other JS files
// 	js.Global.Set("global", map[string]interface{}{
// 		"buildPermissionWidget":    BuildPermissionWidget,
// 		"populatePermissionWidget": PopulatePermissionWidget,
// 	})

// }

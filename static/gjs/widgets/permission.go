package widgets

import (
	"fmt"
	"strconv"

	. "github.com/tbellembois/gochimitheque/static/gjs/globals"
	"github.com/tbellembois/gochimitheque/static/gjs/types"
	"honnef.co/go/js/dom"
)

// Permission return a widget to setup people permissions
func Permission(entityID int, entityName string, ismanager bool) *dom.HTMLDivElement {
	// create main widget div
	widgetdiv := Doc.CreateElement("div").(*dom.HTMLDivElement)
	widgetdiv.SetID(fmt.Sprintf("perm%d", entityID))
	widgetdiv.SetClass("col-sm-12")
	title := Doc.CreateElement("div").(*dom.HTMLDivElement)
	title.SetClass("d-flex")
	title.SetInnerHTML("<span class='mdi mdi-store mdi-24px'/>" + entityName)

	widgetdiv.AppendChild(title)

	if ismanager {
		s := Doc.CreateElement("span").(*dom.HTMLSpanElement)
		s.SetClass("mdi mdi-36px mdi-account-star")
		s.SetAttribute("title", "manager")
		coldiv := Doc.CreateElement("div").(*dom.HTMLDivElement)
		coldiv.SetClass("col-sm-2")
		coldiv.AppendChild(s)

		widgetdiv.AppendChild(coldiv)
		return widgetdiv
	}

	for _, i := range PermItems {
		// products permissions widget is static
		if i != "products" && i != "rproducts" {
			//println(i)
			// building main row
			mainrowdiv := Doc.CreateElement("div").(*dom.HTMLDivElement)
			mainrowdiv.SetClass("form-group row d-flex")
			// building first col for table name
			label := Doc.CreateElement("div").(*dom.HTMLDivElement)
			label.SetClass("iconlabel text-right")
			label.SetInnerHTML(i)
			firstcoldiv := Doc.CreateElement("div").(*dom.HTMLDivElement)
			firstcoldiv.SetClass("col-sm-2")
			firstcoldiv.AppendChild(label)
			// building second col for radios
			noneradioattrs := map[string]string{
				"id":        fmt.Sprintf("permn%s%d", i, entityID),
				"name":      fmt.Sprintf("perm%s%d", i, entityID),
				"value":     "none",
				"checked":   "checked",
				"label":     fmt.Sprintf("<label class=\"form-check-label ml-sm-1 pr-sm-1 pl-sm-1 text-secondary border border-secondary rounded\" for=\"permn%s%d\"><span class=\"mdi mdi-close\"></span></label>", i, entityID),
				"perm_name": "n",
				"item_name": i,
				"entity_id": fmt.Sprintf("%d", entityID),
				"class":     fmt.Sprintf("perm permn permn%s", i)}
			readradioattrs := map[string]string{
				"id":        fmt.Sprintf("permr%s%d", i, entityID),
				"name":      fmt.Sprintf("perm%s%d", i, entityID),
				"value":     "r",
				"label":     fmt.Sprintf("<label class=\"form-check-label ml-sm-1 pr-sm-1 pl-sm-1 text-secondary border border-secondary rounded\" for=\"permn%s%d\"><span class=\"mdi mdi-eye\"></span></label>", i, entityID),
				"perm_name": "r",
				"item_name": i,
				"entity_id": fmt.Sprintf("%d", entityID),
				"class":     fmt.Sprintf("perm permr permr%s", i)}
			writeradioattrs := map[string]string{
				"id":        fmt.Sprintf("permw%s%d", i, entityID),
				"name":      fmt.Sprintf("perm%s%d", i, entityID),
				"value":     "w",
				"label":     fmt.Sprintf("<label class=\"form-check-label ml-sm-1 pr-sm-1 pl-sm-1 text-secondary border border-secondary rounded\" for=\"permn%s%d\"><span class=\"mdi mdi-eye\"></span><span class=\"mdi mdi-creation\"></span><span class=\"mdi mdi-border-color\"></span><span class=\"mdi mdi-delete\"></span></label>", i, entityID),
				"perm_name": "w",
				"item_name": i,
				"entity_id": fmt.Sprintf("%d", entityID),
				"class":     fmt.Sprintf("perm permw permw%s", i)}
			secondcoldiv := Doc.CreateElement("div").(*dom.HTMLDivElement)
			secondcoldiv.SetClass("col-sm-4")
			secondcoldiv.AppendChild(InlineRadio(noneradioattrs))
			secondcoldiv.AppendChild(InlineRadio(readradioattrs))
			secondcoldiv.AppendChild(InlineRadio(writeradioattrs))

			// appending to final div
			mainrowdiv.AppendChild(firstcoldiv)
			mainrowdiv.AppendChild(secondcoldiv)
			widgetdiv.AppendChild(mainrowdiv)
		}
	}
	return widgetdiv
}

// PopulatePermission checks the permissions checkboxes in the person edition page
func PopulatePermission(permissions []types.PersonPermission) {
	// unchecking all permissions
	for _, e := range Doc.GetElementsByClassName("perm") {
		e.(*dom.HTMLInputElement).RemoveAttribute("checked")
	}

	// setting all permissions at none by defaut
	for _, e := range Doc.GetElementsByClassName("permn") {
		e.(*dom.HTMLInputElement).Checked = true
	}

	// then setting up new permissions
	for _, p := range permissions {

		pentityid := strconv.Itoa(p.PermissionEntityID)

		switch p.PermissionItemName {
		case "products":
			switch p.PermissionPermName {
			case "w", "all":
				if pentityid == "-1" {
					for _, e := range Doc.GetElementsByClassName("permwproducts") {
						e.(*dom.HTMLInputElement).Checked = true
					}
				} else {
					Doc.GetElementByID("permw" + p.PermissionItemName + pentityid).(*dom.HTMLInputElement).Checked = true
				}
			case "r":
				if pentityid == "-1" {
					for _, e := range Doc.GetElementsByClassName("permrproducts") {
						// avoid selecting r if w is already selected
						eid := e.GetAttribute("entity_id")
						if !Doc.GetElementByID("permw" + p.PermissionItemName + eid).(*dom.HTMLInputElement).Checked {
							e.(*dom.HTMLInputElement).Checked = true
						}
					}
				} else {
					// avoid selecting r if w is already selected
					if !Doc.GetElementByID("permw" + p.PermissionItemName + pentityid).(*dom.HTMLInputElement).Checked {
						Doc.GetElementByID("permr" + p.PermissionItemName + pentityid).(*dom.HTMLInputElement).Checked = true
					}
				}
			case "n":
				if pentityid == "-1" {
					for _, e := range Doc.GetElementsByClassName("permnproducts") {
						e.(*dom.HTMLInputElement).Checked = true
					}
				} else {
					Doc.GetElementByID("permn" + p.PermissionItemName + pentityid).(*dom.HTMLInputElement).Checked = true
				}
			}
		case "rproducts":
			switch p.PermissionPermName {
			case "w", "all":
				if pentityid == "-1" {
					for _, e := range Doc.GetElementsByClassName("permwrproducts") {
						e.(*dom.HTMLInputElement).Checked = true
					}
				} else {
					Doc.GetElementByID("permw" + p.PermissionItemName + pentityid).(*dom.HTMLInputElement).Checked = true
				}
			case "r":
				if pentityid == "-1" {
					for _, e := range Doc.GetElementsByClassName("permrrproducts") {
						// avoid selecting r if w is already selected
						eid := e.GetAttribute("entity_id")
						if !Doc.GetElementByID("permw" + p.PermissionItemName + eid).(*dom.HTMLInputElement).Checked {
							e.(*dom.HTMLInputElement).Checked = true
						}
					}
				} else {
					// avoid selecting r if w is already selected
					if !Doc.GetElementByID("permw" + p.PermissionItemName + pentityid).(*dom.HTMLInputElement).Checked {
						Doc.GetElementByID("permr" + p.PermissionItemName + pentityid).(*dom.HTMLInputElement).Checked = true
					}
				}
			case "n":
				if pentityid == "-1" {
					for _, e := range Doc.GetElementsByClassName("permnrproducts") {
						e.(*dom.HTMLInputElement).Checked = true
					}
				} else {
					Doc.GetElementByID("permn" + p.PermissionItemName + pentityid).(*dom.HTMLInputElement).Checked = true
				}
			}
		case "storages":
			switch p.PermissionPermName {
			case "w", "all":
				if pentityid == "-1" {
					for _, e := range Doc.GetElementsByClassName("permwstorages") {
						e.(*dom.HTMLInputElement).Checked = true
					}
				} else {
					Doc.GetElementByID("permw" + p.PermissionItemName + pentityid).(*dom.HTMLInputElement).Checked = true
				}
			case "r":
				if pentityid == "-1" {
					for _, e := range Doc.GetElementsByClassName("permrstorages") {
						// avoid selecting r if w is already selected
						eid := e.GetAttribute("entity_id")
						if !Doc.GetElementByID("permw" + p.PermissionItemName + eid).(*dom.HTMLInputElement).Checked {
							e.(*dom.HTMLInputElement).Checked = true
						}
					}
				} else {
					// avoid selecting r if w is already selected
					if !Doc.GetElementByID("permw" + p.PermissionItemName + pentityid).(*dom.HTMLInputElement).Checked {
						Doc.GetElementByID("permr" + p.PermissionItemName + pentityid).(*dom.HTMLInputElement).Checked = true
					}
				}
			}
		case "all":
			switch p.PermissionPermName {
			case "w", "all":
				if pentityid == "-1" {
					// super admin (if "all")
					for _, e := range Doc.GetElementsByClassName("permw") {
						e.(*dom.HTMLInputElement).Checked = true
					}
				}
			case "r":
				for _, e := range Doc.GetElementsByClassName("permr") {
					e.(*dom.HTMLInputElement).Checked = true
				}
			}
		}
	}

}

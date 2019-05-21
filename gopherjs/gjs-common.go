package main

import (
	"fmt"
	"regexp"

	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

var (
	window   dom.Window
	document dom.Document

	// permissions
	permitems = [3]string{
		"rproducts",
		"products",
		//"storelocations",
		"storages"}
)

func init() {
	window = dom.GetWindow()
	document = window.Document()
}

// CreateTitle return a JDiv title wrapped in a js.Object
func CreateTitle(msgText string, msgType string) *dom.HTMLDivElement {
	t := document.CreateElement("div").(*dom.HTMLDivElement)
	s := document.CreateElement("span").(*dom.HTMLSpanElement)
	sp := document.CreateElement("span").(*dom.HTMLSpanElement)
	t.Class().SetString("mt-md-3 mb-md-3 row")
	s.Class().SetString("col-sm-11 align-bottom")
	s.SetTextContent(msgText)

	switch msgType {
	case "history":
		sp.Class().SetString("mdi mdi-24px mdi-alarm")
		break
	case "bookmark":
		sp.Class().SetString("mdi mdi-24px mdi-bookmark")
		break
	case "entity":
		sp.Class().SetString("mdi mdi-24px mdi-store")
		break
	case "storelocation":
		sp.Class().SetString("mdi mdi-24px mdi-docker")
		break
	case "product":
		sp.Class().SetString("mdi mdi-24px mdi-tag")
		break
	case "storage":
		sp.Class().SetString("mdi mdi-24px mdi-cube-unfolded")
		break
	default:
		sp.Class().SetString("mdi mdi-24px mdi-menu-right-outline")
	}

	t.AppendChild(sp)
	t.AppendChild(s)

	return t
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
						// avoid selecting r if w is already selected
						eid := e.GetAttribute("entity_id")
						if !document.GetElementByID("permw" + pitemname + eid).(*dom.HTMLInputElement).Checked {
							e.(*dom.HTMLInputElement).Checked = true
						}
					}
				} else {
					// avoid selecting r if w is already selected
					if document.GetElementByID("permw"+pitemname+pentityid).(*dom.HTMLInputElement).Checked == false {
						document.GetElementByID("permr" + pitemname + pentityid).(*dom.HTMLInputElement).Checked = true
					}
				}
			case "n":
				if pentityid == "-1" {
					for _, e := range document.GetElementsByClassName("permnproducts") {
						e.(*dom.HTMLInputElement).Checked = true
					}
				} else {
					document.GetElementByID("permn" + pitemname + pentityid).(*dom.HTMLInputElement).Checked = true
				}
			}
		case "rproducts":
			switch ppermname {
			case "w", "all":
				if pentityid == "-1" {
					for _, e := range document.GetElementsByClassName("permwrproducts") {
						e.(*dom.HTMLInputElement).Checked = true
					}
				} else {
					document.GetElementByID("permw" + pitemname + pentityid).(*dom.HTMLInputElement).Checked = true
				}
			case "r":
				if pentityid == "-1" {
					for _, e := range document.GetElementsByClassName("permrrproducts") {
						// avoid selecting r if w is already selected
						eid := e.GetAttribute("entity_id")
						if !document.GetElementByID("permw" + pitemname + eid).(*dom.HTMLInputElement).Checked {
							e.(*dom.HTMLInputElement).Checked = true
						}
					}
				} else {
					// avoid selecting r if w is already selected
					if document.GetElementByID("permw"+pitemname+pentityid).(*dom.HTMLInputElement).Checked == false {
						document.GetElementByID("permr" + pitemname + pentityid).(*dom.HTMLInputElement).Checked = true
					}
				}
			case "n":
				if pentityid == "-1" {
					for _, e := range document.GetElementsByClassName("permnrproducts") {
						e.(*dom.HTMLInputElement).Checked = true
					}
				} else {
					document.GetElementByID("permn" + pitemname + pentityid).(*dom.HTMLInputElement).Checked = true
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
						// avoid selecting r if w is already selected
						eid := e.GetAttribute("entity_id")
						if !document.GetElementByID("permw" + pitemname + eid).(*dom.HTMLInputElement).Checked {
							e.(*dom.HTMLInputElement).Checked = true
						}
					}
				} else {
					// avoid selecting r if w is already selected
					if document.GetElementByID("permw"+pitemname+pentityid).(*dom.HTMLInputElement).Checked == false {
						document.GetElementByID("permr" + pitemname + pentityid).(*dom.HTMLInputElement).Checked = true
					}
				}
			}
		// case "storelocations":
		// 	switch ppermname {
		// 	case "w", "all":
		// 		if pentityid == "-1" {
		// 			for _, e := range document.GetElementsByClassName("permwstorelocations") {
		// 				e.(*dom.HTMLInputElement).Checked = true
		// 			}
		// 		} else {
		// 			document.GetElementByID("permw" + pitemname + pentityid).(*dom.HTMLInputElement).Checked = true
		// 		}
		// 	case "r":
		// 		if pentityid == "-1" {
		// 			for _, e := range document.GetElementsByClassName("permrstorelocations") {
		// 				// avoid selecting r if w is already selected
		// 				eid := e.GetAttribute("entity_id")
		// 				if !document.GetElementByID("permw" + pitemname + eid).(*dom.HTMLInputElement).Checked {
		// 					e.(*dom.HTMLInputElement).Checked = true
		// 				}
		// 			}
		// 		} else {
		// 			// avoid selecting r if w is already selected
		// 			if document.GetElementByID("permw"+pitemname+pentityid).(*dom.HTMLInputElement).Checked == false {
		// 				document.GetElementByID("permr" + pitemname + pentityid).(*dom.HTMLInputElement).Checked = true
		// 			}
		// 		}
		// 	}
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

// BuildInlineRadioElement return a radio inline block
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
func BuildPermissionWidget(entityID int, entityName string, ismanager bool) *dom.HTMLDivElement {
	var widgetdiv *dom.HTMLDivElement
	// create main widget div
	widgetdiv = document.CreateElement("div").(*dom.HTMLDivElement)
	widgetdiv.SetID(fmt.Sprintf("perm%d", entityID))
	widgetdiv.SetClass("col-sm-12")
	title := document.CreateElement("div").(*dom.HTMLDivElement)
	title.SetClass("d-flex")
	title.SetInnerHTML("<span class='mdi mdi-store mdi-24px'/>" + entityName)

	widgetdiv.AppendChild(title)

	if ismanager {
		s := document.CreateElement("span").(*dom.HTMLSpanElement)
		s.SetClass("mdi mdi-36px mdi-account-star")
		s.SetAttribute("title", "manager")
		coldiv := document.CreateElement("div").(*dom.HTMLDivElement)
		coldiv.SetClass("col-sm-2")
		coldiv.AppendChild(s)

		widgetdiv.AppendChild(coldiv)
		return widgetdiv
	}

	for _, i := range permitems {
		// products permissions widget is static
		if i != "products" && i != "rproducts" {
			//println(i)
			// building main row
			mainrowdiv := document.CreateElement("div").(*dom.HTMLDivElement)
			mainrowdiv.SetClass("form-group row d-flex")
			// building first col for table name
			label := document.CreateElement("div").(*dom.HTMLDivElement)
			label.SetClass("iconlabel text-right")
			label.SetInnerHTML(i)
			firstcoldiv := document.CreateElement("div").(*dom.HTMLDivElement)
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
				"item_name": fmt.Sprintf("%s", i),
				"entity_id": fmt.Sprintf("%d", entityID),
				"class":     fmt.Sprintf("perm permn permn%s", i)}
			readradioattrs := map[string]string{
				"id":        fmt.Sprintf("permr%s%d", i, entityID),
				"name":      fmt.Sprintf("perm%s%d", i, entityID),
				"value":     "r",
				"label":     fmt.Sprintf("<label class=\"form-check-label ml-sm-1 pr-sm-1 pl-sm-1 text-secondary border border-secondary rounded\" for=\"permn%s%d\"><span class=\"mdi mdi-eye\"></span></label>", i, entityID),
				"perm_name": "r",
				"item_name": fmt.Sprintf("%s", i),
				"entity_id": fmt.Sprintf("%d", entityID),
				"class":     fmt.Sprintf("perm permr permr%s", i)}
			writeradioattrs := map[string]string{
				"id":        fmt.Sprintf("permw%s%d", i, entityID),
				"name":      fmt.Sprintf("perm%s%d", i, entityID),
				"value":     "w",
				"label":     fmt.Sprintf("<label class=\"form-check-label ml-sm-1 pr-sm-1 pl-sm-1 text-secondary border border-secondary rounded\" for=\"permn%s%d\"><span class=\"mdi mdi-eye\"></span><span class=\"mdi mdi-creation\"></span><span class=\"mdi mdi-border-color\"></span><span class=\"mdi mdi-delete\"></span></label>", i, entityID),
				"perm_name": "w",
				"item_name": fmt.Sprintf("%s", i),
				"entity_id": fmt.Sprintf("%d", entityID),
				"class":     fmt.Sprintf("perm permw permw%s", i)}
			secondcoldiv := document.CreateElement("div").(*dom.HTMLDivElement)
			secondcoldiv.SetClass("col-sm-4")
			secondcoldiv.AppendChild(BuildInlineRadioElement(noneradioattrs))
			secondcoldiv.AppendChild(BuildInlineRadioElement(readradioattrs))
			secondcoldiv.AppendChild(BuildInlineRadioElement(writeradioattrs))

			// appending to final div
			mainrowdiv.AppendChild(firstcoldiv)
			mainrowdiv.AppendChild(secondcoldiv)
			widgetdiv.AppendChild(mainrowdiv)
		}
	}
	return widgetdiv
}

// DisplayMessage display fading messages at the
// top of the screen
func DisplayMessage(msgText string, msgType string) {
	d := document.CreateElement("div").(*dom.HTMLDivElement)
	s := document.CreateElement("span").(*dom.HTMLSpanElement)
	d.SetAttribute("role", "alert")
	d.SetAttribute("style", "z-index:2;")
	d.Class().SetString("animated fadeOutUp delay-2s fixed-top w-100 p-3 text-center alert alert-" + msgType)
	s.SetTextContent(msgText)
	d.AppendChild(s)

	document.GetElementByID("message").SetInnerHTML("")
	document.GetElementByID("message").AppendChild(d)
}

// NormalizeSqlNull removes from the obj map keys the .Foo (.String, .Int64...) prefixes
func NormalizeSqlNull(obj map[string]interface{}) *map[string]interface{} {
	// result map
	r := make(map[string]interface{})

	// sqlNull values type detection regexp
	regexS := regexp.MustCompile("(.+)\\.String")
	regexI := regexp.MustCompile("(.+)\\.Int64")
	regexF := regexp.MustCompile("(.+)\\.Float64")
	regexB := regexp.MustCompile("(.+)\\.Bool")

	for k, iv := range obj {
		// trying to match a regex
		mS := regexS.FindStringSubmatch(k)
		mI := regexI.FindStringSubmatch(k)
		mF := regexF.FindStringSubmatch(k)
		mB := regexB.FindStringSubmatch(k)

		// building the new map without
		// the .Foo in the key names
		if len(mS) > 0 {
			r[mS[1]] = iv.(string)
		} else if len(mI) > 0 {
			r[mI[1]] = iv.(float64)
		} else if len(mF) > 0 {
			r[mF[1]] = iv.(float64)
		} else if len(mB) > 0 {
			r[mB[1]] = iv.(bool)
		} else {
			r[k] = iv
		}
	}
	return &r
}

func T(s string, accept string) string {
	if s == "" {
		return ""
	}

	// accept is an accept language http header
	// like fr,fr-FR;q=0.8,en-US;q=0.5,en;q=0.3
	r := regexp.MustCompile("([a-zA-Z\\-]+){0,1},{0,1}([a-zA-Z\\-]+){0,1};q=([0-9\\.]+)")

	ms := r.FindAllStringSubmatch(accept, -1)
	if len(ms) == 0 {
		return "translation error"
	}

	// lazily assuming that the entries are
	// ordered by the preferred language
	for _, m := range ms {
		for _, i := range m {
			js_locale_varname := fmt.Sprintf("locale_%s_%s", i, s)
			translated := js.Global.Get(js_locale_varname)
			if translated != js.Undefined {
				return translated.String()
			}
		}
	}
	return "translation error"
}

// type Foo struct {
// 	*js.Object
// 	Toto                string `js:"toto"`
// 	StorageCreationDate string `js:"storage_creationdate"`
// }

// func (f *Foo) GetToto() string {
// 	return f.Toto
// }

// func Test(obj *js.Object) *js.Object {
// 	t := &Foo{Object: obj}
// 	t.Toto = "toto"
// 	fmt.Println(t.StorageCreationDate)
// 	fmt.Println(t.Toto)
// 	return js.MakeWrapper(t)
// }

func main() {

	// exporting functions to be called from other JS files
	js.Global.Set("global", map[string]interface{}{
		"buildPermissionWidget":    BuildPermissionWidget,
		"populatePermissionWidget": PopulatePermissionWidget,
		"createTitle":              CreateTitle,
		"normalizeSqlNull":         NormalizeSqlNull,
		"displayMessage":           DisplayMessage,
		"t":                        T,
	})

}

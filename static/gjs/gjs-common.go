package main

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"
	. "github.com/tbellembois/gochimitheque/static/gjs/globals"
	"github.com/tbellembois/gochimitheque/static/gjs/utils"
	"github.com/tbellembois/gochimitheque/static/gjs/views/entity"
	"github.com/tbellembois/gochimitheque/static/gjs/widgets"
)

func Test() {
	fmt.Println(URLLocationPathName)
	utils.DisplayGenericErrorMessage()
}

func main() {

	// exporting functions to be called from other JS files
	js.Global.Set("gjsGlobals", map[string]interface{}{
		"test": Test,
	})
	js.Global.Set("gjsWidgets", map[string]interface{}{
		"permission":         widgets.Permission,
		"permissionPopulate": widgets.PopulatePermission,
		"title":              widgets.Title,
	})
	js.Global.Set("gjsUtils", map[string]interface{}{
		"normalizeSqlNull": utils.NormalizeSqlNull,
		"translate":        utils.Translate,
		"message":          utils.DisplayMessage,
		"closeEdit":        utils.CloseEdit,
	})
	js.Global.Set("gjsEntity", map[string]interface{}{
		"getData":           entity.GetTableData,
		"managersFormatter": entity.ManagersFormatter,
		"operateFormatter":  entity.OperateFormatter,
		"saveEntity":        entity.SaveEntity,
	})

	switch URLLocationPathName {
	case "/vc/entities", "/v/entities":
		entity.InsertJS()
	}

}

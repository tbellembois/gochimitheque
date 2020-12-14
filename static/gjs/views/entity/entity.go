package entity

import (
	"github.com/gopherjs/gopherjs/js"
	. "github.com/tbellembois/gochimitheque/static/gjs/globals"
	. "github.com/tbellembois/gochimitheque/static/gjs/types"
	"github.com/tbellembois/gochimitheque/static/gjs/utils"
)

func InsertJS() {

	// validate
	Jq("#entity").Validate(ValidateConfig{
		ErrorClass: "alert alert-danger",
		Rules: map[string]ValidateRule{
			"entity_name": {
				Required: true,
				Remote: ValidateRemote{
					URL:        "",
					Type:       "post",
					BeforeSend: ValidateEntityNameBeforeSend,
				},
			},
		},
		Messages: map[string]ValidateMessage{
			"entity_name": {
				Required: utils.Translate("required_input", HTTPHeaderAcceptLanguage),
			},
		},
	})

	// select2
	Jq("select#managers").Select2(Select2Config{
		Ajax: Select2Ajax{
			URL:            ApplicationProxyPath + "people",
			DataType:       "json",
			Data:           Select2ManagersAjaxData,
			ProcessResults: Select2ManagersAjaxProcessResults,
		},
	})

	Jq("#table").Call("on", "load-success.bs.table", func() {
		showIfAuthorizedActionButtons()
	})

	// Bootstrap table operations.
	var operateEvents = map[string]func(e, value *js.Object, row *Entity, index *js.Object){
		"click .storelocations": OperateEventsStorelocations,
		"click .members":        OperateEventsMembers,
		"click .edit":           OperateEventsEdit,
		"click .delete":         OperateEventsDelete,
	}
	js.Global.Get("window").Set("operateEvents", operateEvents)

}

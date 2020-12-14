package entity

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/gopherjs/gopherjs/js"
	. "github.com/tbellembois/gochimitheque/static/gjs/globals"
	. "github.com/tbellembois/gochimitheque/static/gjs/types"
	"github.com/tbellembois/gochimitheque/static/gjs/utils"
	"github.com/tbellembois/gochimitheque/static/gjs/widgets"
	"github.com/tbellembois/gochimitheque/static/gjs/widgets/themes"
	"honnef.co/go/js/xhr"
)

func OperateEventsStorelocations(e, value *js.Object, row *Entity, index *js.Object) {

	href := fmt.Sprintf("%sv/storelocations?entity=%d", ApplicationProxyPath, row.EntityID)
	utils.RedirectTo(href)

}

func OperateEventsMembers(e, value *js.Object, row *Entity, index *js.Object) {

	href := fmt.Sprintf("%sv/people?entity=%d", ApplicationProxyPath, row.EntityID)
	utils.RedirectTo(href)

}

func OperateEventsDelete(e, value *js.Object, row *Entity, index *js.Object) {

	url := fmt.Sprintf("%sentities/%d", ApplicationProxyPath, row.EntityID)
	method := "delete"

	done := func(data *js.Object) {

		utils.DisplayMessage(utils.Translate("entity_deleted_message", HTTPHeaderAcceptLanguage), "success")
		Jq("#table").Bootstraptable().Refresh()

	}
	fail := func(jqXHR *js.Object) {

		utils.DisplayGenericErrorMessage()

	}

	Ajax{
		Method: method,
		URL:    url,
		Done:   done,
		Fail:   fail,
	}.Send()

}

func OperateEventsEdit(e, value *js.Object, row *Entity, index *js.Object) {

	url := fmt.Sprintf("%sentities/%d", ApplicationProxyPath, row.EntityID)
	method := "get"

	done := func(data *js.Object) {

		o := js.Global.Get("JSON").Call("parse", data)
		entity := &Entity{Object: o}

		FillInEntityForm(entity, "edit-collapse")

		Jq("input#index").SetVal(index.Int())
		Jq("#edit-collapse").Show()

	}
	fail := func(jqXHR *js.Object) {

		utils.DisplayGenericErrorMessage()

	}

	Ajax{
		Method: method,
		URL:    url,
		Done:   done,
		Fail:   fail,
	}.Send()

}

func OperateFormatter(value *js.Object, row *Entity, index *js.Object) string {

	buttonStorelocations := widgets.NewBSButtonWithIcon(
		widgets.ButtonAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:         "storelocations" + strconv.Itoa(row.EntityID),
				Classes:    []string{"storelocations"},
				Visible:    false,
				Attributes: map[string]string{"eid": strconv.Itoa(row.EntityID)},
			},
			Title: utils.Translate("storelocations", HTTPHeaderAcceptLanguage),
		},
		widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible:    true,
				Attributes: map[string]string{"eid": strconv.Itoa(row.EntityID)},
			},
			Text: strconv.Itoa(row.EntitySLC),
			Icon: themes.NewMdiIcon(themes.MDI_STORELOCATION, ""),
		},
		[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
	).OuterHTML()

	buttonMembers := widgets.NewBSButtonWithIcon(
		widgets.ButtonAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:         "members" + strconv.Itoa(row.EntityID),
				Classes:    []string{"members"},
				Visible:    false,
				Attributes: map[string]string{"eid": strconv.Itoa(row.EntityID)},
			},
			Title: utils.Translate("members", HTTPHeaderAcceptLanguage),
		},
		widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible:    true,
				Attributes: map[string]string{"eid": strconv.Itoa(row.EntityID)},
			},
			Text: strconv.Itoa(row.EntitySLC),
			Icon: themes.NewMdiIcon(themes.MDI_PEOPLE, ""),
		},
		[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
	).OuterHTML()

	buttonEdit := widgets.NewBSButtonWithIcon(
		widgets.ButtonAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:         "edit" + strconv.Itoa(row.EntityID),
				Classes:    []string{"edit"},
				Visible:    false,
				Attributes: map[string]string{"eid": strconv.Itoa(row.EntityID)},
			},
			Title: utils.Translate("edit", HTTPHeaderAcceptLanguage),
		},
		widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible:    true,
				Attributes: map[string]string{"eid": strconv.Itoa(row.EntityID)},
			},
			Text: strconv.Itoa(row.EntitySLC),
			Icon: themes.NewMdiIcon(themes.MDI_EDIT, ""),
		},
		[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
	).OuterHTML()

	buttonDelete := widgets.NewBSButtonWithIcon(
		widgets.ButtonAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Id:         "delete" + strconv.Itoa(row.EntityID),
				Classes:    []string{"delete"},
				Visible:    false,
				Attributes: map[string]string{"eid": strconv.Itoa(row.EntityID)},
			},
			Title: utils.Translate("delete", HTTPHeaderAcceptLanguage),
		},
		widgets.IconAttributes{
			BaseAttributes: widgets.BaseAttributes{
				Visible:    true,
				Attributes: map[string]string{"eid": strconv.Itoa(row.EntityID)},
			},
			Text: strconv.Itoa(row.EntitySLC),
			Icon: themes.NewMdiIcon(themes.MDI_DELETE, ""),
		},
		[]themes.BSClass{themes.BS_BTN, themes.BS_BNT_LINK},
	).OuterHTML()

	return buttonStorelocations + buttonMembers + buttonEdit + buttonDelete
}

// Can not set and use value []*person as argument because nil
// managers values lead to JS exception during Go->JS conversion.
// We use the full row (Entity) parameter instead.
func ManagersFormatter(value *js.Object, row *Entity, index, field *js.Object) string {

	ul := widgets.NewUl(widgets.UlAttributes{})

	// Checking row.Managers == nil leads to
	// JS exception when Managers is nil.
	// Checking the JS value instead.
	if row.Object.Get("managers").Bool() {
		for _, manager := range row.Managers {
			li := widgets.NewLi(widgets.LiAttributes{
				Text: manager.PersonEmail,
			})
			ul.AppendChild(li)
		}
	}

	return ul.OuterHTML()
}

// Get remote bootstrap table data - params object defined here:
// https://bootstrap-table.com/docs/api/table-options/#queryparams
func GetTableData(params *QueryParams) {

	go func() {

		u := url.URL{Path: ApplicationProxyPath + "entities"}
		values := url.Values{}
		values.Set("search", params.Data.Search)
		values.Set("page", strconv.Itoa(params.Data.Page))
		values.Set("offset", strconv.Itoa(params.Data.Offset))
		values.Set("limit", strconv.Itoa(params.Data.Limit))
		u.RawQuery = values.Encode()

		req := xhr.NewRequest("get", u.String())
		req.ResponseType = xhr.JSON

		err := req.Send(nil)
		if err != nil {
			println(err)
		}

		if req.Status == 200 {
			params.Call("success", Response{
				Rows:  req.Response.Get("rows"),
				Total: req.Response.Get("total"),
			})
		}
	}()

}

func showIfAuthorizedActionButtons() {

	// Iterating other the button with the class "storelocation"
	// (we could choose "members" or "delete")
	// to retrieve once the entity id.
	for _, button := range Doc.GetElementsByTagName("button") {
		if button.Class().Contains("storelocations") {
			entityId := button.GetAttribute("eid")

			utils.HasPermission("storelocations", entityId, "get", func() {
				Jq("#storelocations" + entityId).FadeIn()
			}, func() {
			})
			utils.HasPermission("people", entityId, "get", func() {
				Jq("#members" + entityId).FadeIn()
			}, func() {
			})
			utils.HasPermission("entities", entityId, "put", func() {
				Jq("#edit" + entityId).FadeIn()
			}, func() {
			})
			utils.HasPermission("entities", entityId, "delete", func() {
				Jq("#delete" + entityId).FadeIn()
			}, func() {
			})
		}
	}

}

package entity

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gopherjs/gopherjs/js"
	. "github.com/tbellembois/gochimitheque/static/gjs/globals"
	"github.com/tbellembois/gochimitheque/static/gjs/localStorage"
	"github.com/tbellembois/gochimitheque/static/gjs/types"
	. "github.com/tbellembois/gochimitheque/static/gjs/types"
	"github.com/tbellembois/gochimitheque/static/gjs/utils"
	"github.com/tbellembois/gochimitheque/static/gjs/widgets"
)

func FillInEntityForm(e *Entity, id string) {

	Jq(fmt.Sprintf("#%s #entity_id", id)).SetVal(e.EntityID)
	Jq(fmt.Sprintf("#%s #entity_name", id)).SetVal(e.EntityName)
	Jq(fmt.Sprintf("#%s #entity_description", id)).SetVal(e.EntityDescription)

	for _, manager := range e.Managers {
		Jq("select#managers").Select2AppendOption(
			widgets.NewOption(widgets.OptionAttributes{
				Text:            manager.PersonEmail,
				Value:           strconv.Itoa(manager.PersonId),
				DefaultSelected: true,
				Selected:        true,
			}),
		)
	}

}

func SaveEntity() {

	var (
		ajaxURL, ajaxMethod string
		entity              *Entity
		dataBytes           []byte
		err                 error
	)

	if !Jq("#entity").Valid() {
		return
	}

	entity = types.NewEntity()
	entity.EntityID = Jq("input#entity_id").GetVal().Int()
	entity.EntityName = Jq("input#entity_name").GetVal().String()
	entity.EntityDescription = Jq("input#entity_description").GetVal().String()

	for _, select2Item := range Jq("select#managers").Select2Data() {
		person := NewPerson()
		if person.PersonId, err = strconv.Atoi(select2Item.Id); err != nil {
			println(err)
		}
		person.PersonEmail = select2Item.Text

		entity.Managers = append(entity.Managers, person)
	}

	if dataBytes, err = json.Marshal(entity); err != nil {
		println(err)
	}

	if Jq("form#entity input#entity_id").Length() > 0 {
		ajaxURL = fmt.Sprintf("%sentities/%d", ApplicationProxyPath, entity.EntityID)
		ajaxMethod = "put"
	} else {
		ajaxURL = fmt.Sprintf("%sentities", ApplicationProxyPath)
		ajaxMethod = "post"
	}

	Ajax{
		URL:    ajaxURL,
		Method: ajaxMethod,
		Data:   dataBytes,
		Done: func(data *js.Object) {

			localStorage.Clear()

			o := js.Global.Get("JSON").Call("parse", data)
			entity := &Entity{Object: o}

			// TODO: use entityId for redirection
			href := fmt.Sprintf("%sv/entities?search=%s", ApplicationProxyPath, entity.EntityName)
			utils.RedirectTo(href)

		},
		Fail: func(jqXHR *js.Object) {

			utils.DisplayGenericErrorMessage()

		},
	}.Send()

}

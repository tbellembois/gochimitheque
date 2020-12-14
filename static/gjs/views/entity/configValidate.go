package entity

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"
	. "github.com/tbellembois/gochimitheque/static/gjs/globals"
)

func ValidateEntityNameBeforeSend(_ *js.Object, settings *js.Object) {

	id := -1
	eid := Jq("input#entity_id")

	if eid.Length() > 0 {
		id = eid.GetVal().Int()
	}

	settings.Set("url", fmt.Sprintf("%svalidate/entity/%d/name/", ApplicationProxyPath, id))

}

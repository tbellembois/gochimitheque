package entity

import (
	"strconv"

	"github.com/gopherjs/gopherjs/js"
	. "github.com/tbellembois/gochimitheque/static/gjs/types"
)

func Select2ManagersAjaxProcessResults(data *People, params *js.Object) *Select2Data {

	page := params.Get("page").Int()
	total := data.Total

	var select2Items []*Select2Item

	for _, person := range data.Rows {
		select2Item := NewSelect2Item(strconv.Itoa(person.PersonId), person.PersonEmail)
		select2Items = append(select2Items, select2Item)
	}

	select2Pagination := NewSelect2Pagination((page * 10) < total)

	return NewSelect2Data(select2Items, select2Pagination)

}

func Select2ManagersAjaxData(params *js.Object) *QueryFilter {

	search := params.Get("term").String()
	page := params.Get("page").Int()
	offset := (page - 1) * 10
	limit := 10

	if offset < 0 {
		offset = 0
	}
	if search == js.Undefined.String() {
		search = ""
	}

	return NewQueryFilter(search, offset, page, limit)

}

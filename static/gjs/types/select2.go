package types

import (
	"github.com/gopherjs/gopherjs/js"
)

// Select2Config is a select2 parameters struct
// as defined https://select2.org/configuration
type Select2Config struct {
	Placeholder string      `json:"placeholder"`
	Tags        bool        `json:"tags"`
	Ajax        Select2Ajax `json:"ajax"`
}

// Select2Ajax is a select2 ajax request
// as defined https://select2.org/data-sources/ajax
type Select2Ajax struct {
	URL            string      `json:"url"`
	DataType       string      `json:"datatype"`
	Data           interface{} `json:"data"`
	ProcessResults interface{} `json:"processResults"`
}

// Select2Data is a select2 data struct
// as defined https://select2.org/data-sources/formats
type Select2Data struct {
	*js.Object
	Results    []*Select2Item     `js:"results" json:"results"`
	Pagination *Select2Pagination `js:"pagination" json:"pagination"`
}

type Select2Pagination struct {
	*js.Object
	More bool `js:"more" json:"more"`
}

type Select2Item struct {
	*js.Object
	Id   string `js:"id" json:"id"`
	Text string `js:"text" json:"text"`
}

func NewSelect2Data(results []*Select2Item, pagination *Select2Pagination) *Select2Data {

	select2Data := &Select2Data{Object: js.Global.Get("Object").New()}
	select2Data.Results = results
	select2Data.Pagination = pagination
	return select2Data

}

func NewSelect2Pagination(more bool) *Select2Pagination {

	select2Pagination := &Select2Pagination{Object: js.Global.Get("Object").New()}
	select2Pagination.More = more
	return select2Pagination

}

func NewSelect2Item(id string, text string) *Select2Item {

	select2Item := &Select2Item{Object: js.Global.Get("Object").New()}
	select2Item.Id = id
	select2Item.Text = text
	return select2Item

}

func (jq Jquery) Select2(config Select2Config) {

	configMap := StructToMap(config)
	jq.Call("select2", configMap)

}

func (jq Jquery) Select2Data() []*Select2Item {

	var (
		select2Items []*Select2Item
	)

	select2Data := jq.Call("select2", "data").Interface().([]interface{})

	for _, select2DataItem := range select2Data {
		select2Item := select2DataItem.(map[string]interface{})

		select2Items = append(select2Items, NewSelect2Item(
			select2Item["id"].(string),
			select2Item["text"].(string),
		))
	}

	return select2Items

}

func (jq Jquery) Select2AppendOption(o interface{}) {

	jq.Call("append", o).Call("trigger", "change")

}

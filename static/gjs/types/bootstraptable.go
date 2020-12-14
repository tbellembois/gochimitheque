package types

import (
	"github.com/gopherjs/gopherjs/js"
)

type Bootstraptable struct {
	Jquery
}

// QueryParams is the data sent while requesting
// remote data as defined
// https://bootstrap-table.com/docs/api/table-options/#queryparams
type QueryParams struct {
	*js.Object
	Data QueryFilter `js:"data"`
}

func (jq Jquery) Bootstraptable() *Bootstraptable {

	return &Bootstraptable{Jquery: jq}

}

func (bt Bootstraptable) Refresh() {

	bt.Call("bootstrapTable", "refresh")

}

func (bt Bootstraptable) ResetSearch(search string) {

	bt.Call("bootstrapTable", "resetSearch", search)

}

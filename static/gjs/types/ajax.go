package types

import (
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/xhr"
)

// QueryFilter contains the parameters sent
// while doing AJAX requests to retrieve multiple
// results ("/entities", "/people"...).
// It is especially used by select2 and bootstraptable.
type QueryFilter struct {
	*js.Object
	Search string `js:"search" json:"search"`
	Page   int    `js:"page" json:"page"`
	Offset int    `js:"offset" json:"offset"`
	Limit  int    `js:"limit" json:"limit"`
}

func NewQueryFilter(search string, offset, page, limit int) *QueryFilter {

	query := &QueryFilter{Object: js.Global.Get("Object").New()}
	query.Search = search
	query.Page = page
	query.Offset = offset
	query.Limit = limit
	return query

}

// Response contains the data retrieved
// from the query above.
type Response struct {
	Rows  *js.Object `js:"rows"`
	Total *js.Object `js:"total"`
}

type Ajax struct {
	URL    string
	Method string
	Data   []byte
	Done   AjaxDone
	Fail   AjaxFail
}

type AjaxDone func(data *js.Object)
type AjaxFail func(jqXHR *js.Object)

func (ajax Ajax) Send() {

	go func() {

		var (
			err error
		)

		req := xhr.NewRequest(ajax.Method, ajax.URL)
		req.SetRequestHeader("Content-Type", "application/json; charset=utf-8")

		if err = req.Send(ajax.Data); err != nil {
			println(err)
			ajax.Fail(req.Response)
		}

		if req.Status == 200 {
			ajax.Done(req.Response)
		} else {
			ajax.Fail(req.Response)
		}

	}()

}

package utils

import (
	"fmt"

	. "github.com/tbellembois/gochimitheque/static/gjs/globals"
	localStorage "github.com/tbellembois/gochimitheque/static/gjs/localStorage"
	"honnef.co/go/js/xhr"
)

func HasPermission(item, id, method string, done, fail func()) {

	go func() {
		cacheKey := fmt.Sprintf("%s:%s:%s", item, id, method)
		cachedPermission := localStorage.GetItem(cacheKey)

		if !DisableCache && cachedPermission != "" {
			if cachedPermission == "true" {
				done()
			} else {
				fail()
			}
		} else {
			var url string
			if id != "" {
				url = ApplicationProxyPath + "f/" + item + "/" + id
			} else {
				url = ApplicationProxyPath + "f/" + item
			}

			req := xhr.NewRequest(method, url)
			err := req.Send(nil)
			if err != nil {
				println(err)
				fail()
			}
			if req.Status == 200 {
				localStorage.SetItem(cacheKey, "true")
				done()
			} else {
				localStorage.SetItem(cacheKey, "false")
				fail()
			}
		}
	}()

}

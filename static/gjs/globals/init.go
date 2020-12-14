package globals

import (
	"strconv"

	"net/url"

	"github.com/gopherjs/gopherjs/js"
	"github.com/tbellembois/gochimitheque/static/gjs/types"
	"honnef.co/go/js/dom"
)

var (
	Win dom.Window
	Doc dom.HTMLDocument
	Jq  func(args ...interface{}) *types.Jquery

	fullUrl                                                             *url.URL
	URLParameters                                                       url.Values
	URLLocationPathName, ApplicationProxyPath, HTTPHeaderAcceptLanguage string
	DisableCache                                                        bool

	// permissions
	PermItems = [3]string{
		"rproducts",
		"products",
		"storages"}

	err error
)

func init() {

	Win = dom.GetWindow()
	Doc = Win.Document().(dom.HTMLDocument)
	Jq = types.NewJquery
	URLLocationPathName = js.Global.Get("location").Get("pathname").String()
	fullUrl, err = url.Parse(js.Global.Get("location").Get("href").String())
	if err != nil {
		println(err)
	}
	URLParameters, err = url.ParseQuery(fullUrl.RawQuery)
	if err != nil {
		println(err)
	}

	// TODO: get the variables from Go instead of JS
	ApplicationProxyPath = js.Global.Get("proxyPath").String()
	HTTPHeaderAcceptLanguage = js.Global.Get("container.PersonLanguage").String()
	DisableCache, err = strconv.ParseBool(js.Global.Get("disableCache").String())
	if err != nil {
		println(err)
	}

	// TODO: move this
	// Common actions for all pages.
	Jq("#table").Bind("load-success.bs.table", func() {
		search := URLParameters.Get("search")
		if search != "" {
			Jq("#table").Bootstraptable().ResetSearch(search)
		}
	})

}

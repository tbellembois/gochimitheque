package utils

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"

	"golang.org/x/text/language"
)

// Translate translates the s string in the "accept" language
func Translate(s string, accept string) string {
	if s == "" {
		return ""
	}

	ts, _, e := language.ParseAcceptLanguage(accept)
	if e != nil {
		// falling back on english if error
		ts = []language.Tag{language.English}
	}

	// the t entries are
	// ordered by the preferred language
	for _, t := range ts {
		js_locale_varname := fmt.Sprintf("locale_%s_%s", t, s)
		translated := js.Global.Get(js_locale_varname)
		if translated != js.Undefined {
			return translated.String()
		}
	}
	return "translation error"
}

// package localStorage wraps the javascritp localStorage api for GOPHERJS.
package localStorage

import (
	"github.com/gopherjs/gopherjs/js"
)

var (
	localStorage = js.Global.Get("localStorage")
)

// Save val into localStorage under key
func SetItem(key string, val string) {
	localStorage.Call("setItem", key, val)
}

// Remove val into localStorage under key
func RemoveItem(key string) {
	localStorage.Call("removeItem", key)
}

// Return "" when no key found
func GetItem(key string) string {
	obj := localStorage.Call("getItem", key)
	if obj == nil {
		return ""
	}
	return obj.String()
}

// Clear the cache
func Clear() {
	localStorage.Call("clear")
}

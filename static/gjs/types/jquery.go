package types

import "github.com/gopherjs/gopherjs/js"

type Jquery struct {
	*js.Object
}

func NewJquery(args ...interface{}) *Jquery {

	return &Jquery{Object: js.Global.Get("jQuery").New(args...)}

}

// Selector does the same thing as the NewJquery method.
// This is a convenient way to avoid importing the types package
// where this would generate import cycles.
func (jq Jquery) Selector(i interface{}) Jquery {

	jq.Object = js.Global.Get("jQuery").New(i)
	return jq

}

func (jq Jquery) SetVal(i interface{}) Jquery {

	jq.Object.Call("val", i)
	return jq

}

func (jq Jquery) GetVal() *js.Object {

	return jq.Object.Call("val")

}

func (jq Jquery) Show() Jquery {

	jq.Object = jq.Object.Call("collapse", "show")
	return jq

}

func (jq Jquery) Hide() Jquery {

	jq.Object = jq.Object.Call("collapse", "hide")
	return jq

}

func (jq Jquery) FadeIn() Jquery {

	jq.Object = jq.Object.Call("fadeIn")
	return jq

}

func (jq Jquery) Bind(event string, i interface{}) Jquery {

	jq.Object = jq.Object.Call("bind", event, i)
	return jq

}

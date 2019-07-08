package components

import (
	h "github.com/azukaar/drago"
)

func App(document *h.DocumentObject) func() h.Node {

	return func() h.Node {
		return h.Div(h.Props{},
			h.H1(h.Props{}, h.Text("Hello World!")),
			h.Div(h.Props{},
				h.Text("This is some fancy shit!"),
				h.C(ToggleButton),
			),
		)
	}
}

package components

import (
	h "github.com/azukaar/drago"
)

func ToggleButton(document *h.DocumentObject) func() h.Node {
	toggled, setToggled := document.UseToggle(true)

	return func() h.Node {
		return h.Div(h.Props{},
			h.Text("Toggle me"),
			h.Button(h.Props{"onclick": setToggled},
				h.If(toggled(),
					h.Text("on"),
					h.Text("off"),
				),
			),
		)
	}
}

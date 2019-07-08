package scaffolding

import (
	"fmt"
	"net/url"

	"./components"

	h "."
	"github.com/zserge/webview"
)

type EventItem struct {
	key       string
	eventName string
	event     interface{}
}

var eventsQueue = make(chan EventItem)

type GoAPI struct {
}

func (c *GoAPI) EventDispatch(key string, eventName string, event interface{}) {
	eventsQueue <- EventItem{
		key:       key,
		eventName: eventName,
		event:     event,
	}
}

var window webview.WebView

func main() {
	const myHTML = `<!doctype html><html>
	<head>
		<script>
		</script>
	</head>
	<body>
	</body>
	</html>`

	window = webview.New(webview.Settings{
		URL:       `data:text/html,` + url.PathEscape(myHTML),
		Debug:     true,
		Title:     "Test",
		Width:     1024,
		Height:    768,
		Resizable: true,
	})

	if window == nil {
		fmt.Println("oh oh")
	}

	window.Bind("GoAPI", &GoAPI{})

	window.Loop(true)
	window.Loop(true)
	window.Loop(true)

	go func() {
		document := h.Document()

		document.OnChange(func() {
			window.Dispatch(func() {
				window.Eval("document.body.innerHTML = `" + document.Render() + "`")
			})
		})

		document.SetRoot(h.Div(h.Props{}, h.C(components.App)))

		go func() {
			for true {
				ev := <-eventsQueue
				document.NewEvent(ev.key, ev.eventName, ev.event)
			}
		}()
	}()

	window.Run()
}

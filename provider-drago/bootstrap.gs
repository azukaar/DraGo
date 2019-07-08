if(fileExists("gupm.json")) {
    console.error("A project already exist in this project. Aborting.")
    exit()
}

var name = waitForInput("Please enter the name of the project: ")
var description = waitForInput("Enter a description: ")
var author = waitForInput("Enter the author: ")
var licence = waitForInput("Enter the licence (ISC): ")

if(name == "") {
    console.error("Name cannot be empty. Try again.")
    exit()
}

var webkitFile = "package main\r\n\r\nimport (\r\n\t\"fmt\"\r\n\t\"net\/url\"\r\n\r\n\t\".\/components\"\r\n\r\n\th \"drago\"\r\n\t\"github.com\/zserge\/webview\"\r\n)\r\n\r\ntype EventItem struct {\r\n\tkey       string\r\n\teventName string\r\n\tevent     interface{}\r\n}\r\n\r\nvar eventsQueue = make(chan EventItem)\r\n\r\ntype GoAPI struct {\r\n}\r\n\r\nfunc (c *GoAPI) EventDispatch(key string, eventName string, event interface{}) {\r\n\teventsQueue <- EventItem{\r\n\t\tkey:       key,\r\n\t\teventName: eventName,\r\n\t\tevent:     event,\r\n\t}\r\n}\r\n\r\nvar window webview.WebView\r\n\r\nfunc main() {\r\n\tconst myHTML = `<!doctype html><html>\r\n\t<head>\r\n\t\t<script>\r\n\t\t<\/script>\r\n\t<\/head>\r\n\t<body>\r\n\t<\/body>\r\n\t<\/html>`\r\n\r\n\twindow = webview.New(webview.Settings{\r\n\t\tURL:       `data:text\/html,` + url.PathEscape(myHTML),\r\n\t\tDebug:     true,\r\n\t\tTitle:     \"Test\",\r\n\t\tWidth:     1024,\r\n\t\tHeight:    768,\r\n\t\tResizable: true,\r\n\t})\r\n\r\n\tif window == nil {\r\n\t\tfmt.Println(\"oh oh\")\r\n\t}\r\n\r\n\twindow.Bind(\"GoAPI\", &GoAPI{})\r\n\r\n\twindow.Loop(true)\r\n\twindow.Loop(true)\r\n\twindow.Loop(true)\r\n\r\n\tgo func() {\r\n\t\tdocument := h.Document()\r\n\r\n\t\tdocument.OnChange(func() {\r\n\t\t\twindow.Dispatch(func() {\r\n\t\t\t\twindow.Eval(\"document.body.innerHTML = `\" + document.Render() + \"`\")\r\n\t\t\t})\r\n\t\t})\r\n\r\n\t\tdocument.SetRoot(h.Div(h.Props{}, h.C(components.App)))\r\n\r\n\t\tgo func() {\r\n\t\t\tfor true {\r\n\t\t\t\tev := <-eventsQueue\r\n\t\t\t\tdocument.NewEvent(ev.key, ev.eventName, ev.event)\r\n\t\t\t}\r\n\t\t}()\r\n\t}()\r\n\r\n\twindow.Run()\r\n}\r\n";

var appFile = "package components\r\n\r\nimport (\r\n\th \"drago\"\r\n)\r\n\r\nfunc App(document *h.DocumentObject) func() h.Node {\r\n\r\n\treturn func() h.Node {\r\n\t\treturn h.Div(h.Props{},\r\n\t\t\th.H1(h.Props{}, h.Text(\"Hello World!\")),\r\n\t\t\th.Div(h.Props{},\r\n\t\t\t\th.Text(\"From Drago\"),\r\n\t\t\t),\r\n\t\t)\r\n\t}\r\n}\r\n";

var result = {
    name: name,
    description: description,
    author: author,
    licence: licence || "ISC",
    dependencies: {
        defaultProvider: "go",
        default: {
            "https://azukaar.github.io/DraGo/repo:drago": "master"
        }
    },
    cli: {
        aliases: {
            "webkit": "build/webkit",
            "browser": "build/browser"
        }
    }
}


writeJsonFile("gupm.json", result)
writeFile(".gupm_rc.gs", 'env("GOPATH", run("go", ["env", "GOROOT"]) + ":" + pwd() + "/go_modules")')

writeFile("build.gs", 
'removeFiles(["build"])\n' +
'exec("go", ["build", "-o", "build/webkit", "src/webkit.go"]);\n' + 
'exec("go", ["build", "-o", "build/browser", "src/browser.go"]);\n' 
)

mkdir('src')
mkdir('src/components')

writeFile("src/webkit.go", webkitFile)
writeFile("src/browser.go", webkitFile)

writeFile("src/components/app.go", webkitFile)

readme = "# "+name + "\n";
readme += "# Installation\n";
readme += "You need [GuPM](https://github.com/azukaar/GuPM) with the [provider-go](https://github.com/azukaar/GuPM-official#provider-go) plugin to run this project.\n";
readme += "```\n";
readme += "g make\n";
readme += "```\n";
readme += "# Add dependencies from a go project\n";
readme += "```\n";
readme += "g i newPackage\n";
readme += "```\n\n";
readme += "# Build and start\n";
readme += "```\n";
readme += "g build\n";
readme += "g start\n";
readme += "```\n\n";

writeFile("readme.md", readme)

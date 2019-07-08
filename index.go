package drago

import (
	"fmt"
	"strconv"
	"strings"
)

type Node struct {
	nodeType int
	index    string
	document *DocumentObject
	content  interface{}

	fun      func(*DocumentObject) func() Node
	renderer func() Node
	rendered *Node
}

type HTMLElement struct {
	tag      string
	props    Props
	children []Node
}

type Props map[string]interface{}

func New(tag string, props Props, children ...Node) Node {
	return Node{
		nodeType: ELEMENT_NODE,
		content: HTMLElement{
			tag:      tag,
			props:    props,
			children: children,
		},
	}
}

func Text(text string) Node {
	return Node{
		nodeType: ELEMENT_TEXT,
		content:  text,
	}
}

func (node *Node) Render(document *DocumentObject, key string) string {
	node.index = key
	node.document = document

	if node.nodeType == ELEMENT_NODE {
		htmlELement := node.content.(HTMLElement)
		res := htmlELement.RenderHTML(document, key)
		mergeContent(node, htmlELement, node.nodeType)
		return res
	} else if node.nodeType == ELEMENT_TEXT {
		return node.content.(string) + "\n"
	} else if node.nodeType == ELEMENT_FUNC {
		if node.renderer == nil {
			node.renderer = node.fun(node.document)
		}

		n := node.renderer()
		n.Render(document, key)
		if node.rendered != nil {
			mergeNode(node.rendered, n)
		} else {
			node.rendered = &n
		}

		res := node.rendered.Render(document, key)
		return res
	}

	fmt.Println("Error rendering", node)
	return ""
}

func mergeContent(node *Node, content interface{}, contentType int) {
	if node.content == nil || content == nil {
		node.content = content
	} else if contentType == ELEMENT_FUNC || contentType == ELEMENT_NODE {
		on := node.content.(HTMLElement)
		nn := content.(HTMLElement)

		if on.tag != nn.tag {
			node.content = content
		} else {
			on.props = nn.props

			for i, v := range on.children {
				mergeNode(&v, nn.children[i])
				on.children[i] = v
			}

			node.content = on
		}

	} else if contentType == ELEMENT_TEXT {
		node.content = content
	}
}

func mergeNode(node *Node, newNode Node) {
	if node.nodeType != newNode.nodeType {
		node.nodeType = newNode.nodeType
		node.index = newNode.index
		node.document = newNode.document
		node.content = newNode.content

		node.rendered = newNode.rendered
		node.fun = newNode.fun
		node.renderer = newNode.renderer
	} else if node.nodeType == ELEMENT_FUNC {
		// p1 := fmt.Sprintf("%v", node.fun)
		// p2 := fmt.Sprintf("%v", newNode.fun)

		// if p1 != p2 {
		// node.rendered = newNode.rendered
		// node.fun = newNode.fun
		// node.renderer = newNode.renderer
		// node.rendered = newNode.rendered
		// } else {
		// }
	} else if node.nodeType == ELEMENT_NODE || node.nodeType == ELEMENT_TEXT {
		mergeContent(node, newNode.content, node.nodeType)
	}
}

func (element *HTMLElement) RenderHTML(document *DocumentObject, key string) string {
	res := "<" + element.tag

	res += ` drago-index="` + key + `"`

	for index, value := range element.props {
		if valueString, ok := value.(string); ok {
			res += " " + index + `="` + valueString + `"`
		} else if valueFunc, ok := value.(func()); ok {
			_ = valueFunc
			res += " " + index + `="GoAPI.eventDispatch('` + key + `', '` + index + `', event)"`
		}
	}

	res += ">\n"

	for index, _ := range element.children {
		res += element.children[index].Render(document, key+"-"+strconv.Itoa(index))
	}

	res += "</" + element.tag + ">\n"

	return res
}

func (n *Node) GetByIndex(path []string) Node {
	if len(path) == 0 {
		return *n
	}
	if n.nodeType == ELEMENT_NODE {
		root := (n.content).(HTMLElement)
		for _, node := range root.children {
			nodePath := strings.Split(node.index, "-")[1:]
			if nodePath[len(nodePath)-1] == path[0] {
				return node.GetByIndex(path[1:])
			}
		}
	} else if n.nodeType == ELEMENT_FUNC {
		return n.rendered.GetByIndex(path)
	}

	return Nul()
}

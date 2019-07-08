package drago

import (
	"strings"
	"sync"
)

type DocumentObject struct {
	root     Node
	onChange func()
	lock     sync.Mutex
}

func Document() DocumentObject {
	return DocumentObject{}
}

func (a *DocumentObject) SetRoot(root Node) {
	a.root = root
	a.onChange()
}

func (a *DocumentObject) OnChange(cb func()) {
	a.onChange = func() {
		a.lock.Lock()
		cb()
		a.lock.Unlock()
	}
}

func (a *DocumentObject) Render() string {
	res := a.root.Render(a, "")
	return res
}

func (a *DocumentObject) UseState(defaultValue interface{}) (func() interface{}, func(value interface{})) {
	currentValue := defaultValue

	return func() interface{} {
			return currentValue
		}, func(value interface{}) {
			currentValue = value
			a.onChange()
		}
}

func (a *DocumentObject) UseToggle(defaultValue bool) (func() bool, func()) {
	value, fun := a.UseState(defaultValue)

	return func() bool {
			return value().(bool)
		}, func() {
			fun(!value().(bool))
		}
}

func (a *DocumentObject) NewEvent(key string, name string, event interface{}) {
	target := a.GetByIndex(strings.Split(key, "-")[1:])
	if target.nodeType == ELEMENT_NODE {
		targetContent := target.content.(HTMLElement)
		targetContent.props[name].(func())()
	}
}

func (a *DocumentObject) GetByIndex(path []string) Node {
	if len(path) == 0 {
		return a.root
	}

	root := (a.root.content).(HTMLElement)
	for _, node := range root.children {
		nodePath := strings.Split(node.index, "-")[1:]
		if nodePath[len(nodePath)-1] == path[0] {
			return node.GetByIndex(path[1:])
		}
	}

	return Nul()
}

package drago

import (
	"reflect"
	"time"
)

var ELEMENT_NUL = 0
var ELEMENT_NODE = 1
var ELEMENT_FUNC = 2
var ELEMENT_TEXT = 3

type Component func() Node

func (a *Node) Compare(b Node) bool {
	if a.nodeType != b.nodeType {
		return false
	}

	if a.nodeType == ELEMENT_NODE {
		aH := a.content.(HTMLElement)
		bH := b.content.(HTMLElement)

		return aH.Compare(bH)
	} else if a.nodeType == ELEMENT_TEXT {
		return a.content.(string) == b.content.(string)
	}

	return false
}

func (a *HTMLElement) Compare(b HTMLElement) bool {
	return a.tag == b.tag &&
		reflect.DeepEqual(a.props, b.props)
}

func SetTimeout(f func(), timeSet int) {
	go func() {
		time.Sleep(time.Duration(timeSet) * time.Millisecond)
		f()
	}()
}

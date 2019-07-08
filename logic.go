package drago

func If(cond bool, a Node, b Node) Node {
	if cond {
		return a
	} else {
		return b
	}
}

func Nul() Node {
	return Node{
		nodeType: ELEMENT_NUL,
	}
}

func C(fun func(*DocumentObject) func() Node) Node {
	return Node{
		nodeType: ELEMENT_FUNC,
		fun:      fun,
	}
}

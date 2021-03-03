package draw

import (
	"fmt"
	"strings"
	"syscall/js"
)

//go:generate go run github.com/abice/go-enum -f=$GOFILE --noprefix

/*
ENUM(
strokeStyle,
fillStyle,
lineWidth,
font,
textBaseline,
globalAlpha
styleEnumSize
)
*/
type attribute uint8

type styleAttribute struct {
	attr   attribute
	value  js.Value
	parent *styleAttribute
}

func (attr *styleAttribute) String() string {
	var b strings.Builder
	b.WriteString(attr.attr.String())
	b.WriteString(" => ")

	for tree := attr; tree != nil; tree = tree.parent {
		fmt.Fprint(&b, tree.value)
		b.WriteString(" -> ")
	}

	return b.String()
}

func (attr *styleAttribute) apply(ctx js.Value) {
	//if attr.parent != nil && attr.value.Equal(attr.parent.value) {
	//	// Do not apply if this value is the same as its parent
	//	return
	//}
	ctx.Set(attr.attr.String(), attr.value)
}

type styleStateMachine [StyleEnumSize]*styleAttribute

func (state *styleStateMachine) push(attr attribute, value interface{}) *styleAttribute {
	parent := state[attr]
	next := &styleAttribute{
		attr:   attr,
		value:  js.ValueOf(value),
		parent: parent,
	}
	state[attr] = next
	return next
}

func (state *styleStateMachine) pop(attr attribute) *styleAttribute {
	current := state[attr]
	if current == nil {
		return nil
	}

	state[attr] = current.parent
	return current.parent
}

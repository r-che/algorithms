package rbtree

import "fmt"

type KeyType int
const FakeNode = KeyType(-1)
func (k KeyType) String() string {
	if k == FakeNode {
		return "<>"
	}
	return fmt.Sprintf("%d", k)
}

type ColorType bool
const (
	Red		=	ColorType(true)
	Black	=	ColorType(false)
)


type RBNode struct {
	key		KeyType
	left	*RBNode
	right	*RBNode
	parent	*RBNode

	color	ColorType

	data	any
}

func NewRBNode(key KeyType, data any) *RBNode {
	return &RBNode{
		key: key,
		data: data,
	}
}

func (n *RBNode) String() string {
	if n == nil {
		return fmt.Sprintf("%s<nil>", Black)
	}
	if n.key == FakeNode {
		return fmt.Sprintf("%s<FAKE>", Black)
	}
	return fmt.Sprintf("%s%d", n.color, n.key)
}
func (n *RBNode) Color() ColorType {
	if n == nil {
		// Leaf always black
		return Black
	}

	return n.color
}

func (n *RBNode) SetColor(color ColorType) {
	if n != nil {
		n.color = color
	}
}

func (n *RBNode) Flip() {
	if n != nil {
		n.color = !n.color
	}
}

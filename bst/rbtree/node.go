package rbtree

import "fmt"

type KeyType int
const (
	FakeNode		=	KeyType(-1)
	strFakeNode		=	`<>`
	strInvalidNode	=	`<invalid>`
)

func (k KeyType) String() string {
	switch {
	case k == FakeNode:
		return strFakeNode
	case k < 0:
		return strInvalidNode
	default:
		return fmt.Sprintf("%d", k)
	}
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
		return Black.String() + "<nil>"
	}
	if n.key == FakeNode {
		return strFakeNode
	}

	return n.color.String() + n.key.String()
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

// Key returns the key value of the node
func (n *RBNode) Key() KeyType {
	if n == nil {
		return FakeNode
	}
	return n.key
}

// Value returns the data associated with the node
func (n *RBNode) Value() any {
	if n == nil {
		return nil
	}

	return n.data
}

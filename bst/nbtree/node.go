package nbtree

import "fmt"


//
// Key type
//
type KeyType int

const (
	FakeNode = KeyType(-1)

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

//
// BSTNode type
//
type BSTNode struct {
	key		KeyType
	left	*BSTNode
	right	*BSTNode
	parent	*BSTNode

	data any
}
func NewBSTNode(k KeyType, data any) *BSTNode {
	return &BSTNode{key: k, data: data}
}
func (n *BSTNode) String() string {
	if n == nil {
		return "<nil>"
	}
	return n.key.String()
}

func (n *BSTNode) Key() KeyType {
	if n == nil {
		return FakeNode
	}
	return n.key
}

func (n *BSTNode) Value() any {
	if n == nil {
		return nil
	}

	return n.data
}

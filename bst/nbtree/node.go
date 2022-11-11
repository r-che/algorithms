package nbtree

import "fmt"


// KeyType represents the key type of a binary search tree node
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

// BSTNode implements a binary search tree node
type BSTNode struct {
	key		KeyType
	left	*BSTNode
	right	*BSTNode
	parent	*BSTNode

	data any
}

// NewBSTNode creates a binary search tree node with key k  and associates the data with it
func NewBSTNode(k KeyType, data any) *BSTNode {
	return &BSTNode{key: k, data: data}
}
func (n *BSTNode) String() string {
	if n == nil {
		return "<nil>"
	}
	return n.key.String()
}

// Key returns the key value of the node
func (n *BSTNode) Key() KeyType {
	if n == nil {
		return FakeNode
	}
	return n.key
}

// Value returns the data associated with the node
func (n *BSTNode) Value() any {
	if n == nil {
		return nil
	}

	return n.data
}

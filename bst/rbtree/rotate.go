package rbtree

import "fmt"

//
// Single rotate values
//
type Rotate int
const (
	Left =	Rotate(iota + 1)
	Right
)

func (r Rotate) String() string {
	switch r {
		case Left:		return "Left"
		case Right:		return "Right"
	}

	panic(fmt.Sprintf("Unexpected rotate value: %d", r))
}

//
// Double rotate values
//
type RotateDouble Rotate
const (
	LeftRight	=	RotateDouble(Right + 1 + iota)
	RightLeft
)

func (r RotateDouble) String() string {
	switch r {
		case LeftRight:	return "LeftRight"
		case RightLeft:	return "RightLeft"
	}

	panic(fmt.Sprintf("Unexpected rotate value: %d", r))
}

//
// RBTree rotation operations
//

func (t *RBTree) rotateDouble(rType RotateDouble, pivot, node *RBNode) {
	// DBG-print: fmt.Printf("[ROTATE] %s - pivot: %s child: %s\n", rType, pivot, node)
	switch rType {
		case LeftRight:
			// Do left rotation using node as pivot
			nextNode := node.right
			t.rotate(Left, node, nextNode)

			// Do right rotation around pivot, using next node as left child node of pivot
			t.rotate(Right, pivot, nextNode)
		case RightLeft:
			// Do right rotation using node as pivot
			nextNode := node.left
			t.rotate(Right, node, nextNode)

			// Do left rotation around pivot, using next node as left child node of pivot
			t.rotate(Left, pivot, nextNode)
		default:
			panic(fmt.Sprintf("Unsupported rotation type: %d", rType))
	}
}

func (t *RBTree) rotate(rType Rotate, pivot, node *RBNode) {
	// Select rotate type
	switch rType {
		case Left:
			// Attach left child of node to right of pivot
			pivot.right = node.left
			if node.left != nil {
				node.left.parent = pivot
			}

			// Make pivot left child of the node
			node.left = pivot

		case Right:
			// Attach right child of node to left of pivot
			pivot.left = node.right
			if node.right != nil {
				node.right.parent = pivot
			}

			// Make pivot right child of the node
			node.right = pivot

		default:
			panic(`Unsupported rotation type "` + rType.String() + `" in rotate(), must be only Left or Right`)
	}

	// Update parents
	node.parent = pivot.parent
	if parent := pivot.parent; parent != nil {
		// Need to update pointer in the pivot's parent
		if parent.left == pivot {
			// Update left pointer
			parent.left = node
		} else {
			// Update right pointer
			parent.right = node
		}
	} else {
		// pivot parent == nil => pivot is the root of the tree,
		// need to update pointer to the tree root
		t.root = node
	}

	pivot.parent = node
}

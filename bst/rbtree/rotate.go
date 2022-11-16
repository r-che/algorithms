package rbtree

import "fmt"

//
// Single rotate values
//
type Rotate int
const (
	L =	Rotate(iota + 1)
	R
)

func (r Rotate) String() string {
	switch r {
		case L:		return "L"
		case R:		return "R"
	}

	panic(fmt.Sprintf("Unexpected rotate value: %d", r))
}

//
// Double rotate values
//
type RotateDouble Rotate
const (
	LR	=	RotateDouble(R + 1 + iota)
	RL
)

func (r RotateDouble) String() string {
	switch r {
		case LR:	return "LR"
		case RL:	return "RL"
	}

	panic(fmt.Sprintf("Unexpected rotate value: %d", r))
}

//
// RBTree rotation operations
//

func (t *RBTree) rotateDouble(rType RotateDouble, pivot, node *RBNode) {
	// DBG-print: fmt.Printf("[ROTATE] %s - pivot: %s child: %s\n", rType, pivot, node)
	switch rType {
		case LR:
			// Do left rotation using node as pivot
			nextNode := node.right
			t.rotate(L, node, nextNode)

			// Do right rotation around pivot, using next node as left child node of pivot
			t.rotate(R, pivot, nextNode)
		case RL:
			// Do right rotation using node as pivot
			nextNode := node.left
			t.rotate(R, node, nextNode)

			// Do left rotation around pivot, using next node as left child node of pivot
			t.rotate(L, pivot, nextNode)
		default:
			panic(fmt.Sprintf("Unsupported rotation type: %d", rType))
	}
}

func (t *RBTree) rotate(rType Rotate, pivot, node *RBNode) {
	// Select rotate type
	switch rType {
		case L:
			// Attach left child of node to right of pivot
			pivot.right = node.left
			if node.left != nil {
				node.left.parent = pivot
			}

			// Make pivot left child of the node
			node.left = pivot

		case R:
			// Attach right child of node to left of pivot
			pivot.left = node.right
			if node.right != nil {
				node.right.parent = pivot
			}

			// Make pivot right child of the node
			node.right = pivot

		default:
			panic(`Unsupported rotation type "` + rType.String() + `" in rotate(), must be only L or R`)
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

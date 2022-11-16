/*
Package nbtree provides an example of a non-balanced binary search tree
implementation.

It supports standard tree procedures, such as: inserting and deleting nodes,
finding nodes by given arbitrary key, finding the root, maximum and minimum
nodes, finding the predecessor and successor of a node.

It supports colored output of graphical representation of the tree using ASCII
graphics. For example, a tree with the keys 20, 10, 30, 5, 15, 25, 35, 8, 17,
37, 33, 13, 2, 23, 27 added sequentially will look like this:

                             20
                ____________/  \____________
               /                            \
             10                              30
        ____/  \____                    ____/  \____
       /            \                  /            \
     5               15              25              35
    /  \            /  \            /  \            /  \
   /    \          /    \          /    \          /    \
 2       8       13      17      23      27      33      37

*/
package nbtree

// BSTree implements a binary search tree.
type BSTree struct {
	root	*BSTNode
}

// NewBSTree returns new empty binary search tree.
func NewBSTree() *BSTree {
	return &BSTree{}
}

// Search returns a tree node with key k or nil if there is no such node.
func (t *BSTree) Search(k KeyType) *BSTNode {
	n := t.root
	for n != nil && n.key != k {
		if k < n.key {
			n = n.left
		} else {
			n = n.right
		}
	}

	return n
}

// Root returns the root node of the binary search tree, or nil if the tree is empty.
func (t *BSTree) Root() *BSTNode {
	return t.root
}

// Min returns the tree node with the minimum key value.
func (t *BSTree) Min() *BSTNode {
	if t.root == nil {
		return nil
	}
	n := t.root
	for n.left != nil {
		n = n.left
	}
	return n
}

// Max returns the tree node with the maximum key value.
func (t *BSTree) Max() *BSTNode {
	if t.root == nil {
		return nil
	}
	n := t.root
	for n.right != nil {
		n = n.right
	}
	return n
}

// Successor returns the tree node following node n in a linear ordering of
// tree nodes in ascending order of their keys. If there is none, i.e. n has a
// maximal key value, then nil is returned.
func (t *BSTree) Successor(n *BSTNode) *BSTNode {
	// If node has right sub-tree
	if n.right != nil {
		// Need to return minimum of the left sub-tree
		n = n.right
		for n.left != nil {
			n = n.left
		}
		return n
	}

	// Need to go up until find parent for which n is the LEFT child
	p := n.parent
	for p != nil && n == p.right {
		n = p
		p = p.parent
	}

	return p
}

// Predecessor returns the tree node following node n in a linear ordering of
// tree nodes in descending order of their keys. If there is none, i.e. n has a
// minimum key value, then nil is returned.
func (t *BSTree) Predecessor(n *BSTNode) *BSTNode {
	// If node has left sub-tree
	if n.left != nil {
		// Need to return maximum of the right sub-tree
		n = n.left
		for n.right != nil {
			n = n.right
		}
		return n
	}

	// Need to go up until find parent for which n is the RIGHT child
	p := n.parent
	for p != nil && n == p.left {
		n = p
		p = p.parent
	}

	return p
}

// SearchWithParent returns as the first value a node with k if found or nil if not found,
// as the second - parent of the found node even if the node was not found.
func (t *BSTree) SearchWithParent(k KeyType) (*BSTNode, *BSTNode) {
	n := t.root
	p := n.parent
	for n != nil && n.key != k {
		p = n
		if k < n.key {
			n = n.left
		} else {
			n = n.right
		}
	}

	return n, p
}

// Insert node n into the tree keeping the properties of the binary search tree.
func (t *BSTree) Insert(n *BSTNode) *BSTNode { //nolint:varnamelen // n is too obvious to make it longer
	// Check for empty tree
	if t.root == nil {
		// Make the node a root and return
		t.root = n
		return n
	}

	// Search node with key k in the tree
	N, p := t.SearchWithParent(n.key)
	if N != nil {
		// Already exists
		// fmt.Printf("Node %d already exists in the tree\n", n.key)
		return nil
	}

	// Assign correct parent of the new node
	n.parent = p

	// Select correct child pointer in the parent
	if n.key < p.key {
		// Assign new node as left child
		p.left = n
	} else {
		// Assign new node as right child
		p.right = n
	}

	return n
}

// Delete deletes the node n from the tree keeping the properties of the binary search tree.
func (t *BSTree) Delete(n *BSTNode) *BSTNode {
	// Choose type of deletion
	switch {
	// Node has TWO children
	case n.left != nil && n.right != nil:
		// Return successor from delChildren, as pointer to really removed node
		return t.delChildren(n)

	// Node is leaf - NO children
	case n.left == nil && n.right == nil:
		return t.delLeaf(n)

	// Node has one child
	default:
		return t.delChild(n)
	}
}

func (t *BSTree) delChildren(n *BSTNode) *BSTNode {
	// Get successor of nPtr - this node will repalce nPtr
	s := t.Successor(n)

	// Remove s from its position - s can point only to
	// leaf node or to node that has only one right-child
	t.Delete(s)

	// if n == t.root {
	// 	fmt.Printf("Tree root updated[2]: %v -> %v\n", n, s)
	// }

	// Now s extracted from tree, need to replace n by s, do inplace update
	n.key = s.key
	n.data = s.data

	// Return successor s as pointer to really removed node
	return s
}

func (t *BSTree) delLeaf(n *BSTNode) *BSTNode {
	// Check for n is root of the tree
	if n == t.root {
		// Cleanup root
		t.root = nil

		// Empty tree, nothing to fixup - return now
		return n
	}

	// Clear parent's pointer to n
	if n.parent.left == n {
		n.parent.left = nil
	} else {
		n.parent.right = nil
	}

	return n
}

func (t *BSTree) delChild(n *BSTNode) *BSTNode { //nolint:varnamelen // n is too obvious to make it longer
	// Get n's single child
	var child *BSTNode
	if n.left != nil {
		child = n.left
	} else {
		child = n.right
	}

	// Replace parent value of the child node
	child.parent = n.parent

	if n.parent == nil {
		// fmt.Printf("Tree root updated[1]: %v -> %v\n", n, child)
		// Replace root
		t.root = child

		// Return n as is
		return n
	}

	// Assign child as child of n's parent
	if n.parent.left == n {
		n.parent.left = child
	} else {
		n.parent.right = child
	}

	return n
}

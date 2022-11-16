package rbtree

// Search returns a tree node with key k or nil if there is no such node.
func (t *RBTree) Search(k KeyType) *RBNode {
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
func (t *RBTree) Root() *RBNode {
	return t.root
}

// Min returns the tree node with the minimum key value.
func (t *RBTree) Min() *RBNode {
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
func (t *RBTree) Max() *RBNode {
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
func (t *RBTree) Successor(n *RBNode) *RBNode {
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
func (t *RBTree) Predecessor(n *RBNode) *RBNode {
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
func (t *RBTree) SearchWithParent(k KeyType) (*RBNode, *RBNode) {
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

// bstInsert inserts node n into the tree keeping the properties of the binary search tree
func (t *RBTree) bstInsert(n *RBNode) (*RBNode, bool) { //nolint:varnamelen // n is too obvious to make it longer
	// Set color
	n.color = Red

	// Check for empty tree
	if t.root == nil {
		// Make the node a root of the tree
		t.root = n

		// Repaint root to black
		n.color = Black

		// Return root node and no fixup required
		return n, false
	}

	// Search node with key k in the tree
	N, p := t.SearchWithParent(n.key)
	if N != nil {
		// Already exists, no insertion or fixup required
		return nil, false
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

	// Return pointer to the inserted node, RB-tree fixup required
	return n, true
}

//nolint:cyclop // Not sure that code will be more clear if this function is split into several
func (t *RBTree) bstDelete(n *RBNode) *RBNode { //nolint:varnamelen // n is too obvious to make it longer
	// Choose type of deletion
	switch {
	// Node has TWO children
	case n.left != nil && n.right != nil:
		// Get successor of nPtr - this node will repalce nPtr
		s := t.Successor(n)

		// Remove s from its position - s can point only to
		// leaf node or to node that has only one right-child
		t.bstDelete(s)

		// if n == t.root {
		// 	fmt.Printf("Tree root updated[2]: %v -> %v\n", n, s)
		// }

		// Now s extracted from tree, need to replace n by s, do inplace update
		n.key = s.key
		n.data = s.data

		// Replace value of n to return s as pointer to really removed node
		n = s

	// Node is leaf - NO children
	case n.left == nil && n.right == nil:
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

	// Node has one child
	default:
		// Get n's single child
		var child *RBNode
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

			// Break to return n as is
			break
		}

		// Assign child as child of n's parent
		if n.parent.left == n {
			n.parent.left = child
		} else {
			n.parent.right = child
		}

		// n is deleted now
	}

	return n
}

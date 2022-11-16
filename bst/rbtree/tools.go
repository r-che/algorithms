package rbtree

import "fmt"

// SelfTest performs a self-test of the red-black tree and returns the black-height,
// and a description of the problem if detected. If an issuse is detected, the
// black-height is zero.
func (t *RBTree) SelfTest() (int, error) {
	if t.root.Color() != Black {
		return 0, fmt.Errorf("v#5: tree root (%v) is NOT black", t.root)
	}

	return t.root.test()
}

func (n *RBNode) test() (int, error) {
	// No errors on empty sub-tree
	if n == nil {
		return 0, nil
	}

	bhl, err := n.left.test()	// bhl - black height left
	if err != nil {
		return 0, err
	}

	bhr, err := n.right.test()
	if err != nil {
		return 0, err
	}

	// Test inequality of black heights of subtrees
	if bhl != bhr {
		return 0, fmt.Errorf(
			"v#4: node %v - black-height left (%d) is not equal black height-right (%d)",
			n, bhl, bhr)
	}

	// Test current node color
	if n.color == Black {
		bhl++
	} else
	// Red node, need to check children colors - both must be Red
	if n.left.Color() != Black || n.right.Color() != Black {
		return 0, fmt.Errorf(
			"v#3: Red node (%v) has non-Black child (left: %v, right: %v)",
			n, n.left, n.right)
	}

	// OK
	return bhl, nil
}

func swapColors(n1, n2 *RBNode) {
	n1.color, n2.color = n2.color, n1.color
}

// straightLine returns true if child c is added to parent f on the same side that f is added as child to g
func straightLine(c, f, g *RBNode) bool {
	if c == f.left && f == g.left ||
	   c == f.right && f == g.right {
		// Straight line c->f->g
		return true
	}

	// Not straight
	return false
}

func determineRelatedness(n *RBNode) (f, u, g *RBNode) {	//nolint:nonamedreturns
	f = n.parent	// father of n
	g = f.parent	// grandfather of n

	// Determine uncle of n
	if g.left == f {
		// Uncle is right child of "grandfather"
		u = g.right
	} else {
		// Uncle is left child of parent
		u = g.left
	}

	return f, u, g
}

func determChildOfDeleted(d *RBNode) (*RBNode, func()) { //nolint:varnamelen // name too obvious to make it longer
	// XXX Deleted node can have only 0 or 1 child
	if n := d.left; n != nil {
		// Return left child of deleted node
		return n, func(){}
	}

	// Left child does not exist, check for right
	if n := d.right; n != nil {
		// Return left child of deleted node
		return n, func(){}
	}

	// XXX Deleted node does not have children, return fake node
	fakeChild := &RBNode{
		parent:	d.parent,
		key:	FakeNode,
	}

	// Need to assign fakeChild to correct side of d.parent
	// XXX We can safely remove fake node from f (assign nil to child pointer)
	// XXX when returning from the function, because there is no combination that
	// XXX can replace the values of f pointers with some other node(s)
	if d.parent.left == nil {
		// d was removed from left side, assign fakeChild as left child
		d.parent.left = fakeChild

		// Return fakeChild with function to cleanup left child of deleted parent
		return fakeChild, func() { d.parent.left = nil }
	}

	// Otherwise - d was right child of its parent
	d.parent.right = fakeChild

	// Return fakeChild with function to cleanup right child of deleted parent
	return fakeChild, func() { d.parent.right = nil }
}

// determParticipants returns participants of fixup:
// f - father of node, b - brother of node, cn - nearside child of b, cf - far side child of b
func determParticipants(n *RBNode) (f, b, cn, cf *RBNode) {	//nolint:nonamedreturns
	f = n.parent	// "father"

	if f.left == n {
		// n left of f => New brother is right node of f
		b = f.right
		// Assign children to left case
		/* n, */ cn, cf = /* n, */ b.left, b.right
	} else {
		// n right of f => new brother is left node of f
		b = f.left
		// Assign children to right case
		cf, cn /*, n */ = b.left, b.right /*, n */
	}

	return
}

func determTurns(n, f *RBNode) (turnCase2or4, turnCase3 Rotate) {	//nolint:nonamedreturns
	// If n is a left child of f
	if f.left == n {
		return Left, Right
	}

	// n is right child of f
	return Right, Left
}

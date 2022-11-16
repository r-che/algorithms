package rbtree

import "fmt"

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
	if n1 == nil {
		n1 = &RBNode{}
	}
	if n2 == nil {
		n2 = &RBNode{}
	}
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

func determineRelatedness(n *RBNode) (f, u, g *RBNode) {
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

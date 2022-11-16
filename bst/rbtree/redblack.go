/*
Package rbtree provides an example of a Red-black search tree implementation.

It supports standard tree procedures, such as: inserting and deleting nodes,
finding nodes by given arbitrary key, finding the root, maximum and minimum
nodes, finding the predecessor and successor of a node.

It supports colored output of graphical representation of the tree using ASCII
graphics. For example, a tree with the keys 20, 10, 30, 5, 15, 25, 35, 8, 17,
37, 33, 13, 2, 23, 27 added sequentially will create tree [like this].

[like this]: https://pkg.go.dev/github.com/r-che/algorithms/bst/rbtree/rbtree.png

*/
package rbtree

import "fmt"

type RBTree struct {
	root	*RBNode
}
func NewRBTree() *RBTree {
	return &RBTree{}
}

// Delete deletes the node n from the tree keeping the properties of the Red-Black tree.
func (t *RBTree) Delete(n *RBNode) *RBNode {
	n = t.bstDelete(n)

	if t.root != nil {
		t.fixupDel(n)
	}

	return n
}

// Insert inserts node n into the tree keeping the properties of the Red-Black tree.
func (t *RBTree) Insert(n *RBNode) *RBNode {
	n, needFixup := t.bstInsert(n)

	if needFixup {
		// RB insertion fixup
		t.fixupIns(n)
	}

	return n
}

func (t *RBTree) fixupIns(n *RBNode) {	//nolint:varnamelen	// variable name too obvious to make it longer
	if n.parent.color == Black {
		// Nothing to fixup
		return
	}

	//
	// Red-violation - red node attached to red parent, fixup is required
	//

	//nolint:varnamelen	// variable names markings too obvious to make them longer
	for {
		// DBG-print: fmt.Printf("[FIXUP INS]: n: %v f: %v\n", n, n.parent)
		// Get n's relatedness
		f, u, g := determineRelatedness(n)

		// Do fixup operations
		switch {
		// Red uncle
		case u.Color() == Red:
			if contFixup := t.fixupRedUncle(f, u, g); !contFixup {
				// Stop fixup
				return
			}

			// Do fixup again, use g as new initiator of red-violation
			n = g

		// Black uncle and n->f->g is a straight line
		case u.Color() == Black && straightLine(n, f, g):
			// Fixup nodes
			t.fixupBlackUncleStraight(f, g)

			// No more fixups required
			return

		// Black uncle and n->f->g is angle (not a straight line)
		case u.Color() == Black && !straightLine(n, f, g):
			// Fixup nodes
			t.fixupBlackUncleAngle(n, f, g)

			// No more fixups required
			return

		default:
			panic(fmt.Sprintf("Unexpected state on nodes: n: %v f: %v g: %v u: %v", n, f, g, u))
		}
	}
}

// fixupRedUncle fixes tree when uncle color is red
func (t *RBTree) fixupRedUncle(f, u, g *RBNode) bool {
	// DBG-print: fmt.Printf("[U RED] n: %v u: %v\n", n, u)
	// Only a repaint is required
	f.color = Black
	u.color = Black

	// Is g root?
	if g == t.root {
		// Root always black, stop repainting and fixup
		return false
	}

	// Else - repaint g to red
	g.color = Red	// this may cause new red-violation

	// Return result of the check for new red-violation
	return g.parent.color == Red
}

// fixupBlackUncleStraight fixes tree when: black uncle and n->f->g is a straight line
func (t *RBTree) fixupBlackUncleStraight(f, g *RBNode) {
	// DBG-print: fmt.Printf("[U BLACK, STRAIGHT] n: %v f: %v g: %v\n", n, f, g)
	// Repaint nodes
	f.color = Black
	g.color = Red

	// Now, need rotate g around f

	// Define rotation direction and rotate
	if f == g.left {
		// f is a left child of g - need to rotate right
		t.rotate(Right, g, f)
	} else {
		// f is a right child of g - need to rotate left
		t.rotate(Left, g, f)
	}
}

// fixupBlackUncleAngle fixes tree when: black uncle and n->f->g is angle (not a straight line)
func (t *RBTree) fixupBlackUncleAngle(n, f, g *RBNode) {
	// DBG-print: fmt.Printf("[U BLACK, ANGLE] n: %v f: %v g: %v\n", n, f, g)
	// Repaint nodes
	g.color = Red
	n.color = Black

	// Now, double rotation is required

	// Define rotation direction and rotate
	if f == g.left {
		// f is a left child of g - need to rotate left+right
		t.rotateDouble(LeftRight, g, f)
	} else {
		// f is a right child of g - need to rotate right+left
		t.rotateDouble(RightLeft, g, f)
	}
}

func (t *RBTree) fixupDel(d *RBNode) {	//nolint:varnamelen	// variable name too obvious to make it longer
	//
	// Simple fixup cases
	//

	if d.color == Red {
		// Violations is not possible if deleted node is red
		// DBG-print: fmt.Printf("[DEL:RED] %v - no fixup required\n", d)
		return
	}

	//nolint:varnamelen	// n has too common meaning to make its name longer
	n, cleanFake := determChildOfDeleted(d)

	// Defer cleanup of fake node if it was created
	defer cleanFake()

	// If child of deleted node is Red - only repaint required
	if n.color == Red {
		// DBG-print: fmt.Printf("[CASE #0:REPAINT] N: %v -> %v\n", n, Black)
		// Repaint to Black and return
		n.color = Black

		return
	}

	//
	// More complex fixup
	//

	//nolint:varnamelen	// variable names markings too obvious to make them longer
	for n != nil {
		// Determine participants of fixup
		f, b, cn, cf := determParticipants(n)

		// Determine turns
		turnCase2or4, turnCase3 := determTurns(n, f)

		// DBG-print: fmt.Printf("[RELS] N:%v F:%v B:%v Cn:%v Cf:%v, turnCase2or4:%v, turnCase3:%v\n",
		// DBG-print:	n, f, b, cn, cf, turnCase2or4, turnCase3)

		switch {
			// 1. f is Red, others are Black
			case t.fixCase1(n, f, b, cn, cf):
				return

			// 2. b is Black, cf is Red
			case t.fixCase2(b, cf, f, turnCase2or4):
				return

			// 3. b is Black, cf is Black, cn is Red
			case t.fixCase3(b, cn, cf, turnCase3):
				// XXX Now situation brought to the case #2, required fixup will be done on the next iteration

			// 4. b is Red
			case t.fixCase4(f, b, turnCase2or4):
				// XXX Now situation brought to the cases #1..3, required fixup will be done on the next iteration

			// 5. All are Black
			case t.allBlack(n, f, b, cn, cf):
				// DBG-print: fmt.Println("[CASE #5]")
				// DBG-print: fmt.Printf("  [REPAINT] B:%v -> %v\n", b, Red)

				// Update n value by result of fixup
				n = t.fixCase5(f, b)

			default:
				panic(fmt.Sprintf("Unexpected state on nodes: D: %v N: %v F: %v B: %v Cn: %v Cf: %v",
					d, n, f, b, cn, cf))
		}
	}
}

func (t *RBTree) fixCase1(n, f, b, cn, cf *RBNode) bool {
	if !(n.color == Black &&
			f.color  == Red &&
			b.Color()  == Black &&
			cn.Color() == Black &&
			cf.Color() == Black) {
		// Another case
		return false
	}

	// DBG-print: fmt.Printf("[CASE #1:SWAP COLORS] F:%v <-> B:%v => ", f, b)

	// Swap colors between f and b
	swapColors(f, b)
	// DBG-print: fmt.Printf("f:%v, b:%v\n", f, b)

	// Fixed
	return true
}

func (t *RBTree) fixCase2(b, cf, f *RBNode, turn Rotate) bool {
	if !(b.Color() == Black && cf.Color() == Red) {
		// Another case
		return false
	}

	// DBG-print: fmt.Println("[CASE #2]")

	// DBG-print: fmt.Printf("  [ROTATE] F:%v around B:%v to %v\n", f, b, turnCase2or4)
	// 1. Rotate f around b
	t.rotate(turn, f, b)
	// DBG-print: fmt.Println(t)

	// DBG-print: fmt.Printf("  [REPAINT] Cf:%v -> %v\n", cf, Black)
	// 2. Repainting cf to Black
	cf.SetColor(Black)

	// DBG-print: fmt.Printf("  [SWAP COLORS] F:%v <-> B:%v => ", f, b)
	// 3. Swap colors between f and b
	swapColors(f, b)
	// DBG-print: fmt.Printf("F:%v, B:%v\n", f, b)

	// Fixed
	return true
}

func (t *RBTree) fixCase3(b, cn, cf *RBNode, turn Rotate) bool {
	if !(b.Color() == Black &&
		cn.Color() == Red &&
		cf.Color() == Black) {
		// Another case
		return false
	}

	// DBG-print: fmt.Println("[CASE #3]")

	// DBG-print: fmt.Printf("  [ROTATE] B:%v around Cn:%v to %v\n", b, cn, turnCase3)
	// Rotate b around cn
	t.rotate(turn, b, cn)
	// DBG-print: fmt.Println(t)

	// DBG-print: fmt.Printf("  [FLIP] B:%v -> %v, Cn:%v -> %v\n", b, !b.Color(), cn, !cn.Color())
	// Flip colors of b and cn
	b.Flip()
	cn.Flip()

	// Fixed
	return true
}

func (t *RBTree) fixCase4(f, b *RBNode, turn Rotate) bool {
	if !(b.Color() ==  Red) {
		// Another case
		return false
	}
	// DBG-print: fmt.Println("[CASE #4]")

	// DBG-print: fmt.Printf("  [ROTATE] F:%v around B:%v to %v\n", f, b, turnCase2or4)
	// Rotate f around b
	t.rotate(turn, f, b)
	// DBG-print: fmt.Println(t)

	// DBG-print: fmt.Printf("  [FLIP] F:%v -> %v, B:%v -> %v\n", f, !f.Color(), b, !b.Color())
	// Flip colors of f and b
	f.color = !f.color
	b.Flip()

	// Fixed
	return true
}

func (t *RBTree) allBlack(nodes ...*RBNode) bool {
	for _, n := range nodes {
		if n.Color() != Black {
			return false
		}
	}

	return true
}

func (t *RBTree) fixCase5(f, b *RBNode) *RBNode {
	b.SetColor(Red)

	// Check for f is tree root
	if f == t.root {
		// Fixup is done
		return nil
	}

	// DBG-print: fmt.Printf("  [UPDATE] N = F:%v (continue fixup) on:\n", f)
	// DBG-print: fmt.Println(t)

	// Return f to use it as next node to fixup
	return f
}

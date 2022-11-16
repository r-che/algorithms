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

func (t *RBTree) fixupIns(n *RBNode) {
	if n.parent.color == Black {
		// Nothing to fixup
		return
	}

	//
	// Red-violation - red node attached to red parent, fixup is required
	//

	// straightLine returns true if child c is added to parent f
	// on the same side that f is added as child to g
	straightLine := func(c, f, g *RBNode) bool {
		if c == f.left && f == g.left ||
		   c == f.right && f == g.right {
			// Straight line c->f->g
			return true
		}

		// Not straight
		return false
	}


	for {
		// DBG-print: fmt.Printf("[FIXUP INS]: n: %v f: %v\n", n, n.parent)
		// Get n's relations
		f := n.parent	// father of n
		g := f.parent	// grandfather of n
		var u *RBNode	// uncle of n
		if g.left == f {
			// Uncle is right child of "grandfather"
			u = g.right
		} else {
			// Uncle is left child of parent
			u = g.left
		}

		// Do fixup operations
		switch {
		// Red uncle
		case u.Color() == Red:
			// DBG-print: fmt.Printf("[U RED] n: %v u: %v\n", n, u)
			// Only a repaint is required
			f.color = Black
			u.color = Black
			// Is g root?
			if g == t.root {
				// Root always black, stop repainting
				return
			}

			// Else - repaint g to red
			g.color = Red	// this may cause new red-violation

			// Check for new red-violation
			if g.parent.color != Red {
				// No violation, stop fixup
				return
			}

			// Do fixup again, use g as new initiator of red-violation
			n = g

		// Black uncle and n->f->g is a straight line
		case u.Color() == Black && straightLine(n, f, g):
			// DBG-print: fmt.Printf("[U BLACK, STRAIGHT] n: %v f: %v g: %v\n", n, f, g)
			// Repaint nodes
			f.color = Black
			g.color = Red

			// Now, need rotate g around f

			// Define rotation direction and rotate
			if f == g.left {
				// f is a left child of g - need to rotate right
				t.rotate(R, g, f)
			} else {
				// f is a right child of g - need to rotate left
				t.rotate(L, g, f)
			}

			// No more fixups required
			return

		// Black uncle and n->f->g is angle (not a straight line)
		case u.Color() == Black && !straightLine(n, f, g):
			// DBG-print: fmt.Printf("[U BLACK, ANGLE] n: %v f: %v g: %v\n", n, f, g)
			// Repaint nodes
			g.color = Red
			n.color = Black

			// Now, double rotation is required

			// Define rotation direction and rotate
			if f == g.left {
				// f is a left child of g - need to rotate left+right
				t.rotateDouble(LR, g, f)
			} else {
				// f is a right child of g - need to rotate right+left
				t.rotateDouble(RL, g, f)
			}

			// No more fixups required
			return

		default:
			panic(fmt.Sprintf("Unexpected state on nodes: n: %v f: %v g: %v u: %v", n, f, g, u))
		}
	}
}

func (t *RBTree) fixupDel(d *RBNode) {
	//
	// Simple fixup cases
	//

	if d.color == Red {
		// Violations is not possible if deleted node is red
		// DBG-print: fmt.Printf("[DEL:RED] %v - no fixup required\n", d)
		return
	}

	// XXX Deleted node can have only 0 or 1 child
	var n *RBNode	// child of deleted node
	if n = d.left; n == nil {
		if n = d.right; n == nil {
			// XXX Eventually n is nil, assign to n fake node
			n = &RBNode{
				parent:	d.parent,
				key:	FakeNode,
			}

			// Need to assign n to correct side of d.parent
			// XXX We can safely remove fake node from f (assign nil to child pointer)
			// XXX when returning from the function, because there is no combination that
			// XXX can replace the values of f pointers with some other node(s)
			if d.parent.left == nil {
				// d was removed from left side, assign n as left child
				d.parent.left = n
				defer func() { d.parent.left = nil }()
			} else {
				// Otherwise - d was right child of its parent
				d.parent.right = n
				defer func() { d.parent.right = nil }()
			}
		}
	}

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

	var f *RBNode

	for {
		// Define "family" relationship
		f = n.parent	// "father"

		var b *RBNode		// new brother of child of deleted node
		// Turns' directions
		var turn24, turn3 Rotate
		// b children to use in conditions
		var cn *RBNode	// nearside child of b
		var cf *RBNode	// far side child of b

		if f.left == n {
			// n left of f => New brother is right node of f
			b = f.right
			// Define turns directions
			turn24 = L
			turn3 = R
			// Assign children to left case
			fmt.Printf("\n  === b %v\n\n", b)
			/* n, */ cn, cf = /* n, */ b.left, b.right
		} else {
			// n right of f => new brother is left node of f
			b = f.left
			// Define turns directions
			turn24 = R
			turn3 = L
			// Assign children to right case
			cf, cn /*, n */ = b.left, b.right /*, n */
		}

		// DBG-print: fmt.Printf("[RELS] N:%v F:%v B:%v Cn:%v Cf:%v, turn24:%v, turn3:%v\n",
		// DBG-print:	n, f, b, cn, cf, turn24, turn3)

		switch {
			// 1. f is Red, others are Black
			case n.color == Black &&
				f.color  == Red &&
				b.Color()  == Black &&
				cn.Color() == Black &&
				cf.Color() == Black:
				// DBG-print: fmt.Printf("[CASE #1:SWAP COLORS] F:%v <-> B:%v => ", f, b)
				// Swap colors between f and b
				swapColors(f, b)
				// DBG-print: fmt.Printf("f:%v, b:%v\n", f, b)

				// Stop fixup
				return

			// 2. b is Black, cf is Red
			case b.Color() == Black &&
				cf.Color() == Red:
				// Label to other cases can call this
				// DBG-print: fmt.Println("[CASE #2]")

				// DBG-print: fmt.Printf("  [ROTATE] F:%v around B:%v to %v\n", f, b, turn24)
				// 1. Rotate f around b
				t.rotate(turn24, f, b)
				// DBG-print: fmt.Println(t)

				// DBG-print: fmt.Printf("  [REPAINT] Cf:%v -> %v\n", cf, Black)
				// 2. Repainting cf to Black
				cf.SetColor(Black)

				// DBG-print: fmt.Printf("  [SWAP COLORS] F:%v <-> B:%v => ", f, b)
				// 3. Swap colors between f and b
				swapColors(f, b)
				// DBG-print: fmt.Printf("F:%v, B:%v\n", f, b)

				// Stop fixup
				return

			// 3. b is Black, cf is Black, cn is Red
			case b.Color() == Black &&
				cn.Color() == Red &&
				cf.Color() == Black:
				// DBG-print: fmt.Println("[CASE #3]")

				// DBG-print: fmt.Printf("  [ROTATE] B:%v around Cn:%v to %v\n", b, cn, turn3)
				// Rotate b around cn
				t.rotate(turn3, b, cn)
				// DBG-print: fmt.Println(t)

				// DBG-print: fmt.Printf("  [FLIP] B:%v -> %v, Cn:%v -> %v\n", b, !b.Color(), cn, !cn.Color())
				// Flip colors of b and cn
				b.Flip()
				cn.Flip()

				// Now situation brought to the case #2, required fixup
				// will be done on the next iteration

			// 4. b is Red
			case b.Color()  == Red:
				// DBG-print: fmt.Println("[CASE #4]")

				// DBG-print: fmt.Printf("  [ROTATE] F:%v around B:%v to %v\n", f, b, turn24)
				// Rotate f around b
				t.rotate(turn24, f, b)
				// DBG-print: fmt.Println(t)

				// DBG-print: fmt.Printf("  [FLIP] F:%v -> %v, B:%v -> %v\n", f, !f.Color(), b, !b.Color())
				// Flip colors of f and b
				f.color = !f.color
				b.Flip()

				// Now situation brought to the cases #1..3, required fixup
				// will be done on the next iteration

			// 5. All are Black
			case n.color == Black &&
				f.color  == Black &&
				b.Color()  == Black &&
				cn.Color() == Black &&
				cf.Color() == Black:

				// DBG-print: fmt.Println("[CASE #5]")
				// DBG-print: fmt.Printf("  [REPAINT] B:%v -> %v\n", b, Red)
				// Repaint b to Red
				b.SetColor(Red)

				if b == t.root {
					// Stop fixup process
					return
				}

				// Check for f is tree root
				if f == t.root {
					// Fixup is done
					return
				}

				// DBG-print: fmt.Printf("  [UPDATE] N = F:%v (continue fixup) on:\n", f)
				// DBG-print: fmt.Println(t)

				// Update n value by f and continue fixup
				n = f

			default:
				panic(fmt.Sprintf("Unexpected state on nodes: D: %v N: %v F: %v B: %v Cn: %v Cf: %v",
					d, n, f, b, cn, cf))
		}
	}
}

package rbtree

import (
	"fmt"
	"testing"
	"math/rand"
	"sort"
	"reflect"
)

const (
	keysCount	=	10240
	MaxItem		=	99999

	skipKeys	=	false
	makeKeys	=	true
)

//nolint:gochecknoglobals // We definitely do not want to
// run initialization for each test separately
var testKeys []KeyType
//nolint:gochecknoinits
func init() {
	// Use static seed for random source
	rand.Seed(2022)

	// Initiate keysCount unique keys...
	uniqs := make(map[KeyType]bool, keysCount)
	testKeys = make([]KeyType, 0, keysCount)
	for len(uniqs) < keysCount {
		n := KeyType(rand.Int() % (MaxItem + 1))	//nolint:gosec
		if _, ok := uniqs[n]; ok {
			// Already exists
			continue
		}

		// Append this item
		uniqs[n] = true
		testKeys = append(testKeys, n)
	}
}

func TestKeyType(t *testing.T) {
	for i, test := range []struct {
		kv		KeyType
		want	string
	} {
		{ FakeNode,  strFakeNode},
		{ -10, strInvalidNode },
		{ 1234, "1234" },
	} {
		if v := test.kv.String(); v != test.want {
			t.Errorf("[%d] KeyType.String() on %d, want - %q, got - %q", i, test.kv, test.want, v)
		}
	}
}

func TestNewRBNode(t *testing.T) {
	for testN, test := range []struct {
		n		*RBNode
		want	any
		wantKey	KeyType
		wantStr	string
	} { {
			NewRBNode(55, &struct{iv int; bv bool; is []int}{17, true, []int{9,8,7,6,5,4,3,2,1,0}}),
			&struct{iv int; bv bool; is []int}{17, true, []int{9,8,7,6,5,4,3,2,1,0}},
			55,
			Black.String() + "55",
		}, {
			NewRBNode(FakeNode, nil),
			nil,
			FakeNode,
			"<>",
		}, {
			nil,
			nil,
			FakeNode,
			Black.String() + "<nil>",
		},
	} {
		if !reflect.DeepEqual(test.n.Value(), test.want) {
			t.Errorf("[%d] got node - %#v, want - %#v", testN, test.n, test.want)
		}

		if sv := test.n.String(); sv != test.wantStr {
			t.Errorf("[%d] String() on node %#v returned %q, want - %q", testN, test.n, sv, test.wantStr)
		}

		if kv := test.n.Key(); kv != test.wantKey {
			t.Errorf("[%d] Key() on node %#v returned %v, want - %v", testN, test.n, kv, test.wantStr)
		}
	}
}

func TestEmpty(t *testing.T) {
	tree := NewRBTree()

	if n := tree.Root(); n != nil {
		t.Errorf("Root returned non-nil value %v (%#v) on empty tree", n, n)
	}

	if n := tree.Min(); n != nil {
		t.Errorf("Min returned non-nil value %v (%#v) on empty tree", n, n)
	}

	if n := tree.Max(); n != nil {
		t.Errorf("Max returned non-nil value %v (%#v) on empty tree", n, n)
	}
}

func TestInsert(t *testing.T) {
	tree := NewRBTree()

	// Insert all keys
	for i, v := range testKeys {
		// Create new node
		n := NewRBNode(v, nil)

		// Insert
		ins := tree.Insert(n)

		// Check for unsuccessful insertion
		if ins == nil {
			t.Errorf("[%d] RBTree.Insert returned nil, want - %p (inseted node: %v)", i, n, n)
			t.FailNow()
		}

		// Check for correctness of returned node
		if ins != n {
			t.Errorf("[%d] RBTree.Insert %p (%v), want - %p (inseted node: %v)", i, ins, ins, n, n)
			t.FailNow()
		}
	}
}

func TestSelfTest(t *testing.T) {
	tree, expHeight := newStaticTree()

	// Run self-testing
	bh, err := tree.SelfTest()

	// Check for tree structure problem
	if err != nil {
		t.Errorf("Red-Black tree structure issue: %v", err)
	}

	// Check for black height
	if bh != expHeight {
		t.Errorf("returned black-height - %d, want - %d", bh, expHeight)
	}
}

func TestSelfTestFail(t *testing.T) {
	for i, test := range treeBreakers() {
		tree, _ := newStaticTree()

		// Apply test to tree
		test(tree)

		// Run self-testing
		bh, err := tree.SelfTest()

		// Check for tree structure problem
		switch {
		// Check for no errors
		case err == nil:
			t.Errorf("[%d] self-test does not return expected issue", i)

		// Check for non-zero black height
		case bh != 0:
			t.Errorf("returned black-height of the invalid tree is not zero - %d", bh)

		// Ok, just print error for information
		default:
			t.Log("Expected self-test error:", err)
		}
	}
}

func treeBreakers() []func(t *RBTree) {
	return []func(t *RBTree) {
		// Repaint root to red
		func(t *RBTree) {
			t.root.color = Red
		},
		// Add red node to create red violation
		func(t *RBTree) {
			for n := t.Min(); n != nil; n = t.Successor(n) {
				if n.left == nil && n.right == nil && n.color == Red {
					n.left = NewRBNode(999, nil)
					n.left.color = Red
					return
				}
			}
			panic("No leaf red nodes were found")
		},
		// Add black node to create black-height violation
		func(t *RBTree) {
			for n := t.Max(); n != nil; n = t.Predecessor(n) {
				if n.left == nil && n.right == nil {
					n.left = NewRBNode(999, nil)
					n.left.color = Black
					return
				}
			}
			panic("No leaf nodes were found")
		},
	}
}

func TestInsertDupes(t *testing.T) {
	tree, _ := newTreeSortedKeys(testKeys, skipKeys)

	// Insert all keys
	for i, v := range testKeys {
		// Create new node
		n := NewRBNode(v, nil)

		// Insert
		ins := tree.Insert(n)

		// Check for successful insertion
		if ins != nil {
			t.Errorf("[%d] RBTree.Insert returned %p (%v), want - nil," +
				" because node with key %v should be already inserted", i, ins, ins, n.Key())
			t.FailNow()
		}
	}
}

func TestSearch(t *testing.T) {
	tree, _ := newTreeSortedKeys(testKeys, skipKeys)

	// Check keys for existence
	for _, k := range testKeys {
		n := tree.Search(k)
		if n == nil {
			t.Errorf("key %s was added but not found in the tree", k)
			t.FailNow()
		}
	}
}

func TestSuccessor(t *testing.T) {
	tree, sKeys := newTreeSortedKeys(testKeys, makeKeys)

	i := 0
	var prev *RBNode
	for s := tree.Min(); s != nil; s, i = tree.Successor(s), i+1 {
		// Check for overrun
		if i == len(sKeys) {
			t.Errorf("[%d] got successor #%d but only %d keys are available", i, i, len(sKeys))
			t.FailNow()
		}

		if s.Key() != sKeys[i] {
			t.Errorf("[%d] successor %v should have key %v, got - %v", i, prev, sKeys[i], s.Key())
			t.FailNow()
		}
		prev = s
	}

	// Check for size
	if i != len(sKeys) {
		t.Errorf("number of tested successor is %d, want - %d", i, len(sKeys))
	}
}

func TestPredecessor(t *testing.T) {
	tree, sKeys := newTreeSortedKeys(testKeys, makeKeys)

	i := 0
	var prev *RBNode
	for m, p := len(sKeys)-1, tree.Max(); p != nil; m, p, i = m-1, tree.Predecessor(p), i+1 {
		// Check for overrun
		if i == len(sKeys) {
			t.Errorf("[%d] got predecessor #%d but only %d keys are available", i, i, len(sKeys))
			t.FailNow()
		}

		if p.Key() != sKeys[m] {
			t.Errorf("[%d] predecessor %v should have key %v, got - %v", i, prev, sKeys[m], p.Key())
			t.FailNow()
		}
		prev = p
	}

	// Check for size
	if i != len(sKeys) {
		t.Errorf("number of tested predecessor is %d, want - %d", i, len(sKeys))
	}
}

func TestDelMax(t *testing.T) {
	tree, sKeys := newTreeSortedKeys(testKeys, makeKeys)

	for m, i, n := len(sKeys)-1, 0, tree.Max(); n != nil; n, i, m = tree.Max(), i+1, m-1 {
		// Check for overrun
		if i == len(sKeys) {
			t.Errorf("[%d] the RBTree.Max() returned %v as %d-th element but only %d keys are available", i, n, i+1, len(sKeys))
			t.FailNow()
		}

		if n.Key() != sKeys[m] {
			t.Errorf("[%d] maximal value in the tree with %d elements should be %v, got - %v", i, m+1, sKeys[m], n)
			t.FailNow()
		}

		// Delete
		if err := delWithChecks(tree, n); err != nil {
			t.Errorf("[%d] %v", i, err)
		}
	}
}

func TestDelMin(t *testing.T) {
	tree, sKeys := newTreeSortedKeys(testKeys, makeKeys)

	for i, n := 0, tree.Min(); n != nil; n, i = tree.Min(), i+1 {
		// Check for overrun
		if i == len(sKeys) {
			t.Errorf("[%d] the RBTree.Min() returned %v as %d-th element but only %d keys are available", i, n, i+1, len(sKeys))
			t.FailNow()
		}

		if n.Key() != sKeys[i] {
			t.Errorf("[%d] minimal value in the tree with %d elements should be %v, got - %v", i, i+1, sKeys[i], n)
			t.FailNow()
		}

		// Delete
		if err := delWithChecks(tree, n); err != nil {
			t.Errorf("[%d] %v", i, err)
		}
	}
}

func TestDelRoot(t *testing.T) {
	tree, _ := newTreeSortedKeys(testKeys, skipKeys)

	for i, n := 0, tree.Root(); n != nil; n, i = tree.Root(), i+1 {
		// Check for overrun
		if i == len(testKeys) {
			t.Errorf("[%d] the RBTree.Root() returned %v as %d-th element," +
				" but only %d keys are available", i, n, i+1, len(testKeys))
			t.FailNow()
		}

		// Delete
		if err := delWithChecks(tree, n); err != nil {
			t.Errorf("[%d] %v", i, err)
		}
	}
}

func TestDelRandom(t *testing.T) {
	tree, sKeys := newTreeSortedKeys(testKeys, makeKeys)

	for i := 0; len(sKeys) != 0; i++ {
		// Get the random element from the sKeys
		idx := rand.Int() % len(sKeys)	//nolint:gosec
		k := sKeys[idx]
		// Remove k from keys slice
		sKeys = append(sKeys[:idx], sKeys[idx+1:]...)

		// Search this key
		n := tree.Search(k)
		if n == nil {
			t.Errorf("the key %v was not found in the tree, but mast", k)
			t.FailNow()
		}

		if err := delWithChecks(tree, n); err != nil {
			t.Errorf("[%d] %v", i, err)
			t.FailNow()
		}
	}

	// Tree now must be empty
	if root := tree.Root(); root != nil {
		t.Errorf("tree must be empty (root == nil), but root is - %v", root)
	}
}

func newTreeSortedKeys(keys []KeyType, makeKeys bool) (*RBTree, []KeyType) {
	tree := NewRBTree()

	// Insert all keys
	for _, v := range keys {
		tree.Insert(NewRBNode(v, nil))
	}

	if !makeKeys {
		return tree, nil
	}

	// Make sorted copy of keys
	sKeys := make([]KeyType, len(keys))
	copy(sKeys, keys)
	sort.Slice(sKeys, func(i, j int) bool { return sKeys[i] < sKeys[j] } )

	return tree, sKeys
}

func newStaticTree() (*RBTree, int) {
	tree := NewRBTree()

	// RB-tree created from this keys MUST have black-height == 4
	const expectedHeight = 4
	keys := []KeyType{
		26, 13, 53, 93, 97, 57, 60, 65, 39, 44, 28, 17, 22, 2, 93, 25, 2, 24, 5, 25, 20, 73,
		4, 89, 27, 60, 48, 20, 62, 22, 92, 14, 52, 90, 36, 6, 50, 44, 68, 2, 89, 87, 64, 19,
		92, 82, 76, 49, 59, 64, 62, 19, 3, 71, 85, 69, 56, 59, 74, 44, 57, 56, 96, 94,
	}

	for _, k := range keys {
		tree.Insert(NewRBNode(k, nil))
	}

	return tree, expectedHeight
}

func delWithChecks(tree *RBTree, n *RBNode) error { //nolint:varnamelen // n is too obvious to make it longer
	if n.right != nil && n.left != nil {
		// Node has two children - successor of n should be returned as deleted node

		// Get successor BEFORE deletion
		s := tree.Successor(n)

		// Delete node n
		if del := tree.Delete(n); del != s {
			return fmt.Errorf("RBTree.Delete returned %v (%#v) as deleted node for %v (%#v), want - %v (%#v) - successor of %v",
				del, del, n, n, s, s, n)
		}

		// Ok
		return nil
	}

	//
	// Leaf or single node deletion - n should be returned
	//

	if del := tree.Delete(n); del != n {
		return fmt.Errorf("RBTree.Delete returned %v (%#v), want - %v (%#v)", del, del, n, n)
	}

	// OK
	return nil
}

func TestRotateString(t *testing.T) {
	for _, test := range []struct{
		r	Rotate
		s	string
	} {
		{ Left, "Left" },
		{ Right, "Right" },
	} {
		if s := test.r.String(); s != test.s {
			t.Errorf("%#v.String() returned %q, want - %q", test.r, s, test.s)
		}
	}
}

func TestRotateDoubleString(t *testing.T) {
	for _, test := range []struct{
		r	RotateDouble
		s	string
	} {
		{ LeftRight, "LeftRight" },
		{ RightLeft, "RightLeft" },
	} {
		if s := test.r.String(); s != test.s {
			t.Errorf("%#v.String() returned %q, want - %q", test.r, s, test.s)
		}
	}
}

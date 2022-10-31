package nbtree

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

func TestNewBSTNode(t *testing.T) {
	for testN, test := range []struct {
		n		*BSTNode
		want	any
		wantKey	KeyType
		wantStr	string
	} { {
			NewBSTNode(55, &struct{iv int; bv bool; is []int}{17, true, []int{9,8,7,6,5,4,3,2,1,0}}),
			&struct{iv int; bv bool; is []int}{17, true, []int{9,8,7,6,5,4,3,2,1,0}},
			55,
			"55",
		}, {
			NewBSTNode(FakeNode, nil),
			nil,
			FakeNode,
			"<>",
		}, {
			nil,
			nil,
			FakeNode,
			"<nil>",
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
	tree := NewBSTree()

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
	tree := NewBSTree()

	// Insert all keys
	for i, v := range testKeys {
		// Create new node
		n := NewBSTNode(v, nil)

		// Insert
		ins := tree.Insert(n)

		// Check for unsuccessful insertion
		if ins == nil {
			t.Errorf("[%d] BSTree.Insert returned nil, want - %p (inseted node: %v)", i, n, n)
			t.FailNow()
		}

		// Check for correctness of returned node
		if ins != n {
			t.Errorf("[%d] BSTree.Insert %p (%v), want - %p (inseted node: %v)", i, ins, ins, n, n)
			t.FailNow()
		}
	}
}

func TestInsertDupes(t *testing.T) {
	tree, _ := newTreeSortedKeys(testKeys, skipKeys)

	// Insert all keys
	for i, v := range testKeys {
		// Create new node
		n := NewBSTNode(v, nil)

		// Insert
		ins := tree.Insert(n)

		// Check for successful insertion
		if ins != nil {
			t.Errorf("[%d] BSTree.Insert returned %p (%v), want - nil," +
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
	var prev *BSTNode
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
	var prev *BSTNode
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
			t.Errorf("[%d] the BSTree.Max() returned %v as %d-th element but only %d keys are available", i, n, i+1, len(sKeys))
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
			t.Errorf("[%d] the BSTree.Min() returned %v as %d-th element but only %d keys are available", i, n, i+1, len(sKeys))
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
			t.Errorf("[%d] the BSTree.Root() returned %v as %d-th element," +
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

func newTreeSortedKeys(keys []KeyType, makeKeys bool) (*BSTree, []KeyType) {
	tree := NewBSTree()

	// Insert all keys
	for _, v := range keys {
		tree.Insert(NewBSTNode(v, nil))
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

func delWithChecks(tree *BSTree, n *BSTNode) error { //nolint:varnamelen // n is too obvious to make it longer
	if n.right != nil && n.left != nil {
		// Node has two children - successor of n should be returned as deleted node

		// Get successor BEFORE deletion
		s := tree.Successor(n)

		// Delete node n
		if del := tree.Delete(n); del != s {
			return fmt.Errorf("BSTree.Delete returned %v (%#v) as deleted node for %v (%#v), want - %v (%#v) - successor of %v",
				del, del, n, n, s, s, n)
		}

		// Ok
		return nil
	}

	//
	// Leaf or single node deletion - n should be returned
	//

	if del := tree.Delete(n); del != n {
		return fmt.Errorf("BSTree.Delete returned %v (%#v), want - %v (%#v)", del, del, n, n)
	}

	// OK
	return nil
}

func TestStringFilled(t *testing.T) {
	//nolint:lll // Predefined colored tree
	want := fmt.Sprintf(
`                          %[1]s20 %[2]s                                                                                           ` + `
                         /   \________________________________________                                                  ` + `
                        /                                             \                                                 ` + `
                     %[1]s10 %[2]s                                               %[1]s30 %[2]s                                              ` + `
                    /                             ____________________/   \____________________                         ` + `
                   /                             /                                             \                        ` + `
                %[1]s5  %[2]s                           %[1]s25 %[2]s                                               %[1]s35 %[2]s                     ` + `
     __________/                             /   \_____                                        /   \                    ` + `
    /                                       /          \                                      /     \                   ` + `
 %[1]s2  %[2]s                                     %[1]s23 %[2]s            %[1]s27 %[2]s                                %[1]s34 %[2]s       %[1]s37 %[2]s                ` + `
    \                              _____/              /   \                    __________/             \__________     ` + `
     \                            /                   /     \                  /                                   \    ` + `
      %[1]s3  %[2]s                      %[1]s21 %[2]s                 %[1]s26 %[2]s       %[1]s28 %[2]s            %[1]s31 %[2]s                                     %[1]s400%[2]s ` + `
         \                        \                             \              \                                   /    ` + `
          \                        \                             \              \                                 /     ` + `
           %[1]s4  %[2]s                      %[1]s22 %[2]s                           %[1]s29 %[2]s            %[1]s32 %[2]s                           %[1]s390%[2]s      ` + `
                                                                                    \                         /         ` + `
                                                                                     \                       /          ` + `
                                                                                      %[1]s33 %[2]s                 %[1]s38 %[2]s           ` + `
`,
	Color, Rst)

	// Make real tree
	tree := NewBSTree()
	for _, k := range []KeyType{
		20, 10, 30, 5, 25, 35, 37, 34, 2, 23, 27, 21, 31,
		3, 4, 28, 29, 400, 390, 38, 26, 22, 31, 32, 33,
	} {
		tree.Insert(NewBSTNode(k, nil))
	}

	// Compare
	if tStr := tree.String(); tStr != want {
		t.Errorf("BSTree.String() returned:\n---\n%s\n---\nWant:\n---\n%s\n---\n", tStr, want)
	}
}

func TestStringEmpty(t *testing.T) {
	tree := NewBSTree()
	if tStr := tree.String(); tStr != strEmptyTree {
		t.Errorf("BSTree.String() returned:\n---\n%s\n---\nWant:\n---\n%s\n---\n", tStr, strEmptyTree)
	}
}

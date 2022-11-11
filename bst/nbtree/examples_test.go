package nbtree

import "fmt"

//nolint:testableexamples
func Example_treeCreation() {
	// Create tree
	tree := NewBSTree()

	// Insert keys and data
	for _, k := range []KeyType{20, 10, 30, 5, 15, 25, 35, 8, 17, 37, 33, 13, 2, 23, 27} {
		tree.Insert(NewBSTNode(k, fmt.Sprintf("Value for key %v", k)))
	}

	// Print graphical representation of the tree
	fmt.Print(tree)
}

func Example_treeSearch() {
	// Tree creation
	tree := NewBSTree()
	for _, k := range []KeyType{20, 10, 30, 5, 15, 25, 35, 8, 17, 37, 33, 13, 2, 23, 27} {
		tree.Insert(NewBSTNode(k, fmt.Sprintf("Value for key %v", k)))
	}

	// Set of keys for search
	lookups := []KeyType{183, 30, 8, 92, 37, 99, 0, 15}

	for _, k := range lookups {
		if n := tree.Search(k); n == nil {
			fmt.Println("Key not found:", k)
		} else {
			fmt.Println("Found key", k, "value:", n.Value())
		}
	}

	// Output:
	// Key not found: 183
	// Found key 30 value: Value for key 30
	// Found key 8 value: Value for key 8
	// Key not found: 92
	// Found key 37 value: Value for key 37
	// Key not found: 99
	// Key not found: 0
	// Found key 15 value: Value for key 15
}

func Example_treeWalkingAscending() {
	// Tree creation
	tree := NewBSTree()
	for _, k := range []KeyType{20, 10, 30, 5, 15, 25, 35, 8, 17, 37, 33, 13, 2, 23, 27} {
		tree.Insert(NewBSTNode(k, fmt.Sprintf("Value for key %v", k)))
	}

	// Get the minimal node
	n := tree.Min()
	// Print it
	fmt.Print(n)

	// Walking through all nodes in ascending order using the Successor method
	for n = tree.Successor(n); n != nil; n = tree.Successor(n) {
		fmt.Print(" -> ", n)
	}

	fmt.Println()

	// Output:
	// 2 -> 5 -> 8 -> 10 -> 13 -> 15 -> 17 -> 20 -> 23 -> 25 -> 27 -> 30 -> 33 -> 35 -> 37
}

func Example_treeWalkingDescending() {
	// Tree creation
	tree := NewBSTree()
	for _, k := range []KeyType{20, 10, 30, 5, 15, 25, 35, 8, 17, 37, 33, 13, 2, 23, 27} {
		tree.Insert(NewBSTNode(k, fmt.Sprintf("Value for key %v", k)))
	}

	// Get the maximal node
	n := tree.Max()
	// Print it
	fmt.Print(n)

	// Walking through all nodes in descending order using the Predecessor method
	for n = tree.Predecessor(n); n != nil; n = tree.Predecessor(n) {
		fmt.Print(" <- ", n)
	}

	fmt.Println()

	// Output:
	// 37 <- 35 <- 33 <- 30 <- 27 <- 25 <- 23 <- 20 <- 17 <- 15 <- 13 <- 10 <- 8 <- 5 <- 2
}

//nolint:testableexamples
func Example_treeDelete() {
	// Tree creation
	tree := NewBSTree()
	keys := []KeyType{20, 10, 30, 5, 15, 25, 35}
	for _, k := range keys {
		tree.Insert(NewBSTNode(k, fmt.Sprintf("Value for key %v", k)))
	}

	fmt.Println("Created tree:\n", tree)

	for _, k := range keys {
		fmt.Println("Delete key:", k)
		tree.Delete(tree.Search(k))
		fmt.Println(tree)
	}
}

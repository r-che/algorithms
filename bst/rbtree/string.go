package rbtree

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	CircleRed	=	"\u2B24"
	CircleBlack	=	"\u25CB"
	TermRed		=	"\u001b[31m"
	TermBlack	=	"\u001b[37;1m"
	TermRst		=	"\u001b[0m"

	// Empty tree stub
	strEmptyTree	= `<tree-is-empty>`
)

func (ct ColorType) String() string {
	if ct == Red {
		return TermRed + CircleRed + TermRst
	}

	// Else - black
	return TermBlack + CircleBlack + TermRst
}

const circlePrintableWidth = 1

func (t *RBTree) String() string {
	if t.root == nil {
		return strEmptyTree
	}

	// Get a map with nodes separated by levels and
	// a map with positions of keys in a linear ordering of keys
	levels, positions := stringPrepareData(t)

	// Tree width
	width := len(positions)

	// Create output matrix
	const linesPerLevel = 3	// each output matrix level contains 3 lines, for:
							// * node keys
							// * initial slope of edge from node + horizontal part of edge
							// * final slanting part of the edge
	// So, the output matrix will have vertical dimension - tree height * linesPerLevel
	oMatrix := make([][]string, len(levels) * linesPerLevel)

	//
	// Calculate string representation sizes
	//

	// Maximal key width
	kw := t.root.maxKeyWidth()

	// Node ouptput format
	nFmt := "%s " + // XXX colored circle, skip length specifier because it has visible length lesser that length in bytes
		fmt.Sprintf("%%-%ds", kw) /* key */

	// Width of cell - colored circle pritable width + length between circle and key + key width
	nw := circlePrintableWidth + len(" ") + kw

	// Summary cell width that contains node
	cellWidth := len(" ") +  nw + len(" ")	// one space left + one space right of the node

	// Short stub that used to print cells that contain part of edges
	stub := strings.Repeat(" ", nw)

	// Fragment of a branch with one cell width
	branchFrag := strings.Repeat("_", cellWidth)

	//
	// Make an output matrix buffer
	//

	oLine := 0	// output matrix line
	level := 0	// levels matrix line
	for ; oLine < len(levels) * linesPerLevel; oLine, level = oLine + linesPerLevel, level+1 {
		// Fill output matrix level
		oMatrix[oLine] = make([]string, width)
		oMatrix[oLine+1] = make([]string, width)
		oMatrix[oLine+2] = make([]string, width)
		for _, node := range levels[level] {
			// Write node key to the output matrix
			oMatrix[oLine][positions[node.key]] = " " + fmt.Sprintf(nFmt, node.color, node.key) + " "

			// Write the initial fragment of the branch from the children to its parent
			stringInitBranchFrag(oMatrix[oLine+1], positions, node, stub)

			// Is it a root node?
			if node.parent == nil {
				// Root has no parents, no need to draw connections to them
				continue
			}

			// Determine direction of drawing
			var step int
			// Get the number of cells between parent and child
			if nc := positions[node.key] - positions[node.parent.key]; nc < 0 {
				// Node - LEFT child of its parent, need to draw branch to the right toward the parent
				oMatrix[oLine-1][positions[node.key]] = ` ` + stub + `/`
				step = 1
			} else {
				// Node - RIGHT child of its parent, need to draw branch to the left toward the parent
				oMatrix[oLine-1][positions[node.key]] = `\` + stub + ` `
				step = -1
			}

			for ni := positions[node.key] + step; ni != positions[node.parent.key]; ni += step {
				oMatrix[oLine-2][ni] = branchFrag
			}
		}
	}

	return stringMakeOutput(oMatrix, cellWidth)
}

// stringPrepareData source data to create string representation of the tree. It returns:
// levels -  map containing a set of levels (starting from the root - 0), each of that level
//           contains list of corresponding nodes in ascending order
// positions - map of key<=>position, when position is the position of corresponding key
//             in the flat ordered list of tree's keys
func stringPrepareData(t *RBTree) (map[int][]*RBNode, map[KeyType]int) {
	// Collect all nodes into the matrix
	levels := map[int][]*RBNode{0: []*RBNode{t.root}}
	t.root.childKeys(1, levels)

	// Map keys<=>position
	positions := map[KeyType]int{}
	for n, pos := t.Min(), 0; n != nil; n, pos = t.Successor(n), pos+1 {
		positions[n.key] = pos
	}

	return levels, positions
}

// stringInitBranchFrag writes the initial fragment of branches to chilldren, if any
func stringInitBranchFrag(row []string, positions map[KeyType]int, node *RBNode, stub string) {
	switch {
	case node.left != nil && node.right != nil:
		row[positions[node.key]] = `/` + stub + `\`
	case node.left != nil:
		row[positions[node.key]] = `/` + stub + ` `
	case node.right != nil:
		row[positions[node.key]] = ` ` + stub + `\`
	}
}

// stringMakeOutput converts matrix-representation of the tree to the multiline string value
func stringMakeOutput(matrix [][]string, cellWidth int) string {
	// Remove last two rows from matrix - it always empty
	matrix = matrix[:len(matrix)-2]

	// Make output buffer
	out := strings.Builder{}

	// Stub that used to print completely empty cells
	stubFull := strings.Repeat(" ", cellWidth)

	for _, level := range matrix {
		for _, n := range level {
			if n == "" {
				out.WriteString(stubFull)
			} else {
				out.WriteString(n)
			}
			// Append new line
		}
		out.WriteString("\n")
	}

	return out.String()
}

func (n *RBNode) childKeys(levelNum int, levels map[int][]*RBNode) {
	if n == nil || (n.left == nil && n.right == nil) {
		return
	}

	var children []*RBNode

	// Add children of the current node
	if n.left != nil {
		children = append(children, n.left)
	}

	if n.right != nil {
		children = append(children, n.right)
	}

	// Select level to appending children
	if _, ok := levels[levelNum]; ok {
		// Update existing level
		levels[levelNum] = append(levels[levelNum], children...)
	} else {
		// Need to assign new level
		levels[levelNum] = children
	}

	// Call recursively
	n.left.childKeys(levelNum+1, levels)
	n.right.childKeys(levelNum+1, levels)
}

func (n *RBNode) maxKeyWidth() int {
	if n == nil {
		return 0
	}

	max := utf8.RuneCountInString(n.key.String())

	if lmax := n.left.maxKeyWidth(); lmax > max {
		// XXX This code cannot be reached with integer keys because keys in the left subtree
		// XXX are always lesser than in the right, consequently their length cannot be greater
		// XXX than the length of the key in the current node
		max = lmax
	}

	if rmax := n.right.maxKeyWidth(); rmax > max {
		max = rmax
	}

	return max
}

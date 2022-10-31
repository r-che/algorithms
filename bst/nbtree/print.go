package nbtree

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	Color = "\u001b[92m"
	Rst = "\u001b[0m"

	strEmptyTree	= `<tree-is-empty>`
)

//nolint:cyclop // Not sure that code will be more clear if this function is split into several
func (t *BSTree) String() string {
	if t.root == nil {
		return strEmptyTree
	}

	// Collect all nodes into the matrix
	levels := map[int][]*BSTNode{0: []*BSTNode{t.root}}
	t.root.childKeys(1, levels)

	// Map keys<=>position
	positions := map[KeyType]int{}
	for n, pos := t.Min(), 0; n != nil; n, pos = t.Successor(n), pos+1 {
		positions[n.key] = pos
	}

	// Tree width
	width := len(positions)

	// Create output matrix
	const linesPerLevel = 3	// each output matrix level contains 3 lines, for:
							// * node keys
							// * initial slope of edge from node + horizontal part of edge
							// * final slanting part of the edge
	// So, the output matrix will have vertical dimension - tree height * linesPerLevel
	oMatrix := make([][]string, len(levels) * linesPerLevel)

	// Maximal key width
	kw := t.root.maxKeyWidth()
	// Node ouptput format
	nFmt := fmt.Sprintf(Color + "%%-%ds" + Rst, kw) // key format
	// Summary cell width that contains node
	cellWidth := len(" ") +  kw + len(" ")	// one space left + one space right of the key value
	// Short stub that used to print cells that contain part of edges
	stub := strings.Repeat(" ", kw)
	// Long stub that used to print completely empty cells
	stubFull := strings.Repeat(" ", cellWidth)

	oLine := 0	// output matrix line
	level := 0	// levels matrix line
	for ; oLine < len(levels) * linesPerLevel; oLine, level = oLine + linesPerLevel, level+1 {
		// Fill output matrix level
		oMatrix[oLine] = make([]string, width)
		oMatrix[oLine+1] = make([]string, width)
		oMatrix[oLine+2] = make([]string, width)
		for _, node := range levels[level] {
			oMatrix[oLine][positions[node.key]] = " " + fmt.Sprintf(nFmt, node.key) + " "
			switch {
			case node.left != nil && node.right != nil:
				oMatrix[oLine+1][positions[node.key]] = `/` + stub + `\`
			case node.left != nil:
				oMatrix[oLine+1][positions[node.key]] = `/` + stub + ` `
			case node.right != nil:
				oMatrix[oLine+1][positions[node.key]] = ` ` + stub + `\`
			}

			if node.parent == nil {
				continue
			}

			// Number of cells between parent and child
			nc := positions[node.key] - positions[node.parent.key]
			switch {
			// Node - LEFT child of its parent
			case nc < 0:
				oMatrix[oLine-1][positions[node.key]] = ` ` + stub + `/`
				for ni := positions[node.key] + 1; ni < positions[node.parent.key]; ni++ {
					oMatrix[oLine-2][ni] = strings.Repeat("_", cellWidth)
				}
			// Node - RIGHT child of its parent
			case nc > 0:
				oMatrix[oLine-1][positions[node.key]] = `\` + stub + ` `
				for ni := positions[node.key] - 1; ni > positions[node.parent.key]; ni-- {
					oMatrix[oLine-2][ni] = strings.Repeat("_", cellWidth)
				}
			}
		}
	}

	// Remove last two rows from oMatrix - it always empty
	oMatrix = oMatrix[:len(oMatrix)-2]

	// Make output buffer
	out := strings.Builder{}

	for _, level := range oMatrix {
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

func (n *BSTNode) childKeys(levelNum int, levels map[int][]*BSTNode) {
	if n == nil || (n.left == nil && n.right == nil) {
		return
	}

	var children []*BSTNode

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

func (n *BSTNode) maxKeyWidth() int {
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

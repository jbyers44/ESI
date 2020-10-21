package trie

import (
	"bytes"
	"fmt"
	"sort"
)

// Node represents a non-terminal node that points to either other nodes or a leaf
type Node struct {
	hash       []byte
	leftLabel  []byte
	left       interface{}
	rightLabel []byte
	right      interface{}
}

// This was created purely for testing whether or not a corrupted block would lead to an invalid result
// when validating the chain
// Should delete when all completed

// SetHash sets the hash of a node to a specific value
func (node *Node) SetHash(hash []byte) {
	node.hash = hash
}

func (node *Node) String() string {
	m := make(map[int]string)
	nodeStrings(node, m)

	nodeIds := make([]int, len(m))
	i := 0
	for id := range m {
		nodeIds[i] = id
		i++
	}
	sort.Ints(nodeIds)

	var b bytes.Buffer
	for i, id := range nodeIds {
		fmt.Fprintf(&b, "%s", m[id])
		if i < len(nodeIds)-1 {
			fmt.Fprintf(&b, "\n\n")
		}
	}
	return b.String()
}

func nodeStrings(node interface{}, nodeMap map[int]string) {
	switch n := node.(type) {
	case *Leaf:
		var b bytes.Buffer
		fmt.Fprintf(&b, "%d\n", len(nodeMap)+1)
		fmt.Fprintf(&b, "%s\n", n.value)
		fmt.Fprintf(&b, "%x\n", n.hash)
		nodeMap[len(nodeMap)+1] = b.String()
		return

	case *Node:
		var b bytes.Buffer
		fmt.Fprintf(&b, "%d\n", len(nodeMap)+1)
		fmt.Fprintf(&b, "%d\n", len(nodeMap)+2)
		fmt.Fprintf(&b, "%s\n", n.leftLabel)
		fmt.Fprintf(&b, "%x\n", n.hash)
		fmt.Fprintf(&b, "%s\n", n.rightLabel)
		fmt.Fprintf(&b, "%d\n", len(nodeMap)+3)
		nodeMap[len(nodeMap)+1] = b.String()

		nodeStrings(n.left, nodeMap)
		nodeStrings(n.right, nodeMap)
		return
	}
}

// GetHash gets the root
func (node *Node) GetHash() []byte {
	return node.hash

}
